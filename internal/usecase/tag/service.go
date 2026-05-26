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
	tag := &domain.HabitTag{Name: input.Name}
	if err := s.tag.Create(ctx, tag); err != nil {
		return nil, err
	}
	return tag, nil
}

func (s *Service) Update(ctx context.Context, input *UpdateTagInput) (*domain.HabitTag, error) {
	tag := &domain.HabitTag{ID: input.TagID, Name: input.NewName}
	if err := s.tag.Update(ctx, tag); err != nil {
		return nil, err
	}
	return tag, nil
}

func (s *Service) Delete(ctx context.Context, tagID uint) error {
	return s.tag.Delete(ctx, tagID)
}

func (s *Service) GetByID(ctx context.Context, tagID uint) (*domain.HabitTag, error) {
	return s.tag.GetByID(ctx, tagID)
}

func (s *Service) GetAll(ctx context.Context) ([]*domain.HabitTag, error) {
	return s.tag.GetAll(ctx)
}
