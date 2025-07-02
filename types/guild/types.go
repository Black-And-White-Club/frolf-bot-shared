// In frolf-bot-shared/types/guild/config.go

package guildtypes

import (
	"time"

	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
)

type GuildConfig struct {
	GuildID              sharedtypes.GuildID
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
	// ...add any other fields you want to expose...
}
