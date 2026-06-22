/*
# PHASE 13 — OTEL INTEGRATION

Inside your telemetry SDK:

Add tracing middleware:

```txt
packages/telemetry/http_middleware.go
```

Install:

```powershell
go get go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp
*/
package telemetry

import (
    "net/http"

    "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func Middleware(service string) func(http.Handler) http.Handler {
    return otelhttp.NewMiddleware(service)
}
