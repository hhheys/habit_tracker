package habit

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"habit-tracker/internal/domain"
	habituc "habit-tracker/internal/usecase/habit"
	"strings"

	"github.com/lib/pq"
	"go.uber.org/zap"
)

type Repository struct {
	db  *sql.DB
	log *zap.Logger
}

func NewRepository(db *sql.DB, log *zap.Logger) *Repository {
	return &Repository{db: db, log: log}
}

func (r *Repository) List(ctx context.Context, params habituc.ListHabitsParams) ([]*domain.Habit, int64, error) {
	search := "%" + params.Search + "%"
	var total int64
	err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM habit h
		WHERE h.title ILIKE $1 OR h.description ILIKE $1`, search,
	).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	orderBy := listOrder(params.SortBy, params.SortOrder)
	query := fmt.Sprintf(`
		SELECT h.id, h.title, h.description, h.created_at, h.image_filename,
		       EXISTS (SELECT 1 FROM user_habit uh WHERE uh.habit_id = h.id AND uh.user_id = $2)
		FROM habit h
		WHERE h.title ILIKE $1 OR h.description ILIKE $1
		ORDER BY %s
		LIMIT $3 OFFSET $4`, orderBy)
	rows, err := r.db.QueryContext(ctx, query, search, params.UserID, params.Limit, params.Offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	habits := make([]*domain.Habit, 0)
	for rows.Next() {
		var h domain.Habit
		if err := rows.Scan(&h.ID, &h.Title, &h.Description, &h.CreatedAt, &h.ImageFilename, &h.IsAdded); err != nil {
			return nil, 0, err
		}
		h.Tags, err = r.tagsByHabitID(ctx, h.ID)
		if err != nil {
			return nil, 0, err
		}
		habits = append(habits, &h)
	}
	return habits, total, rows.Err()
}

func (r *Repository) GetByID(ctx context.Context, id uint) (*domain.Habit, error) {
	var h domain.Habit
	err := r.db.QueryRowContext(ctx, `
		SELECT id, title, description, created_at, image_filename
		FROM habit WHERE id = $1`, id,
	).Scan(&h.ID, &h.Title, &h.Description, &h.CreatedAt, &h.ImageFilename)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrHabitNotFound
	}
	if err != nil {
		return nil, err
	}
	h.Tags, err = r.tagsByHabitID(ctx, h.ID)
	if err != nil {
		return nil, err
	}
	return &h, nil
}

func (r *Repository) Create(ctx context.Context, h *domain.Habit) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = tx.QueryRowContext(ctx, `
		INSERT INTO habit (title, description, image_filename)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`,
		h.Title, h.Description, h.ImageFilename,
	).Scan(&h.ID, &h.CreatedAt)
	if isUniqueViolation(err) {
		return domain.ErrHabitAlreadyExists
	}
	if err != nil {
		return err
	}
	if err := insertTags(ctx, tx, h.ID, tagIDs(h.Tags)); err != nil {
		return err
	}
	return tx.Commit()
}

func (r *Repository) Update(ctx context.Context, input *habituc.UpdateHabitInput) (*domain.Habit, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, `
		UPDATE habit
		SET title = $1, description = $2,
		    image_filename = COALESCE(NULLIF($3, ''), image_filename)
		WHERE id = $4`,
		input.Title, input.Description, input.ImageFilename, input.ID,
	)
	if isUniqueViolation(err) {
		return nil, domain.ErrHabitAlreadyExists
	}
	if err != nil {
		return nil, err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, domain.ErrHabitNotFound
	}
	if len(input.RemoveTagIDs) > 0 {
		_, err = tx.ExecContext(ctx, `
			DELETE FROM habit_tag WHERE habit_id = $1 AND tag_id = ANY($2)`,
			input.ID, pq.Array(pgIDs(input.RemoveTagIDs)),
		)
		if err != nil {
			return nil, err
		}
	}
	if err := insertTags(ctx, tx, input.ID, input.AddTagIDs); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return r.GetByID(ctx, input.ID)
}

func (r *Repository) DeleteByID(ctx context.Context, id uint) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM habit WHERE id = $1`, id)
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return domain.ErrHabitNotFound
	}
	return nil
}

func (r *Repository) tagsByHabitID(ctx context.Context, habitID uint) ([]*domain.HabitTag, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT t.id, t.title
		FROM tag t JOIN habit_tag ht ON ht.tag_id = t.id
		WHERE ht.habit_id = $1
		ORDER BY t.id`, habitID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tags := make([]*domain.HabitTag, 0)
	for rows.Next() {
		var tag domain.HabitTag
		if err := rows.Scan(&tag.ID, &tag.Name); err != nil {
			return nil, err
		}
		tags = append(tags, &tag)
	}
	return tags, rows.Err()
}

func insertTags(ctx context.Context, tx *sql.Tx, habitID uint, ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	_, err := tx.ExecContext(ctx, `
		INSERT INTO habit_tag (habit_id, tag_id)
		SELECT $1, unnest($2::int[])
		ON CONFLICT (habit_id, tag_id) DO NOTHING`, habitID, pq.Array(pgIDs(ids)))
	return err
}

func tagIDs(tags []*domain.HabitTag) []uint {
	ids := make([]uint, 0, len(tags))
	for _, tag := range tags {
		if tag != nil {
			ids = append(ids, tag.ID)
		}
	}
	return ids
}

func pgIDs(ids []uint) []int64 {
	result := make([]int64, len(ids))
	for i, id := range ids {
		result[i] = int64(id)
	}
	return result
}

func listOrder(sortBy, sortOrder string) string {
	column := map[string]string{
		"title":      "h.title",
		"created_at": "h.created_at",
		"id":         "h.id",
	}[strings.ToLower(sortBy)]
	if column == "" {
		column = "h.created_at"
	}
	direction := "DESC"
	if strings.EqualFold(sortOrder, "asc") {
		direction = "ASC"
	}
	return column + " " + direction
}

func isUniqueViolation(err error) bool {
	var pqErr *pq.Error
	return errors.As(err, &pqErr) && pqErr.Code == "23505"
}
