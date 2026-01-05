// Package sharedevents contains cross-module shared events.
//
// This file defines the Shared Tag Lookup Flow - events for cross-module
// tag lookup operations between round and leaderboard modules.
//
// # Flow Sequences
//
// ## Round Tag Lookup Flow
//  1. Request -> RoundTagLookupRequestedV1
//  2. Found -> RoundTagLookupFoundV1
//  3. OR Not found -> RoundTagLookupNotFoundV1
//
// ## Discord Tag Lookup Flow
//  1. Request -> DiscordTagLookupRequestedV1
//  2. Success -> DiscordTagLookupSucceededV1
//  3. OR Not found -> DiscordTagLookupNotFoundV1
//  4. OR Failed -> DiscordTagLookupFailedV1
//
// # Cross-Module Communication
//
// These events enable communication between modules:
//   - Round module publishes tag lookup requests
//   - Leaderboard module responds with tag data
//
// # Versioning Strategy
//
// All events include a V1 suffix for future schema evolution.
package sharedevents

import (
	roundtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/round"
	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
)

// =============================================================================
// ROUND TAG LOOKUP FLOW - Event Constants
// =============================================================================

// RoundTagLookupRequestedV1 is published when round needs a tag lookup.
//
// Pattern: Event Notification
// Subject: round.tag.lookup.requested.v1
// Producer: round-service
// Consumers: leaderboard-service (lookup handler)
// Triggers: RoundTagLookupFoundV1 OR RoundTagLookupNotFoundV1
// Version: v1 (December 2024)
const RoundTagLookupRequestedV1 = "round.tag.lookup.requested.v1"

// RoundTagLookupFoundV1 is published when a tag is found for a round request.
//
// Pattern: Event Notification
// Subject: round.tag.lookup.found.v1
// Producer: leaderboard-service
// Consumers: round-service
// Version: v1 (December 2024)
const RoundTagLookupFoundV1 = "round.tag.lookup.found.v1"

// RoundTagLookupNotFoundV1 is published when a tag is not found.
//
// Pattern: Event Notification
// Subject: round.tag.lookup.not.found.v1
// Producer: leaderboard-service
// Consumers: round-service
// Version: v1 (December 2024)
const RoundTagLookupNotFoundV1 = "round.tag.lookup.not.found.v1"

// TagUpdateForScheduledRoundsV1 is published to update tags for scheduled rounds.
//
// Pattern: Event Notification
// Subject: round.tag.update.for.scheduled.rounds.v1
// Producer: leaderboard-service
// Consumers: round-service
// Version: v1 (December 2024)
const TagUpdateForScheduledRoundsV1 = "round.tag.update.for.scheduled.rounds.v1"

// =============================================================================
// DISCORD TAG LOOKUP FLOW - Event Constants
// =============================================================================

// DiscordTagLookupRequestedV1 is published for Discord-initiated tag lookups.
//
// Pattern: Event Notification
// Subject: leaderboard.tag.lookup.by.user.id.requested.v1
// Producer: discord-service
// Consumers: leaderboard-service (lookup handler)
// Triggers: DiscordTagLookupSucceededV1 OR DiscordTagLookupNotFoundV1
// Version: v1 (December 2024)
const DiscordTagLookupRequestedV1 = "leaderboard.tag.lookup.by.user.id.requested.v1"

// DiscordTagLookupSucceededV1 is published when Discord tag lookup succeeds.
//
// Pattern: Event Notification
// Subject: discord.leaderboard.tag.lookup.by.user.id.success.v1
// Producer: leaderboard-service
// Consumers: discord-service
// Version: v1 (December 2024)
const DiscordTagLookupSucceededV1 = "discord.leaderboard.tag.lookup.by.user.id.success.v1"

// DiscordTagLookupNotFoundV1 is published when Discord tag lookup finds nothing.
//
// Pattern: Event Notification
// Subject: discord.leaderboard.tag.lookup.by.user.id.not.found.v1
// Producer: leaderboard-service
// Consumers: discord-service
// Version: v1 (December 2024)
const DiscordTagLookupNotFoundV1 = "discord.leaderboard.tag.lookup.by.user.id.not.found.v1"

// DiscordTagLookupFailedV1 is published when Discord tag lookup fails.
//
// Pattern: Event Notification
// Subject: discord.leaderboard.tag.lookup.by.user.id.failed.v1
// Producer: leaderboard-service
// Consumers: discord-service
// Version: v1 (December 2024)
const DiscordTagLookupFailedV1 = "discord.leaderboard.tag.lookup.by.user.id.failed.v1"

// =============================================================================
// BATCH TAG ASSIGNMENT FLOW - Event Constants
// =============================================================================

