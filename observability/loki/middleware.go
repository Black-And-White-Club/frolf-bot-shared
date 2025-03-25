package lokifrolfbot

import (
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
)

// LokiLoggingMiddleware creates a middleware for logging Watermill messages
func LokiLoggingMiddleware(logger watermill.LoggerAdapter) message.HandlerMiddleware {
	return func(next message.HandlerFunc) message.HandlerFunc {
		return func(msg *message.Message) ([]*message.Message, error) {
			startTime := time.Now()

			correlationID := msg.Metadata.Get(middleware.CorrelationIDMetadataKey)
			if correlationID == "" {
				correlationID = "unknown"
			}

			ctxLogger := logger.With(watermill.LogFields{
				"correlation_id": correlationID,
				"topic":          msg.Metadata.Get("topic"),
				"message_id":     msg.UUID,
				"handler":        msg.Metadata.Get("handler_name"),
				"domain":         msg.Metadata.Get("domain"),
			})

			ctxLogger.Info("Received message", watermill.LogFields{
				"topic": msg.Metadata.Get("topic"),
			})

			producedMessages, err := next(msg)
			duration := time.Since(startTime)

			if err != nil {
				ctxLogger.Error("Error processing message", err, watermill.LogFields{
					"topic":    msg.Metadata.Get("topic"),
					"duration": duration.String(),
				})
			} else {
				ctxLogger.Info("Message processed successfully", watermill.LogFields{
					"topic":    msg.Metadata.Get("topic"),
					"duration": duration.String(),
				})
			}

			return producedMessages, err
		}
	}
}
