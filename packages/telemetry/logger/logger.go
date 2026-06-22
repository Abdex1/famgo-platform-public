
// STEP 4 — LOGGER FOUNDATION
// packages/telemetry/go/logger/logger.go

package logger

import (
    "context"
    "os"
    "time"

    "github.com/rs/zerolog"
)

type Logger struct {
    z zerolog.Logger
}

func New(service string, env string) *Logger {
    logger := zerolog.New(os.Stdout).
        With().
        Timestamp().
        Str("service", service).
        Str("environment", env).
        Logger()

    zerolog.TimeFieldFormat = time.RFC3339Nano

    return &Logger{z: logger}
}

func (l *Logger) Info(ctx context.Context, message string, fields ...Field) {
    evt := l.z.Info()

    for _, field := range fields {
        field.Apply(evt)
    }

    evt.Msg(message)
}

func (l *Logger) Error(ctx context.Context, err error, message string, fields ...Field) {
    evt := l.z.Error().Err(err)

    for _, field := range fields {
        field.Apply(evt)
    }

    evt.Msg(message)
}
