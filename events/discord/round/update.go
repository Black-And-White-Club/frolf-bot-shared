// Package discordroundevents contains Discord-specific round events.
//
// This file defines the Discord Round Update Flow - events specific to
// updating rounds through Discord modal submissions and deletion requests.
//
// # Flow Sequences
//
// ## Update Flow
//  1. User submits update modal -> RoundUpdateModalSubmittedV1
//  2. Success notification -> RoundUpdatedDiscordV1
//
// ## Delete Flow
//  1. User requests deletion -> RoundDeleteRequestDiscordV1
//  2. Success notification -> RoundDeletedDiscordV1
//
// # Relationship to Domain Events
//
// These Discord events wrap/trigger domain events in events/round/update.go and delete.go:
//   - RoundUpdateModalSubmittedV1 -> publishes RoundUpdateRequestedV1 (domain)
//   - RoundUpdatedDiscordV1 <- subscribes to RoundUpdatedV1 (domain)
//   - RoundDeleteRequestDiscordV1 -> publishes RoundDeleteRequestedV1 (domain)
//   - RoundDeletedDiscordV1 <- subscribes to RoundDeletedV1 (domain)
//
// # Versioning Strategy
//
// All events include a V1 suffix for future schema evolution.
package discordroundevents

import (
	roundtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/round"
	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
)

// =============================================================================
// DISCORD ROUND UPDATE FLOW - Event Constants
// =============================================================================

// -----------------------------------------------------------------------------
// Update Events
// -----------------------------------------------------------------------------

// RoundUpdateModalSubmittedV1 is published when a user submits the round update modal.
//
// Pattern: Event Notification
// Subject: discord.round.update.modal.submitted.v1
// Producer: discord-service (modal handler)
// Consumers: discord-service (update handler)
// Triggers: Domain event RoundUpdateRequestedV1
// Version: v1 (December 2024)
const RoundUpdateModalSubmittedV1 = "discord.round.update.modal.submitted.v1"

// RoundUpdateRequestDiscordV1 is published when a round update is requested via Discord.
//
// Pattern: Event Notification
// Subject: discord.round.update.request.v1
// Producer: discord-service
// Consumers: discord-service (update handler)
// Version: v1 (December 2024)
const RoundUpdateRequestDiscordV1 = "discord.round.update.request.v1"

// RoundUpdatedDiscordV1 is published to notify Discord that a round was updated.
//
// Pattern: Event Notification
// Subject: discord.round.updated.v1
// Producer: discord-service (after receiving domain RoundUpdatedV1)
// Consumers: discord-service (embed updater)
// Triggers: Discord embed updated with new round details
// Version: v1 (December 2024)
const RoundUpdatedDiscordV1 = "discord.round.updated.v1"

// -----------------------------------------------------------------------------
// Delete Events
// -----------------------------------------------------------------------------

// RoundDeleteRequestDiscordV1 is published when a round deletion is requested via Discord.
//
// Pattern: Event Notification
// Subject: discord.round.delete.request.v1
// Producer: discord-service (command handler)
// Consumers: discord-service (delete handler)
// Triggers: Domain event RoundDeleteRequestedV1
// Version: v1 (December 2024)
const RoundDeleteRequestDiscordV1 = "discord.round.delete.request.v1"

// RoundDeletedDiscordV1 is published to notify Discord that a round was deleted.
//
// Pattern: Event Notification
// Subject: discord.round.deleted.v1
// Producer: discord-service (after receiving domain RoundDeletedV1)
// Consumers: discord-service (message handler)
// Triggers: Discord embed removed or marked as deleted
// Version: v1 (December 2024)
const RoundDeletedDiscordV1 = "discord.round.deleted.v1"

// =============================================================================
// DISCORD ROUND UPDATE FLOW - Payload Types
// =============================================================================

// -----------------------------------------------------------------------------
// Update Payloads
// -----------------------------------------------------------------------------

// RoundUpdateModalSubmittedPayloadV1 contains data from the round update modal.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoundUpdateModalSubmittedPayloadV1 struct {
	RoundID     sharedtypes.RoundID     `json:"round_id"`
	UserID      sharedtypes.DiscordID   `json:"user_id"`
	MessageID   string                  `json:"message_id"`
	Title       *roundtypes.Title       `json:"title,omitempty"`
	Description *roundtypes.Description `json:"description,omitempty"`
	StartTime   *string                 `json:"start_time,omitempty"` // Unparsed
	Timezone    *roundtypes.Timezone    `json:"timezone,omitempty"`
	Location    *roundtypes.Location    `json:"location,omitempty"`
	ChannelID   string                  `json:"channel_id"`
	GuildID     sharedtypes.GuildID     `json:"guild_id"`
}

// RoundUpdateRequestDiscordPayloadV1 contains update request data from Discord.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoundUpdateRequestDiscordPayloadV1 struct {
	RoundID     sharedtypes.RoundID     `json:"round_id"`
	UserID      sharedtypes.DiscordID   `json:"user_id"`
	MessageID   string                  `json:"message_id"`
	Title       *roundtypes.Title       `json:"title,omitempty"`
	Description *roundtypes.Description `json:"description,omitempty"`
	StartTime   *sharedtypes.StartTime  `json:"start_time,omitempty"`
	Location    *roundtypes.Location    `json:"location,omitempty"`
	ChannelID   string                  `json:"channel_id"`
	GuildID     sharedtypes.GuildID     `json:"guild_id"`
}

// RoundUpdatedDiscordPayloadV1 contains update confirmation data for Discord.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoundUpdatedDiscordPayloadV1 struct {
	RoundID     sharedtypes.RoundID     `json:"round_id"`
	MessageID   string                  `json:"message_id"`
	ChannelID   string                  `json:"channel_id"`
	Title       *roundtypes.Title       `json:"title,omitempty"`
	Description *roundtypes.Description `json:"description,omitempty"`
	StartTime   *sharedtypes.StartTime  `json:"start_time,omitempty"`
	Location    *roundtypes.Location    `json:"location,omitempty"`
	GuildID     sharedtypes.GuildID     `json:"guild_id"`
}

// -----------------------------------------------------------------------------
// Delete Payloads
// -----------------------------------------------------------------------------

// RoundDeleteRequestDiscordPayloadV1 contains deletion request data from Discord.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoundDeleteRequestDiscordPayloadV1 struct {
	RoundID   sharedtypes.RoundID   `json:"round_id"`
	UserID    sharedtypes.DiscordID `json:"user_id"`
	ChannelID string                `json:"channel_id"`
	MessageID string                `json:"message_id"`
	GuildID   string                `json:"guild_id"`
}

// RoundDeletedDiscordPayloadV1 contains deletion confirmation data for Discord.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoundDeletedDiscordPayloadV1 struct {
	RoundID   sharedtypes.RoundID `json:"round_id"`
	ChannelID string              `json:"channel_id"`
	MessageID string              `json:"message_id"`
	GuildID   string              `json:"guild_id"`
}
