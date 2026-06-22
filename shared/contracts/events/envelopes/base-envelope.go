/*
# STEP 2 — CREATE ENTERPRISE EVENT ENVELOPE


shared/contracts/events/envelopes/base-envelope.go

*/
package envelopes

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type EventEnvelope struct {
	EventID        string          `json:"event_id"`
	EventType      string          `json:"event_type"`
	EventVersion   string          `json:"event_version"`
	OccurredAt     time.Time       `json:"occurred_at"`

	TraceID        string          `json:"trace_id"`
	SpanID         string          `json:"span_id"`
	CorrelationID  string          `json:"correlation_id"`
	CausationID    string          `json:"causation_id"`

	ServiceName    string          `json:"service_name"`
	Environment    string          `json:"environment"`

	UserID         string          `json:"user_id,omitempty"`
	DeviceID       string          `json:"device_id,omitempty"`
	SessionID      string          `json:"session_id,omitempty"`

	IdempotencyKey string          `json:"idempotency_key"`

	SchemaVersion  string          `json:"schema_version"`

	Payload        json.RawMessage `json:"payload"`
	Metadata       map[string]any  `json:"metadata,omitempty"`
}

func NewEnvelope(
	eventType string,
	version string,
	service string,
	payload any,
) (*EventEnvelope, error) {

	raw, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return &EventEnvelope{
		EventID:       uuid.NewString(),
		EventType:     eventType,
		EventVersion:  version,
		SchemaVersion: version,

		ServiceName:   service,
		OccurredAt:    time.Now().UTC(),

		Payload:       raw,
		Metadata:      map[string]any{},
	}, nil
}
