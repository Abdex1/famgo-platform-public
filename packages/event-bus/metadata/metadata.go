// 10. BUILD EVENT METADATA CONTRACTS

// packages/event-bus/metadata/metadata.go

package metadata

type EventMetadata struct {
	TraceID        string
	SpanID         string
	CorrelationID  string
	CausationID    string
	RequestID      string
	IdempotencyKey string
	Service        string
	Environment    string
}
