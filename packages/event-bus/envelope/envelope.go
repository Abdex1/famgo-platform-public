
// 6. BUILD THE CORE EVENT ENVELOPE

// packages/event-bus/envelope/envelope.go

package envelope

import "time"

type EventEnvelope struct {
	EventID        string                 `json:"event_id" validate:"required"`
	EventType      string                 `json:"event_type" validate:"required"`
	EventVersion   string                 `json:"event_version" validate:"required"`

	TraceID        string                 `json:"trace_id"`
	SpanID         string                 `json:"span_id"`

	CorrelationID  string                 `json:"correlation_id"`
	CausationID    string                 `json:"causation_id"`

	RequestID      string                 `json:"request_id"`

	Service        string                 `json:"service"`
	Domain         string                 `json:"domain"`

	Environment    string                 `json:"environment"`

	PartitionKey   string                 `json:"partition_key"`

	IdempotencyKey string                 `json:"idempotency_key"`

	OccurredAt     time.Time              `json:"occurred_at"`

	Headers        map[string]string      `json:"headers"`

	Payload        any                    `json:"payload"`
}
