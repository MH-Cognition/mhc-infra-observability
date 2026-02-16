package logging

import (
	"context"
	"log/slog"
	"os"
)

// Level represents log level.
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelError
)

// Logger wraps slog for structured logging with trace awareness.
type Logger struct {
	inner *slog.Logger
}

// New creates a Logger with the given level. Uses JSON handler for production.
func New(level Level) *Logger {
	var slogLevel slog.Level
	switch level {
	case LevelDebug:
		slogLevel = slog.LevelDebug
	case LevelInfo:
		slogLevel = slog.LevelInfo
	case LevelError:
		slogLevel = slog.LevelError
	default:
		slogLevel = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{Level: slogLevel}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	return &Logger{inner: slog.New(handler)}
}

// WithTrace adds trace_id and span_id to log attributes when present in ctx.
func (l *Logger) WithTrace(ctx context.Context) *slog.Logger {
	tc := FromContext(ctx)
	if tc.TraceID == "" && tc.SpanID == "" {
		return l.inner
	}
	args := make([]any, 0, 4)
	if tc.TraceID != "" {
		args = append(args, "trace_id", tc.TraceID)
	}
	if tc.SpanID != "" {
		args = append(args, "span_id", tc.SpanID)
	}
	return l.inner.With(args...)
}

// Info logs at info level with trace context from ctx.
func (l *Logger) Info(ctx context.Context, msg string, args ...any) {
	l.WithTrace(ctx).Info(msg, args...)
}

// Error logs at error level with trace context from ctx.
func (l *Logger) Error(ctx context.Context, msg string, args ...any) {
	l.WithTrace(ctx).Error(msg, args...)
}

// Debug logs at debug level with trace context from ctx.
func (l *Logger) Debug(ctx context.Context, msg string, args ...any) {
	l.WithTrace(ctx).Debug(msg, args...)
}

// ParseLevel converts string to Level. Defaults to Info for unknown values.
func ParseLevel(s string) Level {
	switch s {
	case "debug", "DEBUG":
		return LevelDebug
	case "info", "INFO":
		return LevelInfo
	case "error", "ERROR":
		return LevelError
	default:
		return LevelInfo
	}
}

// LoggerFromEnv creates a Logger using LOG_LEVEL env var.
func LoggerFromEnv() *Logger {
	levelStr := os.Getenv("LOG_LEVEL")
	return New(ParseLevel(levelStr))
}
