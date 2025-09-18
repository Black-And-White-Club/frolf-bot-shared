package leaderboardevents

import (
	leaderboardtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/leaderboard"
	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
	sharedevents "github.com/Black-And-White-Club/frolf-bot-shared/events/shared"
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
	RoundFinalized              = "leaderboard.round.finalized"
	LeaderboardUpdateRequested  = "leaderboard.update.requested"
	LeaderboardUpdated          = "discord.leaderboard.batch.tag.assigned"
	LeaderboardUpdateFailed     = "discord.leaderboard.update.failed"
	DeactivateOldLeaderboard    = "leaderboard.deactivate"
	TagUpdateForScheduledRounds = "round.tag.update.for.scheduled.rounds"

	// Tag Assignment
	TagAvailabilityCheckRequest            = "leaderboard.tag.availability.check.requested"
	LeaderboardTagAssignmentRequested      = "leaderboard.tag.assignment.requested"
	LeaderboardTagAssignmentFailed         = "discord.leaderboard.tag.assignment.failed"
	LeaderboardTagAssignmentSuccess        = "discord.leaderboard.tag.assignment.success"
	TagAvailable                           = "user.tag.available"
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
	LeaderboardTraceEvent = "discord.leaderboard.trace.event"

	// Request events
	GetTagByUserIDRequest      = "leaderboard.tag.get.by.user.id.request"
	RoundGetTagByUserIDRequest = "leaderboard.round.tag.get.by.user.id.request"

	// Response events
	GetTagNumberResponse   = "round.leaderboard.tag.get.by.user.id.response"
	GetTagByUserIDResponse = "round.leaderboard.tag.get.by.user.id.response"
	GetTagByUserIDNotFound = "round.leaderboard.tag.get.by.user.id.not.found"

	// Round-specific response events
	RoundTagNumberFound    = "round.leaderboard.tag.found"
	RoundTagNumberNotFound = "round.leaderboard.tag.not.found"

	// Failure events
	GetTagNumberFailed = "discord.leaderboard.tag.get.by.user.id.failed"
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
	GuildID    sharedtypes.GuildID    `json:"guild_id"`
	UserID     sharedtypes.DiscordID  `json:"user_id"`
	TagNumber  *sharedtypes.TagNumber `json:"tag_number"`
	Source     string                 `json:"source"`
	UpdateType string                 `json:"update_type"`
	Reason     string                 `json:"reason"`
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
	GuildID     sharedtypes.GuildID   `json:"guild_id"`
	RequestorID sharedtypes.DiscordID `json:"requestor_id"`
	TargetID    sharedtypes.DiscordID `json:"target_id"`
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
