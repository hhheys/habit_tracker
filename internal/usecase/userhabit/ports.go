package userhabit

import (
	"context"
	"habit-tracker/internal/domain"
)

type Repository interface {
	ListUserHabits(ctx context.Context, params ListUserHabitsParams) ([]*domain.UserHabit, int64, error)
	CreateUserHabit(ctx context.Context, habit *domain.UserHabit) error
	//UpdateUserHabit(ctx context.Context, habit *UpdateUserHabitInput) (*domain.UserHabit, error)
	//DeleteUserHabitByID(ctx context.Context, id uint) error
}

type StreakRepository interface {
	GetStreak(ctx context.Context, userID uint) (*domain.Streak, error)
	GetUserStreaks(ctx context.Context, userID uint) ([]*domain.Streak, error)
	GetStreaksFromIDs(ctx context.Context, userHabitIDs []uint) ([]*domain.Streak, error)
}
