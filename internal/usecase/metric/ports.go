package metric

import (
	"context"
	"habit-tracker/internal/domain/achievement"
	"habit-tracker/internal/domain/events"
)

type Repository interface {
	UpdateUserMetric(ctx context.Context, userMetric *achievement.UserMetric) error
	UpdateUserHabitMetric(ctx context.Context, userMetric *achievement.UserHabitMetric) error

	GetMetricByKey(ctx context.Context, key string) (*achievement.Metric, error)
}

type UserHabitRepository interface {
	GetTotalUserHabits(ctx context.Context, userID uint) (int, error)
}

type TXManager interface {
	WithTx(ctx context.Context, fn func(context.Context) error) error
}

type EventPublisher interface {
	Publish(ctx context.Context, event *events.Event) error
}
