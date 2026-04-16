package postgres

import (
	"database/sql"
	"errors"
	appErrors "habit-tracker/internal/errors"
	"habit-tracker/internal/models"

	"github.com/lib/pq"
	"go.uber.org/zap"
)

//go:generate mockery --name=HabitRepository --output=../../../mocks --outpkg=mocks
type StreakRepository interface {
	CreateDailyConfirmation(confirmation *models.DailyConfirmation) error
	GetStreak(userHabitID uint) (*models.Streak, error)
}

type streakRepositoryImpl struct {
	DB     *sql.DB
	logger *zap.Logger
}

func NewStreakRepository(db *sql.DB, logger *zap.Logger) StreakRepository {
	return &streakRepositoryImpl{
		DB:     db,
		logger: logger,
	}
}

func (r *streakRepositoryImpl) CreateDailyConfirmation(confirmation *models.DailyConfirmation) error {
	err := r.DB.QueryRow(
		`INSERT INTO daily_confirmation (user_habit_id) VALUES ($1) RETURNING id, confirmed_at`,
		confirmation.UserHabitID,
	).Scan(
		&confirmation.ID,
		&confirmation.ConfirmedAt,
	)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case "23505": // unique_violation
				if pqErr.Constraint == "unique_user_habit_day" {
					return appErrors.ErrHabitAlreadyConfirmed
				}
			case "23503":
				return appErrors.ErrHabitNotAdded
			}
		}
		return err
	}
	return nil
}

func (r *streakRepositoryImpl) GetStreak(userHabitID uint) (*models.Streak, error) {
	var streak models.Streak
	streak.UserHabitID = userHabitID
	err := r.DB.QueryRow(`
WITH dates AS (
    SELECT DISTINCT DATE(confirmed_at) AS day
    FROM daily_confirmation
    WHERE user_habit_id = $1
),

     ordered AS (
         SELECT
             day,
             ROW_NUMBER() OVER (ORDER BY day) AS rn
         FROM dates
     ),

     groups AS (
         SELECT
             day,
             (day - rn * INTERVAL '1 day') AS grp
         FROM ordered
     ),

     streaks AS (
         SELECT
             COUNT(*) AS streak_length,
             MAX(day) AS last_day
         FROM groups
         GROUP BY grp
     ),

     longest AS (
         SELECT COALESCE(MAX(streak_length), 0) AS longest_streak
         FROM streaks
     ),

     current AS (
         SELECT streak_length AS current_streak
         FROM streaks
         WHERE last_day >= CURRENT_DATE - INTERVAL '1 day'
         ORDER BY streak_length DESC
         LIMIT 1
     ),

     today AS (
         SELECT EXISTS (
             SELECT 1
             FROM dates
             WHERE day = CURRENT_DATE
         ) AS is_confirmed
     )

SELECT
    (SELECT longest_streak FROM longest) AS longest_streak,

    (SELECT current_streak FROM current) AS current_streak,

    (SELECT is_confirmed FROM today) AS is_confirmed;`, userHabitID).Scan(
		&streak.LongestStreak,
		&streak.CurrentStreak,
		&streak.IsConfirmedToday,
	)
	if err != nil {
		return nil, err
	}
	return &streak, nil
}
