package utils

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
)

const (
	MetadataChannelID        = "channel_id"
	MetadataMessageID        = "discord_message_id"
	MetadataGuildID          = "guild_id"
	MetadataInteractionID    = "interaction_id"
	MetadataInteractionToken = "interaction_token"
)

// MiddlewareHelpers defines the interface for handling metadata.
type MiddlewareHelpers interface {
	AddCommonMetadata(msg *message.Message, domain string)
	AddDiscordMetadata(msg *message.Message, guildID, interactionToken string)
	AddRoutingMetadata(from *message.Message, to *message.Message)
	CommonMetadataMiddleware(domain string) message.HandlerMiddleware
	DiscordMetadataMiddleware() message.HandlerMiddleware
	RoutingMetadataMiddleware() message.HandlerMiddleware
}

// Middleware is the default implementation of MiddlewareHelpers.
type Middleware struct{}

// NewMiddlewareHelper creates a new MiddlewareHelper.
func NewMiddlewareHelper() MiddlewareHelpers {
	return &Middleware{}
}

// AddCommonMetadata adds common metadata fields to a message if missing.
func (mh *Middleware) AddCommonMetadata(msg *message.Message, domain string) {
	msg.Metadata.Set("domain", domain)

	if msg.Metadata.Get("event_name") == "" {
		msg.Metadata.Set("event_name", "unknown_event")
	}
}

// AddDiscordMetadata adds Discord-specific metadata fields to a message.
func (mh *Middleware) AddDiscordMetadata(msg *message.Message, guildID, interactionToken string) {
	if guildID != "" {
		msg.Metadata.Set(MetadataGuildID, guildID)
	}
	if interactionToken != "" {
		msg.Metadata.Set(MetadataInteractionToken, interactionToken)
	}
}

// AddRoutingMetadata copies allowed routing metadata keys from one message to another.
func (mh *Middleware) AddRoutingMetadata(from *message.Message, to *message.Message) {
	for _, key := range []string{
		"domain",
		"event_name",
		middleware.CorrelationIDMetadataKey,
		MetadataChannelID,
		MetadataMessageID,
		MetadataGuildID,
		MetadataInteractionID,
		MetadataInteractionToken,
	} {
		if val := from.Metadata.Get(key); val != "" {
			to.Metadata.Set(key, val)
		}
	}
}

// CommonMetadataMiddleware attaches standard common metadata to incoming messages.
func (mh *Middleware) CommonMetadataMiddleware(domain string) message.HandlerMiddleware {
	return func(next message.HandlerFunc) message.HandlerFunc {
		return func(msg *message.Message) ([]*message.Message, error) {
			mh.AddCommonMetadata(msg, domain)
			return next(msg)
		}
	}
}

// DiscordMetadataMiddleware reads Discord-specific metadata from incoming message and re-applies it.
func (mh *Middleware) DiscordMetadataMiddleware() message.HandlerMiddleware {
	return func(next message.HandlerFunc) message.HandlerFunc {
		return func(msg *message.Message) ([]*message.Message, error) {
			guildID := msg.Metadata.Get(MetadataGuildID)
			interactionToken := msg.Metadata.Get(MetadataInteractionToken)

			mh.AddDiscordMetadata(msg, guildID, interactionToken)
			return next(msg)
		}
	}
}

// RoutingMetadataMiddleware propagates a strict set of routing metadata fields from the input message to all output messages.
func (mh *Middleware) RoutingMetadataMiddleware() message.HandlerMiddleware {
	return func(next message.HandlerFunc) message.HandlerFunc {
		return func(msg *message.Message) ([]*message.Message, error) {
			outMsgs, err := next(msg)
			if err != nil {
				return nil, err
			}

			for _, outMsg := range outMsgs {
				mh.AddRoutingMetadata(msg, outMsg)
			}

			return outMsgs, nil
		}
	}
}
