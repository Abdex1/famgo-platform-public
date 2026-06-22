
// 5. EVENT GOVERNANCE PLATFORM

//shared/contracts/events/common/v1/envelope.go

package v1

import "time"

type EventEnvelope[T any] struct {
    EventID        string            `json:"event_id"`
    EventType      string            `json:"event_type"`
    EventVersion   string            `json:"event_version"`
    CorrelationID  string            `json:"correlation_id"`
    CausationID    string            `json:"causation_id"`
    TenantID       string            `json:"tenant_id"`
    Source         string            `json:"source"`
    Timestamp      time.Time         `json:"timestamp"`
    Metadata       map[string]string `json:"metadata"`
    Payload        T                 `json:"payload"`
}
