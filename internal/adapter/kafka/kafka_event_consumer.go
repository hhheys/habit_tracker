package kafka

import (
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

func (p *Consumer) getOrCreateReader(topic, groupID string) (*kafka.Reader, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	reader, ok := p.readers[topic]
	if ok {
		return reader, nil
	}

	reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers: p.brokers,
		Topic:   topic,
		GroupID: groupID,
	})

	p.readers[topic] = reader

	return reader, nil
}
