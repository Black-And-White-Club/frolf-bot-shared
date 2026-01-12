// Package score contains Discord-specific score events.
//
// This file defines the Discord Score Update Flow - events specific to
// handling score updates through Discord interactions.
//
// # Flow Sequence
//
//  1. User requests score update -> ScoreUpdateRequestDiscordV1
//  2. Success response -> ScoreUpdateResponseDiscordV1
//  3. OR Failure response -> ScoreUpdateFailedDiscordV1
//
// # Relationship to Domain Events
//
// These Discord events wrap/trigger domain events in events/score/:
//   - ScoreUpdateRequestDiscordV1 -> publishes ScoreUpdateRequestedV1 (domain)
//   - ScoreUpdateResponseDiscordV1 <- subscribes to ScoreUpdatedV1 (domain)
//
// # Versioning Strategy
//
// All events include a V1 suffix for future schema evolution.
package score

import (
	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
)

// =============================================================================
// DISCORD SCORE UPDATE FLOW - Event Constants
// =============================================================================

// ScoreUpdateRequestDiscordV1 is published when a score update is requested via Discord.
//
// Pattern: Event Notification
// Subject: discord.score.update.request.v1
// Producer: discord-service (command/button handler)
// Consumers: discord-service (score handler)
// Triggers: Domain ScoreUpdateRequestedV1
// Version: v1 (December 2024)
const ScoreUpdateRequestDiscordV1 = "discord.score.update.request.v1"

// ScoreUpdateResponseDiscordV1 is published to notify Discord of successful score update.
//
// Pattern: Event Notification
// Subject: discord.score.update.response.v1
// Producer: discord-service (after domain response)
// Consumers: discord-service (embed/message handler)
// Version: v1 (December 2024)
const ScoreUpdateResponseDiscordV1 = "discord.score.update.response.v1"

// ScoreUpdateFailedDiscordV1 is published when a score update fails.
//
// Pattern: Event Notification
// Subject: discord.score.update.failed.v1
// Producer: discord-service
// Consumers: discord-service (error handler)
// Version: v1 (December 2024)
const ScoreUpdateFailedDiscordV1 = "discord.score.update.failed.v1"

// ScoreBulkUpdateRequestDiscordV1 is published when bulk score updates are requested via Discord.
//
// Pattern: Event Notification
// Subject: discord.score.bulk.update.request.v1
// Producer: discord-service (admin command handler)
// Consumers: discord-service (bulk score handler)
// Triggers: Domain ScoreBulkUpdateRequestedV1
// Version: v1 (December 2024)
const ScoreBulkUpdateRequestDiscordV1 = "discord.score.bulk.update.request.v1"

// ScoreBulkUpdateResponseDiscordV1 is published to notify Discord of bulk score update completion.
//
// Pattern: Event Notification
// Subject: discord.score.bulk.update.response.v1
// Producer: discord-service (after domain response)
// Consumers: discord-service (summary message handler)
// Version: v1 (December 2024)
const ScoreBulkUpdateResponseDiscordV1 = "discord.score.bulk.update.response.v1"

// =============================================================================
// DISCORD SCORE UPDATE FLOW - Payload Types
// =============================================================================

// ScoreUpdateRequestDiscordPayloadV1 contains score update request data from Discord.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type ScoreUpdateRequestDiscordPayloadV1 struct {
	GuildID   sharedtypes.GuildID    `json:"guild_id"`
	RoundID   sharedtypes.RoundID    `json:"round_id"`
	UserID    sharedtypes.DiscordID  `json:"user_id"`
	Score     sharedtypes.Score      `json:"score"`
	TagNumber *sharedtypes.TagNumber `json:"tag_number,omitempty"`
	ChannelID string                 `json:"channel_id"`
	MessageID string                 `json:"message_id"`
}

// ScoreUpdateResponseDiscordPayloadV1 contains score update success data for Discord.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type ScoreUpdateResponseDiscordPayloadV1 struct {
	GuildID   sharedtypes.GuildID   `json:"guild_id"`
	RoundID   sharedtypes.RoundID   `json:"round_id"`
	UserID    sharedtypes.DiscordID `json:"user_id"`
	Score     sharedtypes.Score     `json:"score"`
	ChannelID string                `json:"channel_id"`
	MessageID string                `json:"message_id"`
}

// ScoreUpdateFailedDiscordPayloadV1 contains score update failure data for Discord.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type ScoreUpdateFailedDiscordPayloadV1 struct {
	GuildID   sharedtypes.GuildID   `json:"guild_id"`
	RoundID   sharedtypes.RoundID   `json:"round_id"`
	UserID    sharedtypes.DiscordID `json:"user_id"`
	Reason    string                `json:"reason"`
	ChannelID string                `json:"channel_id"`
	MessageID string                `json:"message_id"`
}

// ScoreBulkUpdateRequestDiscordPayloadV1 contains bulk score update request data from Discord.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type ScoreBulkUpdateRequestDiscordPayloadV1 struct {
	GuildID   sharedtypes.GuildID                  `json:"guild_id"`
	RoundID   sharedtypes.RoundID                  `json:"round_id"`
	Updates   []ScoreUpdateRequestDiscordPayloadV1 `json:"updates"`
	ChannelID string                               `json:"channel_id"`
	MessageID string                               `json:"message_id"`
}

// ScoreBulkUpdateResponseDiscordPayloadV1 contains bulk score update result data for Discord.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type ScoreBulkUpdateResponseDiscordPayloadV1 struct {
	GuildID        sharedtypes.GuildID     `json:"guild_id"`
	RoundID        sharedtypes.RoundID     `json:"round_id"`
	AppliedCount   int                     `json:"applied_count"`
	FailedCount    int                     `json:"failed_count"`
	TotalRequested int                     `json:"total_requested"`
	UserIDsApplied []sharedtypes.DiscordID `json:"user_ids_applied"`
	ChannelID      string                  `json:"channel_id"`
	MessageID      string                  `json:"message_id"`
}
