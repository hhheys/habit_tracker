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

func (r *habitRepositoryImpl) GetAllHabits(req *request.GetAllHabitsRequest, requestUserID uint) (int64, []*models.Habit, error) {
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

	var total int64
	err := r.DB.QueryRow(
		`SELECT COUNT(*) FROM habit 
		 WHERE title ILIKE $1 OR description ILIKE $1`,
		search,
	).Scan(&total)
	if err != nil {
		r.logger.Error("Error counting habits", zap.Error(err))
		return 0, nil, err
	}

	offset := (req.Page - 1) * req.PageSize

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
		AND uh.user_id = $4
	WHERE 
		h.title ILIKE $1 
		OR h.description ILIKE $1
	ORDER BY %s
	LIMIT $2 OFFSET $3;
`, req.Sort)

	rows, err := r.DB.Query(
		query,
		search,
		req.PageSize,
		offset,
		requestUserID,
	)
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
