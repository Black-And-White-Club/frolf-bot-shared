// Package guildevents contains guild-related domain events.
//
// MIGRATION NOTICE: This file contains legacy event constants.
// New code should use the versioned events from the flow-based files:
//   - config.go: GuildConfigCreationRequestedV1, GuildConfigCreatedV1, etc.
//
// See each file for detailed flow documentation and versioning information.
package guildevents

import (
	"time"

	guildtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/guild"
	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
)

// Event topics for guild config creation
// Deprecated: Use versioned constants from config.go
const (
	// Deprecated: Use GuildConfigCreationRequestedV1 from config.go
	GuildConfigCreationRequested = "guild.config.creation_requested"
	// Deprecated: Use GuildConfigCreatedV1 from config.go
	GuildConfigCreated = "guild.config.created"
	// Deprecated: Use GuildConfigCreationFailedV1 from config.go
	GuildConfigCreationFailed = "guild.config.creation_failed"
)

// Event topics for guild config retrieval
// Deprecated: Use versioned constants from config.go
const (
	// Deprecated: Use GuildConfigRetrievalRequestedV1 from config.go
	GuildConfigRetrievalRequested = "guild.config.retrieval_requested"
	// Deprecated: Use GuildConfigRetrievedV1 from config.go
	GuildConfigRetrieved = "guild.config.retrieved"
	// Deprecated: Use GuildConfigRetrievalFailedV1 from config.go
	GuildConfigRetrievalFailed = "guild.config.retrieval_failed"
)

// Event topics for guild config update
// Deprecated: Use versioned constants from config.go
const (
	// Deprecated: Use GuildConfigUpdateRequestedV1 from config.go
	GuildConfigUpdateRequested = "guild.config.update_requested"
	// Deprecated: Use GuildConfigUpdatedV1 from config.go
	GuildConfigUpdated = "guild.config.updated"
	// Deprecated: Use GuildConfigUpdateFailedV1 from config.go
	GuildConfigUpdateFailed = "guild.config.update_failed"
)

// Event topics for guild config deletion
// Deprecated: Use versioned constants from config.go
const (
	// Deprecated: Use GuildConfigDeletionRequestedV1 from config.go
	GuildConfigDeletionRequested = "guild.config.deletion_requested"
	// Deprecated: Use GuildConfigDeletedV1 from config.go
	GuildConfigDeleted = "guild.config.deleted"
	// Deprecated: Use GuildConfigDeletionFailedV1 from config.go
	GuildConfigDeletionFailed = "guild.config.deletion_failed"
)

// Emitted when a new guild config is created
type GuildConfigCreatedPayload struct {
	GuildID sharedtypes.GuildID    `json:"guild_id"`
	Config  guildtypes.GuildConfig `json:"config"`
}

// Emitted when guild config creation fails
type GuildConfigCreationFailedPayload struct {
	GuildID sharedtypes.GuildID `json:"guild_id"`
	Reason  string              `json:"reason"`
}

// Emitted when a guild config is updated
type GuildConfigUpdatedPayload struct {
	GuildID       sharedtypes.GuildID    `json:"guild_id"`
	Config        guildtypes.GuildConfig `json:"config"`
	UpdatedFields []string               `json:"updated_fields,omitempty"`
}

// Emitted when a guild config update fails
type GuildConfigUpdateFailedPayload struct {
	GuildID sharedtypes.GuildID `json:"guild_id"`
	Reason  string              `json:"reason"`
}

// Emitted when a guild config is retrieved
type GuildConfigRetrievedPayload struct {
	GuildID sharedtypes.GuildID    `json:"guild_id"`
	Config  guildtypes.GuildConfig `json:"config"`
}

// Emitted when a guild config retrieval fails
type GuildConfigRetrievalFailedPayload struct {
	GuildID sharedtypes.GuildID `json:"guild_id"`
	Reason  string              `json:"reason"`
}

// Emitted when a guild config deletion is requested
// This is the payload for the deletion request event
type GuildConfigDeletionRequestedPayload struct {
	GuildID sharedtypes.GuildID `json:"guild_id"`
}

// Emitted when a guild config is deleted
type GuildConfigDeletedPayload struct {
	GuildID sharedtypes.GuildID `json:"guild_id"`
}

// Emitted when a guild config deletion fails
type GuildConfigDeletionFailedPayload struct {
	GuildID sharedtypes.GuildID `json:"guild_id"`
	Reason  string              `json:"reason"`
}

// Emitted when a create guild config is requested
// This is the payload for the create config request event
// (used as input to the create handler)
type GuildConfigRequestedPayload struct {
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
	SetupCompletedAt     *time.Time          `json:"setup_completed_at"`
	// Add more fields as needed
}

// Emitted when a guild config retrieval is requested
// This is the payload for the retrieval request event
type GuildConfigRetrievalRequestedPayload struct {
	GuildID sharedtypes.GuildID `json:"guild_id"`
}

// Emitted when a guild config update is requested
// This is the payload for the update request event
type GuildConfigUpdateRequestedPayload struct {
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
	SetupCompletedAt     *time.Time          `json:"setup_completed_at"`
	// Add more fields as needed
}
