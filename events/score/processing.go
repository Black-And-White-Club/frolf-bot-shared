// Package scoreevents contains score-related domain events.
//
// This file defines the Score Processing Flow - events for processing round scores
// and triggering leaderboard updates.
//
// # Flow Sequence
//
//  1. Round finalized -> ProcessRoundScoresRequestedV1
//  2. Processing complete -> ProcessRoundScoresSucceededV1
//  3. OR Processing failed -> ProcessRoundScoresFailedV1
//
// # Relationship to Other Modules
//
// Score processing events trigger leaderboard events:
//   - ProcessRoundScoresSucceededV1 -> publishes LeaderboardBatchTagAssignmentRequested
//
// # Versioning Strategy
//
// All events include a V1 suffix for future schema evolution.
package scoreevents

import (
	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
)

// =============================================================================
// SCORE PROCESSING FLOW - Event Constants
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

// =============================================================================
// SCORE PROCESSING FLOW - Payload Types
// =============================================================================

// ProcessRoundScoresRequestedPayloadV1 contains round scores to be processed.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type ProcessRoundScoresRequestedPayloadV1 struct {
	GuildID   sharedtypes.GuildID     `json:"guild_id"`
	RoundID   sharedtypes.RoundID     `json:"round_id"`
	Scores    []sharedtypes.ScoreInfo `json:"scores"`
	Overwrite bool                    `json:"overwrite"`
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
