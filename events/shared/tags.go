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
	discordleaderboard "github.com/Black-And-White-Club/frolf-bot-shared/events/discord/leaderboard"
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

// -----------------------------------------------------------------------------
// DEPRECATED: Discord-prefixed lookup response topics
// -----------------------------------------------------------------------------
// These constants are deprecated aliases that remain for compatibility while the
// leaderboard-owned canonical topics are adopted. They are now defined in the
// Discord-specific package so that discord-prefixed subjects live under
// `events/discord` per ownership rules.
const (
	// DiscordTagLookupSucceededV1 is the deprecated alias for the Discord-prefixed success topic.
	DiscordTagLookupSucceededV1 = discordleaderboard.LeaderboardTagLookupSucceededV1

	// DiscordTagLookupNotFoundV1 is the deprecated alias for the Discord-prefixed not-found topic.
	DiscordTagLookupNotFoundV1 = discordleaderboard.LeaderboardTagLookupNotFoundV1

	// DiscordTagLookupFailedV1 is the deprecated alias for the Discord-prefixed failed topic.
	DiscordTagLookupFailedV1 = discordleaderboard.LeaderboardTagLookupFailedV1
)

// -----------------------------------------------------------------------------
// Canonical Leaderboard-owned Topics (for backend responses)
// -----------------------------------------------------------------------------

// LeaderboardTagLookupRequestedV1 is the canonical request topic for tag lookups
// where the consumer is the leaderboard service. This is an alias of the
// previously named DiscordTagLookupRequestedV1 to smooth migration.
const LeaderboardTagLookupRequestedV1 = DiscordTagLookupRequestedV1

// LeaderboardTagLookupSucceededV1 is published when a tag lookup succeeds.
// Pattern: Event Notification
// Subject: leaderboard.tag.lookup.by.user.id.success.v1
// Producer: leaderboard-service
// Consumers: requesters (round or discord via subscription to leaderboard)
// Version: v1 (January 2026)
const LeaderboardTagLookupSucceededV1 = "leaderboard.tag.lookup.by.user.id.success.v1"

// LeaderboardTagLookupNotFoundV1 is published when a lookup finds nothing.
// Pattern: Event Notification
// Subject: leaderboard.tag.lookup.by.user.id.not.found.v1
// Producer: leaderboard-service
// Consumers: requesters (round or discord via subscription to leaderboard)
// Version: v1 (January 2026)
const LeaderboardTagLookupNotFoundV1 = "leaderboard.tag.lookup.by.user.id.not.found.v1"

// LeaderboardTagLookupFailedV1 is published when a lookup fails with an error.
// Pattern: Event Notification
// Subject: leaderboard.tag.lookup.by.user.id.failed.v1
// Producer: leaderboard-service
// Consumers: requesters (round or discord via subscription to leaderboard)
// Version: v1 (January 2026)
const LeaderboardTagLookupFailedV1 = "leaderboard.tag.lookup.by.user.id.failed.v1"

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
