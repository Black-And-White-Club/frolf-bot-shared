package sharedevents

import (
    guildtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/guild"
    sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
)

// GuildConfigFragment is a minimal, serializable subset of guild configuration data
// that Discord-side consumers need to avoid a secondary config lookup. Keep fields
// optional (omitempty) so older producers without this enrichment do not break decoding.
type GuildConfigFragment struct {
    GuildID              sharedtypes.GuildID `json:"guild_id"`
    SignupChannelID      string              `json:"signup_channel_id,omitempty"`
    SignupMessageID      string              `json:"signup_message_id,omitempty"`
    EventChannelID       string              `json:"event_channel_id,omitempty"`
    LeaderboardChannelID string              `json:"leaderboard_channel_id,omitempty"`
    UserRoleID           string              `json:"user_role_id,omitempty"`
    EditorRoleID         string              `json:"editor_role_id,omitempty"`
    AdminRoleID          string              `json:"admin_role_id,omitempty"`
    SignupEmoji          string              `json:"signup_emoji,omitempty"`
}

// NewGuildConfigFragment builds a fragment from a full GuildConfig. Nil-safe.
func NewGuildConfigFragment(cfg *guildtypes.GuildConfig) *GuildConfigFragment {
    if cfg == nil {
        return nil
    }
    return &GuildConfigFragment{
        GuildID:              cfg.GuildID,
        SignupChannelID:      cfg.SignupChannelID,
        SignupMessageID:      cfg.SignupMessageID,
        EventChannelID:       cfg.EventChannelID,
        LeaderboardChannelID: cfg.LeaderboardChannelID,
        UserRoleID:           cfg.UserRoleID,
        EditorRoleID:         cfg.EditorRoleID,
        AdminRoleID:          cfg.AdminRoleID,
        SignupEmoji:          cfg.SignupEmoji,
    }
}
