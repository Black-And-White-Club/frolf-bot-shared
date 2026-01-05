// Package roundevents contains all round-related domain events.
//
// This file defines the Round Creation Flow - the complete lifecycle of creating
// a new round from initial user request through validation, parsing, storage,
// and completion.
//
// # Flow Sequence
//
//  1. User initiates creation (Discord slash command) -> RoundCreationRequestedV1
//  2. Backend validates input -> RoundValidationPassedV1 OR RoundValidationFailedV1
//  3. Backend parses date/time -> RoundDateTimeParsedV1
//  4. Backend creates entity -> RoundEntityCreatedV1
//  5. Backend stores to database -> RoundStoredV1
//  6. Backend schedules reminders -> RoundScheduledV1
//  7. Success notification -> RoundCreatedV1
//  8. Failure notification -> RoundCreationFailedV1
//
// # Pattern Reference
//
// This flow follows the Event Notification pattern (Martin Fowler) where each step
// publishes an event to notify downstream consumers of state changes. Events use
// the Snapshot Message style, containing all data needed for consumers to process
// without additional lookups.
//
// # Versioning Strategy
//
// All events include a V1 suffix in the constant name and .v1 suffix in the topic
// string. This allows for future schema evolution while maintaining backward
// compatibility. When breaking changes are needed, a V2 version can be introduced
// with consumers explicitly opting in.
package roundevents

import (
	sharedevents "github.com/Black-And-White-Club/frolf-bot-shared/events/shared"
	roundtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/round"
	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
)

// =============================================================================
// ROUND CREATION FLOW - Event Constants
// =============================================================================

// -----------------------------------------------------------------------------
// Step 1: Creation Requested
// -----------------------------------------------------------------------------

// RoundCreationRequestedV1 is published when a user initiates round creation via Discord.
//
// Pattern: Event Notification
// Subject: round.creation.requested.v1
// Producer: discord-service
// Consumers: backend-service (validation handler)
// Triggers: RoundValidationPassedV1 OR RoundValidationFailedV1
// Version: v1 (December 2024)
// Breaking Changes: None (initial version)
const RoundCreationRequestedV1 = "round.creation.requested.v1"

// -----------------------------------------------------------------------------
// Step 2: Validation Results
// -----------------------------------------------------------------------------

// RoundValidationPassedV1 is published when round creation input passes validation.
//
// Pattern: Event Notification
// Subject: round.validation.passed.v1
// Producer: backend-service (validation handler)
// Consumers: backend-service (datetime parser handler)
// Triggers: RoundDateTimeParsedV1
// Version: v1 (December 2024)
// Breaking Changes: None (initial version)
const RoundValidationPassedV1 = "round.validation.passed.v1"

// RoundValidationFailedV1 is published when round creation input fails validation.
//
// Pattern: Event Notification
// Subject: round.validation.failed.v1
// Producer: backend-service (validation handler)
// Consumers: discord-service (error notification handler)
// Triggers: Discord error message to user
// Version: v1 (December 2024)
// Breaking Changes: None (initial version)
const RoundValidationFailedV1 = "round.validation.failed.v1"

// -----------------------------------------------------------------------------
// Step 3: DateTime Parsing
// -----------------------------------------------------------------------------

// RoundDateTimeParsedV1 is published when the natural language datetime is parsed.
//
// Pattern: Event Notification
// Subject: round.datetime.parsed.v1
// Producer: backend-service (datetime parser handler)
// Consumers: backend-service (entity creation handler)
// Triggers: RoundEntityCreatedV1
// Version: v1 (December 2024)
// Breaking Changes: None (initial version)
const RoundDateTimeParsedV1 = "round.datetime.parsed.v1"

// -----------------------------------------------------------------------------
// Step 4: Entity Creation
// -----------------------------------------------------------------------------

// RoundEntityCreatedV1 is published when the round entity is created in memory.
//
// Pattern: Event Notification
// Subject: round.entity.created.v1
// Producer: backend-service (entity creation handler)
// Consumers: backend-service (storage handler)
// Triggers: RoundStoredV1
// Version: v1 (December 2024)
// Breaking Changes: None (initial version)
const RoundEntityCreatedV1 = "round.entity.created.v1"

// -----------------------------------------------------------------------------
// Step 5: Storage
// -----------------------------------------------------------------------------

// RoundStoredV1 is published when the round is persisted to the database.
//
// Pattern: Event Notification
// Subject: round.stored.v1
// Producer: backend-service (storage handler)
// Consumers: backend-service (scheduling handler)
// Triggers: RoundScheduledV1
// Version: v1 (December 2024)
// Breaking Changes: None (initial version)
const RoundStoredV1 = "round.stored.v1"

