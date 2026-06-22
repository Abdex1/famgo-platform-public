package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

// Logger is the shared structured logging interface used by FamGo services.
type Logger interface {
	Info(msg string, fields map[string]interface{})
	Warn(msg string, fields map[string]interface{})
	Error(msg string, fields map[string]interface{})
	Debug(msg string, fields map[string]interface{})
}

type stdLogger struct {
	level int
	l     *log.Logger
}

const (
	levelDebug = iota
	levelInfo
	levelWarn
	levelError
)

func parseLevel(level string) int {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug":
		return levelDebug
	case "warn", "warning":
		return levelWarn
	case "error":
		return levelError
	default:
		return levelInfo
	}
}

// New creates a stdout logger at the requested level.
func New(level string) Logger {
	return &stdLogger{
		level: parseLevel(level),
		l:     log.New(os.Stdout, "", log.LstdFlags|log.LUTC),
	}
}

func (l *stdLogger) Info(msg string, fields map[string]interface{}) {
	if l.level <= levelInfo {
		l.write("INFO", msg, fields)
	}
}

func (l *stdLogger) Warn(msg string, fields map[string]interface{}) {
	if l.level <= levelWarn {
		l.write("WARN", msg, fields)
	}
}

func (l *stdLogger) Error(msg string, fields map[string]interface{}) {
	if l.level <= levelError {
		l.write("ERROR", msg, fields)
	}
}

func (l *stdLogger) Debug(msg string, fields map[string]interface{}) {
	if l.level <= levelDebug {
		l.write("DEBUG", msg, fields)
	}
}

func (l *stdLogger) write(level, msg string, fields map[string]interface{}) {
	if len(fields) == 0 {
		l.l.Printf("[%s] %s", level, msg)
		return
	}

	payload, err := json.Marshal(fields)
	if err != nil {
		l.l.Printf("[%s] %s fields=%v", level, msg, fields)
		return
	}

	l.l.Printf("[%s] %s %s", level, msg, fmt.Sprintf("%s", payload))
}
