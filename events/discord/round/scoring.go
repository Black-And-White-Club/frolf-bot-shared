// Package round contains Discord-specific round events.
//
// This file defines the Discord Round Scoring Flow - events specific to
// submitting and updating scores through Discord interactions.
//
// # Flow Sequence
//
//  1. User submits score -> RoundScoreUpdateRequestDiscordV1
//  2. Success notification -> RoundParticipantScoreUpdatedDiscordV1
//
// # Relationship to Domain Events
//
// These Discord events wrap/trigger domain events in events/round/scoring.go:
//   - RoundScoreUpdateRequestDiscordV1 -> publishes RoundScoreUpdateRequestedV1 (domain)
//   - RoundParticipantScoreUpdatedDiscordV1 <- subscribes to RoundParticipantScoreUpdatedV1 (domain)
//
// # Versioning Strategy
//
// All events include a V1 suffix for future schema evolution.
package round

import (
	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
)

// =============================================================================
// DISCORD ROUND SCORING FLOW - Event Constants
// =============================================================================

// RoundScoreUpdateRequestDiscordV1 is published when a user submits a score via Discord.
//
// Pattern: Event Notification
// Subject: discord.round.score.update.request.v1
// Producer: discord-service (score input handler)
// Consumers: discord-service (score handler)
// Triggers: Domain event RoundScoreUpdateRequestedV1
// Version: v1 (December 2024)
const RoundScoreUpdateRequestDiscordV1 = "discord.round.score.update.request.v1"

// RoundParticipantScoreUpdatedDiscordV1 is published to notify Discord that a score was updated.
//
// Pattern: Event Notification
// Subject: discord.round.participant.score.updated.v1
// Producer: discord-service (after receiving domain event)
// Consumers: discord-service (embed updater)
// Triggers: Discord embed updated with new score
// Version: v1 (December 2024)
const RoundParticipantScoreUpdatedDiscordV1 = "discord.round.participant.score.updated.v1"

// =============================================================================
// DISCORD ROUND SCORING FLOW - Payload Types
// =============================================================================

// RoundScoreUpdateRequestDiscordPayloadV1 contains score update request data from Discord.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoundScoreUpdateRequestDiscordPayloadV1 struct {
	RoundID   sharedtypes.RoundID   `json:"round_id"`
	UserID    sharedtypes.DiscordID `json:"user_id"` // User submitting the score
	Score     sharedtypes.Score     `json:"score"`
	ChannelID string                `json:"channel_id"`
	MessageID string                `json:"message_id"`
	GuildID   string                `json:"guild_id"`
}

// RoundParticipantScoreUpdatedDiscordPayloadV1 contains score update confirmation for Discord.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoundParticipantScoreUpdatedDiscordPayloadV1 struct {
	RoundID   sharedtypes.RoundID   `json:"round_id"`
	UserID    sharedtypes.DiscordID `json:"user_id"`
	Score     sharedtypes.Score     `json:"score"`
	ChannelID string                `json:"channel_id"`
	MessageID string                `json:"message_id"`
	GuildID   string                `json:"guild_id"`
}
