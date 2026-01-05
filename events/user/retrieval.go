// Package userevents contains user-related domain events.
//
// This file defines the User Retrieval Flow - events for retrieving user data.
//
// # Flow Sequence
//
//  1. Request -> GetUserRequestedV1
//  2. Success -> GetUserResponseV1
//  3. OR Failure -> GetUserFailedV1
//
// # Versioning Strategy
//
// All events include a V1 suffix for future schema evolution.
package userevents

import (
	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
	usertypes "github.com/Black-And-White-Club/frolf-bot-shared/types/user"
)

// =============================================================================
// USER RETRIEVAL FLOW - Event Constants
// =============================================================================

// GetUserRequestedV1 is published when user data is requested.
//
// Pattern: Event Notification
// Subject: user.get.requested.v1
// Producer: any service needing user data
// Consumers: user-service (retrieval handler)
// Triggers: GetUserResponseV1 OR GetUserFailedV1
// Version: v1 (December 2024)
const GetUserRequestedV1 = "user.get.requested.v1"

// GetUserResponseV1 is published with the user data.
//
// Pattern: Event Notification
// Subject: user.get.response.v1
// Producer: user-service
// Consumers: requesting service
// Version: v1 (December 2024)
const GetUserResponseV1 = "user.get.response.v1"

// GetUserFailedV1 is published when user retrieval fails.
//
// Pattern: Event Notification
// Subject: user.get.failed.v1
// Producer: user-service
// Consumers: requesting service
// Version: v1 (December 2024)
const GetUserFailedV1 = "user.get.failed.v1"

// =============================================================================
// USER RETRIEVAL FLOW - Payload Types
// =============================================================================

// GetUserRequestedPayloadV1 contains user retrieval request data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type GetUserRequestedPayloadV1 struct {
	GuildID sharedtypes.GuildID   `json:"guild_id"`
	UserID  sharedtypes.DiscordID `json:"user_id"`
}

// GetUserResponsePayloadV1 contains user retrieval response data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type GetUserResponsePayloadV1 struct {
	GuildID sharedtypes.GuildID `json:"guild_id"`
	User    *usertypes.UserData `json:"user"`
}

// GetUserFailedPayloadV1 contains user retrieval failure data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type GetUserFailedPayloadV1 struct {
	GuildID sharedtypes.GuildID   `json:"guild_id"`
	UserID  sharedtypes.DiscordID `json:"user_id"`
	Reason  string                `json:"reason"`
}
