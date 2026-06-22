
// STEP 12 — CREATE TRACE MIDDLEWARE

// packages/telemetry/middleware/http.go

package middleware

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func HTTP(service string, next http.Handler) http.Handler {
	return otelhttp.NewHandler(next, service)
}
