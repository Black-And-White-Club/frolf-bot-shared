package utils

import (
	"encoding/json"
	"fmt"

	lokifrolfbot "github.com/Black-And-White-Club/frolf-bot-shared/observability/loki"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
)

// Helpers defines the interface for utility functions.
type Helpers interface {
	CreateResultMessage(originalMsg *message.Message, payload interface{}, topic string) (*message.Message, error)
	CreateNewMessage(payload interface{}, topic string) (*message.Message, error)
	UnmarshalPayload(msg *message.Message, payload interface{}) error
}

// DefaultHelper is the default implementation of Helpers.
type DefaultHelper struct {
	Logger lokifrolfbot.Logger
}

// NewHelper creates a new DefaultHelper.
func NewHelper(logger lokifrolfbot.Logger) Helpers {
	return &DefaultHelper{Logger: logger}
}

// CreateResultMessage creates a new Watermill message, carrying over metadata from the original message.
func (h *DefaultHelper) CreateResultMessage(originalMsg *message.Message, payload interface{}, topic string) (*message.Message, error) {
	if originalMsg == nil {
		return h.CreateNewMessage(payload, topic)
	}

	newEvent := message.NewMessage(watermill.NewUUID(), nil)

	// Copy metadata from the original message. This ensures correlation ID is propagated.
	for key, value := range originalMsg.Metadata {
		newEvent.Metadata.Set(key, value)
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload for topic %s: %w", topic, err)
	}
	newEvent.Payload = payloadBytes

	newEvent.Metadata.Set("handler_name", "CreateResultMessage")
	newEvent.Metadata.Set("topic", topic)
	newEvent.Metadata.Set("domain", "discord")

	return newEvent, nil
}

// CreateNewMessage creates a new Watermill message **without** an original message.
func (h *DefaultHelper) CreateNewMessage(payload interface{}, topic string) (*message.Message, error) {
	newEvent := message.NewMessage(watermill.NewUUID(), nil)

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload for topic %s: %w", topic, err)
	}
	newEvent.Payload = payloadBytes

	newEvent.Metadata.Set("handler_name", "CreateNewMessage")
	newEvent.Metadata.Set("topic", topic)
	newEvent.Metadata.Set("domain", "discord")

	return newEvent, nil
}

// UnmarshalPayload is a generic helper to unmarshal message payloads.
func (h *DefaultHelper) UnmarshalPayload(msg *message.Message, payload interface{}) error {
	if err := json.Unmarshal(msg.Payload, payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload for message %s: %w", msg.UUID, err)
	}
	return nil
}
