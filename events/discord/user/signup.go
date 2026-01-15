// Package discorduserevents contains Discord-specific user events.
//
// This file defines the Discord User Signup Flow - events specific to
// user registration through Discord interactions (button clicks, modals, forms).
//
// # Flow Sequence
//
//  1. User clicks signup button -> SignupStartedV1
//  2. Signup modal shown -> SignupFormSubmittedV1
//  3. User submits form -> SignupSubmissionV1
//  4. Success -> SignupSuccessV1 + SignupAddRoleV1
//  5. OR Failure -> SignupFailedV1 / SignupCanceledV1
//
// # Relationship to Domain Events
//
// These Discord events wrap/trigger domain events in events/user/:
//   - SignupSubmissionV1 -> publishes UserSignupRequest (domain)
//   - SignupSuccessV1 <- subscribes to UserCreated (domain)
//
// # Versioning Strategy
//
// All events include a V1 suffix for future schema evolution.
package discorduserevents

import (
	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
)

// =============================================================================
// DISCORD USER SIGNUP FLOW - Event Constants
// =============================================================================

// -----------------------------------------------------------------------------
// Signup Initiation
// -----------------------------------------------------------------------------

// SignupStartedV1 is published when a user initiates signup via Discord.
//
// Pattern: Event Notification
// Subject: discord.user.signup.started.v1
// Producer: discord-service (button handler)
// Consumers: discord-service (signup flow handler)
// Triggers: Signup modal displayed
// Version: v1 (December 2024)
const SignupStartedV1 = "discord.user.signup.started.v1"

// -----------------------------------------------------------------------------
// Signup Form Events
// -----------------------------------------------------------------------------

// SignupFormSubmittedV1 is published when a user submits the signup form.
//
// Pattern: Event Notification
// Subject: discord.user.signup.form.submitted.v1
// Producer: discord-service (modal handler)
// Consumers: discord-service (signup processor)
// Triggers: User creation flow
// Version: v1 (December 2024)
const SignupFormSubmittedV1 = "discord.user.signup.form.submitted.v1"

// SignupSubmissionV1 is published for signup submission processing.
//
// Pattern: Event Notification
// Subject: discord.user.signup.submission.v1
// Producer: discord-service
// Consumers: discord-service (signup handler)
// Version: v1 (December 2024)
const SignupSubmissionV1 = "discord.user.signup.submission.v1"

// -----------------------------------------------------------------------------
// Signup Tag Events
// -----------------------------------------------------------------------------

// SignupTagAskV1 is published to ask the user for their tag number.
//
// Pattern: Event Notification
// Subject: discord.user.signup.tag.ask.v1
// Producer: discord-service (signup handler)
// Consumers: discord-service (tag prompt handler)
// Triggers: Tag number prompt displayed
// Version: v1 (December 2024)
const SignupTagAskV1 = "discord.user.signup.tag.ask.v1"

// SignupTagSkipV1 is published when a user skips tag entry.
//
// Pattern: Event Notification
// Subject: discord.user.signup.tag.skip.v1
// Producer: discord-service (button handler)
// Consumers: discord-service (signup handler)
// Version: v1 (December 2024)
const SignupTagSkipV1 = "discord.user.signup.tag.skip.v1"

// SignupTagIncludeRequestedV1 is published when a user wants to include a tag.
//
// Pattern: Event Notification
// Subject: discord.user.signup.tag.include.requested.v1
// Producer: discord-service
// Consumers: discord-service (tag handler)
// Version: v1 (December 2024)
const SignupTagIncludeRequestedV1 = "discord.user.signup.tag.include.requested.v1"

// SignupTagPromptSentV1 is published when the tag prompt is sent to user.
//
// Pattern: Event Notification
// Subject: discord.user.signup.tag.prompt.sent.v1
// Producer: discord-service (tag handler)
// Consumers: discord-service (response waiter)
// Version: v1 (December 2024)
const SignupTagPromptSentV1 = "discord.user.signup.tag.prompt.sent.v1"

// -----------------------------------------------------------------------------
// Signup Completion Events
// -----------------------------------------------------------------------------

// SignupSuccessV1 is published when signup completes successfully.
//
// Pattern: Event Notification
// Subject: discord.user.signup.success.v1
// Producer: discord-service (after domain UserCreated)
// Consumers: discord-service (success message handler)
// Triggers: Discord success message, role assignment
// Version: v1 (December 2024)
const SignupSuccessV1 = "discord.user.signup.success.v1"

// SignupFailedV1 is published when signup fails.
//
// Pattern: Event Notification
// Subject: discord.user.signup.failed.v1
// Producer: discord-service (error handler)
// Consumers: discord-service (error message handler)
// Triggers: Discord error message
// Version: v1 (December 2024)
const SignupFailedV1 = "discord.user.signup.failed.v1"

