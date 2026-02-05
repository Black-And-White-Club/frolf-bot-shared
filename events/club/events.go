// Package clubevents contains club-related domain events.
//
// This file defines Club Info Flow - events for retrieving and updating
// club information (name, icon).
//
// # Flow Sequences
//
// ## Club Info Retrieval (Request-Reply)
//  1. Request -> ClubInfoRequestV1
//  2. Response -> ClubInfoResponseV1
//
// ## Club Updated (Broadcast)
//  1. Notification -> ClubUpdatedV1
//
// # Versioning Strategy
//
// All events include a V1 suffix for future schema evolution.
package clubevents

// =============================================================================
// CLUB INFO FLOW - Event Constants
// =============================================================================

// ClubInfoRequestV1 is published when club info is requested.
//
// Pattern: Request-Reply
// Subject: club.info.request.v1.{club_uuid}
// Producer: PWA (via NATS WebSocket), any service needing club info
// Consumers: club-service
// Version: v1 (February 2026)
const ClubInfoRequestV1 = "club.info.request.v1"

// ClubInfoResponseV1 is the reply containing club info.
//
// Pattern: Request-Reply Response
// Subject: _INBOX.{reply}
// Producer: club-service
// Consumers: requesting service/PWA
// Version: v1 (February 2026)
const ClubInfoResponseV1 = "club.info.response.v1"

// ClubUpdatedV1 is published when club details change (name, icon).
//
// Pattern: Event Notification (Broadcast)
// Subject: club.updated.v1.{club_uuid}
// Producer: club-service (after setup or admin update)
// Consumers: PWA (cache invalidation), any interested services
// Version: v1 (February 2026)
const ClubUpdatedV1 = "club.updated.v1"

// =============================================================================
// CLUB INFO FLOW - Payload Types
// =============================================================================

// ClubInfoRequestPayloadV1 contains club info request data.
//
// Schema History:
//   - v1.0 (February 2026): Initial version
type ClubInfoRequestPayloadV1 struct {
	ClubUUID string `json:"club_uuid"`
}

// ClubInfoResponsePayloadV1 contains club info response data.
//
// Schema History:
//   - v1.0 (February 2026): Initial version
type ClubInfoResponsePayloadV1 struct {
	UUID           string  `json:"uuid"`
	Name           string  `json:"name"`
	IconURL        *string `json:"icon_url,omitempty"`
	DiscordGuildID *string `json:"discord_guild_id,omitempty"`
}

// ClubUpdatedPayloadV1 contains updated club data.
//
// Schema History:
//   - v1.0 (February 2026): Initial version
type ClubUpdatedPayloadV1 struct {
	UUID    string  `json:"uuid"`
	Name    string  `json:"name"`
	IconURL *string `json:"icon_url,omitempty"`
}
