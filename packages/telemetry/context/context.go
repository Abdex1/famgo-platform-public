
// STEP 3 — CONTEXT PACKAGE
//packages/telemetry/go/context/context.go

package context

import "context"

type ctxKey string

const (
    TraceIDKey       ctxKey = "trace_id"
    CorrelationIDKey ctxKey = "correlation_id"
    RequestIDKey     ctxKey = "request_id"
    UserIDKey        ctxKey = "user_id"
    DeviceIDKey      ctxKey = "device_id"
    RideIDKey        ctxKey = "ride_id"
)

func SetValue(ctx context.Context, key ctxKey, value string) context.Context {
    return context.WithValue(ctx, key, value)
}

func GetValue(ctx context.Context, key ctxKey) string {
    val, ok := ctx.Value(key).(string)
    if !ok {
        return ""
    }

    return val
}
