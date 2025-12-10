package userevents

import (
	"time"

	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
	usertypes "github.com/Black-And-White-Club/frolf-bot-shared/types/user"
)

// Stream names
const (
	UserUDiscMatchConfirmationRequired = "user.udisc.match.confirmation_required"
	UserUDiscMatchConfirmed            = "user.udisc.match.confirmed"
	UserStreamName                     = "user"
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
	UserPermissionsCheckRequest   = "user.permissions.check.request"
	UserPermissionsCheckResponse  = "user.permissions.check.response"
	UserPermissionsCheckFailed    = "user.permissions.check.failed"
	TagAvailable                  = "user.tag.available"
	TagUnavailable                = "user.tag.unavailable"
	TagAvailabilityCheckRequested = "leaderboard.tag.availability.check.requested"
	TagAssignmentRequested        = "user.tag.assignment.requested"
	TagAssignedForUserCreation    = "user.tag.assigned.for.creation"
	TagAssignmentFailed           = "user.tag.assignment.failed"
	// Discord-specific topics for user role updates
	DiscordUserRoleUpdateRequest = "discord.user.role.update.request"
	DiscordUserRoleUpdated       = "discord.user.role.updated"
	DiscordUserRoleUpdateFailed  = "discord.user.role.update.failed"

	// User retrieval flow
	GetUserRequest  = "user.get.request"
	GetUserResponse = "user.get.response"
	GetUserFailed   = "user.get.failed"

	// User role retrieval flow
	GetUserRoleRequest  = "user.role.get.request"
	GetUserRoleResponse = "discord.user.role.get.response"
	GetUserRoleFailed   = "discord.user.role.get.failed"

	// UDisc identity management
	UpdateUDiscIdentityRequest = "user.udisc.identity.update.request"
	UDiscIdentityUpdated       = "user.udisc.identity.updated"
	UDiscIdentityUpdateFailed  = "user.udisc.identity.update.failed"
)

// BaseEventPayload is a struct that can be embedded in other event structs to provide common fields.
type BaseEventPayload struct {
	EventID   string    `json:"event_id"`
	Timestamp time.Time `json:"timestamp"`
}

// Payload types
type CreateUserRequestedPayload struct {
	GuildID       sharedtypes.GuildID    `json:"guild_id"`
	UserID        sharedtypes.DiscordID  `json:"user_id"`
	TagNumber     *sharedtypes.TagNumber `json:"tag_number"`
	UDiscUsername *string                `json:"udisc_username,omitempty"`
	UDiscName     *string                `json:"udisc_name,omitempty"`
}

type UserSignupRequestPayload struct {
	GuildID       sharedtypes.GuildID    `json:"guild_id"`
	UserID        sharedtypes.DiscordID  `json:"user_id"`
	TagNumber     *sharedtypes.TagNumber `json:"tag_number,omitempty"`
	UDiscUsername *string                `json:"udisc_username,omitempty"`
	UDiscName     *string                `json:"udisc_name,omitempty"`
}

type UserSignupFailedPayload struct {
	GuildID sharedtypes.GuildID `json:"guild_id"`
	Reason  string              `json:"reason"`
}

type UserCreatedPayload struct {
	GuildID   sharedtypes.GuildID    `json:"guild_id"`
	UserID    sharedtypes.DiscordID  `json:"user_id"`
	TagNumber *sharedtypes.TagNumber `json:"tag_number,omitempty"`
}

type UserCreationFailedPayload struct {
	GuildID   sharedtypes.GuildID    `json:"guild_id"`
	UserID    sharedtypes.DiscordID  `json:"user_id"`
	Reason    string                 `json:"reason"`
	TagNumber *sharedtypes.TagNumber `json:"tag_number,omitempty"`
}

type UserRoleUpdateRequestPayload struct {
	GuildID     sharedtypes.GuildID      `json:"guild_id"`
	RequesterID sharedtypes.DiscordID    `json:"requester_id"`
	UserID      sharedtypes.DiscordID    `json:"user_id"`
	Role        sharedtypes.UserRoleEnum `json:"role"`
}

