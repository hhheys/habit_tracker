package metric

import (
	"context"
	"database/sql"
	"habit-tracker/internal/adapter/postgres/txmanager"
	"habit-tracker/internal/domain/achievement"

	"go.uber.org/zap"
)

type Repository struct {
	db  *sql.DB
	log *zap.Logger
}

func NewRepository(db *sql.DB, log *zap.Logger) *Repository {
	return &Repository{db: db, log: log}
}

func (r *Repository) UpdateUserMetric(ctx context.Context, userMetric *achievement.UserMetric) error {
	executor := txmanager.ExecutorFromContext(ctx, r.db)

	return executor.QueryRowContext(ctx, `
		INSERT INTO user_metrics (user_id, metric_id, value)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, metric_id)
		DO UPDATE SET value = EXCLUDED.value, updated_at = NOW()
		RETURNING id, updated_at`,
		userMetric.UserID,
		userMetric.MetricID,
		userMetric.Value,
	).Scan(&userMetric.ID, &userMetric.UpdatedAt)
}

func (r *Repository) UpdateUserHabitMetric(ctx context.Context, userMetric *achievement.UserHabitMetric) error {
	executor := txmanager.ExecutorFromContext(ctx, r.db)

	return executor.QueryRowContext(ctx, `
		INSERT INTO user_habit_metrics (user_habit_id, metric_id, value)
		SELECT $1, m.id, $3
		FROM metric m
		WHERE m.metric_key = $2
		ON CONFLICT (user_habit_id, metric_id)
		DO UPDATE SET value = EXCLUDED.value, updated_at = NOW()
		RETURNING id, updated_at`,
		userMetric.UserHabitID,
		userMetric.MetricKey,
		userMetric.Value,
	).Scan(&userMetric.ID, &userMetric.UpdatedAt)
}

func (r *Repository) GetMetricByKey(ctx context.Context, key string) (*achievement.Metric, error) {
	metric := &achievement.Metric{}
	err := r.db.QueryRowContext(ctx, `
		SELECT id, metric_key
		FROM metric
		WHERE metric_key = $1`,
		key,
	).Scan(&metric.ID, &metric.Key)
	if err != nil {
		return nil, err
	}

	return metric, nil
}

func (r *Repository) GetUserMetricByKeyAndUserID(ctx context.Context, userID uint, key string) (*achievement.UserMetric, error) {
	executor := txmanager.ExecutorFromContext(ctx, r.db)

	userMetric := &achievement.UserMetric{}
	err := executor.QueryRowContext(ctx, `
		SELECT um.id, um.user_id, um.metric_id, um.value, um.updated_at
		FROM user_metrics um
		JOIN metric m ON m.id = um.metric_id
		WHERE um.user_id = $1
		  AND m.metric_key = $2`,
		userID,
		key,
	).Scan(
		&userMetric.ID,
		&userMetric.UserID,
		&userMetric.MetricID,
		&userMetric.Value,
		&userMetric.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return userMetric, nil
}

func (r *Repository) GetUserHabitMetricByKeyAndUserHabitID(ctx context.Context, userHabitID uint, key string) (*achievement.UserHabitMetric, error) {
	executor := txmanager.ExecutorFromContext(ctx, r.db)

	userHabitMetric := &achievement.UserHabitMetric{}
	err := executor.QueryRowContext(ctx, `
		SELECT uhm.id, uhm.user_habit_id, m.metric_key, uhm.value, uhm.updated_at
		FROM user_habit_metrics uhm
		JOIN metric m ON m.id = uhm.metric_id
		WHERE uhm.user_habit_id = $1
		  AND m.metric_key = $2`,
		userHabitID,
		key,
	).Scan(
		&userHabitMetric.ID,
		&userHabitMetric.UserHabitID,
		&userHabitMetric.MetricKey,
		&userHabitMetric.Value,
		&userHabitMetric.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return userHabitMetric, nil
}
