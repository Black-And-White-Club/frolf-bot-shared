// Package roundevents contains all round-related domain events.
//
// This file defines the Scorecard Import Flow - events related to
// importing scores from UDisc scorecards.
//
// # Flow Sequence
//
//  1. User uploads scorecard -> ScorecardUploadedV1
//  2. OR User provides URL -> ScorecardURLRequestedV1
//  3. Parse request sent -> ScorecardParseRequestedV1
//  4. Scorecard parsed -> ScorecardParsedV1
//  5. OR Parse fails -> ScorecardParseFailedV1
//  6. Conflict detected (optional) -> ImportConflictDetectedV1
//  7. Overwrite confirmed (optional) -> ImportOverwriteConfirmedV1
//  8. Scores imported -> RoundScoresImportedV1
//  9. Participant auto-added (optional) -> RoundParticipantAutoAddedV1
// 10. Import completed -> ImportCompletedV1
// 11. Scores finalized -> RoundScoresFinalizedV1
// 12. OR Import fails -> ImportFailedV1
//
// # Pattern Reference
//
// This flow follows the Event Notification pattern (Martin Fowler).
//
// # Versioning Strategy
//
// All events include a V1 suffix for future schema evolution.
package roundevents

import (
	"time"

	roundtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/round"
	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
)

// =============================================================================
// SCORECARD IMPORT FLOW - Event Constants
// =============================================================================

// -----------------------------------------------------------------------------
// Upload Events
// -----------------------------------------------------------------------------

// ScorecardUploadedV1 is published when a user uploads a scorecard file.
//
// Pattern: Event Notification
// Subject: round.scorecard.uploaded.v1
// Producer: discord-service
// Consumers: backend-service (parse handler)
// Triggers: ScorecardParseRequestedV1
// Version: v1 (December 2024)
const ScorecardUploadedV1 = "round.scorecard.uploaded.v1"

// ScorecardURLRequestedV1 is published when a user provides a UDisc URL.
//
// Pattern: Event Notification
// Subject: round.scorecard.url.requested.v1
// Producer: discord-service
// Consumers: backend-service (URL fetch handler)
// Triggers: ScorecardParseRequestedV1
// Version: v1 (December 2024)
const ScorecardURLRequestedV1 = "round.scorecard.url.requested.v1"

// -----------------------------------------------------------------------------
// Parse Events
// -----------------------------------------------------------------------------

// ScorecardParseRequestedV1 is published to request scorecard parsing.
//
// Pattern: Event Notification
// Subject: round.scorecard.parse.requested.v1
// Producer: backend-service (upload handler)
// Consumers: backend-service (parse service)
// Triggers: ScorecardParsedV1 OR ScorecardParseFailedV1
// Version: v1 (December 2024)
const ScorecardParseRequestedV1 = "round.scorecard.parse.requested.v1"

// ScorecardParsedV1 is published when a scorecard is successfully parsed.
//
// Pattern: Event Notification
// Subject: round.scorecard.parsed.v1
// Producer: backend-service (parse service)
// Consumers: backend-service (import handler), user module
// Triggers: Import flow continuation
// Version: v1 (December 2024)
const ScorecardParsedV1 = "round.scorecard.parsed.v1"

// ScorecardParsedForUserV1 is a fan-out copy of ScorecardParsedV1 for the user module.
// We intentionally use a different subject so the backend can keep a single durable
// consumer per subject (queue semantics) while still letting both modules receive
// the same parsed payload.
//
// Pattern: Event Notification
// Subject: round.scorecard.parsed.user.v1
// Producer: backend-service (parse service)
// Consumers: backend-service (user module)
// Version: v1 (December 2024)
const ScorecardParsedForUserV1 = "round.scorecard.parsed.user.v1"

// ScorecardParseFailedV1 is published when scorecard parsing fails.
//
// Pattern: Event Notification
// Subject: round.scorecard.parse.failed.v1
// Producer: backend-service (parse service)
// Consumers: discord-service (error handler)
// Version: v1 (December 2024)
const ScorecardParseFailedV1 = "round.scorecard.parse.failed.v1"

// -----------------------------------------------------------------------------
// Import Events
// -----------------------------------------------------------------------------

// ImportConflictDetectedV1 is published when an import conflict is detected.
//
// Pattern: Event Notification
// Subject: round.import.conflict.detected.v1
// Producer: backend-service (import handler)
// Consumers: discord-service (conflict resolution UI)
// Version: v1 (December 2024)
const ImportConflictDetectedV1 = "round.import.conflict.detected.v1"

