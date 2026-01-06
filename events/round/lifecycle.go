// Package roundevents contains all round-related domain events.
//
// This file defines the Round Lifecycle Flow - events related to the
// round's state transitions from scheduled through started, finalized,
// and completed states.
//
// # Flow Sequences
//
// ## Round Start Flow
//  1. Scheduled time reached -> RoundStartedV1
//  2. Discord notified -> RoundStartedDiscordV1
//
// ## Round Finalization Flow
//  1. All scores submitted -> RoundFinalizedV1 (backend state updated)
//  2. Discord UI updated -> RoundFinalizedDiscordV1
//  3. OR Finalization fails -> RoundFinalizationErrorV1
//
// ## Reminder Flow
//  1. Reminder scheduled -> RoundReminderScheduledV1
//  2. Reminder sent to Discord -> RoundReminderSentV1
//
// # Pattern Reference
//
// This flow follows the Event Notification pattern (Martin Fowler) where each step
// publishes an event to notify downstream consumers of state changes.
//
// # Versioning Strategy
//
// All events include a V1 suffix in the constant name and .v1 suffix in the topic
// string for future schema evolution while maintaining backward compatibility.
package roundevents

import (
	sharedevents "github.com/Black-And-White-Club/frolf-bot-shared/events/shared"
	roundtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/round"
	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
)

// =============================================================================
// ROUND LIFECYCLE FLOW - Event Constants
// =============================================================================

// -----------------------------------------------------------------------------
// Round Start Events
// -----------------------------------------------------------------------------

// RoundStartedV1 is published when a round's scheduled start time is reached.
//
// Pattern: Event Notification
// Subject: round.started.v1
// Producer: backend-service (scheduler)
// Consumers: backend-service (round state handler), discord-service
// Triggers: RoundStartedDiscordV1
// Version: v1 (December 2024)
const RoundStartedV1 = "round.started.v1"

// RoundStartedDiscordV1 is published to notify Discord that a round has started.
//
// Pattern: Event Notification
// Subject: round.started.discord.v1
// Producer: backend-service
// Consumers: discord-service (start notification handler)
// Triggers: Discord start notification/embed
// Version: v1 (December 2024)
const RoundStartedDiscordV1 = "round.started.discord.v1"

// -----------------------------------------------------------------------------
// Round Finalization Events
// -----------------------------------------------------------------------------

// RoundFinalizedV1 is published when a round is finalized (all scores in, DB updated).
// This is the backend-only event indicating database state has been updated.
//
// Pattern: Event Notification
// Subject: round.finalized.v1
// Producer: backend-service (finalization handler)
// Consumers: backend-service (score processor), discord-service
// Triggers: RoundFinalizedDiscordV1, score processing
// Version: v1 (December 2024)
const RoundFinalizedV1 = "round.finalized.v1"

// RoundFinalizedDiscordV1 is published for Discord-specific finalization consumers.
// It contains the fields required to update/patch the Discord embed.
// This is kept separate from RoundFinalizedV1 to avoid mixing domain and
// integration concerns.
//
// Pattern: Event Notification
// Subject: round.finalized.discord.v1
// Producer: backend-service
// Consumers: discord-service (embed update handler)
// Triggers: Discord embed updated with final scores
// Version: v1 (December 2024)
const RoundFinalizedDiscordV1 = "round.finalized.discord.v1"

// RoundCompletedV1 was intended for external app notification after round processing.
//
// DEPRECATED: This event is redundant with RoundFinalizedV1 and RoundFinalizedDiscordV1.
// RoundFinalizedV1 already signals completion of all backend processing to domain consumers.
// RoundFinalizedDiscordV1 provides Discord-specific finalization data.
// External apps (PWA, mobile apps) should subscribe to RoundFinalizedV1 instead.
// This constant will be removed in v2.0.
//
// Pattern: Event Notification
// Subject: round.completed.v1
// Producer: NONE - Never implemented
// Consumers: NONE - Never implemented
// Version: v1 (December 2024) - Deprecated before implementation
const RoundCompletedV1 = "round.completed.v1"

