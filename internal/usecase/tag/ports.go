package tag

import (
	"context"
	"habit-tracker/internal/domain"
)

type Repository interface {
	Create(ctx context.Context, tag *domain.HabitTag) error
	GetByID(ctx context.Context, tagID uint) (*domain.HabitTag, error)
	Update(ctx context.Context, tag *domain.HabitTag) error
	Delete(ctx context.Context, tagID uint) error
	GetAll(ctx context.Context) ([]*domain.HabitTag, error)
}
