// 11. CREATE TRACE PROPAGATION

// packages/event-bus/tracing/propagation.go

package tracing

const (
	HeaderTraceID       = "trace_id"
	HeaderSpanID        = "span_id"
	HeaderCorrelationID = "correlation_id"
	HeaderRequestID     = "request_id"
)
