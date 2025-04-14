package usermetrics

import (
	"context"
	"time"

	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
)

// NoOpMetrics is a metrics collector that does nothing. Useful for unit tests.
type NoOpMetrics struct{}

func NewNoop() UserMetrics {
	return &NoOpMetrics{}
}

func (n *NoOpMetrics) RecordUserCreationAttempt(ctx context.Context, userType string, source string) {
}

func (n *NoOpMetrics) RecordUserCreationSuccess(ctx context.Context, userType string, source string) {
}

func (n *NoOpMetrics) RecordUserCreationFailure(ctx context.Context, userType string, source string) {
}

func (n *NoOpMetrics) RecordUserCreationDuration(ctx context.Context, userType string, source string, duration time.Duration) {
}
func (n *NoOpMetrics) RecordUserRetrievalAttempt(ctx context.Context, userID sharedtypes.DiscordID) {}
func (n *NoOpMetrics) RecordUserRetrievalSuccess(ctx context.Context, userID sharedtypes.DiscordID) {}
func (n *NoOpMetrics) RecordUserRetrievalFailure(ctx context.Context, userID sharedtypes.DiscordID) {}
func (n *NoOpMetrics) RecordUserRetrievalDuration(ctx context.Context, userID sharedtypes.DiscordID, duration time.Duration) {
}

func (n *NoOpMetrics) RecordUserRoleRetrievalAttempt(ctx context.Context, userID sharedtypes.DiscordID) {
}

func (n *NoOpMetrics) RecordUserRoleRetrievalSuccess(ctx context.Context, userID sharedtypes.DiscordID) {
}

func (n *NoOpMetrics) RecordUserRoleRetrievalFailure(ctx context.Context, userID sharedtypes.DiscordID) {
}

func (n *NoOpMetrics) RecordUserRoleRetrievalDuration(ctx context.Context, userID sharedtypes.DiscordID, duration time.Duration) {
}

func (n *NoOpMetrics) RecordTagAvailabilityCheck(ctx context.Context, available bool, tag sharedtypes.TagNumber) {
}

func (n *NoOpMetrics) RecordPermissionCheckAttempt(ctx context.Context, role sharedtypes.UserRoleEnum, action string, resource string) {
}

func (n *NoOpMetrics) RecordPermissionCheckSuccess(ctx context.Context, role sharedtypes.UserRoleEnum, action string, resource string) {
}

func (n *NoOpMetrics) RecordPermissionCheckFailure(ctx context.Context, role sharedtypes.UserRoleEnum, action string, resource string) {
}

func (n *NoOpMetrics) RecordPermissionCheckDuration(ctx context.Context, role sharedtypes.UserRoleEnum, action string, resource string, duration time.Duration) {
}

func (n *NoOpMetrics) RecordRoleUpdateAttempt(ctx context.Context, userID sharedtypes.DiscordID, oldRole, newRole sharedtypes.UserRoleEnum) {
}

func (n *NoOpMetrics) RecordRoleUpdateSuccess(ctx context.Context, userID sharedtypes.DiscordID, oldRole, newRole sharedtypes.UserRoleEnum) {
}

func (n *NoOpMetrics) RecordRoleUpdateFailure(ctx context.Context, userID sharedtypes.DiscordID, oldRole, newRole sharedtypes.UserRoleEnum) {
}
func (n *NoOpMetrics) RecordHandlerAttempt(ctx context.Context, handlerName string) {}
func (n *NoOpMetrics) RecordHandlerSuccess(ctx context.Context, handlerName string) {}
func (n *NoOpMetrics) RecordHandlerFailure(ctx context.Context, handlerName string) {}
func (n *NoOpMetrics) RecordHandlerDuration(ctx context.Context, handlerName string, duration time.Duration) {
}

func (n *NoOpMetrics) RecordRoleUpdateDuration(ctx context.Context, userID sharedtypes.DiscordID, oldRole, newRole sharedtypes.UserRoleEnum, duration time.Duration) {
}
func (n *NoOpMetrics) RecordDBQueryDuration(ctx context.Context, duration time.Duration)      {}
func (n *NoOpMetrics) RecordUserCreationByTag(ctx context.Context, tag sharedtypes.TagNumber) {}
func (n *NoOpMetrics) RecordOperationAttempt(ctx context.Context, operationName string, userID sharedtypes.DiscordID) {
}

func (n *NoOpMetrics) RecordOperationSuccess(ctx context.Context, operationName string, userID sharedtypes.DiscordID) {
}

func (n *NoOpMetrics) RecordOperationFailure(ctx context.Context, operationName string, userID sharedtypes.DiscordID) {
}

func (n *NoOpMetrics) RecordOperationDuration(ctx context.Context, operationName string, duration time.Duration, userID sharedtypes.DiscordID) {
}

func (n *NoOpMetrics) RecordUserRole(ctx context.Context, userID sharedtypes.DiscordID, role sharedtypes.UserRoleEnum) {
}
