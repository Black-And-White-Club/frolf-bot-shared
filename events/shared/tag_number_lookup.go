// Package sharedevents contains cross-module shared events.
//
// This file defines shared tag number lookup events used between round, leaderboard, and discord modules.
package sharedevents

import (
	roundtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/round"
	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
)

// =============================================================================
// TAG NUMBER LOOKUP FLOW - Shared Event Constants
// =============================================================================

// GetTagByUserIDRequestedV1 is published when tag lookup by user ID is requested.
//
// Pattern: Event Notification
// Subject: leaderboard.tag.get.by.user.id.requested.v1
// Producer: round-service, discord-service
// Consumers: leaderboard-service (lookup handler)
// Triggers: GetTagNumberResponseV1 OR GetTagNumberFailedV1
// Version: v1 (December 2024)
const GetTagByUserIDRequestedV1 = "leaderboard.tag.get.by.user.id.requested.v1"

// RoundGetTagByUserIDRequestedV1 is published for round-specific tag lookup.
//
// Pattern: Event Notification
// Subject: leaderboard.round.tag.get.by.user.id.requested.v1
// Producer: round-service
// Consumers: leaderboard-service (lookup handler)
// Version: v1 (December 2024)
const RoundGetTagByUserIDRequestedV1 = "leaderboard.round.tag.get.by.user.id.requested.v1"

// GetTagNumberResponseV1 is published with the tag number result.
//
// Pattern: Event Notification
// Subject: leaderboard.tag.get.response.v1
// Producer: leaderboard-service
// Consumers: requesting service
// Version: v1 (December 2024)
const GetTagNumberResponseV1 = "leaderboard.tag.get.response.v1"

// GetTagNumberFailedV1 is published when tag lookup fails.
//
// Pattern: Event Notification
// Subject: leaderboard.tag.get.failed.v1
// Producer: leaderboard-service
// Consumers: requesting service
// Version: v1 (December 2024)
const GetTagNumberFailedV1 = "leaderboard.tag.get.failed.v1"

// RoundTagNumberFoundV1 is published when a round-specific tag is found.
//
// DEPRECATED: This event is superseded by RoundTagLookupFoundV1.
// Use round.tag.lookup.found.v1 (from events/shared/tags.go) instead.
//
// Pattern: Event Notification
// Subject: round.leaderboard.tag.found.v1
// Producer: leaderboard-service
// Consumers: round-service
// Version: v1 (December 2024)
const RoundTagNumberFoundV1 = "round.leaderboard.tag.found.v1"

// RoundTagNumberNotFoundV1 is published when a round-specific tag is not found.
//
// DEPRECATED: This event is superseded by RoundTagLookupNotFoundV1.
// Use round.tag.lookup.not.found.v1 (from events/shared/tags.go) instead.
//
// Pattern: Event Notification
// Subject: round.leaderboard.tag.not.found.v1
// Producer: leaderboard-service
// Consumers: round-service
// Version: v1 (December 2024)
const RoundTagNumberNotFoundV1 = "round.leaderboard.tag.not.found.v1"

// =============================================================================
// TAG NUMBER LOOKUP FLOW - Shared Payload Types
// =============================================================================

// TagNumberRequestPayloadV1 contains tag number lookup request data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type TagNumberRequestPayloadV1 struct {
	GuildID sharedtypes.GuildID   `json:"guild_id"`
	UserID  sharedtypes.DiscordID `json:"user_id"`
	RoundID sharedtypes.RoundID   `json:"round_id"`
}

// UserTagNumberRequestPayloadV1 contains User tag number lookup request data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type UserTagNumberRequestPayloadV1 struct {
	GuildID   sharedtypes.GuildID    `json:"guild_id"`
	TagNumber *sharedtypes.TagNumber `json:"tag_number"`
	UserID    sharedtypes.DiscordID  `json:"user_id"`
}

// GetTagNumberResponsePayloadV1 contains tag number lookup response data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type GetTagNumberResponsePayloadV1 struct {
	GuildID   sharedtypes.GuildID    `json:"guild_id"`
	TagNumber *sharedtypes.TagNumber `json:"tag_number"`
	UserID    sharedtypes.DiscordID  `json:"user_id"`
	RoundID   sharedtypes.RoundID    `json:"round_id"`
	Found     bool                   `json:"found"`
}

