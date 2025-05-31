package leaderboardtypes

import (
	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
)

// Core leaderboard domain types
type LeaderboardEntry struct {
	TagNumber sharedtypes.TagNumber `json:"tag_number"`
	UserID    sharedtypes.DiscordID `json:"user_id"`
}

type LeaderboardData []LeaderboardEntry

// Event-specific types (used in event payloads)
type TagAssignmentInfo struct {
	UserID    sharedtypes.DiscordID `json:"user_id"`
	TagNumber sharedtypes.TagNumber `json:"tag_number"`
}

// User interface for leaderboard operations
type User interface {
	GetID() int64
	GetUserID() sharedtypes.DiscordID
}
