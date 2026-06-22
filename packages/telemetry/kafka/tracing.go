
// STEP 9 — KAFKA TRACING

 
//packages/telemetry/go/kafka/tracing.go
  
package kafka

import (
    "context"

    "go.opentelemetry.io/otel"
)

func StartProducerSpan(ctx context.Context, topic string) (context.Context, func()) {
    tracer := otel.Tracer("kafka-producer")

    ctx, span := tracer.Start(ctx, "kafka.produce")

    return ctx, func() {
        span.End()
    }
}
