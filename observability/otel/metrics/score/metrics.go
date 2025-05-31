package scoremetrics

import (
	"context"
	"time"

	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared" // Assuming this path is correct
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// RecordScoreProcessingAttempt records a score processing attempt
func (m *scoreMetrics) RecordScoreProcessingAttempt(ctx context.Context, roundID sharedtypes.RoundID) {
	m.scoreProcessingAttemptCounter.Add(ctx, 1, metric.WithAttributes(roundAttrs(roundID)))
}

// RecordScoreProcessingSuccess records a successful score processing attempt
func (m *scoreMetrics) RecordScoreProcessingSuccess(ctx context.Context, roundID sharedtypes.RoundID) {
	m.scoreProcessingSuccessCounter.Add(ctx, 1, metric.WithAttributes(roundAttrs(roundID)))
}

// RecordScoreProcessingFailure records a failed score processing attempt
func (m *scoreMetrics) RecordScoreProcessingFailure(ctx context.Context, roundID sharedtypes.RoundID) {
	m.scoreProcessingFailureCounter.Add(ctx, 1, metric.WithAttributes(roundAttrs(roundID)))
}

// RecordScoreProcessingDuration records the duration of score processing
func (m *scoreMetrics) RecordScoreProcessingDuration(ctx context.Context, roundID sharedtypes.RoundID, duration time.Duration) {
	m.scoreProcessingDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(roundAttrs(roundID)))
}

// RecordScoreCorrectionAttempt records a score correction attempt
func (m *scoreMetrics) RecordScoreCorrectionAttempt(ctx context.Context, roundID sharedtypes.RoundID) {
	m.scoreCorrectionAttemptCounter.Add(ctx, 1, metric.WithAttributes(roundAttrs(roundID)))
}

// RecordScoreCorrectionSuccess records a successful score correction
func (m *scoreMetrics) RecordScoreCorrectionSuccess(ctx context.Context, roundID sharedtypes.RoundID) {
	m.scoreCorrectionSuccessCounter.Add(ctx, 1, metric.WithAttributes(roundAttrs(roundID)))
}

// RecordScoreCorrectionFailure records a failed score correction
func (m *scoreMetrics) RecordScoreCorrectionFailure(ctx context.Context, roundID sharedtypes.RoundID) {
	m.scoreCorrectionFailureCounter.Add(ctx, 1, metric.WithAttributes(roundAttrs(roundID)))
}

// RecordScoreCorrectionDuration records the duration of score correction
func (m *scoreMetrics) RecordScoreCorrectionDuration(ctx context.Context, roundID sharedtypes.RoundID, duration time.Duration) {
	m.scoreCorrectionDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(roundAttrs(roundID)))
}

// RecordLeaderboardUpdateAttempt records a leaderboard update attempt
func (m *scoreMetrics) RecordLeaderboardUpdateAttempt(ctx context.Context, roundID sharedtypes.RoundID) {
	m.leaderboardUpdateAttemptCounter.Add(ctx, 1, metric.WithAttributes(roundAttrs(roundID)))
}

// RecordLeaderboardUpdateSuccess records a successful leaderboard update
func (m *scoreMetrics) RecordLeaderboardUpdateSuccess(ctx context.Context, roundID sharedtypes.RoundID) {
	m.leaderboardUpdateSuccessCounter.Add(ctx, 1, metric.WithAttributes(roundAttrs(roundID)))
}

// RecordLeaderboardUpdateFailure records a failed leaderboard update
func (m *scoreMetrics) RecordLeaderboardUpdateFailure(ctx context.Context, roundID sharedtypes.RoundID) {
	m.leaderboardUpdateFailureCounter.Add(ctx, 1, metric.WithAttributes(roundAttrs(roundID)))
}

// RecordLeaderboardUpdateDuration records the duration of a leaderboard update
func (m *scoreMetrics) RecordLeaderboardUpdateDuration(ctx context.Context, roundID sharedtypes.RoundID, duration time.Duration) {
	m.leaderboardUpdateDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(roundAttrs(roundID)))
}

// RecordDBQueryDuration records the duration of a database query
func (m *scoreMetrics) RecordDBQueryDuration(ctx context.Context, duration time.Duration) {
	// No specific attributes mentioned in the original, add if needed
	m.dbQueryDuration.Record(ctx, duration.Seconds())
}

// RecordOperationAttempt records an operation attempt
func (m *scoreMetrics) RecordOperationAttempt(ctx context.Context, operationName string, roundID sharedtypes.RoundID) {
	m.scoreOperationAttemptCounter.Add(ctx, 1, metric.WithAttributes(operationRoundAttrs(operationName, roundID)...))
}

