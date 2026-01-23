// Package sharedevents contains cross-module shared events.
//
// This file defines shared score processing events used across modules.
package sharedevents

import (
	roundtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/round"
	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
)

// =============================================================================
// SCORE PROCESSING FLOW - Shared Event Constants
// =============================================================================

// ProcessRoundScoresRequestedV1 is published when a round's scores need processing.
//
// Pattern: Event Notification
// Subject: score.process.round.scores.requested.v1
// Producer: round-service (after round finalization)
// Consumers: score-service (processing handler)
// Triggers: ProcessRoundScoresSucceededV1 OR ProcessRoundScoresFailedV1
// Version: v1 (December 2024)
const ProcessRoundScoresRequestedV1 = "score.process.round.scores.requested.v1"

// ProcessRoundScoresSucceededV1 is published when round scores are processed successfully.
//
// Pattern: Event Notification
// Subject: score.process.round.scores.succeeded.v1
// Producer: score-service
// Consumers: leaderboard-service (batch tag assignment handler)
// Triggers: LeaderboardBatchTagAssignmentRequested
// Version: v1 (December 2024)
const ProcessRoundScoresSucceededV1 = "score.process.round.scores.succeeded.v1"

// ProcessRoundScoresFailedV1 is published when round score processing fails.
//
// Pattern: Event Notification
// Subject: score.process.round.scores.failed.v1
// Producer: score-service
// Consumers: monitoring, error handlers
// Version: v1 (December 2024)
const ProcessRoundScoresFailedV1 = "score.process.round.scores.failed.v1"

// ScoreModuleNotificationErrorV1 is published when score module notification fails.
//
// Pattern: Event Notification
// Subject: score.module.notification.error.v1
// Producer: backend-service (round module)
// Consumers: monitoring systems
// Version: v1 (December 2024)
const ScoreModuleNotificationErrorV1 = "score.module.notification.error.v1"

// =============================================================================
// SCORE PROCESSING FLOW - Shared Payload Types
// =============================================================================

// ProcessRoundScoresRequestedPayloadV1 contains round scores to be processed.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
//   - v1.1 (January 2026): Added RoundMode and Participants for doubles support
type ProcessRoundScoresRequestedPayloadV1 struct {
	GuildID   sharedtypes.GuildID     `json:"guild_id"`
	RoundID   sharedtypes.RoundID     `json:"round_id"`
	Scores    []sharedtypes.ScoreInfo `json:"scores"`
	Overwrite bool                    `json:"overwrite"`
	RoundMode sharedtypes.RoundMode   `json:"round_mode,omitempty"` // "SINGLES" or "DOUBLES"

	// Participants includes TeamID and other metadata needed to group scores for Doubles
	Participants []roundtypes.Participant `json:"participants,omitempty"`
}

// ProcessRoundScoresSucceededPayloadV1 contains processed score results.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type ProcessRoundScoresSucceededPayloadV1 struct {
	GuildID     sharedtypes.GuildID      `json:"guild_id"`
	RoundID     sharedtypes.RoundID      `json:"round_id"`
	TagMappings []sharedtypes.TagMapping `json:"tag_mappings"`
}

// ProcessRoundScoresFailedPayloadV1 contains score processing failure details.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type ProcessRoundScoresFailedPayloadV1 struct {
	GuildID sharedtypes.GuildID `json:"guild_id"`
	RoundID sharedtypes.RoundID `json:"round_id"`
	Reason  string              `json:"reason"`
}

// ScoreModuleNotificationErrorPayloadV1 contains score module error details.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type ScoreModuleNotificationErrorPayloadV1 struct {
	GuildID sharedtypes.GuildID `json:"guild_id"`
	RoundID sharedtypes.RoundID `json:"round_id"`
	Error   string              `json:"error"`
}
