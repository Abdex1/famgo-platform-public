
// STEP 6 — CREATE CENTRAL CONFIG
//packages/kafka-sdk/config/config.go

package config

import "time"

type Config struct {
	Brokers               []string
	ClientID              string
	GroupID               string
	RetryMax              int
	RetryBackoff          time.Duration
	DialTimeout           time.Duration
	ReadTimeout           time.Duration
	WriteTimeout          time.Duration
	EnableTLS             bool
	EnableSASL            bool
	Username              string
	Password              string
	DLQSuffix             string
	ConsumerConcurrency   int
	BatchSize             int
	BatchTimeout          time.Duration
	RequiredAcks          int
	CompressionCodec      string
	AutoTopicCreation     bool
	EnableIdempotency     bool
	EnableTracing         bool
	EnableMetrics         bool
}
