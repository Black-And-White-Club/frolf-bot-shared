// Package roundevents contains all round-related domain events.
//
// This file defines the Round Scoring Flow - all events related to
// updating participant scores and processing score submissions.
//
// # Flow Sequence
//
//  1. User submits score -> RoundScoreUpdateRequestedV1
//  2. Backend validates score -> RoundScoreUpdateValidatedV1
//  3. Score saved to database -> RoundParticipantScoreUpdatedV1
//  4. Check if all scores submitted:
//     a. All scores in -> RoundAllScoresSubmittedV1 (triggers finalization)
//     b. Scores remaining -> RoundScoresPartiallySubmittedV1
//  5. OR Score update fails -> RoundScoreUpdateErrorV1
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
// ROUND SCORING FLOW - Event Constants
// =============================================================================

// -----------------------------------------------------------------------------
// Score Update Flow
// -----------------------------------------------------------------------------

// RoundScoreUpdateRequestedV1 is published when a user submits a score.
//
// Pattern: Event Notification
// Subject: round.score.update.requested.v1
// Producer: discord-service (score submission handler)
// Consumers: backend-service (score validation handler)
// Triggers: RoundScoreUpdateValidatedV1 OR RoundScoreUpdateErrorV1
// Version: v1 (December 2024)
const RoundScoreUpdateRequestedV1 = "round.score.update.requested.v1"

// RoundScoreBulkUpdateRequestedV1 is published when bulk score overrides are submitted.
//
// Pattern: Event Notification
// Subject: round.score.bulk.update.requested.v1
// Producer: backend-service (score override handler)
// Consumers: backend-service (round bulk update handler)
// Triggers: RoundParticipantScoreUpdatedV1 OR RoundScoreUpdateErrorV1
// Version: v1 (January 2026)
const RoundScoreBulkUpdateRequestedV1 = "round.score.bulk.update.requested.v1"

// RoundScoreUpdateValidatedV1 is published when a score update passes validation.
//
// Pattern: Event Notification
// Subject: round.score.update.validated.v1
// Producer: backend-service (validation handler)
// Consumers: backend-service (score storage handler)
// Triggers: RoundParticipantScoreUpdatedV1
// Version: v1 (December 2024)
const RoundScoreUpdateValidatedV1 = "round.score.update.validated.v1"

// RoundParticipantScoreUpdatedV1 is published when a participant's score is saved.
//
// Pattern: Event Notification
// Subject: round.participant.score.updated.v1
// Producer: backend-service (score storage handler)
// Consumers: backend-service (score completion checker), discord-service (embed update)
// Triggers: RoundAllScoresSubmittedV1 OR RoundScoresPartiallySubmittedV1
// Version: v1 (December 2024)
const RoundParticipantScoreUpdatedV1 = "round.participant.score.updated.v1"

// RoundScoresBulkUpdatedV1 is published when bulk score overrides have been applied.
//
// Pattern: Event Notification
// Subject: round.scores.bulk.updated.v1
// Producer: backend-service (bulk score handler)
// Consumers: discord-service (embed update)
// Version: v1 (January 2026)
const RoundScoresBulkUpdatedV1 = "round.scores.bulk.updated.v1"

// RoundScoreUpdateErrorV1 is published when a score update fails.
//
// Pattern: Event Notification
// Subject: round.score.update.error.v1
// Producer: backend-service (any handler in score flow)
// Consumers: discord-service (error handler)
// Version: v1 (December 2024)
const RoundScoreUpdateErrorV1 = "round.score.update.error.v1"

// -----------------------------------------------------------------------------
// Team / Group-Specific Score Update Flow
// -----------------------------------------------------------------------------

// RoundTeamScoreUpdatedV1 is published when a team-based score update is saved.
//
// This event is emitted for doubles, triples, or other grouped play modes
// and allows downstream consumers (e.g. Discord) to render team-aware views
// without inferring grouping logic themselves.
//
// Pattern: Event Notification
// Subject: round.team.score.updated.v1
// Producer: backend-service (team score storage handler)
// Consumers: discord-service (team embed update)
// Version: v1 (January 2026)
const RoundTeamScoreUpdatedV1 = "round.team.score.updated.v1"

