// Package scoreevents contains score-related domain events.
//
// This file defines the Score Update Flow - events for updating individual
// and bulk scores.
//
// # Flow Sequences
//
// ## Individual Score Update
//  1. Score update requested -> ScoreUpdateRequestedV1
//  2. Success -> ScoreUpdatedV1
//  3. OR Failure -> ScoreUpdateFailedV1
//
// ## Bulk Score Update
//  1. Bulk update requested -> ScoreBulkUpdateRequestedV1
//  2. Completion -> ScoreBulkUpdatedV1
//
// # Versioning Strategy
//
// All events include a V1 suffix for future schema evolution.
package scoreevents

// =============================================================================
// SCORE UPDATE FLOW - Event Constants
// =============================================================================

// -----------------------------------------------------------------------------
// Individual Score Update Events
// -----------------------------------------------------------------------------

// NOTE: All score update topics are defined in events/shared/score_updates.go.

// =============================================================================
// SCORE UPDATE FLOW - Payload Types
// =============================================================================

// -----------------------------------------------------------------------------
// Individual Score Update Payloads
// -----------------------------------------------------------------------------

// NOTE: All score update payloads are defined in events/shared/score_updates.go.