// ImportOverwriteConfirmedV1 is published when user confirms overwrite.
//
// Pattern: Event Notification
// Subject: round.import.overwrite.confirmed.v1
// Producer: discord-service
// Consumers: backend-service (import handler)
// Triggers: Continue import flow
// Version: v1 (December 2024)
const ImportOverwriteConfirmedV1 = "round.import.overwrite.confirmed.v1"

// ImportCompletedV1 is published when import completes successfully.
//
// Pattern: Event Notification
// Subject: round.import.completed.v1
// Producer: backend-service (import handler)
// Consumers: discord-service (success notification)
// Version: v1 (December 2024)
const ImportCompletedV1 = "round.import.completed.v1"

// ImportFailedV1 is published when import fails.
//
// Pattern: Event Notification
// Subject: round.import.failed.v1
// Producer: backend-service (import handler)
// Consumers: discord-service (error handler)
// Version: v1 (December 2024)
const ImportFailedV1 = "round.import.failed.v1"

// RoundParticipantAutoAddedV1 is published when a player is auto-added during import.
//
// Pattern: Event Notification
// Subject: round.participant.auto.added.v1
// Producer: backend-service (import handler)
// Consumers: discord-service (notification)
// Version: v1 (December 2024)
const RoundParticipantAutoAddedV1 = "round.participant.auto.added.v1"

// RoundScoresImportedV1 is published when scores are imported into a round.
//
// Pattern: Event Notification
// Subject: round.scores.imported.v1
// Producer: backend-service (import handler)
// Consumers: backend-service (score processing)
// Version: v1 (December 2024)
const RoundScoresImportedV1 = "round.scores.imported.v1"

// RoundScoresFinalizedV1 is published when imported scores are finalized.
//
// Pattern: Event Notification
// Subject: round.scores.finalized.v1
// Producer: backend-service (import handler)
// Consumers: discord-service (embed update)
// Version: v1 (December 2024)
const RoundScoresFinalizedV1 = "round.scores.finalized.v1"

// =============================================================================
// SCORECARD IMPORT FLOW - Payload Types
// =============================================================================

// -----------------------------------------------------------------------------
// Upload Payloads
// -----------------------------------------------------------------------------

// ScorecardUploadedPayloadV1 contains uploaded scorecard data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type ScorecardUploadedPayloadV1 struct {
	GuildID   sharedtypes.GuildID   `json:"guild_id"`
	RoundID   sharedtypes.RoundID   `json:"round_id"`
	ImportID  string                `json:"import_id"`
	UserID    sharedtypes.DiscordID `json:"user_id"`
	ChannelID string                `json:"channel_id"`
	MessageID string                `json:"message_id"`
	FileData  []byte                `json:"file_data,omitempty"`
	FileURL   string                `json:"file_url,omitempty"`
	FileName  string                `json:"file_name,omitempty"`
	UDiscURL  string                `json:"udisc_url,omitempty"`
	Notes     string                `json:"notes,omitempty"`
	Timestamp time.Time             `json:"timestamp"`
}

// ScorecardURLRequestedPayloadV1 contains URL import request data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type ScorecardURLRequestedPayloadV1 struct {
	GuildID   sharedtypes.GuildID   `json:"guild_id"`
	RoundID   sharedtypes.RoundID   `json:"round_id"`
	ImportID  string                `json:"import_id"`
	UserID    sharedtypes.DiscordID `json:"user_id"`
	ChannelID string                `json:"channel_id"`
	MessageID string                `json:"message_id"`
	UDiscURL  string                `json:"udisc_url"`
	Notes     string                `json:"notes,omitempty"`
	Timestamp time.Time             `json:"timestamp"`
}

// -----------------------------------------------------------------------------
// Parse Payloads
// -----------------------------------------------------------------------------

// ScorecardParseFailedPayloadV1 contains parse failure data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type ScorecardParseFailedPayloadV1 struct {
	GuildID        sharedtypes.GuildID   `json:"guild_id"`
	RoundID        sharedtypes.RoundID   `json:"round_id"`
	ImportID       string                `json:"import_id"`
	EventMessageID string                `json:"event_message_id"`
	UserID         sharedtypes.DiscordID `json:"user_id"`
	ChannelID      string                `json:"channel_id"`
	Error          string                `json:"error"`
	Timestamp      time.Time             `json:"timestamp"`
}

