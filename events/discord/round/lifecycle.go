// Package round contains Discord-specific round events.
//
// This file defines the Discord Round Lifecycle Flow - events specific to
// round state changes that need Discord UI updates (start, finalize, reminders).
//
// # Flow Sequences
//
// ## Round Start Flow
//  1. Round starts -> RoundStartedDiscordV1
//  2. Discord embed updated with start notification
//
// ## Round Finalize Flow
//  1. All scores submitted -> RoundFinalizedDiscordNotifyV1
//  2. Discord embed updated with final results
//
// ## Reminder Flow
//  1. Reminder triggered -> RoundReminderDiscordV1
//  2. Discord reminder sent to participants
//
// # Relationship to Domain Events
//
// These Discord events subscribe to domain events in events/round/lifecycle.go:
//   - RoundStartedDiscordV1 <- subscribes to RoundStartedV1 (domain)
//   - RoundFinalizedDiscordNotifyV1 <- subscribes to RoundFinalizedV1 (domain)
//   - RoundReminderDiscordV1 <- subscribes to RoundReminderScheduledV1 (domain)
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
// DISCORD ROUND LIFECYCLE FLOW - Event Constants
// =============================================================================

// -----------------------------------------------------------------------------
// Round Start Events
// -----------------------------------------------------------------------------

// RoundStartedDiscordV1 is published to notify Discord that a round has started.
//
// Pattern: Event Notification
// Subject: discord.round.started.v1
// Producer: discord-service (after receiving domain RoundStartedV1)
// Consumers: discord-service (embed updater)
// Triggers: Discord embed updated to show round started
// Version: v1 (December 2024)
const RoundStartedDiscordV1 = "discord.round.started.v1"

// -----------------------------------------------------------------------------
// Round Finalize Events
// -----------------------------------------------------------------------------

// RoundFinalizedDiscordNotifyV1 is published to notify Discord of round finalization.
//
// Pattern: Event Notification
// Subject: discord.round.finalized.v1
// Producer: discord-service (after receiving domain RoundFinalizedV1)
// Consumers: discord-service (embed updater)
// Triggers: Discord embed updated with final scores
// Version: v1 (December 2024)
const RoundFinalizedDiscordNotifyV1 = "discord.round.finalized.v1"

// -----------------------------------------------------------------------------
// Reminder Events
// -----------------------------------------------------------------------------

// RoundReminderDiscordV1 is published to send a reminder via Discord.
//
// Pattern: Event Notification
// Subject: discord.round.reminder.v1
// Producer: discord-service (after receiving domain RoundReminderScheduledV1)
// Consumers: discord-service (reminder sender)
// Triggers: Discord reminder message sent to participants
// Version: v1 (December 2024)
const RoundReminderDiscordV1 = "discord.round.reminder.v1"

// =============================================================================
// DISCORD ROUND LIFECYCLE FLOW - Payload Types
// =============================================================================

// -----------------------------------------------------------------------------
// Round Start Payloads
// -----------------------------------------------------------------------------

// RoundStartedDiscordPayloadV1 contains round start data for Discord.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoundStartedDiscordPayloadV1 struct {
	RoundID   sharedtypes.RoundID    `json:"round_id"`
	Title     roundtypes.Title       `json:"title"`
	Location  *roundtypes.Location   `json:"location"`
	StartTime *sharedtypes.StartTime `json:"start_time"`
	ChannelID string                 `json:"channel_id"`
	GuildID   string                 `json:"guild_id"`
}

// -----------------------------------------------------------------------------
// Round Finalize Payloads
// -----------------------------------------------------------------------------

// RoundFinalizedDiscordPayloadV1 contains finalization data for Discord.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoundFinalizedDiscordPayloadV1 struct {
	RoundID   sharedtypes.RoundID `json:"round_id"`
	ChannelID string              `json:"channel_id"`
	MessageID string              `json:"discord_message_id"`
	GuildID   string              `json:"guild_id"`
}

// -----------------------------------------------------------------------------
// Reminder Payloads
// -----------------------------------------------------------------------------

// RoundReminderDiscordPayloadV1 contains reminder data for Discord.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoundReminderDiscordPayloadV1 struct {
	RoundID      sharedtypes.RoundID `json:"round_id"`
	RoundTitle   roundtypes.Title    `json:"round_title"`
	UserIDs      []string            `json:"user_ids"`
	ReminderType string              `json:"reminder_type"`
	ChannelID    string              `json:"channel_id"`
	GuildID      string              `json:"guild_id"`
}
