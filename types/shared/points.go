package sharedtypes

// PointsAwardedPayloadV1 contains point awards for a round.
// This is used across modules to update displays after points calculation.
//
// Schema History:
//   - v1.0 (February 2026): Initial version (moved from leaderboard domain)
type PointsAwardedPayloadV1 struct {
	GuildID GuildID           `json:"guild_id"`
	RoundID RoundID           `json:"round_id"`
	Points  map[DiscordID]int `json:"points"`
}