// ParsedScorecardPayloadV1 contains successfully parsed scorecard data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type ParsedScorecardPayloadV1 struct {
	GuildID        sharedtypes.GuildID         `json:"guild_id"`
	RoundID        sharedtypes.RoundID         `json:"round_id"`
	ImportID       string                      `json:"import_id"`
	EventMessageID string                      `json:"event_message_id"`
	UserID         sharedtypes.DiscordID       `json:"user_id"`
	ChannelID      string                      `json:"channel_id"`
	ParsedData     *roundtypes.ParsedScorecard `json:"parsed_data"`
	Timestamp      time.Time                   `json:"timestamp"`
}

// -----------------------------------------------------------------------------
// Import Payloads
// -----------------------------------------------------------------------------

// ImportFailedPayloadV1 contains import failure data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type ImportFailedPayloadV1 struct {
	GuildID        sharedtypes.GuildID   `json:"guild_id"`
	RoundID        sharedtypes.RoundID   `json:"round_id"`
	ImportID       string                `json:"import_id"`
	EventMessageID string                `json:"event_message_id"`
	UserID         sharedtypes.DiscordID `json:"user_id"`
	ChannelID      string                `json:"channel_id"`
	Error          string                `json:"error"`
	ErrorCode      string                `json:"error_code"`
	Timestamp      time.Time             `json:"timestamp"`
}

// ImportCompletedPayloadV1 contains successful import data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type ImportCompletedPayloadV1 struct {
	GuildID            sharedtypes.GuildID        `json:"guild_id"`
	RoundID            sharedtypes.RoundID        `json:"round_id"`
	ImportID           string                     `json:"import_id"`
	EventMessageID     string                     `json:"event_message_id"`
	UserID             sharedtypes.DiscordID      `json:"user_id"`
	ChannelID          string                     `json:"channel_id"`
	ScoresIngested     int                        `json:"scores_ingested"`
	MatchedPlayers     int                        `json:"matched_players"`
	UnmatchedPlayers   int                        `json:"unmatched_players"`
	PlayersAutoAdded   int                        `json:"players_auto_added"`
	MatchedPlayersList []roundtypes.MatchedPlayer `json:"matched_players_list,omitempty"`
	SkippedPlayers     []string                   `json:"skipped_players,omitempty"`
	AutoAddedUserIDs   []sharedtypes.DiscordID    `json:"auto_added_user_ids,omitempty"`
	Scores             []sharedtypes.ScoreInfo    `json:"scores,omitempty"`
	Timestamp          time.Time                  `json:"timestamp"`
}

// RoundParticipantAutoAddedPayloadV1 contains auto-add event data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoundParticipantAutoAddedPayloadV1 struct {
	GuildID   sharedtypes.GuildID   `json:"guild_id"`
	RoundID   sharedtypes.RoundID   `json:"round_id"`
	ImportID  string                `json:"import_id"`
	UserID    sharedtypes.DiscordID `json:"user_id"`
	ChannelID string                `json:"channel_id"`
	AddedUser sharedtypes.DiscordID `json:"added_user"`
	Timestamp time.Time             `json:"timestamp"`
}

// RoundScoresImportedPayloadV1 contains imported scores data.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type RoundScoresImportedPayloadV1 struct {
	GuildID   sharedtypes.GuildID   `json:"guild_id"`
	RoundID   sharedtypes.RoundID   `json:"round_id"`
	ImportID  string                `json:"import_id"`
	UserID    sharedtypes.DiscordID `json:"user_id"`
	ChannelID string                `json:"channel_id"`
	Count     int                   `json:"count"`
	Timestamp time.Time             `json:"timestamp"`
}

// ImportScoresAppliedPayloadV1 is emitted after imported scores are applied.
//
// Schema History:
//   - v1.0 (December 2024): Initial version
type ImportScoresAppliedPayloadV1 struct {
	GuildID        sharedtypes.GuildID      `json:"guild_id"`
	RoundID        sharedtypes.RoundID      `json:"round_id"`
	ImportID       string                   `json:"import_id"`
	Participants   []roundtypes.Participant `json:"participants"`
	EventMessageID string                   `json:"event_message_id"`
	Timestamp      time.Time                `json:"timestamp"`
}
