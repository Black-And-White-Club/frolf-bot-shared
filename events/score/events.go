package scoreevents

import (
	"github.com/Black-And-White-Club/frolf-bot-shared/events"
	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
)

// Stream names
const (
	ScoreStreamName       = "score"
	RoundStreamName       = "round"
	UserStreamName        = "user"        // You probably have this defined elsewhere
	LeaderboardStreamName = "leaderboard" // You probably have this defined elsewhere
)

// Event subjects.
const (
	ProcessRoundScoresRequest  = "score.process_round_scores.request" // From round module (backend internal)
	LeaderboardUpdateRequested = "score.leaderboard.update.requested" // To leaderboard module (backend internal)
	ScoreUpdateRequest         = "score.update.request"               // From Discord
	ScoreUpdateResponse        = "score.update.response"              // To Discord
	ScoreUpdateError           = "score.update.error"                 //Consider using later
)

// --- Payloads ---

// ProcessRoundScoresRequestPayload is the payload for the ProcessRoundScoresRequest event.
type ProcessRoundScoresRequestPayload struct {
	events.CommonMetadata
	RoundID sharedtypes.RoundID `json:"round_id"`
	Scores  []ParticipantScore  `json:"scores"`
}

// ParticipantScore represents a single score entry with UserID, TagNumber, and Score.
type ParticipantScore struct {
	events.CommonMetadata
	UserID    sharedtypes.DiscordID `json:"user_id"` // Consistent naming
	TagNumber sharedtypes.TagNumber `json:"tag_number"`
	Score     sharedtypes.Score     `json:"score"`
}

// LeaderboardUpdateRequestedPayload is the payload for the LeaderboardUpdateRequested event.
type LeaderboardUpdateRequestedPayload struct {
	events.CommonMetadata
	RoundID sharedtypes.RoundID `json:"round_id"`
	Scores  []ParticipantScore  `json:"scores"`
}

type Participant struct {
	UserID    sharedtypes.DiscordID  `json:"user_id"`
	TagNumber *sharedtypes.TagNumber `json:"tag_number,omitempty"`
	Score     *sharedtypes.Score     `json:"score"`
}

// ScoreUpdateRequestPayload is the payload for the ScoreUpdateRequest event.
type ScoreUpdateRequestPayload struct {
	events.CommonMetadata
	RoundID     sharedtypes.RoundID `json:"round_id"`
	Participant Participant         `json:"participant"` // Discord ID of the participant

}

// ScoreUpdateResponsePayload is the payload for the ScoreUpdateResponse
type ScoreUpdateResponsePayload struct {
	events.CommonMetadata
	Success     bool                `json:"success"`
	RoundID     sharedtypes.RoundID `json:"round_id,omitempty"`
	Participant Participant         `json:"participant,omitempty"`
	Error       string              `json:"error,omitempty"` // Include error details
}

// ScoreUpdateErrorPayload (Consider using later for more detailed error handling)
type ScoreUpdateErrorPayload struct {
	events.CommonMetadata
	CorrelationID string                     `json:"correlation_id"`
	Request       *ScoreUpdateRequestPayload `json:"score_update_request"` // Include original request
	Error         string                     `json:"error"`
}
