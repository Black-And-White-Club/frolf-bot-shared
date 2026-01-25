package loggerfrolfbot

import (
	"context"
	"log/slog"

	"github.com/Black-And-White-Club/frolf-bot-shared/observability/attr"
	"github.com/ThreeDotsLabs/watermill"
)

// loggerToWatermillAdapter converts our Logger to watermill.LoggerAdapter.
// It preserves trace context for Grafana Trace ↔ Logs correlation.
type loggerToWatermillAdapter struct {
	logger *slog.Logger
	ctx    context.Context
}

// ToWatermillAdapter converts a Logger to a watermill.LoggerAdapter.
func ToWatermillAdapter(logger *slog.Logger) watermill.LoggerAdapter {
	return &loggerToWatermillAdapter{
		logger: logger,
		ctx:    context.Background(),
	}
}

// ToWatermillAdapterWithContext converts a Logger to a watermill.LoggerAdapter with trace context.
// Use this when you have a context with trace information (e.g., from a message handler).
func ToWatermillAdapterWithContext(ctx context.Context, logger *slog.Logger) watermill.LoggerAdapter {
	return &loggerToWatermillAdapter{
		logger: logger,
		ctx:    ctx,
	}
}

func (a *loggerToWatermillAdapter) Error(msg string, err error, fields watermill.LogFields) {
	attrs := convertFieldsToAttrs(fields)
	attrs = append(attrs, attr.Error(err))
	a.logger.LogAttrs(a.ctx, slog.LevelError, msg, attrs...)
}

func (a *loggerToWatermillAdapter) Info(msg string, fields watermill.LogFields) {
	attrs := convertFieldsToAttrs(fields)
	a.logger.LogAttrs(a.ctx, slog.LevelInfo, msg, attrs...)
}

func (a *loggerToWatermillAdapter) Debug(msg string, fields watermill.LogFields) {
	attrs := convertFieldsToAttrs(fields)
	a.logger.LogAttrs(a.ctx, slog.LevelDebug, msg, attrs...)
}

func (a *loggerToWatermillAdapter) Trace(msg string, fields watermill.LogFields) {
	attrs := convertFieldsToAttrs(fields)
	a.logger.LogAttrs(a.ctx, slog.LevelDebug, msg, attrs...) // Use LevelDebug instead of LevelTrace
}

func (a *loggerToWatermillAdapter) With(fields watermill.LogFields) watermill.LoggerAdapter {
	attrs := convertFieldsToAttrs(fields)

	// Convert slog.Attr to []any
	var anyAttrs []any
	for _, attr := range attrs {
		anyAttrs = append(anyAttrs, attr)
	}

	return &loggerToWatermillAdapter{
		logger: a.logger.With(anyAttrs...),
		ctx:    a.ctx,
	}
}

func (a *loggerToWatermillAdapter) Close() {}

// Helper function to convert watermill.LogFields to slog.Attr
func convertFieldsToAttrs(fields watermill.LogFields) []slog.Attr {
	attrs := make([]slog.Attr, 0, len(fields))
	for k, v := range fields {
		attrs = append(attrs, attr.Any(k, v)) // Assuming attr.Any returns slog.Attr
	}
	return attrs
}