// RoundFinalizationErrorV1 is published when round finalization fails.
//
// Pattern: Event Notification
// Subject: round.finalization.error.v1
// Producer: backend-service (finalization handler)
// Consumers: discord-service (error handler), monitoring systems
// Version: v1 (December 2024)
const RoundFinalizationErrorV1 = "round.finalization.error.v1"

// -----------------------------------------------------------------------------
// Reminder Events
// -----------------------------------------------------------------------------

// RoundReminderScheduledV1 is published when a round reminder is scheduled.
// Note: Renamed from RoundReminder for clarity.
//
// Pattern: Event Notification
// Subject: round.reminder.scheduled.v1
// Producer: backend-service (scheduler)
// Consumers: backend-service (reminder processor)
// Triggers: RoundReminderSentV1
// Version: v1 (December 2024)
const RoundReminderScheduledV1 = "round.reminder.scheduled.v1"

// RoundReminderSentV1 is published when a reminder is sent to Discord.
// Note: Renamed from RoundDiscordReminder for clarity.
//
// Pattern: Event Notification
// Subject: round.reminder.sent.v1
// Producer: backend-service (reminder handler)
// Consumers: discord-service (reminder notification handler)
// Triggers: Discord reminder message sent
// Version: v1 (December 2024)
const RoundReminderSentV1 = "round.reminder.sent.v1"

// =============================================================================
// ROUND LIFECYCLE FLOW - Payload Types
// =============================================================================

// -----------------------------------------------------------------------------
// Round Start Payloads
// -----------------------------------------------------------------------------

// RoundStartedPayloadV1 contains the round start event data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoundStartedPayloadV1 struct {
	GuildID   sharedtypes.GuildID               `json:"guild_id"`
	RoundID   sharedtypes.RoundID               `json:"round_id"`
	Title     roundtypes.Title                  `json:"title"`
	Location  *roundtypes.Location              `json:"location"`
	StartTime *sharedtypes.StartTime            `json:"start_time"`
	ChannelID string                            `json:"channel_id"`
	Config    *sharedevents.GuildConfigFragment `json:"config_fragment,omitempty"`
}

// DiscordRoundStartPayloadV1 contains Discord-specific round start data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type DiscordRoundStartPayloadV1 struct {
	GuildID          sharedtypes.GuildID               `json:"guild_id"`
	RoundID          sharedtypes.RoundID               `json:"round_id"`
	Title            roundtypes.Title                  `json:"title"`
	Location         *roundtypes.Location              `json:"location"`
	StartTime        *sharedtypes.StartTime            `json:"start_time"`
	Participants     []RoundParticipantV1              `json:"participants"`
	DiscordChannelID string                            `json:"discord_channel_id,omitempty"`
	DiscordGuildID   string                            `json:"discord_guild_id,omitempty"`
	EventMessageID   string                            `json:"event_message_id"`
	Config           *sharedevents.GuildConfigFragment `json:"config_fragment,omitempty"`
}

// -----------------------------------------------------------------------------
// Round Finalization Payloads
// -----------------------------------------------------------------------------

// RoundFinalizedPayloadV1 contains the backend finalization event data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoundFinalizedPayloadV1 struct {
	GuildID   sharedtypes.GuildID               `json:"guild_id"`
	RoundID   sharedtypes.RoundID               `json:"round_id"`
	RoundData roundtypes.Round                  `json:"round_data"`
	Config    *sharedevents.GuildConfigFragment `json:"config_fragment,omitempty"`
}

