package userevents

import (
	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
	"github.com/google/uuid"
)

const (
	// UserProfileUpdatedV1 is emitted when Discord profile data is observed
	UserProfileUpdatedV1 = "user.profile.updated.v1"
)

// UserProfileUpdatedPayloadV1 contains Discord user profile data.
// Emitted by discord-frolf-bot when it observes user interactions.
//
// Schema History:
//   - v1.0 (January 2026): Initial version
type UserProfileUpdatedPayloadV1 struct {
	UserID      sharedtypes.DiscordID `json:"user_id"`
	UserUUID    *uuid.UUID            `json:"user_uuid,omitempty"`
	GuildID     sharedtypes.GuildID   `json:"guild_id"`
	ClubUUID    *uuid.UUID            `json:"club_uuid,omitempty"`
	Username    string                `json:"username"`     // Discord username
	DisplayName string                `json:"display_name"` // Nickname in guild, or username if no nickname
	AvatarHash  string                `json:"avatar_hash"`  // Discord avatar hash (empty if default)
}

const (
	// UserProfileSyncRequestTopicV1 is used to request a profile sync from discord-bot.
	UserProfileSyncRequestTopicV1 = "user.profile.sync.request.v1"
)

// UserProfileSyncRequestPayloadV1 contains data for syncing a user profile.
type UserProfileSyncRequestPayloadV1 struct {
	UserID  sharedtypes.DiscordID `json:"user_id"`
	GuildID sharedtypes.GuildID   `json:"guild_id"`
}
