// Package authevents contains specific PWA magic link events.
//
// This file defines the Discord PWA Magic Link Flow - events for requesting
// and handling magic link generation for PWA dashboard access.
//
// # Flow Sequence
//
//  1. User runs /dashboard command -> MagicLinkRequestedV1
//  2. Backend generates magic link -> MagicLinkGeneratedV1
//  3. Discord bot sends DM with link
//
// # Relationship to Domain Events
//
// These Discord events are request/response events:
//   - MagicLinkRequestedV1 -> sent to backend for processing
//   - MagicLinkGeneratedV1 <- response from backend with generated link
//
// # Versioning Strategy
//
// All events include a V1 suffix for future schema evolution.
package authevents

// =============================================================================
// DISCORD PWA MAGIC LINK FLOW - Event Constants
// =============================================================================

// MagicLinkRequestedV1 is published when a user requests a PWA magic link.
//
// Pattern: Event Request
// Subject: pwa.magic-link.requested.v1
// Producer: discord-service (dashboard command handler)
// Consumers: backend (magic link generation service)
// Triggers: Magic link generation in backend
// Version: v1 (January 2025)
const MagicLinkRequestedV1 = "auth.magic-link.requested.v1"

// MagicLinkGeneratedV1 is published when the backend generates a magic link.
//
// Pattern: Event Response
// Subject: pwa.magic-link.generated.v1
// Producer: backend (magic link service)
// Consumers: discord-service (response handler)
// Triggers: Discord DM sent to user with magic link
// Version: v1 (January 2025)
const MagicLinkGeneratedV1 = "auth.magic-link.generated.v1"

// =============================================================================
// DISCORD PWA MAGIC LINK FLOW - Payload Types
// =============================================================================

// MagicLinkRequestedPayload is published when user requests a magic link.
//
// Schema History:
//   - v1.0 (January 2025): Initial version
type MagicLinkRequestedPayload struct {
	UserID        string `json:"user_id"`
	GuildID       string `json:"guild_id"`
	Role          string `json:"role"` // viewer | player | editor
	CorrelationID string `json:"correlation_id"`
}

// MagicLinkGeneratedPayload is published by backend with the generated link.
//
// Schema History:
//   - v1.0 (January 2025): Initial version
type MagicLinkGeneratedPayload struct {
	Success       bool   `json:"success"`
	URL           string `json:"url,omitempty"`
	Error         string `json:"error,omitempty"`
	UserID        string `json:"user_id"`
	GuildID       string `json:"guild_id"`
	CorrelationID string `json:"correlation_id"`
}
