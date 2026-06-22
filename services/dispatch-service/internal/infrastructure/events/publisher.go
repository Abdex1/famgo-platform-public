package events

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/Abdex1/FamGo-platform/packages/event-bus/envelope"
	dispatchcontracts "github.com/Abdex1/FamGo-platform/packages/event-bus/contracts/dispatch"
	"github.com/Abdex1/FamGo-platform/services/dispatch-service/internal/domain/ports"
)

const (
	EventMatchingStarted = "dispatch.matching.started.v1"
	EventDriverMatched   = "dispatch.driver.matched.v1"
	EventDriverAssigned  = "dispatch.driver.assigned.v1"
	EventMatchingFailed  = "dispatch.matching.failed.v1"
	EventMatchingExpired = "dispatch.matching.expired.v1"
)

// KafkaPublisher publishes dispatch events to Kafka topics.
type KafkaPublisher struct {
	publish func(ctx context.Context, topic, key string, payload []byte, headers map[string]string) error
	service string
	env     string
}

func NewKafkaPublisher(
	publish func(ctx context.Context, topic, key string, payload []byte, headers map[string]string) error,
	service, env string,
) *KafkaPublisher {
	return &KafkaPublisher{publish: publish, service: service, env: env}
}

func (p *KafkaPublisher) PublishMatchingStarted(ctx context.Context, event ports.MatchingStartedEvent) error {
	payload := dispatchcontracts.MatchingStarted{
		DispatchRequestID: event.DispatchRequestID,
		RideID:            event.RideID,
		RiderID:           event.RiderID,
		SearchRadiusKm:    event.SearchRadiusKm,
	}
	return p.publishEnvelope(ctx, EventMatchingStarted, event.RideID, event.TraceID, event.CorrelationID, event.RequestID, payload)
}

func (p *KafkaPublisher) PublishDriverMatched(ctx context.Context, event ports.DriverMatchedEvent) error {
	payload := dispatchcontracts.DriverMatched{
		DispatchRequestID: event.DispatchRequestID,
		RideID:            event.RideID,
		DriverID:          event.DriverID,
		ProposedDrivers:   event.ProposedDrivers,
		MatchScore:        event.MatchScore,
	}
	return p.publishEnvelope(ctx, EventDriverMatched, event.RideID, event.TraceID, event.CorrelationID, event.RequestID, payload)
}

func (p *KafkaPublisher) PublishDriverAssigned(ctx context.Context, event ports.DriverAssignedEvent) error {
	payload := dispatchcontracts.DriverAssigned{
		DispatchRequestID: event.DispatchRequestID,
		RideID:            event.RideID,
		DriverID:          event.DriverID,
	}
	return p.publishEnvelope(ctx, EventDriverAssigned, event.RideID, event.TraceID, event.CorrelationID, event.RequestID, payload)
}

func (p *KafkaPublisher) PublishMatchingFailed(ctx context.Context, event ports.MatchingFailedEvent) error {
	payload := dispatchcontracts.MatchingFailed{
		DispatchRequestID: event.DispatchRequestID,
		RideID:            event.RideID,
		Reason:            event.Reason,
		AttemptCount:      event.AttemptCount,
	}
	return p.publishEnvelope(ctx, EventMatchingFailed, event.RideID, event.TraceID, event.CorrelationID, event.RequestID, payload)
}

func (p *KafkaPublisher) PublishMatchingExpired(ctx context.Context, event ports.MatchingExpiredEvent) error {
	payload := dispatchcontracts.MatchingExpired{
		DispatchRequestID: event.DispatchRequestID,
		RideID:            event.RideID,
	}
	return p.publishEnvelope(ctx, EventMatchingExpired, event.RideID, event.TraceID, event.CorrelationID, event.RequestID, payload)
}

func (p *KafkaPublisher) publishEnvelope(
	ctx context.Context,
	eventType, partitionKey, traceID, correlationID, requestID string,
	payload any,
) error {
	if p.publish == nil {
		return nil
	}

	env := envelope.EventEnvelope{
		EventID:       uuid.NewString(),
		EventType:     eventType,
		EventVersion:  "v1",
		TraceID:       traceID,
		CorrelationID: correlationID,
		RequestID:     requestID,
		Service:       p.service,
		Domain:        "dispatch",
		Environment:   p.env,
		PartitionKey:  partitionKey,
		OccurredAt:    time.Now().UTC(),
		Payload:       payload,
	}

	body, err := json.Marshal(env)
	if err != nil {
		return fmt.Errorf("marshal event envelope: %w", err)
	}

	headers := map[string]string{
		"trace_id":       traceID,
		"correlation_id": correlationID,
		"request_id":     requestID,
		"event_type":     eventType,
	}
	return p.publish(ctx, eventType, partitionKey, body, headers)
}

// NoOpPublisher discards events (used in tests/local without Kafka).
type NoOpPublisher struct{}

func (NoOpPublisher) PublishMatchingStarted(context.Context, ports.MatchingStartedEvent) error { return nil }
func (NoOpPublisher) PublishDriverMatched(context.Context, ports.DriverMatchedEvent) error     { return nil }
func (NoOpPublisher) PublishDriverAssigned(context.Context, ports.DriverAssignedEvent) error   { return nil }
func (NoOpPublisher) PublishMatchingFailed(context.Context, ports.MatchingFailedEvent) error   { return nil }
func (NoOpPublisher) PublishMatchingExpired(context.Context, ports.MatchingExpiredEvent) error { return nil }