type UpdateUserRoleRequestedPayload struct {
	GuildID     sharedtypes.GuildID      `json:"guild_id"`
	UserID      sharedtypes.DiscordID    `json:"user_id"`
	Role        sharedtypes.UserRoleEnum `json:"role"`
	RequesterID sharedtypes.DiscordID    `json:"requester_id"`
}

// UserRoleUpdateResultPayload is the result of the role update operation (from backend).
type UserRoleUpdateResultPayload struct {
	GuildID sharedtypes.GuildID      `json:"guild_id"`
	UserID  sharedtypes.DiscordID    `json:"user_id"`
	Role    sharedtypes.UserRoleEnum `json:"role"`
	Success bool                     `json:"success"`
	Error   string                   `json:"error,omitempty"`
}

// DEPRECATED, going to use UserRoleUpdateResultPayload
type UserRoleUpdatedPayload struct {
	GuildID sharedtypes.GuildID      `json:"guild_id"`
	UserID  sharedtypes.DiscordID    `json:"user_id"`
	Role    sharedtypes.UserRoleEnum `json:"role"`
}

// DEPRECATED, going to use UserRoleUpdateResultPayload
type UserRoleUpdateFailedPayload struct {
	GuildID sharedtypes.GuildID      `json:"guild_id"`
	UserID  sharedtypes.DiscordID    `json:"user_id"`
	Role    sharedtypes.UserRoleEnum `json:"role"`
	Reason  string                   `json:"reason"`
}

type GetUserRoleRequestPayload struct {
	GuildID sharedtypes.GuildID   `json:"guild_id"`
	UserID  sharedtypes.DiscordID `json:"user_id"`
}

type GetUserRoleResponsePayload struct {
	GuildID sharedtypes.GuildID      `json:"guild_id"`
	UserID  sharedtypes.DiscordID    `json:"user_id"`
	Role    sharedtypes.UserRoleEnum `json:"role"`
}

type GetUserRoleFailedPayload struct {
	GuildID sharedtypes.GuildID   `json:"guild_id"`
	UserID  sharedtypes.DiscordID `json:"user_id"`
	Reason  string                `json:"reason"`
}

type GetUserRequestPayload struct {
	GuildID sharedtypes.GuildID   `json:"guild_id"`
	UserID  sharedtypes.DiscordID `json:"user_id"`
}

type GetUserResponsePayload struct {
	GuildID sharedtypes.GuildID `json:"guild_id"`
	User    *usertypes.UserData `json:"user"`
}

type GetUserFailedPayload struct {
	GuildID sharedtypes.GuildID   `json:"guild_id"`
	UserID  sharedtypes.DiscordID `json:"user_id"`
	Reason  string                `json:"reason"`
}

type UpdateUDiscIdentityRequestPayload struct {
	GuildID  sharedtypes.GuildID   `json:"guild_id"`
	UserID   sharedtypes.DiscordID `json:"user_id"`
	Username *string               `json:"username,omitempty"`
	Name     *string               `json:"name,omitempty"`
}

type UDiscIdentityUpdatedPayload struct {
	GuildID  sharedtypes.GuildID   `json:"guild_id"`
	UserID   sharedtypes.DiscordID `json:"user_id"`
	Username *string               `json:"username,omitempty"`
	Name     *string               `json:"name,omitempty"`
}

type UDiscIdentityUpdateFailedPayload struct {
	GuildID sharedtypes.GuildID   `json:"guild_id"`
	UserID  sharedtypes.DiscordID `json:"user_id"`
	Reason  string                `json:"reason"`
}

type UserPermissionsCheckRequestPayload struct {
	GuildID     sharedtypes.GuildID      `json:"guild_id"`
	UserID      sharedtypes.DiscordID    `json:"user_id"`
	Role        sharedtypes.UserRoleEnum `json:"role"`
	RequesterID sharedtypes.DiscordID    `json:"requester_id"`
}

