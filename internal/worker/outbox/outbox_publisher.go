package outbox

import (
	"context"
	"habit-tracker/internal/domain/events"
	"sync"
	"time"

	"go.uber.org/zap"
)

var AttemptSleepDuration []time.Duration = []time.Duration{
	0,                // First attempt, no delay
	time.Second * 10, // Second attempt, 3-second delay
	time.Second * 30, // Third attempt, 5-second delay
	time.Minute,      // Fourth attempt, 10-second delay
	time.Minute * 5,
	time.Minute * 15,
}

type EventPublisher struct {
	outbox     Repository
	producer   EventProducer
	topic      string
	maxWorkers int

	maxAttempts int

	logger *zap.Logger

	WG *sync.WaitGroup
}

func NewEventPublisher(outbox Repository, producer EventProducer, topic string, logger *zap.Logger) EventPublisher {
	return EventPublisher{
		outbox:      outbox,
		producer:    producer,
		topic:       topic,
		maxWorkers:  5,
		maxAttempts: 5,
		WG:          &sync.WaitGroup{},
		logger:      logger,
	}
}

func (p *EventPublisher) Run(ctx context.Context, interval time.Duration) {
	if interval <= 0 {
		interval = 5 * time.Second
	}

	p.runOnceAndLog(ctx)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.runOnceAndLog(ctx)
		}
	}
}

func (p *EventPublisher) runOnceAndLog(ctx context.Context) {
	if err := p.RunOnce(ctx); err != nil {
		p.logger.Error("failed to publish outbox events", zap.Error(err))
	}
}

func (p *EventPublisher) RunOnce(ctx context.Context) error {
	processEvents, err := p.outbox.GetCreated(ctx, 100)
	if err != nil {
		return err
	}

	sem := make(chan struct{}, p.maxWorkers)
	var wg sync.WaitGroup

	for _, event := range processEvents {
		sem <- struct{}{}
		wg.Add(1)

		go func(processEvent *events.Event) {
			defer func() {
				<-sem
				wg.Done()
			}()

			publishErr := p.producer.Publish(
				ctx,
				p.topic,
				processEvent.PartitionKey,
				processEvent.Payload,
			)

			if publishErr != nil {
				nextAttempt := processEvent.AttemptCount + 1

				if nextAttempt >= p.maxAttempts {
					p.logger.Info("event reached maximum attempts, marking as dead", zap.String("event_id", processEvent.EventID.String()))

					if deadErr := p.outbox.MarkDead(ctx, processEvent.EventID); deadErr != nil {
						p.logger.Error("failed to mark event as dead", zap.Error(deadErr), zap.String("event_id", processEvent.EventID.String()))
					}

					return
				}

				p.logger.Error("failed to publish event", zap.Error(publishErr), zap.String("event_id", processEvent.EventID.String()))

				incrementErr := p.outbox.IncrementAttemptCountAndNextTime(
					ctx,
					processEvent.EventID,
					time.Now().UTC().Add(AttemptSleepDuration[processEvent.AttemptCount]),
				)
				if incrementErr != nil {
					p.logger.Error("failed to increment attempt count and next time", zap.Error(incrementErr), zap.String("event_id", processEvent.EventID.String()))
				}

				return
			}

			if markErr := p.outbox.MarkSent(ctx, processEvent.EventID); markErr != nil {
				p.logger.Error("failed to mark event as sent", zap.Error(markErr), zap.String("event_id", processEvent.EventID.String()))
			}
		}(event)
	}

	wg.Wait()
	return nil
}
