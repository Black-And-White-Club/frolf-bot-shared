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

// CreateInstance creates a new zero-value instance of type T.
// REPLACEMENT: Replaces the old reflection-based NewInstance.
// Usage: utils.CreateInstance[MyStruct]() instead of utils.NewInstance(&MyStruct{})
func CreateInstance[T any]() *T {
	return new(T)
}

// NewInstance is a compatibility helper that creates a new instance based on
// the migration approach. It accepts either:
// - a factory function: func() interface{} and will call it, or
// - a pointer to a zero value (old behavior): &T{} and will allocate via reflection.
// This allows gradual migration from reflection to generics.
func NewInstance(ptr interface{}) interface{} {
	if ptr == nil {
		return nil
	}

	// If caller passed a factory function, call it and return the result.
	if fn, ok := ptr.(func() interface{}); ok {
		return fn()
	}

	// Fallback to reflection-based allocation for old-style calls like &T{}
	t := reflect.TypeOf(ptr)
	if t.Kind() == reflect.Ptr {
		return reflect.New(t.Elem()).Interface()
	}
	return nil
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
		// BLOCK LIST: prevent copying IDs and old routing topics
		case "message_id", "Nats-Msg-Id", "_watermill_message_uuid", "topic", "Topic":
			continue
		}
		newEvent.Metadata.Set(key, value)
	}

	// Set topic metadata used by the router/publisher for routing
	// Keep topic_hint for debugging/logging as well.
	newEvent.Metadata.Set("topic", topic)
	newEvent.Metadata.Set("Topic", topic)
	newEvent.Metadata.Set("topic_hint", topic)

	// Ensure correlation ID exists
	if newEvent.Metadata.Get(middleware.CorrelationIDMetadataKey) == "" {
		newCID := watermill.NewUUID()
		newEvent.Metadata.Set(middleware.CorrelationIDMetadataKey, newCID)
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

	// Set handler name and event name for observability
	newEvent.Metadata.Set("handler_name", "CreateResultMessage")
	newEvent.Metadata.Set("event_name", topic)

	// Set domain from topic (e.g., "user.tag.available.v1" → "user")
	if newEvent.Metadata.Get("domain") == "" {
		domain := extractDomainFromTopic(topic)
		if domain != "" {
			newEvent.Metadata.Set("domain", domain)
		}
	}

	// Debug: log the outgoing message metadata so router/publisher mapping can be observed
	if h.Logger != nil {
		h.Logger.Debug("created result message metadata",
			slog.String("event_name", newEvent.Metadata.Get("event_name")),
			slog.String("topic", newEvent.Metadata.Get("topic")),
			slog.String("Topic", newEvent.Metadata.Get("Topic")),
			slog.String("topic_hint", newEvent.Metadata.Get("topic_hint")),
			slog.String("correlation_id", newEvent.Metadata.Get(middleware.CorrelationIDMetadataKey)),
		)
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
	// Ensure publisher can discover the subject to publish to
	newEvent.Metadata.Set("topic", topic)
	newEvent.Metadata.Set("Topic", topic)
	newEvent.Metadata.Set("topic_hint", topic)
	newEvent.Metadata.Set("event_name", topic)

	// Set domain from topic
	domain := extractDomainFromTopic(topic)
	if domain != "" {
		newEvent.Metadata.Set("domain", domain)
	}

	// Ensure a correlation id exists on newly created messages
	if newEvent.Metadata.Get(middleware.CorrelationIDMetadataKey) == "" {
		newCID := watermill.NewUUID()
		newEvent.Metadata.Set(middleware.CorrelationIDMetadataKey, newCID)
		if h.Logger != nil {
			h.Logger.Debug("generated correlation id for new message", slog.String("correlation_id", newCID))
		}
	}

	// Debug: log the outgoing new message metadata
	if h.Logger != nil {
		h.Logger.Debug("created new message metadata",
			slog.String("event_name", newEvent.Metadata.Get("event_name")),
			slog.String("topic", newEvent.Metadata.Get("topic")),
			slog.String("Topic", newEvent.Metadata.Get("Topic")),
			slog.String("topic_hint", newEvent.Metadata.Get("topic_hint")),
			slog.String("correlation_id", newEvent.Metadata.Get(middleware.CorrelationIDMetadataKey)),
		)
	}

	return newEvent, nil
}

// UnmarshalPayload unmarshals the message payload into the provided struct.
// It detects "poison pill" messages (malformed JSON) and terminates them in NATS to prevent infinite redelivery.
func (h *DefaultHelper) UnmarshalPayload(msg *message.Message, payload interface{}) error {
	if err := json.Unmarshal(msg.Payload, payload); err != nil {
		// POISON PILL HANDLING
		// If unmarshal fails, we must terminate the message so NATS doesn't retry it forever.
		if jsMsg := getNATSMessage(msg); jsMsg != nil {
			terminateReason := fmt.Sprintf("Unmarshal failed: %s", err.Error())

			// Terminate the message (stops redelivery)
			if termErr := jsMsg.TermWithReason(terminateReason); termErr != nil {
				h.Logger.Error("failed to terminate poison message",
					slog.String("msg_uuid", msg.UUID),
					slog.String("term_error", termErr.Error()),
				)
			} else {
				h.Logger.Warn("terminated poison message",
					slog.String("msg_uuid", msg.UUID),
					slog.String("reason", terminateReason),
				)
			}
		} else {
			h.Logger.Debug("no NATS JetStream message found to terminate", slog.String("msg_uuid", msg.UUID))
		}
		return fmt.Errorf("failed to unmarshal payload for message %s: %w", msg.UUID, err)
	}
	return nil
}

// getNATSMessage attempts to extract the underlying NATS JetStream message from Watermill metadata.
func getNATSMessage(msg *message.Message) jetstream.Msg {
	// 1. Direct metadata reference (common in newer Watermill versions)
	if natsMsg := msg.Metadata.Get("_nats_jetstream_msg"); natsMsg != "" {
		if jsMsg, ok := msg.Context().Value(natsMsg).(jetstream.Msg); ok {
			return jsMsg
		}
	}

	// 2. Context key fallback
	if jsMsg, ok := msg.Context().Value("nats_jetstream_msg").(jetstream.Msg); ok {
		return jsMsg
	}

	return nil
}

// extractDomainFromTopic extracts the domain from a topic string.
func extractDomainFromTopic(topic string) string {
	for i := 0; i < len(topic); i++ {
		if topic[i] == '.' {
			return topic[:i]
		}
	}
	return topic
}
