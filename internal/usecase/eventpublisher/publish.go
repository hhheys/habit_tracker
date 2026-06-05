package eventpublisher

import (
	"context"
	"encoding/json"
	"habit-tracker/internal/domain/events"
	"time"
)

type Publisher interface {
	Publish(ctx context.Context, event *events.Event) error
}

func Publish(
	ctx context.Context,
	publisher Publisher,
	eventType events.EventType,
	partitionKey string,
	payload any,
) error {
	eventPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	event := events.Event{
		OccurredAt:       time.Now().UTC(),
		EventType:        eventType,
		EventTypeVersion: 1,
		PartitionKey:     partitionKey,
		Payload:          eventPayload,
	}

	return publisher.Publish(ctx, &event)
}
