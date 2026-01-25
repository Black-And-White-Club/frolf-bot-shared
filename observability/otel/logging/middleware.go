package loggerfrolfbot

import (
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"go.opentelemetry.io/otel/trace"
)

// LoggingMiddleware creates a middleware for logging Watermill messages with trace correlation.
// This middleware extracts trace_id and span_id from the message context for Grafana Trace ↔ Logs correlation.
func LoggingMiddleware(logger watermill.LoggerAdapter) message.HandlerMiddleware {
	return func(next message.HandlerFunc) message.HandlerFunc {
		return func(msg *message.Message) ([]*message.Message, error) {
			startTime := time.Now()

			topic := msg.Metadata.Get("topic")
			correlationID := msg.Metadata.Get(middleware.CorrelationIDMetadataKey)
			if correlationID == "" {
				correlationID = "unknown"
			}

			// Base fields for all log entries
			baseFields := watermill.LogFields{
				"correlation_id": correlationID,
				"topic":          topic,
				"message_id":     msg.UUID,
				"handler":        msg.Metadata.Get("handler_name"),
				"domain":         msg.Metadata.Get("domain"),
			}

			// Extract trace context for Grafana correlation
			spanCtx := trace.SpanContextFromContext(msg.Context())
			if spanCtx.IsValid() {
				baseFields["trace_id"] = spanCtx.TraceID().String()
				baseFields["span_id"] = spanCtx.SpanID().String()
			}

			ctxLogger := logger.With(baseFields)

			ctxLogger.Info("Received message", nil)

			producedMessages, err := next(msg)
			duration := time.Since(startTime)

			fields := watermill.LogFields{
				"duration_ms": duration.Milliseconds(),
			}

			if err != nil {
				ctxLogger.Error("Error processing message", err, fields)
			} else {
				ctxLogger.Info("Message processed successfully", fields)
			}

			return producedMessages, err
		}
	}
}
