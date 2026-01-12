// Package guildevents contains guild-related domain events.
//
// This file defines the Guild Config Flow - events for creating, updating,
// retrieving, and deleting guild configurations.
//
// # Flow Sequences
//
// ## Config Creation Flow
//  1. Request -> GuildConfigCreationRequestedV1
//  2. Success -> GuildConfigCreatedV1
//  3. OR Failure -> GuildConfigCreationFailedV1
//
// ## Config Update Flow
//  1. Request -> GuildConfigUpdateRequestedV1
//  2. Success -> GuildConfigUpdatedV1
//  3. OR Failure -> GuildConfigUpdateFailedV1
//
// ## Config Retrieval Flow
//  1. Request -> GuildConfigRetrievalRequestedV1
//  2. Success -> GuildConfigRetrievedV1
//  3. OR Failure -> GuildConfigRetrievalFailedV1
//
// ## Config Deletion Flow
//  1. Request -> GuildConfigDeletionRequestedV1
//  2. Success -> GuildConfigDeletedV1
//  3. OR Failure -> GuildConfigDeletionFailedV1
//
// # Versioning Strategy
//
// All events include a V1 suffix for future schema evolution.
package guildevents

import (
	"time"

	guildtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/guild"
	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
)

// =============================================================================
// GUILD SETUP FLOW - Event Constants
// =============================================================================

// GuildSetupRequestedV1 is published when initial guild setup is triggered.
//
// Pattern: Event Notification
// Subject: guild.setup.v1
// Producer: discord-service (guild join handler)
// Consumers: guild-service (setup handler)
// Triggers: GuildConfigCreatedV1 OR GuildConfigCreationFailedV1
// Version: v1 (January 2026)
const GuildSetupRequestedV1 = "guild.setup.v1"

// =============================================================================
// GUILD CONFIG CREATION FLOW - Event Constants
// =============================================================================

// GuildConfigCreationRequestedV1 is published when guild config creation is requested.
//
// Pattern: Event Notification
// Subject: guild.config.creation.requested.v1
// Producer: discord-service (setup handler)
// Consumers: guild-service (creation handler)
// Triggers: GuildConfigCreatedV1 OR GuildConfigCreationFailedV1
// Version: v1 (December 2024)
const GuildConfigCreationRequestedV1 = "guild.config.creation.requested.v1"

// GuildConfigCreatedV1 is published when guild config is successfully created.
//
// Pattern: Event Notification
// Subject: guild.config.created.v1
// Producer: guild-service
// Consumers: discord-service (setup completion)
// Version: v1 (December 2024)
const GuildConfigCreatedV1 = "guild.config.created.v1"

// GuildConfigCreationFailedV1 is published when guild config creation fails.
//
// Pattern: Event Notification
// Subject: guild.config.creation.failed.v1
// Producer: guild-service
// Consumers: discord-service (error handler)
// Version: v1 (December 2024)
const GuildConfigCreationFailedV1 = "guild.config.creation.failed.v1"

// =============================================================================
// GUILD CONFIG UPDATE FLOW - Event Constants
// =============================================================================

// GuildConfigUpdateRequestedV1 is published when guild config update is requested.
//
// Pattern: Event Notification
// Subject: guild.config.update.requested.v1
// Producer: discord-service (settings handler)
// Consumers: guild-service (update handler)
// Triggers: GuildConfigUpdatedV1 OR GuildConfigUpdateFailedV1
// Version: v1 (December 2024)
const GuildConfigUpdateRequestedV1 = "guild.config.update.requested.v1"

// GuildConfigUpdatedV1 is published when guild config is successfully updated.
//
// Pattern: Event Notification
// Subject: guild.config.updated.v1
// Producer: guild-service
// Consumers: discord-service (confirmation)
// Version: v1 (December 2024)
const GuildConfigUpdatedV1 = "guild.config.updated.v1"

// GuildConfigUpdateFailedV1 is published when guild config update fails.
//
// Pattern: Event Notification
// Subject: guild.config.update.failed.v1
// Producer: guild-service
// Consumers: discord-service (error handler)
// Version: v1 (December 2024)
const GuildConfigUpdateFailedV1 = "guild.config.update.failed.v1"

// =============================================================================
// GUILD CONFIG RETRIEVAL FLOW - Event Constants
// =============================================================================

