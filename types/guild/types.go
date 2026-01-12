// In frolf-bot-shared/types/guild/config.go

package guildtypes

import (
	"time"

	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
)

type GuildConfig struct {
	GuildID              sharedtypes.GuildID `json:"guild_id"`
	SignupChannelID      string
	SignupMessageID      string
	EventChannelID       string
	LeaderboardChannelID string
	UserRoleID           string
	EditorRoleID         string
	AdminRoleID          string
	SignupEmoji          string
	AutoSetupCompleted   bool
	SetupCompletedAt     *time.Time
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
