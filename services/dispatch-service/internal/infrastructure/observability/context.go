package observability

import (
	"context"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"

	telemetryctx "github.com/Abdex1/FamGo-platform/packages/telemetry/context"
)

const (
	headerTraceID       = "trace_id"
	headerCorrelationID = "correlation_id"
	headerRequestID     = "request_id"
)

func ContextFromGRPC(ctx context.Context) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return EnrichContext(ctx, "", "", "")
	}
	return EnrichContext(ctx,
		firstMetadata(md, headerTraceID),
		firstMetadata(md, headerCorrelationID),
		firstMetadata(md, headerRequestID),
	)
}

func EnrichContext(ctx context.Context, traceID, correlationID, requestID string) context.Context {
	if traceID == "" {
		span := trace.SpanFromContext(ctx)
		if span.SpanContext().HasTraceID() {
			traceID = span.SpanContext().TraceID().String()
		} else {
			traceID = uuid.NewString()
		}
	}
	if correlationID == "" {
		correlationID = traceID
	}
	if requestID == "" {
		requestID = uuid.NewString()
	}

	ctx = telemetryctx.SetValue(ctx, telemetryctx.TraceIDKey, traceID)
	ctx = telemetryctx.SetValue(ctx, telemetryctx.CorrelationIDKey, correlationID)
	ctx = telemetryctx.SetValue(ctx, telemetryctx.RequestIDKey, requestID)
	return ctx
}

func IDsFromContext(ctx context.Context) (traceID, correlationID, requestID string) {
	return telemetryctx.GetValue(ctx, telemetryctx.TraceIDKey),
		telemetryctx.GetValue(ctx, telemetryctx.CorrelationIDKey),
		telemetryctx.GetValue(ctx, telemetryctx.RequestIDKey)
}

func firstMetadata(md metadata.MD, key string) string {
	values := md.Get(key)
	if len(values) == 0 {
		return ""
	}
	return values[0]
}
