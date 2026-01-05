// Package leaderboardevents contains leaderboard-related domain events.
//
// MIGRATION NOTICE: This file contains legacy event constants.
// New code should use the versioned events from the flow-based files:
//   - updates.go: LeaderboardUpdateRequestedV1, LeaderboardUpdatedV1, etc.
//   - tags.go: TagAvailabilityCheckRequestedV1, LeaderboardTagAssignedV1, TagSwapRequestedV1, etc.
//
// See each file for detailed flow documentation and versioning information.
package leaderboardevents

import (
	sharedevents "github.com/Black-And-White-Club/frolf-bot-shared/events/shared"
	leaderboardtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/leaderboard"
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
// Deprecated: Use versioned constants from updates.go and tags.go
const (
	// Leaderboard Update
	// Deprecated: Use LeaderboardRoundFinalizedV1 from updates.go
	RoundFinalized = "leaderboard.round.finalized"
	// Deprecated: Use LeaderboardUpdateRequestedV1 from updates.go
	LeaderboardUpdateRequested = "leaderboard.update.requested"
	// Deprecated: Use LeaderboardUpdatedV1 from updates.go
	LeaderboardUpdated = "leaderboard.batch.tag.assigned"
	// Deprecated: Use LeaderboardUpdateFailedV1 from updates.go
	LeaderboardUpdateFailed = "leaderboard.update.failed"
	// Deprecated: Use DeactivateOldLeaderboardV1 from updates.go
	DeactivateOldLeaderboard = "leaderboard.deactivate"
	// Deprecated: Use TagUpdateForScheduledRoundsV1 from updates.go
	TagUpdateForScheduledRounds = "round.tag.update.for.scheduled.rounds"

	// Tag Assignment
	// Deprecated: Use TagAvailabilityCheckRequestedV1 from tags.go
	TagAvailabilityCheckRequest = "leaderboard.tag.availability.check.requested"
	// Deprecated: Use LeaderboardTagAssignmentRequestedV1 from tags.go
	LeaderboardTagAssignmentRequested = "leaderboard.tag.assignment.requested"
	// Deprecated: Use LeaderboardTagAssignmentFailedV1 from tags.go
	LeaderboardTagAssignmentFailed = "leaderboard.tag.assignment.failed"
	// Deprecated: Use LeaderboardTagAssignedV1 from tags.go
	LeaderboardTagAssignmentSuccess = "leaderboard.tag.assignment.success"
	// Deprecated: Use LeaderboardTagAvailableV1 from tags.go
	TagAvailable = "user.tag.available"
	// Deprecated: Use LeaderboardTagUnavailableV1 from tags.go
	TagUnavailable = "user.tag.unavailable"
	// Deprecated: Use TagAvailabilityCheckFailedV1 from tags.go
	TagAvailableCheckFailure = "leaderboard.tag.availability.failure"
	// Deprecated: Use LeaderboardBatchTagAssignmentRequestedV1 from tags.go
	LeaderboardBatchTagAssignmentRequested = "leaderboard.batch.tag.assignment.requested"
	// Deprecated: Use LeaderboardBatchTagAssignmentFailedV1 from tags.go
	LeaderboardBatchTagAssignmentFailed = "leaderboard.batch.tag.assignment.failed"
	// Deprecated: Use LeaderboardBatchTagAssignedV1 from tags.go
	LeaderboardBatchTagAssigned = "leaderboard.batch.tag.assigned"

	// Tag Swap
	// Deprecated: Use TagSwapRequestedV1 from tags.go
	TagSwapRequested = "leaderboard.tag.swap.requested"
	// Deprecated: Use TagSwapInitiatedV1 from tags.go
	TagSwapInitiated = "leaderboard.tag.swap.initiated"
	// Deprecated: Use TagSwapFailedV1 from tags.go
	TagSwapFailed = "leaderboard.tag.swap.failed"
	// Deprecated: Use TagSwapProcessedV1 from tags.go
	TagSwapProcessed = "leaderboard.tag.swap.processed"

	// Leaderboard Requests
	// Deprecated: Use GetLeaderboardRequestedV1 from updates.go
	GetLeaderboardRequest = "leaderboard.get.request"
	// Deprecated: Use GetLeaderboardResponseV1 from updates.go
	GetLeaderboardResponse = "leaderboard.get.success"
	// Deprecated: Use GetLeaderboardFailedV1 from updates.go
	GetLeaderboardFailed = "leaderboard.get.failed"

	// Tag Requests
	// Deprecated: Use LeaderboardTraceEventV1 from tags.go
	LeaderboardTraceEvent = "leaderboard.trace.event"

	// Request events
	// Deprecated: Use GetTagByUserIDRequestedV1 from tags.go
	GetTagByUserIDRequest = "leaderboard.tag.get.by.user.id.request"
	// Deprecated: Use RoundGetTagByUserIDRequestedV1 from tags.go
	RoundGetTagByUserIDRequest = "leaderboard.round.tag.get.by.user.id.request"

	// Response events
	// Deprecated: Use GetTagNumberResponseV1 from tags.go
	GetTagNumberResponse = "round.leaderboard.tag.get.by.user.id.response"
	// Deprecated: Use GetTagNumberResponseV1 from tags.go
	GetTagByUserIDResponse = "round.leaderboard.tag.get.by.user.id.response"
	// Deprecated: Use RoundTagNumberNotFoundV1 from tags.go
	GetTagByUserIDNotFound = "round.leaderboard.tag.get.by.user.id.not.found"

	// Round-specific response events
	// Deprecated: Use RoundTagNumberFoundV1 from tags.go
	RoundTagNumberFound = "round.leaderboard.tag.found"
	// Deprecated: Use RoundTagNumberNotFoundV1 from tags.go
	RoundTagNumberNotFound = "round.leaderboard.tag.not.found"

	// Failure events
	// Deprecated: Use GetTagNumberFailedV1 from tags.go
	GetTagNumberFailed = "leaderboard.tag.get.by.user.id.failed"
)

// -- Event Payloads --

// RoundFinalizedPayload is the payload for the RoundFinalized event.
type RoundFinalizedPayload struct {
	GuildID               sharedtypes.GuildID `json:"guild_id"`
	RoundID               sharedtypes.RoundID `json:"round_id"`
	SortedParticipantTags TagOrder            `json:"sorted_participant_tags"` // Slice of "tag:UserID" strings
}

// TagAssignedPayload is the payload for the TagAssigned event.
type TagAssignedPayload struct {
	GuildID      sharedtypes.GuildID    `json:"guild_id"`
	UserID       sharedtypes.DiscordID  `json:"user_id"`
	TagNumber    *sharedtypes.TagNumber `json:"tag_number"`
	AssignmentID sharedtypes.RoundID    `json:"assignment_id"`
	Source       string                 `json:"source"`
}

// TagAvailablePayload is the payload for the TagAvailable event.
type TagAvailablePayload struct {
	GuildID      sharedtypes.GuildID    `json:"guild_id"`
	UserID       sharedtypes.DiscordID  `json:"user_id"`
	TagNumber    *sharedtypes.TagNumber `json:"tag_number"`
	AssignmentID string                 `json:"assignment_id"`
}

// TagUnavailablePayload is the payload for the TagUnavailable event.
type TagUnavailablePayload struct {
	GuildID   sharedtypes.GuildID    `json:"guild_id"`
	UserID    sharedtypes.DiscordID  `json:"user_id"`
	TagNumber *sharedtypes.TagNumber `json:"tag_number"`
	Reason    string                 `json:"reason"` // Reason why the tag is unavailable
}

// TagAssignmentRequestedPayload is the payload for the TagAssignmentRequested event.
type TagAssignmentRequestedPayload struct {
	GuildID    sharedtypes.GuildID    `json:"guild_id"`
	UserID     sharedtypes.DiscordID  `json:"user_id"`
	TagNumber  *sharedtypes.TagNumber `json:"tag_number"`
	UpdateID   sharedtypes.RoundID    `json:"update_id"`
	Source     string                 `json:"source"`
	UpdateType string                 `json:"update_type"`
}

type LeaderboardTagAssignRequestPayload struct {
	GuildID          sharedtypes.GuildID   `json:"guild_id"`
	RequestorID      sharedtypes.DiscordID `json:"requestor_id"`
	TagNumber        sharedtypes.TagNumber `json:"tag_number"`
	ChannelID        string                `json:"channel_id"`
	RequestID        string                `json:"request_id"`
	InteractionToken string                `json:"interaction_token"`
}

// LeaderboardUpdateRequestedPayload is the payload for the LeaderboardUpdateRequested event.
type LeaderboardUpdateRequestedPayload struct {
	GuildID               sharedtypes.GuildID `json:"guild_id"`
	RoundID               sharedtypes.RoundID `json:"round_id"`
	SortedParticipantTags []string            `json:"sorted_participant_tags"`
	Source                string              `json:"source"`    // "round", "manual"
	UpdateID              string              `json:"update_id"` // round ID or manual update identifier
}

// LeaderboardUpdatedPayload is the payload for the LeaderboardUpdated event.
type LeaderboardUpdatedPayload struct {
	GuildID         sharedtypes.GuildID                             `json:"guild_id"`
	LeaderboardID   int64                                           `json:"leaderboard_id"`
	RoundID         sharedtypes.RoundID                             `json:"round_id"`
	LeaderboardData map[sharedtypes.TagNumber]sharedtypes.DiscordID `json:"leaderboard_data"`
	Config          *sharedevents.GuildConfigFragment               `json:"config_fragment,omitempty"`
}

// LeaderboardUpdateFailedPayload is the payload for the LeaderboardUpdateFailed event.
type LeaderboardUpdateFailedPayload struct {
	GuildID sharedtypes.GuildID `json:"guild_id"`
	RoundID sharedtypes.RoundID `json:"round_id"`
	Reason  string              `json:"reason"` // Reason for the failure
}

// DeactivateOldLeaderboardPayload is the payload for the DeactivateOldLeaderboard event.
type DeactivateOldLeaderboardPayload struct {
	GuildID       sharedtypes.GuildID `json:"guild_id"`
	LeaderboardID int64               `json:"leaderboard_id"`
}

// TagAssignmentFailedPayload is the payload for the TagAssignmentFailed event.
type TagAssignmentFailedPayload struct {
	GuildID    sharedtypes.GuildID               `json:"guild_id"`
	UserID     sharedtypes.DiscordID             `json:"user_id"`
	TagNumber  *sharedtypes.TagNumber            `json:"tag_number"`
	Source     string                            `json:"source"`
	UpdateType string                            `json:"update_type"`
	Reason     string                            `json:"reason"`
	Config     *sharedevents.GuildConfigFragment `json:"config_fragment,omitempty"`
}

// TagAvailabilityCheckResultPayload is the payload for the result of a tag availability check.
type TagAvailabilityCheckResultPayload struct {
	GuildID   sharedtypes.GuildID    `json:"guild_id"`
	UserID    sharedtypes.DiscordID  `json:"user_id"`
	TagNumber *sharedtypes.TagNumber `json:"tag_number"`
	Available bool                   `json:"tag_available"`
	Reason    string                 `json:"reason,omitempty"` // Reason for unavailability (empty if available)
}

// TagAvailabilityCheckFailedPayload is the payload for the failure of a tag availability check.
type TagAvailabilityCheckFailedPayload struct {
	GuildID   sharedtypes.GuildID    `json:"guild_id"`
	UserID    sharedtypes.DiscordID  `json:"user_id"`
	TagNumber *sharedtypes.TagNumber `json:"tag_number"`
	Reason    string                 `json:"reason"`
}

// TagSwapRequestedPayload is the payload for the TagSwapRequested event.
type TagSwapRequestedPayload struct {
	GuildID     sharedtypes.GuildID   `json:"guild_id"`
	RequestorID sharedtypes.DiscordID `json:"requestor_id"`
	TargetID    sharedtypes.DiscordID `json:"target_id"`
}

// TagSwapInitiatedPayload is the payload for the TagSwapInitiated event.
type TagSwapInitiatedPayload struct {
	GuildID     sharedtypes.GuildID   `json:"guild_id"`
	RequestorID sharedtypes.DiscordID `json:"requestor_id"`
	TargetID    sharedtypes.DiscordID `json:"target_id"`
}

// TagSwapFailedPayload is the payload for the TagSwapFailed event.
type TagSwapFailedPayload struct {
	GuildID     sharedtypes.GuildID   `json:"guild_id"`
	RequestorID sharedtypes.DiscordID `json:"requestor_id"`
	TargetID    sharedtypes.DiscordID `json:"target_id"`
	Reason      string                `json:"reason"`
}

// TagSwapProcessedPayload is the payload for the TagSwapProcessed event.
type TagSwapProcessedPayload struct {
	GuildID     sharedtypes.GuildID               `json:"guild_id"`
	RequestorID sharedtypes.DiscordID             `json:"requestor_id"`
	TargetID    sharedtypes.DiscordID             `json:"target_id"`
	Config      *sharedevents.GuildConfigFragment `json:"config_fragment,omitempty"`
}

// GetLeaderboardRequestPayload is the payload for the GetLeaderboardRequest event.
type GetLeaderboardRequestPayload struct {
	GuildID sharedtypes.GuildID `json:"guild_id"`
}

// GetLeaderboardResponsePayload is the payload for the GetLeaderboardResponse event.
type GetLeaderboardResponsePayload struct {
	GuildID     sharedtypes.GuildID              `json:"guild_id"`
	Leaderboard leaderboardtypes.LeaderboardData `json:"leaderboard"`
}
type SoloTagNumberRequestPayload struct {
	GuildID sharedtypes.GuildID   `json:"guild_id"`
	UserID  sharedtypes.DiscordID `json:"user_id"`
}

type SoloTagNumberResponsePayload struct {
	GuildID   sharedtypes.GuildID    `json:"guild_id"`
	TagNumber *sharedtypes.TagNumber `json:"tag_number"`
	UserID    sharedtypes.DiscordID  `json:"user_id"`
}

// GetTagByUserIDRequestPayload is the payload for the GetTagByUserIDRequest event.
type TagNumberRequestPayload struct {
	GuildID sharedtypes.GuildID   `json:"guild_id"`
	UserID  sharedtypes.DiscordID `json:"user_id"`
	RoundID sharedtypes.RoundID   `json:"round_id"`
}

// GetTagByUserIDResponsePayload is the payload for the GetTagByUserIDResponse event.
type GetTagNumberResponsePayload struct {
	GuildID   sharedtypes.GuildID    `json:"guild_id"`
	TagNumber *sharedtypes.TagNumber `json:"tag_number"`
	UserID    sharedtypes.DiscordID  `json:"user_id"`
	RoundID   sharedtypes.RoundID    `json:"round_id"`
	Found     bool                   `json:"found"`
}

// TagAvailabilityCheckRequestedPayload is the payload for the TagAvailabilityCheckRequested event.
type TagAvailabilityCheckRequestedPayload struct {
	GuildID   sharedtypes.GuildID    `json:"guild_id"`
	TagNumber *sharedtypes.TagNumber `json:"tag_number"`
	UserID    sharedtypes.DiscordID  `json:"user_id"`
}

// -- Helper Types --

// TagOrder represents the order of tags.
type TagOrder []string

type GetLeaderboardFailedPayload struct {
	GuildID sharedtypes.GuildID `json:"guild_id"`
	Reason  string              `json:"reason"` // Reason for the failure
}

// GetTagNumberFailedPayload is the payload for the GetTagNumberFailed event.
type GetTagNumberFailedPayload struct {
	GuildID sharedtypes.GuildID `json:"guild_id"`
	Reason  string              `json:"reason"` // Reason for the failure
}

type BatchTagAssignmentRequestedPayload struct {
	GuildID          sharedtypes.GuildID `json:"guild_id"`
	RequestingUserID sharedtypes.DiscordID
	BatchID          string
	Assignments      []TagAssignmentInfo
}

type TagAssignmentInfo struct {
	GuildID   sharedtypes.GuildID `json:"guild_id"`
	UserID    sharedtypes.DiscordID
	TagNumber sharedtypes.TagNumber
}

type BatchTagAssignedPayload struct {
	GuildID          sharedtypes.GuildID `json:"guild_id"`
	RequestingUserID sharedtypes.DiscordID
	BatchID          string
	AssignmentCount  int
	Assignments      []TagAssignmentInfo
	Config           *sharedevents.GuildConfigFragment `json:"config_fragment,omitempty"`
}

type BatchTagAssignmentFailedPayload struct {
	GuildID          sharedtypes.GuildID `json:"guild_id"`
	RequestingUserID sharedtypes.DiscordID
	BatchID          string
	Reason           string
}

// GetUserID returns the UserID of the payload
func (p *TagAssignmentRequestedPayload) GetUserID() sharedtypes.DiscordID {
	return p.UserID
}

// GetTagNumber returns the TagNumber of the payload
func (p *TagAssignmentRequestedPayload) GetTagNumber() *sharedtypes.TagNumber {
	return p.TagNumber
}

func (p *TagAssignedPayload) GetUserID() sharedtypes.DiscordID {
	return p.UserID
}

// GetTagNumber returns the TagNumber of the payload
func (p *TagAssignedPayload) GetTagNumber() *sharedtypes.TagNumber {
	return p.TagNumber
}
