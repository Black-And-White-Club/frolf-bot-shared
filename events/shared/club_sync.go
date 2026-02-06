package sharedevents

import sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"

// =============================================================================
// CLUB SYNC FLOW - Event Constants
// =============================================================================

// ClubSyncFromDiscordRequestedV1 is published when a user signup includes
// guild metadata (name, icon) that should be synced to the club module.
//
// This is a cross-module event: the user module publishes it after processing
// a signup request, and the club module consumes it to upsert club info.
//
// Pattern: Event Notification
// Subject: club.sync.from.discord.requested.v1
// Producer: user-service (signup handler)
// Consumers: club-service (club sync handler)
// Version: v1 (February 2026)
const ClubSyncFromDiscordRequestedV1 = "club.sync.from.discord.requested.v1"

// =============================================================================
// CLUB SYNC FLOW - Payload Types
// =============================================================================

// ClubSyncFromDiscordRequestedPayloadV1 contains the guild metadata needed
// to create or update a club record from Discord guild information.
//
// Schema History:
//   - v1.0 (February 2026): Initial version
type ClubSyncFromDiscordRequestedPayloadV1 struct {
	GuildID   sharedtypes.GuildID `json:"guild_id"`
	GuildName string              `json:"guild_name"`
	IconURL   *string             `json:"icon_url,omitempty"`
}
