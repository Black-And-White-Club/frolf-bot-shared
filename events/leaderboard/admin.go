// Package leaderboardevents contains leaderboard-related domain events.
//
// This file defines the Leaderboard Admin Flow - events for administrative
// operations such as point history, manual adjustments, round recalculation,
// and season management.
//
// # Flow Sequences
//
// ## Point History Flow
//  1. Request -> LeaderboardPointHistoryRequestedV1
//  2. Response -> LeaderboardPointHistoryResponseV1
//  3. OR Failure -> LeaderboardPointHistoryFailedV1
//
// ## Manual Point Adjustment Flow
//  1. Request -> LeaderboardManualPointAdjustmentV1
//  2. Success -> LeaderboardManualPointAdjustmentSuccessV1
//  3. OR Failure -> LeaderboardManualPointAdjustmentFailedV1
//
// ## Recalculate Round Flow
//  1. Request -> LeaderboardRecalculateRoundV1
//  2. Success -> LeaderboardRecalculateRoundSuccessV1
//  3. OR Failure -> LeaderboardRecalculateRoundFailedV1
//
// ## New Season Flow
//  1. Request -> LeaderboardStartNewSeasonV1
//  2. Success -> LeaderboardStartNewSeasonSuccessV1
//  3. OR Failure -> LeaderboardStartNewSeasonFailedV1
//
// ## Season Standings Flow
//  1. Request -> LeaderboardGetSeasonStandingsV1
//  2. Response -> LeaderboardGetSeasonStandingsResponseV1
//  3. OR Failure -> LeaderboardGetSeasonStandingsFailedV1
//
// # Versioning Strategy
//
// All events include a V1 suffix for future schema evolution.
// Current admin flow remains V1-only during migration.
package leaderboardevents

import sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"

// =============================================================================
// POINT HISTORY FLOW - Event Constants
// =============================================================================

// LeaderboardPointHistoryRequestedV1 is published when point history is requested.
//
// Pattern: Event Notification
// Subject: leaderboard.point.history.requested.v1
// Producer: discord-service (admin command)
// Consumers: leaderboard-service (history handler)
// Triggers: LeaderboardPointHistoryResponseV1 OR LeaderboardPointHistoryFailedV1
// Version: v1 (February 2026)
const LeaderboardPointHistoryRequestedV1 = "leaderboard.point.history.requested.v1"

// LeaderboardPointHistoryResponseV1 is published with the point history data.
//
// Pattern: Event Notification
// Subject: leaderboard.point.history.response.v1
// Producer: leaderboard-service
// Consumers: discord-service (response handler)
// Version: v1 (February 2026)
const LeaderboardPointHistoryResponseV1 = "leaderboard.point.history.response.v1"

// LeaderboardPointHistoryFailedV1 is published when point history retrieval fails.
//
// Pattern: Event Notification
// Subject: leaderboard.point.history.failed.v1
// Producer: leaderboard-service
// Consumers: discord-service (error handler)
// Version: v1 (February 2026)
const LeaderboardPointHistoryFailedV1 = "leaderboard.point.history.failed.v1"

// =============================================================================
// MANUAL POINT ADJUSTMENT FLOW - Event Constants
// =============================================================================

// LeaderboardManualPointAdjustmentV1 is published when a manual point adjustment is requested.
//
// Pattern: Event Notification
// Subject: leaderboard.manual.point.adjustment.v1
// Producer: discord-service (admin command)
// Consumers: leaderboard-service (adjustment handler)
// Triggers: LeaderboardManualPointAdjustmentSuccessV1 OR LeaderboardManualPointAdjustmentFailedV1
// Version: v1 (February 2026)
const LeaderboardManualPointAdjustmentV1 = "leaderboard.manual.point.adjustment.v1"

// LeaderboardManualPointAdjustmentSuccessV1 is published when a manual point adjustment succeeds.
//
// Pattern: Event Notification
// Subject: leaderboard.manual.point.adjustment.success.v1
// Producer: leaderboard-service
// Consumers: discord-service (confirmation)
// Version: v1 (February 2026)
const LeaderboardManualPointAdjustmentSuccessV1 = "leaderboard.manual.point.adjustment.success.v1"

