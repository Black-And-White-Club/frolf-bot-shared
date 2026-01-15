// Package scoreevents contains score-related domain events.
//
// This file defines the Score Processing Flow - events for processing round scores
// and triggering leaderboard updates.
//
// # Flow Sequence
//
//  1. Round finalized -> ProcessRoundScoresRequestedV1
//  2. Processing complete -> ProcessRoundScoresSucceededV1
//  3. OR Processing failed -> ProcessRoundScoresFailedV1
//
// # Relationship to Other Modules
//
// Score processing events trigger leaderboard events:
//   - ProcessRoundScoresSucceededV1 -> publishes LeaderboardBatchTagAssignmentRequested
//
// # Versioning Strategy
//
// All events include a V1 suffix for future schema evolution.
package scoreevents

// =============================================================================
// SCORE PROCESSING FLOW - Event Constants
// =============================================================================

// NOTE: All score processing topics are defined in events/shared/score_processing.go.

// =============================================================================
// SCORE PROCESSING FLOW - Payload Types
// =============================================================================

// NOTE: All score processing payloads are defined in events/shared/score_processing.go.
