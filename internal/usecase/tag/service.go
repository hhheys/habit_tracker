package tag

import (
	"context"
	"habit-tracker/internal/domain"
)

type Service struct {
	tag Repository
}

func NewService(tag Repository) Service {
	return Service{tag: tag}
}

func (s *Service) Create(ctx context.Context, input *CreateTagInput) (*domain.HabitTag, error) {
	return s.tag.Create(ctx, input)
}

func (s *Service) Update(ctx context.Context, input *UpdateTagInput) (*domain.HabitTag, error) {
	return s.tag.Update(ctx, input)
}

func (s *Service) Delete(ctx context.Context, tagID uint) error {
	return s.tag.Delete(ctx, tagID)
}

func (s *Service) GetByID(ctx context.Context, tagID string) (*domain.HabitTag, error) {
	return s.tag.GetByID(ctx, tagID)
}

func (s *Service) GetAll(ctx context.Context) ([]*domain.HabitTag, error) {
	return s.tag.GetAll(ctx)
}