// GuildConfigRetrievalRequestedV1 is published when guild config retrieval is requested.
//
// Pattern: Event Notification
// Subject: guild.config.retrieval.requested.v1
// Producer: any service needing guild config
// Consumers: guild-service (retrieval handler)
// Triggers: GuildConfigRetrievedV1 OR GuildConfigRetrievalFailedV1
// Version: v1 (December 2024)
const GuildConfigRetrievalRequestedV1 = "guild.config.retrieval.requested.v1"

// GuildConfigRetrievedV1 is published with the guild config data.
//
// Pattern: Event Notification
// Subject: guild.config.retrieved.v1
// Producer: guild-service
// Consumers: requesting service
// Version: v1 (December 2024)
const GuildConfigRetrievedV1 = "guild.config.retrieved.v1"

// GuildConfigRetrievalFailedV1 is published when guild config retrieval fails.
//
// Pattern: Event Notification
// Subject: guild.config.retrieval.failed.v1
// Producer: guild-service
// Consumers: requesting service
// Version: v1 (December 2024)
const GuildConfigRetrievalFailedV1 = "guild.config.retrieval.failed.v1"

// =============================================================================
// GUILD CONFIG DELETION FLOW - Event Constants
// =============================================================================

// GuildConfigDeletionRequestedV1 is published when guild config deletion is requested.
//
// Pattern: Event Notification
// Subject: guild.config.deletion.requested.v1
// Producer: discord-service (admin command)
// Consumers: guild-service (deletion handler)
// Triggers: GuildConfigDeletedV1 OR GuildConfigDeletionFailedV1
// Version: v1 (December 2024)
const GuildConfigDeletionRequestedV1 = "guild.config.deletion.requested.v1"

// GuildConfigDeletedV1 is published when guild config is successfully deleted.
//
// Pattern: Event Notification
// Subject: guild.config.deleted.v1
// Producer: guild-service
// Consumers: discord-service (confirmation)
// Version: v1 (December 2024)
const GuildConfigDeletedV1 = "guild.config.deleted.v1"

// GuildConfigDeletionFailedV1 is published when guild config deletion fails.
//
// Pattern: Event Notification
// Subject: guild.config.deletion.failed.v1
// Producer: guild-service
// Consumers: discord-service (error handler)
// Version: v1 (December 2024)
const GuildConfigDeletionFailedV1 = "guild.config.deletion.failed.v1"

// GuildConfigDeletionResultsV1 is published after the Discord worker attempts
// to delete resources that were snapshot at config deletion time. It contains
// per-resource outcomes so integrators and admins can audit what happened.
//
// Pattern: Event Notification
// Subject: guild.config.deletion.results.v1
// Producer: discord-service (discord worker)
// Consumers: ops/monitoring, any interested services
// Version: v1 (January 2026)
const GuildConfigDeletionResultsV1 = "guild.config.deletion.results.v1"

// =============================================================================
// GUILD CONFIG FLOW - Payload Types
// =============================================================================

// -----------------------------------------------------------------------------
// Config Creation Payloads
// -----------------------------------------------------------------------------

// GuildConfigCreationRequestedPayloadV1 contains guild config creation request data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type GuildConfigCreationRequestedPayloadV1 struct {
	GuildID              sharedtypes.GuildID `json:"guild_id"`
	SignupChannelID      string              `json:"signup_channel_id"`
	SignupMessageID      string              `json:"signup_message_id"`
	EventChannelID       string              `json:"event_channel_id"`
	LeaderboardChannelID string              `json:"leaderboard_channel_id"`
	UserRoleID           string              `json:"user_role_id"`
	EditorRoleID         string              `json:"editor_role_id"`
	AdminRoleID          string              `json:"admin_role_id"`
	SignupEmoji          string              `json:"signup_emoji"`
	AutoSetupCompleted   bool                `json:"auto_setup_completed"`
	SetupCompletedAt     *time.Time          `json:"setup_completed_at,omitempty"`
}

// GuildConfigCreatedPayloadV1 contains created guild config data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type GuildConfigCreatedPayloadV1 struct {
	GuildID sharedtypes.GuildID    `json:"guild_id"`
	Config  guildtypes.GuildConfig `json:"config"`
}

