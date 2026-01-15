// Package roundevents contains all round-related domain events.
//
// This file defines the Round Participant Flow - all events related to
// participants joining, declining, and being removed from rounds.
//
// # Flow Sequences
//
// ## Join Flow
//  1. User clicks join button -> RoundParticipantJoinRequestedV1
//  2. Backend validates join request -> RoundParticipantJoinValidatedV1
//  3. Backend looks up tag number -> (tag lookup events)
//  4. Participant added to round -> RoundParticipantJoinedV1
//  5. OR Join fails -> RoundParticipantJoinErrorV1
//
// ## Decline Flow
//  1. User clicks decline button -> RoundParticipantDeclinedV1
//
// ## Removal Flow
//  1. Admin/user requests removal -> RoundParticipantRemovalRequestedV1
//  2. Participant removed -> RoundParticipantRemovedV1
//  3. OR Removal fails -> RoundParticipantRemovalErrorV1
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
	"github.com/Black-And-White-Club/frolf-bot-shared/events"
	sharedevents "github.com/Black-And-White-Club/frolf-bot-shared/events/shared"
	roundtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/round"
	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
)

// =============================================================================
// ROUND PARTICIPANT FLOW - Event Constants
// =============================================================================

// -----------------------------------------------------------------------------
// Join Flow Events
// -----------------------------------------------------------------------------

// RoundParticipantJoinRequestedV1 is published when a user requests to join a round.
//
// Pattern: Event Notification
// Subject: round.participant.join.requested.v1
// Producer: discord-service (button handler)
// Consumers: backend-service (join validation handler)
// Triggers: RoundParticipantJoinValidatedV1 OR RoundParticipantJoinErrorV1
// Version: v1 (December 2024)
const RoundParticipantJoinRequestedV1 = "round.participant.join.requested.v1"

// RoundParticipantJoinValidatedV1 is published when a join request passes validation.
//
// Pattern: Event Notification
// Subject: round.participant.join.validated.v1
// Producer: backend-service (validation handler)
// Consumers: backend-service (tag lookup handler)
// Triggers: Tag lookup flow, then RoundParticipantJoinedV1
// Version: v1 (December 2024)
const RoundParticipantJoinValidatedV1 = "round.participant.join.validated.v1"

// RoundParticipantJoinValidationRequestedV1 is published to request join validation.
//
// Pattern: Event Notification
// Subject: round.participant.join.validation.requested.v1
// Producer: backend-service
// Consumers: backend-service (validation handler)
// Version: v1 (December 2024)
const RoundParticipantJoinValidationRequestedV1 = "round.participant.join.validation.requested.v1"

// RoundParticipantJoinedV1 is published when a participant successfully joins a round.
//
// Pattern: Event Notification
// Subject: round.participant.joined.v1
// Producer: backend-service (participant handler)
// Consumers: discord-service (embed update handler)
// Triggers: Discord embed updated with new participant
// Version: v1 (December 2024)
const RoundParticipantJoinedV1 = "round.participant.joined.v1"

// RoundParticipantJoinErrorV1 is published when a join request fails.
//
// Pattern: Event Notification
// Subject: round.participant.join.error.v1
// Producer: backend-service (any handler in join flow)
// Consumers: discord-service (error handler)
// Triggers: Discord error message to user
// Version: v1 (December 2024)
const RoundParticipantJoinErrorV1 = "round.participant.join.error.v1"

// -----------------------------------------------------------------------------
// Decline Flow Events
// -----------------------------------------------------------------------------

// RoundParticipantDeclinedV1 is published when a user declines to join a round.
//
// Pattern: Event Notification
// Subject: round.participant.declined.v1
// Producer: discord-service (button handler) or backend-service
// Consumers: discord-service (embed update handler)
// Triggers: Discord embed updated showing declined status
// Version: v1 (December 2024)
const RoundParticipantDeclinedV1 = "round.participant.declined.v1"

// -----------------------------------------------------------------------------
// Removal Flow Events
// -----------------------------------------------------------------------------

// RoundParticipantRemovalRequestedV1 is published when a participant removal is requested.
//
// Pattern: Event Notification
// Subject: round.participant.removal.requested.v1
// Producer: discord-service (admin action) or backend-service
// Consumers: backend-service (removal handler)
// Triggers: RoundParticipantRemovedV1 OR RoundParticipantRemovalErrorV1
// Version: v1 (December 2024)
const RoundParticipantRemovalRequestedV1 = "round.participant.removal.requested.v1"

