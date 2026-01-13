// Package roundevents contains all round-related domain events.
//
// This file defines the Tag Lookup Flow - events related to
// looking up and managing participant tag numbers.
//
// # Flow Sequence
//
//  1. Tag lookup requested -> RoundTagNumberRequestedV1
//  2. Leaderboard module queried -> LeaderboardGetTagNumberRequestedV1
//  3. Response received -> LeaderboardGetTagNumberResponseV1
//  4. Tag found -> RoundTagNumberFoundV1
//  5. OR Tag not found -> RoundTagNumberNotFoundV1
//
// # Scheduled Round Tag Updates
//
//  1. Tag change detected -> TagUpdateForScheduledRoundsV1
//  2. Rounds updated -> TagsUpdatedForScheduledRoundsV1
//
// # Pattern Reference
//
// This flow follows the Event Notification pattern (Martin Fowler).
//
// # Versioning Strategy
//
// All events include a V1 suffix for future schema evolution.
package roundevents

import (
	roundtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/round"
	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
)

// =============================================================================
// TAG LOOKUP FLOW - Event Constants
// =============================================================================

// RoundTagNumberRequestedV1 is published when a tag number lookup is requested.
//
// Pattern: Event Notification
// Subject: round.tag.number.requested.v1
// Producer: backend-service (participant handler)
// Consumers: backend-service (tag lookup handler)
// Triggers: LeaderboardGetTagNumberRequestedV1
// Version: v1 (December 2024)
const RoundTagNumberRequestedV1 = "round.tag.number.requested.v1"

// RoundTagNumberFoundV1 is published when a tag number is found.
//
// DEPRECATED: This event is superseded by sharedevents.RoundTagLookupFoundV1.
// Use round.tag.lookup.found.v1 (from events/shared/tags.go) instead.
// This constant will be removed in v2.0.
//
// Pattern: Event Notification
// Subject: round.leaderboard.tag.found.v1
// Producer: backend-service (tag lookup handler)
// Consumers: backend-service (participant update handler)
// Version: v1 (December 2024)
const RoundTagNumberFoundV1 = "round.leaderboard.tag.found.v1"

// RoundTagNumberNotFoundV1 is published when a tag number is not found.
//
// DEPRECATED: This event is superseded by sharedevents.RoundTagLookupNotFoundV1.
// Use round.tag.lookup.not.found.v1 (from events/shared/tags.go) instead.
// This constant will be removed in v2.0.
//
// Pattern: Event Notification
// Subject: round.leaderboard.tag.not.found.v1
// Producer: backend-service (tag lookup handler)
// Consumers: backend-service (participant update handler)
// Version: v1 (December 2024)
const RoundTagNumberNotFoundV1 = "round.leaderboard.tag.not.found.v1"

// LeaderboardGetTagNumberRequestedV1 is published to request tag number from leaderboard.
//
// Pattern: Event Notification
// Subject: leaderboard.round.tag.get.by.user.id.requested.v1
// Producer: backend-service (round module)
// Consumers: backend-service (leaderboard module)
// Triggers: LeaderboardGetTagNumberResponseV1
// Version: v1 (December 2024)
const LeaderboardGetTagNumberRequestedV1 = "leaderboard.round.tag.get.by.user.id.requested.v1"

// LeaderboardGetTagNumberResponseV1 is published with tag number response.
//
// Pattern: Event Notification
// Subject: round.get.tag.number.response.v1
// Producer: backend-service (leaderboard module)
// Consumers: backend-service (round module)
// Version: v1 (December 2024)
const LeaderboardGetTagNumberResponseV1 = "round.get.tag.number.response.v1"

// TagUpdateForScheduledRoundsV1 is published to update tags for scheduled rounds.
//
// Pattern: Event Notification
// Subject: round.tag.update.for.scheduled.rounds.v1
// Producer: backend-service (leaderboard module)
// Consumers: backend-service (round module)
// Version: v1 (December 2024)
const TagUpdateForScheduledRoundsV1 = "round.tag.update.for.scheduled.rounds.v1"

// TagsUpdatedForScheduledRoundsV1 is published when scheduled round tags are updated.
//
// Pattern: Event Notification
// Subject: round.tags.updated.for.scheduled.rounds.v1
// Producer: backend-service (round module)
// Consumers: discord-service (embed update handler)
// Version: v1 (December 2024)
const TagsUpdatedForScheduledRoundsV1 = "round.tags.updated.for.scheduled.rounds.v1"

// =============================================================================
// TAG LOOKUP FLOW - Payload Types
// =============================================================================

// TagLookupRequestPayloadV1 contains tag lookup request data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type TagLookupRequestPayloadV1 struct {
	GuildID          sharedtypes.GuildID   `json:"guild_id"`
	UserID           sharedtypes.DiscordID `json:"user_id"`
	RoundID          sharedtypes.RoundID   `json:"round_id"`
	Response         roundtypes.Response   `json:"response"`
	OriginalResponse roundtypes.Response   `json:"original_response"`
	JoinedLate       *bool                 `json:"joined_late,omitempty"`
}

// TagNumberRequestPayloadV1 contains tag number request data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type TagNumberRequestPayloadV1 struct {
	GuildID sharedtypes.GuildID   `json:"guild_id"`
	UserID  sharedtypes.DiscordID `json:"user_id"`
	RoundID sharedtypes.RoundID   `json:"round_id"`
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

// GetTagNumberResponsePayloadV1 contains tag number response data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type GetTagNumberResponsePayloadV1 struct {
	GuildID   sharedtypes.GuildID    `json:"guild_id"`
	UserID    sharedtypes.DiscordID  `json:"user_id"`
	TagNumber *sharedtypes.TagNumber `json:"tag_number"`
	Error     string                 `json:"error,omitempty"`
	RoundID   sharedtypes.RoundID    `json:"round_id"`
}

// ScheduledRoundTagUpdatePayloadV1 contains scheduled round tag update data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type ScheduledRoundTagUpdatePayloadV1 struct {
	GuildID     sharedtypes.GuildID                              `json:"guild_id"`
	ChangedTags map[sharedtypes.DiscordID]*sharedtypes.TagNumber `json:"changed_tags"`
}

// RoundUpdateInfoV1 contains round update information.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoundUpdateInfoV1 struct {
	GuildID             sharedtypes.GuildID      `json:"guild_id"`
	RoundID             sharedtypes.RoundID      `json:"round_id"`
	EventMessageID      string                   `json:"event_message_id"`
	Title               roundtypes.Title         `json:"title"`
	StartTime           *sharedtypes.StartTime   `json:"start_time"`
	Location            *roundtypes.Location     `json:"location"`
	UpdatedParticipants []roundtypes.Participant `json:"updated_participants"`
	ParticipantsChanged int                      `json:"participants_changed"`
}

// UpdateSummaryV1 contains update summary information.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type UpdateSummaryV1 struct {
	GuildID              sharedtypes.GuildID `json:"guild_id"`
	TotalRoundsProcessed int                 `json:"total_rounds_processed"`
	RoundsUpdated        int                 `json:"rounds_updated"`
	ParticipantsUpdated  int                 `json:"participants_updated"`
}

// TagsUpdatedForScheduledRoundsPayloadV1 contains scheduled rounds update result.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type TagsUpdatedForScheduledRoundsPayloadV1 struct {
	GuildID       sharedtypes.GuildID `json:"guild_id"`
	UpdatedRounds []RoundUpdateInfoV1 `json:"updated_rounds"`
	Summary       UpdateSummaryV1     `json:"summary"`
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
