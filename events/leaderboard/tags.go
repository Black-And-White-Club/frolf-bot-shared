// Package leaderboardevents contains leaderboard-related domain events.
//
// This file defines the Leaderboard Tag Flow - events for tag assignment,
// availability checking, swapping, and batch operations.
//
// # Flow Sequences
//
// ## Tag Availability Check Flow
//  1. Request -> TagAvailabilityCheckRequestedV1
//  2. Available -> LeaderboardTagAvailableV1
//  3. OR Unavailable -> LeaderboardTagUnavailableV1
//  4. OR Failure -> TagAvailabilityCheckFailedV1
//
// ## Tag Assignment Flow
//  1. Request -> LeaderboardTagAssignmentRequestedV1
//  2. Success -> LeaderboardTagAssignedV1
//  3. OR Failure -> LeaderboardTagAssignmentFailedV1
//
// ## Batch Tag Assignment Flow
//  1. Request -> LeaderboardBatchTagAssignmentRequestedV1
//  2. Success -> LeaderboardBatchTagAssignedV1
//  3. OR Failure -> LeaderboardBatchTagAssignmentFailedV1
//
// ## Tag Swap Flow
//  1. Request -> TagSwapRequestedV1
//  2. Initiated -> TagSwapInitiatedV1
//  3. Processed -> TagSwapProcessedV1
//  4. OR Failure -> TagSwapFailedV1
//
// # Versioning Strategy
//
// All events include a V1 suffix for future schema evolution.
package leaderboardevents

import (
	sharedevents "github.com/Black-And-White-Club/frolf-bot-shared/events/shared"
	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
)

// =============================================================================
// TAG AVAILABILITY CHECK FLOW - Event Constants
// =============================================================================

// NOTE: Tag availability events are defined in events/shared/tag_availability.go.

// =============================================================================
// TAG ASSIGNMENT FLOW - Event Constants
// =============================================================================

// LeaderboardTagAssignmentRequestedV1 is published when tag assignment is requested.
//
// Pattern: Event Notification
// Subject: leaderboard.tag.assignment.requested.v1
// Producer: user-service, discord-service
// Consumers: leaderboard-service (assignment handler)
// Triggers: LeaderboardTagAssignedV1 OR LeaderboardTagAssignmentFailedV1
// Version: v1 (December 2024)
const LeaderboardTagAssignmentRequestedV1 = "leaderboard.tag.assignment.requested.v1"

// LeaderboardTagAssignedV1 is published when a tag is successfully assigned.
//
// Pattern: Event Notification
// Subject: leaderboard.tag.assigned.v1
// Producer: leaderboard-service
// Consumers: user-service, discord-service
// Version: v1 (December 2024)
const LeaderboardTagAssignedV1 = "leaderboard.tag.assigned.v1"

// LeaderboardTagAssignmentFailedV1 is published when tag assignment fails.
//
// Pattern: Event Notification
// Subject: leaderboard.tag.assignment.failed.v1
// Producer: leaderboard-service
// Consumers: user-service, discord-service, error handlers
// Version: v1 (December 2024)
const LeaderboardTagAssignmentFailedV1 = "leaderboard.tag.assignment.failed.v1"

// =============================================================================
// BATCH TAG ASSIGNMENT FLOW - Event Constants
// =============================================================================

// LeaderboardBatchTagAssignmentRequestedV1 is published for batch tag assignment.
//
// Pattern: Event Notification
// Subject: leaderboard.batch.tag.assignment.requested.v1
// Producer: score-service (after round processing)
// Consumers: leaderboard-service (batch handler)
// Triggers: LeaderboardBatchTagAssignedV1 OR LeaderboardBatchTagAssignmentFailedV1
// Version: v1 (December 2024)
const LeaderboardBatchTagAssignmentRequestedV1 = "leaderboard.batch.tag.assignment.requested.v1"

// LeaderboardBatchTagAssignedV1 is published when batch tag assignment completes.
//
// Pattern: Event Notification
// Subject: leaderboard.batch.tag.assigned.v1
// Producer: leaderboard-service
// Consumers: discord-service (confirmation)
// Version: v1 (December 2024)
const LeaderboardBatchTagAssignedV1 = "leaderboard.batch.tag.assigned.v1"

// LeaderboardBatchTagAssignmentFailedV1 is published when batch assignment fails.
//
// Pattern: Event Notification
// Subject: leaderboard.batch.tag.assignment.failed.v1
// Producer: leaderboard-service
// Consumers: discord-service, error handlers
// Version: v1 (December 2024)
const LeaderboardBatchTagAssignmentFailedV1 = "leaderboard.batch.tag.assignment.failed.v1"

// =============================================================================
// TAG SWAP FLOW - Event Constants
// =============================================================================

// TagSwapRequestedV1 is published when a tag swap is requested.
//
// Pattern: Event Notification
// Subject: leaderboard.tag.swap.requested.v1
// Producer: discord-service (admin command)
// Consumers: leaderboard-service (swap handler)
// Triggers: TagSwapInitiatedV1
// Version: v1 (December 2024)
const TagSwapRequestedV1 = "leaderboard.tag.swap.requested.v1"

