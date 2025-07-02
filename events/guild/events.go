package guildevents

import (
	"time"

	guildtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/guild"
	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
)

// Event topics for guild config creation
const (
	GuildConfigCreationRequested = "guild.config.creation_requested"
	GuildConfigCreated           = "guild.config.created"
	GuildConfigCreationFailed    = "guild.config.creation_failed"
)

// Event topics for guild config retrieval
const (
	GuildConfigRetrievalRequested = "guild.config.retrieval_requested"
	GuildConfigRetrieved          = "guild.config.retrieved"
	GuildConfigRetrievalFailed    = "guild.config.retrieval_failed"
)

// Event topics for guild config update
const (
	GuildConfigUpdateRequested = "guild.config.update_requested"
	GuildConfigUpdated         = "guild.config.updated"
	GuildConfigUpdateFailed    = "guild.config.update_failed"
)

// Event topics for guild config deletion
const (
	GuildConfigDeletionRequested = "guild.config.deletion_requested"
	GuildConfigDeleted           = "guild.config.deleted"
	GuildConfigDeletionFailed    = "guild.config.deletion_failed"
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
