// Package user contains Discord-specific user events.
//
// This file defines the Discord User Tag Flow - events specific to
// prompting users for tag numbers through Discord interactions.
//
// # Flow Sequence
//
//  1. Tag number requested -> TagNumberRequestedV1
//  2. User responds -> TagNumberResponseV1
//
// # Relationship to Domain Events
//
// These Discord events are UI-layer events that wrap domain events:
//   - TagNumberResponseV1 -> feeds into TagAssignmentRequested (domain)
//
// # Versioning Strategy
//
// All events include a V1 suffix for future schema evolution.
package user

import (
	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
)

// =============================================================================
// DISCORD USER TAG FLOW - Event Constants
// =============================================================================

// TagNumberRequestedV1 is published to request a user's tag number via Discord.
//
// Pattern: Event Notification
// Subject: discord.user.tag.number.requested.v1
// Producer: discord-service (signup/tag handler)
// Consumers: discord-service (tag prompt handler)
// Triggers: Discord prompt for tag number
// Version: v1 (December 2024)
const TagNumberRequestedV1 = "discord.user.tag.number.requested.v1"

// TagNumberResponseV1 is published when a user provides their tag number.
//
// Pattern: Event Notification
// Subject: discord.user.tag.number.response.v1
// Producer: discord-service (modal/message handler)
// Consumers: discord-service (tag processor)
// Triggers: Tag validation/assignment flow
// Version: v1 (December 2024)
const TagNumberResponseV1 = "discord.user.tag.number.response.v1"

// =============================================================================
// DISCORD USER TAG FLOW - Payload Types
// =============================================================================

// TagNumberRequestedPayloadV1 contains tag number request data.
// Note: Uses interface{} for Interaction to avoid discordgo dependency in shared package.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type TagNumberRequestedPayloadV1 struct {
	UserID      sharedtypes.DiscordID `json:"user_id"`
	Interaction interface{}           `json:"interaction"` // *discordgo.Interaction in discord-frolf-bot
	GuildID     string                `json:"guild_id"`
}

// TagNumberResponsePayloadV1 contains user's tag number response.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type TagNumberResponsePayloadV1 struct {
	TagNumber string                `json:"tag_number"`
	UserID    sharedtypes.DiscordID `json:"user_id"`
	MessageID string                `json:"message_id"`
	GuildID   string                `json:"guild_id"`
}

// TagNumberProvidedPayloadV1 contains validated tag number data.
// Note: Uses interface{} for Interaction to avoid discordgo dependency in shared package.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type TagNumberProvidedPayloadV1 struct {
	UserID      sharedtypes.DiscordID `json:"user_id"`
	TagNumber   sharedtypes.TagNumber `json:"tag_number"`
	Interaction interface{}           `json:"interaction"` // *discordgo.Interaction in discord-frolf-bot
	GuildID     string                `json:"guild_id"`
}