// RoundParticipantRemovedV1 is published when a participant is removed from a round.
//
// Pattern: Event Notification
// Subject: round.participant.removed.v1
// Producer: backend-service (removal handler)
// Consumers: discord-service (embed update handler)
// Triggers: Discord embed updated without the removed participant
// Version: v1 (December 2024)
const RoundParticipantRemovedV1 = "round.participant.removed.v1"

// RoundParticipantRemovalErrorV1 is published when a participant removal fails.
//
// Pattern: Event Notification
// Subject: round.participant.removal.error.v1
// Producer: backend-service (removal handler)
// Consumers: discord-service (error handler)
// Version: v1 (December 2024)
const RoundParticipantRemovalErrorV1 = "round.participant.removal.error.v1"

// -----------------------------------------------------------------------------
// Status Events
// -----------------------------------------------------------------------------

// RoundParticipantStatusErrorV1 is published for general participant status errors.
//
// Pattern: Event Notification
// Subject: round.participant.status.error.v1
// Producer: backend-service
// Consumers: discord-service (error handler)
// Version: v1 (December 2024)
const RoundParticipantStatusErrorV1 = "round.participant.status.error.v1"

// RoundParticipantStatusFoundV1 is published when a participant's status is found.
//
// Pattern: Event Notification
// Subject: round.participant.status.found.v1
// Producer: backend-service (status lookup handler)
// Consumers: Various handlers
// Version: v1 (December 2024)
const RoundParticipantStatusFoundV1 = "round.participant.status.found.v1"

// RoundParticipantStatusCheckErrorV1 is published when a status check fails.
//
// Pattern: Event Notification
// Subject: round.participant.status.check.error.v1
// Producer: backend-service (status handler)
// Consumers: discord-service (error handler)
// Version: v1 (December 2024)
const RoundParticipantStatusCheckErrorV1 = "round.participant.status.check.error.v1"

// RoundParticipantUpdateErrorV1 is published when a participant update fails.
//
// Pattern: Event Notification
// Subject: round.participant.update.error.v1
// Producer: backend-service (update handler)
// Consumers: discord-service (error handler)
// Version: v1 (December 2024)
const RoundParticipantUpdateErrorV1 = "round.participant.update.error.v1"

// RoundParticipantStatusUpdateRequestedV1 is published to request a participant status update.
//
// Pattern: Event Notification
// Subject: round.participant.status.update.requested.v1
// Producer: backend-service
// Consumers: backend-service (status update handler)
// Version: v1 (December 2024)
const RoundParticipantStatusUpdateRequestedV1 = "round.participant.status.update.requested.v1"

// RoundParticipantsUpdatedV1 is published when round participants are updated due to tag changes.
//
// Pattern: Event Notification
// Subject: round.participants.updated.v1
// Producer: round-service (tag update handler)
// Consumers: discord-bot (embed update handler)
// Triggers: Discord embed updated with new tag numbers
// Version: v1 (January 2026)
const RoundParticipantsUpdatedV1 = "round.participants.updated.v1"

// =============================================================================
// ROUND PARTICIPANT FLOW - Payload Types
// =============================================================================

// -----------------------------------------------------------------------------
// Join Flow Payloads
// -----------------------------------------------------------------------------

// ParticipantJoinRequestPayloadV1 contains the join request details.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type ParticipantJoinRequestPayloadV1 struct {
	GuildID    sharedtypes.GuildID               `json:"guild_id"`
	RoundID    sharedtypes.RoundID               `json:"round_id"`
	UserID     sharedtypes.DiscordID             `json:"user_id"`
	Response   roundtypes.Response               `json:"response"`
	TagNumber  *sharedtypes.TagNumber            `json:"tag_number"`
	JoinedLate *bool                             `json:"joined_late,omitempty"`
	Config     *sharedevents.GuildConfigFragment `json:"config_fragment,omitempty"`
}

// ParticipantJoinValidatedPayloadV1 contains validated join request data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type ParticipantJoinValidatedPayloadV1 struct {
	GuildID                       sharedtypes.GuildID             `json:"guild_id"`
	ParticipantJoinRequestPayload ParticipantJoinRequestPayloadV1 `json:"participant_join_request_payload"`
}