// TagSwapInitiatedV1 is published when tag swap is initiated.
//
// Pattern: Event Notification
// Subject: leaderboard.tag.swap.initiated.v1
// Producer: leaderboard-service
// Consumers: leaderboard-service (swap processor)
// Triggers: TagSwapProcessedV1 OR TagSwapFailedV1
// Version: v1 (December 2024)
const TagSwapInitiatedV1 = "leaderboard.tag.swap.initiated.v1"

// TagSwapProcessedV1 is published when tag swap completes successfully.
//
// Pattern: Event Notification
// Subject: leaderboard.tag.swap.processed.v1
// Producer: leaderboard-service
// Consumers: discord-service (confirmation)
// Version: v1 (December 2024)
const TagSwapProcessedV1 = "leaderboard.tag.swap.processed.v1"

// TagSwapFailedV1 is published when tag swap fails.
//
// Pattern: Event Notification
// Subject: leaderboard.tag.swap.failed.v1
// Producer: leaderboard-service
// Consumers: discord-service (error handler)
// Version: v1 (December 2024)
const TagSwapFailedV1 = "leaderboard.tag.swap.failed.v1"

// =============================================================================
// TAG NUMBER LOOKUP FLOW - Event Constants
// =============================================================================

// NOTE: Tag number lookup events are defined in events/shared/tag_number_lookup.go.

// LeaderboardTraceEventV1 is published for leaderboard tracing/observability.
//
// Pattern: Event Notification
// Subject: leaderboard.trace.event.v1
// Producer: leaderboard-service
// Consumers: observability systems
// Version: v1 (December 2024)
const LeaderboardTraceEventV1 = "leaderboard.trace.event.v1"

// LeaderboardTagUpdatedV1 is published whenever a user's tag changes.
//
// Pattern: Event Notification
// Subject: leaderboard.tag.updated.v1
// Producer: leaderboard-service (after tag mutation)
// Consumers: round-service (update projections), discord-bot (update embeds)
// Triggers: Round participant updates, Discord embed updates
// Version: v1 (January 2026)
const LeaderboardTagUpdatedV1 = "leaderboard.tag.updated.v1"

// =============================================================================
// TAG FLOW - Payload Types
// =============================================================================

// -----------------------------------------------------------------------------
// Tag Availability Payloads
// -----------------------------------------------------------------------------

// NOTE: Tag availability payloads are defined in events/shared/tag_availability.go.

// -----------------------------------------------------------------------------
// Tag Assignment Payloads
// -----------------------------------------------------------------------------

// LeaderboardTagAssignmentRequestedPayloadV1 contains tag assignment request data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type LeaderboardTagAssignmentRequestedPayloadV1 struct {
	GuildID    sharedtypes.GuildID    `json:"guild_id"`
	UserID     sharedtypes.DiscordID  `json:"user_id"`
	TagNumber  *sharedtypes.TagNumber `json:"tag_number"`
	UpdateID   sharedtypes.RoundID    `json:"update_id"`
	Source     string                 `json:"source"`
	UpdateType string                 `json:"update_type"`
}

// LeaderboardTagAssignedPayloadV1 contains tag assignment success data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type LeaderboardTagAssignedPayloadV1 struct {
	GuildID      sharedtypes.GuildID    `json:"guild_id"`
	UserID       sharedtypes.DiscordID  `json:"user_id"`
	TagNumber    *sharedtypes.TagNumber `json:"tag_number"`
	AssignmentID sharedtypes.RoundID    `json:"assignment_id"`
	Source       string                 `json:"source"`
}

// LeaderboardTagAssignmentFailedPayloadV1 contains tag assignment failure data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type LeaderboardTagAssignmentFailedPayloadV1 struct {
	GuildID    sharedtypes.GuildID               `json:"guild_id"`
	UserID     sharedtypes.DiscordID             `json:"user_id"`
	TagNumber  *sharedtypes.TagNumber            `json:"tag_number"`
	Source     string                            `json:"source"`
	UpdateType string                            `json:"update_type"`
	Reason     string                            `json:"reason"`
	Config     *sharedevents.GuildConfigFragment `json:"config_fragment,omitempty"`
}

// -----------------------------------------------------------------------------
// Batch Tag Assignment Payloads
// -----------------------------------------------------------------------------

// TagAssignmentInfoV1 contains individual tag assignment info for batch operations.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type TagAssignmentInfoV1 struct {
	GuildID   sharedtypes.GuildID   `json:"guild_id"`
	UserID    sharedtypes.DiscordID `json:"user_id"`
	TagNumber sharedtypes.TagNumber `json:"tag_number"`
}

// LeaderboardBatchTagAssignmentRequestedPayloadV1 contains batch assignment request data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type LeaderboardBatchTagAssignmentRequestedPayloadV1 struct {
	GuildID          sharedtypes.GuildID   `json:"guild_id"`
	RequestingUserID sharedtypes.DiscordID `json:"requesting_user_id"`
	BatchID          string                `json:"batch_id"`
	Assignments      []TagAssignmentInfoV1 `json:"assignments"`
}

