package leaderboardevents

import (
	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
)

// Stream names
const (
	LeaderboardStreamName = "leaderboard"
	UserStreamName        = "user"
	RoundStreamName       = "round"
	ScoreStreamName       = "score"
)

// Leaderboard-related events
const (
	// Leaderboard Update
	RoundFinalized             = "leaderboard.round.finalized"
	LeaderboardUpdateRequested = "leaderboard.update.requested"
	LeaderboardUpdated         = "discord.leaderboard.updated"
	LeaderboardUpdateFailed    = "discord.leaderboard.update.failed"
	DeactivateOldLeaderboard   = "leaderboard.deactivate"

	// Tag Assignment
	TagAvailabilityCheckRequest            = "leaderboard.tag.availability.check.request"
	LeaderboardTagAssignmentRequested      = "leaderboard.tag.assignment.requested"
	LeaderboardTagAssignmentFailed         = "discord.leaderboard.tag.assignment.failed"
	LeaderboardTagAssignmentSuccess        = "discord.leaderboard.tag.assignment.success"
	TagAvailable                           = "leaderboard.tag.available"
	TagUnavailable                         = "user.tag.unavailable"
	TagAvailableCheckFailure               = "leaderboard.tag.availability.failure"
	LeaderboardBatchTagAssignmentRequested = "leaderboard.batch.tag.assignment.requested"
	LeaderboardBatchTagAssignmentFailed    = "discord.leaderboard.batch.tag.assignment.failed"
	LeaderboardBatchTagAssigned            = "discord.leaderboard.batch.tag.assigned"

	// Tag Swap
	TagSwapRequested = "leaderboard.tag.swap.requested"
	TagSwapInitiated = "leaderboard.tag.swap.initiated"
	TagSwapFailed    = "discord.leaderboard.tag.swap.failed"
	TagSwapProcessed = "discord.leaderboard.tag.swap.processed"

	// Leaderboard Requests
	GetLeaderboardRequest  = "leaderboard.get.request"
	GetLeaderboardResponse = "discord.get.leaderboard.success"
	GetLeaderboardFailed   = "discord.get.leaderboard.failed"

	// Tag Requests
	GetTagByUserIDRequest  = "leaderboard.get.tag.number.request"
	GetTagByUserIDResponse = "round.get.tag.number.response"
	LeaderboardTraceEvent  = "discord.leaderboard.trace.event"
	GetTagNumberFailed     = "discord.get.tag.by.discordid.failed"
	GetTagNumberResponse   = "discord.get.tag.by.discordid.success"
)

// -- Event Payloads --

// RoundFinalizedPayload is the payload for the RoundFinalized event.
type RoundFinalizedPayload struct {
	RoundID               sharedtypes.RoundID `json:"round_id"`
	SortedParticipantTags TagOrder            `json:"sorted_participant_tags"` // Slice of "tag:UserID" strings
}

// TagAssignedPayload is the payload for the TagAssigned event.
type TagAssignedPayload struct {
	UserID       sharedtypes.DiscordID  `json:"user_id"`
	TagNumber    *sharedtypes.TagNumber `json:"tag_number"`
	AssignmentID sharedtypes.RoundID    `json:"assignment_id"`
	Source       string                 `json:"source"`
}

// TagAvailablePayload is the payload for the TagAvailable event.
type TagAvailablePayload struct {
	UserID       sharedtypes.DiscordID  `json:"user_id"`
	TagNumber    *sharedtypes.TagNumber `json:"tag_number"`
	AssignmentID string                 `json:"assignment_id"`
}

// TagUnavailablePayload is the payload for the TagUnavailable event.
type TagUnavailablePayload struct {
	UserID    sharedtypes.DiscordID  `json:"user_id"`
	TagNumber *sharedtypes.TagNumber `json:"tag_number"`
}

// TagAssignmentRequestedPayload is the payload for the TagAssignmentRequested event.
type TagAssignmentRequestedPayload struct {
	UserID     sharedtypes.DiscordID  `json:"user_id"`
	TagNumber  *sharedtypes.TagNumber `json:"tag_number"`
	UpdateID   sharedtypes.RoundID    `json:"update_id"`
	Source     string                 `json:"source"`
	UpdateType string                 `json:"update_type"`
}

// LeaderboardUpdateRequestedPayload is the payload for the LeaderboardUpdateRequested event.
type LeaderboardUpdateRequestedPayload struct {
	RoundID               sharedtypes.RoundID `json:"round_id"`
	SortedParticipantTags []string            `json:"sorted_participant_tags"`
	Source                string              `json:"source"`    // "round", "manual"
	UpdateID              string              `json:"update_id"` // round ID or manual update identifier
}

// LeaderboardUpdatedPayload is the payload for the LeaderboardUpdated event.
type LeaderboardUpdatedPayload struct {
	LeaderboardID   int64               `json:"leaderboard_id"`
	RoundID         sharedtypes.RoundID `json:"round_id"`
	LeaderboardData map[int]string      `json:"leaderboard_data"`
}

// LeaderboardUpdateFailedPayload is the payload for the LeaderboardUpdateFailed event.
type LeaderboardUpdateFailedPayload struct {
	RoundID sharedtypes.RoundID `json:"round_id"`
	Reason  string              `json:"reason"` // Reason for the failure
}