// ParticipantJoinValidationRequestPayloadV1 contains the validation request.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type ParticipantJoinValidationRequestPayloadV1 struct {
	GuildID  sharedtypes.GuildID   `json:"guild_id"`
	RoundID  sharedtypes.RoundID   `json:"round_id"`
	UserID   sharedtypes.DiscordID `json:"user_id"`
	Response roundtypes.Response   `json:"response"`
}

// ParticipantJoinedPayloadV1 contains the successful join result.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type ParticipantJoinedPayloadV1 struct {
	GuildID               sharedtypes.GuildID               `json:"guild_id"`
	RoundID               sharedtypes.RoundID               `json:"round_id"`
	AcceptedParticipants  []roundtypes.Participant          `json:"accepted_participants"`
	DeclinedParticipants  []roundtypes.Participant          `json:"declined_participants"`
	TentativeParticipants []roundtypes.Participant          `json:"tentative_participants"`
	EventMessageID        string                            `json:"discord_message_id"`
	JoinedLate            *bool                             `json:"joined_late,omitempty"`
	Config                *sharedevents.GuildConfigFragment `json:"config_fragment,omitempty"`
}

// RoundParticipantJoinErrorPayloadV1 contains join failure details.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoundParticipantJoinErrorPayloadV1 struct {
	GuildID                sharedtypes.GuildID              `json:"guild_id"`
	ParticipantJoinRequest *ParticipantJoinRequestPayloadV1 `json:"participant_join_request"`
	Error                  string                           `json:"error"`
	EventMessageID         string                           `json:"discord_message_id"`
}

// -----------------------------------------------------------------------------
// Decline Flow Payloads
// -----------------------------------------------------------------------------

// ParticipantDeclinedPayloadV1 contains the decline event details.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type ParticipantDeclinedPayloadV1 struct {
	GuildID        sharedtypes.GuildID               `json:"guild_id"`
	RoundID        sharedtypes.RoundID               `json:"round_id"`
	UserID         sharedtypes.DiscordID             `json:"user_id"`
	EventMessageID string                            `json:"discord_message_id"`
	Config         *sharedevents.GuildConfigFragment `json:"config_fragment,omitempty"`
}

// ParticipantDeclineErrorPayloadV1 contains decline error details.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type ParticipantDeclineErrorPayloadV1 struct {
	GuildID sharedtypes.GuildID   `json:"guild_id"`
	RoundID sharedtypes.RoundID   `json:"round_id"`
	UserID  sharedtypes.DiscordID `json:"user_id"`
	Error   string                `json:"error"`
}

// -----------------------------------------------------------------------------
// Removal Flow Payloads
// -----------------------------------------------------------------------------

// ParticipantRemovalRequestPayloadV1 contains the removal request details.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type ParticipantRemovalRequestPayloadV1 struct {
	GuildID sharedtypes.GuildID   `json:"guild_id"`
	RoundID sharedtypes.RoundID   `json:"round_id"`
	UserID  sharedtypes.DiscordID `json:"user_id"`
}

// ParticipantRemovedPayloadV1 contains the successful removal result.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type ParticipantRemovedPayloadV1 struct {
	GuildID               sharedtypes.GuildID               `json:"guild_id"`
	RoundID               sharedtypes.RoundID               `json:"round_id"`
	UserID                sharedtypes.DiscordID             `json:"user_id"`
	AcceptedParticipants  []roundtypes.Participant          `json:"accepted_participants"`
	DeclinedParticipants  []roundtypes.Participant          `json:"declined_participants"`
	TentativeParticipants []roundtypes.Participant          `json:"tentative_participants"`
	EventMessageID        string                            `json:"discord_message_id"`
	Config                *sharedevents.GuildConfigFragment `json:"config_fragment,omitempty"`
}

// ParticipantRemovalErrorPayloadV1 contains removal error details.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type ParticipantRemovalErrorPayloadV1 struct {
	GuildID sharedtypes.GuildID   `json:"guild_id"`
	RoundID sharedtypes.RoundID   `json:"round_id"`
	UserID  sharedtypes.DiscordID `json:"user_id"`
	Error   string                `json:"error"`
}

// -----------------------------------------------------------------------------
// Status Payloads
// -----------------------------------------------------------------------------

// ParticipantStatusRequestPayloadV1 contains a status request.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type ParticipantStatusRequestPayloadV1 struct {
	GuildID sharedtypes.GuildID   `json:"guild_id"`
	RoundID sharedtypes.RoundID   `json:"round_id"`
	UserID  sharedtypes.DiscordID `json:"user_id"`
}