// LeaderboardBatchTagAssignedPayloadV1 contains batch assignment success data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type LeaderboardBatchTagAssignedPayloadV1 struct {
	GuildID          sharedtypes.GuildID               `json:"guild_id"`
	RequestingUserID sharedtypes.DiscordID             `json:"requesting_user_id"`
	BatchID          string                            `json:"batch_id"`
	AssignmentCount  int                               `json:"assignment_count"`
	Assignments      []TagAssignmentInfoV1             `json:"assignments"`
	Config           *sharedevents.GuildConfigFragment `json:"config_fragment,omitempty"`
}

// LeaderboardBatchTagAssignmentFailedPayloadV1 contains batch assignment failure data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type LeaderboardBatchTagAssignmentFailedPayloadV1 struct {
	GuildID          sharedtypes.GuildID   `json:"guild_id"`
	RequestingUserID sharedtypes.DiscordID `json:"requesting_user_id"`
	BatchID          string                `json:"batch_id"`
	Reason           string                `json:"reason"`
}

// -----------------------------------------------------------------------------
// Tag Swap Payloads
// -----------------------------------------------------------------------------

// TagSwapRequestedPayloadV1 contains tag swap request data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type TagSwapRequestedPayloadV1 struct {
	GuildID     sharedtypes.GuildID   `json:"guild_id"`
	RequestorID sharedtypes.DiscordID `json:"requestor_id"`
	TargetID    sharedtypes.DiscordID `json:"target_id"`
}

// TagSwapInitiatedPayloadV1 contains tag swap initiation data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type TagSwapInitiatedPayloadV1 struct {
	GuildID     sharedtypes.GuildID   `json:"guild_id"`
	RequestorID sharedtypes.DiscordID `json:"requestor_id"`
	TargetID    sharedtypes.DiscordID `json:"target_id"`
}

// TagSwapProcessedPayloadV1 contains tag swap success data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type TagSwapProcessedPayloadV1 struct {
	GuildID     sharedtypes.GuildID               `json:"guild_id"`
	RequestorID sharedtypes.DiscordID             `json:"requestor_id"`
	TargetID    sharedtypes.DiscordID             `json:"target_id"`
	Config      *sharedevents.GuildConfigFragment `json:"config_fragment,omitempty"`
}

// TagSwapFailedPayloadV1 contains tag swap failure data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type TagSwapFailedPayloadV1 struct {
	GuildID     sharedtypes.GuildID   `json:"guild_id"`
	RequestorID sharedtypes.DiscordID `json:"requestor_id"`
	TargetID    sharedtypes.DiscordID `json:"target_id"`
	Reason      string                `json:"reason"`
}

// -----------------------------------------------------------------------------
// Tag Number Lookup Payloads
// -----------------------------------------------------------------------------

// NOTE: Tag number lookup payloads are defined in events/shared/tag_number_lookup.go.

// =============================================================================
// TAG UPDATE EVENT - Canonical Event for Tag Synchronization
// =============================================================================

// LeaderboardTagUpdatedPayloadV1 is published whenever a user's tag changes in the leaderboard.
//
// Pattern: Event Notification
// Subject: leaderboard.tag.updated.v1
// Producer: leaderboard-service (after successful tag mutation)
// Consumers: round-service (update denormalized projections), discord-bot (update embeds)
// Triggers: Round participant tag updates, Discord embed updates
//
// Schema History:
//   - v1.0 (January 2026): Initial version
type LeaderboardTagUpdatedPayloadV1 struct {
	GuildID sharedtypes.GuildID   `json:"guild_id"`
	UserID  sharedtypes.DiscordID `json:"user_id"`

	OldTag *sharedtypes.TagNumber `json:"old_tag,omitempty"`
	NewTag *sharedtypes.TagNumber `json:"new_tag,omitempty"`

	// swap | assign | update | revoke | import | system | admin
	Reason string `json:"reason"`
}

// -----------------------------------------------------------------------------
// Interface Implementations for Legacy Compatibility
// -----------------------------------------------------------------------------

// GetUserID returns the UserID of the payload.
func (p *LeaderboardTagAssignmentRequestedPayloadV1) GetUserID() sharedtypes.DiscordID {
	return p.UserID
}

// GetTagNumber returns the TagNumber of the payload.
func (p *LeaderboardTagAssignmentRequestedPayloadV1) GetTagNumber() *sharedtypes.TagNumber {
	return p.TagNumber
}

// GetUserID returns the UserID of the payload.
func (p *LeaderboardTagAssignedPayloadV1) GetUserID() sharedtypes.DiscordID {
	return p.UserID
}

// GetTagNumber returns the TagNumber of the payload.
func (p *LeaderboardTagAssignedPayloadV1) GetTagNumber() *sharedtypes.TagNumber {
	return p.TagNumber
}
