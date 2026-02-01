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

// Leaderboard represents the state of a leaderboard at a point in time.
type Leaderboard struct {
	ID              int64                           `json:"id"`
	LeaderboardData LeaderboardData                 `json:"leaderboard_data"`
	IsActive        bool                            `json:"is_active"`
	UpdateSource    sharedtypes.ServiceUpdateSource `json:"update_source"`
	UpdateID        sharedtypes.RoundID             `json:"update_id"`
	GuildID         sharedtypes.GuildID             `json:"guild_id"`
}

// User interface for leaderboard operations
type User interface {
	GetID() int64
	GetUserID() sharedtypes.DiscordID
}

func (l *Leaderboard) HasTagNumber(tagNumber sharedtypes.TagNumber) bool {
	for _, entry := range l.LeaderboardData {
		if entry.TagNumber != 0 && entry.TagNumber == tagNumber {
			return true
		}
	}
	return false
}

func (l *Leaderboard) HasUserID(userID sharedtypes.DiscordID) bool {
	for _, entry := range l.LeaderboardData {
		if entry.UserID == userID {
			return true
		}
	}
	return false
}

func (l *Leaderboard) FindEntryForUser(userID sharedtypes.DiscordID) *LeaderboardEntry {
	for i := range l.LeaderboardData {
		if l.LeaderboardData[i].UserID == userID {
			return &l.LeaderboardData[i]
		}
	}
	return nil
}
