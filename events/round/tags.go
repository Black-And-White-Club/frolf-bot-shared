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

// NOTE: Cross-module tag lookup and scheduled tag update events are defined in:
// - events/shared/tag_number_lookup.go
// - events/shared/tags.go

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

// NOTE: Cross-module tag lookup payloads are defined in events/shared/tag_number_lookup.go.

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