// DeactivateOldLeaderboardPayload is the payload for the DeactivateOldLeaderboard event.
type DeactivateOldLeaderboardPayload struct {
	LeaderboardID int64 `json:"leaderboard_id"`
}

// TagAssignmentFailedPayload is the payload for the TagAssignmentFailed event.
type TagAssignmentFailedPayload struct {
	UserID     sharedtypes.DiscordID  `json:"user_id"`
	TagNumber  *sharedtypes.TagNumber `json:"tag_number"`
	Source     string                 `json:"source"`
	UpdateType string                 `json:"update_type"`
	Reason     string                 `json:"reason"`
}

// TagAvailabilityCheckResultPayload is the payload for the result of a tag availability check.
type TagAvailabilityCheckResultPayload struct {
	UserID    sharedtypes.DiscordID  `json:"user_id"`
	TagNumber *sharedtypes.TagNumber `json:"tag_number"`
	Available bool                   `json:"tag_available"`
}

// TagAvailabilityCheckFailedPayload is the payload for the failure of a tag availability check.
type TagAvailabilityCheckFailedPayload struct {
	UserID    sharedtypes.DiscordID  `json:"user_id"`
	TagNumber *sharedtypes.TagNumber `json:"tag_number"`
	Reason    string                 `json:"reason"`
}

// TagSwapRequestedPayload is the payload for the TagSwapRequested event.
type TagSwapRequestedPayload struct {
	RequestorID sharedtypes.DiscordID `json:"requestor_id"`
	TargetID    sharedtypes.DiscordID `json:"target_id"`
}

// TagSwapInitiatedPayload is the payload for the TagSwapInitiated event.
type TagSwapInitiatedPayload struct {
	RequestorID sharedtypes.DiscordID `json:"requestor_id"`
	TargetID    sharedtypes.DiscordID `json:"target_id"`
}

// TagSwapFailedPayload is the payload for the TagSwapFailed event.
type TagSwapFailedPayload struct {
	RequestorID sharedtypes.DiscordID `json:"requestor_id"`
	TargetID    sharedtypes.DiscordID `json:"target_id"`
	Reason      string                `json:"reason"`
}

// TagSwapProcessedPayload is the payload for the TagSwapProcessed event.
type TagSwapProcessedPayload struct {
	RequestorID sharedtypes.DiscordID `json:"requestor_id"`
	TargetID    sharedtypes.DiscordID `json:"target_id"`
}

// GetLeaderboardRequestPayload is the payload for the GetLeaderboardRequest event.
type GetLeaderboardRequestPayload struct{} // Empty, as no data is needed for this request

// LeaderboardEntry represents an entry on the leaderboard.
type LeaderboardEntry struct {
	TagNumber *sharedtypes.TagNumber `json:"tag_number"`
	UserID    sharedtypes.DiscordID  `json:"user_id"`
}

// GetLeaderboardResponsePayload is the payload for the GetLeaderboardResponse event.
type GetLeaderboardResponsePayload struct {
	Leaderboard []LeaderboardEntry `json:"leaderboard"`
}

// GetTagByUserIDRequestPayload is the payload for the GetTagByUserIDRequest event.
type TagNumberRequestPayload struct {
	UserID  sharedtypes.DiscordID `json:"user_id"`
	RoundID sharedtypes.RoundID   `json:"round_id"`
}

// GetTagByUserIDResponsePayload is the payload for the GetTagByUserIDResponse event.
type GetTagNumberResponsePayload struct {
	TagNumber *sharedtypes.TagNumber `json:"tag_number"`
	UserID    sharedtypes.DiscordID  `json:"user_id"`
	RoundID   sharedtypes.RoundID    `json:"round_id"`
}

// TagAvailabilityCheckRequestedPayload is the payload for the TagAvailabilityCheckRequested event.
type TagAvailabilityCheckRequestedPayload struct {
	TagNumber *sharedtypes.TagNumber `json:"tag_number"`
	UserID    sharedtypes.DiscordID  `json:"user_id"`
}

// -- Helper Types --

// TagOrder represents the order of tags.
type TagOrder []string

type GetLeaderboardFailedPayload struct {
	Reason string `json:"reason"` // Reason for the failure
}

// GetTagNumberFailedPayload is the payload for the GetTagNumberFailed event.
type GetTagNumberFailedPayload struct {
	Reason string `json:"reason"` // Reason for the failure
}

type BatchTagAssignmentRequestedPayload struct {
	RequestingUserID sharedtypes.DiscordID
	BatchID          string
	Assignments      []TagAssignmentInfo
}

type TagAssignmentInfo struct {
	UserID    sharedtypes.DiscordID
	TagNumber sharedtypes.TagNumber
}

type BatchTagAssignedPayload struct {
	RequestingUserID sharedtypes.DiscordID
	BatchID          string
	AssignmentCount  int
	Assignments      []TagAssignmentInfo
}

type BatchTagAssignmentFailedPayload struct {
	RequestingUserID sharedtypes.DiscordID
	BatchID          string
	Reason           string
}
