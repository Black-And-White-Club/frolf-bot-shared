// Package leaderboardevents contains leaderboard-related domain events.
//
// This file defines the Leaderboard Update Flow - events for updating the
// leaderboard after round finalization.
//
// # Flow Sequence
//
//  1. Round finalized -> LeaderboardRoundFinalizedV1
//  2. Update requested -> LeaderboardUpdateRequestedV1
//  3. Success -> LeaderboardUpdatedV1
//  4. OR Failure -> LeaderboardUpdateFailedV1
//  5. Old leaderboard deactivated -> DeactivateOldLeaderboardV1
//
// # Relationship to Round Module
//
// Leaderboard updates are triggered by round finalization:
//   - RoundFinalizedV1 -> publishes LeaderboardRoundFinalizedV1
//
// # Versioning Strategy
//
// All events include a V1 suffix for future schema evolution.
// Current leaderboard update flow remains V1-only during migration.
package leaderboardevents

import (
	"time"

	sharedevents "github.com/Black-And-White-Club/frolf-bot-shared/events/shared"
	leaderboardtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/leaderboard"
	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
	usertypes "github.com/Black-And-White-Club/frolf-bot-shared/types/user"
)

// =============================================================================
// LEADERBOARD UPDATE FLOW - Event Constants
// =============================================================================

// LeaderboardRoundFinalizedV1 is published when a round is finalized for leaderboard update.
//
// Pattern: Event Notification
// Subject: leaderboard.round.finalized.v1
// Producer: round-service
// Consumers: leaderboard-service (update handler)
// Triggers: LeaderboardUpdateRequestedV1
// Version: v1 (December 2024)
const LeaderboardRoundFinalizedV1 = "leaderboard.round.finalized.v1"

// LeaderboardUpdateRequestedV1 is published when a leaderboard update is requested.
//
// Pattern: Event Notification
// Subject: leaderboard.update.requested.v1
// Producer: leaderboard-service (after round finalization)
// Consumers: leaderboard-service (update processor)
// Triggers: LeaderboardUpdatedV1 OR LeaderboardUpdateFailedV1
// Version: v1 (December 2024)
const LeaderboardUpdateRequestedV1 = "leaderboard.update.requested.v1"

// LeaderboardUpdatedV1 is published when the leaderboard is successfully updated.
//
// Pattern: Event Notification
// Subject: leaderboard.updated.v1
// Producer: leaderboard-service
// Consumers: discord-service (embed update), round-service (tag updates)
// Version: v1 (December 2024)
const LeaderboardUpdatedV1 = "leaderboard.updated.v1"

// LeaderboardUpdateFailedV1 is published when leaderboard update fails.
//
// Pattern: Event Notification
// Subject: leaderboard.update.failed.v1
// Producer: leaderboard-service
// Consumers: monitoring, error handlers
// Version: v1 (December 2024)
const LeaderboardUpdateFailedV1 = "leaderboard.update.failed.v1"

// DeactivateOldLeaderboardV1 is published to deactivate a previous leaderboard.
//
// Pattern: Event Notification
// Subject: leaderboard.deactivate.v1
// Producer: leaderboard-service
// Consumers: leaderboard-service (cleanup handler)
// Version: v1 (December 2024)
const DeactivateOldLeaderboardV1 = "leaderboard.deactivate.v1"

// TagUpdateForScheduledRoundsV1 is published to update tags for scheduled rounds.
//
// Pattern: Event Notification
// Subject: round.tag.update.for.scheduled.rounds.v1
// Producer: leaderboard-service (after update)
// Consumers: round-service (tag update handler)
// Version: v1 (December 2024)
const TagUpdateForScheduledRoundsV1 = "round.tag.update.for.scheduled.rounds.v1"

// =============================================================================
// LEADERBOARD UPDATE FLOW - Payload Types
// =============================================================================

// LeaderboardRoundFinalizedPayloadV1 contains round finalization data for leaderboard.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type LeaderboardRoundFinalizedPayloadV1 struct {
	GuildID               sharedtypes.GuildID `json:"guild_id"`
	RoundID               sharedtypes.RoundID `json:"round_id"`
	SortedParticipantTags TagOrderV1          `json:"sorted_participant_tags"`
}

// TagOrderV1 represents the order of tags.
type TagOrderV1 []string

// TagUpdateForScheduledRoundsPayloadV1 contains tag changes that may affect scheduled rounds.
//
// Notes:
//   - ChangedTags may be empty
//   - Payload size is unbounded
//   - This is an integration event, not a domain entity
//
// Schema History:
//   - v1.0 (January 2026): Initial version
type TagUpdateForScheduledRoundsPayloadV1 struct {
	GuildID     sharedtypes.GuildID                             `json:"guild_id"`
	RoundID     sharedtypes.RoundID                             `json:"round_id,omitempty"`
	Source      string                                          `json:"source,omitempty"`
	UpdatedAt   time.Time                                       `json:"updated_at"`
	ChangedTags map[sharedtypes.DiscordID]sharedtypes.TagNumber `json:"changed_tags,omitempty"`
}

