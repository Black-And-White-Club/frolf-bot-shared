package loggerfrolfbot

import (
	"context"
	"fmt"
	"log/slog"
)

// TestHandler is a custom handler for capturing log messages.
type TestHandler struct {
	Messages []string
}

// NewTestHandler creates a new instance of TestHandler.
func NewTestHandler() *TestHandler {
	return &TestHandler{
		Messages: []string{},
	}
}

// Handle logs a message and captures it.
func (h *TestHandler) Handle(ctx context.Context, r slog.Record) error {
	h.Messages = append(h.Messages, fmt.Sprintf("%s: %s", r.Level.String(), r.Message))
	return nil
}

// WithGroup is a no-op for this handler.
func (h *TestHandler) WithGroup(name string) slog.Handler {
	return h
}

// WithAttrs is a no-op for this handler.
func (h *TestHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

// Enabled checks if the log level is enabled for logging.
func (h *TestHandler) Enabled(ctx context.Context, level slog.Level) bool {
	// You can customize this logic based on your needs.
	// For example, you might want to log everything:
	return true
}