// RoundAllTeamScoresSubmittedV1 is published when all team scores have been submitted.
//
// This event mirrors RoundAllScoresSubmittedV1 but is explicitly scoped to
// team-based rounds so consumers can branch behavior without inspecting payloads.
//
// Pattern: Event Notification
// Subject: round.all.team.scores.submitted.v1
// Producer: backend-service (team score completion checker)
// Consumers: backend-service (finalization handler), discord-service
// Version: v1 (January 2026)
const RoundAllTeamScoresSubmittedV1 = "round.all.team.scores.submitted.v1"

// -----------------------------------------------------------------------------
// Score Completion Events
// -----------------------------------------------------------------------------

// RoundAllScoresSubmittedV1 is published when all participants have submitted scores.
//
// Pattern: Event Notification
// Subject: round.all.scores.submitted.v1
// Producer: backend-service (score completion checker)
// Consumers: backend-service (finalization handler)
// Triggers: RoundFinalizedV1
// Version: v1 (December 2024)
const RoundAllScoresSubmittedV1 = "round.all.scores.submitted.v1"

// RoundScoresPartiallySubmittedV1 is published when some scores remain to be submitted.
//
// Pattern: Event Notification
// Subject: round.scores.partially.submitted.v1
// Producer: backend-service (score completion checker)
// Consumers: discord-service (status update handler)
// Version: v1 (December 2024)
const RoundScoresPartiallySubmittedV1 = "round.scores.partially.submitted.v1"

// -----------------------------------------------------------------------------
// Score Processing Events (Cross-Module)
// -----------------------------------------------------------------------------

// NOTE: Score processing topics are defined in events/shared/score_processing.go.

// RoundScoresNotificationV1 is published to notify about round scores.
//
// Pattern: Event Notification
// Subject: round.scores.notification.v1
// Producer: backend-service
// Consumers: discord-service, other modules
// Version: v1 (December 2024)
const RoundScoresNotificationV1 = "round.scores.notification.v1"

// =============================================================================
// ROUND SCORING FLOW - Payload Types
// =============================================================================

// -----------------------------------------------------------------------------
// Score Update Payloads
// -----------------------------------------------------------------------------

// ScoreUpdateRequestPayloadV1 contains the score update request details.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type ScoreUpdateRequestPayloadV1 struct {
	GuildID   sharedtypes.GuildID   `json:"guild_id"`
	RoundID   sharedtypes.RoundID   `json:"round_id"`
	UserID    sharedtypes.DiscordID `json:"user_id"`
	Score     *sharedtypes.Score    `json:"score"`
	ChannelID string                `json:"channel_id"`
	MessageID string                `json:"message_id"`
}

// ScoreBulkUpdateRequestPayloadV1 contains the bulk score override request.
//
// Schema History:
//   - v1.0 (January 2026): Initial version
type ScoreBulkUpdateRequestPayloadV1 struct {
	GuildID   sharedtypes.GuildID           `json:"guild_id"`
	RoundID   sharedtypes.RoundID           `json:"round_id"`
	ChannelID string                        `json:"channel_id"`
	MessageID string                        `json:"message_id"`
	Updates   []ScoreUpdateRequestPayloadV1 `json:"updates"`
}

// ScoreUpdateValidatedPayloadV1 contains validated score update data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type ScoreUpdateValidatedPayloadV1 struct {
	GuildID                   sharedtypes.GuildID         `json:"guild_id"`
	ScoreUpdateRequestPayload ScoreUpdateRequestPayloadV1 `json:"score_update_request_payload"`
}

// ParticipantScoreUpdatedPayloadV1 contains the updated score result.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type ParticipantScoreUpdatedPayloadV1 struct {
	GuildID        sharedtypes.GuildID               `json:"guild_id"`
	RoundID        sharedtypes.RoundID               `json:"round_id"`
	UserID         sharedtypes.DiscordID             `json:"user_id"`
	Score          sharedtypes.Score                 `json:"score"`
	ChannelID      string                            `json:"channel_id"`
	EventMessageID string                            `json:"discord_message_id"`
	Participants   []roundtypes.Participant          `json:"participants,omitempty"`
	Config         *sharedevents.GuildConfigFragment `json:"config_fragment,omitempty"`
}

// TeamScoreUpdatedPayloadV1 contains the updated score result for a team.
//
// This payload is used exclusively for team-based rounds and allows
// consumers to render grouped scores deterministically.
//
// Schema History:
//   - v1.0 (January 2026): Initial version
type TeamScoreUpdatedPayloadV1 struct {
	GuildID        sharedtypes.GuildID       `json:"guild_id"`
	RoundID        sharedtypes.RoundID       `json:"round_id"`
	Team           roundtypes.NormalizedTeam `json:"team"`
	EventMessageID string                    `json:"discord_message_id"`
	Participants   []roundtypes.Participant  `json:"participants"`
}

