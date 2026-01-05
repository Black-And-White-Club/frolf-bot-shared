// Package round contains Discord-specific round events.
//
// This file defines the Discord Round Participant Flow - events specific to
// participants joining or leaving rounds through Discord button interactions.
//
// # Flow Sequence
//
//  1. User clicks join button -> RoundParticipantJoinRequestDiscordV1
//  2. Success notification -> RoundParticipantJoinedDiscordV1
//
// # Relationship to Domain Events
//
// These Discord events wrap/trigger domain events in events/round/participants.go:
//   - RoundParticipantJoinRequestDiscordV1 -> publishes RoundParticipantJoinRequestedV1 (domain)
//   - RoundParticipantJoinedDiscordV1 <- subscribes to RoundParticipantJoinedV1 (domain)
//
// # Versioning Strategy
//
// All events include a V1 suffix for future schema evolution.
package round

import (
	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
)

// =============================================================================
// DISCORD ROUND PARTICIPANT FLOW - Event Constants
// =============================================================================

// RoundParticipantJoinRequestDiscordV1 is published when a user clicks the join button.
//
// Pattern: Event Notification
// Subject: discord.round.participant.join.request.v1
// Producer: discord-service (button handler)
// Consumers: discord-service (participant handler)
// Triggers: Domain event RoundParticipantJoinRequestedV1
// Version: v1 (December 2024)
const RoundParticipantJoinRequestDiscordV1 = "discord.round.participant.join.request.v1"

// RoundParticipantJoinedDiscordV1 is published to notify Discord that a participant joined.
//
// Pattern: Event Notification
// Subject: discord.round.participant.joined.v1
// Producer: discord-service (after receiving domain event)
// Consumers: discord-service (embed updater)
// Triggers: Discord embed updated with new participant
// Version: v1 (December 2024)
const RoundParticipantJoinedDiscordV1 = "discord.round.participant.joined.v1"

// =============================================================================
// DISCORD ROUND PARTICIPANT FLOW - Payload Types
// =============================================================================

// RoundParticipantJoinRequestDiscordPayloadV1 contains join request data from Discord.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoundParticipantJoinRequestDiscordPayloadV1 struct {
	RoundID    sharedtypes.RoundID   `json:"round_id"`
	UserID     sharedtypes.DiscordID `json:"user_id"`
	ChannelID  string                `json:"channel_id"`
	JoinedLate *bool                 `json:"joined_late,omitempty"`
	GuildID    string                `json:"guild_id"`
}

// RoundParticipantJoinedDiscordPayloadV1 contains successful join data for Discord.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoundParticipantJoinedDiscordPayloadV1 struct {
	RoundID   sharedtypes.RoundID   `json:"round_id"`
	UserID    sharedtypes.DiscordID `json:"user_id"`
	TagNumber sharedtypes.TagNumber `json:"tag_number"`
	ChannelID string                `json:"channel_id"`
	GuildID   string                `json:"guild_id"`
}
