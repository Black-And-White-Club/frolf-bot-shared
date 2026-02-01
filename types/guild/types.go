// In frolf-bot-shared/types/guild/config.go

package guildtypes

import (
	"errors"
	"time"

	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
)

type GuildConfig struct {
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
	// Optional snapshot of resources (used for deletion flow)
	ResourceState ResourceState `json:"resource_state,omitempty"`
	// ...add any other fields you want to expose...
}

// ResourceState represents a snapshot of bot-created Discord resources and
// per-resource deletion results. This is intended to be used in events so
// consumers (Discord worker) know exactly which resources to delete.
type ResourceState struct {
	SignupChannelID      string                    `json:"signup_channel_id,omitempty"`
	SignupMessageID      string                    `json:"signup_message_id,omitempty"`
	EventChannelID       string                    `json:"event_channel_id,omitempty"`
	LeaderboardChannelID string                    `json:"leaderboard_channel_id,omitempty"`
	UserRoleID           string                    `json:"user_role_id,omitempty"`
	EditorRoleID         string                    `json:"editor_role_id,omitempty"`
	AdminRoleID          string                    `json:"admin_role_id,omitempty"`
	Results              map[string]DeletionResult `json:"results,omitempty"`
}

// DeletionResult records the outcome of an attempted deletion for a single resource.
type DeletionResult struct {
	Status    string     `json:"status"`          // e.g., "pending", "success", "failed"
	Error     string     `json:"error,omitempty"` // error message if any
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// IsEmpty reports whether the ResourceState contains any meaningful data.
// A nil receiver is considered empty.
func (rs *ResourceState) IsEmpty() bool {
	if rs == nil {
		return true
	}
	if rs.SignupChannelID != "" || rs.SignupMessageID != "" || rs.EventChannelID != "" || rs.LeaderboardChannelID != "" || rs.UserRoleID != "" || rs.EditorRoleID != "" || rs.AdminRoleID != "" {
		return false
	}
	return len(rs.Results) == 0
}

// Validate checks if the core fields required for a functional guild are present.
func (c *GuildConfig) Validate() error {
	if c.GuildID == "" {
		return errors.New("guild ID is required")
	}
	if c.SignupChannelID == "" {
		return errors.New("signup channel ID required")
	}
	if c.EventChannelID == "" {
		return errors.New("event channel ID required")
	}
	if c.LeaderboardChannelID == "" {
		return errors.New("leaderboard channel ID required")
	}
	if c.UserRoleID == "" {
		return errors.New("user role ID required")
	}
	if c.SignupEmoji == "" {
		return errors.New("signup emoji required")
	}
	return nil
}
