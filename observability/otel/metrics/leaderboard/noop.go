package leaderboardmetrics

import (
	"context"
	"time"

	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
)

// NoOpMetrics is a metrics collector that does nothing. Useful for unit tests.
type NoOpMetrics struct{}

func NewNoop() LeaderboardMetrics {
	return &NoOpMetrics{}
}

func (n *NoOpMetrics) RecordLeaderboardUpdate(ctx context.Context, success bool, source string, roundID sharedtypes.RoundID) {
}

func (n *NoOpMetrics) RecordTagAssignment(ctx context.Context, success bool, tagNumber sharedtypes.TagNumber, operationName string) {
}

func (n *NoOpMetrics) RecordTagAvailabilityCheck(ctx context.Context, available bool, tagNumber sharedtypes.TagNumber, serviceName string) {
}

func (n *NoOpMetrics) RecordOperationAttempt(ctx context.Context, operationName string, serviceName string) {
}

func (n *NoOpMetrics) RecordOperationSuccess(ctx context.Context, operationName string, serviceName string) {
}

func (n *NoOpMetrics) RecordOperationFailure(ctx context.Context, operationName string, serviceName string) {
}

func (n *NoOpMetrics) RecordOperationDuration(ctx context.Context, operationName string, serviceName string, duration time.Duration) {
}
func (n *NoOpMetrics) RecordServiceAttempt(ctx context.Context, serviceName string) {}
func (n *NoOpMetrics) RecordServiceSuccess(ctx context.Context, serviceName string) {}
func (n *NoOpMetrics) RecordServiceFailure(ctx context.Context, serviceName string) {}
func (n *NoOpMetrics) RecordServiceDuration(ctx context.Context, serviceName string, duration time.Duration) {
}

func (n *NoOpMetrics) RecordLeaderboardUpdateAttempt(ctx context.Context, roundID sharedtypes.RoundID, serviceName string) {
}

func (n *NoOpMetrics) RecordLeaderboardUpdateSuccess(ctx context.Context, roundID sharedtypes.RoundID, serviceName string) {
}

func (n *NoOpMetrics) RecordLeaderboardUpdateFailure(ctx context.Context, roundID sharedtypes.RoundID, serviceName string) {
}

func (n *NoOpMetrics) RecordLeaderboardUpdateDuration(ctx context.Context, serviceName string, duration time.Duration) {
}
func (n *NoOpMetrics) RecordLeaderboardGetAttempt(ctx context.Context, serviceName string) {}
func (n *NoOpMetrics) RecordLeaderboardGetSuccess(ctx context.Context, serviceName string) {}
func (n *NoOpMetrics) RecordLeaderboardGetFailure(ctx context.Context, serviceName string) {}
func (n *NoOpMetrics) RecordLeaderboardGetDuration(ctx context.Context, serviceName string, duration time.Duration) {
}
func (n *NoOpMetrics) RecordTagGetAttempt(ctx context.Context, serviceName string) {}
func (n *NoOpMetrics) RecordTagGetSuccess(ctx context.Context, serviceName string) {}
func (n *NoOpMetrics) RecordTagGetFailure(ctx context.Context, serviceName string) {}
func (n *NoOpMetrics) RecordTagGetDuration(ctx context.Context, serviceName string, duration time.Duration) {
}
func (n *NoOpMetrics) RecordTagAssignmentAttempt(ctx context.Context, operationName string)    {}
func (n *NoOpMetrics) RecordTagAssignmentSuccess(ctx context.Context, operationName string)    {}
func (n *NoOpMetrics) RecordTagAssignmentFailure(ctx context.Context, operationName string)    {}
func (n *NoOpMetrics) RecordTagAssignmentDuration(ctx context.Context, duration time.Duration) {}
func (n *NoOpMetrics) RecordTagSwapAttempt(ctx context.Context, requestorID, targetID sharedtypes.DiscordID) {
}

func (n *NoOpMetrics) RecordTagSwapSuccess(ctx context.Context, requestorID, targetID sharedtypes.DiscordID) {
}

func (n *NoOpMetrics) RecordTagSwapFailure(ctx context.Context, requestorID, targetID sharedtypes.DiscordID, reason string) {
}

func (n *NoOpMetrics) RecordTagAssignmentUpdate(ctx context.Context, oldTag, newTag sharedtypes.TagNumber, userID sharedtypes.DiscordID) {
}

func (n *NoOpMetrics) RecordNewTagAssignment(ctx context.Context, tag sharedtypes.TagNumber, userID sharedtypes.DiscordID) {
}

func (n *NoOpMetrics) RecordTagRemoval(ctx context.Context, tag sharedtypes.TagNumber, userID sharedtypes.DiscordID) {
}
func (n *NoOpMetrics) RecordHandlerAttempt(ctx context.Context, handlerName string) {}
func (n *NoOpMetrics) RecordHandlerSuccess(ctx context.Context, handlerName string) {}
func (n *NoOpMetrics) RecordHandlerFailure(ctx context.Context, handlerName string) {}
func (n *NoOpMetrics) RecordHandlerDuration(ctx context.Context, handlerName string, duration time.Duration) {
}
