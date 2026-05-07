package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"habit-tracker/internal/dto/request"
	appErrors "habit-tracker/internal/errors"
	"habit-tracker/internal/models"

	"github.com/lib/pq"
	"go.uber.org/zap"
)

//go:generate mockery --name=HabitRepository --output=../../../mocks --outpkg=mocks
type HabitRepository interface {
	CreateHabit(habit *models.Habit) error
	GetHabitByID(id uint) (*models.Habit, error)
	UpdateHabit(habit *models.Habit) error
	DeleteHabit(id uint) error
	GetAllHabits(req *request.GetAllHabitsRequest, requestUserID uint) (int64, []*models.Habit, error)

	GetAllUserHabits(userID uint, query *request.GetUserHabitsRequest) ([]*models.UserHabit, error)
	GetUserHabit(userID, habitID uint) (*models.UserHabit, error)
	AddHabit(userID, habitID uint) (*models.UserHabit, error)
	RemoveHabit(userID, habitID uint) error

	CreateTag(habit *models.HabitTag) error
	EditTag(habit *models.HabitTag) error
	DeleteTag(habitID uint) error
	GetTagByID(habitID uint) (*models.HabitTag, error)
	GetAllTags(req request.GetAllHabitTagsRequest) ([]*models.HabitTag, int64, error)
	GetTagsByHabitID(habitID uint) ([]*models.HabitTag, error)

	ReplaceHabitTags(habitID uint, tagIDs []int) error
	AddHabitTags(habitID uint, habitTagIDs []int) error
}

type habitRepositoryImpl struct {
	DB     *sql.DB
	logger *zap.Logger

	StreakRepository
}

func NewHabitRepository(db *sql.DB, logger *zap.Logger, streakRep StreakRepository) HabitRepository {
	return &habitRepositoryImpl{
		DB:               db,
		logger:           logger,
		StreakRepository: streakRep,
	}
}

func (r *habitRepositoryImpl) CreateHabit(habit *models.Habit) error {
	err := r.DB.QueryRow(
		`
INSERT INTO habit(title, description, image_filename) VALUES ($1, $2, $3) RETURNING id, created_at`,
		habit.Title,
		habit.Description,
		habit.ImageFilename,
	).Scan(
		&habit.ID,
		&habit.CreatedAt,
	)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" { // unique_violation
				return appErrors.ErrHabitAlreadyExists
			}
		}
		return err
	}
	return nil
}

func (r *habitRepositoryImpl) GetHabitByID(id uint) (*models.Habit, error) {
	var habit models.Habit
	err := r.DB.QueryRow(
		`SELECT id, created_at, title, description, image_filename FROM habit WHERE id = $1`, id).Scan(
		&habit.ID,
		&habit.CreatedAt,
		&habit.Title,
		&habit.Description,
		&habit.ImageFilename,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, appErrors.ErrHabitNotFound
		}
		return nil, err
	}
	tags, err := r.GetTagsByHabitID(habit.ID)
	if err != nil {
		return nil, err
	}

	habit.Tags = tags
	return &habit, nil
}

func (r *habitRepositoryImpl) UpdateHabit(habit *models.Habit) error {
	err := r.DB.QueryRow(
		`UPDATE habit SET title = $1, description = $2, image_filename = $3 WHERE id = $4 RETURNING id, created_at`,
		habit.Title,
		habit.Description,
		habit.ImageFilename,
		habit.ID,
	).Scan(
		&habit.ID,
		&habit.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return appErrors.ErrHabitNotFound
		}
		return err
	}
	return nil
}

