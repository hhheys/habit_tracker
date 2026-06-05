package streak

import (
	"context"
	"habit-tracker/internal/domain"
	"habit-tracker/internal/domain/events"
	"habit-tracker/internal/usecase/eventpublisher"
	"strconv"
)

type Service struct {
	streak         Repository
	eventPublisher EventPublisher
	txManager      TXManager
}

func NewService(streak Repository, publisher EventPublisher, txManager TXManager) Service {
	return Service{streak: streak, eventPublisher: publisher, txManager: txManager}
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

			return eventpublisher.Publish(
				customContext,
				s.eventPublisher,
				events.EventTypeHabitStreakConfirmed,
				strconv.Itoa(int(input.UserID)),
				streak,
			)
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
