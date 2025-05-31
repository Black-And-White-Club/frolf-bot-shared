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
const (
	ProcessRoundScoresRequest  = "score.process.round.scores.request"
	LeaderboardUpdateRequested = "leaderboard.update.requested"
	ScoreUpdateRequest         = "score.update.request"
	ScoreUpdateSuccess         = "discord.score.update.success"
	ScoreUpdateFailure         = "discord.score.update.fail"
	ProcessRoundScoresSuccess  = "leaderboard.batch.tag.assignment.requested"
	ProcessRoundScoresFailure  = "score.process.round.scores.fail"
)

// ProcessRoundScoresRequestPayload is the payload for the ProcessRoundScoresRequest event.
type ProcessRoundScoresRequestPayload struct {
	RoundID sharedtypes.RoundID     `json:"round_id"`
	Scores  []sharedtypes.ScoreInfo `json:"scores"`
}

// ProcessRoundScoresSuccessPayload is the payload for the ProcessRoundScoresSuccess event.
type ProcessRoundScoresSuccessPayload struct {
	RoundID     sharedtypes.RoundID      `json:"round_id"`
	TagMappings []sharedtypes.TagMapping `json:"tag_mappings"`
}

// ProcessRoundScoresFailurePayload is the payload for the ProcessRoundScoresFailure event.
type ProcessRoundScoresFailurePayload struct {
	RoundID sharedtypes.RoundID `json:"round_id"`
	Error   string              `json:"error"`
}

// ScoreUpdateRequestPayload is the payload for the ScoreUpdateRequest event.
type ScoreUpdateRequestPayload struct {
	RoundID   sharedtypes.RoundID    `json:"round_id"`
	UserID    sharedtypes.DiscordID  `json:"user_id"`
	Score     sharedtypes.Score      `json:"score"`
	TagNumber *sharedtypes.TagNumber `json:"tag_number,omitempty"`
}

// ScoreUpdateSuccessPayload is the payload for successful score updates.
type ScoreUpdateSuccessPayload struct {
	RoundID sharedtypes.RoundID   `json:"round_id"`
	UserID  sharedtypes.DiscordID `json:"user_id"`
	Score   sharedtypes.Score     `json:"score"`
}

// ScoreUpdateFailurePayload is the payload for failed score updates.
type ScoreUpdateFailurePayload struct {
	RoundID sharedtypes.RoundID   `json:"round_id"`
	UserID  sharedtypes.DiscordID `json:"user_id"`
	Error   string                `json:"error"`
}
