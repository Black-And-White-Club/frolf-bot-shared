// Package sharedevents contains cross-module shared events.
//
// MIGRATION NOTICE: This file contains legacy event constants.
// New code should use the versioned events from the flow-based files:
//   - tags.go: RoundTagLookupRequestedV1, DiscordTagLookupRequestedV1, etc.
//
// See each file for detailed flow documentation and versioning information.
package sharedevents

import (
	roundtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/round"
	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
)

// ScopedGuildID can be embedded to ensure all events include a GuildID for multi-tenancy.
type ScopedGuildID struct {
	GuildID sharedtypes.GuildID `json:"guild_id"`
}

// Cross-module event constants
// Deprecated: Use versioned constants from tags.go
const (
	// Deprecated: Use RoundTagLookupRequestedV1 from tags.go
	RoundTagLookupRequest = "round.tag.lookup.request"

	// Deprecated: Use RoundTagLookupFoundV1 from tags.go
	RoundTagLookupFound = "round.tag.lookup.found"

	// Deprecated: Use RoundTagLookupNotFoundV1 from tags.go
	RoundTagLookupNotFound = "round.tag.lookup.not.found"

	// Deprecated: Use TagUpdateForScheduledRoundsV1 from tags.go
	TagUpdateForScheduledRounds = "round.tag.update.for.scheduled.rounds"

	// Deprecated: Use DiscordTagLookupRequestedV1 from tags.go
	DiscordTagLookUpByUserIDRequest = "leaderboard.tag.lookup.by.user.id.request"
	// Deprecated: Use DiscordTagLookupFailedV1 from tags.go
	DiscordTagLookupByUserIDFailed = "discord.leaderboard.tag.lookup.by.user.id.failed"
	// Deprecated: Use DiscordTagLookupSucceededV1 from tags.go
	DiscordTagLookupByUserIDSuccess = "discord.leaderboard.tag.lookup.by.user.id.success"
	// Deprecated: Use DiscordTagLookupNotFoundV1 from tags.go
	DiscordTagLookupByUserIDNotFound = "discord.leaderboard.tag.lookup.by.user.id.not.found"

	// Deprecated: Use LeaderboardBatchTagAssignmentRequestedV1 from tags.go
	LeaderboardBatchTagAssignmentRequested = "leaderboard.batch.tag.assignment.requested"
)

// RoundTagLookupRequestPayload is the payload for requesting a tag number lookup from the Leaderboard module.
// Published by the Round module, consumed by the Leaderboard module.
type RoundTagLookupRequestPayload struct {
	ScopedGuildID
	UserID     sharedtypes.DiscordID `json:"user_id"`
	RoundID    sharedtypes.RoundID   `json:"round_id"`
	Response   roundtypes.Response   `json:"response"`
	JoinedLate *bool                 `json:"joined_late,omitempty"`
}

// RoundTagLookupResultPayload is the payload for the result of a tag number lookup request.
// Published by the Leaderboard module, consumed by the Round module.
type RoundTagLookupResultPayload struct {
	ScopedGuildID
	UserID             sharedtypes.DiscordID  `json:"user_id"`
	RoundID            sharedtypes.RoundID    `json:"round_id"`
	TagNumber          *sharedtypes.TagNumber `json:"tag_number,omitempty"`
	Found              bool                   `json:"found"`
	OriginalResponse   roundtypes.Response    `json:"original_response"`
	OriginalJoinedLate *bool                  `json:"original_joined_late,omitempty"`
	Error              string                 `json:"error,omitempty"`
}

type RoundTagLookupFailedPayload struct {
	ScopedGuildID
	UserID  sharedtypes.DiscordID `json:"user_id"`
	RoundID sharedtypes.RoundID   `json:"round_id"`
	Reason  string                `json:"reason"`
}

// TagAssignment represents a user and the tag number they should be assigned.
type TagAssignmentInfo struct {
	UserID    sharedtypes.DiscordID `json:"user_id"`
	TagNumber sharedtypes.TagNumber `json:"tag_number"`
}

// BatchTagAssignmentRequestedPayload is published by the score module after processing a round.
type BatchTagAssignmentRequestedPayload struct {
	ScopedGuildID
	RequestingUserID sharedtypes.DiscordID `json:"requester_user_id"`
	BatchID          string                `json:"batch_id"`
	Assignments      []TagAssignmentInfo   `json:"assignments"`
}

type DiscordTagLookupRequestPayload struct {
	ScopedGuildID
	RequestingUserID sharedtypes.DiscordID `json:"requester_user_id"`
	UserID           sharedtypes.DiscordID `json:"user_id"`
}

type DiscordTagLookupResultPayload struct {
	ScopedGuildID
	RequestingUserID sharedtypes.DiscordID  `json:"requester_user_id"`
	UserID           sharedtypes.DiscordID  `json:"user_id"`
	TagNumber        *sharedtypes.TagNumber `json:"tag_number,omitempty"`
	Found            bool                   `json:"found"`
}

type DiscordTagLookupByUserIDFailedPayload struct {
	ScopedGuildID
	RequestingUserID sharedtypes.DiscordID `json:"requester_user_id"`
	UserID           sharedtypes.DiscordID `json:"user_id"`
	Reason           string                `json:"reason"`
}
