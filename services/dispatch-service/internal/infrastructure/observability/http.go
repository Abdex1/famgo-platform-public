package observability

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// HTTPMiddleware enriches request context with trace/correlation/request IDs and
// instruments the handler with OpenTelemetry HTTP tracing.
func HTTPMiddleware(service string, next http.Handler) http.Handler {
	instrumented := otelhttp.NewHandler(next, service)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := EnrichContext(
			r.Context(),
			r.Header.Get(headerTraceID),
			r.Header.Get(headerCorrelationID),
			firstHeader(r, "X-Request-ID", headerRequestID),
		)
		traceID, correlationID, requestID := IDsFromContext(ctx)
		w.Header().Set(headerTraceID, traceID)
		w.Header().Set(headerCorrelationID, correlationID)
		w.Header().Set("X-Request-ID", requestID)
		instrumented.ServeHTTP(w, r.WithContext(ctx))
	})
}

func firstHeader(r *http.Request, keys ...string) string {
	for _, key := range keys {
		if v := r.Header.Get(key); v != "" {
			return v
		}
	}
	return ""
}