// SignupCanceledV1 is published when a user cancels signup.
//
// Pattern: Event Notification
// Subject: discord.user.signup.canceled.v1
// Producer: discord-service (cancel handler)
// Consumers: discord-service (cleanup handler)
// Version: v1 (December 2024)
const SignupCanceledV1 = "discord.user.signup.canceled.v1"

// -----------------------------------------------------------------------------
// Discord Role Events (Signup-Related)
// -----------------------------------------------------------------------------

// SignupAddRoleV1 is published to request adding a Discord role after signup.
//
// Pattern: Event Notification
// Subject: discord.user.signup.role.add.v1
// Producer: discord-service (signup success handler)
// Consumers: discord-service (role manager)
// Triggers: SignupRoleAddedV1 OR SignupRoleAdditionFailedV1
// Version: v1 (December 2024)
const SignupAddRoleV1 = "discord.user.signup.role.add.v1"

// SignupRoleAddedV1 is published when a Discord role is successfully added.
//
// Pattern: Event Notification
// Subject: discord.user.signup.role.added.v1
// Producer: discord-service (role manager)
// Consumers: discord-service (completion handler)
// Version: v1 (December 2024)
const SignupRoleAddedV1 = "discord.user.signup.role.added.v1"

// SignupRoleAdditionFailedV1 is published when Discord role addition fails.
//
// Pattern: Event Notification
// Subject: discord.user.signup.role.addition.failed.v1
// Producer: discord-service (role manager)
// Consumers: discord-service (error handler)
// Version: v1 (December 2024)
const SignupRoleAdditionFailedV1 = "discord.user.signup.role.addition.failed.v1"

// -----------------------------------------------------------------------------
// Trace Events
// -----------------------------------------------------------------------------

// SignupResponseTraceV1 is published for signup response tracing.
//
// Pattern: Event Notification
// Subject: discord.user.signup.response.trace.v1
// Producer: discord-service
// Consumers: Observability systems
// Version: v1 (December 2024)
const SignupResponseTraceV1 = "discord.user.signup.response.trace.v1"

// =============================================================================
// DISCORD USER SIGNUP FLOW - Payload Types
// =============================================================================

// -----------------------------------------------------------------------------
// Signup Initiation Payloads
// -----------------------------------------------------------------------------

// SignupStartedPayloadV1 contains signup initiation data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type SignupStartedPayloadV1 struct {
	UserID    sharedtypes.DiscordID `json:"user_id"`
	ChannelID string                `json:"channel_id"`
	MessageID string                `json:"message_id,omitempty"`
	GuildID   string                `json:"guild_id"`
}

// -----------------------------------------------------------------------------
// Signup Form Payloads
// -----------------------------------------------------------------------------

// SignupFormSubmittedPayloadV1 contains form submission data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type SignupFormSubmittedPayloadV1 struct {
	UserID           sharedtypes.DiscordID  `json:"user_id"`
	InteractionID    string                 `json:"interaction_id"`
	InteractionToken string                 `json:"interaction_token"`
	TagNumber        *sharedtypes.TagNumber `json:"tag_number,omitempty"`
	GuildID          string                 `json:"guild_id"`
}

// -----------------------------------------------------------------------------
// Signup Completion Payloads
// -----------------------------------------------------------------------------

// SignupSuccessPayloadV1 contains signup success data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type SignupSuccessPayloadV1 struct {
	UserID        sharedtypes.DiscordID `json:"user_id"`
	CorrelationID string                `json:"correlation_id"`
	GuildID       string                `json:"guild_id"`
}

// SignupFailedPayloadV1 contains signup failure data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type SignupFailedPayloadV1 struct {
	Reason        string                `json:"reason"`
	Detail        string                `json:"detail"`
	UserID        sharedtypes.DiscordID `json:"user_id"`
	CorrelationID string                `json:"correlation_id"`
	GuildID       string                `json:"guild_id"`
}

// CancelPayloadV1 contains signup cancellation data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type CancelPayloadV1 struct {
	UserID  sharedtypes.DiscordID `json:"user_id"`
	GuildID string                `json:"guild_id"`
}

// -----------------------------------------------------------------------------
// Role Payloads
// -----------------------------------------------------------------------------

// AddRolePayloadV1 contains role addition request data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type AddRolePayloadV1 struct {
	UserID  sharedtypes.DiscordID `json:"user_id"`
	RoleID  string                `json:"role_id"`
	GuildID string                `json:"guild_id"`
}

// RoleAddedPayloadV1 contains role addition confirmation data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoleAddedPayloadV1 struct {
	UserID  sharedtypes.DiscordID `json:"user_id"`
	GuildID string                `json:"guild_id"`
}

// RoleAdditionFailedPayloadV1 contains role addition failure data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoleAdditionFailedPayloadV1 struct {
	UserID  sharedtypes.DiscordID `json:"user_id"`
	Reason  string                `json:"reason"`
	GuildID string                `json:"guild_id"`
}
