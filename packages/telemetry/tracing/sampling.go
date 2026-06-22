/*
# STEP 10 — ENTERPRISE TRACE SAMPLING

# WHY THIS MATTERS

Without sampling:

* telemetry costs explode
* Jaeger crashes
* Tempo storage explodes

packages/telemetry/tracing/sampling.go
*/
package tracing

import (
	"go.opentelemetry.io/otel/sdk/trace"
)

func BuildSampler() trace.Sampler {
	return trace.ParentBased(
		trace.TraceIDRatioBased(0.1),
	)
}
