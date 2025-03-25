package userevents

import (
	"time"

	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
	usertypes "github.com/Black-And-White-Club/frolf-bot-shared/types/user"
)

// Stream names
const (
	UserStreamName = "user"
	// Note: Changed to use the request/response pattern consistently
	LeaderboardTagAvailabilityCheckRequest = "leaderboard.tag.availability.check.request"
	UserCreatedDLQ                         = "user.created.dlq"            // DLQ for UserCreated
	UserCreationFailedDLQ                  = "user.creation.failed.dlq"    // DLQ for UserCreationFailed
	UserRoleUpdatedDLQ                     = "user.role.updated.dlq"       // DLQ for UserRoleUpdated
	UserRoleUpdateFailedDLQ                = "user.role.update.failed.dlq" // DLQ for UserRoleUpdateFailed
	GetUserRoleResponseDLQ                 = "user.get.role.response.dlq"  // DLQ for GetUserRoleResponse
	GetUserResponseDLQ                     = "user.get.response.dlq"       // DLQ for GetUserResponse
	GetUserRoleFailedDLQ                   = "user.get.role.failed.dlq"    // DLQ for GetUserRoleFailed
	GetUserFailedDLQ                       = "user.get.failed.dlq"         // DLQ for GetUserFailed
)

// Event names
const (
	CreateUserRequested           = "user.create.requested"
	UserSignupRequest             = "user.signup.request"
	UserSignupFailed              = "discord.user.signup.failed"
	UserCreated                   = "user.created"
	UserSignupSuccess             = "discord.user.signup.success"
	UserCreationFailed            = "user.creation.failed"
	UserRoleUpdateRequest         = "user.role.update.request"
	UpdateUserRoleRequested       = "user.role.update.requested"
	UserRoleUpdated               = "user.role.updated"
	UserRoleUpdateFailed          = "user.role.update.failed"
	GetUserRoleRequest            = "user.get.role.request"
	GetUserRoleResponse           = "user.get.role.response"
	GetUserRoleFailed             = "user.get.role.failed"
	GetUserRequest                = "user.get.request"
	GetUserResponse               = "user.get.response"
	GetUserFailed                 = "user.get.failed"
	UserPermissionsCheckRequest   = "user.permissions.check.request"
	UserPermissionsCheckResponse  = "user.permissions.check.response"
	UserPermissionsCheckFailed    = "user.permissions.check.failed"
	TagAvailable                  = "user.tag.available"
	TagUnavailable                = "discord.user.tag.unavailable"
	TagAvailabilityCheckRequested = "user.tag.availability.check.requested"
)

// BaseEventPayload is a struct that can be embedded in other event structs to provide common fields.
// DEPRECATED. Use CommonMetadata now
type BaseEventPayload struct {
	EventID   string    `json:"event_id"`
	Timestamp time.Time `json:"timestamp"`
}

// Payload types

type CreateUserRequestedPayload struct {
	UserID    sharedtypes.DiscordID  `json:"user_id"`
	TagNumber *sharedtypes.TagNumber `json:"tag_number"`
}

type UserSignupRequestPayload struct {
	UserID    sharedtypes.DiscordID  `json:"user_id"`
	TagNumber *sharedtypes.TagNumber `json:"tag_number,omitempty"`
}

type UserSignupFailedPayload struct {
	Reason string `json:"reason"`
}

type UserCreatedPayload struct {
	UserID    sharedtypes.DiscordID  `json:"user_id"`
	TagNumber *sharedtypes.TagNumber `json:"tag_number,omitempty"`
}

type UserCreationFailedPayload struct {
	UserID    sharedtypes.DiscordID  `json:"user_id"`
	Reason    string                 `json:"reason"`
	TagNumber *sharedtypes.TagNumber `json:"tag_number,omitempty"`
}

