package usermetrics

import (
	"context"
	"time"

	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
)

// UserMetrics defines metrics specific to user operations using OpenTelemetry
type UserMetrics interface {
	RecordUserCreationAttempt(ctx context.Context, userType string, source string)
	RecordUserCreationSuccess(ctx context.Context, userType string, source string)
	RecordUserCreationFailure(ctx context.Context, userType string, source string)
	RecordUserCreationDuration(ctx context.Context, userType string, source string, duration time.Duration)
	RecordUserRetrievalAttempt(ctx context.Context, userID sharedtypes.DiscordID)
	RecordUserRetrievalSuccess(ctx context.Context, userID sharedtypes.DiscordID)
	RecordUserRetrievalFailure(ctx context.Context, userID sharedtypes.DiscordID)
	RecordUserRetrievalDuration(ctx context.Context, userID sharedtypes.DiscordID, duration time.Duration)
	RecordUserRoleRetrievalAttempt(ctx context.Context, userID sharedtypes.DiscordID)
	RecordUserRoleRetrievalSuccess(ctx context.Context, userID sharedtypes.DiscordID)
	RecordUserRoleRetrievalFailure(ctx context.Context, userID sharedtypes.DiscordID)
	RecordUserRoleRetrievalDuration(ctx context.Context, userID sharedtypes.DiscordID, duration time.Duration)
	RecordTagAvailabilityCheck(ctx context.Context, available bool, tag sharedtypes.TagNumber)
	RecordPermissionCheckAttempt(ctx context.Context, role sharedtypes.UserRoleEnum, action string, resource string)
	RecordPermissionCheckSuccess(ctx context.Context, role sharedtypes.UserRoleEnum, action string, resource string)
	RecordPermissionCheckFailure(ctx context.Context, role sharedtypes.UserRoleEnum, action string, resource string)
	RecordPermissionCheckDuration(ctx context.Context, role sharedtypes.UserRoleEnum, action string, resource string, duration time.Duration)
	RecordRoleUpdateAttempt(ctx context.Context, userID sharedtypes.DiscordID, oldRole, newRole sharedtypes.UserRoleEnum)
	RecordRoleUpdateSuccess(ctx context.Context, userID sharedtypes.DiscordID, oldRole, newRole sharedtypes.UserRoleEnum)
	RecordRoleUpdateFailure(ctx context.Context, userID sharedtypes.DiscordID, oldRole, newRole sharedtypes.UserRoleEnum)
	RecordRoleUpdateDuration(ctx context.Context, userID sharedtypes.DiscordID, oldRole, newRole sharedtypes.UserRoleEnum, duration time.Duration)
	RecordDBQueryDuration(ctx context.Context, duration time.Duration)
	RecordHandlerAttempt(ctx context.Context, handlerName string)
	RecordHandlerSuccess(ctx context.Context, handlerName string)
	RecordHandlerFailure(ctx context.Context, handlerName string)
	RecordHandlerDuration(ctx context.Context, handlerName string, duration time.Duration)
	RecordUserCreationByTag(ctx context.Context, tag sharedtypes.TagNumber)
	RecordOperationAttempt(ctx context.Context, operationName string, userID sharedtypes.DiscordID)
	RecordOperationSuccess(ctx context.Context, operationName string, userID sharedtypes.DiscordID)
	RecordOperationFailure(ctx context.Context, operationName string, userID sharedtypes.DiscordID)
	RecordOperationDuration(ctx context.Context, operationName string, duration time.Duration, userID sharedtypes.DiscordID)
	RecordUserRole(ctx context.Context, userID sharedtypes.DiscordID, role sharedtypes.UserRoleEnum)
}
