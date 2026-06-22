
// STEP 7 — CREATE PRODUCER

// packages/kafka-sdk/producer/producer.go

package producer

import (
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"

	"github.com/Abdex1/FamGo-platform/packages/kafka-sdk/config"
	"github.com/Abdex1/FamGo-platform/packages/kafka-sdk/contracts"
)

type Producer struct {
	writer *kafka.Writer
	logger *zap.Logger
}

func New(cfg config.Config, logger *zap.Logger) *Producer {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(cfg.Brokers...),
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: cfg.BatchTimeout,
		BatchSize:    cfg.BatchSize,
		RequiredAcks: kafka.RequireAll,
		Async:        false,
	}

	return &Producer{
		writer: writer,
		logger: logger,
	}
}

func (p *Producer) Publish(
	ctx context.Context,
	topic string,
	key string,
	envelope contracts.EventEnvelope,
) error {
	payload, err := json.Marshal(envelope)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: payload,
		Time:  time.Now().UTC(),
		Headers: []kafka.Header{
			{
				Key:   "trace_id",
				Value: []byte(envelope.TraceID),
			},
			{
				Key:   "correlation_id",
				Value: []byte(envelope.CorrelationID),
			},
		},
	}

	return p.writer.WriteMessages(ctx, msg)
}
