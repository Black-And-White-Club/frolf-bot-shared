// Package userevents contains user-related domain events.
//
// This file defines the UDisc Identity Flow - events for managing UDisc
// identity associations and match confirmations.
//
// # Flow Sequences
//
// ## UDisc Identity Update Flow
//  1. Request -> UpdateUDiscIdentityRequestedV1
//  2. Success -> UDiscIdentityUpdatedV1
//  3. OR Failure -> UDiscIdentityUpdateFailedV1
//
// ## UDisc Match Confirmation Flow
//  1. Match required -> UDiscMatchConfirmationRequiredV1
//  2. Admin confirms -> UDiscMatchConfirmedV1
//
// # Versioning Strategy
//
// All events include a V1 suffix for future schema evolution.
package userevents

import (
	"time"

	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
)

// =============================================================================
// UDISC IDENTITY UPDATE FLOW - Event Constants
// =============================================================================

// UpdateUDiscIdentityRequestedV1 is published when UDisc identity update is requested.
//
// Pattern: Event Notification
// Subject: user.udisc.identity.update.requested.v1
// Producer: discord-service (command handler)
// Consumers: user-service (identity handler)
// Triggers: UDiscIdentityUpdatedV1 OR UDiscIdentityUpdateFailedV1
// Version: v1 (December 2024)
const UpdateUDiscIdentityRequestedV1 = "user.udisc.identity.update.requested.v1"

// UDiscIdentityUpdatedV1 is published when UDisc identity is updated.
//
// Pattern: Event Notification
// Subject: user.udisc.identity.updated.v1
// Producer: user-service
// Consumers: discord-service (confirmation)
// Version: v1 (December 2024)
const UDiscIdentityUpdatedV1 = "user.udisc.identity.updated.v1"

// UDiscIdentityUpdateFailedV1 is published when UDisc identity update fails.
//
// Pattern: Event Notification
// Subject: user.udisc.identity.update.failed.v1
// Producer: user-service
// Consumers: discord-service (error handler)
// Version: v1 (December 2024)
const UDiscIdentityUpdateFailedV1 = "user.udisc.identity.update.failed.v1"

// =============================================================================
// UDISC MATCH CONFIRMATION FLOW - Event Constants
// =============================================================================

// UDiscMatchConfirmationRequiredV1 is published when player matches require confirmation.
//
// Pattern: Event Notification
// Subject: user.udisc.match.confirmation.required.v1
// Producer: round-service (import handler)
// Consumers: discord-service (confirmation UI)
// Triggers: UDiscMatchConfirmedV1
// Version: v1 (December 2024)
const UDiscMatchConfirmationRequiredV1 = "user.udisc.match.confirmation.required.v1"

// UDiscMatchConfirmedV1 is published when an admin confirms player matches.
//
// Pattern: Event Notification
// Subject: user.udisc.match.confirmed.v1
// Producer: discord-service (confirmation handler)
// Consumers: round-service (import completion)
// Version: v1 (December 2024)
const UDiscMatchConfirmedV1 = "user.udisc.match.confirmed.v1"

// =============================================================================
// UDISC FLOW - Payload Types
// =============================================================================

// -----------------------------------------------------------------------------
// UDisc Identity Update Payloads
// -----------------------------------------------------------------------------

// UpdateUDiscIdentityRequestedPayloadV1 contains UDisc identity update request data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type UpdateUDiscIdentityRequestedPayloadV1 struct {
	GuildID  sharedtypes.GuildID   `json:"guild_id"`
	UserID   sharedtypes.DiscordID `json:"user_id"`
	Username *string               `json:"username,omitempty"`
	Name     *string               `json:"name,omitempty"`
}

// UDiscIdentityUpdatedPayloadV1 contains UDisc identity update success data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type UDiscIdentityUpdatedPayloadV1 struct {
	GuildID  sharedtypes.GuildID   `json:"guild_id"`
	UserID   sharedtypes.DiscordID `json:"user_id"`
	Username *string               `json:"username,omitempty"`
	Name     *string               `json:"name,omitempty"`
}

// UDiscIdentityUpdateFailedPayloadV1 contains UDisc identity update failure data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type UDiscIdentityUpdateFailedPayloadV1 struct {
	GuildID sharedtypes.GuildID   `json:"guild_id"`
	UserID  sharedtypes.DiscordID `json:"user_id"`
	Reason  string                `json:"reason"`
}

// -----------------------------------------------------------------------------
// UDisc Match Confirmation Payloads
// -----------------------------------------------------------------------------

// UDiscMatchConfirmationRequiredPayloadV1 contains match confirmation request data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type UDiscMatchConfirmationRequiredPayloadV1 struct {
	GuildID          sharedtypes.GuildID   `json:"guild_id"`
	RoundID          sharedtypes.RoundID   `json:"round_id"`
	ImportID         string                `json:"import_id"`
	UserID           sharedtypes.DiscordID `json:"user_id"`
	ChannelID        string                `json:"channel_id"`
	UnmatchedPlayers []string              `json:"unmatched_players"`
	Timestamp        time.Time             `json:"timestamp"`
}

// UDiscConfirmedMappingV1 represents a resolved player match.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type UDiscConfirmedMappingV1 struct {
	PlayerName    string                `json:"player_name"`
	DiscordUserID sharedtypes.DiscordID `json:"discord_user_id"`
}

// UDiscMatchConfirmedPayloadV1 contains confirmed player match data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type UDiscMatchConfirmedPayloadV1 struct {
	GuildID      sharedtypes.GuildID       `json:"guild_id"`
	RoundID      sharedtypes.RoundID       `json:"round_id"`
	ImportID     string                    `json:"import_id"`
	UserID       sharedtypes.DiscordID     `json:"user_id"`
	ChannelID    string                    `json:"channel_id"`
	Timestamp    time.Time                 `json:"timestamp"`
	Mappings     []UDiscConfirmedMappingV1 `json:"mappings"`
	ParsedScores interface{}               `json:"parsed_scores,omitempty"`
}