// LeaderboardManualPointAdjustmentFailedV1 is published when a manual point adjustment fails.
//
// Pattern: Event Notification
// Subject: leaderboard.manual.point.adjustment.failed.v1
// Producer: leaderboard-service
// Consumers: discord-service (error handler)
// Version: v1 (February 2026)
const LeaderboardManualPointAdjustmentFailedV1 = "leaderboard.manual.point.adjustment.failed.v1"

// =============================================================================
// RECALCULATE ROUND FLOW - Event Constants
// =============================================================================

// LeaderboardRecalculateRoundV1 is published when round recalculation is requested.
//
// Pattern: Event Notification
// Subject: leaderboard.recalculate.round.v1
// Producer: discord-service (admin command)
// Consumers: leaderboard-service (recalculate handler)
// Triggers: LeaderboardRecalculateRoundSuccessV1 OR LeaderboardRecalculateRoundFailedV1
// Version: v1 (February 2026)
const LeaderboardRecalculateRoundV1 = "leaderboard.recalculate.round.v1"

// LeaderboardRecalculateRoundSuccessV1 is published when round recalculation succeeds.
//
// Pattern: Event Notification
// Subject: leaderboard.recalculate.round.success.v1
// Producer: leaderboard-service
// Consumers: discord-service (confirmation)
// Version: v1 (February 2026)
const LeaderboardRecalculateRoundSuccessV1 = "leaderboard.recalculate.round.success.v1"

// LeaderboardRecalculateRoundFailedV1 is published when round recalculation fails.
//
// Pattern: Event Notification
// Subject: leaderboard.recalculate.round.failed.v1
// Producer: leaderboard-service
// Consumers: discord-service (error handler)
// Version: v1 (February 2026)
const LeaderboardRecalculateRoundFailedV1 = "leaderboard.recalculate.round.failed.v1"

// =============================================================================
// NEW SEASON FLOW - Event Constants
// =============================================================================

// LeaderboardStartNewSeasonV1 is published when a new season start is requested.
//
// Pattern: Event Notification
// Subject: leaderboard.start.new.season.v1
// Producer: discord-service (admin command)
// Consumers: leaderboard-service (season handler)
// Triggers: LeaderboardStartNewSeasonSuccessV1 OR LeaderboardStartNewSeasonFailedV1
// Version: v1 (February 2026)
const LeaderboardStartNewSeasonV1 = "leaderboard.start.new.season.v1"

// LeaderboardStartNewSeasonSuccessV1 is published when a new season start succeeds.
//
// Pattern: Event Notification
// Subject: leaderboard.start.new.season.success.v1
// Producer: leaderboard-service
// Consumers: discord-service (confirmation)
// Version: v1 (February 2026)
const LeaderboardStartNewSeasonSuccessV1 = "leaderboard.start.new.season.success.v1"

// LeaderboardStartNewSeasonFailedV1 is published when a new season start fails.
//
// Pattern: Event Notification
// Subject: leaderboard.start.new.season.failed.v1
// Producer: leaderboard-service
// Consumers: discord-service (error handler)
// Version: v1 (February 2026)
const LeaderboardStartNewSeasonFailedV1 = "leaderboard.start.new.season.failed.v1"

// LeaderboardEndSeasonV1 is published when a season end is requested.
//
// Pattern: Event Notification
// Subject: leaderboard.end.season.v1
// Producer: discord-service (admin command)
// Consumers: leaderboard-service (season handler)
// Triggers: LeaderboardEndSeasonSuccessV1 OR LeaderboardEndSeasonFailedV1
// Version: v1 (February 2026)
const LeaderboardEndSeasonV1 = "leaderboard.end.season.v1"

// LeaderboardEndSeasonSuccessV1 is published when a season end succeeds.
//
// Pattern: Event Notification
// Subject: leaderboard.end.season.success.v1
// Producer: leaderboard-service
// Consumers: discord-service (confirmation)
// Version: v1 (February 2026)
const LeaderboardEndSeasonSuccessV1 = "leaderboard.end.season.success.v1"

// LeaderboardEndSeasonFailedV1 is published when a season end fails.
//
// Pattern: Event Notification
// Subject: leaderboard.end.season.failed.v1
// Producer: leaderboard-service
// Consumers: discord-service (error handler)
// Version: v1 (February 2026)
const LeaderboardEndSeasonFailedV1 = "leaderboard.end.season.failed.v1"

