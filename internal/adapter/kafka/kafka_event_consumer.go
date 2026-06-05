package kafka

import (
	"context"
	"encoding/json"
	"habit-tracker/internal/domain/events"
	"strings"
	"sync"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Consumer struct {
	readers map[string]*kafka.Reader
	brokers []string

	logger *zap.Logger

	mutex sync.Mutex
}

func NewConsumer(brokers []string, logger *zap.Logger) *Consumer {
	return &Consumer{
		readers: make(map[string]*kafka.Reader),
		brokers: brokers,
		logger:  logger,
		mutex:   sync.Mutex{},
	}
}

func (p *Consumer) getOrCreateReader(topic, groupID string) (*kafka.Reader, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	readerKey := strings.Join([]string{topic, groupID}, "\x00")
	reader, ok := p.readers[readerKey]
	if ok {
		return reader, nil
	}

	reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers: p.brokers,
		Topic:   topic,
		GroupID: groupID,
	})

	p.readers[readerKey] = reader

	return reader, nil
}

func (p *Consumer) ConsumeEvents(
	ctx context.Context,
	topic string,
	groupID string,
	handler func(context.Context, events.Event) error,
) error {
	reader, err := p.getOrCreateReader(topic, groupID)
	if err != nil {
		return err
	}

	for {
		message, err := reader.FetchMessage(ctx)
		if err != nil {
			return err
		}

		var event events.Event
		if err := json.Unmarshal(message.Value, &event); err != nil {
			return err
		}

		if err := handler(ctx, event); err != nil {
			return err
		}

		if err := reader.CommitMessages(ctx, message); err != nil {
			return err
		}
	}
}

func (p *Consumer) Close() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	for _, reader := range p.readers {
		if err := reader.Close(); err != nil {
			p.logger.Error("failed to close Kafka reader", zap.Error(err))
			return err
		}
	}

	p.readers = make(map[string]*kafka.Reader)
	return nil
}
