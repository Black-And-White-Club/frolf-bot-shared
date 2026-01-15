// Package userevents contains user-related domain events.
//
// This file defines the User Tag Flow - events for tag availability checks
// and tag assignments during user operations.
//
// # Flow Sequences
//
// ## Tag Availability Check Flow
//  1. Request -> TagAvailabilityCheckRequestedV1
//  2. Available -> TagAvailableV1
//  3. OR Unavailable -> TagUnavailableV1
//
// ## Tag Assignment Flow (for user creation)
//  1. Request -> TagAssignmentRequestedV1
//  2. Success -> TagAssignedForUserCreationV1
//  3. OR Failure -> TagAssignmentFailedV1
//
// # Relationship to Leaderboard Module
//
// Tag operations are coordinated with the leaderboard module:
//   - TagAvailabilityCheckRequestedV1 -> triggers leaderboard tag check
//   - TagAssignmentRequestedV1 -> triggers leaderboard tag assignment
//
// # Versioning Strategy
//
// All events include a V1 suffix for future schema evolution.
package userevents

// =============================================================================
// TAG AVAILABILITY CHECK FLOW - Event Constants
// =============================================================================

// NOTE: Tag availability events are defined in events/shared/tag_availability.go.

// =============================================================================
// TAG ASSIGNMENT FLOW - Event Constants
// =============================================================================

// NOTE: Tag assignment events are defined in events/shared/user_tag_assignments.go.

// =============================================================================
// TAG FLOW - Payload Types
// =============================================================================

// -----------------------------------------------------------------------------
// Tag Availability Payloads
// -----------------------------------------------------------------------------

// NOTE: Tag availability payloads are defined in events/shared/tag_availability.go.

// -----------------------------------------------------------------------------
// Tag Assignment Payloads
// -----------------------------------------------------------------------------

// NOTE: Tag assignment payloads are defined in events/shared/user_tag_assignments.go.
