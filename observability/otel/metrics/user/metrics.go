// observability/prometheus/user/metrics.go
package usermetrics

import (
	"context"
	"time"

	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// RecordUserCreationAttempt records a user creation attempt
func (m *userMetrics) RecordUserCreationAttempt(ctx context.Context, userType string, source string) {
	m.userCreationAttemptCounter.Add(ctx, 1, metric.WithAttributes(
		userTypeAttr(userType),
		sourceAttr(source),
	))
}

// RecordUserCreationSuccess records a successful user creation
func (m *userMetrics) RecordUserCreationSuccess(ctx context.Context, userType string, source string) {
	m.userCreationSuccessCounter.Add(ctx, 1, metric.WithAttributes(
		userTypeAttr(userType),
		sourceAttr(source),
	))
}

// RecordUserCreationFailure records a failed user creation
func (m *userMetrics) RecordUserCreationFailure(ctx context.Context, userType string, source string) {
	m.userCreationFailureCounter.Add(ctx, 1, metric.WithAttributes(
		userTypeAttr(userType),
		sourceAttr(source),
	))
}

// RecordUserCreationDuration records the duration of user creation
func (m *userMetrics) RecordUserCreationDuration(ctx context.Context, userType string, source string, duration time.Duration) {
	m.userCreationDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(
		userTypeAttr(userType),
		sourceAttr(source),
	))
}

// RecordUserRetrievalAttempt records a user retrieval attempt
func (m *userMetrics) RecordUserRetrievalAttempt(ctx context.Context, userID sharedtypes.DiscordID) {
	m.userRetrievalAttemptCounter.Add(ctx, 1, metric.WithAttributes(
		userIDAttr(userID),
	))
}

// RecordUserRetrievalSuccess records a successful user retrieval
func (m *userMetrics) RecordUserRetrievalSuccess(ctx context.Context, userID sharedtypes.DiscordID) {
	m.userRetrievalSuccessCounter.Add(ctx, 1, metric.WithAttributes(
		userIDAttr(userID),
	))
}

// RecordUserRetrievalFailure records a failed user retrieval
func (m *userMetrics) RecordUserRetrievalFailure(ctx context.Context, userID sharedtypes.DiscordID) {
	m.userRetrievalFailureCounter.Add(ctx, 1, metric.WithAttributes(
		userIDAttr(userID),
	))
}

// RecordUserRetrievalDuration records the duration of user retrieval
func (m *userMetrics) RecordUserRetrievalDuration(ctx context.Context, userID sharedtypes.DiscordID, duration time.Duration) {
	m.userRetrievalDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(
		userIDAttr(userID),
	))
}

// RecordUserRoleRetrievalAttempt records a user role retrieval attempt
func (m *userMetrics) RecordUserRoleRetrievalAttempt(ctx context.Context, userID sharedtypes.DiscordID) {
	m.userRoleRetrievalAttemptCounter.Add(ctx, 1, metric.WithAttributes(
		userIDAttr(userID),
	))
}

// RecordUserRoleRetrievalSuccess records a successful user role retrieval
func (m *userMetrics) RecordUserRoleRetrievalSuccess(ctx context.Context, userID sharedtypes.DiscordID) {
	m.userRoleRetrievalSuccessCounter.Add(ctx, 1, metric.WithAttributes(
		userIDAttr(userID),
	))
}

// RecordUserRoleRetrievalFailure records a failed user role retrieval
func (m *userMetrics) RecordUserRoleRetrievalFailure(ctx context.Context, userID sharedtypes.DiscordID) {
	m.userRoleRetrievalFailureCounter.Add(ctx, 1, metric.WithAttributes(
		userIDAttr(userID),
	))
}

// RecordUserRoleRetrievalDuration records the duration of user role retrieval
func (m *userMetrics) RecordUserRoleRetrievalDuration(ctx context.Context, userID sharedtypes.DiscordID, duration time.Duration) {
	m.userRoleRetrievalDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(
		userIDAttr(userID),
	))
}

// RecordTagAvailabilityCheck records the availability status of a tag
func (m *userMetrics) RecordTagAvailabilityCheck(ctx context.Context, available bool, tag sharedtypes.TagNumber) {
	attrs := []attribute.KeyValue{tagNumberAttr(tag)}
	if available {
		attrs = append(attrs, attribute.Bool("available", true))
		m.tagAvailableCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
	} else {
		attrs = append(attrs, attribute.Bool("available", false))
		m.tagUnavailableCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
	}
}

// RecordPermissionCheckAttempt records a permission check attempt
func (m *userMetrics) RecordPermissionCheckAttempt(ctx context.Context, role sharedtypes.UserRoleEnum, action string, resource string) {
	m.permissionCheckAttemptCounter.Add(ctx, 1, metric.WithAttributes(
		roleAttr(role),
		actionAttr(action),
		resourceAttr(resource),
	))
}

// RecordPermissionCheckSuccess records a successful permission check
func (m *userMetrics) RecordPermissionCheckSuccess(ctx context.Context, role sharedtypes.UserRoleEnum, action string, resource string) {
	m.permissionCheckSuccessCounter.Add(ctx, 1, metric.WithAttributes(
		roleAttr(role),
		actionAttr(action),
		resourceAttr(resource),
	))
}

// RecordPermissionCheckFailure records a failed permission check
func (m *userMetrics) RecordPermissionCheckFailure(ctx context.Context, role sharedtypes.UserRoleEnum, action string, resource string) {
	m.permissionCheckFailureCounter.Add(ctx, 1, metric.WithAttributes(
		roleAttr(role),
		actionAttr(action),
		resourceAttr(resource),
	))
}

