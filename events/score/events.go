// Package scoreevents contains score-related domain events.
//
// MIGRATION NOTICE: This file contains legacy event constants.
// New code should use the versioned events from the flow-based files:
//   - processing.go: ProcessRoundScoresRequestedV1, ProcessRoundScoresSucceededV1, etc.
//   - updates.go: ScoreUpdateRequestedV1, ScoreUpdatedV1, ScoreBulkUpdateRequestedV1, etc.
//
// See each file for detailed flow documentation and versioning information.
package scoreevents

import (
	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
)

// Stream names
const (
	ScoreStreamName       = "score"
	RoundStreamName       = "round"
	UserStreamName        = "user"
	LeaderboardStreamName = "leaderboard"
)

// Event subjects
// Deprecated: Use versioned constants from processing.go and updates.go
const (
	// Deprecated: Use ProcessRoundScoresRequestedV1 from processing.go
	ProcessRoundScoresRequest = "score.process.round.scores.request"
	// Deprecated: Use LeaderboardUpdateRequestedV1 from leaderboard module
	LeaderboardUpdateRequested = "leaderboard.update.requested"
	// Deprecated: Use ScoreUpdateRequestedV1 from updates.go
	ScoreUpdateRequest = "score.update.request"
	// Deprecated: Use ScoreBulkUpdateRequestedV1 from updates.go
	ScoreBulkUpdateRequest = "score.update.bulk.request"
	// Deprecated: Use ScoreUpdatedV1 from updates.go
	ScoreUpdateSuccess = "score.update.success"
	// Deprecated: Use ScoreUpdateFailedV1 from updates.go
	ScoreUpdateFailure = "score.update.fail"
	// Deprecated: Use ScoreBulkUpdatedV1 from updates.go
	ScoreBulkUpdateSuccess = "score.update.bulk.success"
	// Deprecated: Use ProcessRoundScoresSucceededV1 from processing.go
	ProcessRoundScoresSuccess = "leaderboard.batch.tag.assignment.requested"
	// Deprecated: Use ProcessRoundScoresFailedV1 from processing.go
	ProcessRoundScoresFailure = "score.process.round.scores.fail"
)

// Event payloads
type ProcessRoundScoresRequestPayload struct {
	GuildID   sharedtypes.GuildID     `json:"guild_id"`
	RoundID   sharedtypes.RoundID     `json:"round_id"`
	Scores    []sharedtypes.ScoreInfo `json:"scores"`
	Overwrite bool                    `json:"overwrite"`
}

// ProcessRoundScoresSuccessPayload is the payload for the ProcessRoundScoresSuccess event.
type ProcessRoundScoresSuccessPayload struct {
	GuildID     sharedtypes.GuildID      `json:"guild_id"`
	RoundID     sharedtypes.RoundID      `json:"round_id"`
	TagMappings []sharedtypes.TagMapping `json:"tag_mappings"`
}

// ProcessRoundScoresFailurePayload is the payload for the ProcessRoundScoresFailure event.
type ProcessRoundScoresFailurePayload struct {
	GuildID sharedtypes.GuildID `json:"guild_id"`
	RoundID sharedtypes.RoundID `json:"round_id"`
	Error   string              `json:"error"`
}

// ScoreUpdateRequestPayload is the payload for the ScoreUpdateRequest event.
type ScoreUpdateRequestPayload struct {
	GuildID   sharedtypes.GuildID    `json:"guild_id"`
	RoundID   sharedtypes.RoundID    `json:"round_id"`
	UserID    sharedtypes.DiscordID  `json:"user_id"`
	Score     sharedtypes.Score      `json:"score"`
	TagNumber *sharedtypes.TagNumber `json:"tag_number,omitempty"`
}

// ScoreBulkUpdateRequestPayload batches multiple score corrections for a single round.
type ScoreBulkUpdateRequestPayload struct {
	GuildID sharedtypes.GuildID         `json:"guild_id"`
	RoundID sharedtypes.RoundID         `json:"round_id"`
	Updates []ScoreUpdateRequestPayload `json:"updates"`
}

// ScoreUpdateSuccessPayload is the payload for successful score updates.
type ScoreUpdateSuccessPayload struct {
	GuildID sharedtypes.GuildID   `json:"guild_id"`
	RoundID sharedtypes.RoundID   `json:"round_id"`
	UserID  sharedtypes.DiscordID `json:"user_id"`
	Score   sharedtypes.Score     `json:"score"`
}

// ScoreUpdateFailurePayload is the payload for failed score updates.
type ScoreUpdateFailurePayload struct {
	GuildID sharedtypes.GuildID   `json:"guild_id"`
	RoundID sharedtypes.RoundID   `json:"round_id"`
	UserID  sharedtypes.DiscordID `json:"user_id"`
	Error   string                `json:"error"`
}

// ScoreBulkUpdateSuccessPayload summarises a completed bulk update.
type ScoreBulkUpdateSuccessPayload struct {
	GuildID        sharedtypes.GuildID     `json:"guild_id"`
	RoundID        sharedtypes.RoundID     `json:"round_id"`
	AppliedCount   int                     `json:"applied_count"`
	FailedCount    int                     `json:"failed_count"`
	TotalRequested int                     `json:"total_requested"`
	UserIDsApplied []sharedtypes.DiscordID `json:"user_ids_applied"`
}
