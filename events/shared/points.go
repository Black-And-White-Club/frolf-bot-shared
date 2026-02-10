package sharedevents

import (
	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
)

// PointsAwardedV1 is published after round points have been calculated and persisted.
// This is a shared event consumed by display services (like Discord) to update embeds.
//
// Pattern: Event Notification
// Subject: round.points.awarded.v1
// Producer: leaderboard-service (after ProcessRound)
// Consumers: discord-service (update finalized embed with point values)
// Version: v1 (February 2026)
const PointsAwardedV1 = "round.points.awarded.v1"

// PointsAwardedPayloadV1 contains point awards for a round.
// This is used across modules to update displays after points calculation.
//
// Schema History:
//   - v1.0 (February 2026): Initial version
type PointsAwardedPayloadV1 struct {
	GuildID sharedtypes.GuildID           `json:"guild_id"`
	RoundID sharedtypes.RoundID           `json:"round_id"`
	Points  map[sharedtypes.DiscordID]int `json:"points"`
}
