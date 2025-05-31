package leaderboardmetrics

import (
	"context"
	"time"

	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
	"go.opentelemetry.io/otel/metric"
)

func (m *leaderboardMetrics) RecordLeaderboardUpdate(ctx context.Context, success bool, source string, roundID sharedtypes.RoundID) {
	m.leaderboardUpdateCounter.Add(ctx, 1, metric.WithAttributes(
		successAttr(success),
		sourceAttr(source),
		roundIDAttr(roundID),
	))
}

func (m *leaderboardMetrics) RecordTagAssignment(ctx context.Context, success bool, tagNumber sharedtypes.TagNumber, operationName string) {
	m.tagAssignmentCounter.Add(ctx, 1, metric.WithAttributes(
		successAttr(success),
		tagNumberAttr(tagNumber),
		operationAttr(operationName),
	))
}

func (m *leaderboardMetrics) RecordTagAvailabilityCheck(ctx context.Context, available bool, tagNumber sharedtypes.TagNumber, serviceName string) {
	m.tagAvailabilityCounter.Add(ctx, 1, metric.WithAttributes(
		availableAttr(available),
		tagNumberAttr(tagNumber),
		serviceAttr(serviceName),
	))
}

func (m *leaderboardMetrics) RecordOperationAttempt(ctx context.Context, operationName string, serviceName string) {
	m.operationAttemptCounter.Add(ctx, 1, metric.WithAttributes(
		operationAttr(operationName),
		serviceAttr(serviceName),
	))
}

func (m *leaderboardMetrics) RecordOperationSuccess(ctx context.Context, operationName string, serviceName string) {
	m.operationSuccessCounter.Add(ctx, 1, metric.WithAttributes(
		operationAttr(operationName),
		serviceAttr(serviceName),
	))
}

func (m *leaderboardMetrics) RecordOperationFailure(ctx context.Context, operationName string, serviceName string) {
	m.operationFailureCounter.Add(ctx, 1, metric.WithAttributes(
		operationAttr(operationName),
		serviceAttr(serviceName),
	))
}

func (m *leaderboardMetrics) RecordOperationDuration(ctx context.Context, operationName string, serviceName string, duration time.Duration) {
	m.operationDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(
		operationAttr(operationName),
		serviceAttr(serviceName),
	))
}

func (m *leaderboardMetrics) RecordServiceAttempt(ctx context.Context, serviceName string) {
	m.serviceAttemptCounter.Add(ctx, 1, metric.WithAttributes(serviceAttr(serviceName)))
}

func (m *leaderboardMetrics) RecordServiceSuccess(ctx context.Context, serviceName string) {
	m.serviceSuccessCounter.Add(ctx, 1, metric.WithAttributes(serviceAttr(serviceName)))
}

func (m *leaderboardMetrics) RecordServiceFailure(ctx context.Context, serviceName string) {
	m.serviceFailureCounter.Add(ctx, 1, metric.WithAttributes(serviceAttr(serviceName)))
}

func (m *leaderboardMetrics) RecordServiceDuration(ctx context.Context, serviceName string, duration time.Duration) {
	m.serviceDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(serviceAttr(serviceName)))
}

func (m *leaderboardMetrics) RecordLeaderboardUpdateAttempt(ctx context.Context, roundID sharedtypes.RoundID, serviceName string) {
	m.leaderboardUpdateAttemptCounter.Add(ctx, 1, metric.WithAttributes(
		roundIDAttr(roundID),
		serviceAttr(serviceName),
	))
}

func (m *leaderboardMetrics) RecordLeaderboardUpdateSuccess(ctx context.Context, roundID sharedtypes.RoundID, serviceName string) {
	m.leaderboardUpdateSuccessCounter.Add(ctx, 1, metric.WithAttributes(
		roundIDAttr(roundID),
		serviceAttr(serviceName),
	))
}

func (m *leaderboardMetrics) RecordLeaderboardUpdateFailure(ctx context.Context, roundID sharedtypes.RoundID, serviceName string) {
	m.leaderboardUpdateFailureCounter.Add(ctx, 1, metric.WithAttributes(
		roundIDAttr(roundID),
		serviceAttr(serviceName),
	))
}

func (m *leaderboardMetrics) RecordLeaderboardUpdateDuration(ctx context.Context, serviceName string, duration time.Duration) {
	m.leaderboardUpdateDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(serviceAttr(serviceName)))
}

func (m *leaderboardMetrics) RecordLeaderboardGetAttempt(ctx context.Context, serviceName string) {
	m.leaderboardGetAttemptCounter.Add(ctx, 1, metric.WithAttributes(serviceAttr(serviceName)))
}

func (m *leaderboardMetrics) RecordLeaderboardGetSuccess(ctx context.Context, serviceName string) {
	m.leaderboardGetSuccessCounter.Add(ctx, 1, metric.WithAttributes(serviceAttr(serviceName)))
}

func (m *leaderboardMetrics) RecordLeaderboardGetFailure(ctx context.Context, serviceName string) {
	m.leaderboardGetFailureCounter.Add(ctx, 1, metric.WithAttributes(serviceAttr(serviceName)))
}

