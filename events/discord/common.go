// Package discord contains Discord-specific integration events.
//
// These events are used for Discord bot operations that are specific to the
// Discord platform, such as sending DMs, handling interactions, and managing
// Discord-specific UI elements like embeds and buttons.
//
// This file defines common Discord events shared across all modules.
//
// # Distinction from Domain Events
//
// Domain events (in events/round/, events/user/, etc.) represent business logic
// and can be consumed by any service (Discord bot, PWA, backend, etc.).
//
// Discord events (in events/discord/) are specific to Discord integration:
//   - Modal submissions
//   - Button interactions
//   - Embed updates
//   - DM sending
//   - Discord role management
//
// # Versioning Strategy
//
// All events include a V1 suffix in the constant name and .v1 suffix in the topic
// string for future schema evolution while maintaining backward compatibility.
package discord

// =============================================================================
// COMMON DISCORD EVENTS - Event Constants
// =============================================================================

// -----------------------------------------------------------------------------
// Direct Message Events
// -----------------------------------------------------------------------------

// SendDMV1 is published to request sending a DM to a user.
//
// Pattern: Event Notification
// Subject: discord.send.dm.v1
// Producer: Any service needing to notify a user
// Consumers: discord-service (DM sender)
// Triggers: DMSentV1 OR DMErrorV1
// Version: v1 (December 2024)
const SendDMV1 = "discord.send.dm.v1"

// DMSentV1 is published when a DM is successfully sent.
//
// Pattern: Event Notification
// Subject: discord.dm.sent.v1
// Producer: discord-service (DM sender)
// Consumers: Requesting service, monitoring
// Version: v1 (December 2024)
const DMSentV1 = "discord.dm.sent.v1"

// DMErrorV1 is published when sending a DM fails.
//
// Pattern: Event Notification
// Subject: discord.dm.error.v1
// Producer: discord-service (DM sender)
// Consumers: Requesting service, error handling
// Version: v1 (December 2024)
const DMErrorV1 = "discord.dm.error.v1"

// -----------------------------------------------------------------------------
// Interaction Events
// -----------------------------------------------------------------------------

// InteractionRespondedV1 is published when a Discord interaction is responded to.
//
// Pattern: Event Notification
// Subject: discord.interaction.responded.v1
// Producer: discord-service (interaction handler)
// Consumers: Monitoring, audit logging
// Version: v1 (December 2024)
const InteractionRespondedV1 = "discord.interaction.responded.v1"

// -----------------------------------------------------------------------------
// Trace Events
// -----------------------------------------------------------------------------

// DiscordEventTraceV1 is published for Discord event tracing/observability.
//
// Pattern: Event Notification
// Subject: discord.event.trace.v1
// Producer: discord-service
// Consumers: Observability systems
// Version: v1 (December 2024)
const DiscordEventTraceV1 = "discord.event.trace.v1"

// =============================================================================
// Status Constants
// =============================================================================

const (
	StatusSuccess = "success"
	StatusFail    = "fail"
)

// =============================================================================
// COMMON DISCORD EVENTS - Payload Types
// =============================================================================

// -----------------------------------------------------------------------------
// Direct Message Payloads
// -----------------------------------------------------------------------------

// SendDMPayloadV1 contains the request to send a DM.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type SendDMPayloadV1 struct {
	UserID  string `json:"user_id"`
	Message string `json:"message"`
	GuildID string `json:"guild_id"`
}

// DMSentPayloadV1 confirms a DM was successfully sent.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type DMSentPayloadV1 struct {
	UserID  string `json:"user_id"`
	GuildID string `json:"guild_id"`
}

// DMErrorPayloadV1 contains DM failure details.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type DMErrorPayloadV1 struct {
	UserID      string `json:"user_id"`
	ErrorDetail string `json:"error_detail"`
	GuildID     string `json:"guild_id"`
}

// -----------------------------------------------------------------------------
// Interaction Payloads
// -----------------------------------------------------------------------------

// InteractionRespondedPayloadV1 contains interaction response tracking data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type InteractionRespondedPayloadV1 struct {
	InteractionID string `json:"interaction_id"`
	UserID        string `json:"user_id"`
	Status        string `json:"status"`
	ErrorDetail   string `json:"error_detail,omitempty"`
	GuildID       string `json:"guild_id"`
}

// InteractionResponseV1 contains Discord interaction response data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type InteractionResponseV1 struct {
	InteractionID string        `json:"interaction_id"`
	Token         string        `json:"token"`
	Message       string        `json:"message"`
	RetryData     *RetryDataV1  `json:"retry_data,omitempty"`
}

// RetryDataV1 contains data for retrying a failed round creation.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RetryDataV1 struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	StartTime   string `json:"start_time"`
	Location    string `json:"location"`
}

// -----------------------------------------------------------------------------
// Trace Payloads
// -----------------------------------------------------------------------------

// TracePayloadV1 contains trace event data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type TracePayloadV1 struct {
	Message string `json:"message"`
	GuildID string `json:"guild_id"`
}