// -----------------------------------------------------------------------------
// Step 6: Scheduling
// -----------------------------------------------------------------------------

// RoundScheduledV1 is published when reminders are scheduled for the round.
//
// Pattern: Event Notification
// Subject: round.scheduled.v1
// Producer: backend-service (scheduling handler)
// Consumers: backend-service (completion handler), discord-service (embed creator)
// Triggers: RoundCreatedV1
// Version: v1 (December 2024)
// Breaking Changes: None (initial version)
const RoundScheduledV1 = "round.scheduled.v1"

// -----------------------------------------------------------------------------
// Step 7: Creation Complete
// -----------------------------------------------------------------------------

// RoundCreatedV1 is published when round creation completes successfully.
//
// Pattern: Event Notification
// Subject: round.created.v1
// Producer: backend-service (completion handler)
// Consumers: discord-service (success notification handler)
// Triggers: Discord success embed displayed to user
// Version: v1 (December 2024)
// Breaking Changes: None (initial version)
const RoundCreatedV1 = "round.created.v1"

// -----------------------------------------------------------------------------
// Step 8: Creation Failed
// -----------------------------------------------------------------------------

// RoundCreationFailedV1 is published when round creation fails at any step.
//
// Pattern: Event Notification
// Subject: round.creation.failed.v1
// Producer: backend-service (any handler in the flow)
// Consumers: discord-service (error notification handler)
// Triggers: Discord error message to user
// Version: v1 (December 2024)
// Breaking Changes: None (initial version)
const RoundCreationFailedV1 = "round.creation.failed.v1"

// -----------------------------------------------------------------------------
// Supporting Events
// -----------------------------------------------------------------------------

// RoundErrorV1 is published for generic round operation errors.
//
// Pattern: Event Notification
// Subject: round.error.v1
// Producer: backend-service (various handlers)
// Consumers: discord-service (error handler), monitoring systems
// Version: v1 (December 2024)
const RoundErrorV1 = "round.error.v1"

// RoundEventMessageIDUpdatedV1 is published when the Discord message ID is updated.
//
// Pattern: Event Notification
// Subject: round.event.message.id.updated.v1
// Producer: discord-service
// Consumers: backend-service (message ID storage handler)
// Version: v1 (December 2024)
const RoundEventMessageIDUpdatedV1 = "round.event.message.id.updated.v1"

// RoundEventMessageIDUpdateV1 is the request to update a message ID.
//
// Pattern: Event Notification
// Subject: round.event.message.id.update.v1
// Producer: backend-service
// Consumers: backend-service (message ID handler)
// Version: v1 (December 2024)
const RoundEventMessageIDUpdateV1 = "round.event.message.id.update.v1"

// RoundTraceEventV1 is published for distributed tracing purposes.
//
// Pattern: Event Notification
// Subject: round.trace.event.v1
// Producer: Any service
// Consumers: Observability systems
// Version: v1 (December 2024)
const RoundTraceEventV1 = "round.trace.event.v1"

// =============================================================================
// ROUND CREATION FLOW - Payload Types
// =============================================================================

// -----------------------------------------------------------------------------
// Step 1: Creation Requested Payload
// -----------------------------------------------------------------------------

// CreateRoundRequestedPayloadV1 contains raw user input before validation.
// This is a snapshot message containing all data needed for validation.
//
// Backward Compatibility Rules:
//   - Never remove fields
//   - Never rename fields
//   - Never change field types
//   - All new fields must be optional pointers with omitempty
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type CreateRoundRequestedPayloadV1 struct {
	// v1.0 fields (required, never change these)
	GuildID     sharedtypes.GuildID    `json:"guild_id"`
	Title       roundtypes.Title       `json:"title"`
	Description roundtypes.Description `json:"description"`
	StartTime   string                 `json:"start_time"` // Unparsed natural language, e.g., "tomorrow 3pm"
	Location    roundtypes.Location    `json:"location"`
	UserID      sharedtypes.DiscordID  `json:"user_id"`
	ChannelID   string                 `json:"channel_id"`
	Timezone    roundtypes.Timezone    `json:"timezone"`

	// Future additions go here, always optional with omitempty
	// Example: MaxParticipants *int `json:"max_participants,omitempty"`
}

// RoundCreateRequestPayloadV1 contains the round creation request with base payload.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoundCreateRequestPayloadV1 struct {
	GuildID sharedtypes.GuildID `json:"guild_id"`
	roundtypes.BaseRoundPayload
	Timezone roundtypes.Timezone `json:"timezone"`
}

// -----------------------------------------------------------------------------
// Step 2: Validation Result Payloads
// -----------------------------------------------------------------------------