func (r *habitRepositoryImpl) DeleteHabit(id uint) error {
	_, err := r.DB.Exec(`DELETE FROM habit WHERE id = $1`, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *habitRepositoryImpl) GetAllHabits(
	req *request.GetAllHabitsRequest,
	requestUserID uint,
) (int64, []*models.Habit, error) {

	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 10
	}

	search := "%"
	if req.Search != "" {
		search = "%" + req.Search + "%"
	}

	offset := (req.Page - 1) * req.PageSize

	// =========================
	// BASE WHERE
	// =========================
	where := `
		(h.title ILIKE $1 OR h.description ILIKE $1)
	`

	args := []any{search}
	argIndex := 2
	// =========================
	// TAG FILTER
	// =========================
	if len(req.TagIDs) > 0 {
		where += fmt.Sprintf(`
			AND EXISTS (
				SELECT 1
				FROM habit_tag ht
				WHERE ht.habit_id = h.id
				AND ht.tag_id = ANY($%d)
			)
		`, argIndex)

		args = append(args, pq.Array(req.TagIDs))
		argIndex++
	}

	// =========================
	// COUNT QUERY
	// =========================
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM habit h
		WHERE %s
	`, where)

	var total int64
	err := r.DB.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		r.logger.Error("Error counting habits", zap.Error(err))
		return 0, nil, err
	}

	// =========================
	// SORT
	// =========================
	sort := req.Sort

	// =========================
	// MAIN QUERY
	// =========================
	query := fmt.Sprintf(`
		SELECT 
			h.id,
			h.created_at,
			h.title,
			h.description,
			h.image_filename,
			CASE 
				WHEN uh.habit_id IS NOT NULL THEN true
				ELSE false
			END AS is_added
		FROM habit h
		LEFT JOIN user_habit uh 
			ON uh.habit_id = h.id 
			AND uh.user_id = $%d
		WHERE %s
		ORDER BY %s
		LIMIT $%d OFFSET $%d;
	`, argIndex, where, sort, argIndex+1, argIndex+2)

	// добавляем userID + pagination args
	args = append(args, requestUserID, req.PageSize, offset)

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		r.logger.Error("Error querying habits", zap.Error(err))
		return 0, nil, err
	}
	defer rows.Close()

	habits := make([]*models.Habit, 0)

	for rows.Next() {
		h := new(models.Habit)

		if err := rows.Scan(
			&h.ID,
			&h.CreatedAt,
			&h.Title,
			&h.Description,
			&h.ImageFilename,
			&h.IsAdded,
		); err != nil {
			r.logger.Error("Error scanning habit row", zap.Error(err))
			return 0, nil, err
		}

		tags, err := r.GetTagsByHabitID(h.ID)
		if err != nil {
			return 0, nil, err
		}

		h.Tags = tags
		habits = append(habits, h)
	}

	return total, habits, nil
}

func (r *habitRepositoryImpl) GetAllUserHabits(userID uint, query *request.GetUserHabitsRequest) ([]*models.UserHabit, error) {
	rows, err := r.DB.Query(`SELECT uh.id, h.id, h.title, h.description, h.created_at, h.image_filename, added_at FROM habit h
																						   JOIN user_habit uh ON h.id = uh.habit_id
	WHERE user_id = $1 ORDER BY $2 `, userID, query.Sort)
	if err != nil {
		r.logger.Error("Error querying habits", zap.Error(err))
		return nil, err
	}
	defer rows.Close()
	habits := make([]*models.UserHabit, 0)
	for rows.Next() {
		userHabit := new(models.UserHabit)
		scanErr := rows.Scan(
			&userHabit.ID,
			&userHabit.Habit.ID,
			&userHabit.Habit.Title,
			&userHabit.Habit.Description,
			&userHabit.Habit.CreatedAt,
			&userHabit.Habit.ImageFilename,
			&userHabit.AddedAt,
		)
		if scanErr != nil {
			if !errors.Is(scanErr, sql.ErrNoRows) {
				r.logger.Error("Error scanning habits", zap.Error(scanErr))
				return nil, err
			}
		}

		streak, err := r.GetStreak(userHabit.ID)
		if err != nil {
			return nil, err
		}

		userHabit.Streak = streak
		habits = append(habits, userHabit)
	}
	return habits, nil
}

func (r *habitRepositoryImpl) AddHabit(userID, habitID uint) (*models.UserHabit, error) {
	habit, err := r.GetHabitByID(habitID)
	if err != nil {
		return nil, err
	}
	var userHabit models.UserHabit
	err = r.DB.QueryRow(`INSERT INTO user_habit(user_id, habit_id) VALUES ($1, $2) RETURNING added_at`, userID, habitID).Scan(&userHabit.AddedAt)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" { // unique_violation
				r.logger.Warn("Habit already added", zap.Error(err))
				return nil, appErrors.ErrHabitAlreadyAdded
			}
		}
		r.logger.Error("Failed to insert user habit", zap.Error(err))
		return nil, err
	}
	userHabit.Habit = *habit
	userHabit.UserID = userID
	return &userHabit, nil
}

func (r *habitRepositoryImpl) RemoveHabit(userID, habitID uint) error {
	_, err := r.DB.Exec(`DELETE FROM user_habit WHERE user_id = $1 AND habit_id=$2`, userID, habitID)
	if err != nil {
		return err
	}
	return nil
}

func (r *habitRepositoryImpl) GetUserHabit(userID, habitID uint) (*models.UserHabit, error) {
	var userHabit models.UserHabit
	userHabit.UserID = userID
	err := r.DB.QueryRow(
		`SELECT uh.id, h.id, h.title, h.description, h.created_at, h.image_filename, added_at FROM habit h
                                                                                       JOIN user_habit uh ON h.id = uh.habit_id
WHERE user_id = $1 AND habit_id = $2;`,
		userID,
		habitID,
	).Scan(
		&userHabit.ID,
		&userHabit.Habit.ID,
		&userHabit.Habit.Title,
		&userHabit.Habit.Description,
		&userHabit.Habit.CreatedAt,
		&userHabit.Habit.ImageFilename,
		&userHabit.AddedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, appErrors.ErrHabitNotFound
		}
		r.logger.Error("Error querying habits", zap.Error(err))
		return nil, err
	}

	streak, err := r.GetStreak(userHabit.ID)
	if err != nil {
		return nil, err
	}

	userHabit.Streak = streak

	return &userHabit, nil

}

func (r *habitRepositoryImpl) AddHabitTags(habitID uint, habitTagIDs []int) error {
	if len(habitTagIDs) == 0 {
		return nil
	}

	_, err := r.DB.Exec(`
		INSERT INTO habit_tag (habit_id, tag_id)
		SELECT $1, unnest($2::int[])
		ON CONFLICT (habit_id, tag_id) DO NOTHING
	`, habitID, pq.Array(habitTagIDs))

	return err
}

func (r *habitRepositoryImpl) GetAllTags(req request.GetAllHabitTagsRequest) ([]*models.HabitTag, int64, error) {
	baseQuery := `FROM tag`
	where := ``
	args := make([]interface{}, 0)
	argPos := 1

	// Поиск
	if req.Search != "" {
		where = fmt.Sprintf(" WHERE title ILIKE $%d", argPos)
		args = append(args, "%"+req.Search+"%")
		argPos++
	}

	// ---- total ----
	var total int64
	countQuery := `SELECT COUNT(*) ` + baseQuery + where
	if err := r.DB.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// ---- пагинация ----
	page := req.Page
	if page <= 0 {
		page = 1
	}

	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	query := `SELECT id, title ` + baseQuery + where +
		fmt.Sprintf(" ORDER BY id LIMIT $%d OFFSET $%d", argPos, argPos+1)

	argsWithPagination := append(args, pageSize, offset)

	rows, err := r.DB.Query(query, argsWithPagination...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	habitTags := make([]*models.HabitTag, 0, pageSize)

	for rows.Next() {
		var habitTag models.HabitTag
		if err := rows.Scan(&habitTag.ID, &habitTag.Title); err != nil {
			return nil, 0, err
		}
		habitTags = append(habitTags, &habitTag)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return habitTags, total, nil
}

func (r *habitRepositoryImpl) EditTag(habit *models.HabitTag) error {
	err := r.DB.QueryRow(`UPDATE tag SET title = $1 WHERE id = $2 RETURNING id, title`, habit.Title, habit.ID).Scan(&habit.ID, &habit.Title)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return appErrors.ErrHabitNotFound
		}
		return err
	}
	return nil
}

func (r *habitRepositoryImpl) DeleteTag(habitID uint) error {
	_, err := r.DB.Exec(`DELETE FROM tag WHERE id = $1`, habitID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return appErrors.ErrHabitNotFound
		}
		return err
	}
	return nil
}

func (r *habitRepositoryImpl) CreateTag(habit *models.HabitTag) error {
	err := r.DB.QueryRow(`INSERT INTO tag(title) VALUES ($1) RETURNING id, title`, habit.Title).Scan(&habit.ID, &habit.Title)
	if err != nil {
		return err
	}
	return nil
}

func (r *habitRepositoryImpl) GetTagByID(habitID uint) (*models.HabitTag, error) {
	var habitTag models.HabitTag
	err := r.DB.QueryRow(`SELECT id, title FROM tag WHERE id = $1`, habitID).Scan(&habitTag.ID, &habitTag.Title)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, appErrors.ErrHabitNotFound
		}
		return nil, err
	}
	return &habitTag, nil
}

func (r *habitRepositoryImpl) GetTagsByHabitID(habitID uint) ([]*models.HabitTag, error) {
	rows, err := r.DB.Query(
		`SELECT t.id, t.title
		 FROM tag t
		 JOIN habit_tag ht ON ht.tag_id = t.id
		 WHERE ht.habit_id = $1`,
		habitID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []*models.HabitTag

	for rows.Next() {
		var tag models.HabitTag
		if err := rows.Scan(&tag.ID, &tag.Title); err != nil {
			return nil, err
		}
		tags = append(tags, &tag)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tags, nil
}

func (r *habitRepositoryImpl) ReplaceHabitTags(habitID uint, tagIDs []int) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 1. удалить старые связи
	_, err = tx.Exec(
		`DELETE FROM habit_tag WHERE habit_id = $1`,
		habitID,
	)
	if err != nil {
		return err
	}

	// 2. добавить новые
	if len(tagIDs) > 0 {
		stmt, err := tx.Prepare(
			`INSERT INTO habit_tag (habit_id, tag_id) VALUES ($1, $2)`,
		)
		if err != nil {
			return err
		}
		defer stmt.Close()

		for _, tagID := range tagIDs {
			_, err = stmt.Exec(habitID, tagID)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}
