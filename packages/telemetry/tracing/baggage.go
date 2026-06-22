/*
# STEP 11 — ENTERPRISE BAGGAGE PROPAGATION


packages/telemetry/tracing/baggage.go
*/
package tracing

import (
	"go.opentelemetry.io/otel/propagation"
)

func BuildPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}
