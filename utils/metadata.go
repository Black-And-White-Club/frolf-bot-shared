package utils

import (
	"encoding/json"
	"fmt"

	"github.com/Black-And-White-Club/frolf-bot-shared/events"
	"github.com/Black-And-White-Club/frolf-bot-shared/observability/attr"
	lokifrolfbot "github.com/Black-And-White-Club/frolf-bot-shared/observability/loki"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
)

// WithMetadata sets common metadata on a Watermill message and marshals the payload.
func (e *EventUtilImpl) WithMetadata(msg *message.Message, payload interface{}, logger lokifrolfbot.Logger) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		// Use the provided logger (no nil check needed, as it is a dependency).
		logger.Error("Failed to marshal payload", attr.Error(err))
		return fmt.Errorf("failed to marshal payload: %w", err)
	}
	msg.Payload = payloadBytes

	// Handle the common metadata, checking if there is an override
	if p, ok := payload.(events.MetadataCarrier); ok {
		msg.Metadata.Set("event_name", p.GetEventName())
		msg.Metadata.Set("domain", p.GetDomain())
	} else if p, ok := payload.(events.CommonMetadata); ok {
		msg.Metadata.Set("event_name", p.GetEventName())
		msg.Metadata.Set("domain", p.GetDomain())
	} else if p, ok := payload.(*events.CommonMetadata); ok {
		msg.Metadata.Set("event_name", p.GetEventName())
		msg.Metadata.Set("domain", p.GetDomain())
	}
	return nil
}

type EventUtil interface {
	PropagateMetadata(srcMsg *message.Message, dstMsg *message.Message)
	WithMetadata(msg *message.Message, payload interface{}, logger lokifrolfbot.Logger) error
}

type EventUtilImpl struct{}

// NewEventUtil creates a new EventUtil.
func NewEventUtil() EventUtil {
	return &EventUtilImpl{}
}

// PropagateMetadata copies the correlation ID from the source message to the destination message.
func (e *EventUtilImpl) PropagateMetadata(srcMsg *message.Message, dstMsg *message.Message) {
	correlationID := srcMsg.Metadata.Get(middleware.CorrelationIDMetadataKey)
	if correlationID != "" {
		dstMsg.Metadata.Set(middleware.CorrelationIDMetadataKey, correlationID)
	}
}
