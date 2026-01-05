// Package userevents contains user-related domain events.
//
// This file defines the User Role Flow - events for updating and retrieving
// user roles and permissions.
//
// # Flow Sequences
//
// ## Role Update Flow
//  1. Update requested -> UserRoleUpdateRequestedV1
//  2. Success -> UserRoleUpdatedV1
//  3. OR Failure -> UserRoleUpdateFailedV1
//
// ## Role Retrieval Flow
//  1. Request -> GetUserRoleRequestedV1
//  2. Success -> GetUserRoleResponseV1
//  3. OR Failure -> GetUserRoleFailedV1
//
// ## Permissions Check Flow
//  1. Request -> UserPermissionsCheckRequestedV1
//  2. Response -> UserPermissionsCheckResponseV1
//  3. OR Failure -> UserPermissionsCheckFailedV1
//
// # Versioning Strategy
//
// All events include a V1 suffix for future schema evolution.
package userevents

import (
	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
)

// =============================================================================
// USER ROLE UPDATE FLOW - Event Constants
// =============================================================================

// UserRoleUpdateRequestedV1 is published when a role update is requested.
//
// Pattern: Event Notification
// Subject: user.role.update.requested.v1
// Producer: discord-service (admin command)
// Consumers: user-service (role handler)
// Triggers: UserRoleUpdatedV1 OR UserRoleUpdateFailedV1
// Version: v1 (December 2024)
const UserRoleUpdateRequestedV1 = "user.role.update.requested.v1"

// UserRoleUpdatedV1 is published when a user's role is successfully updated.
//
// Pattern: Event Notification
// Subject: user.role.updated.v1
// Producer: user-service
// Consumers: discord-service (confirmation)
// Version: v1 (December 2024)
const UserRoleUpdatedV1 = "user.role.updated.v1"

// UserRoleUpdateFailedV1 is published when a role update fails.
//
// Pattern: Event Notification
// Subject: user.role.update.failed.v1
// Producer: user-service
// Consumers: discord-service (error handler)
// Version: v1 (December 2024)
const UserRoleUpdateFailedV1 = "user.role.update.failed.v1"

// =============================================================================
// USER ROLE RETRIEVAL FLOW - Event Constants
// =============================================================================

// GetUserRoleRequestedV1 is published when a user's role is requested.
//
// Pattern: Event Notification
// Subject: user.role.get.requested.v1
// Producer: any service needing role info
// Consumers: user-service (role retrieval handler)
// Triggers: GetUserRoleResponseV1 OR GetUserRoleFailedV1
// Version: v1 (December 2024)
const GetUserRoleRequestedV1 = "user.role.get.requested.v1"

// GetUserRoleResponseV1 is published with the user's role.
//
// Pattern: Event Notification
// Subject: user.role.get.response.v1
// Producer: user-service
// Consumers: requesting service
// Version: v1 (December 2024)
const GetUserRoleResponseV1 = "user.role.get.response.v1"

// GetUserRoleFailedV1 is published when role retrieval fails.
//
// Pattern: Event Notification
// Subject: user.role.get.failed.v1
// Producer: user-service
// Consumers: requesting service
// Version: v1 (December 2024)
const GetUserRoleFailedV1 = "user.role.get.failed.v1"

// =============================================================================
// USER PERMISSIONS CHECK FLOW - Event Constants
// =============================================================================

// UserPermissionsCheckRequestedV1 is published when permissions check is requested.
//
// Pattern: Event Notification
// Subject: user.permissions.check.requested.v1
// Producer: any service needing authorization
// Consumers: user-service (permissions handler)
// Triggers: UserPermissionsCheckResponseV1 OR UserPermissionsCheckFailedV1
// Version: v1 (December 2024)
const UserPermissionsCheckRequestedV1 = "user.permissions.check.requested.v1"

// UserPermissionsCheckResponseV1 is published with permissions check result.
//
// Pattern: Event Notification
// Subject: user.permissions.check.response.v1
// Producer: user-service
// Consumers: requesting service
// Version: v1 (December 2024)
const UserPermissionsCheckResponseV1 = "user.permissions.check.response.v1"

// UserPermissionsCheckFailedV1 is published when permissions check fails.
//
// Pattern: Event Notification
// Subject: user.permissions.check.failed.v1
// Producer: user-service
// Consumers: requesting service
// Version: v1 (December 2024)
const UserPermissionsCheckFailedV1 = "user.permissions.check.failed.v1"

