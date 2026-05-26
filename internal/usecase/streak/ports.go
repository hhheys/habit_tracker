package streak

import (
	"context"
	"habit-tracker/internal/domain"
)

type Repository interface {
	CreateDailyConfirmation(ctx context.Context, userID, habitID uint) (*domain.Streak, error)
	GetHeatmap(ctx context.Context, input HeatmapInput) ([]*domain.HeatmapDay, error)
}
