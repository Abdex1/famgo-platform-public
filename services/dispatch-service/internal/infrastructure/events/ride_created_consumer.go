package events

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"

	kafkaconsumer "github.com/Abdex1/FamGo-platform/packages/kafka-sdk/consumer"
	"github.com/Abdex1/FamGo-platform/packages/kafka-sdk/governance"
	ridecontracts "github.com/Abdex1/FamGo-platform/packages/event-bus/contracts/ride"
	"github.com/Abdex1/FamGo-platform/packages/event-bus/envelope"
	"github.com/Abdex1/FamGo-platform/services/dispatch-service/internal/application/saga"
	"go.uber.org/zap"
)

// RideCreatedConsumer consumes ride.created.v1 and triggers dispatch matching.
type RideCreatedConsumer struct {
	consumer *kafkaconsumer.Consumer
	handler  *saga.DispatchSagaHandler
	logger   *zap.Logger
}

func NewRideCreatedConsumer(
	brokers []string,
	groupID string,
	handler *saga.DispatchSagaHandler,
	logger *zap.Logger,
) *RideCreatedConsumer {
	return &RideCreatedConsumer{
		consumer: kafkaconsumer.New(brokers, groupID, governance.TopicRideCreated),
		handler:  handler,
		logger:   logger,
	}
}

// Start launches the consumer loop in a background goroutine.
func (c *RideCreatedConsumer) Start(ctx context.Context) {
	go c.consumer.Consume(ctx, c.handleMessage)
}

func (c *RideCreatedConsumer) handleMessage(ctx context.Context, message kafka.Message) error {
	var env envelope.EventEnvelope
	if err := json.Unmarshal(message.Value, &env); err != nil {
		return fmt.Errorf("decode event envelope: %w", err)
	}

	payloadBytes, err := json.Marshal(env.Payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	var event ridecontracts.RideCreated
	if err := json.Unmarshal(payloadBytes, &event); err != nil {
		return fmt.Errorf("decode ride.created payload: %w", err)
	}

	if event.RideID == "" {
		return fmt.Errorf("ride.created missing ride_id")
	}

	c.logger.Info("processing ride.created.v1",
		zap.String("ride_id", event.RideID),
		zap.String("trace_id", env.TraceID),
	)

	_, err = c.handler.HandleRideCreated(ctx, event, env.TraceID, env.CorrelationID, env.RequestID)
	if err != nil {
		return fmt.Errorf("dispatch saga failed: %w", err)
	}

	return nil
}
