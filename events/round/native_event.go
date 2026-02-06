// Package roundevents contains all round-related domain events.
//
// This file defines the Native Event Flow - events related to creating,
// updating, and resolving Discord Guild Scheduled Events ("native events")
// that mirror round lifecycle states.
//
// # Flow Sequences
//
// ## Native Event Creation Flow
//  1. RoundCreatedV1 triggers discord-service (parallel to embed creation)
//  2. discord-service creates Guild Scheduled Event -> NativeEventCreatedV1
//  3. backend-service stores discord_event_id on the round
//  4. OR creation fails -> NativeEventCreateFailedV1
//
// ## Native Event Update Flow
//  1. RoundUpdatedV1 triggers discord-service
//  2. discord-service updates Guild Scheduled Event -> NativeEventUpdatedV1
//  3. OR update fails -> NativeEventUpdateFailedV1
//
// ## RSVP Resolution Flow (Request-Reply)
//  1. discord-service needs to resolve DiscordEventID -> RoundID
//  2. Publishes NativeEventLookupRequestV1
//  3. backend-service responds with NativeEventLookupResultV1
//
// # Pattern Reference
//
// This flow follows the Event Notification pattern (Martin Fowler). The native
// event creation runs as an independent side effect of RoundCreatedV1, parallel
// to the existing embed creation flow. Failure of native event creation does
// not block the embed flow.
//
// # Versioning Strategy
//
// All events include a V1 suffix in the constant name and .v1 suffix in the topic
// string for future schema evolution while maintaining backward compatibility.
package roundevents

import (
	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
)

// =============================================================================
// NATIVE EVENT FLOW - Event Constants
// =============================================================================

// -----------------------------------------------------------------------------
// Native Event Creation Events
// -----------------------------------------------------------------------------

// NativeEventCreatedV1 is published when the Discord Guild Scheduled Event is
// successfully created for a round.
//
// Pattern: Event Notification
// Subject: round.native.event.created.v1
// Producer: discord-service (native event creation handler)
// Consumers: backend-service (stores discord_event_id on round)
// Version: v1 (February 2026)
const NativeEventCreatedV1 = "round.native.event.created.v1"

// NativeEventCreateFailedV1 is published when creating the Discord Guild
// Scheduled Event fails (e.g., API rate limit, missing MANAGE_EVENTS permission).
//
// Pattern: Event Notification
// Subject: round.native.event.create.failed.v1
// Producer: discord-service (native event creation handler)
// Consumers: monitoring systems
// Version: v1 (February 2026)
const NativeEventCreateFailedV1 = "round.native.event.create.failed.v1"

// -----------------------------------------------------------------------------
// Native Event Update Events
// -----------------------------------------------------------------------------

// NativeEventUpdatedV1 is published when a Guild Scheduled Event is updated
// (e.g., title, time, or location changes from RoundUpdatedV1).
//
// Pattern: Event Notification
// Subject: round.native.event.updated.v1
// Producer: discord-service (native event update handler)
// Consumers: monitoring systems (informational)
// Version: v1 (February 2026)
const NativeEventUpdatedV1 = "round.native.event.updated.v1"

// NativeEventUpdateFailedV1 is published when updating a Guild Scheduled Event fails.
//
// Pattern: Event Notification
// Subject: round.native.event.update.failed.v1
// Producer: discord-service (native event update handler)
// Consumers: monitoring systems
// Version: v1 (February 2026)
const NativeEventUpdateFailedV1 = "round.native.event.update.failed.v1"

// -----------------------------------------------------------------------------
// Native Event Lookup Events (Request-Reply)
// -----------------------------------------------------------------------------

// NativeEventLookupRequestV1 is a request-reply event published by
// discord-service to resolve a DiscordEventID to a RoundID. This is the
// fallback path for post-restart RSVP resolution when the in-memory map
// has been lost.
//
// Pattern: Command/Request
// Subject: round.native.event.lookup.request.v1
// Producer: discord-service (RSVP gateway listener)
// Consumers: backend-service (lookup handler)
// Triggers: NativeEventLookupResultV1
// Version: v1 (February 2026)
const NativeEventLookupRequestV1 = "round.native.event.lookup.request.v1"

