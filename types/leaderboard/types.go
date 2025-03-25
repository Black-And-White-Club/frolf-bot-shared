package leaderboardtypes

import (
	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
)

// LeaderboardData represents the data of a leaderboard.
type LeaderboardData map[int]string

// User interface
type User interface {
	GetID() int64
	GetUserID() sharedtypes.DiscordID // Update to use sharedtypes.DiscordID
}