// RecordPermissionCheckDuration records the duration of a permission check
func (m *userMetrics) RecordPermissionCheckDuration(ctx context.Context, role sharedtypes.UserRoleEnum, action string, resource string, duration time.Duration) {
	m.permissionCheckDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(
		roleAttr(role),
		actionAttr(action),
		resourceAttr(resource),
	))
}

// RecordRoleUpdateAttempt records a role update attempt
func (m *userMetrics) RecordRoleUpdateAttempt(ctx context.Context, userID sharedtypes.DiscordID, oldRole, newRole sharedtypes.UserRoleEnum) {
	m.roleUpdateAttemptCounter.Add(ctx, 1, metric.WithAttributes(
		userIDAttr(userID),
		attribute.String("old_role", oldRole.String()),
		attribute.String("new_role", newRole.String()),
	))
}

// RecordRoleUpdateSuccess records a successful role update
func (m *userMetrics) RecordRoleUpdateSuccess(ctx context.Context, userID sharedtypes.DiscordID, oldRole, newRole sharedtypes.UserRoleEnum) {
	m.roleUpdateSuccessCounter.Add(ctx, 1, metric.WithAttributes(
		userIDAttr(userID),
		attribute.String("old_role", oldRole.String()),
		attribute.String("new_role", newRole.String()),
	))
}

// RecordRoleUpdateFailure records a failed role update
func (m *userMetrics) RecordRoleUpdateFailure(ctx context.Context, userID sharedtypes.DiscordID, oldRole, newRole sharedtypes.UserRoleEnum) {
	m.roleUpdateFailureCounter.Add(ctx, 1, metric.WithAttributes(
		userIDAttr(userID),
		attribute.String("old_role", oldRole.String()),
		attribute.String("new_role", newRole.String()),
	))
}

// RecordRoleUpdateDuration records the duration of a role update
func (m *userMetrics) RecordRoleUpdateDuration(ctx context.Context, userID sharedtypes.DiscordID, oldRole, newRole sharedtypes.UserRoleEnum, duration time.Duration) {
	m.roleUpdateDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(
		userIDAttr(userID),
		attribute.String("old_role", oldRole.String()),
		attribute.String("new_role", newRole.String()),
	))
}

// RecordDBQueryDuration records the duration of a database query
func (m *userMetrics) RecordDBQueryDuration(ctx context.Context, duration time.Duration) {
	m.dbQueryDuration.Record(ctx, duration.Seconds())
}

// RecordUserCreationByTag records a user creation event associated with a tag
func (m *userMetrics) RecordUserCreationByTag(ctx context.Context, tag sharedtypes.TagNumber) {
	m.userCreationByTagCounter.Add(ctx, 1, metric.WithAttributes(
		tagNumberAttr(tag),
	))
}

// RecordOperationAttempt records an operation attempt
func (m *userMetrics) RecordOperationAttempt(ctx context.Context, operationName string, userID sharedtypes.DiscordID) {
	m.userOperationAttemptCounter.Add(ctx, 1, metric.WithAttributes(
		operationAttr(operationName),
		userIDAttr(userID),
	))
}

// RecordOperationSuccess records a successful operation
func (m *userMetrics) RecordOperationSuccess(ctx context.Context, operationName string, userID sharedtypes.DiscordID) {
	m.userOperationSuccessCounter.Add(ctx, 1, metric.WithAttributes(
		operationAttr(operationName),
		userIDAttr(userID),
	))
}

// RecordOperationFailure records a failed operation
func (m *userMetrics) RecordOperationFailure(ctx context.Context, operationName string, userID sharedtypes.DiscordID) {
	m.userOperationFailureCounter.Add(ctx, 1, metric.WithAttributes(
		operationAttr(operationName),
		userIDAttr(userID),
	))
}

// RecordOperationDuration records the duration of an operation
func (m *userMetrics) RecordOperationDuration(ctx context.Context, operationName string, duration time.Duration, userID sharedtypes.DiscordID) {
	m.userOperationDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(
		operationAttr(operationName),
		userIDAttr(userID),
	))
}

// RecordUserRole records the current role of a user (using a counter for changes)
func (m *userMetrics) RecordUserRole(ctx context.Context, userID sharedtypes.DiscordID, role sharedtypes.UserRoleEnum) {
	m.roleChangedTotal.Add(ctx, 1, metric.WithAttributes(
		userIDAttr(userID),
		attribute.String("new_role", role.String()),
	))
}

// RecordHandlerAttempt records a handler attempt
func (m *userMetrics) RecordHandlerAttempt(ctx context.Context, handlerName string) {
	m.handlerAttemptCounter.Add(ctx, 1, metric.WithAttributes(handlerAttrs(handlerName)))
}

// RecordHandlerSuccess records a successful handler attempt
func (m *userMetrics) RecordHandlerSuccess(ctx context.Context, handlerName string) {
	m.handlerSuccessCounter.Add(ctx, 1, metric.WithAttributes(handlerAttrs(handlerName)))
}

// RecordHandlerFailure records a failed handler attempt
func (m *userMetrics) RecordHandlerFailure(ctx context.Context, handlerName string) {
	m.handlerFailureCounter.Add(ctx, 1, metric.WithAttributes(handlerAttrs(handlerName)))
}

// RecordHandlerDuration records the duration of a handler execution
func (m *userMetrics) RecordHandlerDuration(ctx context.Context, handlerName string, duration time.Duration) {
	m.handlerDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(handlerAttrs(handlerName)))
}
