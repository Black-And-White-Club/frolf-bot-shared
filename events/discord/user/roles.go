// Package discorduserevents contains Discord-specific user events.
//
// This file defines the Discord User Role Update Flow - events specific to
// updating user roles through Discord button interactions and commands.
//
// # Flow Sequence
//
//  1. Admin issues command -> RoleUpdateCommandV1
//  2. Role options shown -> RoleOptionsRequestedV1
//  3. User/admin clicks button -> RoleUpdateButtonPressV1
//  4. Response received -> RoleResponseV1
//  5. OR Timeout -> RoleUpdateTimeoutV1
//  6. OR Failure -> RoleResponseFailedV1
//
// # Relationship to Domain Events
//
// These Discord events wrap/trigger domain events in events/user/:
//   - RoleUpdateButtonPressV1 -> publishes UserRoleUpdateRequest (domain)
//   - RoleResponseV1 <- subscribes to UserRoleUpdated (domain)
//
// # Versioning Strategy
//
// All events include a V1 suffix for future schema evolution.
package discorduserevents

import (
	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
)

// =============================================================================
// DISCORD USER ROLE UPDATE FLOW - Event Constants
// =============================================================================

// RoleUpdateCommandV1 is published when an admin issues a role update command.
//
// Pattern: Event Notification
// Subject: discord.user.role.update.command.v1
// Producer: discord-service (command handler)
// Consumers: discord-service (role update flow handler)
// Triggers: Role options displayed
// Version: v1 (December 2024)
const RoleUpdateCommandV1 = "discord.user.role.update.command.v1"

// RoleUpdateButtonPressV1 is published when a user presses a role update button.
//
// Pattern: Event Notification
// Subject: discord.user.role.update.button.press.v1
// Producer: discord-service (button handler)
// Consumers: discord-service (role processor)
// Triggers: Domain UserRoleUpdateRequest
// Version: v1 (December 2024)
const RoleUpdateButtonPressV1 = "discord.user.role.update.button.press.v1"

// RoleUpdateTimeoutV1 is published when role update times out.
//
// Pattern: Event Notification
// Subject: discord.user.role.update.timeout.v1
// Producer: discord-service (timeout handler)
// Consumers: discord-service (cleanup handler)
// Version: v1 (December 2024)
const RoleUpdateTimeoutV1 = "discord.user.role.update.timeout.v1"

// RoleOptionsRequestedV1 is published to request role options display.
//
// Pattern: Event Notification
// Subject: discord.user.role.options.requested.v1
// Producer: discord-service
// Consumers: discord-service (options handler)
// Version: v1 (December 2024)
const RoleOptionsRequestedV1 = "discord.user.role.options.requested.v1"

// RoleResponseV1 is published when user responds to role selection.
//
// Pattern: Event Notification
// Subject: discord.user.role.response.v1
// Producer: discord-service (selection handler)
// Consumers: discord-service (role processor)
// Version: v1 (December 2024)
const RoleResponseV1 = "discord.user.role.response.v1"

// RoleResponseFailedV1 is published when role response processing fails.
//
// Pattern: Event Notification
// Subject: discord.user.role.response.failed.v1
// Producer: discord-service
// Consumers: discord-service (error handler)
// Version: v1 (December 2024)
const RoleResponseFailedV1 = "discord.user.role.response.failed.v1"

// RoleUpdateResponseTraceV1 is published for role update response tracing.
//
// Pattern: Event Notification
// Subject: discord.user.role.update.response.trace.v1
// Producer: discord-service
// Consumers: Observability systems
// Version: v1 (December 2024)
const RoleUpdateResponseTraceV1 = "discord.user.role.update.response.trace.v1"

// =============================================================================
// DISCORD USER ROLE UPDATE FLOW - Payload Types
// =============================================================================

// RoleUpdateCommandPayloadV1 contains role update command data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoleUpdateCommandPayloadV1 struct {
	TargetUserID sharedtypes.DiscordID `json:"target_user_id"`
	GuildID      string                `json:"guild_id"`
}

// RoleUpdateButtonPressPayloadV1 contains button press data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoleUpdateButtonPressPayloadV1 struct {
	RequesterID         sharedtypes.DiscordID    `json:"requester_id"`
	TargetUserID        sharedtypes.DiscordID    `json:"target_user_id"`
	SelectedRole        sharedtypes.UserRoleEnum `json:"selected_role"`
	InteractionID       string                   `json:"interaction_id"`
	InteractionToken    string                   `json:"interaction_token"`
	InteractionCustomID string                   `json:"custom_id"`
	GuildID             string                   `json:"guild_id"`
}

// RoleUpdateTimeoutPayloadV1 contains timeout data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoleUpdateTimeoutPayloadV1 struct {
	UserID  string `json:"user_id"`
	GuildID string `json:"guild_id"`
}

// RoleUpdateResponsePayloadV1 contains role selection response data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoleUpdateResponsePayloadV1 struct {
	Response  string                `json:"response"` // The chosen role (or "cancel")
	UserID    sharedtypes.DiscordID `json:"user_id"`
	MessageID string                `json:"message_id"`
	GuildID   string                `json:"guild_id"`
}
