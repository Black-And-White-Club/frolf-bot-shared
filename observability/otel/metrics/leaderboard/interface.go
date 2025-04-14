package leaderboardmetrics

import (
	"context"
	"time"

	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
)

// LeaderboardMetrics defines metrics specific to leaderboard operations
type LeaderboardMetrics interface {
	RecordLeaderboardUpdate(ctx context.Context, success bool, source string, roundID sharedtypes.RoundID)
	RecordTagAssignment(ctx context.Context, success bool, tagNumber sharedtypes.TagNumber, operationName string)
	RecordTagAvailabilityCheck(ctx context.Context, available bool, tagNumber sharedtypes.TagNumber, serviceName string)
	RecordOperationAttempt(ctx context.Context, operationName string, serviceName string)
	RecordOperationSuccess(ctx context.Context, operationName string, serviceName string)
	RecordOperationFailure(ctx context.Context, operationName string, serviceName string)
	RecordOperationDuration(ctx context.Context, operationName string, serviceName string, duration time.Duration)
	RecordServiceAttempt(ctx context.Context, serviceName string)
	RecordServiceSuccess(ctx context.Context, serviceName string)
	RecordServiceFailure(ctx context.Context, serviceName string)
	RecordServiceDuration(ctx context.Context, serviceName string, duration time.Duration)
	RecordLeaderboardUpdateAttempt(ctx context.Context, roundID sharedtypes.RoundID, serviceName string)
	RecordLeaderboardUpdateSuccess(ctx context.Context, roundID sharedtypes.RoundID, serviceName string)
	RecordLeaderboardUpdateFailure(ctx context.Context, roundID sharedtypes.RoundID, serviceName string)
	RecordLeaderboardUpdateDuration(ctx context.Context, serviceName string, duration time.Duration)
	RecordLeaderboardGetAttempt(ctx context.Context, serviceName string)
	RecordLeaderboardGetSuccess(ctx context.Context, serviceName string)
	RecordLeaderboardGetFailure(ctx context.Context, serviceName string)
	RecordLeaderboardGetDuration(ctx context.Context, serviceName string, duration time.Duration)
	RecordTagGetAttempt(ctx context.Context, serviceName string)
	RecordTagGetSuccess(ctx context.Context, serviceName string)
	RecordTagGetFailure(ctx context.Context, serviceName string)
	RecordTagGetDuration(ctx context.Context, serviceName string, duration time.Duration)
	RecordTagAssignmentAttempt(ctx context.Context, operationName string)
	RecordTagAssignmentSuccess(ctx context.Context, operationName string)
	RecordTagAssignmentFailure(ctx context.Context, operationName string)
	RecordTagAssignmentDuration(ctx context.Context, duration time.Duration)
	RecordTagSwapAttempt(ctx context.Context, requestorID, targetID sharedtypes.DiscordID)
	RecordTagSwapSuccess(ctx context.Context, requestorID, targetID sharedtypes.DiscordID)
	RecordTagSwapFailure(ctx context.Context, requestorID, targetID sharedtypes.DiscordID, reason string)
	RecordTagAssignmentUpdate(ctx context.Context, oldTag, newTag sharedtypes.TagNumber, userID sharedtypes.DiscordID)
	RecordNewTagAssignment(ctx context.Context, tag sharedtypes.TagNumber, userID sharedtypes.DiscordID)
	RecordTagRemoval(ctx context.Context, tag sharedtypes.TagNumber, userID sharedtypes.DiscordID)
	RecordHandlerAttempt(ctx context.Context, handlerName string)
	RecordHandlerSuccess(ctx context.Context, handlerName string)
	RecordHandlerFailure(ctx context.Context, handlerName string)
	RecordHandlerDuration(ctx context.Context, handlerName string, duration time.Duration)
}
