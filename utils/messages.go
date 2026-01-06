package utils

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"reflect"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/nats-io/nats.go/jetstream"
)

// Helpers defines the interface for utility functions.
type Helpers interface {
	CreateResultMessage(originalMsg *message.Message, payload interface{}, topic string) (*message.Message, error)
	CreateNewMessage(payload interface{}, topic string) (*message.Message, error)
	UnmarshalPayload(msg *message.Message, payload interface{}) error
}

// DefaultHelper is the default implementation of Helpers.
type DefaultHelper struct {
	Logger *slog.Logger
}

// NewHelper creates a new DefaultHelper.
func NewHelper(logger *slog.Logger) Helpers {
	return &DefaultHelper{Logger: logger}
}

// NewInstance creates a new instance of the type of the provided interface, which is expected to be a pointer to a struct.
// It returns a pointer to the newly created instance.
func NewInstance(ptr interface{}) interface{} {
	if ptr == nil {
		return nil
	}
	// Get the type of the element pointed to by ptr
	elemType := reflect.TypeOf(ptr).Elem()
	// Create a new instance of that type and return a pointer to it
	return reflect.New(elemType).Interface()
}

// CreateResultMessage creates a new Watermill message, carrying over metadata from the original message.
func (h *DefaultHelper) CreateResultMessage(originalMsg *message.Message, payload interface{}, topic string) (*message.Message, error) {
	if originalMsg == nil {
		return h.CreateNewMessage(payload, topic)
	}

	newEvent := message.NewMessage(watermill.NewUUID(), nil)

	// Copy original metadata first
	for key, value := range originalMsg.Metadata {
		switch key {
		case "message_id", "Nats-Msg-Id", "_watermill_message_uuid":
			continue
		}
		newEvent.Metadata.Set(key, value)
	}

	// Overwrite topic deterministically
	newEvent.Metadata.Set("topic", topic)
	// Remove conflicting alternate casing like "Topic", if present
	delete(newEvent.Metadata, "Topic")

	// Ensure correlation ID exists
	if newEvent.Metadata.Get(middleware.CorrelationIDMetadataKey) == "" {
		newCID := watermill.NewUUID()
		newEvent.Metadata.Set(middleware.CorrelationIDMetadataKey, newCID)
		// This is a normal path; keep level at Debug to avoid noise
		h.Logger.Debug("generated correlation id for result message",
			slog.String("correlation_id", newCID),
		)
	}

	// Marshal payload
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload for topic %s: %w", topic, err)
	}
	newEvent.Payload = payloadBytes

	// Set handler
	newEvent.Metadata.Set("handler_name", "CreateResultMessage")
	
	// Set event_name to match the topic
	newEvent.Metadata.Set("event_name", topic)
	
	// Set domain from topic (e.g., "user.tag.available.v1" → "user")
	// Only set domain if not already present in metadata
	if newEvent.Metadata.Get("domain") == "" {
		// Extract domain from topic (first segment before first dot)
		domain := extractDomainFromTopic(topic)
		if domain != "" {
			newEvent.Metadata.Set("domain", domain)
		}
	}

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
	newEvent.Metadata.Set("event_name", topic)
	
	// Set domain from topic (e.g., "user.tag.available.v1" → "user")
	domain := extractDomainFromTopic(topic)
	if domain != "" {
		newEvent.Metadata.Set("domain", domain)
	}

	return newEvent, nil
}

// UnmarshalPayload is a generic helper to unmarshal message payloads.
// This updated version checks for NATS JetStream message and terminates it on unmarshal errors
func (h *DefaultHelper) UnmarshalPayload(msg *message.Message, payload interface{}) error {
	if err := json.Unmarshal(msg.Payload, payload); err != nil {
		// Try to get the underlying NATS JetStream message if it exists
		if jsMsg := getNATSMessage(msg); jsMsg != nil {
			terminateReason := fmt.Sprintf("Failed to unmarshal payload: %s", err.Error())
			termErr := jsMsg.TermWithReason(terminateReason)
			if termErr != nil {
				h.Logger.Error("failed to terminate message after unmarshal error",
					slog.String("msg_uuid", msg.UUID),
					slog.String("unmarshal_error", err.Error()),
					slog.String("term_error", termErr.Error()),
				)
			} else {
				h.Logger.Warn("message terminated due to unmarshal error",
					slog.String("msg_uuid", msg.UUID),
					slog.String("reason", terminateReason),
				)
			}
		} else {
			// Avoid extra noise; keep as Debug since we'll return the error
			h.Logger.Debug("no NATS JetStream message found to terminate",
				slog.String("msg_uuid", msg.UUID),
			)
		}
		return fmt.Errorf("failed to unmarshal payload for message %s: %w", msg.UUID, err)
	}
	return nil
}

// getNATSMessage attempts to extract the underlying NATS JetStream message
// from a Watermill message based on expected metadata structure
func getNATSMessage(msg *message.Message) jetstream.Msg {
	// Check if the original message is stored in the internal metadata
	// The exact key depends on how you're implementing the NATS subscriber
	// Common places to check:

	// 1. Check if there's a direct message reference stored
	if natsMsg := msg.Metadata.Get("_nats_jetstream_msg"); natsMsg != "" {
		if jsMsg, ok := msg.Context().Value(natsMsg).(jetstream.Msg); ok {
			return jsMsg
		}
	}

	// 2. Check if the original message is stored with a specific key in context
	if jsMsg, ok := msg.Context().Value("nats_jetstream_msg").(jetstream.Msg); ok {
		return jsMsg
	}

	// 3. Check if there's a NATS message ID that could be used to look up the message
	// This would require access to the NATS connection, which we don't have here

	return nil
}

// extractDomainFromTopic extracts the domain from a topic string.
// For example: "user.tag.available.v1" → "user", "score.update.requested.v1" → "score"
func extractDomainFromTopic(topic string) string {
	// Find the first dot
	for i := 0; i < len(topic); i++ {
		if topic[i] == '.' {
			return topic[:i]
		}
	}
	// If no dot found, return the whole topic
	return topic
}
