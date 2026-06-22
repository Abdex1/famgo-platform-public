
// STEP 12 — CREATE OBSERVABILITY HOOKS

//packages/kafka-sdk/observability/tracing.go

package observability

import (
	"go.opentelemetry.io/otel"
)

var Tracer = otel.Tracer("famgo-kafka-sdk")
