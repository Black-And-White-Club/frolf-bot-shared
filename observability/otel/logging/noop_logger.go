package loggerfrolfbot

import (
	"context"
	"log/slog"
)

// NoOpLogger is a logger that does nothing. Useful for unit tests.
type NoOpLogger struct{}

// NewNoOpLogger creates a new instance of NoOpLogger.
func NewNoOpLogger() *NoOpLogger {
	return &NoOpLogger{}
}

// Debug does nothing.
func (n *NoOpLogger) Debug(msg string, attrs ...slog.Attr) {}

// Info does nothing.
func (n *NoOpLogger) Info(msg string, attrs ...slog.Attr) {}

// Warn does nothing.
func (n *NoOpLogger) Warn(msg string, attrs ...slog.Attr) {}

// Error does nothing.
func (n *NoOpLogger) Error(msg string, attrs ...slog.Attr) {}

// Close does nothing.
func (n *NoOpLogger) Close() {}

// WithContext returns the same NoOpLogger, as it does not log anything.
func (n *NoOpLogger) WithContext(ctx context.Context) *NoOpLogger {
	return n
}
