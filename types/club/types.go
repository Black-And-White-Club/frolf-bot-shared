// Package clubtypes contains club-related domain types.
package clubtypes

// Club represents a disc golf club/league.
// Clubs are platform-agnostic and may or may not be linked to Discord.
type Club struct {
	UUID           string  `json:"uuid"`
	Name           string  `json:"name"`
	IconURL        *string `json:"icon_url,omitempty"`
	DiscordGuildID *string `json:"discord_guild_id,omitempty"`
}

// ClubInfo is a lightweight view of club data for API responses.
type ClubInfo struct {
	UUID    string  `json:"uuid"`
	Name    string  `json:"name"`
	IconURL *string `json:"icon_url,omitempty"`
}
