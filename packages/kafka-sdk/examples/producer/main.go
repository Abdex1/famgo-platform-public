
// STEP 14 — CREATE EXAMPLE PRODUCER

// packages/kafka-sdk/examples/producer/main.go


// IMPLEMENTATION

package main

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/Abdex1/FamGo-platform/packages/kafka-sdk/config"
	"github.com/Abdex1/FamGo-platform/packages/kafka-sdk/contracts"
	"github.com/Abdex1/FamGo-platform/packages/kafka-sdk/governance"
	"github.com/Abdex1/FamGo-platform/packages/kafka-sdk/producer"
)

func main() {
	logger, _ := zap.NewProduction()

	cfg := config.Config{
		Brokers: []string{"localhost:9092"},
		BatchSize: 100,
		BatchTimeout: time.Second,
	}

	p := producer.New(cfg, logger)

	event := contracts.EventEnvelope{
		EventID: uuid.NewString(),
		EventType: "ride.created",
		EventVersion: "v1",
		TraceID: uuid.NewString(),
		CorrelationID: uuid.NewString(),
		OccurredAt: time.Now().UTC(),
		Payload: map[string]any{
			"ride_id": "ride_123",
			"rider_id": "user_456",
		},
	}

	_ = p.Publish(
		context.Background(),
		governance.TopicRideCreated,
		"ride_123",
		event,
	)
}
