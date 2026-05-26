package tag

import (
	"context"
	"habit-tracker/internal/domain"
)

type Repository interface {
	Create(ctx context.Context, tag *CreateTagInput) (*domain.HabitTag, error)
	GetByID(ctx context.Context, tagID string) (*domain.HabitTag, error)
	Update(ctx context.Context, tag *UpdateTagInput) (*domain.HabitTag, error)
	Delete(ctx context.Context, tagID uint) error
	GetAll(ctx context.Context) ([]*domain.HabitTag, error)
}