func (m *leaderboardMetrics) RecordLeaderboardGetDuration(ctx context.Context, serviceName string, duration time.Duration) {
	m.leaderboardGetDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(serviceAttr(serviceName)))
}

func (m *leaderboardMetrics) RecordTagGetAttempt(ctx context.Context, serviceName string) {
	m.tagGetAttemptCounter.Add(ctx, 1, metric.WithAttributes(serviceAttr(serviceName)))
}

func (m *leaderboardMetrics) RecordTagGetSuccess(ctx context.Context, serviceName string) {
	m.tagGetSuccessCounter.Add(ctx, 1, metric.WithAttributes(serviceAttr(serviceName)))
}

func (m *leaderboardMetrics) RecordTagGetFailure(ctx context.Context, serviceName string) {
	m.tagGetFailureCounter.Add(ctx, 1, metric.WithAttributes(serviceAttr(serviceName)))
}

func (m *leaderboardMetrics) RecordTagGetDuration(ctx context.Context, serviceName string, duration time.Duration) {
	m.tagGetDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(serviceAttr(serviceName)))
}

func (m *leaderboardMetrics) RecordTagAssignmentAttempt(ctx context.Context, operationName string) {
	m.tagAssignmentAttemptCounter.Add(ctx, 1, metric.WithAttributes(operationAttr(operationName)))
}

func (m *leaderboardMetrics) RecordTagAssignmentSuccess(ctx context.Context, operationName string) {
	m.tagAssignmentSuccessCounter.Add(ctx, 1, metric.WithAttributes(operationAttr(operationName)))
}

func (m *leaderboardMetrics) RecordTagAssignmentFailure(ctx context.Context, operationName string) {
	m.tagAssignmentFailureCounter.Add(ctx, 1, metric.WithAttributes(operationAttr(operationName)))
}

func (m *leaderboardMetrics) RecordTagAssignmentDuration(ctx context.Context, duration time.Duration) {
	m.tagAssignmentDuration.Record(ctx, duration.Seconds())
}

func (m *leaderboardMetrics) RecordTagSwapAttempt(ctx context.Context, requestorID, targetID sharedtypes.DiscordID) {
	m.tagSwapAttemptCounter.Add(ctx, 1, metric.WithAttributes(
		requestorIDAttr(requestorID),
		targetIDAttr(targetID),
	))
}

func (m *leaderboardMetrics) RecordTagSwapSuccess(ctx context.Context, requestorID, targetID sharedtypes.DiscordID) {
	m.tagSwapSuccessCounter.Add(ctx, 1, metric.WithAttributes(
		requestorIDAttr(requestorID),
		targetIDAttr(targetID),
	))
}

func (m *leaderboardMetrics) RecordTagSwapFailure(ctx context.Context, requestorID, targetID sharedtypes.DiscordID, reason string) {
	m.tagSwapFailureCounter.Add(ctx, 1, metric.WithAttributes(
		requestorIDAttr(requestorID),
		targetIDAttr(targetID),
		reasonAttr(reason),
	))
}

func (m *leaderboardMetrics) RecordTagAssignmentUpdate(ctx context.Context, oldTag, newTag sharedtypes.TagNumber, userID sharedtypes.DiscordID) {
	m.tagAssignmentUpdateCounter.Add(ctx, 1, metric.WithAttributes(
		oldTagAttr(oldTag),
		newTagAttr(newTag),
		userIDAttr(userID),
	))
}

func (m *leaderboardMetrics) RecordNewTagAssignment(ctx context.Context, tag sharedtypes.TagNumber, userID sharedtypes.DiscordID) {
	m.newTagAssignmentCounter.Add(ctx, 1, metric.WithAttributes(
		tagNumberAttr(tag),
		userIDAttr(userID),
	))
}

func (m *leaderboardMetrics) RecordTagRemoval(ctx context.Context, tag sharedtypes.TagNumber, userID sharedtypes.DiscordID) {
	m.tagRemovalCounter.Add(ctx, 1, metric.WithAttributes(
		tagNumberAttr(tag),
		userIDAttr(userID),
	))
}

func (m *leaderboardMetrics) RecordHandlerAttempt(ctx context.Context, handlerName string) {
	m.handlerAttemptCounter.Add(ctx, 1, metric.WithAttributes(handlerAttr(handlerName)))
}

func (m *leaderboardMetrics) RecordHandlerSuccess(ctx context.Context, handlerName string) {
	m.handlerSuccessCounter.Add(ctx, 1, metric.WithAttributes(handlerAttr(handlerName)))
}

func (m *leaderboardMetrics) RecordHandlerFailure(ctx context.Context, handlerName string) {
	m.handlerFailureCounter.Add(ctx, 1, metric.WithAttributes(handlerAttr(handlerName)))
}

func (m *leaderboardMetrics) RecordHandlerDuration(ctx context.Context, handlerName string, duration time.Duration) {
	m.handlerDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(handlerAttr(handlerName)))
}
