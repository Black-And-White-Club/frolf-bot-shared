package sharedevents

import (
	roundtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/round"
	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
)

const (
	// RoundTagLookupRequest is the topic for requesting a tag number lookup from the Leaderboard module.
	RoundTagLookupRequest = "round.tag.lookup.request"

	// RoundTagLookupFound is the topic for the result when a tag number is found.
	RoundTagLookupFound = "round.tag.lookup.found"

	// RoundTagLookupNotFound is the topic for the result when a tag number is not found.
	RoundTagLookupNotFound = "round.tag.lookup.not.found"

	DiscordTagLookUpByUserIDRequest  = "leaderboard.tag.lookup.by.user.id.request"
	DiscordTagLookupByUserIDFailed   = "discord.leaderboard.tag.lookup.by.user.id.failed"
	DiscordTagLookupByUserIDSuccess  = "discord.leaderboard.tag.lookup.by.user.id.success"
	DiscordTagLoopupByUserIDNotFound = "discord.leaderboard.tag.lookup.by.user.id.not.found"
)

// RoundTagLookupRequestPayload is the payload for requesting a tag number lookup from the Leaderboard module.
// Published by the Round module, consumed by the Leaderboard module.
type RoundTagLookupRequestPayload struct {
	UserID     sharedtypes.DiscordID `json:"user_id"`
	RoundID    sharedtypes.RoundID   `json:"round_id"`
	Response   roundtypes.Response   `json:"response"`
	JoinedLate *bool                 `json:"joined_late,omitempty"`
}

// RoundTagLookupResultPayload is the payload for the result of a tag number lookup request.
// Published by the Leaderboard module, consumed by the Round module.
type RoundTagLookupResultPayload struct {
	UserID             sharedtypes.DiscordID  `json:"user_id"`
	RoundID            sharedtypes.RoundID    `json:"round_id"`
	TagNumber          *sharedtypes.TagNumber `json:"tag_number,omitempty"`
	Found              bool                   `json:"found"`
	OriginalResponse   roundtypes.Response    `json:"original_response"`
	OriginalJoinedLate *bool                  `json:"original_joined_late,omitempty"`
	Error              string                 `json:"error,omitempty"`
}

type RoundTagLookupFailedPayload struct {
	UserID  sharedtypes.DiscordID `json:"user_id"`
	RoundID sharedtypes.RoundID   `json:"round_id"`
	Reason  string                `json:"reason"`
}

// LeaderboardBatchTagAssignmentRequested is the topic for triggering batch tag updates to the leaderboard.
const LeaderboardBatchTagAssignmentRequested = "leaderboard.batch.tag.assignment.requested"

// TagAssignment represents a user and the tag number they should be assigned.
type TagAssignmentInfo struct {
	UserID    sharedtypes.DiscordID
	TagNumber sharedtypes.TagNumber
}

// BatchTagAssignmentRequestedPayload is published by the score module after processing a round.
type BatchTagAssignmentRequestedPayload struct {
	RequestingUserID sharedtypes.DiscordID
	BatchID          string
	Assignments      []TagAssignmentInfo
}

type DiscordTagLookupRequestPayload struct {
	RequestingUserID sharedtypes.DiscordID `json:"requester_user_id"`
	UserID           sharedtypes.DiscordID `json:"user_id"`
}

type DiscordTagLookupResultPayload struct {
	RequestingUserID sharedtypes.DiscordID  `json:"requester_user_id"`
	UserID           sharedtypes.DiscordID  `json:"user_id"`
	TagNumber        *sharedtypes.TagNumber `json:"tag_number,omitempty"`
	Found            bool                   `json:"found"`
}

type DiscordTagLookupByUserIDFailedPayload struct {
	RequestingUserID sharedtypes.DiscordID `json:"requester_user_id"`
	UserID           sharedtypes.DiscordID `json:"user_id"`
	Reason           string                `json:"reason"`
}
