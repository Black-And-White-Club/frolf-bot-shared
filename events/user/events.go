// Package userevents contains user-related domain events.
//
// MIGRATION NOTICE: This file contains legacy event constants.
// New code should use the versioned events from the flow-based files:
//   - creation.go: UserCreationRequestedV1, UserCreatedV1, UserSignupRequestedV1, etc.
//   - roles.go: UserRoleUpdateRequestedV1, UserRoleUpdatedV1, etc.
//   - retrieval.go: GetUserRequestedV1, GetUserResponseV1, etc.
//   - tags.go: TagAvailabilityCheckRequestedV1, TagAvailableV1, etc.
//   - udisc.go: UpdateUDiscIdentityRequestedV1, UDiscMatchConfirmedV1, etc.
//
// See each file for detailed flow documentation and versioning information.
package userevents

import (
	"time"

	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
	usertypes "github.com/Black-And-White-Club/frolf-bot-shared/types/user"
)

// Stream names and DLQ topics
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
// Deprecated: Use versioned constants from creation.go, roles.go, retrieval.go, tags.go, udisc.go
const (
	// Deprecated: Use UserCreationRequestedV1 from creation.go
	CreateUserRequested = "user.create.requested"
	// Deprecated: Use UserSignupRequestedV1 from creation.go
	UserSignupRequest = "user.signup.request"
	// Deprecated: Use UserSignupFailedV1 from creation.go (Discord-specific: use discord/user/signup.go)
	UserSignupFailed = "discord.user.signup.failed"
	// Deprecated: Use UserCreatedV1 from creation.go
	UserCreated = "user.created"
	// Deprecated: Use UserSignupSucceededV1 from creation.go (Discord-specific: use discord/user/signup.go)
	UserSignupSuccess = "discord.user.signup.success"
	// Deprecated: Use UserCreationFailedV1 from creation.go
	UserCreationFailed = "user.creation.failed"
	// Deprecated: Use UserRoleUpdateRequestedV1 from roles.go
	UserRoleUpdateRequest = "user.role.update.request"
	// Deprecated: Use UserRoleUpdateRequestedV1 from roles.go
	UpdateUserRoleRequested = "user.role.update.requested"
	// Deprecated: Use UserRoleUpdatedV1 from roles.go
	UserRoleUpdated = "user.role.updated"
	// Deprecated: Use UserRoleUpdateFailedV1 from roles.go
	UserRoleUpdateFailed = "user.role.update.failed"
	// Deprecated: Use UserPermissionsCheckRequestedV1 from roles.go
	UserPermissionsCheckRequest = "user.permissions.check.request"
	// Deprecated: Use UserPermissionsCheckResponseV1 from roles.go
	UserPermissionsCheckResponse = "user.permissions.check.response"
	// Deprecated: Use UserPermissionsCheckFailedV1 from roles.go
	UserPermissionsCheckFailed = "user.permissions.check.failed"
	// Deprecated: Use TagAvailableV1 from tags.go
	TagAvailable = "user.tag.available"
	// Deprecated: Use TagUnavailableV1 from tags.go
	TagUnavailable = "user.tag.unavailable"
	// Deprecated: Use TagAvailabilityCheckRequestedV1 from tags.go
	TagAvailabilityCheckRequested = "leaderboard.tag.availability.check.requested"
	// Deprecated: Use TagAssignmentRequestedV1 from tags.go
	TagAssignmentRequested = "user.tag.assignment.requested"
	// Deprecated: Use TagAssignedForUserCreationV1 from tags.go
	TagAssignedForUserCreation = "user.tag.assigned.for.creation"
	// Deprecated: Use TagAssignmentFailedV1 from tags.go
	TagAssignmentFailed = "user.tag.assignment.failed"

	// Discord-specific topics for user role updates
	// Deprecated: Use discord/user/roles.go events
	DiscordUserRoleUpdateRequest = "discord.user.role.update.request"
	// Deprecated: Use discord/user/roles.go events
	DiscordUserRoleUpdated = "discord.user.role.updated"
	// Deprecated: Use discord/user/roles.go events
	DiscordUserRoleUpdateFailed = "discord.user.role.update.failed"

	// User retrieval flow
	// Deprecated: Use GetUserRequestedV1 from retrieval.go
	GetUserRequest = "user.get.request"
	// Deprecated: Use GetUserResponseV1 from retrieval.go
	GetUserResponse = "user.get.response"
	// Deprecated: Use GetUserFailedV1 from retrieval.go
	GetUserFailed = "user.get.failed"

	// User role retrieval flow
	// Deprecated: Use GetUserRoleRequestedV1 from roles.go
	GetUserRoleRequest = "user.role.get.request"
	// Deprecated: Use GetUserRoleResponseV1 from roles.go
	GetUserRoleResponse = "discord.user.role.get.response"
	// Deprecated: Use GetUserRoleFailedV1 from roles.go
	GetUserRoleFailed = "discord.user.role.get.failed"

	// UDisc identity management
	// Deprecated: Use UpdateUDiscIdentityRequestedV1 from udisc.go
	UpdateUDiscIdentityRequest = "user.udisc.identity.update.request"
	// Deprecated: Use UDiscIdentityUpdatedV1 from udisc.go
	UDiscIdentityUpdated = "user.udisc.identity.updated"
	// Deprecated: Use UDiscIdentityUpdateFailedV1 from udisc.go
	UDiscIdentityUpdateFailed = "user.udisc.identity.update.failed"
	// Deprecated: Use UDiscMatchConfirmationRequiredV1 from udisc.go
	UserUDiscMatchConfirmationRequired = "user.udisc.match.confirmation_required"
	// Deprecated: Use UDiscMatchConfirmedV1 from udisc.go
	UserUDiscMatchConfirmed = "user.udisc.match.confirmed"
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
	GuildID          sharedtypes.GuildID   `json:"guild_id"`
	RoundID          sharedtypes.RoundID   `json:"round_id"`
	ImportID         string                `json:"import_id"`
	UserID           sharedtypes.DiscordID `json:"user_id"`
	ChannelID        string                `json:"channel_id"`
	UnmatchedPlayers []string              `json:"unmatched_players"`
	Timestamp        time.Time             `json:"timestamp"`
}

// UDiscMatchConfirmedPayload is published when an admin confirms player matches
type UDiscMatchConfirmedPayload struct {
	GuildID      sharedtypes.GuildID     `json:"guild_id"`
	RoundID      sharedtypes.RoundID     `json:"round_id"`
	ImportID     string                  `json:"import_id"`
	UserID       sharedtypes.DiscordID   `json:"user_id"`
	ChannelID    string                  `json:"channel_id"`
	Timestamp    time.Time               `json:"timestamp"`
	Mappings     []UDiscConfirmedMapping `json:"mappings"`
	ParsedScores interface{}             `json:"parsed_scores,omitempty"` // Contains the parsed scorecard data for round module
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
