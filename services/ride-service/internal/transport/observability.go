// services/ride-service/internal/transport/observability.go
// Observability: Prometheus Metrics, Jaeger Traces, Structured Logging

package transport

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.uber.org/zap"

	"github.com/Abdex1/FamGo-platform/packages/telemetry"
)

// MetricsCollector collects Prometheus metrics
type MetricsCollector struct {
	requestCount      prometheus.CounterVec
	requestDuration   prometheus.HistogramVec
	requestErrors     prometheus.CounterVec
	ridesCreated      prometheus.Counter
	ridesCompleted    prometheus.Counter
	ridesCancelled    prometheus.Counter
	activeRides       prometheus.Gauge
	gRPCCallCount     prometheus.CounterVec
	gRPCCallDuration  prometheus.HistogramVec
	circuitBreakerStatus prometheus.GaugeVec
}

// NewMetricsCollector creates and registers all metrics
func NewMetricsCollector(serviceName string) *MetricsCollector {
	mc := &MetricsCollector{
		requestCount: *promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_request_count",
				Help: "Total HTTP requests",
			},
			[]string{"method", "path", "status"},
		),

		requestDuration: *promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "HTTP request latency in seconds",
				Buckets: []float64{.001, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
			},
			[]string{"method", "path"},
		),

		requestErrors: *promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_request_errors_total",
				Help: "Total HTTP request errors",
			},
			[]string{"method", "path", "error_type"},
		),

		ridesCreated: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "rides_created_total",
				Help: "Total rides created",
			},
		),

		ridesCompleted: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "rides_completed_total",
				Help: "Total rides completed",
			},
		),

		ridesCancelled: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "rides_cancelled_total",
				Help: "Total rides cancelled",
			},
		),

		activeRides: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "active_rides",
				Help: "Current active rides",
			},
		),

		gRPCCallCount: *promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "grpc_call_count",
				Help: "Total gRPC calls",
			},
			[]string{"service", "method", "status"},
		),

		gRPCCallDuration: *promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "grpc_call_duration_seconds",
				Help:    "gRPC call latency in seconds",
				Buckets: []float64{.001, .01, .05, .1, .5, 1, 5},
			},
			[]string{"service", "method"},
		),

		circuitBreakerStatus: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "circuit_breaker_status",
				Help: "Circuit breaker status (0=closed, 1=open, 2=half-open)",
			},
			[]string{"service"},
		),
	}

	return mc
}

// RecordHTTPRequest records HTTP request metrics
func (mc *MetricsCollector) RecordHTTPRequest(method, path string, status int, duration time.Duration) {
	mc.requestCount.WithLabelValues(method, path, string(rune(status))).Inc()
	mc.requestDuration.WithLabelValues(method, path).Observe(duration.Seconds())

	if status >= 400 {
		errorType := "server_error"
		if status >= 400 && status < 500 {
			errorType = "client_error"
		}
		mc.requestErrors.WithLabelValues(method, path, errorType).Inc()
	}
}

// RecordRideCreated increments rides created counter
func (mc *MetricsCollector) RecordRideCreated() {
	mc.ridesCreated.Inc()
	mc.activeRides.Inc()
}

// RecordRideCompleted increments rides completed counter
func (mc *MetricsCollector) RecordRideCompleted() {
	mc.ridesCompleted.Inc()
	mc.activeRides.Dec()
}

// RecordRideCancelled increments rides cancelled counter
func (mc *MetricsCollector) RecordRideCancelled() {
	mc.ridesCancelled.Inc()
	mc.activeRides.Dec()
}

// RecordGRPCCall records gRPC call metrics
func (mc *MetricsCollector) RecordGRPCCall(service, method, status string, duration time.Duration) {
	mc.gRPCCallCount.WithLabelValues(service, method, status).Inc()
	mc.gRPCCallDuration.WithLabelValues(service, method).Observe(duration.Seconds())
}

// SetCircuitBreakerStatus updates circuit breaker status gauge
func (mc *MetricsCollector) SetCircuitBreakerStatus(service string, status int) {
	mc.circuitBreakerStatus.WithLabelValues(service).Set(float64(status))
}

// ==================== JAEGER TRACING ====================

// InitJaeger initializes Jaeger tracer
func InitJaeger(serviceName string, jaegerEndpoint string, logger *zap.Logger) (*tracesdk.TracerProvider, error) {
	exp, err := jaeger.New(
		jaeger.WithAgentHost(jaegerEndpoint),
	)
	if err != nil {
		logger.Error("failed to initialize jaeger exporter", zap.Error(err))
		return nil, err
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(serviceName),
			),
		),
	)

	otel.SetTracerProvider(tp)
	logger.Info("jaeger tracer initialized", zap.String("endpoint", jaegerEndpoint))

	return tp, nil
}

