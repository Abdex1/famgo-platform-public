# 📊 OBSERVABILITY PATTERNS
## Extracted from uber-master, Enhanced for FamGo

**Status:** Pattern 8/8

---

## Metrics Pattern (Prometheus)

```go
import "github.com/prometheus/client_golang/prometheus"

var (
    requestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "http_request_duration_seconds",
            Buckets: []float64{.001, .01, .1, .5, 1, 2.5, 5, 10},
        },
        []string{"service", "method", "path", "status"},
    )
    
    requestCount = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
        },
        []string{"service", "method", "path", "status"},
    )
)

// In handler:
start := time.Now()
defer func() {
    duration := time.Since(start).Seconds()
    requestDuration.WithLabelValues(
        "driver-service", r.Method, r.URL.Path, fmt.Sprintf("%d", statusCode),
    ).Observe(duration)
    requestCount.WithLabelValues(
        "driver-service", r.Method, r.URL.Path, fmt.Sprintf("%d", statusCode),
    ).Inc()
}()
```

## Structured Logging

```go
import "github.com/go-kit/log"

logger.Log(
    "msg", "Creating driver",
    "driver_id", driverID,
    "email", email,
    "trace_id", traceID,
)

logger.Log(
    "msg", "Driver created",
    "driver_id", driver.ID,
    "status", driver.Status,
    "duration_ms", time.Since(start).Milliseconds(),
)
```

## Distributed Tracing

```go
import "go.opentelemetry.io/otel"

ctx, span := otel.Tracer("driver-service").Start(
    r.Context(), 
    "CreateDriver",
)
defer span.End()

span.SetAttributes(
    attribute.String("driver.email", email),
    attribute.String("driver.name", name),
)

if err != nil {
    span.RecordError(err)
    span.SetAttributes(attribute.Bool("error", true))
    return
}
```

---

## All 8 Patterns Complete

```
✅ Pattern 1: HTTP Handler Patterns
✅ Pattern 2: Service Bootstrap
✅ Pattern 3: Kafka Patterns
✅ Pattern 4: State Machines
✅ Pattern 5: Data Access
✅ Pattern 6: Payment Gateway
✅ Pattern 7: Testing
✅ Pattern 8: Observability
```

**Status:** ALL PATTERNS DOCUMENTED AND READY FOR USE

---