// =============================================================================
// SEASON STANDINGS FLOW - Event Constants
// =============================================================================

// LeaderboardGetSeasonStandingsV1 is published when season standings are requested.
//
// Pattern: Event Notification
// Subject: leaderboard.get.season.standings.v1
// Producer: discord-service (admin command)
// Consumers: leaderboard-service (standings handler)
// Triggers: LeaderboardGetSeasonStandingsResponseV1 OR LeaderboardGetSeasonStandingsFailedV1
// Version: v1 (February 2026)
const LeaderboardGetSeasonStandingsV1 = "leaderboard.get.season.standings.v1"

// LeaderboardGetSeasonStandingsResponseV1 is published with season standings data.
//
// Pattern: Event Notification
// Subject: leaderboard.get.season.standings.response.v1
// Producer: leaderboard-service
// Consumers: discord-service (response handler)
// Version: v1 (February 2026)
const LeaderboardGetSeasonStandingsResponseV1 = "leaderboard.get.season.standings.response.v1"

// LeaderboardGetSeasonStandingsFailedV1 is published when season standings retrieval fails.
//
// Pattern: Event Notification
// Subject: leaderboard.get.season.standings.failed.v1
// Producer: leaderboard-service
// Consumers: discord-service (error handler)
// Version: v1 (February 2026)
const LeaderboardGetSeasonStandingsFailedV1 = "leaderboard.get.season.standings.failed.v1"

// =============================================================================
// LIST SEASONS REQUEST-REPLY FLOW - Event Constants
// =============================================================================

// LeaderboardListSeasonsRequestedV1 is the wildcard subject for listing seasons
// via NATS request-reply. The PWA publishes to this subject with a reply_to inbox.
//
// Pattern: Request-Reply (wildcard)
// Subject: leaderboard.seasons.list.request.v1.>
// Producer: PWA, discord-service
// Consumers: leaderboard-service (list handler)
// Triggers: LeaderboardListSeasonsResponseV1 OR LeaderboardListSeasonsFailedV1
// Version: v1 (February 2026)
const LeaderboardListSeasonsRequestedV1 = "leaderboard.seasons.list.request.v1"

// LeaderboardListSeasonsResponseV1 is published as a reply with the seasons list.
//
// Pattern: Request-Reply response
// Subject: (reply_to inbox)
// Producer: leaderboard-service
// Consumers: PWA, discord-service
// Version: v1 (February 2026)
const LeaderboardListSeasonsResponseV1 = "leaderboard.seasons.list.response.v1"

// LeaderboardListSeasonsFailedV1 is published when listing seasons fails.
//
// Pattern: Request-Reply response
// Subject: (reply_to inbox)
// Producer: leaderboard-service
// Consumers: PWA, discord-service
// Version: v1 (February 2026)
const LeaderboardListSeasonsFailedV1 = "leaderboard.seasons.list.failed.v1"

// =============================================================================
// SEASON STANDINGS REQUEST-REPLY FLOW - Event Constants
// =============================================================================

// LeaderboardSeasonStandingsRequestV1 is the wildcard subject for retrieving
// season standings via NATS request-reply. Separate from the event-driven
// GetSeasonStandings flow — this is for PWA historical season browsing.
//
// Pattern: Request-Reply (wildcard)
// Subject: leaderboard.season.standings.request.v1.>
// Producer: PWA, discord-service
// Consumers: leaderboard-service (standings handler)
// Triggers: LeaderboardSeasonStandingsResponseV1 OR LeaderboardSeasonStandingsFailedV1
// Version: v1 (February 2026)
const LeaderboardSeasonStandingsRequestV1 = "leaderboard.season.standings.request.v1"

// LeaderboardSeasonStandingsResponseV1 is published as a reply with the standings.
//
// Pattern: Request-Reply response
// Subject: (reply_to inbox)
// Producer: leaderboard-service
// Consumers: PWA, discord-service
// Version: v1 (February 2026)
const LeaderboardSeasonStandingsResponseV1 = "leaderboard.season.standings.response.v1"

// LeaderboardSeasonStandingsFailedV1 is published when season standings retrieval fails.
//
// Pattern: Request-Reply response
// Subject: (reply_to inbox)
// Producer: leaderboard-service
// Consumers: PWA, discord-service
// Version: v1 (February 2026)
const LeaderboardSeasonStandingsFailedV1 = "leaderboard.season.standings.failed.v1"

