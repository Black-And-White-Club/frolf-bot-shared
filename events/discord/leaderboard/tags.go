// Package leaderboard contains Discord-specific leaderboard events.
//
// This file defines the Discord Leaderboard Tag Operations Flow - events specific to
// assigning, swapping, and checking tag availability through Discord interactions.
//
// # Flow Sequences
//
// ## Tag Assignment Flow
//  1. Admin requests assignment -> LeaderboardTagAssignRequestV1
//  2. Success -> LeaderboardTagAssignedV1
//  3. OR Failure -> LeaderboardTagAssignFailedV1
//
// ## Tag Availability Flow
//  1. User checks availability -> LeaderboardTagAvailabilityRequestV1
//  2. Response -> LeaderboardTagAvailabilityResponseV1
//
// ## Tag Swap Flow
//  1. Admin requests swap -> LeaderboardTagSwapRequestV1
//  2. Success -> LeaderboardTagSwappedV1
//  3. OR Failure -> LeaderboardTagSwapFailedV1
//
// # Relationship to Domain Events
//
// These Discord events wrap/trigger domain events in events/leaderboard/:
//   - LeaderboardTagAssignRequestV1 -> publishes LeaderboardTagAssignmentRequested (domain)
//   - LeaderboardTagSwapRequestV1 -> publishes TagSwapRequested (domain)
//
// # Versioning Strategy
//
// All events include a V1 suffix for future schema evolution.
package leaderboard

import (
	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
)

// =============================================================================
// DISCORD LEADERBOARD TAG FLOW - Event Constants
// =============================================================================

// -----------------------------------------------------------------------------
// Tag Assignment Events
// -----------------------------------------------------------------------------

// LeaderboardTagAssignRequestV1 is published when an admin requests tag assignment.
//
// Pattern: Event Notification
// Subject: discord.leaderboard.tag.assign.request.v1
// Producer: discord-service (command handler)
// Consumers: discord-service (tag handler)
// Triggers: Domain LeaderboardTagAssignmentRequested
// Version: v1 (December 2024)
const LeaderboardTagAssignRequestV1 = "discord.leaderboard.tag.assign.request.v1"

// LeaderboardTagAssignedV1 is published when tag assignment succeeds.
//
// Pattern: Event Notification
// Subject: discord.leaderboard.tag.assigned.v1
// Producer: discord-service (after domain response)
// Consumers: discord-service (success message handler)
// Version: v1 (December 2024)
const LeaderboardTagAssignedV1 = "discord.leaderboard.tag.assigned.v1"

// LeaderboardTagAssignFailedV1 is published when tag assignment fails.
//
// Pattern: Event Notification
// Subject: discord.leaderboard.tag.assign.failed.v1
// Producer: discord-service
// Consumers: discord-service (error message handler)
// Version: v1 (December 2024)
const LeaderboardTagAssignFailedV1 = "discord.leaderboard.tag.assign.failed.v1"

// LeaderboardBatchTagAssignedV1 is published when batch tag assignment completes.
//
// Pattern: Event Notification
// Subject: discord.leaderboard.batch.tag.assigned.v1
// Producer: discord-service
// Consumers: discord-service (batch result handler)
// Version: v1 (December 2024)
const LeaderboardBatchTagAssignedV1 = "discord.leaderboard.batch.tag.assigned.v1"

// -----------------------------------------------------------------------------
// Tag Availability Events
// -----------------------------------------------------------------------------

// LeaderboardTagAvailabilityRequestV1 is published to check tag availability.
//
// Pattern: Event Notification
// Subject: discord.leaderboard.tag.availability.request.v1
// Producer: discord-service (command handler)
// Consumers: discord-service (availability handler)
// Triggers: Domain TagAvailabilityCheckRequested
// Version: v1 (December 2024)
const LeaderboardTagAvailabilityRequestV1 = "discord.leaderboard.tag.availability.request.v1"

// LeaderboardTagAvailabilityResponseV1 is published with availability result.
//
// Pattern: Event Notification
// Subject: discord.leaderboard.tag.availability.response.v1
// Producer: discord-service (after domain response)
// Consumers: discord-service (response handler)
// Version: v1 (December 2024)
const LeaderboardTagAvailabilityResponseV1 = "discord.leaderboard.tag.availability.response.v1"

// -----------------------------------------------------------------------------
// Tag Swap Events
// -----------------------------------------------------------------------------

// LeaderboardTagSwapRequestV1 is published when an admin requests tag swap.
//
// Pattern: Event Notification
// Subject: discord.leaderboard.tag.swap.request.v1
// Producer: discord-service (command handler)
// Consumers: discord-service (swap handler)
// Triggers: Domain TagSwapRequested
// Version: v1 (December 2024)
const LeaderboardTagSwapRequestV1 = "discord.leaderboard.tag.swap.request.v1"

// LeaderboardTagSwappedV1 is published when tag swap succeeds.
//
// Pattern: Event Notification
// Subject: discord.leaderboard.tag.swapped.v1
// Producer: discord-service (after domain response)
// Consumers: discord-service (success message handler)
// Version: v1 (December 2024)
const LeaderboardTagSwappedV1 = "discord.leaderboard.tag.swapped.v1"

