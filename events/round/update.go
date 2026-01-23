// Package roundevents contains all round-related domain events.
//
// This file defines the Round Update Flow - events related to
// modifying existing rounds.
//
// # Flow Sequence
//
//  1. User requests update -> RoundUpdateRequestedV1
//  2. Backend validates update -> RoundUpdateValidatedV1
//  3. Backend fetches round -> RoundFetchedV1
//  4. Backend updates entity -> RoundEntityUpdatedV1
//  5. Schedule updated if needed -> RoundScheduleUpdatedV1
//  6. Success notification -> RoundUpdatedV1 OR RoundUpdateSuccessV1
//  7. OR Update fails -> RoundUpdateErrorV1
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
// ROUND UPDATE FLOW - Event Constants
// =============================================================================

// RoundUpdateRequestedV1 is published when a user requests to update a round.
//
// Pattern: Event Notification
// Subject: round.update.requested.v1
// Producer: discord-service
// Consumers: backend-service (validation handler)
// Triggers: RoundUpdateValidatedV1 OR RoundUpdateErrorV1
// Version: v1 (December 2024)
const RoundUpdateRequestedV1 = "round.update.requested.v1"

// RoundUpdateValidatedV1 is published when an update request passes validation.
//
// Pattern: Event Notification
// Subject: round.update.validated.v1
// Producer: backend-service (validation handler)
// Consumers: backend-service (fetch handler)
// Triggers: RoundFetchedV1
// Version: v1 (December 2024)
const RoundUpdateValidatedV1 = "round.update.validated.v1"

// RoundFetchedV1 is published when a round is fetched for update.
//
// Pattern: Event Notification
// Subject: round.fetched.v1
// Producer: backend-service (fetch handler)
// Consumers: backend-service (entity update handler)
// Triggers: RoundEntityUpdatedV1
// Version: v1 (December 2024)
const RoundFetchedV1 = "round.fetched.v1"

// RoundEntityUpdatedV1 is published when the round entity is updated.
//
// Pattern: Event Notification
// Subject: round.entity.updated.v1
// Producer: backend-service (entity update handler)
// Consumers: backend-service (storage handler)
// Triggers: RoundUpdatedV1
// Version: v1 (December 2024)
const RoundEntityUpdatedV1 = "round.entity.updated.v1"

// RoundUpdatedV1 is published when a round update completes successfully.
//
// Pattern: Event Notification
// Subject: round.updated.v1
// Producer: backend-service
// Consumers: discord-service (embed update handler)
// Version: v1 (December 2024)
const RoundUpdatedV1 = "round.updated.v1"

// RoundUpdateSuccessV1 is published to confirm successful round update.
//
// Pattern: Event Notification
// Subject: round.update.success.v1
// Producer: backend-service
// Consumers: discord-service
// Version: v1 (December 2024)
const RoundUpdateSuccessV1 = "round.update.success.v1"

// RoundUpdateErrorV1 is published when a round update fails.
//
// Pattern: Event Notification
// Subject: round.update.error.v1
// Producer: backend-service
// Consumers: discord-service (error handler)
// Version: v1 (December 2024)
const RoundUpdateErrorV1 = "round.update.error.v1"

// RoundScheduleUpdatedV1 is published when round schedule is updated.
//
// Pattern: Event Notification
// Subject: round.schedule.updated.v1
// Producer: backend-service (scheduler)
// Consumers: backend-service (reminder scheduler)
// Version: v1 (December 2024)
const RoundScheduleUpdatedV1 = "round.schedule.updated.v1"

// RoundStateUpdatedV1 is published when round state changes.
//
// Pattern: Event Notification
// Subject: round.state.updated.v1
// Producer: backend-service
// Consumers: discord-service, other modules
// Version: v1 (December 2024)
const RoundStateUpdatedV1 = "round.state.updated.v1"

// RoundsUpdatedV1 is published when multiple rounds are updated.
//
// Pattern: Event Notification
// Subject: round.rounds.updated.v1
// Producer: backend-service
// Consumers: discord-service
// Version: v1 (December 2024)
const RoundsUpdatedV1 = "round.rounds.updated.v1"

// RoundUpdateRescheduleV1 is published when round update requires rescheduling.
//
// Pattern: Event Notification
// Subject: round.update.reschedule.v1
// Producer: backend-service
// Consumers: backend-service (scheduler)
// Version: v1 (December 2024)
const RoundUpdateRescheduleV1 = "round.update.reschedule.v1"

// =============================================================================
// ROUND UPDATE FLOW - Payload Types
// =============================================================================

