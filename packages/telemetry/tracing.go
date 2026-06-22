/*
# PHASE 3 — OPENTELEMETRY DISTRIBUTED TRACING

# =========================================================

# STEP 10 — TRACE PROPAGATION


packages/telemetry/tracing.go
*/
package telemetry

import (
	"context"

	"go.opentelemetry.io/otel"
)

func StartSpan(
	ctx context.Context,
	name string,
) (context.Context, func()) {

	tracer := otel.Tracer("famgo")

	ctx, span := tracer.Start(ctx, name)

	return ctx, func() {
		span.End()
	}
}
