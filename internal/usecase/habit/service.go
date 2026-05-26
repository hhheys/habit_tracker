package habit

import (
	"context"
	"habit-tracker/internal/domain"
)

type Service struct {
	habit  Repository
	streak StreakRepository
}

func NewService(habit Repository, streak StreakRepository) *Service {
	return &Service{habit: habit, streak: streak}
}

// Поскольку бизнес логика может быть разной в засимости от сущности, сделаю валидацию отдельной для каждого юзкейса
func validatePagination(limit, offset int) (int, int) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}

func (s *Service) ListHabits(ctx context.Context, input ListHabitsParams) (ListHabitsOutput, int64, error) {
	input.Limit, input.Offset = validatePagination(input.Limit, input.Offset)

	habits, total, err := s.habit.List(ctx, input)
	if err != nil {
		return ListHabitsOutput{}, 0, err
	}

	return ListHabitsOutput{
		Habits: habits,
		Total:  total,
		Limit:  input.Limit,
		Offset: input.Offset,
	}, total, nil
}

func (s *Service) AddHabit(ctx context.Context, id uint) (*domain.Streak, error) {
	return s.streak.GetStreak(ctx, id)
}

func (s *Service) GetByID(ctx context.Context, id uint) (*domain.Habit, error) {
	return s.habit.GetByID(ctx, id)
}

func (s *Service) Create(ctx context.Context, habit *CreateHabitInput) (*domain.Habit, error) {
	habitTags := make([]*domain.HabitTag, len(habit.Tags))
	for i, tagID := range habit.Tags {
		habitTags[i] = &domain.HabitTag{
			ID: tagID,
		}
	}

	h := domain.Habit{
		Title:         habit.Title,
		Description:   habit.Description,
		Tags:          habitTags,
		ImageFilename: habit.ImageFilename,
	}
	err := s.habit.Create(ctx, &h)
	if err != nil {
		return nil, err
	}
	return &h, nil
}

func (s *Service) Update(ctx context.Context, habit *UpdateHabitInput) (*domain.Habit, error) {
	return s.habit.Update(ctx, habit)
}

func (s *Service) DeleteByID(ctx context.Context, id uint) error {
	return s.habit.DeleteByID(ctx, id)
}
