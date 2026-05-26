package streak

import (
	"context"
	"habit-tracker/internal/domain"
)

type Service struct {
	streak Repository
}

func NewService(streak Repository) Service {
	return Service{streak: streak}
}

func (s Service) CreateDailyConfirmation(ctx context.Context, input DailyConfirmationInput) (*domain.Streak, error) {
	return s.streak.CreateDailyConfirmation(ctx, input.UserID, input.HabitID)
}

func (s Service) GetHeatmap(ctx context.Context, input HeatmapInput) ([]*domain.HeatmapDay, error) {
	return s.streak.GetHeatmap(ctx, input)
}
