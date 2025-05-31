package loggerfrolfbot

import (
	"context"
	"log/slog"

	"github.com/Black-And-White-Club/frolf-bot-shared/observability/attr"
	"github.com/ThreeDotsLabs/watermill"
)

// loggerToWatermillAdapter converts our Logger to watermill.LoggerAdapter
type loggerToWatermillAdapter struct {
	logger *slog.Logger
}

// ToWatermillAdapter converts a Logger to a watermill.LoggerAdapter
func ToWatermillAdapter(logger *slog.Logger) watermill.LoggerAdapter {
	return &loggerToWatermillAdapter{
		logger: logger,
	}
}

func (a *loggerToWatermillAdapter) Error(msg string, err error, fields watermill.LogFields) {
	attrs := convertFieldsToAttrs(fields)
	attrs = append(attrs, attr.Error(err))
	a.logger.LogAttrs(context.Background(), slog.LevelError, msg, attrs...)
}

func (a *loggerToWatermillAdapter) Info(msg string, fields watermill.LogFields) {
	attrs := convertFieldsToAttrs(fields)
	a.logger.LogAttrs(context.Background(), slog.LevelInfo, msg, attrs...)
}

func (a *loggerToWatermillAdapter) Debug(msg string, fields watermill.LogFields) {
	attrs := convertFieldsToAttrs(fields)
	a.logger.LogAttrs(context.Background(), slog.LevelDebug, msg, attrs...)
}

func (a *loggerToWatermillAdapter) Trace(msg string, fields watermill.LogFields) {
	attrs := convertFieldsToAttrs(fields)
	a.logger.LogAttrs(context.Background(), slog.LevelDebug, msg, attrs...) // Use LevelDebug instead of LevelTrace
}

func (a *loggerToWatermillAdapter) With(fields watermill.LogFields) watermill.LoggerAdapter {
	attrs := convertFieldsToAttrs(fields)

	// Convert slog.Attr to []any
	var anyAttrs []any
	for _, attr := range attrs {
		anyAttrs = append(anyAttrs, attr) // attr is of type slog.Attr
	}

	return &loggerToWatermillAdapter{
		logger: a.logger.With(anyAttrs...), // Use the converted slice
	}
}
func (s *loggerToWatermillAdapter) Close() {}

// Helper function to convert watermill.LogFields to slog.Attr
func convertFieldsToAttrs(fields watermill.LogFields) []slog.Attr {
	attrs := make([]slog.Attr, 0, len(fields))
	for k, v := range fields {
		attrs = append(attrs, attr.Any(k, v)) // Assuming attr.Any returns slog.Attr
	}
	return attrs
}
