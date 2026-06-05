package achievement

import (
	"context"
	"habit-tracker/internal/domain/achievement"
	"habit-tracker/internal/domain/events"
)

type UserAchievementRepository interface {
	// GetExpectedAchievements returns the expected achievements for a user based on their current metric
	GetExpectedAchievements(ctx context.Context, userID uint, metricKey achievement.MetricKey) ([]*achievement.Achievement, error)
	ListUserAchievements(ctx context.Context, userID uint, limit, offset int) ([]*achievement.UserAchievementListItem, int64, error)
	UnlockUserAchievementByCode(ctx context.Context, userID uint, code string) (*achievement.UserAchievement, error)
}

type UserMetricRepository interface {
	GetUserMetricByKeyAndUserID(ctx context.Context, userID uint, key string) (*achievement.UserMetric, error)
	GetUserHabitMetricByKeyAndUserHabitID(ctx context.Context, userHabitID uint, key string) (*achievement.UserHabitMetric, error)
}

type UserHabitRepository interface {
	GetUserIDByUserHabitID(ctx context.Context, userHabitID uint) (uint, error)
}

type TXManager interface {
	WithTx(ctx context.Context, fn func(context.Context) error) error
}

type EventPublisher interface {
	Publish(ctx context.Context, event *events.Event) error
}
