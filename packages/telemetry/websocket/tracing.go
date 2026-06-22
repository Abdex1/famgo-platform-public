// STEP 10 — WEBSOCKET TELEMETRY

 
// packages/telemetry/go/websocket/tracing.go
  
package websocket

import (
    "context"

    "go.opentelemetry.io/otel"
)

func StartSocketSpan(ctx context.Context, event string) (context.Context, func()) {
    tracer := otel.Tracer("websocket")

    ctx, span := tracer.Start(ctx, event)

    return ctx, func() {
        span.End()
    }
}