// =============================================================================
// USER ROLE FLOW - Payload Types
// =============================================================================

// -----------------------------------------------------------------------------
// Role Update Payloads
// -----------------------------------------------------------------------------

// UserRoleUpdateRequestedPayloadV1 contains role update request data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type UserRoleUpdateRequestedPayloadV1 struct {
	GuildID     sharedtypes.GuildID      `json:"guild_id"`
	UserID      sharedtypes.DiscordID    `json:"user_id"`
	Role        sharedtypes.UserRoleEnum `json:"role"`
	RequesterID sharedtypes.DiscordID    `json:"requester_id"`
}

// UserRoleUpdatedPayloadV1 contains role update success data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type UserRoleUpdatedPayloadV1 struct {
	GuildID sharedtypes.GuildID      `json:"guild_id"`
	UserID  sharedtypes.DiscordID    `json:"user_id"`
	Role    sharedtypes.UserRoleEnum `json:"role"`
}

// UserRoleUpdateFailedPayloadV1 contains role update failure data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type UserRoleUpdateFailedPayloadV1 struct {
	GuildID sharedtypes.GuildID      `json:"guild_id"`
	UserID  sharedtypes.DiscordID    `json:"user_id"`
	Role    sharedtypes.UserRoleEnum `json:"role"`
	Reason  string                   `json:"reason"`
}

// UserRoleUpdateResultPayloadV1 contains the combined result of a role update.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type UserRoleUpdateResultPayloadV1 struct {
	GuildID sharedtypes.GuildID      `json:"guild_id"`
	UserID  sharedtypes.DiscordID    `json:"user_id"`
	Role    sharedtypes.UserRoleEnum `json:"role"`
	Success bool                     `json:"success"`
	Reason  string                   `json:"reason,omitempty"`
}

// -----------------------------------------------------------------------------
// Role Retrieval Payloads
// -----------------------------------------------------------------------------

// GetUserRoleRequestedPayloadV1 contains role retrieval request data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type GetUserRoleRequestedPayloadV1 struct {
	GuildID sharedtypes.GuildID   `json:"guild_id"`
	UserID  sharedtypes.DiscordID `json:"user_id"`
}

// GetUserRoleResponsePayloadV1 contains role retrieval response data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type GetUserRoleResponsePayloadV1 struct {
	GuildID sharedtypes.GuildID      `json:"guild_id"`
	UserID  sharedtypes.DiscordID    `json:"user_id"`
	Role    sharedtypes.UserRoleEnum `json:"role"`
}

// GetUserRoleFailedPayloadV1 contains role retrieval failure data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type GetUserRoleFailedPayloadV1 struct {
	GuildID sharedtypes.GuildID   `json:"guild_id"`
	UserID  sharedtypes.DiscordID `json:"user_id"`
	Reason  string                `json:"reason"`
}

// -----------------------------------------------------------------------------
// Permissions Check Payloads
// -----------------------------------------------------------------------------

// UserPermissionsCheckRequestedPayloadV1 contains permissions check request data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type UserPermissionsCheckRequestedPayloadV1 struct {
	GuildID     sharedtypes.GuildID      `json:"guild_id"`
	UserID      sharedtypes.DiscordID    `json:"user_id"`
	Role        sharedtypes.UserRoleEnum `json:"role"`
	RequesterID sharedtypes.DiscordID    `json:"requester_id"`
}

// UserPermissionsCheckResponsePayloadV1 contains permissions check result.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type UserPermissionsCheckResponsePayloadV1 struct {
	GuildID       sharedtypes.GuildID      `json:"guild_id"`
	HasPermission bool                     `json:"has_permission"`
	UserID        sharedtypes.DiscordID    `json:"user_id"`
	Role          sharedtypes.UserRoleEnum `json:"role"`
	RequesterID   sharedtypes.DiscordID    `json:"requester_id"`
}

// UserPermissionsCheckFailedPayloadV1 contains permissions check failure data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type UserPermissionsCheckFailedPayloadV1 struct {
	GuildID     sharedtypes.GuildID      `json:"guild_id"`
	UserID      sharedtypes.DiscordID    `json:"user_id"`
	Role        sharedtypes.UserRoleEnum `json:"role"`
	RequesterID sharedtypes.DiscordID    `json:"requester_id"`
	Reason      string                   `json:"reason"`
}
