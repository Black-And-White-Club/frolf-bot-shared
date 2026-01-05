// Package round contains Discord-specific round events.
//
// This file defines the Discord Round Creation Flow - events specific to
// creating rounds through Discord modal submissions and receiving Discord
// notifications about creation success/failure.
//
// # Flow Sequence
//
//  1. User submits modal -> RoundCreateModalSubmittedV1
//  2. Success notification -> RoundCreatedDiscordV1
//  3. OR Failure notification -> RoundCreationFailedDiscordV1
//
// # Relationship to Domain Events
//
// These Discord events wrap/trigger domain events in events/round/creation.go:
//   - RoundCreateModalSubmittedV1 -> publishes RoundCreationRequestedV1 (domain)
//   - RoundCreatedDiscordV1 <- subscribes to RoundCreatedV1 (domain)
//
// # Versioning Strategy
//
// All events include a V1 suffix for future schema evolution.
package round

import (
	roundtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/round"
	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
)

// =============================================================================
// DISCORD ROUND CREATION FLOW - Event Constants
// =============================================================================

// RoundCreateModalSubmittedV1 is published when a user submits the round creation modal.
//
// Pattern: Event Notification
// Subject: discord.round.modal.submitted.v1
// Producer: discord-service (modal handler)
// Consumers: discord-service (round creation handler)
// Triggers: Domain event RoundCreationRequestedV1
// Version: v1 (December 2024)
const RoundCreateModalSubmittedV1 = "discord.round.modal.submitted.v1"

// RoundCreatedDiscordV1 is published to notify Discord that a round was created.
//
// Pattern: Event Notification
// Subject: discord.round.created.v1
// Producer: discord-service (after receiving domain RoundCreatedV1)
// Consumers: discord-service (embed creator)
// Triggers: Discord embed created in channel
// Version: v1 (December 2024)
const RoundCreatedDiscordV1 = "discord.round.created.v1"

// RoundCreationFailedDiscordV1 is published when round creation fails.
//
// Pattern: Event Notification
// Subject: discord.round.creation.failed.v1
// Producer: discord-service (error handler)
// Consumers: discord-service (error message sender)
// Triggers: Discord error message to user
// Version: v1 (December 2024)
const RoundCreationFailedDiscordV1 = "discord.round.creation.failed.v1"

// RoundValidationFailedDiscordV1 is published when round validation fails.
//
// Pattern: Event Notification
// Subject: discord.round.validation.failed.v1
// Producer: discord-service
// Consumers: discord-service (error message sender)
// Version: v1 (December 2024)
const RoundValidationFailedDiscordV1 = "discord.round.validation.failed.v1"

// RoundCreatedTraceV1 is published for tracing round creation.
//
// Pattern: Event Notification
// Subject: discord.round.created.trace.v1
// Producer: discord-service
// Consumers: Observability systems
// Version: v1 (December 2024)
const RoundCreatedTraceV1 = "discord.round.created.trace.v1"

// =============================================================================
// DISCORD ROUND CREATION FLOW - Payload Types
// =============================================================================

// CreateRoundModalPayloadV1 contains data from the round creation modal.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type CreateRoundModalPayloadV1 struct {
	Title       roundtypes.Title       `json:"title"`
	Description roundtypes.Description `json:"description"`
	StartTime   string                 `json:"start_time"` // Unparsed natural language
	Location    roundtypes.Location    `json:"location"`
	UserID      sharedtypes.DiscordID  `json:"user_id"`
	ChannelID   string                 `json:"channel_id"`
	Timezone    roundtypes.Timezone    `json:"timezone"`
	GuildID     sharedtypes.GuildID    `json:"guild_id"`
}

// RoundCreatedDiscordPayloadV1 contains round creation success data for Discord.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoundCreatedDiscordPayloadV1 struct {
	RoundID     sharedtypes.RoundID   `json:"round_id"`
	Title       roundtypes.Title      `json:"title"`
	StartTime   sharedtypes.StartTime `json:"start_time"`
	Location    roundtypes.Location   `json:"location"`
	RequesterID sharedtypes.DiscordID `json:"requester_id"`
	ChannelID   string                `json:"channel_id"`
	GuildID     sharedtypes.GuildID   `json:"guild_id"`
}

// RoundCreationFailedDiscordPayloadV1 contains round creation failure data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoundCreationFailedDiscordPayloadV1 struct {
	UserID  sharedtypes.DiscordID `json:"user_id"`
	Reason  string                `json:"reason"`
	GuildID string                `json:"guild_id"`
}

// RoundCreatedTracePayloadV1 contains round creation trace data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoundCreatedTracePayloadV1 struct {
	RoundID   sharedtypes.RoundID   `json:"round_id"`
	Title     roundtypes.Title      `json:"title"`
	CreatedBy sharedtypes.DiscordID `json:"created_by"`
	GuildID   string                `json:"guild_id"`
}
