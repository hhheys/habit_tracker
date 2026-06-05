package streak

import (
	"context"
	"database/sql"
	"errors"
	"habit-tracker/internal/adapter/postgres/txmanager"
	"habit-tracker/internal/domain"

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

func (r *Repository) CreateDailyConfirmation(ctx context.Context, userID, habitID uint) (*domain.Streak, error) {
	executorType := txmanager.ExecutorFromContext(ctx, r.db)

	var userHabitID uint
	err := executorType.QueryRowContext(ctx, `
		SELECT id FROM user_habit WHERE user_id = $1 AND habit_id = $2`,
		userID, habitID,
	).Scan(&userHabitID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrHabitNotAdded
	}
	if err != nil {
		return nil, err
	}
	_, err = executorType.ExecContext(ctx, `INSERT INTO daily_confirmation (user_habit_id) VALUES ($1)`, userHabitID)
	var pqErr *pq.Error
	if errors.As(err, &pqErr) && pqErr.Code == "23505" {
		return nil, domain.ErrHabitAlreadyConfirmed
	}
	if err != nil {
		return nil, err
	}
	return r.GetStreak(ctx, userHabitID)
}

func (r *Repository) GetHeatmap(ctx context.Context, filter domain.HeatmapFilter) ([]*domain.HeatmapDay, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT DATE(dc.confirmed_at), COUNT(*)
		FROM daily_confirmation dc
		JOIN user_habit uh ON uh.id = dc.user_habit_id
		WHERE uh.user_id = $1
		  AND ($2 = '' OR DATE(dc.confirmed_at) >= $2::date)
		  AND ($3 = '' OR DATE(dc.confirmed_at) <= $3::date)
		GROUP BY DATE(dc.confirmed_at)
		ORDER BY DATE(dc.confirmed_at)`, filter.UserID, filter.StartDate, filter.EndDate)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		closeErr := rows.Close()
		if closeErr != nil {
			r.log.Error("failed to close rows", zap.Error(closeErr))
		}
	}(rows)
	days := make([]*domain.HeatmapDay, 0)
	for rows.Next() {
		var day domain.HeatmapDay
		if err := rows.Scan(&day.Date, &day.Count); err != nil {
			return nil, err
		}
		days = append(days, &day)
	}
	return days, rows.Err()
}

func (r *Repository) GetStreak(ctx context.Context, userHabitID uint) (*domain.Streak, error) {
	executorType := txmanager.ExecutorFromContext(ctx, r.db)

	var streak domain.Streak
	var current sql.NullInt64
	streak.UserHabitID = userHabitID
	err := executorType.QueryRowContext(ctx, `
		WITH dates AS (
			SELECT DISTINCT DATE(confirmed_at) AS day
			FROM daily_confirmation WHERE user_habit_id = $1
		), grouped AS (
			SELECT day, day - ROW_NUMBER() OVER (ORDER BY day) * INTERVAL '1 day' AS grp
			FROM dates
		), streaks AS (
			SELECT COUNT(*) AS streak_length, MAX(day) AS last_day
			FROM grouped GROUP BY grp
		)
		SELECT COALESCE(MAX(streak_length), 0),
		       MAX(streak_length) FILTER (WHERE last_day >= CURRENT_DATE - INTERVAL '1 day'),
		       EXISTS (SELECT 1 FROM dates WHERE day = CURRENT_DATE)
		FROM streaks`, userHabitID,
	).Scan(&streak.LongestStreak, &current, &streak.IsConfirmedToday)
	if err != nil {
		return nil, err
	}
	if current.Valid {
		value := int(current.Int64)
		streak.CurrentStreak = &value
	}
	return &streak, nil
}

func (r *Repository) GetUserStreaks(ctx context.Context, userID uint) ([]*domain.Streak, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id FROM user_habit WHERE user_id = $1 ORDER BY id`, userID)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			r.log.Error("failed to close rows", zap.Error(err))
		}
	}(rows)
	var ids []uint
	for rows.Next() {
		var id uint
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return r.GetStreaksFromIDs(ctx, ids)
}

func (r *Repository) GetStreaksFromIDs(ctx context.Context, ids []uint) ([]*domain.Streak, error) {
	streaks := make([]*domain.Streak, 0, len(ids))
	for _, id := range ids {
		streak, err := r.GetStreak(ctx, id)
		if err != nil {
			return nil, err
		}
		streaks = append(streaks, streak)
	}
	return streaks, nil
}
