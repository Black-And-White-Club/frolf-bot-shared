// Package roundevents contains all round-related domain events.
//
// This file defines the Round Retrieval Flow - events related to
// fetching round data and authorization checks.
//
// # Flow Sequence
//
//  1. Round data requested -> GetRoundRequestedV1
//  2. Round retrieved -> RoundRetrievedV1
//
// # Authorization Flow
//
//  1. User role check requested -> RoundUserRoleCheckRequestedV1
//  2. Role check completed -> RoundUserRoleCheckResultV1
//  3. OR Check fails -> RoundUserRoleCheckErrorV1
//
// # Pattern Reference
//
// This flow follows the Event Notification pattern (Martin Fowler).
//
// # Versioning Strategy
//
// All events include a V1 suffix for future schema evolution.
package roundevents

import (
	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
)

// =============================================================================
// ROUND RETRIEVAL FLOW - Event Constants
// =============================================================================

// GetRoundRequestedV1 is published when round data is requested.
//
// Pattern: Event Notification
// Subject: round.get.requested.v1
// Producer: Various services
// Consumers: backend-service (retrieval handler)
// Triggers: RoundRetrievedV1
// Version: v1 (December 2024)
const GetRoundRequestedV1 = "round.get.requested.v1"

// RoundRetrievedV1 is published when round data is retrieved.
//
// Pattern: Event Notification
// Subject: round.retrieved.v1
// Producer: backend-service (retrieval handler)
// Consumers: Requesting service
// Version: v1 (December 2024)
const RoundRetrievedV1 = "round.retrieved.v1"

// =============================================================================
// AUTHORIZATION FLOW - Event Constants
// =============================================================================

// RoundUserRoleCheckRequestedV1 is published to request user role check.
//
// Pattern: Event Notification
// Subject: round.user.role.check.requested.v1
// Producer: backend-service
// Consumers: backend-service (authorization handler)
// Triggers: RoundUserRoleCheckResultV1 OR RoundUserRoleCheckErrorV1
// Version: v1 (December 2024)
const RoundUserRoleCheckRequestedV1 = "round.user.role.check.requested.v1"

// RoundUserRoleCheckResultV1 is published with role check result.
//
// Pattern: Event Notification
// Subject: round.user.role.check.result.v1
// Producer: backend-service (authorization handler)
// Consumers: backend-service (requesting handler)
// Version: v1 (December 2024)
const RoundUserRoleCheckResultV1 = "round.user.role.check.result.v1"

// RoundUserRoleCheckErrorV1 is published when role check fails.
//
// Pattern: Event Notification
// Subject: round.user.role.check.error.v1
// Producer: backend-service (authorization handler)
// Consumers: discord-service (error handler)
// Version: v1 (December 2024)
const RoundUserRoleCheckErrorV1 = "round.user.role.check.error.v1"

// =============================================================================
// SPECIAL EVENTS - Event Constants
// =============================================================================

// RoundEventsSubjectV1 is the general round events subject.
//
// Subject: round.event.v1
// Version: v1 (December 2024)
const RoundEventsSubjectV1 = "round.event.v1"

// DelayedMessagesSubjectV1 is the delayed messages subject.
//
// Subject: delayed.messages.v1
// Version: v1 (December 2024)
const DelayedMessagesSubjectV1 = "delayed.messages.v1"

// =============================================================================
// ROUND RETRIEVAL FLOW - Payload Types
// =============================================================================

// GetRoundRequestPayloadV1 contains round retrieval request data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type GetRoundRequestPayloadV1 struct {
	GuildID        sharedtypes.GuildID   `json:"guild_id"`
	RoundID        sharedtypes.RoundID   `json:"round_id"`
	EventMessageID string                `json:"event_message_id"`
	UserID         sharedtypes.DiscordID `json:"user_id"`
}

// =============================================================================
// AUTHORIZATION FLOW - Payload Types
// =============================================================================

// UserRoleCheckRequestPayloadV1 contains role check request data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type UserRoleCheckRequestPayloadV1 struct {
	GuildID sharedtypes.GuildID   `json:"guild_id"`
	UserID  sharedtypes.DiscordID `json:"user_id"`
	RoundID sharedtypes.RoundID   `json:"round_id"`
}

// UserRoleCheckResultPayloadV1 contains role check result data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type UserRoleCheckResultPayloadV1 struct {
	GuildID sharedtypes.GuildID   `json:"guild_id"`
	UserID  sharedtypes.DiscordID `json:"user_id"`
	RoundID sharedtypes.RoundID   `json:"round_id"`
	HasRole bool                  `json:"has_role"`
	Error   string                `json:"error,omitempty"`
}

// RoundUserRoleCheckErrorPayloadV1 contains role check error data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoundUserRoleCheckErrorPayloadV1 struct {
	GuildID sharedtypes.GuildID   `json:"guild_id"`
	UserID  sharedtypes.DiscordID `json:"user_id"`
	RoundID sharedtypes.RoundID   `json:"round_id"`
	Error   string                `json:"error"`
}
