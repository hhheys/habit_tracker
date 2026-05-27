package userhabit

import (
	"context"
	"habit-tracker/internal/domain"
)

type Service struct {
	userHabit Repository
	streak    StreakRepository
}

func NewService(userHabit Repository, streak StreakRepository) *Service {
	return &Service{userHabit: userHabit, streak: streak}
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

func (s *Service) List(ctx context.Context, input ListUserHabitsParams) (*ListUserHabitsOutput, error) {
	input.Limit, input.Offset = validatePagination(input.Limit, input.Offset)

	uHabits, total, err := s.userHabit.ListUserHabits(ctx, domain.UserHabitListFilter{
		UserID: input.UserID, Search: input.Search, SortBy: input.SortBy,
		SortOrder: input.SortOrder, Limit: input.Limit, Offset: input.Offset,
	})
	if err != nil {
		return nil, err
	}

	habitIDsFromPage := make([]uint, len(uHabits))
	for i, uHabit := range uHabits {
		habitIDsFromPage[i] = uHabit.ID
	}

	streaks, err := s.streak.GetStreaksFromIDs(ctx, habitIDsFromPage)
	if err != nil {
		return nil, err
	}

	streaksMap := make(map[uint]*domain.Streak)
	for _, streak := range streaks {
		streaksMap[streak.UserHabitID] = streak
	}

	res := make([]*Output, len(uHabits))
	for i, uHabit := range uHabits {
		res[i] = &Output{
			Habit:  uHabit,
			Streak: streaksMap[uHabit.ID],
		}
	}

	return &ListUserHabitsOutput{
		UserHabits: res,
		Total:      total,
		Limit:      input.Limit,
		Offset:     input.Offset,
	}, nil
}

func (s *Service) Add(ctx context.Context, input AddUserHabitInput) (*domain.UserHabit, error) {
	h := domain.UserHabit{
		UserID:  input.UserID,
		HabitID: input.HabitID,
	}
	err := s.userHabit.CreateUserHabit(ctx, &h)
	if err != nil {
		return nil, err
	}
	return &h, nil
}