// UpdateRoundRequestedPayloadV1 contains raw update request data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type UpdateRoundRequestedPayloadV1 struct {
	GuildID     sharedtypes.GuildID     `json:"guild_id"`
	RoundID     sharedtypes.RoundID     `json:"round_id"`
	UserID      sharedtypes.DiscordID   `json:"user_id"`
	ChannelID   string                  `json:"channel_id"`
	MessageID   string                  `json:"message_id"`
	Title       *roundtypes.Title       `json:"title,omitempty"`
	Description *roundtypes.Description `json:"description,omitempty"`
	StartTime   *string                 `json:"start_time,omitempty"`
	Timezone    *roundtypes.Timezone    `json:"timezone"`
	Location    *roundtypes.Location    `json:"location,omitempty"`
}

// RoundUpdateRequestPayloadV1 contains the update request details.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoundUpdateRequestPayloadV1 struct {
	GuildID     sharedtypes.GuildID     `json:"guild_id"`
	RoundID     sharedtypes.RoundID     `json:"round_id"`
	Title       *roundtypes.Title       `json:"title,omitempty"`
	Description *roundtypes.Description `json:"description,omitempty"`
	Location    *roundtypes.Location    `json:"location,omitempty"`
	StartTime   *sharedtypes.StartTime  `json:"start_time,omitempty"`
	EventType   *roundtypes.EventType   `json:"event_type,omitempty"`
	UserID      sharedtypes.DiscordID   `json:"user_id"`
}

// RoundUpdateValidatedPayloadV1 contains validated update data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoundUpdateValidatedPayloadV1 struct {
	GuildID                   sharedtypes.GuildID         `json:"guild_id"`
	RoundUpdateRequestPayload RoundUpdateRequestPayloadV1 `json:"round_update_request_payload"`
}

// RoundFetchedPayloadV1 contains fetched round data for update.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoundFetchedPayloadV1 struct {
	GuildID                   sharedtypes.GuildID         `json:"guild_id"`
	Round                     roundtypes.Round            `json:"round"`
	RoundUpdateRequestPayload RoundUpdateRequestPayloadV1 `json:"round_update_request_payload"`
}

// RoundEntityUpdatedPayloadV1 contains the updated entity.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoundEntityUpdatedPayloadV1 struct {
	GuildID sharedtypes.GuildID `json:"guild_id"`
	Round   roundtypes.Round    `json:"round"`
}

// RoundUpdateSuccessPayloadV1 contains update success confirmation.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoundUpdateSuccessPayloadV1 struct {
	GuildID sharedtypes.GuildID `json:"guild_id"`
	RoundID sharedtypes.RoundID `json:"round_id"`
}

// RoundUpdateErrorPayloadV1 contains update error details.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoundUpdateErrorPayloadV1 struct {
	GuildID            sharedtypes.GuildID          `json:"guild_id"`
	RoundUpdateRequest *RoundUpdateRequestPayloadV1 `json:"round_update_request"`
	Error              string                       `json:"error"`
}

// RoundScheduleUpdatePayloadV1 contains schedule update data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoundScheduleUpdatePayloadV1 struct {
	GuildID   sharedtypes.GuildID    `json:"guild_id"`
	RoundID   sharedtypes.RoundID    `json:"round_id"`
	Title     roundtypes.Title       `json:"title"`
	StartTime *sharedtypes.StartTime `json:"start_time"`
	Location  roundtypes.Location    `json:"location"`
}

// RoundUpdateReschedulePayloadV1 contains rescheduling details.
//
// Schema History:
//   - v1.0 (January 2026): Initial version
type RoundUpdateReschedulePayloadV1 struct {
	GuildID   sharedtypes.GuildID    `json:"guild_id"`
	RoundID   sharedtypes.RoundID    `json:"round_id"`
	Title     roundtypes.Title       `json:"title"`
	StartTime *sharedtypes.StartTime `json:"start_time"`
	Location  roundtypes.Location    `json:"location"`
}

// RoundStateUpdatedPayloadV1 contains state change data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoundStateUpdatedPayloadV1 struct {
	GuildID sharedtypes.GuildID   `json:"guild_id"`
	RoundID sharedtypes.RoundID   `json:"round_id"`
	State   roundtypes.RoundState `json:"state"`
}

// RoundsUpdatedPayloadV1 contains bulk round update results.
//
// Schema History:
//   - v1.0 (January 2026): Initial version
type RoundsUpdatedPayloadV1 struct {
	GuildID  sharedtypes.GuildID   `json:"guild_id"`
	RoundIDs []sharedtypes.RoundID `json:"round_ids"`
	Reason   string                `json:"reason,omitempty"`
}