// GuildConfigCreationFailedPayloadV1 contains guild config creation failure data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type GuildConfigCreationFailedPayloadV1 struct {
	GuildID sharedtypes.GuildID `json:"guild_id"`
	Reason  string              `json:"reason"`
}

// -----------------------------------------------------------------------------
// Config Update Payloads
// -----------------------------------------------------------------------------

// GuildConfigUpdateRequestedPayloadV1 contains guild config update request data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type GuildConfigUpdateRequestedPayloadV1 struct {
	GuildID              sharedtypes.GuildID `json:"guild_id"`
	SignupChannelID      string              `json:"signup_channel_id,omitempty"`
	SignupMessageID      string              `json:"signup_message_id,omitempty"`
	EventChannelID       string              `json:"event_channel_id,omitempty"`
	LeaderboardChannelID string              `json:"leaderboard_channel_id,omitempty"`
	UserRoleID           string              `json:"user_role_id,omitempty"`
	EditorRoleID         string              `json:"editor_role_id,omitempty"`
	AdminRoleID          string              `json:"admin_role_id,omitempty"`
	SignupEmoji          string              `json:"signup_emoji,omitempty"`
	AutoSetupCompleted   bool                `json:"auto_setup_completed,omitempty"`
	SetupCompletedAt     *time.Time          `json:"setup_completed_at,omitempty"`
}

// GuildConfigUpdatedPayloadV1 contains updated guild config data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type GuildConfigUpdatedPayloadV1 struct {
	GuildID       sharedtypes.GuildID    `json:"guild_id"`
	Config        guildtypes.GuildConfig `json:"config"`
	UpdatedFields []string               `json:"updated_fields,omitempty"`
}

// GuildConfigUpdateFailedPayloadV1 contains guild config update failure data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type GuildConfigUpdateFailedPayloadV1 struct {
	GuildID sharedtypes.GuildID `json:"guild_id"`
	Reason  string              `json:"reason"`
}

// -----------------------------------------------------------------------------
// Config Retrieval Payloads
// -----------------------------------------------------------------------------

// GuildConfigRetrievalRequestedPayloadV1 contains guild config retrieval request data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type GuildConfigRetrievalRequestedPayloadV1 struct {
	GuildID sharedtypes.GuildID `json:"guild_id"`
}

// GuildConfigRetrievedPayloadV1 contains retrieved guild config data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type GuildConfigRetrievedPayloadV1 struct {
	GuildID sharedtypes.GuildID    `json:"guild_id"`
	Config  guildtypes.GuildConfig `json:"config"`
}

// GuildConfigRetrievalFailedPayloadV1 contains guild config retrieval failure data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type GuildConfigRetrievalFailedPayloadV1 struct {
	GuildID sharedtypes.GuildID `json:"guild_id"`
	Reason  string              `json:"reason"`
}

// -----------------------------------------------------------------------------
// Config Deletion Payloads
// -----------------------------------------------------------------------------

// GuildConfigDeletionRequestedPayloadV1 contains guild config deletion request data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type GuildConfigDeletionRequestedPayloadV1 struct {
	GuildID sharedtypes.GuildID `json:"guild_id"`
}

// GuildConfigDeletedPayloadV1 contains guild config deletion success data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type GuildConfigDeletedPayloadV1 struct {
	GuildID sharedtypes.GuildID `json:"guild_id"`
	// ResourceState is an optional snapshot of the resources that were present
	// at deletion time. Consumers (Discord worker) should use this to perform
	// deletions and record per-resource outcomes.
	ResourceState guildtypes.ResourceState `json:"resource_state,omitempty"`
}

// GuildConfigDeletionFailedPayloadV1 contains guild config deletion failure data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type GuildConfigDeletionFailedPayloadV1 struct {
	GuildID sharedtypes.GuildID `json:"guild_id"`
	Reason  string              `json:"reason"`
}

// GuildConfigDeletionResultsPayloadV1 contains the per-resource deletion
// results produced by the Discord worker after processing a
// GuildConfigDeletedV1 event.
//
// Schema History:
//   - v1.0 (January 2026): Initial version
type GuildConfigDeletionResultsPayloadV1 struct {
	GuildID       sharedtypes.GuildID                  `json:"guild_id"`
	ResourceState guildtypes.ResourceState             `json:"resource_state,omitempty"`
	Results       map[string]guildtypes.DeletionResult `json:"results,omitempty"`
}
