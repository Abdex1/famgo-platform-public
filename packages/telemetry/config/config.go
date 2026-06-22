// packages/telemetry/config/config.go
package config

import "time"

type Config struct {
	ServiceName            string
	ServiceVersion         string
	Environment            string
	OTLPEndpoint           string
	MetricsEnabled         bool
	TracingEnabled         bool
	LoggingEnabled         bool
	SamplingRatio          float64
	MetricsInterval        time.Duration
	ExportTimeout          time.Duration
	EnableRuntimeMetrics   bool
	EnableHostMetrics      bool
	EnableKafkaTracing     bool
	EnableHTTPTracing      bool
	EnableGRPCTracing      bool
	EnableWebsocketTracing bool
}