// RecordOperationSuccess records a successful operation
func (m *scoreMetrics) RecordOperationSuccess(ctx context.Context, operationName string, roundID sharedtypes.RoundID) {
	m.scoreOperationSuccessCounter.Add(ctx, 1, metric.WithAttributes(operationRoundAttrs(operationName, roundID)...))
}

// RecordOperationFailure records a failed operation
func (m *scoreMetrics) RecordOperationFailure(ctx context.Context, operationName string, roundID sharedtypes.RoundID) {
	m.scoreOperationFailureCounter.Add(ctx, 1, metric.WithAttributes(operationRoundAttrs(operationName, roundID)...))
}

// RecordOperationDuration records the duration of an operation
func (m *scoreMetrics) RecordOperationDuration(ctx context.Context, operationName string, duration time.Duration) {
	m.scoreOperationDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(operationAttrs(operationName)))
}

// RecordHandlerAttempt records a handler attempt
func (m *scoreMetrics) RecordHandlerAttempt(ctx context.Context, handlerName string) {
	m.handlerAttemptCounter.Add(ctx, 1, metric.WithAttributes(handlerAttrs(handlerName)))
}

// RecordHandlerSuccess records a successful handler attempt
func (m *scoreMetrics) RecordHandlerSuccess(ctx context.Context, handlerName string) {
	m.handlerSuccessCounter.Add(ctx, 1, metric.WithAttributes(handlerAttrs(handlerName)))
}

// RecordHandlerFailure records a failed handler attempt
func (m *scoreMetrics) RecordHandlerFailure(ctx context.Context, handlerName string) {
	m.handlerFailureCounter.Add(ctx, 1, metric.WithAttributes(handlerAttrs(handlerName)))
}

// RecordHandlerDuration records the duration of a handler execution
func (m *scoreMetrics) RecordHandlerDuration(ctx context.Context, handlerName string, duration time.Duration) {
	m.handlerDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(handlerAttrs(handlerName)))
}

// RecordScoreUpdateAttempt records a score update attempt
func (m *scoreMetrics) RecordScoreUpdateAttempt(ctx context.Context, success bool, roundID sharedtypes.RoundID, userID sharedtypes.DiscordID) {
	// Add success as an attribute if needed: attribute.Bool("success", success)
	attrs := roundUserAttrs(roundID, userID)
	if success {
		attrs = append(attrs, attribute.Bool("success", true))
	} else {
		attrs = append(attrs, attribute.Bool("success", false))
	}
	m.scoreUpdateAttemptCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
}

// RecordRoundScoresProcessingAttempt records a round scores processing attempt
func (m *scoreMetrics) RecordRoundScoresProcessingAttempt(ctx context.Context, success bool, roundID sharedtypes.RoundID) {
	// Add success as an attribute if needed: attribute.Bool("success", success)
	attrs := []attribute.KeyValue{roundIDAttr(roundID)}
	if success {
		attrs = append(attrs, attribute.Bool("success", true))
	} else {
		attrs = append(attrs, attribute.Bool("success", false))
	}
	m.roundScoresProcessingAttemptCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
}

// RecordScoresProcessed records the number of scores processed for a round
func (m *scoreMetrics) RecordScoresProcessed(ctx context.Context, roundID sharedtypes.RoundID, numScores int) {
	m.scoresProcessedCounter.Add(ctx, int64(numScores), metric.WithAttributes(roundAttrs(roundID)))
}

// RecordTaggedPlayersProcessed records the number of tagged players processed for a round
func (m *scoreMetrics) RecordTaggedPlayersProcessed(ctx context.Context, roundID sharedtypes.RoundID, numTaggedPlayers int) {
	m.taggedPlayersProcessedCounter.Add(ctx, int64(numTaggedPlayers), metric.WithAttributes(roundAttrs(roundID)))
}

// RecordUntaggedPlayersProcessed records the number of untagged players processed for a round
func (m *scoreMetrics) RecordUntaggedPlayersProcessed(ctx context.Context, roundID sharedtypes.RoundID, numUntaggedPlayers int) {
	m.untaggedPlayersProcessedCounter.Add(ctx, int64(numUntaggedPlayers), metric.WithAttributes(roundAttrs(roundID)))
}

// RecordScoreSortingDuration records the duration of score sorting for a round
func (m *scoreMetrics) RecordScoreSortingDuration(ctx context.Context, roundID sharedtypes.RoundID, duration time.Duration) {
	m.scoreSortingDurationHistogram.Record(ctx, duration.Seconds(), metric.WithAttributes(roundAttrs(roundID)))
}

