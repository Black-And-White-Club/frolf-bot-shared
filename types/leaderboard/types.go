package leaderboardtypes

import (
	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
)

// LeaderboardData represents the data of a leaderboard.
type LeaderboardData []LeaderboardEntry

// User interface
type User interface {
	GetID() int64
	GetUserID() sharedtypes.DiscordID // Update to use sharedtypes.DiscordID
}

type LeaderboardEntry struct {
	TagNumber sharedtypes.TagNumber `json:"tag_number"`
	UserID    sharedtypes.DiscordID `json:"user_id"`
}

func (e LeaderboardEntry) IsValid() bool {
	return e.TagNumber > 0 && e.UserID != ""
}