// UserTagNumberResponsePayloadV1 contains User tag number lookup response data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type UserTagNumberResponsePayloadV1 struct {
	GuildID   sharedtypes.GuildID    `json:"guild_id"`
	TagNumber *sharedtypes.TagNumber `json:"tag_number"`
	UserID    sharedtypes.DiscordID  `json:"user_id"`
}

// GetTagNumberFailedPayloadV1 contains tag number lookup failure data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type GetTagNumberFailedPayloadV1 struct {
	GuildID sharedtypes.GuildID `json:"guild_id"`
	Reason  string              `json:"reason"`
}

// RoundTagNumberFoundPayloadV1 contains found tag number data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoundTagNumberFoundPayloadV1 struct {
	GuildID            sharedtypes.GuildID    `json:"guild_id"`
	RoundID            sharedtypes.RoundID    `json:"round_id"`
	UserID             sharedtypes.DiscordID  `json:"user_id"`
	TagNumber          *sharedtypes.TagNumber `json:"tag_number"`
	OriginalResponse   roundtypes.Response    `json:"original_response"`
	OriginalJoinedLate *bool                  `json:"original_joined_late,omitempty"`
}

// RoundTagNumberNotFoundPayloadV1 contains tag not found data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoundTagNumberNotFoundPayloadV1 struct {
	GuildID            sharedtypes.GuildID   `json:"guild_id"`
	RoundID            sharedtypes.RoundID   `json:"round_id"`
	UserID             sharedtypes.DiscordID `json:"user_id"`
	OriginalResponse   roundtypes.Response   `json:"original_response"`
	OriginalJoinedLate *bool                 `json:"original_joined_late,omitempty"`
	Reason             string                `json:"reason,omitempty"`
}

// =============================================================================
// Interface Implementations for V1 Payloads
// =============================================================================

// GetRoundID implements ParticipantUpdatePayload for RoundTagNumberNotFoundPayloadV1.
func (p *RoundTagNumberNotFoundPayloadV1) GetRoundID() sharedtypes.RoundID {
	return p.RoundID
}

// GetUserID implements ParticipantUpdatePayload for RoundTagNumberNotFoundPayloadV1.
func (p *RoundTagNumberNotFoundPayloadV1) GetUserID() sharedtypes.DiscordID {
	return p.UserID
}

// GetTagNumber implements ParticipantUpdatePayload for RoundTagNumberNotFoundPayloadV1.
func (p *RoundTagNumberNotFoundPayloadV1) GetTagNumber() *sharedtypes.TagNumber {
	return nil // No tag number found
}

// GetJoinedLate implements ParticipantUpdatePayload for RoundTagNumberNotFoundPayloadV1.
func (p *RoundTagNumberNotFoundPayloadV1) GetJoinedLate() *bool {
	return p.OriginalJoinedLate
}

// GetRoundID implements ParticipantUpdatePayload for RoundTagNumberFoundPayloadV1.
func (p *RoundTagNumberFoundPayloadV1) GetRoundID() sharedtypes.RoundID {
	return p.RoundID
}

// GetUserID implements ParticipantUpdatePayload for RoundTagNumberFoundPayloadV1.
func (p *RoundTagNumberFoundPayloadV1) GetUserID() sharedtypes.DiscordID {
	return p.UserID
}

// GetTagNumber implements ParticipantUpdatePayload for RoundTagNumberFoundPayloadV1.
func (p *RoundTagNumberFoundPayloadV1) GetTagNumber() *sharedtypes.TagNumber {
	return p.TagNumber
}

// GetJoinedLate implements ParticipantUpdatePayload for RoundTagNumberFoundPayloadV1.
func (p *RoundTagNumberFoundPayloadV1) GetJoinedLate() *bool {
	return p.OriginalJoinedLate
}

// ScheduledRoundTagUpdatePayloadV1 contains scheduled round tag update data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type ScheduledRoundTagUpdatePayloadV1 struct {
	GuildID     sharedtypes.GuildID                              `json:"guild_id"`
	ChangedTags map[sharedtypes.DiscordID]*sharedtypes.TagNumber `json:"changed_tags"`
}