// LeaderboardBatchTagAssignmentRequestedV1 is published for batch tag assignment.
//
// Pattern: Event Notification
// Subject: leaderboard.batch.tag.assignment.requested.v1
// Producer: score-service
// Consumers: leaderboard-service
// Version: v1 (December 2024)
const LeaderboardBatchTagAssignmentRequestedV1 = "leaderboard.batch.tag.assignment.requested.v1"

// =============================================================================
// TAG LOOKUP FLOW - Payload Types
// =============================================================================

// -----------------------------------------------------------------------------
// Round Tag Lookup Payloads
// -----------------------------------------------------------------------------

// RoundTagLookupRequestedPayloadV1 contains tag lookup request data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoundTagLookupRequestedPayloadV1 struct {
	ScopedGuildID
	UserID     sharedtypes.DiscordID `json:"user_id"`
	RoundID    sharedtypes.RoundID   `json:"round_id"`
	Response   roundtypes.Response   `json:"response"`
	JoinedLate *bool                 `json:"joined_late,omitempty"`
}

// Deprecated: Use RoundTagLookupRequestedPayloadV1.
// This alias exists to smooth over older code that accidentally double-suffixed versions (e.g., V1PayloadV1).
type RoundTagLookupRequestedV1PayloadV1 = RoundTagLookupRequestedPayloadV1

// RoundTagLookupResultPayloadV1 contains tag lookup result data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoundTagLookupResultPayloadV1 struct {
	ScopedGuildID
	UserID             sharedtypes.DiscordID  `json:"user_id"`
	RoundID            sharedtypes.RoundID    `json:"round_id"`
	TagNumber          *sharedtypes.TagNumber `json:"tag_number,omitempty"`
	Found              bool                   `json:"found"`
	OriginalResponse   roundtypes.Response    `json:"original_response"`
	OriginalJoinedLate *bool                  `json:"original_joined_late,omitempty"`
	Error              string                 `json:"error,omitempty"`
}

// RoundTagLookupFailedPayloadV1 contains tag lookup failure data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoundTagLookupFailedPayloadV1 struct {
	ScopedGuildID
	UserID  sharedtypes.DiscordID `json:"user_id"`
	RoundID sharedtypes.RoundID   `json:"round_id"`
	Reason  string                `json:"reason"`
}

// Deprecated: Use RoundTagLookupFailedPayloadV1.
// This alias exists to smooth over older code that accidentally double-suffixed versions (e.g., V1PayloadV1).
type RoundTagLookupFailedV1PayloadV1 = RoundTagLookupFailedPayloadV1

// -----------------------------------------------------------------------------
// Discord Tag Lookup Payloads
// -----------------------------------------------------------------------------

// DiscordTagLookupRequestedPayloadV1 contains Discord tag lookup request data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type DiscordTagLookupRequestedPayloadV1 struct {
	ScopedGuildID
	RequestingUserID sharedtypes.DiscordID `json:"requester_user_id"`
	UserID           sharedtypes.DiscordID `json:"user_id"`
}

// DiscordTagLookupResultPayloadV1 contains Discord tag lookup result data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type DiscordTagLookupResultPayloadV1 struct {
	ScopedGuildID
	RequestingUserID sharedtypes.DiscordID  `json:"requester_user_id"`
	UserID           sharedtypes.DiscordID  `json:"user_id"`
	TagNumber        *sharedtypes.TagNumber `json:"tag_number,omitempty"`
	Found            bool                   `json:"found"`
}

// DiscordTagLookupFailedPayloadV1 contains Discord tag lookup failure data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type DiscordTagLookupFailedPayloadV1 struct {
	ScopedGuildID
	RequestingUserID sharedtypes.DiscordID `json:"requester_user_id"`
	UserID           sharedtypes.DiscordID `json:"user_id"`
	Reason           string                `json:"reason"`
}

// -----------------------------------------------------------------------------
// Batch Tag Assignment Payloads
// -----------------------------------------------------------------------------

// TagAssignmentInfoV1 represents a user and the tag number they should be assigned.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type TagAssignmentInfoV1 struct {
	UserID    sharedtypes.DiscordID `json:"user_id"`
	TagNumber sharedtypes.TagNumber `json:"tag_number"`
}

// BatchTagAssignmentRequestedPayloadV1 contains batch tag assignment request data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type BatchTagAssignmentRequestedPayloadV1 struct {
	ScopedGuildID
	RequestingUserID sharedtypes.DiscordID `json:"requester_user_id"`
	BatchID          string                `json:"batch_id"`
	Assignments      []TagAssignmentInfoV1 `json:"assignments"`
}
