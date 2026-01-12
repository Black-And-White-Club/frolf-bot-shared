// Package leaderboard contains Discord-specific leaderboard events.
//
// This file defines the Discord Leaderboard Retrieval Flow - events specific to
// displaying leaderboards through Discord embeds.
//
// # Flow Sequence
//
//  1. User requests leaderboard -> LeaderboardRetrieveRequestV1
//  2. Leaderboard data received -> LeaderboardRetrievedV1
//
// # Relationship to Domain Events
//
// These Discord events wrap/trigger domain events in events/leaderboard/:
//   - LeaderboardRetrieveRequestV1 -> publishes GetLeaderboardRequest (domain)
//   - LeaderboardRetrievedV1 <- subscribes to GetLeaderboardResponse (domain)
//
// # Versioning Strategy
//
// All events include a V1 suffix for future schema evolution.
package leaderboard

import (
	leaderboardtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/leaderboard"
	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
)

// =============================================================================
// DISCORD LEADERBOARD RETRIEVAL FLOW - Event Constants
// =============================================================================

// LeaderboardRetrieveRequestV1 is published when a user requests to view the leaderboard.
//
// Pattern: Event Notification
// Subject: discord.leaderboard.retrieve.request.v1
// Producer: discord-service (command handler)
// Consumers: discord-service (leaderboard handler)
// Triggers: Domain GetLeaderboardRequest
// Version: v1 (December 2024)
const LeaderboardRetrieveRequestV1 = "discord.leaderboard.retrieve.request.v1"

// LeaderboardRetrievedV1 is published when leaderboard data is ready for display.
//
// Pattern: Event Notification
// Subject: discord.leaderboard.retrieved.v1
// Producer: discord-service (after domain response)
// Consumers: discord-service (embed builder)
// Triggers: Discord embed with leaderboard
// Version: v1 (December 2024)
const LeaderboardRetrievedV1 = "discord.leaderboard.retrieved.v1"

// =============================================================================
// DISCORD LEADERBOARD RETRIEVAL FLOW - Payload Types
// =============================================================================

// LeaderboardRetrieveRequestPayloadV1 contains leaderboard request data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type LeaderboardRetrieveRequestPayloadV1 struct {
	UserID    sharedtypes.DiscordID `json:"user_id"`
	ChannelID string                `json:"channel_id"`
	MessageID string                `json:"discord_message_id"`
	GuildID   string                `json:"guild_id"`
}

// LeaderboardRetrievedPayloadV1 contains retrieved leaderboard data for Discord.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type LeaderboardRetrievedPayloadV1 struct {
	Leaderboard []leaderboardtypes.LeaderboardEntry `json:"leaderboard"`
	ChannelID   string                              `json:"channel_id"`
	MessageID   string                              `json:"discord_message_id"`
	GuildID     string                              `json:"guild_id"`
}
