package utils

import (
	"github.com/ThreeDotsLabs/watermill/message"
)

const (
	MetadataChannelID        = "channel_id"
	MetadataMessageID        = "message_id"
	MetadataGuildID          = "guild_id"
	MetadataInteractionID    = "interaction_id"
	MetadataInteractionToken = "interaction_token"
)

// MiddlewareHelpers defines the interface for handling metadata.
type MiddlewareHelpers interface {
	AddCommonMetadata(msg *message.Message, domain string)
	AddDiscordMetadata(msg *message.Message, guildID, interactionToken string)
	AddRoutingMetadata(msg *message.Message)
	CommonMetadataMiddleware(domain string) message.HandlerMiddleware
	DiscordMetadataMiddleware() message.HandlerMiddleware
	RoutingMetadataMiddleware() message.HandlerMiddleware
}

// Middleware is the default implementation of Helper.
type Middleware struct{}

// NewMiddlewareHelper creates a new MiddlewareHelper.
func NewMiddlewareHelper() MiddlewareHelpers {
	return &Middleware{}
}

// AddCommonMetadata adds common metadata fields to messages.
func (mh *Middleware) AddCommonMetadata(msg *message.Message, domain string) {
	msg.Metadata.Set("domain", domain)
	if msg.Metadata.Get("event_name") == "" {
		msg.Metadata.Set("event_name", "unknown_event")
	}
}

// AddDiscordMetadata adds Discord-specific metadata.
func (mh *Middleware) AddDiscordMetadata(msg *message.Message, guildID, interactionToken string) {
	if guildID != "" {
		msg.Metadata.Set(MetadataGuildID, guildID)
	}
	if interactionToken != "" {
		msg.Metadata.Set(MetadataInteractionToken, interactionToken)
	}
}

// AddRoutingMetadata ensures routing metadata is carried through the flow.
func (mh *Middleware) AddRoutingMetadata(msg *message.Message) {
	for _, key := range []string{MetadataChannelID, MetadataMessageID} {
		if value := msg.Metadata.Get(key); value != "" {
			msg.Metadata.Set(key, value)
		}
	}
}

// CommonMetadataMiddleware returns a middleware function that adds common metadata.
func (mh *Middleware) CommonMetadataMiddleware(domain string) message.HandlerMiddleware {
	return func(next message.HandlerFunc) message.HandlerFunc {
		return func(msg *message.Message) ([]*message.Message, error) {
			mh.AddCommonMetadata(msg, domain)
			return next(msg)
		}
	}
}

// DiscordMetadataMiddleware sets Discord-specific metadata.
func (mh *Middleware) DiscordMetadataMiddleware() message.HandlerMiddleware {
	return func(next message.HandlerFunc) message.HandlerFunc {
		return func(msg *message.Message) ([]*message.Message, error) {
			// Extract metadata *from the message*, not as arguments
			guildID := msg.Metadata.Get(MetadataGuildID)
			interactionToken := msg.Metadata.Get(MetadataInteractionToken)

			mh.AddDiscordMetadata(msg, guildID, interactionToken)
			return next(msg)
		}
	}
}

// RoutingMetadataMiddleware returns a middleware that copies routing-related metadata.
func (mh *Middleware) RoutingMetadataMiddleware() message.HandlerMiddleware {
	return func(next message.HandlerFunc) message.HandlerFunc {
		return func(msg *message.Message) ([]*message.Message, error) {
			mh.AddRoutingMetadata(msg)
			return next(msg)
		}
	}
}