// RoundScoresBulkUpdatedPayloadV1 contains the updated participant list after a bulk override.
//
// Schema History:
//   - v1.0 (January 2026): Initial version
type RoundScoresBulkUpdatedPayloadV1 struct {
	GuildID        sharedtypes.GuildID      `json:"guild_id"`
	RoundID        sharedtypes.RoundID      `json:"round_id"`
	EventMessageID string                   `json:"discord_message_id"`
	ChannelID      string                   `json:"channel_id"`
	Participants   []roundtypes.Participant `json:"participants"`
}

// RoundScoreUpdateErrorPayloadV1 contains score update error details.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoundScoreUpdateErrorPayloadV1 struct {
	GuildID            sharedtypes.GuildID          `json:"guild_id"`
	ScoreUpdateRequest *ScoreUpdateRequestPayloadV1 `json:"score_update_request"`
	Error              string                       `json:"error"`
}

// -----------------------------------------------------------------------------
// Score Completion Payloads
// -----------------------------------------------------------------------------

// AllScoresSubmittedPayloadV1 contains the all-scores-submitted event data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type AllScoresSubmittedPayloadV1 struct {
	GuildID        sharedtypes.GuildID               `json:"guild_id"`
	RoundID        sharedtypes.RoundID               `json:"round_id"`
	EventMessageID string                            `json:"discord_message_id"`
	RoundData      roundtypes.Round                  `json:"round_data"`
	RoundMode      sharedtypes.RoundMode              `json:"round_mode"`
	Participants   []roundtypes.Participant          `json:"participants,omitempty"`
	Teams          []roundtypes.NormalizedTeam       `json:"teams,omitempty"`
	Config         *sharedevents.GuildConfigFragment `json:"config_fragment,omitempty"`
}

// ScoresPartiallySubmittedPayloadV1 contains partial scores submission data.
// Note: Renamed from NotAllScoresSubmittedPayload for clarity.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type ScoresPartiallySubmittedPayloadV1 struct {
	GuildID        sharedtypes.GuildID         `json:"guild_id"`
	RoundID        sharedtypes.RoundID         `json:"round_id"`
	UserID         sharedtypes.DiscordID       `json:"user_id"`
	Score          sharedtypes.Score           `json:"score"`
	EventMessageID string                      `json:"event_message_id"`
	Scores         []ParticipantScoreV1        `json:"scores"`
	Participants   []roundtypes.Participant    `json:"participants"`
	Teams          []roundtypes.NormalizedTeam `json:"teams,omitempty"`
}

// -----------------------------------------------------------------------------
// Score Processing Payloads (Cross-Module)
// -----------------------------------------------------------------------------

// ParticipantScoreV1 represents a participant's score (versioned type).
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type ParticipantScoreV1 struct {
	GuildID   sharedtypes.GuildID    `json:"guild_id"`
	UserID    sharedtypes.DiscordID  `json:"user_id"`
	TagNumber *sharedtypes.TagNumber `json:"tag_number"`
	Score     sharedtypes.Score      `json:"score"`
	TeamID    string                 `json:"team_id,omitempty"`
}

// RoundScoresNotificationPayloadV1 contains scores notification data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoundScoresNotificationPayloadV1 struct {
	GuildID sharedtypes.GuildID  `json:"guild_id"`
	RoundID sharedtypes.RoundID  `json:"round_id"`
	Scores  []ParticipantScoreV1 `json:"scores"`
}

// NOTE: Score processing payloads are defined in events/shared/score_processing.go.

// =============================================================================
// Helper Methods
// =============================================================================

// ToRoundFinalizedPayload converts AllScoresSubmittedPayloadV1 to RoundFinalizedPayloadV1.
func (p *AllScoresSubmittedPayloadV1) ToRoundFinalizedPayloadV1(round roundtypes.Round) RoundFinalizedPayloadV1 {
	// Make sure Teams are copied into the RoundData
	round.Teams = p.Teams
	return RoundFinalizedPayloadV1{
		GuildID:   p.GuildID,
		RoundID:   p.RoundID,
		RoundData: round,
	}
}
