package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"

	kafkacontracts "github.com/Abdex1/FamGo-platform/packages/kafka-sdk/contracts"
	kafkaconfig "github.com/Abdex1/FamGo-platform/packages/kafka-sdk/config"
	kafkaproducer "github.com/Abdex1/FamGo-platform/packages/kafka-sdk/producer"
	ridecontracts "github.com/Abdex1/FamGo-platform/packages/event-bus/contracts/ride"
	"github.com/Abdex1/FamGo-platform/packages/event-bus/envelope"
	"github.com/Abdex1/FamGo-platform/services/dispatch-service/internal/application/saga"
	"github.com/Abdex1/FamGo-platform/services/dispatch-service/internal/infrastructure/events"
	"github.com/Abdex1/FamGo-platform/services/dispatch-service/internal/infrastructure/observability"
)

type ProducerBridge struct {
	producer *kafkaproducer.Producer
}

func NewProducerBridge(brokers []string, logger *zap.Logger) *ProducerBridge {
	return &ProducerBridge{
		producer: kafkaproducer.New(kafkaconfig.Config{
			Brokers:      brokers,
			BatchTimeout: 10,
			BatchSize:    100,
		}, logger),
	}
}

func (p *ProducerBridge) Publish(
	ctx context.Context,
	topic, key string,
	payload []byte,
	headers map[string]string,
) error {
	var env envelope.EventEnvelope
	if err := json.Unmarshal(payload, &env); err != nil {
		return fmt.Errorf("invalid envelope payload: %w", err)
	}
	return p.producer.Publish(ctx, topic, key, toKafkaEnvelope(env))
}

func toKafkaEnvelope(env envelope.EventEnvelope) kafkacontracts.EventEnvelope {
	return kafkacontracts.EventEnvelope{
		EventID:        env.EventID,
		EventType:      env.EventType,
		EventVersion:   env.EventVersion,
		TraceID:        env.TraceID,
		SpanID:         env.SpanID,
		CorrelationID:  env.CorrelationID,
		RequestID:      env.RequestID,
		Producer:       env.Service,
		Environment:    env.Environment,
		OccurredAt:     env.OccurredAt,
		PartitionKey:   env.PartitionKey,
		IdempotencyKey: env.IdempotencyKey,
		Headers:        env.Headers,
		Payload:        env.Payload,
	}
}

func (p *ProducerBridge) AsDispatchPublisher(service, env string) *events.KafkaPublisher {
	return events.NewKafkaPublisher(p.Publish, service, env)
}

type RideCreatedConsumer struct {
	reader  *kafka.Reader
	handler *saga.DispatchSagaHandler
	logger  *zap.Logger
}

func NewRideCreatedConsumer(brokers []string, groupID string, handler *saga.DispatchSagaHandler, logger *zap.Logger) *RideCreatedConsumer {
	return &RideCreatedConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			GroupID: groupID,
			Topic:   "ride.created.v1",
			MinBytes: 1,
			MaxBytes: 10e6,
		}),
		handler: handler,
		logger:  logger,
	}
}

func (c *RideCreatedConsumer) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		message, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			c.logger.Warn("kafka read error", zap.Error(err))
			continue
		}

		if err := c.handleMessage(ctx, message); err != nil {
			c.logger.Error("ride.created handler failed", zap.Error(err))
		}
	}
}

func (c *RideCreatedConsumer) handleMessage(ctx context.Context, message kafka.Message) error {
	var env envelope.EventEnvelope
	if err := json.Unmarshal(message.Value, &env); err != nil {
		return err
	}

	payloadBytes, err := json.Marshal(env.Payload)
	if err != nil {
		return err
	}

	var event ridecontracts.RideCreated
	if err := json.Unmarshal(payloadBytes, &event); err != nil {
		return err
	}

	ctx = observability.EnrichContext(ctx, env.TraceID, env.CorrelationID, env.RequestID)
	_, err = c.handler.HandleRideCreated(ctx, event, env.TraceID, env.CorrelationID, env.RequestID)
	return err
}

func (c *RideCreatedConsumer) Close() error {
	return c.reader.Close()
}
