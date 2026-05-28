package kafka

import (
	"context"
	"sync"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Producer struct {
	writers map[string]*kafka.Writer
	brokers []string

	logger *zap.Logger

	mutex sync.Mutex
}

func NewProducer(brokers []string, logger *zap.Logger) *Producer {
	return &Producer{
		writers: make(map[string]*kafka.Writer),
		brokers: brokers,
		logger:  logger,
		mutex:   sync.Mutex{},
	}
}

func (p *Producer) getOrCreateWriter(topic string) (*kafka.Writer, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	writer, ok := p.writers[topic]
	if ok {
		return writer, nil
	}

	writer = &kafka.Writer{
		Addr:         kafka.TCP(p.brokers...),
		Topic:        topic,
		RequiredAcks: kafka.RequireAll,
		Async:        false,

		Balancer: &kafka.Hash{},
	}

	p.writers[topic] = writer

	return writer, nil
}

func (p *Producer) Publish(ctx context.Context, topic string, key string, value []byte) error {
	writer, err := p.getOrCreateWriter(topic)
	if err != nil {
		return err
	}
	return writer.WriteMessages(ctx, kafka.Message{Key: []byte(key), Value: value})

}

func (p *Producer) Close() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	for _, writer := range p.writers {
		err := writer.Close()
		if err != nil {
			p.logger.Error("Failed to close Kafka writer", zap.Error(err))
			return err
		}
	}
	p.writers = make(map[string]*kafka.Writer)
	return nil
}