// =============================================================================
// ADMIN EVENT PAYLOAD TYPES
// =============================================================================

// -----------------------------------------------------------------------------
// Point History Payloads
// -----------------------------------------------------------------------------

// PointHistoryRequestedPayloadV1 requests point history for a member.
//
// Schema History:
//   - v1.0 (February 2026): Initial version
type PointHistoryRequestedPayloadV1 struct {
	GuildID  sharedtypes.GuildID   `json:"guild_id"`
	MemberID sharedtypes.DiscordID `json:"member_id"`
	Limit    int                   `json:"limit,omitempty"`
}

// PointHistoryResponsePayloadV1 contains the point history response.
//
// Schema History:
//   - v1.0 (February 2026): Initial version
type PointHistoryResponsePayloadV1 struct {
	GuildID  sharedtypes.GuildID   `json:"guild_id"`
	MemberID sharedtypes.DiscordID `json:"member_id"`
	History  []PointHistoryItemV1  `json:"history"`
}

// PointHistoryItemV1 represents a single point history entry.
type PointHistoryItemV1 struct {
	RoundID   sharedtypes.RoundID `json:"round_id"`
	SeasonID  string              `json:"season_id"`
	Points    int                 `json:"points"`
	Reason    string              `json:"reason"`
	Tier      string              `json:"tier"`
	Opponents int                 `json:"opponents"`
	CreatedAt string              `json:"created_at"`
}

// -----------------------------------------------------------------------------
// Manual Point Adjustment Payloads
// -----------------------------------------------------------------------------

// ManualPointAdjustmentPayloadV1 requests a manual point adjustment.
//
// Schema History:
//   - v1.0 (February 2026): Initial version
type ManualPointAdjustmentPayloadV1 struct {
	GuildID     sharedtypes.GuildID   `json:"guild_id"`
	MemberID    sharedtypes.DiscordID `json:"member_id"`
	PointsDelta int                   `json:"points_delta"`
	Reason      string                `json:"reason"`
	AdminID     sharedtypes.DiscordID `json:"admin_id"`
}

// ManualPointAdjustmentSuccessPayloadV1 confirms a point adjustment.
//
// Schema History:
//   - v1.0 (February 2026): Initial version
type ManualPointAdjustmentSuccessPayloadV1 struct {
	GuildID     sharedtypes.GuildID   `json:"guild_id"`
	MemberID    sharedtypes.DiscordID `json:"member_id"`
	PointsDelta int                   `json:"points_delta"`
	Reason      string                `json:"reason"`
}

// -----------------------------------------------------------------------------
// Recalculate Round Payloads
// -----------------------------------------------------------------------------

// RecalculateRoundPayloadV1 requests recalculation of a round.
//
// Schema History:
//   - v1.0 (February 2026): Initial version
type RecalculateRoundPayloadV1 struct {
	GuildID sharedtypes.GuildID `json:"guild_id"`
	RoundID sharedtypes.RoundID `json:"round_id"`
}

// RecalculateRoundSuccessPayloadV1 confirms round recalculation.
//
// Schema History:
//   - v1.0 (February 2026): Initial version
type RecalculateRoundSuccessPayloadV1 struct {
	GuildID       sharedtypes.GuildID           `json:"guild_id"`
	RoundID       sharedtypes.RoundID           `json:"round_id"`
	PointsAwarded map[sharedtypes.DiscordID]int `json:"points_awarded"`
}

// -----------------------------------------------------------------------------
// New Season Payloads
// -----------------------------------------------------------------------------

// StartNewSeasonPayloadV1 requests a new season to be started.
//
// Schema History:
//   - v1.0 (February 2026): Initial version
type StartNewSeasonPayloadV1 struct {
	GuildID    sharedtypes.GuildID `json:"guild_id"`
	SeasonID   string              `json:"season_id"`
	SeasonName string              `json:"season_name"`
}

// StartNewSeasonSuccessPayloadV1 confirms a new season was started.
//
// Schema History:
//   - v1.0 (February 2026): Initial version
type StartNewSeasonSuccessPayloadV1 struct {
	GuildID    sharedtypes.GuildID `json:"guild_id"`
	SeasonID   string              `json:"season_id"`
	SeasonName string              `json:"season_name"`
}