// RoundValidationPassedPayloadV1 contains validated round creation data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoundValidationPassedPayloadV1 struct {
	GuildID                     sharedtypes.GuildID           `json:"guild_id"`
	CreateRoundRequestedPayload CreateRoundRequestedPayloadV1 `json:"round_create_request_payload"`
}

// RoundValidationFailedPayloadV1 contains validation failure details.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoundValidationFailedPayloadV1 struct {
	GuildID       sharedtypes.GuildID   `json:"guild_id"`
	UserID        sharedtypes.DiscordID `json:"user_id"`
	ErrorMessages []string              `json:"error_messages"`
}

// -----------------------------------------------------------------------------
// Step 3: DateTime Parsed Payload
// -----------------------------------------------------------------------------

// RoundDateTimeParsedPayloadV1 contains the parsed datetime result.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoundDateTimeParsedPayloadV1 struct {
	GuildID                   sharedtypes.GuildID         `json:"guild_id"`
	RoundCreateRequestPayload RoundCreateRequestPayloadV1 `json:"round_create_request_payload"`
	StartTime                 *sharedtypes.StartTime      `json:"start_time"`
}

// -----------------------------------------------------------------------------
// Step 4: Entity Created Payload
// -----------------------------------------------------------------------------

// RoundEntityCreatedPayloadV1 contains the created round entity.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoundEntityCreatedPayloadV1 struct {
	GuildID          sharedtypes.GuildID               `json:"guild_id"`
	Round            roundtypes.Round                  `json:"round"`
	DiscordChannelID string                            `json:"discord_channel_id"`
	DiscordGuildID   string                            `json:"discord_guild_id"`
	Config           *sharedevents.GuildConfigFragment `json:"config_fragment,omitempty"`
}

// -----------------------------------------------------------------------------
// Step 5: Storage Payload
// -----------------------------------------------------------------------------

// RoundStoredPayloadV1 contains the stored round data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoundStoredPayloadV1 struct {
	GuildID sharedtypes.GuildID               `json:"guild_id"`
	Round   roundtypes.Round                  `json:"round"`
	Config  *sharedevents.GuildConfigFragment `json:"config_fragment,omitempty"`
}

// -----------------------------------------------------------------------------
// Step 6: Scheduling Payload
// -----------------------------------------------------------------------------

// RoundScheduledPayloadV1 contains scheduled round details.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoundScheduledPayloadV1 struct {
	GuildID sharedtypes.GuildID `json:"guild_id"`
	roundtypes.BaseRoundPayload
	EventMessageID string                            `json:"discord_message_id"`
	Config         *sharedevents.GuildConfigFragment `json:"config_fragment,omitempty"`
	ChannelID      string                            `json:"channel_id,omitempty"`
}

// -----------------------------------------------------------------------------
// Step 7: Creation Complete Payload
// -----------------------------------------------------------------------------

// RoundCreatedPayloadV1 contains the successfully created round.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoundCreatedPayloadV1 struct {
	GuildID sharedtypes.GuildID `json:"guild_id"`
	roundtypes.BaseRoundPayload
	ChannelID string                            `json:"channel_id"`
	Config    *sharedevents.GuildConfigFragment `json:"config_fragment,omitempty"`
}

// -----------------------------------------------------------------------------
// Step 8: Creation Failed Payload
// -----------------------------------------------------------------------------

// RoundCreationFailedPayloadV1 contains creation failure details.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoundCreationFailedPayloadV1 struct {
	GuildID      sharedtypes.GuildID   `json:"guild_id"`
	UserID       sharedtypes.DiscordID `json:"user_id"`
	ErrorMessage string                `json:"error_message"`
	ChannelID    string                `json:"channel_id"`
}

// -----------------------------------------------------------------------------
// Supporting Payloads
// -----------------------------------------------------------------------------

// RoundEventCreatedPayloadV1 contains the round event creation confirmation.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoundEventCreatedPayloadV1 struct {
	GuildID        sharedtypes.GuildID `json:"guild_id"`
	RoundID        sharedtypes.RoundID `json:"round_id"`
	EventMessageID string              `json:"discord_message_id"`
}

// RoundMessageIDUpdatePayloadV1 contains the message ID update request.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoundMessageIDUpdatePayloadV1 struct {
	GuildID sharedtypes.GuildID `json:"guild_id"`
	RoundID sharedtypes.RoundID `json:"round_id"`
}

// RoundErrorPayloadV1 contains generic round error details.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoundErrorPayloadV1 struct {
	GuildID sharedtypes.GuildID `json:"guild_id"`
	RoundID sharedtypes.RoundID `json:"round_id"`
	Error   string              `json:"error"`
}
