// Package sharedevents contains cross-module shared events.
//
// This file defines shared tag availability events used between user and leaderboard modules.
package sharedevents

import (
	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
)

// =============================================================================
// TAG AVAILABILITY CHECK FLOW - Shared Event Constants
// =============================================================================

// TagAvailabilityCheckRequestedV1 is published to check tag availability.
//
// Pattern: Event Notification
// Subject: leaderboard.tag.availability.check.requested.v1
// Producer: user-service (during signup)
// Consumers: leaderboard-service (availability handler)
// Triggers: TagAvailableV1 OR TagUnavailableV1
// Version: v1 (December 2024)
const TagAvailabilityCheckRequestedV1 = "leaderboard.tag.availability.check.requested.v1"

// TagAvailableV1 is published when a tag is available.
//
// Pattern: Event Notification
// Subject: user.tag.available.v1
// Producer: leaderboard-service
// Consumers: user-service (creation handler)
// Version: v1 (December 2024)
const TagAvailableV1 = "user.tag.available.v1"

// TagUnavailableV1 is published when a tag is not available.
//
// Pattern: Event Notification
// Subject: user.tag.unavailable.v1
// Producer: leaderboard-service
// Consumers: user-service (error handler)
// Version: v1 (December 2024)
const TagUnavailableV1 = "user.tag.unavailable.v1"

// TagAvailabilityCheckFailedV1 is published when availability check fails.
//
// Pattern: Event Notification
// Subject: leaderboard.tag.availability.check.failed.v1
// Producer: leaderboard-service
// Consumers: user-service, error handlers
// Version: v1 (December 2024)
const TagAvailabilityCheckFailedV1 = "leaderboard.tag.availability.check.failed.v1"

// =============================================================================
// TAG AVAILABILITY CHECK FLOW - Shared Payload Types
// =============================================================================

// TagAvailabilityCheckRequestedPayloadV1 contains tag availability check request data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type TagAvailabilityCheckRequestedPayloadV1 struct {
	GuildID   sharedtypes.GuildID    `json:"guild_id"`
	TagNumber *sharedtypes.TagNumber `json:"tag_number"`
	UserID    sharedtypes.DiscordID  `json:"user_id"`
}

// TagAvailablePayloadV1 contains tag availability success data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type TagAvailablePayloadV1 struct {
	GuildID   sharedtypes.GuildID   `json:"guild_id"`
	UserID    sharedtypes.DiscordID `json:"user_id"`
	TagNumber sharedtypes.TagNumber `json:"tag_number"`
}

// TagUnavailablePayloadV1 contains tag unavailability data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type TagUnavailablePayloadV1 struct {
	GuildID   sharedtypes.GuildID   `json:"guild_id"`
	UserID    sharedtypes.DiscordID `json:"user_id"`
	TagNumber sharedtypes.TagNumber `json:"tag_number"`
	Reason    string                `json:"reason"`
}

// TagAvailabilityCheckFailedPayloadV1 contains availability check failure data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type TagAvailabilityCheckFailedPayloadV1 struct {
	GuildID   sharedtypes.GuildID    `json:"guild_id"`
	UserID    sharedtypes.DiscordID  `json:"user_id"`
	TagNumber *sharedtypes.TagNumber `json:"tag_number"`
	Reason    string                 `json:"reason"`
}

// TagAvailabilityCheckResultPayloadV1 contains availability check result data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type TagAvailabilityCheckResultPayloadV1 struct {
	GuildID   sharedtypes.GuildID    `json:"guild_id"`
	UserID    sharedtypes.DiscordID  `json:"user_id"`
	TagNumber *sharedtypes.TagNumber `json:"tag_number"`
	Available bool                   `json:"tag_available"`
	Reason    string                 `json:"reason,omitempty"`
}
