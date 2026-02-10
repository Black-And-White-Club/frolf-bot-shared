package sharedevents

import (
	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
)

// Actor represents a specific functional unit in your system
type Actor struct {
	Service string // e.g., "frolf-bot"
	Module  string // e.g., "round", "user", "leaderboard"
}

// EventInfo is the metadata container used by the AsyncAPI generator
type EventInfo struct {
	Payload     interface{}
	Summary     string
	Description string
	Producer    Actor
	Consumers   []Actor
}

// Service constants to avoid magic strings
const (
	ServiceBackend = "frolf-bot-backend"
	ServiceDiscord = "discord"
	ServicePWA     = "pwa"
)

// GetV1Registry returns all modern events for the shared functional area
func GetV1Registry() map[string]EventInfo {
	return map[string]EventInfo{
		// Round tag lookup flow
		RoundTagLookupRequestedV1: {
			Payload:     &RoundTagLookupRequestedPayloadV1{},
			Summary:     "Round Tag Lookup Requested",
			Description: "Request (from round) to lookup a user's tag in leaderboard.",
			Producer:    Actor{Service: ServiceBackend, Module: "round"},
			Consumers:   []Actor{{Service: ServiceBackend, Module: "leaderboard"}},
		},
		RoundTagLookupFoundV1: {
			Payload:     &RoundTagLookupResultPayloadV1{},
			Summary:     "Round Tag Lookup Found",
			Description: "Leaderboard returned a tag lookup result for a round request.",
			Producer:    Actor{Service: ServiceBackend, Module: "leaderboard"},
			Consumers:   []Actor{{Service: ServiceBackend, Module: "round"}},
		},
		RoundTagLookupNotFoundV1: {
			Payload:     &RoundTagLookupResultPayloadV1{},
			Summary:     "Round Tag Lookup Not Found",
			Description: "Leaderboard did not find a tag for the requested player.",
			Producer:    Actor{Service: ServiceBackend, Module: "leaderboard"},
			Consumers:   []Actor{{Service: ServiceBackend, Module: "round"}},
		},
		RoundTagLookupFailedV1: {
			Payload:     &RoundTagLookupFailedPayloadV1{},
			Summary:     "Round Tag Lookup Failed",
			Description: "An error occurred while looking up a tag for a round request.",
			Producer:    Actor{Service: ServiceBackend, Module: "leaderboard"},
			Consumers:   []Actor{{Service: ServiceBackend, Module: "round"}},
		},

		// SyncRoundsTagRequestV1 is the cross-module trigger from Leaderboard to Round
		SyncRoundsTagRequestV1: {
			Payload:     &SyncRoundsTagRequestPayloadV1{},
			Summary:     "Synchronize Scheduled Round Tags",
			Description: "Triggered by the leaderboard after tag changes to ensure all upcoming rounds reflect the new source-of-truth tags.",
			Producer:    Actor{Service: ServiceBackend, Module: "leaderboard"},
			Consumers:   []Actor{{Service: ServiceBackend, Module: "round"}},
		},
		// Discord / Leaderboard lookup flow
		DiscordTagLookupRequestedV1: {
			Payload:     &DiscordTagLookupRequestedPayloadV1{},
			Summary:     "Discord Tag Lookup Requested",
			Description: "Discord-initiated request to lookup a user's tag in leaderboard.",
			Producer:    Actor{Service: ServiceDiscord, Module: "leaderboard"},
			Consumers:   []Actor{{Service: ServiceBackend, Module: "leaderboard"}},
		},

		// Canonical leaderboard-owned response topics
		LeaderboardTagLookupSucceededV1: {
			Payload:     &DiscordTagLookupResultPayloadV1{},
			Summary:     "Leaderboard Tag Lookup Succeeded",
			Description: "Leaderboard returned a successful tag lookup result.",
			Producer:    Actor{Service: ServiceBackend, Module: "leaderboard"},
			Consumers:   []Actor{{Service: ServiceBackend, Module: "round"}, {Service: ServiceDiscord, Module: "leaderboard"}}},
		LeaderboardTagLookupNotFoundV1: {
			Payload:     &DiscordTagLookupResultPayloadV1{},
			Summary:     "Leaderboard Tag Lookup Not Found",
			Description: "Leaderboard returned not-found for the lookup.",
			Producer:    Actor{Service: ServiceBackend, Module: "leaderboard"},
			Consumers:   []Actor{{Service: ServiceBackend, Module: "round"}, {Service: ServiceDiscord, Module: "leaderboard"}}},
		LeaderboardTagLookupFailedV1: {
			Payload:     &DiscordTagLookupFailedPayloadV1{},
			Summary:     "Leaderboard Tag Lookup Failed",
			Description: "Leaderboard failed to perform the tag lookup.",
			Producer:    Actor{Service: ServiceBackend, Module: "leaderboard"},
			Consumers:   []Actor{{Service: ServiceBackend, Module: "round"}, {Service: ServiceDiscord, Module: "leaderboard"}}},

		// Batch tag assignment
		LeaderboardBatchTagAssignmentRequestedV1: {
			Payload:     &BatchTagAssignmentRequestedPayloadV1{},
			Summary:     "Leaderboard Batch Tag Assignment Requested",
			Description: "Batch request to assign tags to multiple users.",
			Producer:    Actor{Service: ServiceBackend, Module: "score"},
			Consumers:   []Actor{{Service: ServiceBackend, Module: "leaderboard"}},
		},

		// Tag availability flow
		TagAvailabilityCheckRequestedV1: {
			Payload:     &TagAvailabilityCheckRequestedPayloadV1{},
			Summary:     "Tag Availability Check Requested",
			Description: "Request to check if a tag is available.",
			Producer:    Actor{Service: ServiceBackend, Module: "user"},
			Consumers:   []Actor{{Service: ServiceBackend, Module: "leaderboard"}},
		},
		TagAvailableV1: {
			Payload:     &TagAvailablePayloadV1{},
			Summary:     "Tag Available",
			Description: "Tag is available.",
			Producer:    Actor{Service: ServiceBackend, Module: "leaderboard"},
			Consumers:   []Actor{{Service: ServiceBackend, Module: "user"}},
		},
		TagUnavailableV1: {
			Payload:     &TagUnavailablePayloadV1{},
			Summary:     "Tag Unavailable",
			Description: "Tag is not available.",
			Producer:    Actor{Service: ServiceBackend, Module: "leaderboard"},
			Consumers:   []Actor{{Service: ServiceBackend, Module: "user"}},
		},
		TagAvailabilityCheckFailedV1: {
			Payload:     &TagAvailabilityCheckFailedPayloadV1{},
			Summary:     "Tag Availability Check Failed",
			Description: "Tag availability check failed.",
			Producer:    Actor{Service: ServiceBackend, Module: "leaderboard"},
			Consumers:   []Actor{{Service: ServiceBackend, Module: "user"}},
		},

		// Tag assignment flow (user creation)
		TagAssignmentRequestedV1: {
			Payload:     &TagAssignmentRequestedPayloadV1{},
			Summary:     "Tag Assignment Requested",
			Description: "Request to assign a tag to a user.",
			Producer:    Actor{Service: ServiceBackend, Module: "user"},
			Consumers:   []Actor{{Service: ServiceBackend, Module: "leaderboard"}},
		},
		TagAssignedForUserCreationV1: {
			Payload:     &TagAssignedForUserCreationPayloadV1{},
			Summary:     "Tag Assigned For User Creation",
			Description: "Tag assigned during user creation.",
			Producer:    Actor{Service: ServiceBackend, Module: "leaderboard"},
			Consumers:   []Actor{{Service: ServiceBackend, Module: "user"}},
		},
		TagAssignmentFailedV1: {
			Payload:     &TagAssignmentFailedPayloadV1{},
			Summary:     "Tag Assignment Failed",
			Description: "Tag assignment failed.",
			Producer:    Actor{Service: ServiceBackend, Module: "leaderboard"},
			Consumers:   []Actor{{Service: ServiceBackend, Module: "user"}},
		},

		// Tag number lookup flow
		GetTagByUserIDRequestedV1: {
			Payload:     &TagNumberRequestPayloadV1{},
			Summary:     "Get Tag By User ID Requested",
			Description: "Request to lookup a tag number by user ID.",
			Producer:    Actor{Service: ServiceBackend, Module: "round"},
			Consumers:   []Actor{{Service: ServiceBackend, Module: "leaderboard"}},
		},
		RoundGetTagByUserIDRequestedV1: {
			Payload:     &TagNumberRequestPayloadV1{},
			Summary:     "Round Get Tag By User ID Requested",
			Description: "Round-specific request to lookup a tag number by user ID.",
			Producer:    Actor{Service: ServiceBackend, Module: "round"},
			Consumers:   []Actor{{Service: ServiceBackend, Module: "leaderboard"}},
		},
		GetTagNumberResponseV1: {
			Payload:     &GetTagNumberResponsePayloadV1{},
			Summary:     "Get Tag Number Response",
			Description: "Response containing a tag number lookup result.",
			Producer:    Actor{Service: ServiceBackend, Module: "leaderboard"},
			Consumers:   []Actor{{Service: ServiceBackend, Module: "round"}, {Service: ServiceDiscord, Module: "leaderboard"}}},
		GetTagNumberFailedV1: {
			Payload:     &GetTagNumberFailedPayloadV1{},
			Summary:     "Get Tag Number Failed",
			Description: "Tag number lookup failed.",
			Producer:    Actor{Service: ServiceBackend, Module: "leaderboard"},
			Consumers:   []Actor{{Service: ServiceBackend, Module: "round"}, {Service: ServiceDiscord, Module: "leaderboard"}}},
		RoundTagNumberFoundV1: {
			Payload:     &RoundTagNumberFoundPayloadV1{},
			Summary:     "Round Tag Number Found",
			Description: "Legacy round tag number found.",
			Producer:    Actor{Service: ServiceBackend, Module: "leaderboard"},
			Consumers:   []Actor{{Service: ServiceBackend, Module: "round"}},
		},
		RoundTagNumberNotFoundV1: {
			Payload:     &RoundTagNumberNotFoundPayloadV1{},
			Summary:     "Round Tag Number Not Found",
			Description: "Legacy round tag number not found.",
			Producer:    Actor{Service: ServiceBackend, Module: "leaderboard"},
			Consumers:   []Actor{{Service: ServiceBackend, Module: "round"}},
		},

		// Score update flow (shared across modules)
		ScoreUpdateRequestedV1: {
			Payload:     &ScoreUpdateRequestedPayloadV1{},
			Summary:     "Score Update Requested",
			Description: "Request to update an individual score.",
			Producer:    Actor{Service: ServiceDiscord, Module: "score"},
			Consumers:   []Actor{{Service: ServiceBackend, Module: "score"}},
		},
		ScoreUpdatedV1: {
			Payload:     &ScoreUpdatedPayloadV1{},
			Summary:     "Score Updated",
			Description: "Individual score successfully updated.",
			Producer:    Actor{Service: ServiceBackend, Module: "score"},
			Consumers:   []Actor{{Service: ServiceDiscord, Module: "score"}, {Service: ServiceDiscord, Module: "round"}}},
		ScoreUpdateFailedV1: {
			Payload:     &ScoreUpdateFailedPayloadV1{},
			Summary:     "Score Update Failed",
			Description: "Individual score update failed.",
			Producer:    Actor{Service: ServiceBackend, Module: "score"},
			Consumers:   []Actor{{Service: ServiceDiscord, Module: "score"}},
		},
		ScoreBulkUpdateRequestedV1: {
			Payload:     &ScoreBulkUpdateRequestedPayloadV1{},
			Summary:     "Score Bulk Update Requested",
			Description: "Bulk score update requested.",
			Producer:    Actor{Service: ServiceBackend, Module: "score"},
			Consumers:   []Actor{{Service: ServiceBackend, Module: "score"}},
		},
		ScoreBulkUpdatedV1: {
			Payload:     &ScoreBulkUpdatedPayloadV1{},
			Summary:     "Score Bulk Updated",
			Description: "Bulk score update completed.",
			Producer:    Actor{Service: ServiceBackend, Module: "score"},
			Consumers:   []Actor{},
		},

		// Score processing flow
		ProcessRoundScoresRequestedV1: {
			Payload:     &ProcessRoundScoresRequestedPayloadV1{},
			Summary:     "Process Round Scores Requested",
			Description: "Request to process a round's scores.",
			Producer:    Actor{Service: ServiceBackend, Module: "round"},
			Consumers:   []Actor{{Service: ServiceBackend, Module: "score"}},
		},
		ProcessRoundScoresSucceededV1: {
			Payload:     &ProcessRoundScoresSucceededPayloadV1{},
			Summary:     "Process Round Scores Succeeded",
			Description: "Round score processing completed successfully.",
			Producer:    Actor{Service: ServiceBackend, Module: "score"},
			Consumers:   []Actor{{Service: ServiceBackend, Module: "leaderboard"}},
		},
		ProcessRoundScoresFailedV1: {
			Payload:     &ProcessRoundScoresFailedPayloadV1{},
			Summary:     "Process Round Scores Failed",
			Description: "Round score processing failed.",
			Producer:    Actor{Service: ServiceBackend, Module: "score"},
			Consumers:   []Actor{{Service: ServiceDiscord, Module: "score"}},
		},
		ScoreModuleNotificationErrorV1: {
			Payload:     &ScoreModuleNotificationErrorPayloadV1{},
			Summary:     "Score Module Notification Error",
			Description: "Failed to notify score module.",
			Producer:    Actor{Service: ServiceBackend, Module: "round"},
			Consumers:   []Actor{},
		},

		// Club sync flow (user → club cross-module)
		ClubSyncFromDiscordRequestedV1: {
			Payload:     &ClubSyncFromDiscordRequestedPayloadV1{},
			Summary:     "Club Sync From Discord Requested",
			Description: "User signup included guild metadata; sync club record.",
			Producer:    Actor{Service: ServiceBackend, Module: "user"},
			Consumers:   []Actor{{Service: ServiceBackend, Module: "club"}},
		},
		PointsAwardedV1: {
			Payload:     &sharedtypes.PointsAwardedPayloadV1{},
			Summary:     "Points Awarded",
			Description: "Points awarded after round processing.",
			Producer:    Actor{Service: ServiceBackend, Module: "leaderboard"},
			Consumers:   []Actor{{Service: ServiceDiscord, Module: "round"}},
		},
	}
}
