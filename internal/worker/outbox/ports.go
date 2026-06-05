package outbox

import (
	"context"
	domain "habit-tracker/internal/domain/events"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	GetCreated(ctx context.Context, limit int) ([]*domain.Event, error)
	MarkDead(ctx context.Context, id uuid.UUID) error
	IncrementAttemptCountAndNextTime(ctx context.Context, id uuid.UUID, nextTime time.Time) error
	MarkSent(ctx context.Context, id uuid.UUID) error
}

type EventProducer interface {
	Publish(ctx context.Context, topic string, key string, value []byte) error
}
