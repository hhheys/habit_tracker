package userhabit

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"habit-tracker/internal/adapter/postgres/txmanager"
	"habit-tracker/internal/domain"
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

func (r *Repository) ListUserHabits(ctx context.Context, filter domain.UserHabitListFilter) ([]*domain.UserHabit, int64, error) {
	search := "%" + filter.Search + "%"
	var total int64
	err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM user_habit uh JOIN habit h ON h.id = uh.habit_id
		WHERE uh.user_id = $1 AND (h.title ILIKE $2 OR h.description ILIKE $2)`,
		filter.UserID, search,
	).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(`
		SELECT uh.id, uh.user_id, uh.habit_id, uh.added_at,
		       h.id, h.title, h.description, h.created_at, h.image_filename
		FROM user_habit uh JOIN habit h ON h.id = uh.habit_id
		WHERE uh.user_id = $1 AND (h.title ILIKE $2 OR h.description ILIKE $2)
		ORDER BY %s
		LIMIT $3 OFFSET $4`, listOrder(filter.SortBy, filter.SortOrder))
	rows, err := r.db.QueryContext(ctx, query, filter.UserID, search, filter.Limit, filter.Offset)
	if err != nil {
		return nil, 0, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			r.log.Error("failed to close rows", zap.Error(err))
		}
	}(rows)

	habits := make([]*domain.UserHabit, 0)
	for rows.Next() {
		h := domain.UserHabit{Habit: &domain.Habit{IsAdded: true}}
		if err := rows.Scan(
			&h.ID, &h.UserID, &h.HabitID, &h.AddedAt,
			&h.Habit.ID, &h.Habit.Title, &h.Habit.Description, &h.Habit.CreatedAt, &h.Habit.ImageFilename,
		); err != nil {
			return nil, 0, err
		}
		habits = append(habits, &h)
	}
	return habits, total, rows.Err()
}

func (r *Repository) CreateUserHabit(ctx context.Context, h *domain.UserHabit) error {
	executorType := txmanager.ExecutorFromContext(ctx, r.db)
	err := executorType.QueryRowContext(ctx, `
		INSERT INTO user_habit (user_id, habit_id)
		VALUES ($1, $2)
		RETURNING id, added_at`, h.UserID, h.HabitID,
	).Scan(&h.ID, &h.AddedAt)
	if isUniqueViolation(err) {
		return domain.ErrHabitAlreadyAdded
	}
	var pqErr *pq.Error
	if errors.As(err, &pqErr) && pqErr.Code == "23503" {
		return domain.ErrHabitNotFound
	}
	return err
}

func (r *Repository) GetTotalUserHabits(ctx context.Context, userID uint) (int, error) {
	executorType := txmanager.ExecutorFromContext(ctx, r.db)

	var total int
	err := executorType.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM user_habit
		WHERE user_id = $1`,
		userID,
	).Scan(&total)

	return total, err
}

func listOrder(sortBy, sortOrder string) string {
	column := map[string]string{
		"title":      "h.title",
		"added_at":   "uh.added_at",
		"created_at": "h.created_at",
	}[strings.ToLower(sortBy)]
	if column == "" {
		column = "uh.added_at"
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
