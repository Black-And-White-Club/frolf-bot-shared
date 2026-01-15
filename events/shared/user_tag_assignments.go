// Package sharedevents contains cross-module shared events.
//
// This file defines shared user tag assignment events used between user and leaderboard modules.
package sharedevents

import (
	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
)

// =============================================================================
// TAG ASSIGNMENT FLOW - Shared Event Constants
// =============================================================================

// TagAssignmentRequestedV1 is published to request tag assignment for a user.
//
// Pattern: Event Notification
// Subject: user.tag.assignment.requested.v1
// Producer: user-service (during creation)
// Consumers: leaderboard-service (assignment handler)
// Triggers: TagAssignedForUserCreationV1 OR TagAssignmentFailedV1
// Version: v1 (December 2024)
const TagAssignmentRequestedV1 = "user.tag.assignment.requested.v1"

// TagAssignedForUserCreationV1 is published when a tag is assigned during user creation.
//
// Pattern: Event Notification
// Subject: user.tag.assigned.for.creation.v1
// Producer: leaderboard-service
// Consumers: user-service (creation completion)
// Version: v1 (December 2024)
const TagAssignedForUserCreationV1 = "user.tag.assigned.for.creation.v1"

// TagAssignmentFailedV1 is published when tag assignment fails.
//
// Pattern: Event Notification
// Subject: user.tag.assignment.failed.v1
// Producer: leaderboard-service
// Consumers: user-service (error handler)
// Version: v1 (December 2024)
const TagAssignmentFailedV1 = "user.tag.assignment.failed.v1"

// =============================================================================
// TAG ASSIGNMENT FLOW - Shared Payload Types
// =============================================================================

// TagAssignmentRequestedPayloadV1 contains tag assignment request data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type TagAssignmentRequestedPayloadV1 struct {
	GuildID   sharedtypes.GuildID   `json:"guild_id"`
	UserID    sharedtypes.DiscordID `json:"user_id"`
	TagNumber sharedtypes.TagNumber `json:"tag_number"`
}

// TagAssignedForUserCreationPayloadV1 contains tag assignment success data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type TagAssignedForUserCreationPayloadV1 struct {
	GuildID   sharedtypes.GuildID    `json:"guild_id"`
	UserID    sharedtypes.DiscordID  `json:"user_id"`
	TagNumber *sharedtypes.TagNumber `json:"tag_number"`
}

// TagAssignmentFailedPayloadV1 contains tag assignment failure data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type TagAssignmentFailedPayloadV1 struct {
	GuildID   sharedtypes.GuildID   `json:"guild_id"`
	UserID    sharedtypes.DiscordID `json:"user_id"`
	TagNumber sharedtypes.TagNumber `json:"tag_number"`
	Reason    string                `json:"reason"`
}
