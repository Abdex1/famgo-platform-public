package observability

import (
	"context"
	"os"

	telemetryconfig "github.com/Abdex1/FamGo-platform/packages/telemetry/config"
	"github.com/Abdex1/FamGo-platform/packages/telemetry/sdk"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace/noop"
)

// Bootstrap initializes OpenTelemetry tracing when OTLP_ENDPOINT is configured.
func Bootstrap(ctx context.Context, serviceName, environment string) (*sdk.SDK, error) {
	endpoint := os.Getenv("OTLP_ENDPOINT")
	if endpoint == "" {
		otel.SetTracerProvider(noop.NewTracerProvider())
		return nil, nil
	}

	return sdk.Bootstrap(ctx, telemetryconfig.Config{
		ServiceName:    serviceName,
		Environment:    environment,
		OTLPEndpoint:   endpoint,
		TracingEnabled: true,
		LoggingEnabled: true,
	})
}
