package streak

import (
	"context"
	"habit-tracker/internal/domain"
	"habit-tracker/internal/domain/events"
)

type Repository interface {
	CreateDailyConfirmation(ctx context.Context, userID, habitID uint) (*domain.Streak, error)
	GetHeatmap(ctx context.Context, filter domain.HeatmapFilter) ([]*domain.HeatmapDay, error)
}

type TXManager interface {
	WithTx(ctx context.Context, fn func(context.Context) error) error
}

type EventPublisher interface {
	Publish(ctx context.Context, event *events.Event) error
}
