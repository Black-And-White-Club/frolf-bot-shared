// Package sharedevents contains cross-module shared events.
//
// This file defines shared score update events used across modules and services.
package sharedevents

import (
	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
)

// =============================================================================
// SCORE UPDATE FLOW - Shared Event Constants
// =============================================================================

// ScoreUpdateRequestedV1 is published when a score update is requested.
//
// Pattern: Event Notification
// Subject: score.update.requested.v1
// Producer: round-service, discord-service
// Consumers: score-service (update handler)
// Triggers: ScoreUpdatedV1 OR ScoreUpdateFailedV1
// Version: v1 (December 2024)
const ScoreUpdateRequestedV1 = "score.update.requested.v1"

// ScoreUpdatedV1 is published when a score is successfully updated.
//
// Pattern: Event Notification
// Subject: score.updated.v1
// Producer: score-service
// Consumers: round-service, discord-service
// Version: v1 (December 2024)
const ScoreUpdatedV1 = "score.updated.v1"

// ScoreUpdateFailedV1 is published when a score update fails.
//
// Pattern: Event Notification
// Subject: score.update.failed.v1
// Producer: score-service
// Consumers: requesting service, error handlers
// Version: v1 (December 2024)
const ScoreUpdateFailedV1 = "score.update.failed.v1"

// ScoreBulkUpdateRequestedV1 is published when bulk score updates are requested.
//
// Pattern: Event Notification
// Subject: score.bulk.update.requested.v1
// Producer: round-service (import flow), discord-service (admin command)
// Consumers: score-service (bulk update handler)
// Triggers: ScoreBulkUpdatedV1
// Version: v1 (December 2024)
const ScoreBulkUpdateRequestedV1 = "score.bulk.update.requested.v1"

// ScoreBulkUpdatedV1 is published when bulk score updates complete.
//
// Pattern: Event Notification
// Subject: score.bulk.updated.v1
// Producer: score-service
// Consumers: round-service, discord-service
// Version: v1 (December 2024)
const ScoreBulkUpdatedV1 = "score.bulk.updated.v1"

// =============================================================================
// SCORE UPDATE FLOW - Shared Payload Types
// =============================================================================

// ScoreUpdateRequestedPayloadV1 contains score update request data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type ScoreUpdateRequestedPayloadV1 struct {
	GuildID   sharedtypes.GuildID    `json:"guild_id"`
	RoundID   sharedtypes.RoundID    `json:"round_id"`
	UserID    sharedtypes.DiscordID  `json:"user_id"`
	Score     sharedtypes.Score      `json:"score"`
	TagNumber *sharedtypes.TagNumber `json:"tag_number,omitempty"`
}

// ScoreUpdatedPayloadV1 contains successful score update data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type ScoreUpdatedPayloadV1 struct {
	GuildID sharedtypes.GuildID   `json:"guild_id"`
	RoundID sharedtypes.RoundID   `json:"round_id"`
	UserID  sharedtypes.DiscordID `json:"user_id"`
	Score   sharedtypes.Score     `json:"score"`
}

// ScoreUpdateFailedPayloadV1 contains score update failure data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type ScoreUpdateFailedPayloadV1 struct {
	GuildID sharedtypes.GuildID   `json:"guild_id"`
	RoundID sharedtypes.RoundID   `json:"round_id"`
	UserID  sharedtypes.DiscordID `json:"user_id"`
	Reason  string                `json:"reason"`
}

// ScoreBulkUpdateRequestedPayloadV1 contains bulk score update request data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type ScoreBulkUpdateRequestedPayloadV1 struct {
	GuildID sharedtypes.GuildID             `json:"guild_id"`
	RoundID sharedtypes.RoundID             `json:"round_id"`
	Updates []ScoreUpdateRequestedPayloadV1 `json:"updates"`
}

// ScoreBulkUpdatedPayloadV1 contains bulk score update completion data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type ScoreBulkUpdatedPayloadV1 struct {
	GuildID        sharedtypes.GuildID     `json:"guild_id"`
	RoundID        sharedtypes.RoundID     `json:"round_id"`
	AppliedCount   int                     `json:"applied_count"`
	FailedCount    int                     `json:"failed_count"`
	TotalRequested int                     `json:"total_requested"`
	UserIDsApplied []sharedtypes.DiscordID `json:"user_ids_applied"`
}
