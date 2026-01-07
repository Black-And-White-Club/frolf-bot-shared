// Package userevents contains user-related domain events.
//
// This file defines the User Creation Flow - events for creating new users
// and handling signup processes.
//
// # Flow Sequences
//
// ## User Creation Flow
//  1. Creation requested -> UserCreationRequestedV1
//  2. Success -> UserCreatedV1
//  3. OR Failure -> UserCreationFailedV1
//
// ## User Signup Flow
//  1. Signup requested -> UserSignupRequestedV1
//  2. Success -> UserSignupSucceededV1
//  3. OR Failure -> UserSignupFailedV1
//
// # Versioning Strategy
//
// All events include a V1 suffix for future schema evolution.
package userevents

import (
	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
)

// =============================================================================
// USER CREATION FLOW - Event Constants
// =============================================================================

// UserCreationRequestedV1 is published when user creation is requested.
//
// Pattern: Event Notification
// Subject: user.creation.requested.v1
// Producer: discord-service, api-service
// Consumers: user-service (creation handler)
// Triggers: UserCreatedV1 OR UserCreationFailedV1
// Version: v1 (December 2024)
const UserCreationRequestedV1 = "user.creation.requested.v1"

// UserCreatedV1 is published when a user is successfully created.
//
// Pattern: Event Notification
// Subject: user.created.v1
// Producer: user-service
// Consumers: discord-service (welcome message), leaderboard-service (tag assignment)
// Version: v1 (December 2024)
const UserCreatedV1 = "user.created.v1"

// UserCreationFailedV1 is published when user creation fails.
//
// Pattern: Event Notification
// Subject: user.creation.failed.v1
// Producer: user-service
// Consumers: discord-service (error handler), monitoring
// Version: v1 (December 2024)
const UserCreationFailedV1 = "user.creation.failed.v1"

// UserSignupRequestedV1 is published when a user initiates signup.
//
// Pattern: Event Notification
// Subject: user.signup.requested.v1
// Producer: discord-service (signup handler)
// Consumers: user-service (signup processor)
// Triggers: UserSignupSucceededV1 OR UserSignupFailedV1
// Version: v1 (December 2024)
const UserSignupRequestedV1 = "user.signup.requested.v1"

// UserSignupSucceededV1 is published when user signup succeeds.
//
// Pattern: Event Notification
// Subject: user.signup.succeeded.v1
// Producer: user-service
// Consumers: discord-service (success handler)
// Version: v1 (December 2024)
const UserSignupSucceededV1 = "user.signup.succeeded.v1"

// UserSignupFailedV1 is published when user signup fails.
//
// Pattern: Event Notification
// Subject: user.signup.failed.v1
// Producer: user-service
// Consumers: discord-service (error handler)
// Version: v1 (December 2024)
const UserSignupFailedV1 = "user.signup.failed.v1"

// =============================================================================
// USER CREATION FLOW - Payload Types
// =============================================================================

// UserCreationRequestedPayloadV1 contains user creation request data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type UserCreationRequestedPayloadV1 struct {
	GuildID       sharedtypes.GuildID    `json:"guild_id"`
	UserID        sharedtypes.DiscordID  `json:"user_id"`
	TagNumber     *sharedtypes.TagNumber `json:"tag_number,omitempty"`
	UDiscUsername *string                `json:"udisc_username,omitempty"`
	UDiscName     *string                `json:"udisc_name,omitempty"`
}

// UserCreatedPayloadV1 contains created user data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
//   - v1.1 (January 2026): Added IsReturningUser flag to distinguish new vs returning users
type UserCreatedPayloadV1 struct {
	GuildID         sharedtypes.GuildID    `json:"guild_id"`
	UserID          sharedtypes.DiscordID  `json:"user_id"`
	TagNumber       *sharedtypes.TagNumber `json:"tag_number,omitempty"`
	IsReturningUser bool                   `json:"is_returning_user"`
}

// UserCreationFailedPayloadV1 contains user creation failure data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type UserCreationFailedPayloadV1 struct {
	GuildID   sharedtypes.GuildID    `json:"guild_id"`
	UserID    sharedtypes.DiscordID  `json:"user_id"`
	Reason    string                 `json:"reason"`
	TagNumber *sharedtypes.TagNumber `json:"tag_number,omitempty"`
}

// UserSignupRequestedPayloadV1 contains user signup request data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type UserSignupRequestedPayloadV1 struct {
	GuildID       sharedtypes.GuildID    `json:"guild_id"`
	UserID        sharedtypes.DiscordID  `json:"user_id"`
	TagNumber     *sharedtypes.TagNumber `json:"tag_number,omitempty"`
	UDiscUsername *string                `json:"udisc_username,omitempty"`
	UDiscName     *string                `json:"udisc_name,omitempty"`
}

// UserSignupSucceededPayloadV1 contains signup success data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type UserSignupSucceededPayloadV1 struct {
	GuildID   sharedtypes.GuildID    `json:"guild_id"`
	UserID    sharedtypes.DiscordID  `json:"user_id"`
	TagNumber *sharedtypes.TagNumber `json:"tag_number,omitempty"`
}

// UserSignupFailedPayloadV1 contains signup failure data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type UserSignupFailedPayloadV1 struct {
	GuildID sharedtypes.GuildID `json:"guild_id"`
	Reason  string              `json:"reason"`
}
