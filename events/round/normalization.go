package roundevents

import (
	"time"

	roundtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/round"
	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
)

// =============================================================================
// SCORECARD NORMALIZATION FLOW - Event Constants
// =============================================================================

// ScorecardNormalizedV1 is published when a parsed scorecard has been normalized
// into a deterministic structure suitable for ingestion.
//
// Pattern: Event Notification
// Subject: round.scorecard.normalized.v1
// Producer: backend-service (round import handler)
// Consumers: backend-service (round ingest handler)
// Triggers: ScorecardParsedForUserV1
// Version: v1 (January 2026)
const ScorecardNormalizedV1 = "round.scorecard.normalized.v1"

// ScorecardParsedForNormalizationV1 is published when raw parsing is complete
// and the data is ready to be transformed into a normalized structure.
//
// Pattern: Event Notification
// Subject: round.scorecard.parsed.for_normalization.v1
// Producer: backend-service (round parse handler)
// Consumers: backend-service (round normalization handler)
// Triggers: ScorecardNormalizedV1
// Version: v1 (January 2026)
const ScorecardParsedForNormalizationV1 = "round.scorecard.parsed.for.normalization.v1"

// =============================================================================
// SCORECARD NORMALIZATION FLOW - Payload Types
// =============================================================================

// ScorecardNormalizedPayloadV1 contains normalized scorecard data derived from
// a parsed scorecard.
//
// This event represents the final transformation step before user matching
// and score ingestion. Normalization ensures consistent handling of singles,
// doubles, and future team-based formats.
//
// Schema History:
//   - v1.0 (January 2026): Initial version
type ScorecardNormalizedPayloadV1 struct {
	// Identifiers
	ImportID string `json:"import_id"`

	GuildID   sharedtypes.GuildID   `json:"guild_id"`
	RoundID   sharedtypes.RoundID   `json:"round_id"`
	UserID    sharedtypes.DiscordID `json:"user_id"`
	ChannelID string                `json:"channel_id"`

	// Normalized domain representation of the scorecard
	Normalized roundtypes.NormalizedScorecard `json:"normalized"`

	// Discord / observability metadata
	EventMessageID string    `json:"event_message_id"`
	Timestamp      time.Time `json:"timestamp"`
}
