package scoremetrics

import (
	"context"
	"time"

	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
)

// ScoreMetrics defines metrics specific to score operations using OpenTelemetry
type ScoreMetrics interface {
	RecordScoreProcessingAttempt(ctx context.Context, roundID sharedtypes.RoundID)
	RecordScoreProcessingSuccess(ctx context.Context, roundID sharedtypes.RoundID)
	RecordScoreProcessingFailure(ctx context.Context, roundID sharedtypes.RoundID)
	RecordScoreProcessingDuration(ctx context.Context, roundID sharedtypes.RoundID, duration time.Duration)
	RecordScoreCorrectionAttempt(ctx context.Context, roundID sharedtypes.RoundID)
	RecordScoreCorrectionSuccess(ctx context.Context, roundID sharedtypes.RoundID)
	RecordScoreCorrectionFailure(ctx context.Context, roundID sharedtypes.RoundID)
	RecordScoreCorrectionDuration(ctx context.Context, roundID sharedtypes.RoundID, duration time.Duration)
	RecordLeaderboardUpdateAttempt(ctx context.Context, roundID sharedtypes.RoundID)
	RecordLeaderboardUpdateSuccess(ctx context.Context, roundID sharedtypes.RoundID)
	RecordLeaderboardUpdateFailure(ctx context.Context, roundID sharedtypes.RoundID)
	RecordLeaderboardUpdateDuration(ctx context.Context, roundID sharedtypes.RoundID, duration time.Duration)
	RecordOperationAttempt(ctx context.Context, operationName string, roundID sharedtypes.RoundID)
	RecordOperationSuccess(ctx context.Context, operationName string, roundID sharedtypes.RoundID)
	RecordOperationFailure(ctx context.Context, operationName string, roundID sharedtypes.RoundID)
	RecordOperationDuration(ctx context.Context, operationName string, duration time.Duration)
	RecordDBQueryDuration(ctx context.Context, duration time.Duration)
	RecordHandlerAttempt(ctx context.Context, handlerName string)
	RecordHandlerSuccess(ctx context.Context, handlerName string)
	RecordHandlerFailure(ctx context.Context, handlerName string)
	RecordHandlerDuration(ctx context.Context, handlerName string, duration time.Duration)
	RecordScoreUpdateAttempt(ctx context.Context, success bool, roundID sharedtypes.RoundID, userID sharedtypes.DiscordID)
	RecordRoundScoresProcessingAttempt(ctx context.Context, success bool, roundID sharedtypes.RoundID)
	RecordScoresProcessed(ctx context.Context, roundID sharedtypes.RoundID, numScores int)
	RecordTaggedPlayersProcessed(ctx context.Context, roundID sharedtypes.RoundID, numTaggedPlayers int)
	RecordUntaggedPlayersProcessed(ctx context.Context, roundID sharedtypes.RoundID, numUntaggedPlayers int)
	RecordScoreSortingDuration(ctx context.Context, roundID sharedtypes.RoundID, duration time.Duration)
	RecordTagExtractionDuration(ctx context.Context, roundID sharedtypes.RoundID, duration time.Duration)
	RecordPlayerScore(ctx context.Context, roundID sharedtypes.RoundID, userID sharedtypes.DiscordID, score sharedtypes.Score)
	RecordPlayerTag(ctx context.Context, roundID sharedtypes.RoundID, userID sharedtypes.DiscordID, tagNumber *sharedtypes.TagNumber)
	RecordTagPerformance(ctx context.Context, roundID sharedtypes.RoundID, tagNumber *sharedtypes.TagNumber, score sharedtypes.Score)
	RecordTagMovement(ctx context.Context, roundID sharedtypes.RoundID, tagNumber *sharedtypes.TagNumber, fromUserID, toUserID sharedtypes.DiscordID)
	RecordUntaggedPlayer(ctx context.Context, roundID sharedtypes.RoundID, userID sharedtypes.DiscordID)
	RecordPlayerScoresBatch(ctx context.Context, roundID sharedtypes.RoundID, playerScores map[sharedtypes.DiscordID]sharedtypes.Score)
	RecordPlayerTagsBatch(ctx context.Context, roundID sharedtypes.RoundID, playerTags map[sharedtypes.DiscordID]*sharedtypes.TagNumber)
	RecordTagPerformanceBatch(ctx context.Context, roundID sharedtypes.RoundID, tagPerformance map[*sharedtypes.TagNumber]sharedtypes.Score)
	RecordUntaggedPlayersBatch(ctx context.Context, roundID sharedtypes.RoundID, untaggedPlayers []sharedtypes.DiscordID)
}