// NativeEventLookupResultV1 is the reply to NativeEventLookupRequestV1,
// containing the resolved RoundID (or Found=false if the round no longer exists).
//
// Pattern: Event Notification (Reply)
// Subject: round.native.event.lookup.result.v1
// Producer: backend-service (lookup handler)
// Consumers: discord-service (RSVP gateway listener)
// Version: v1 (February 2026)
const NativeEventLookupResultV1 = "round.native.event.lookup.result.v1"

// =============================================================================
// NATIVE EVENT FLOW - Payload Types
// =============================================================================

// -----------------------------------------------------------------------------
// Native Event Creation Payloads
// -----------------------------------------------------------------------------

// NativeEventCreatedPayloadV1 contains the data published when a Discord
// Guild Scheduled Event is successfully created for a round.
//
// Schema History:
//   - v1.0 (February 2026): Initial version
type NativeEventCreatedPayloadV1 struct {
	GuildID        sharedtypes.GuildID `json:"guild_id"`
	RoundID        sharedtypes.RoundID `json:"round_id"`
	DiscordEventID string              `json:"discord_event_id"`
}

// NativeEventCreateFailedPayloadV1 contains the data published when creating
// a Discord Guild Scheduled Event fails.
//
// Schema History:
//   - v1.0 (February 2026): Initial version
type NativeEventCreateFailedPayloadV1 struct {
	GuildID sharedtypes.GuildID `json:"guild_id"`
	RoundID sharedtypes.RoundID `json:"round_id"`
	Error   string              `json:"error"`
}

// -----------------------------------------------------------------------------
// Native Event Update Payloads
// -----------------------------------------------------------------------------

// NativeEventUpdatedPayloadV1 contains the data published when a Discord
// Guild Scheduled Event is successfully updated.
//
// Schema History:
//   - v1.0 (February 2026): Initial version
type NativeEventUpdatedPayloadV1 struct {
	GuildID        sharedtypes.GuildID `json:"guild_id"`
	RoundID        sharedtypes.RoundID `json:"round_id"`
	DiscordEventID string              `json:"discord_event_id"`
}

// NativeEventUpdateFailedPayloadV1 contains the data published when updating
// a Discord Guild Scheduled Event fails.
//
// Schema History:
//   - v1.0 (February 2026): Initial version
type NativeEventUpdateFailedPayloadV1 struct {
	GuildID        sharedtypes.GuildID `json:"guild_id"`
	RoundID        sharedtypes.RoundID `json:"round_id"`
	DiscordEventID string              `json:"discord_event_id"`
	Error          string              `json:"error"`
}

// -----------------------------------------------------------------------------
// Native Event Lookup Payloads (Request-Reply)
// -----------------------------------------------------------------------------

// NativeEventLookupRequestPayloadV1 contains the data for a request to resolve
// a DiscordEventID to a RoundID. Used as the fallback resolution path when
// the in-memory NativeEventMap is empty (e.g., after bot restart).
//
// Schema History:
//   - v1.0 (February 2026): Initial version
type NativeEventLookupRequestPayloadV1 struct {
	GuildID        sharedtypes.GuildID `json:"guild_id"`
	DiscordEventID string              `json:"discord_event_id"`
}

// NativeEventLookupResultPayloadV1 contains the result of a DiscordEventID
// to RoundID resolution. If Found is false, the round no longer exists
// (it may have been deleted).
//
// Schema History:
//   - v1.0 (February 2026): Initial version
type NativeEventLookupResultPayloadV1 struct {
	GuildID        sharedtypes.GuildID `json:"guild_id"`
	RoundID        sharedtypes.RoundID `json:"round_id"`
	DiscordEventID string              `json:"discord_event_id"`
	Found          bool                `json:"found"`
}
