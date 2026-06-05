package outbox

import (
	"context"
	domain "habit-tracker/internal/domain/events"
)

type Repository interface {
	Publish(ctx context.Context, event *domain.Event) error
}