// LeaderboardUpdateRequestedPayloadV1 contains leaderboard update request data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type LeaderboardUpdateRequestedPayloadV1 struct {
	GuildID               sharedtypes.GuildID       `json:"guild_id"`
	RoundID               sharedtypes.RoundID       `json:"round_id"`
	SortedParticipantTags []string                  `json:"sorted_participant_tags"`
	Participants          []RoundParticipantInputV1 `json:"participants,omitempty"`
	Source                string                    `json:"source"`
	UpdateID              string                    `json:"update_id"`
}

// RoundParticipantInputV1 provides raw finish data for ProcessRound-style handling.
// When present, this should be preferred over SortedParticipantTags.
type RoundParticipantInputV1 struct {
	MemberID   sharedtypes.DiscordID `json:"member_id"`
	FinishRank int                   `json:"finish_rank"`
}

// LeaderboardUpdatedPayloadV1 contains leaderboard update success data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type LeaderboardUpdatedPayloadV1 struct {
	GuildID         sharedtypes.GuildID                             `json:"guild_id"`
	LeaderboardID   int64                                           `json:"leaderboard_id"`
	RoundID         sharedtypes.RoundID                             `json:"round_id"`
	LeaderboardData map[sharedtypes.TagNumber]sharedtypes.DiscordID `json:"leaderboard_data"`
	Config          *sharedevents.GuildConfigFragment               `json:"config_fragment,omitempty"`
}

// LeaderboardUpdateFailedPayloadV1 contains leaderboard update failure data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type LeaderboardUpdateFailedPayloadV1 struct {
	GuildID sharedtypes.GuildID `json:"guild_id"`
	RoundID sharedtypes.RoundID `json:"round_id"`
	Reason  string              `json:"reason"`
}

// DeactivateOldLeaderboardPayloadV1 contains leaderboard deactivation data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type DeactivateOldLeaderboardPayloadV1 struct {
	GuildID       sharedtypes.GuildID `json:"guild_id"`
	LeaderboardID int64               `json:"leaderboard_id"`
}

// =============================================================================
// LEADERBOARD RETRIEVAL - Event Constants
// =============================================================================

// GetLeaderboardRequestedV1 is published when leaderboard data is requested.
//
// Pattern: Event Notification
// Subject: leaderboard.get.requested.v1
// Producer: any service needing leaderboard data
// Consumers: leaderboard-service (retrieval handler)
// Triggers: GetLeaderboardResponseV1 OR GetLeaderboardFailedV1
// Version: v1 (December 2024)
const GetLeaderboardRequestedV1 = "leaderboard.get.requested.v1"

// GetLeaderboardResponseV1 is published with the leaderboard data.
//
// Pattern: Event Notification
// Subject: leaderboard.get.response.v1
// Producer: leaderboard-service
// Consumers: requesting service
// Version: v1 (December 2024)
const GetLeaderboardResponseV1 = "leaderboard.get.response.v1"

// GetLeaderboardFailedV1 is published when leaderboard retrieval fails.
//
// Pattern: Event Notification
// Subject: leaderboard.get.failed.v1
// Producer: leaderboard-service
// Consumers: requesting service
// Version: v1 (December 2024)
const GetLeaderboardFailedV1 = "leaderboard.get.failed.v1"

// =============================================================================
// LEADERBOARD RETRIEVAL - Payload Types
// =============================================================================

// GetLeaderboardRequestedPayloadV1 contains leaderboard retrieval request data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type GetLeaderboardRequestedPayloadV1 struct {
	GuildID  sharedtypes.GuildID `json:"guild_id"`
	SeasonID string              `json:"season_id,omitempty"`
}

// GetLeaderboardResponsePayloadV1 contains leaderboard retrieval response data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type GetLeaderboardResponsePayloadV1 struct {
	GuildID     sharedtypes.GuildID                              `json:"guild_id"`
	Leaderboard leaderboardtypes.LeaderboardData                 `json:"leaderboard"`
	Profiles    map[sharedtypes.DiscordID]*usertypes.UserProfile `json:"profiles,omitempty"`
}

// GetLeaderboardFailedPayloadV1 contains leaderboard retrieval failure data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type GetLeaderboardFailedPayloadV1 struct {
	GuildID sharedtypes.GuildID `json:"guild_id"`
	Reason  string              `json:"reason"`
}