// RecordTagExtractionDuration records the duration of tag extraction for a round
func (m *scoreMetrics) RecordTagExtractionDuration(ctx context.Context, roundID sharedtypes.RoundID, duration time.Duration) {
	m.tagExtractionDurationHistogram.Record(ctx, duration.Seconds(), metric.WithAttributes(roundAttrs(roundID)))
}

// RecordPlayerScore records a player's score for a round
func (m *scoreMetrics) RecordPlayerScore(ctx context.Context, roundID sharedtypes.RoundID, userID sharedtypes.DiscordID, score sharedtypes.Score) {
	m.playerScoreGauge.Record(ctx, int64(score), metric.WithAttributes(roundUserAttrs(roundID, userID)...))
}

// RecordPlayerTag records a player's tag number for a round
func (m *scoreMetrics) RecordPlayerTag(ctx context.Context, roundID sharedtypes.RoundID, userID sharedtypes.DiscordID, tagNumber *sharedtypes.TagNumber) {
	if tagNumber == nil {
		return // Don't record if tagNumber is nil
	}
	attrs := roundUserAttrs(roundID, userID)
	m.playerTagGauge.Record(ctx, int64(*tagNumber), metric.WithAttributes(attrs...))
}

// RecordTagPerformance records the score achieved by a specific tag in a round
func (m *scoreMetrics) RecordTagPerformance(ctx context.Context, roundID sharedtypes.RoundID, tagNumber *sharedtypes.TagNumber, score sharedtypes.Score) {
	if tagNumber == nil {
		return // Don't record if tagNumber is nil
	}
	attrs := tagAttrs(roundID, tagNumber)
	m.tagPerformanceGauge.Record(ctx, int64(score), metric.WithAttributes(attrs...))
}

// RecordTagMovement records the movement of a tag from one player to another
func (m *scoreMetrics) RecordTagMovement(ctx context.Context, roundID sharedtypes.RoundID, tagNumber *sharedtypes.TagNumber, fromUserID, toUserID sharedtypes.DiscordID) {
	if tagNumber == nil {
		// Handle nil tagNumber if necessary, maybe record without tag_number attribute
		// Or skip recording as done here.
		return
	}
	attrs := tagMovementAttrs(roundID, tagNumber, fromUserID, toUserID)
	m.tagMovementCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
}

// RecordUntaggedPlayer records a player without a tag in a round
func (m *scoreMetrics) RecordUntaggedPlayer(ctx context.Context, roundID sharedtypes.RoundID, userID sharedtypes.DiscordID) {
	m.untaggedPlayerCounter.Add(ctx, 1, metric.WithAttributes(roundUserAttrs(roundID, userID)...))
}

// --- Batch Recording Methods ---

func (m *scoreMetrics) RecordPlayerScoresBatch(ctx context.Context, roundID sharedtypes.RoundID, playerScores map[sharedtypes.DiscordID]sharedtypes.Score) {
	for userID, score := range playerScores {
		m.playerScoreGauge.Record(ctx, int64(score), metric.WithAttributes(roundUserAttrs(roundID, userID)...))
	}
}

func (m *scoreMetrics) RecordPlayerTagsBatch(ctx context.Context, roundID sharedtypes.RoundID, playerTags map[sharedtypes.DiscordID]*sharedtypes.TagNumber) {
	for userID, tagNumber := range playerTags {
		if tagNumber == nil {
			continue // Skip nil tags
		}
		attrs := roundUserAttrs(roundID, userID)
		m.playerTagGauge.Record(ctx, int64(*tagNumber), metric.WithAttributes(attrs...))
	}
}

func (m *scoreMetrics) RecordTagPerformanceBatch(ctx context.Context, roundID sharedtypes.RoundID, tagPerformance map[*sharedtypes.TagNumber]sharedtypes.Score) {
	for tagNumber, score := range tagPerformance {
		if tagNumber == nil {
			continue // Skip nil tags
		}
		attrs := tagAttrs(roundID, tagNumber)
		m.tagPerformanceGauge.Record(ctx, int64(score), metric.WithAttributes(attrs...))
	}
}

func (m *scoreMetrics) RecordUntaggedPlayersBatch(ctx context.Context, roundID sharedtypes.RoundID, untaggedPlayers []sharedtypes.DiscordID) {
	for _, userID := range untaggedPlayers {
		m.untaggedPlayerCounter.Add(ctx, 1, metric.WithAttributes(roundUserAttrs(roundID, userID)...))
	}
}