// ==================== STRUCTURED LOGGING ====================

// StructuredLogger wraps zap logger for structured logging
type StructuredLogger struct {
	logger *zap.Logger
}

func NewStructuredLogger(logger *zap.Logger) *StructuredLogger {
	return &StructuredLogger{logger: logger}
}

// LogOperation logs a domain operation (Create, Update, Delete, etc.)
func (sl *StructuredLogger) LogOperation(
	ctx context.Context,
	operation string,
	userID string,
	resource string,
	resourceID string,
	status string,
	duration time.Duration,
	error error,
) {
	fields := []zap.Field{
		zap.String("operation", operation),
		zap.String("user_id", userID),
		zap.String("resource", resource),
		zap.String("resource_id", resourceID),
		zap.String("status", status),
		zap.Duration("duration_ms", duration),
	}

	if traceID := extractTraceID(ctx); traceID != "" {
		fields = append(fields, zap.String("trace_id", traceID))
	}

	if error != nil {
		fields = append(fields, zap.Error(error))
		sl.logger.Error("operation failed", fields...)
	} else {
		sl.logger.Info("operation completed", fields...)
	}
}

// LogGRPCCall logs a gRPC external call
func (sl *StructuredLogger) LogGRPCCall(
	ctx context.Context,
	service string,
	method string,
	status string,
	duration time.Duration,
	error error,
) {
	fields := []zap.Field{
		zap.String("type", "grpc_call"),
		zap.String("service", service),
		zap.String("method", method),
		zap.String("status", status),
		zap.Duration("duration_ms", duration),
	}

	if traceID := extractTraceID(ctx); traceID != "" {
		fields = append(fields, zap.String("trace_id", traceID))
	}

	if error != nil {
		fields = append(fields, zap.Error(error))
		sl.logger.Error("grpc call failed", fields...)
	} else {
		sl.logger.Info("grpc call succeeded", fields...)
	}
}

// LogSecurityEvent logs security-related events (auth, authz, input validation)
func (sl *StructuredLogger) LogSecurityEvent(
	ctx context.Context,
	eventType string, // AUTH, AUTHZ, VALIDATION, RATE_LIMIT
	userID string,
	resource string,
	action string,
	result string, // ALLOW, DENY, INVALID
	details map[string]string,
) {
	fields := []zap.Field{
		zap.String("event_type", "security"),
		zap.String("security_event", eventType),
		zap.String("user_id", userID),
		zap.String("resource", resource),
		zap.String("action", action),
		zap.String("result", result),
	}

	if traceID := extractTraceID(ctx); traceID != "" {
		fields = append(fields, zap.String("trace_id", traceID))
	}

	for key, value := range details {
		fields = append(fields, zap.String(key, value))
	}

	sl.logger.Info("security event", fields...)
}

// extractTraceID extracts trace ID from context
func extractTraceID(ctx context.Context) string {
	if span := extractSpanContext(ctx); span != "" {
		return span
	}
	return ""
}

// extractSpanContext extracts span context from context (simplified)
func extractSpanContext(ctx context.Context) string {
	// In production, use opentelemetry.io/otel/trace
	// For now, return empty string
	return ""
}

// ==================== MIDDLEWARE ====================

// ObservabilityMiddleware wraps HTTP handlers with observability
type ObservabilityMiddleware struct {
	metrics *MetricsCollector
	logger  *StructuredLogger
}

func NewObservabilityMiddleware(metrics *MetricsCollector, logger *StructuredLogger) *ObservabilityMiddleware {
	return &ObservabilityMiddleware{
		metrics: metrics,
		logger:  logger,
	}
}

// WrapHandler wraps an HTTP handler with metrics and logging
func (om *ObservabilityMiddleware) WrapHandler(
	method string,
	path string,
	handler func(context.Context) (interface{}, int, error),
) func(context.Context) (interface{}, int, error) {
	return func(ctx context.Context) (interface{}, int, error) {
		start := time.Now()

		result, status, err := handler(ctx)

		duration := time.Since(start)
		om.metrics.RecordHTTPRequest(method, path, status, duration)
		om.logger.LogOperation(
			ctx,
			method+" "+path,
			"", // user_id would be extracted from auth context
			path,
			"",
			string(rune(status)),
			duration,
			err,
		)

		return result, status, err
	}
}
