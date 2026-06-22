// STEP 4 — CREATE ENTERPRISE EVENT ENVELOPE

//packages/kafka-sdk/contracts/envelope.go


//## PURPOSE

//This becomes the MANDATORY enterprise event format.

//EVERY service MUST publish this envelope.


//# IMPLEMENTATION

package contracts

import "time"

type EventEnvelope struct {
	EventID        string            `json:"event_id"`
	EventType      string            `json:"event_type"`
	EventVersion   string            `json:"event_version"`
	TraceID        string            `json:"trace_id"`
	SpanID         string            `json:"span_id"`
	CorrelationID  string            `json:"correlation_id"`
	RequestID      string            `json:"request_id"`
	Producer       string            `json:"producer"`
	Environment    string            `json:"environment"`
	OccurredAt     time.Time         `json:"occurred_at"`
	PartitionKey   string            `json:"partition_key"`
	IdempotencyKey string            `json:"idempotency_key"`
	Headers        map[string]string `json:"headers"`
	Payload        any               `json:"payload"`
}
