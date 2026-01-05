// Package userevents contains user-related domain events.
//
// This file defines the User Tag Flow - events for tag availability checks
// and tag assignments during user operations.
//
// # Flow Sequences
//
// ## Tag Availability Check Flow
//  1. Request -> TagAvailabilityCheckRequestedV1
//  2. Available -> TagAvailableV1
//  3. OR Unavailable -> TagUnavailableV1
//
// ## Tag Assignment Flow (for user creation)
//  1. Request -> TagAssignmentRequestedV1
//  2. Success -> TagAssignedForUserCreationV1
//  3. OR Failure -> TagAssignmentFailedV1
//
// # Relationship to Leaderboard Module
//
// Tag operations are coordinated with the leaderboard module:
//   - TagAvailabilityCheckRequestedV1 -> triggers leaderboard tag check
//   - TagAssignmentRequestedV1 -> triggers leaderboard tag assignment
//
// # Versioning Strategy
//
// All events include a V1 suffix for future schema evolution.
package userevents

import (
	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
)

// =============================================================================
// TAG AVAILABILITY CHECK FLOW - Event Constants
// =============================================================================

// TagAvailabilityCheckRequestedV1 is published to check if a tag is available.
//
// Pattern: Event Notification
// Subject: user.tag.availability.check.requested.v1
// Producer: user-service (during signup)
// Consumers: leaderboard-service (availability handler)
// Triggers: TagAvailableV1 OR TagUnavailableV1
// Version: v1 (December 2024)
const TagAvailabilityCheckRequestedV1 = "user.tag.availability.check.requested.v1"

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

// =============================================================================
// TAG ASSIGNMENT FLOW - Event Constants
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
// TAG FLOW - Payload Types
// =============================================================================

// -----------------------------------------------------------------------------
// Tag Availability Payloads
// -----------------------------------------------------------------------------

// TagAvailabilityCheckRequestedPayloadV1 contains tag availability check request data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type TagAvailabilityCheckRequestedPayloadV1 struct {
	GuildID   sharedtypes.GuildID   `json:"guild_id"`
	TagNumber sharedtypes.TagNumber `json:"tag_number"`
	UserID    sharedtypes.DiscordID `json:"user_id"`
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

// -----------------------------------------------------------------------------
// Tag Assignment Payloads
// -----------------------------------------------------------------------------

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