// RoundFinalizedDiscordPayloadV1 is a Discord-specific payload emitted when a round
// has been finalized. It contains the snapshot of the round needed by Discord
// consumers to update or finalize the scorecard embed.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoundFinalizedDiscordPayloadV1 struct {
	GuildID          sharedtypes.GuildID      `json:"guild_id"`
	RoundID          sharedtypes.RoundID      `json:"round_id"`
	Title            roundtypes.Title         `json:"title"`
	StartTime        *sharedtypes.StartTime   `json:"start_time"`
	Location         *roundtypes.Location     `json:"location"`
	Participants     []roundtypes.Participant `json:"participants"`
	EventMessageID   string                   `json:"event_message_id"`
	DiscordChannelID string                   `json:"discord_channel_id,omitempty"`
}

// RoundFinalizedEmbedUpdatePayloadV1 contains embed update data for finalization.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoundFinalizedEmbedUpdatePayloadV1 struct {
	GuildID          sharedtypes.GuildID      `json:"guild_id"`
	RoundID          sharedtypes.RoundID      `json:"round_id"`
	Title            roundtypes.Title         `json:"title"`
	StartTime        *sharedtypes.StartTime   `json:"start_time"`
	Location         *roundtypes.Location     `json:"location"`
	Participants     []roundtypes.Participant `json:"participants"`
	EventMessageID   string                   `json:"discord_message_id"`
	DiscordChannelID string                   `json:"discord_channel_id,omitempty"`
}

// RoundCompletedPayloadV1 is published after backend completes all score processing.
// This is what Discord/PWA should consume to update their UI.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoundCompletedPayloadV1 struct {
	GuildID   sharedtypes.GuildID               `json:"guild_id"`
	RoundID   sharedtypes.RoundID               `json:"round_id"`
	RoundData roundtypes.Round                  `json:"round_data"`
	Config    *sharedevents.GuildConfigFragment `json:"config_fragment,omitempty"`
}

// RoundFinalizationErrorPayloadV1 contains finalization error details.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoundFinalizationErrorPayloadV1 struct {
	GuildID sharedtypes.GuildID `json:"guild_id"`
	RoundID sharedtypes.RoundID `json:"round_id"`
	Error   string              `json:"error"`
}

// -----------------------------------------------------------------------------
// Reminder Payloads
// -----------------------------------------------------------------------------

// DiscordReminderPayloadV1 contains Discord reminder event data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type DiscordReminderPayloadV1 struct {
	GuildID          sharedtypes.GuildID     `json:"guild_id"`
	RoundID          sharedtypes.RoundID     `json:"round_id"`
	ReminderType     string                  `json:"reminder_type"`
	RoundTitle       roundtypes.Title        `json:"round_title"`
	StartTime        *sharedtypes.StartTime  `json:"start_time"`
	Location         *roundtypes.Location    `json:"location"`
	UserIDs          []sharedtypes.DiscordID `json:"user_ids"`
	DiscordChannelID string                  `json:"discord_channel_id,omitempty"`
	DiscordGuildID   string                  `json:"discord_guild_id,omitempty"`
	EventMessageID   string                  `json:"event_message_id"`
}

// RoundReminderProcessedPayloadV1 contains reminder processing confirmation.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoundReminderProcessedPayloadV1 struct {
	GuildID sharedtypes.GuildID `json:"guild_id"`
	RoundID sharedtypes.RoundID `json:"round_id"`
}

// -----------------------------------------------------------------------------
// Discord Update Payloads
// -----------------------------------------------------------------------------

// DiscordRoundParticipantV1 represents a participant for Discord display.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type DiscordRoundParticipantV1 struct {
	GuildID   sharedtypes.GuildID    `json:"guild_id"`
	UserID    sharedtypes.DiscordID  `json:"user_id"`
	TagNumber *sharedtypes.TagNumber `json:"tag_number"`
	Score     *sharedtypes.Score     `json:"score"`
}

// DiscordRoundUpdatePayloadV1 contains Discord round update data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type DiscordRoundUpdatePayloadV1 struct {
	GuildID         sharedtypes.GuildID      `json:"guild_id"`
	Participants    []roundtypes.Participant `json:"participants"`
	RoundIDs        []sharedtypes.RoundID    `json:"round_ids"`
	EventMessageIDs []string                 `json:"event_message_ids"`
}