type UserRoleUpdateRequestPayload struct {
	RequesterID sharedtypes.DiscordID    `json:"requester_id"`
	UserID      sharedtypes.DiscordID    `json:"user_id"`
	Role        sharedtypes.UserRoleEnum `json:"role"`
}

type UpdateUserRoleRequestedPayload struct { //DEPRECATED
	UserID      sharedtypes.DiscordID    `json:"user_id"`
	Role        sharedtypes.UserRoleEnum `json:"role"`
	RequesterID sharedtypes.DiscordID    `json:"requester_id"`
}

// UserRoleUpdateResultPayload is the result of the role update operation (from backend).
type UserRoleUpdateResultPayload struct {
	UserID  sharedtypes.DiscordID    `json:"user_id"`
	Role    sharedtypes.UserRoleEnum `json:"role"`
	Success bool                     `json:"success"`
	Error   string                   `json:"error,omitempty"`
}

// DEPRECATED, going to use UserRoleUpdateResultPayload
type UserRoleUpdatedPayload struct {
	UserID sharedtypes.DiscordID    `json:"user_id"`
	Role   sharedtypes.UserRoleEnum `json:"role"`
}

// DEPRECATED, going to use UserRoleUpdateResultPayload
type UserRoleUpdateFailedPayload struct {
	UserID sharedtypes.DiscordID    `json:"user_id"`
	Role   sharedtypes.UserRoleEnum `json:"role"`
	Reason string                   `json:"reason"`
}

type GetUserRoleRequestPayload struct {
	UserID sharedtypes.DiscordID `json:"user_id"`
}

type GetUserRoleResponsePayload struct {
	UserID sharedtypes.DiscordID    `json:"user_id"`
	Role   sharedtypes.UserRoleEnum `json:"role"`
}

type GetUserRoleFailedPayload struct {
	UserID sharedtypes.DiscordID `json:"user_id"`
	Reason string                `json:"reason"`
}

type GetUserRequestPayload struct {
	UserID sharedtypes.DiscordID `json:"user_id"`
}

type GetUserResponsePayload struct {
	User *usertypes.UserData `json:"user"`
}

type GetUserFailedPayload struct {
	UserID sharedtypes.DiscordID `json:"user_id"`
	Reason string                `json:"reason"`
}

type UserPermissionsCheckRequestPayload struct {
	UserID      sharedtypes.DiscordID    `json:"user_id"`
	Role        sharedtypes.UserRoleEnum `json:"role"`
	RequesterID sharedtypes.DiscordID    `json:"requester_id"`
}

type UserPermissionsCheckResponsePayload struct {
	HasPermission bool                     `json:"has_permission"`
	UserID        sharedtypes.DiscordID    `json:"user_id"`
	Role          sharedtypes.UserRoleEnum `json:"role"`
	RequesterID   string                   `json:"requester_id"`
}

type UserPermissionsCheckFailedPayload struct {
	Reason      string                   `json:"reason"`
	UserID      sharedtypes.DiscordID    `json:"user_id"`
	Role        sharedtypes.UserRoleEnum `json:"role"`
	RequesterID sharedtypes.DiscordID    `json:"requester_id"`
}

// TagAvailabilityCheckRequestedPayload is the payload for the TagAvailabilityCheckRequested event.
type TagAvailabilityCheckRequestedPayload struct {
	TagNumber sharedtypes.TagNumber `json:"tag_number"`
	UserID    sharedtypes.DiscordID `json:"user_id"`
}

// TagAvailablePayload is the payload for the TagAvailable event.
type TagAvailablePayload struct {
	UserID    sharedtypes.DiscordID `json:"user_id"`
	TagNumber sharedtypes.TagNumber `json:"tag_number"`
}

// TagUnavailablePayload is the payload for the TagUnavailable event.
type TagUnavailablePayload struct {
	UserID    sharedtypes.DiscordID `json:"user_id"`
	TagNumber sharedtypes.TagNumber `json:"tag_number"`
	Reason    string                `json:"reason"`
}
