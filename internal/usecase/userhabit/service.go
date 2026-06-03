package userhabit

import (
	"context"
	"encoding/json"
	"habit-tracker/internal/domain"
	"habit-tracker/internal/domain/events"
	"strconv"
	"time"
)

type Service struct {
	userHabit      Repository
	streak         StreakRepository
	eventPublisher EventPublisher
	txManager      TXManager
}

func NewService(userHabit Repository, streak StreakRepository, publisher EventPublisher, txManager TXManager) *Service {
	return &Service{userHabit: userHabit, streak: streak, eventPublisher: publisher, txManager: txManager}
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

	err := s.txManager.WithTx(
		ctx,
		func(customContext context.Context) error {
			createErr := s.userHabit.CreateUserHabit(customContext, &h)
			if createErr != nil {
				return createErr
			}

			eventPayload, jsonErr := json.Marshal(h)
			if jsonErr != nil {
				return jsonErr
			}

			event := events.Event{
				OccurredAt:       time.Now().UTC(),
				EventType:        events.EventTypeUserHabitAdded,
				EventTypeVersion: 1,
				PartitionKey:     strconv.Itoa(int(input.UserID)),
				Payload:          eventPayload,
			}

			createErr = s.eventPublisher.Publish(customContext, &event)
			if createErr != nil {
				return createErr
			}
			return nil
		},
	)

	if err != nil {
		return nil, err
	}

	return &h, nil
}
