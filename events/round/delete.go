// Package roundevents contains all round-related domain events.
//
// This file defines the Round Delete Flow - events related to
// deleting rounds.
//
// # Flow Sequence
//
//  1. User requests deletion -> RoundDeleteRequestedV1
//  2. Backend validates request -> RoundDeleteValidatedV1
//  3. Backend fetches round -> RoundToDeleteFetchedV1
//  4. Backend checks authorization -> RoundDeleteAuthorizedV1 OR RoundDeleteUnauthorizedV1
//  5. Round deleted -> RoundDeletedV1
//  6. OR Delete fails -> RoundDeleteErrorV1
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
	roundtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/round"
	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
)

// =============================================================================
// ROUND DELETE FLOW - Event Constants
// =============================================================================

// RoundDeleteRequestedV1 is published when a user requests to delete a round.
//
// Pattern: Event Notification
// Subject: round.delete.requested.v1
// Producer: discord-service
// Consumers: backend-service (validation handler)
// Triggers: RoundDeleteValidatedV1 OR RoundDeleteErrorV1
// Version: v1 (December 2024)
const RoundDeleteRequestedV1 = "round.delete.requested.v1"

// RoundDeleteValidatedV1 is published when a delete request passes validation.
//
// Pattern: Event Notification
// Subject: round.delete.validated.v1
// Producer: backend-service (validation handler)
// Consumers: backend-service (fetch handler)
// Triggers: RoundToDeleteFetchedV1
// Version: v1 (December 2024)
const RoundDeleteValidatedV1 = "round.delete.validated.v1"

// RoundToDeleteFetchedV1 is published when the round to delete is fetched.
//
// Pattern: Event Notification
// Subject: round.to.delete.fetched.v1
// Producer: backend-service (fetch handler)
// Consumers: backend-service (authorization handler)
// Triggers: RoundDeleteAuthorizedV1 OR RoundDeleteUnauthorizedV1
// Version: v1 (December 2024)
const RoundToDeleteFetchedV1 = "round.to.delete.fetched.v1"

// RoundDeleteAuthorizedV1 is published when deletion is authorized.
//
// Pattern: Event Notification
// Subject: round.delete.authorized.v1
// Producer: backend-service (authorization handler)
// Consumers: backend-service (delete handler)
// Triggers: RoundDeletedV1
// Version: v1 (December 2024)
const RoundDeleteAuthorizedV1 = "round.delete.authorized.v1"

// RoundDeleteUnauthorizedV1 is published when deletion is not authorized.
//
// Pattern: Event Notification
// Subject: round.delete.unauthorized.v1
// Producer: backend-service (authorization handler)
// Consumers: discord-service (error handler)
// Version: v1 (December 2024)
const RoundDeleteUnauthorizedV1 = "round.delete.unauthorized.v1"

// RoundDeletedV1 is published when a round is successfully deleted.
//
// Pattern: Event Notification
// Subject: round.deleted.v1
// Producer: backend-service (delete handler)
// Consumers: discord-service (notification handler)
// Version: v1 (December 2024)
const RoundDeletedV1 = "round.deleted.v1"

// RoundDeleteErrorV1 is published when round deletion fails.
//
// Pattern: Event Notification
// Subject: round.delete.error.v1
// Producer: backend-service
// Consumers: discord-service (error handler)
// Version: v1 (December 2024)
const RoundDeleteErrorV1 = "round.delete.error.v1"

// =============================================================================
// ROUND DELETE FLOW - Payload Types
// =============================================================================

// RoundDeleteRequestPayloadV1 contains the delete request details.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoundDeleteRequestPayloadV1 struct {
	GuildID              sharedtypes.GuildID   `json:"guild_id"`
	RoundID              sharedtypes.RoundID   `json:"round_id" validate:"required"`
	RequestingUserUserID sharedtypes.DiscordID `json:"requesting_user_user_id" validate:"required"`
}

// RoundDeleteRequestedPayloadV1 contains delete request data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoundDeleteRequestedPayloadV1 struct {
	GuildID              sharedtypes.GuildID   `json:"guild_id"`
	RoundID              sharedtypes.RoundID   `json:"round_id" validate:"required"`
	RequestingUserUserID sharedtypes.DiscordID `json:"requesting_user_user_id" validate:"required"`
}

// RoundDeleteValidatedPayloadV1 contains validated delete request data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoundDeleteValidatedPayloadV1 struct {
	GuildID                   sharedtypes.GuildID         `json:"guild_id"`
	RoundDeleteRequestPayload RoundDeleteRequestPayloadV1 `json:"round_delete_request_payload"`
}

// RoundToDeleteFetchedPayloadV1 contains the fetched round for deletion.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoundToDeleteFetchedPayloadV1 struct {
	GuildID                   sharedtypes.GuildID         `json:"guild_id"`
	Round                     roundtypes.Round            `json:"round"`
	RoundDeleteRequestPayload RoundDeleteRequestPayloadV1 `json:"round_delete_request_payload"`
}

// RoundDeleteAuthorizedPayloadV1 contains authorization confirmation.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoundDeleteAuthorizedPayloadV1 struct {
	GuildID sharedtypes.GuildID `json:"guild_id"`
	RoundID sharedtypes.RoundID `json:"round_id"`
}

// RoundDeleteUnauthorizedPayloadV1 contains unauthorized deletion details.
//
// Schema History:
//   - v1.0 (January 2026): Initial version
type RoundDeleteUnauthorizedPayloadV1 struct {
	GuildID              sharedtypes.GuildID   `json:"guild_id"`
	RoundID              sharedtypes.RoundID   `json:"round_id"`
	RequestingUserUserID sharedtypes.DiscordID `json:"requesting_user_user_id"`
	Reason               string                `json:"reason"`
}

// RoundDeletedPayloadV1 contains deletion confirmation.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoundDeletedPayloadV1 struct {
	GuildID        sharedtypes.GuildID `json:"guild_id"`
	RoundID        sharedtypes.RoundID `json:"round_id"`
	EventMessageID string              `json:"discord_message_id"`
}

// RoundDeleteErrorPayloadV1 contains deletion error details.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoundDeleteErrorPayloadV1 struct {
	GuildID            sharedtypes.GuildID          `json:"guild_id"`
	RoundDeleteRequest *RoundDeleteRequestPayloadV1 `json:"round_delete_request"`
	Error              string                       `json:"error"`
}