// ParticipantStatusFoundPayloadV1 contains a found participant status.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type ParticipantStatusFoundPayloadV1 struct {
	GuildID sharedtypes.GuildID   `json:"guild_id"`
	RoundID sharedtypes.RoundID   `json:"round_id"`
	UserID  sharedtypes.DiscordID `json:"user_id"`
	Status  string                `json:"status"`
}

// ParticipantStatusCheckErrorPayloadV1 contains status check error details.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type ParticipantStatusCheckErrorPayloadV1 struct {
	GuildID sharedtypes.GuildID   `json:"guild_id"`
	RoundID sharedtypes.RoundID   `json:"round_id"`
	UserID  sharedtypes.DiscordID `json:"user_id"`
	Error   string                `json:"error"`
}

// ParticipantUpdateErrorPayloadV1 contains participant update error details.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type ParticipantUpdateErrorPayloadV1 struct {
	GuildID sharedtypes.GuildID   `json:"guild_id"`
	RoundID sharedtypes.RoundID   `json:"round_id"`
	UserID  sharedtypes.DiscordID `json:"user_id"`
	Error   string                `json:"error"`
}

// -----------------------------------------------------------------------------
// PARTICIPANT UPDATE EVENT - For Tag Synchronization
// -----------------------------------------------------------------------------

// RoundParticipantsUpdatedV1 is published when round participants are updated due to tag changes.
//
// Pattern: Event Notification
// Subject: round.participants.updated.v1
// Producer: round-service (after updating participant tags from leaderboard changes)
// Consumers: discord-bot (update embeds with new tag numbers)
// Triggers: Discord embed updates
//
// Schema History:
//   - v1.0 (January 2026): Initial version
type RoundParticipantsUpdatedPayloadV1 struct {
	GuildID sharedtypes.GuildID `json:"guild_id"`
	RoundID sharedtypes.RoundID `json:"round_id"`
	Round   roundtypes.Round    `json:"round"`

	Metadata events.CommonMetadata `json:"metadata"`
}

// -----------------------------------------------------------------------------
// Shared Types
// -----------------------------------------------------------------------------

// RoundParticipantV1 represents a participant in a round (versioned type).
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoundParticipantV1 struct {
	GuildID   sharedtypes.GuildID    `json:"guild_id"`
	UserID    sharedtypes.DiscordID  `json:"user_id"`
	TagNumber *sharedtypes.TagNumber `json:"tag_number"`
	Response  roundtypes.Response    `json:"response"`
	Score     *sharedtypes.Score     `json:"score"`
}

// =============================================================================
// Interface Implementations for V1 Payloads
// =============================================================================

// GetRoundID implements ParticipantUpdatePayload for ParticipantDeclinedPayloadV1.
func (p *ParticipantDeclinedPayloadV1) GetRoundID() sharedtypes.RoundID {
	return p.RoundID
}

// GetUserID implements ParticipantUpdatePayload for ParticipantDeclinedPayloadV1.
func (p *ParticipantDeclinedPayloadV1) GetUserID() sharedtypes.DiscordID {
	return p.UserID
}

// GetTagNumber implements ParticipantUpdatePayload for ParticipantDeclinedPayloadV1.
func (p *ParticipantDeclinedPayloadV1) GetTagNumber() *sharedtypes.TagNumber {
	return nil // No tag number for declined participants
}

// GetJoinedLate implements ParticipantUpdatePayload for ParticipantDeclinedPayloadV1.
func (p *ParticipantDeclinedPayloadV1) GetJoinedLate() *bool {
	return nil
}

// GetRoundID implements ParticipantUpdatePayload for ParticipantJoinRequestPayloadV1.
func (p *ParticipantJoinRequestPayloadV1) GetRoundID() sharedtypes.RoundID {
	return p.RoundID
}

// GetUserID implements ParticipantUpdatePayload for ParticipantJoinRequestPayloadV1.
func (p *ParticipantJoinRequestPayloadV1) GetUserID() sharedtypes.DiscordID {
	return p.UserID
}

// GetTagNumber implements ParticipantUpdatePayload for ParticipantJoinRequestPayloadV1.
func (p *ParticipantJoinRequestPayloadV1) GetTagNumber() *sharedtypes.TagNumber {
	return p.TagNumber
}

// GetJoinedLate implements ParticipantUpdatePayload for ParticipantJoinRequestPayloadV1.
func (p *ParticipantJoinRequestPayloadV1) GetJoinedLate() *bool {
	return p.JoinedLate
}