// EndSeasonPayloadV1 requests the current season to be ended.
//
// Schema History:
//   - v1.0 (February 2026): Initial version
type EndSeasonPayloadV1 struct {
	GuildID sharedtypes.GuildID `json:"guild_id"`
}

// EndSeasonSuccessPayloadV1 confirms the season was ended.
//
// Schema History:
//   - v1.0 (February 2026): Initial version
type EndSeasonSuccessPayloadV1 struct {
	GuildID sharedtypes.GuildID `json:"guild_id"`
}

// -----------------------------------------------------------------------------
// Season Standings Payloads
// -----------------------------------------------------------------------------

// GetSeasonStandingsPayloadV1 requests standings for a season.
//
// Schema History:
//   - v1.0 (February 2026): Initial version
type GetSeasonStandingsPayloadV1 struct {
	GuildID  sharedtypes.GuildID `json:"guild_id"`
	SeasonID string              `json:"season_id"`
}

// GetSeasonStandingsResponsePayloadV1 contains the standings response.
//
// Schema History:
//   - v1.0 (February 2026): Initial version
type GetSeasonStandingsResponsePayloadV1 struct {
	GuildID   sharedtypes.GuildID    `json:"guild_id"`
	SeasonID  string                 `json:"season_id"`
	Standings []SeasonStandingItemV1 `json:"standings"`
}

// SeasonStandingItemV1 represents a single season standing entry.
type SeasonStandingItemV1 struct {
	MemberID      sharedtypes.DiscordID `json:"member_id"`
	TotalPoints   int                   `json:"total_points"`
	CurrentTier   string                `json:"current_tier"`
	SeasonBestTag int                   `json:"season_best_tag"`
	RoundsPlayed  int                   `json:"rounds_played"`
}

// -----------------------------------------------------------------------------
// List Seasons Request-Reply Payloads
// -----------------------------------------------------------------------------

// ListSeasonsRequestPayloadV1 requests the list of seasons for a guild.
// Used by the PWA and Discord bot via NATS request-reply.
//
// Schema History:
//   - v1.0 (February 2026): Initial version
type ListSeasonsRequestPayloadV1 struct {
	GuildID sharedtypes.GuildID `json:"guild_id"`
}

// ListSeasonsResponsePayloadV1 contains the seasons list response.
//
// Schema History:
//   - v1.0 (February 2026): Initial version
type ListSeasonsResponsePayloadV1 struct {
	GuildID sharedtypes.GuildID `json:"guild_id"`
	Seasons []SeasonInfoV1      `json:"seasons"`
}

// SeasonInfoV1 represents a single season summary for listing.
type SeasonInfoV1 struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	IsActive  bool    `json:"is_active"`
	StartDate string  `json:"start_date"`
	EndDate   *string `json:"end_date"`
}

// -----------------------------------------------------------------------------
// Season Standings Request-Reply Payloads
// -----------------------------------------------------------------------------

// SeasonStandingsRequestPayloadV1 requests standings for a specific season
// via NATS request-reply. Used by the PWA for historical season browsing.
//
// Schema History:
//   - v1.0 (February 2026): Initial version
type SeasonStandingsRequestPayloadV1 struct {
	GuildID  sharedtypes.GuildID `json:"guild_id"`
	SeasonID string              `json:"season_id"`
}

// SeasonStandingsResponsePayloadV1 contains the standings response for request-reply.
// Includes SeasonName for display purposes (unlike the event-driven response).
//
// Schema History:
//   - v1.0 (February 2026): Initial version
type SeasonStandingsResponsePayloadV1 struct {
	GuildID    sharedtypes.GuildID    `json:"guild_id"`
	SeasonID   string                 `json:"season_id"`
	SeasonName string                 `json:"season_name"`
	Standings  []SeasonStandingItemV1 `json:"standings"`
}

// -----------------------------------------------------------------------------
// Failure Payloads
// -----------------------------------------------------------------------------

// AdminFailedPayloadV1 is a generic failure payload for admin operations.
//
// Schema History:
//   - v1.0 (February 2026): Initial version
type AdminFailedPayloadV1 struct {
	GuildID sharedtypes.GuildID `json:"guild_id"`
	Reason  string              `json:"reason"`
}
