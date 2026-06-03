package streak

import (
	"context"
	"encoding/json"
	"habit-tracker/internal/domain"
	"habit-tracker/internal/domain/events"
	"strconv"
	"time"
)

type Service struct {
	streak         Repository
	eventPublisher EventPublisher
	txManager      TXManager
}

func NewService(streak Repository, txManager TXManager) Service {
	return Service{streak: streak, txManager: txManager}
}

func (s Service) CreateDailyConfirmation(ctx context.Context, input DailyConfirmationInput) (*domain.Streak, error) {
	var streak *domain.Streak
	var createErr error
	err := s.txManager.WithTx(
		ctx,
		func(customContext context.Context) error {
			streak, createErr = s.streak.CreateDailyConfirmation(ctx, input.UserID, input.HabitID)
			if createErr != nil {
				return createErr
			}

			eventPayload, jsonErr := json.Marshal(streak)
			if jsonErr != nil {
				return jsonErr
			}

			event := events.Event{
				OccurredAt:       time.Now().UTC(),
				EventType:        events.EventTypeHabitStreakConfirmed,
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

	return streak, nil
}

func (s Service) GetHeatmap(ctx context.Context, input HeatmapInput) ([]*domain.HeatmapDay, error) {
	return s.streak.GetHeatmap(ctx, domain.HeatmapFilter{
		UserID: input.UserID, StartDate: input.StartDate, EndDate: input.EndDate,
	})
}
