// Package logger provides a structured zerolog-based logger with context support.
package logger

import (
	"context"
	"io"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

type contextKey struct{}

// New creates a configured zerolog.Logger based on the log level string.
func New(level string) zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339

	var lvl zerolog.Level
	switch strings.ToLower(level) {
	case "debug":
		lvl = zerolog.DebugLevel
	case "warn":
		lvl = zerolog.WarnLevel
	case "error":
		lvl = zerolog.ErrorLevel
	default:
		lvl = zerolog.InfoLevel
	}

	var output io.Writer = os.Stdout
	if os.Getenv("APP_ENV") == "development" {
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "15:04:05",
		}
	}

	return zerolog.New(output).
		Level(lvl).
		With().
		Timestamp().
		Caller().
		Logger()
}

// WithContext returns a copy of ctx with the logger attached.
func WithContext(ctx context.Context, l zerolog.Logger) context.Context {
	return context.WithValue(ctx, contextKey{}, l)
}

// FromContext retrieves the logger from the context, or returns a default logger.
func FromContext(ctx context.Context) zerolog.Logger {
	if l, ok := ctx.Value(contextKey{}).(zerolog.Logger); ok {
		return l
	}
	return zerolog.Nop()
}