// LeaderboardTagSwapFailedV1 is published when tag swap fails.
//
// Pattern: Event Notification
// Subject: discord.leaderboard.tag.swap.failed.v1
// Producer: discord-service
// Consumers: discord-service (error message handler)
// Version: v1 (December 2024)
const LeaderboardTagSwapFailedV1 = "discord.leaderboard.tag.swap.failed.v1"

// =============================================================================
// DISCORD LEADERBOARD TAG FLOW - Payload Types
// =============================================================================

// -----------------------------------------------------------------------------
// Tag Assignment Payloads
// -----------------------------------------------------------------------------

// LeaderboardTagAssignRequestPayloadV1 contains tag assignment request data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type LeaderboardTagAssignRequestPayloadV1 struct {
	TargetUserID sharedtypes.DiscordID `json:"target_user_id"`
	TagNumber    sharedtypes.TagNumber `json:"tag_number"`
	RequestorID  sharedtypes.DiscordID `json:"requestor_id"`
	ChannelID    string                `json:"channel_id"`
	MessageID    string                `json:"message_id"`
	GuildID      string                `json:"guild_id"`
}

// LeaderboardTagAssignedPayloadV1 contains assignment success data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type LeaderboardTagAssignedPayloadV1 struct {
	TargetUserID string                `json:"target_user_id"`
	TagNumber    sharedtypes.TagNumber `json:"tag_number"`
	ChannelID    string                `json:"channel_id"`
	MessageID    string                `json:"message_id"`
	GuildID      string                `json:"guild_id"`
}

// LeaderboardTagAssignFailedPayloadV1 contains assignment failure data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type LeaderboardTagAssignFailedPayloadV1 struct {
	TargetUserID string                `json:"target_user_id"`
	TagNumber    sharedtypes.TagNumber `json:"tag_number"`
	Reason       string                `json:"reason"`
	ChannelID    string                `json:"channel_id"`
	MessageID    string                `json:"message_id"`
	GuildID      string                `json:"guild_id"`
}

// TagAssignmentInfoV1 contains individual tag assignment info.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type TagAssignmentInfoV1 struct {
	UserID    sharedtypes.DiscordID `json:"user_id"`
	TagNumber sharedtypes.TagNumber `json:"tag_number"`
}

// BatchTagAssignedPayloadV1 contains batch assignment result data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type BatchTagAssignedPayloadV1 struct {
	RequestingUserID sharedtypes.DiscordID `json:"requesting_user_id"`
	BatchID          string                `json:"batch_id"`
	AssignmentCount  int                   `json:"assignment_count"`
	Assignments      []TagAssignmentInfoV1 `json:"assignments"`
	GuildID          string                `json:"guild_id"`
}

// -----------------------------------------------------------------------------
// Tag Availability Payloads
// -----------------------------------------------------------------------------

// LeaderboardTagAvailabilityRequestPayloadV1 contains availability check request.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type LeaderboardTagAvailabilityRequestPayloadV1 struct {
	TagNumber sharedtypes.TagNumber `json:"tag_number"`
	UserID    sharedtypes.DiscordID `json:"user_id"`
	ChannelID string                `json:"channel_id"`
	MessageID string                `json:"message_id"`
	GuildID   string                `json:"guild_id"`
}

// LeaderboardTagAvailabilityResponsePayloadV1 contains availability check result.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type LeaderboardTagAvailabilityResponsePayloadV1 struct {
	TagNumber sharedtypes.TagNumber `json:"tag_number"`
	Available bool                  `json:"available"`
	ChannelID string                `json:"channel_id"`
	MessageID string                `json:"message_id"`
	GuildID   string                `json:"guild_id"`
}

// -----------------------------------------------------------------------------
// Tag Swap Payloads
// -----------------------------------------------------------------------------

// LeaderboardTagSwapRequestPayloadV1 contains tag swap request data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type LeaderboardTagSwapRequestPayloadV1 struct {
	User1ID     sharedtypes.DiscordID `json:"user1_id"`
	User2ID     sharedtypes.DiscordID `json:"user2_id"`
	RequestorID sharedtypes.DiscordID `json:"requestor_id"`
	ChannelID   string                `json:"channel_id"`
	MessageID   string                `json:"message_id"`
	GuildID     string                `json:"guild_id"`
}

// LeaderboardTagSwappedPayloadV1 contains swap success data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type LeaderboardTagSwappedPayloadV1 struct {
	User1ID   sharedtypes.DiscordID `json:"user1_id"`
	User2ID   sharedtypes.DiscordID `json:"user2_id"`
	ChannelID string                `json:"channel_id"`
	MessageID string                `json:"message_id"`
	GuildID   string                `json:"guild_id"`
}

// LeaderboardTagSwapFailedPayloadV1 contains swap failure data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type LeaderboardTagSwapFailedPayloadV1 struct {
	User1ID   sharedtypes.DiscordID `json:"user1_id"`
	User2ID   sharedtypes.DiscordID `json:"user2_id"`
	Reason    string                `json:"reason"`
	ChannelID string                `json:"channel_id"`
	MessageID string                `json:"message_id"`
	GuildID   string                `json:"guild_id"`
}
