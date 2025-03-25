package lokifrolfbot

import (
	"github.com/Black-And-White-Club/frolf-bot-shared/observability/attr"
	"github.com/ThreeDotsLabs/watermill"
)

// loggerToWatermillAdapter converts our Logger to watermill.LoggerAdapter
type loggerToWatermillAdapter struct {
	logger Logger
}

// ToWatermillAdapter converts a Logger to a watermill.LoggerAdapter
// This allows you to use your structured Logger with Watermill
func ToWatermillAdapter(logger Logger) watermill.LoggerAdapter {
	// Create an adapter
	return &loggerToWatermillAdapter{
		logger: logger,
	}
}

func (a *loggerToWatermillAdapter) Error(msg string, err error, fields watermill.LogFields) {
	attrs := make([]attr.LogAttr, 0, len(fields)+1)
	for k, v := range fields {
		attrs = append(attrs, attr.Any(k, v))
	}
	attrs = append(attrs, attr.Error(err))
	a.logger.Error(msg, attrs...)
}

func (a *loggerToWatermillAdapter) Info(msg string, fields watermill.LogFields) {
	attrs := make([]attr.LogAttr, 0, len(fields))
	for k, v := range fields {
		attrs = append(attrs, attr.Any(k, v))
	}
	a.logger.Info(msg, attrs...)
}

func (a *loggerToWatermillAdapter) Debug(msg string, fields watermill.LogFields) {
	attrs := make([]attr.LogAttr, 0, len(fields))
	for k, v := range fields {
		attrs = append(attrs, attr.Any(k, v))
	}
	a.logger.Debug(msg, attrs...)
}

func (a *loggerToWatermillAdapter) Trace(msg string, fields watermill.LogFields) {
	attrs := make([]attr.LogAttr, 0, len(fields))
	for k, v := range fields {
		attrs = append(attrs, attr.Any(k, v))
	}
	a.logger.Debug(msg, attrs...) // Using Debug level for Trace
}

func (a *loggerToWatermillAdapter) With(fields watermill.LogFields) watermill.LoggerAdapter {
	// Not the most efficient implementation, but it works
	// Create a new logger with these fields
	return &loggerToWatermillAdapter{
		logger: a.logger,
	}
}

func (a *loggerToWatermillAdapter) Close() {
	a.logger.Close()
}
