package observability

import "context"

// MetricsRegistry holds service metrics hooks. Sprint wiring uses a no-op registry.
type MetricsRegistry struct {
	Service string
}

// Tracer wraps OpenTelemetry setup. Sprint wiring uses a no-op tracer.
type Tracer struct {
	Service string
}

// InitMetrics registers baseline metrics for a service.
func InitMetrics(service string) (*MetricsRegistry, error) {
	return &MetricsRegistry{Service: service}, nil
}

// InitTracer configures distributed tracing for a service.
func InitTracer(service, endpoint string) (*Tracer, error) {
	_ = endpoint
	return &Tracer{Service: service}, nil
}

// Shutdown stops tracer exporters.
func (t *Tracer) Shutdown(ctx context.Context) error {
	_ = ctx
	return nil
}
