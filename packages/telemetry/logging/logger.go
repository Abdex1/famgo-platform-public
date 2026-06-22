
// STEP 11 — CREATE ENTERPRISE LOGGER
// packages/telemetry/logging/logger.go

package logging

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

func New(service string) zerolog.Logger {
	return zerolog.New(os.Stdout).
		With().
		Timestamp().
		Str("service", service).
		Logger().
		Level(zerolog.InfoLevel).
		Output(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		})
}
