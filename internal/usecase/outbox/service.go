package outbox

import (
	"context"
	domain "habit-tracker/internal/domain/events"
)

type Outbox struct {
	outbox Repository
}

func NewEventService(outbox Repository) Outbox {
	return Outbox{
		outbox: outbox,
	}
}

// PublishEventToOutbox publishes any event to outbox table. Use subjectKey to transfer userID or habitID.
func (s *Outbox) PublishEventToOutbox(ctx context.Context, event *domain.Event) error {
	return s.outbox.Publish(ctx, event)
}
