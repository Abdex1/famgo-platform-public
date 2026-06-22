package events

import (
	"context"

	"github.com/segmentio/kafka-go"
)

// RawKafkaPublisher writes pre-serialized event envelopes to Kafka topics.
type RawKafkaPublisher struct {
	writer *kafka.Writer
}

func NewRawKafkaPublisher(brokers []string) *RawKafkaPublisher {
	return &RawKafkaPublisher{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Balancer:     &kafka.LeastBytes{},
			RequiredAcks: kafka.RequireAll,
			Async:        false,
		},
	}
}

func (p *RawKafkaPublisher) Publish(
	ctx context.Context,
	topic, key string,
	payload []byte,
	headers map[string]string,
) error {
	msg := kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: payload,
	}
	for k, v := range headers {
		msg.Headers = append(msg.Headers, kafka.Header{Key: k, Value: []byte(v)})
	}
	return p.writer.WriteMessages(ctx, msg)
}

func (p *RawKafkaPublisher) Close() error {
	if p.writer == nil {
		return nil
	}
	return p.writer.Close()
}
