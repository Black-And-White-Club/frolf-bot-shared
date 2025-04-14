package scoremetrics

import (
	"context"
	"time"

	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
)

// NoOpMetrics is a metrics collector that does nothing. Useful for unit tests.
type NoOpMetrics struct{}

func NewNoop() ScoreMetrics {
	return &NoOpMetrics{}
}

func (n *NoOpMetrics) RecordScoreProcessingAttempt(ctx context.Context, roundID sharedtypes.RoundID) {
}

func (n *NoOpMetrics) RecordScoreProcessingSuccess(ctx context.Context, roundID sharedtypes.RoundID) {
}

func (n *NoOpMetrics) RecordScoreProcessingFailure(ctx context.Context, roundID sharedtypes.RoundID) {
}

func (n *NoOpMetrics) RecordScoreProcessingDuration(ctx context.Context, roundID sharedtypes.RoundID, duration time.Duration) {
}

func (n *NoOpMetrics) RecordScoreCorrectionAttempt(ctx context.Context, roundID sharedtypes.RoundID) {
}

func (n *NoOpMetrics) RecordScoreCorrectionSuccess(ctx context.Context, roundID sharedtypes.RoundID) {
}

func (n *NoOpMetrics) RecordScoreCorrectionFailure(ctx context.Context, roundID sharedtypes.RoundID) {
}

func (n *NoOpMetrics) RecordScoreCorrectionDuration(ctx context.Context, roundID sharedtypes.RoundID, duration time.Duration) {
}

func (n *NoOpMetrics) RecordLeaderboardUpdateAttempt(ctx context.Context, roundID sharedtypes.RoundID) {
}

func (n *NoOpMetrics) RecordLeaderboardUpdateSuccess(ctx context.Context, roundID sharedtypes.RoundID) {
}

func (n *NoOpMetrics) RecordLeaderboardUpdateFailure(ctx context.Context, roundID sharedtypes.RoundID) {
}

func (n *NoOpMetrics) RecordLeaderboardUpdateDuration(ctx context.Context, roundID sharedtypes.RoundID, duration time.Duration) {
}

func (n *NoOpMetrics) RecordOperationAttempt(ctx context.Context, operationName string, roundID sharedtypes.RoundID) {
}

func (n *NoOpMetrics) RecordOperationSuccess(ctx context.Context, operationName string, roundID sharedtypes.RoundID) {
}

func (n *NoOpMetrics) RecordOperationFailure(ctx context.Context, operationName string, roundID sharedtypes.RoundID) {
}

func (n *NoOpMetrics) RecordOperationDuration(ctx context.Context, operationName string, duration time.Duration) {
}
func (n *NoOpMetrics) RecordDBQueryDuration(ctx context.Context, duration time.Duration) {}
func (n *NoOpMetrics) RecordHandlerAttempt(ctx context.Context, handlerName string)      {}
func (n *NoOpMetrics) RecordHandlerSuccess(ctx context.Context, handlerName string)      {}
func (n *NoOpMetrics) RecordHandlerFailure(ctx context.Context, handlerName string)      {}
func (n *NoOpMetrics) RecordHandlerDuration(ctx context.Context, handlerName string, duration time.Duration) {
}

func (n *NoOpMetrics) RecordScoreUpdateAttempt(ctx context.Context, success bool, roundID sharedtypes.RoundID, userID sharedtypes.DiscordID) {
}

func (n *NoOpMetrics) RecordRoundScoresProcessingAttempt(ctx context.Context, success bool, roundID sharedtypes.RoundID) {
}

func (n *NoOpMetrics) RecordScoresProcessed(ctx context.Context, roundID sharedtypes.RoundID, numScores int) {
}

func (n *NoOpMetrics) RecordTaggedPlayersProcessed(ctx context.Context, roundID sharedtypes.RoundID, numTaggedPlayers int) {
}

func (n *NoOpMetrics) RecordUntaggedPlayersProcessed(ctx context.Context, roundID sharedtypes.RoundID, numUntaggedPlayers int) {
}

func (n *NoOpMetrics) RecordScoreSortingDuration(ctx context.Context, roundID sharedtypes.RoundID, duration time.Duration) {
}

func (n *NoOpMetrics) RecordTagExtractionDuration(ctx context.Context, roundID sharedtypes.RoundID, duration time.Duration) {
}

func (n *NoOpMetrics) RecordPlayerScore(ctx context.Context, roundID sharedtypes.RoundID, userID sharedtypes.DiscordID, score sharedtypes.Score) {
}

func (n *NoOpMetrics) RecordPlayerTag(ctx context.Context, roundID sharedtypes.RoundID, userID sharedtypes.DiscordID, tagNumber *sharedtypes.TagNumber) {
}

func (n *NoOpMetrics) RecordTagPerformance(ctx context.Context, roundID sharedtypes.RoundID, tagNumber *sharedtypes.TagNumber, score sharedtypes.Score) {
}

func (n *NoOpMetrics) RecordTagMovement(ctx context.Context, roundID sharedtypes.RoundID, tagNumber *sharedtypes.TagNumber, fromUserID, toUserID sharedtypes.DiscordID) {
}

func (n *NoOpMetrics) RecordUntaggedPlayer(ctx context.Context, roundID sharedtypes.RoundID, userID sharedtypes.DiscordID) {
}

func (n *NoOpMetrics) RecordPlayerScoresBatch(ctx context.Context, roundID sharedtypes.RoundID, playerScores map[sharedtypes.DiscordID]sharedtypes.Score) {
}

func (n *NoOpMetrics) RecordPlayerTagsBatch(ctx context.Context, roundID sharedtypes.RoundID, playerTags map[sharedtypes.DiscordID]*sharedtypes.TagNumber) {
}

func (n *NoOpMetrics) RecordTagPerformanceBatch(ctx context.Context, roundID sharedtypes.RoundID, tagPerformance map[*sharedtypes.TagNumber]sharedtypes.Score) {
}

func (n *NoOpMetrics) RecordUntaggedPlayersBatch(ctx context.Context, roundID sharedtypes.RoundID, untaggedPlayers []sharedtypes.DiscordID) {
}
