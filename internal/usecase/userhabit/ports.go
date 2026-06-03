package userhabit

import (
	"context"
	"habit-tracker/internal/domain"
	"habit-tracker/internal/domain/events"
)

type Repository interface {
	ListUserHabits(ctx context.Context, filter domain.UserHabitListFilter) ([]*domain.UserHabit, int64, error)
	CreateUserHabit(ctx context.Context, habit *domain.UserHabit) error
}

type StreakRepository interface {
	GetStreak(ctx context.Context, userID uint) (*domain.Streak, error)
	GetUserStreaks(ctx context.Context, userID uint) ([]*domain.Streak, error)
	GetStreaksFromIDs(ctx context.Context, userHabitIDs []uint) ([]*domain.Streak, error)
}

type EventPublisher interface {
	Publish(ctx context.Context, event *events.Event) error
}

type TXManager interface {
	WithTx(ctx context.Context, fn func(context.Context) error) error
}