type UserPermissionsCheckResponsePayload struct {
	GuildID       sharedtypes.GuildID      `json:"guild_id"`
	HasPermission bool                     `json:"has_permission"`
	UserID        sharedtypes.DiscordID    `json:"user_id"`
	Role          sharedtypes.UserRoleEnum `json:"role"`
	RequesterID   string                   `json:"requester_id"`
}

type UserPermissionsCheckFailedPayload struct {
	GuildID     sharedtypes.GuildID      `json:"guild_id"`
	Reason      string                   `json:"reason"`
	UserID      sharedtypes.DiscordID    `json:"user_id"`
	Role        sharedtypes.UserRoleEnum `json:"role"`
	RequesterID sharedtypes.DiscordID    `json:"requester_id"`
}

// UDiscMatchConfirmationRequiredPayload is published when player matches require admin confirmation
type UDiscMatchConfirmationRequiredPayload struct {
	ImportID         string                `json:"import_id"`
	GuildID          sharedtypes.GuildID   `json:"guild_id"`
	RoundID          sharedtypes.RoundID   `json:"round_id"`
	UserID           sharedtypes.DiscordID `json:"user_id"`
	ChannelID        string                `json:"channel_id"`
	UnmatchedPlayers []string              `json:"unmatched_players"`
	Timestamp        time.Time             `json:"timestamp"`
}

// UDiscMatchConfirmedPayload is published when an admin confirms player matches
type UDiscMatchConfirmedPayload struct {
	ImportID  string                  `json:"import_id"`
	GuildID   sharedtypes.GuildID     `json:"guild_id"`
	RoundID   sharedtypes.RoundID     `json:"round_id"`
	UserID    sharedtypes.DiscordID   `json:"user_id"`
	ChannelID string                  `json:"channel_id"`
	Timestamp time.Time               `json:"timestamp"`
	Mappings  []UDiscConfirmedMapping `json:"mappings"`
}

// UDiscConfirmedMapping represents a resolved player match
type UDiscConfirmedMapping struct {
	PlayerName    string                `json:"player_name"`
	DiscordUserID sharedtypes.DiscordID `json:"discord_user_id"`
}

// TagAvailabilityCheckRequestedPayload is the payload for the TagAvailabilityCheckRequested event.
type TagAvailabilityCheckRequestedPayload struct {
	GuildID   sharedtypes.GuildID   `json:"guild_id"`
	TagNumber sharedtypes.TagNumber `json:"tag_number"`
	UserID    sharedtypes.DiscordID `json:"user_id"`
}

// TagAvailablePayload is the payload for the TagAvailable event.
type TagAvailablePayload struct {
	GuildID   sharedtypes.GuildID   `json:"guild_id"`
	UserID    sharedtypes.DiscordID `json:"user_id"`
	TagNumber sharedtypes.TagNumber `json:"tag_number"`
}

// TagUnavailablePayload is the payload for the TagUnavailable event.
type TagUnavailablePayload struct {
	GuildID   sharedtypes.GuildID   `json:"guild_id"`
	UserID    sharedtypes.DiscordID `json:"user_id"`
	TagNumber sharedtypes.TagNumber `json:"tag_number"`
	Reason    string                `json:"reason"`
}

type TagAssignmentRequestedPayload struct {
	GuildID   sharedtypes.GuildID   `json:"guild_id"`
	UserID    sharedtypes.DiscordID `json:"user_id"`
	TagNumber sharedtypes.TagNumber `json:"tag_number"`
}

type TagAssignedForUserCreationPayload struct {
	GuildID   sharedtypes.GuildID    `json:"guild_id"`
	UserID    sharedtypes.DiscordID  `json:"user_id"`
	TagNumber *sharedtypes.TagNumber `json:"tag_number"`
}
