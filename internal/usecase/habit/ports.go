package habit

import (
	"context"
	"habit-tracker/internal/domain"
)

type Repository interface {
	List(ctx context.Context, filter domain.HabitListFilter) ([]*domain.Habit, int64, error)
	GetByID(ctx context.Context, id, userID uint) (*domain.Habit, error)
	Create(ctx context.Context, habit *domain.Habit) error
	Update(ctx context.Context, habit *domain.Habit, addTagIDs, removeTagIDs []uint) error
	DeleteByID(ctx context.Context, id uint) error
}

type StreakRepository interface {
	GetStreak(ctx context.Context, userID uint) (*domain.Streak, error)
	GetUserStreaks(ctx context.Context, userID uint) ([]*domain.Streak, error)
	GetStreaksFromIDs(ctx context.Context, userHabitIDs []uint) ([]*domain.Streak, error)
}
