package loggerfrolfbot

import (
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
)

// LokiLoggingMiddleware creates a middleware for logging Watermill messages
func LoggingMiddleware(logger watermill.LoggerAdapter) message.HandlerMiddleware {
	return func(next message.HandlerFunc) message.HandlerFunc {
		return func(msg *message.Message) ([]*message.Message, error) {
			startTime := time.Now()

			topic := msg.Metadata.Get("topic")
			correlationID := msg.Metadata.Get(middleware.CorrelationIDMetadataKey)
			if correlationID == "" {
				correlationID = "unknown"
			}

			ctxLogger := logger.With(watermill.LogFields{
				"correlation_id": correlationID,
				"topic":          topic,
				"discord_message_id":     msg.UUID,
				"handler":        msg.Metadata.Get("handler_name"),
				"domain":         msg.Metadata.Get("domain"),
			})

			ctxLogger.Info("Received message", watermill.LogFields{
				"duration": "0s",
			})

			producedMessages, err := next(msg)
			duration := time.Since(startTime)

			fields := watermill.LogFields{
				"duration": duration.String(),
				"topic":    topic,
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
