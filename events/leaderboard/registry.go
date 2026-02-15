package leaderboardevents

import sharedevents "github.com/Black-And-White-Club/frolf-bot-shared/events/shared"

// GetV1Registry returns all modern events for the leaderboard functional area
func GetV1Registry() map[string]sharedevents.EventInfo {
	return map[string]sharedevents.EventInfo{
		LeaderboardRoundFinalizedV1: {
			Payload:     &LeaderboardRoundFinalizedPayloadV1{},
			Summary:     "Leaderboard Round Finalized",
			Description: "Round finalized and ready for leaderboard update.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "round"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceBackend, Module: "leaderboard"}},
		},
		LeaderboardUpdateRequestedV1: {
			Payload:     &LeaderboardUpdateRequestedPayloadV1{},
			Summary:     "Leaderboard Update Requested",
			Description: "Request to update the leaderboard.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "leaderboard"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceBackend, Module: "leaderboard"}},
		},
		LeaderboardUpdatedV1: {
			Payload:     &LeaderboardUpdatedPayloadV1{},
			Summary:     "Leaderboard Updated",
			Description: "Leaderboard successfully updated.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "leaderboard"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceDiscord, Module: "leaderboard"}, {Service: sharedevents.ServiceBackend, Module: "round"}}},
		LeaderboardUpdateFailedV1: {
			Payload:     &LeaderboardUpdateFailedPayloadV1{},
			Summary:     "Leaderboard Update Failed",
			Description: "Leaderboard update failed.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "leaderboard"},
			Consumers:   []sharedevents.Actor{},
		},
		DeactivateOldLeaderboardV1: {
			Payload:     &DeactivateOldLeaderboardPayloadV1{},
			Summary:     "Deactivate Old Leaderboard",
			Description: "Deactivate a previous leaderboard after update.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "leaderboard"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceBackend, Module: "leaderboard"}},
		},
		sharedevents.PointsAwardedV1: {
			Payload:     &sharedevents.PointsAwardedPayloadV1{},
			Summary:     "Leaderboard Points Awarded",
			Description: "Points awarded for a round.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "leaderboard"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceDiscord, Module: "round"}},
		},

		GetLeaderboardRequestedV1: {
			Payload:     &GetLeaderboardRequestedPayloadV1{},
			Summary:     "Get Leaderboard Requested",
			Description: "Request to retrieve leaderboard.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceDiscord, Module: "leaderboard"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceBackend, Module: "leaderboard"}},
		},
		GetLeaderboardResponseV1: {
			Payload:     &GetLeaderboardResponsePayloadV1{},
			Summary:     "Get Leaderboard Response",
			Description: "Leaderboard retrieval response.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "leaderboard"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceDiscord, Module: "leaderboard"}},
		},
		GetLeaderboardFailedV1: {
			Payload:     &GetLeaderboardFailedPayloadV1{},
			Summary:     "Get Leaderboard Failed",
			Description: "Leaderboard retrieval failed.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "leaderboard"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceDiscord, Module: "leaderboard"}},
		},

		LeaderboardTagAssignmentRequestedV1: {
			Payload:     &LeaderboardTagAssignmentRequestedPayloadV1{},
			Summary:     "Leaderboard Tag Assignment Requested",
			Description: "Request to assign a tag to a user.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "user"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceBackend, Module: "leaderboard"}},
		},
		LeaderboardTagAssignedV1: {
			Payload:     &LeaderboardTagAssignedPayloadV1{},
			Summary:     "Leaderboard Tag Assigned",
			Description: "Tag successfully assigned to a user.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "leaderboard"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceBackend, Module: "user"}, {Service: sharedevents.ServiceDiscord, Module: "leaderboard"}}},
		LeaderboardTagAssignmentFailedV1: {
			Payload:     &LeaderboardTagAssignmentFailedPayloadV1{},
			Summary:     "Leaderboard Tag Assignment Failed",
			Description: "Tag assignment failed.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "leaderboard"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceBackend, Module: "user"}, {Service: sharedevents.ServiceDiscord, Module: "leaderboard"}}},
		LeaderboardBatchTagAssignmentRequestedV1: {
			Payload:     &sharedevents.BatchTagAssignmentRequestedPayloadV1{},
			Summary:     "Leaderboard Batch Tag Assignment Requested",
			Description: "Batch tag assignment requested.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "score"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceBackend, Module: "leaderboard"}},
		},
		LeaderboardBatchTagAssignedV1: {
			Payload:     &LeaderboardBatchTagAssignedPayloadV1{},
			Summary:     "Leaderboard Batch Tag Assigned",
			Description: "Batch tag assignment completed.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "leaderboard"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceDiscord, Module: "leaderboard"}},
		},
		LeaderboardBatchTagAssignmentFailedV1: {
			Payload:     &LeaderboardBatchTagAssignmentFailedPayloadV1{},
			Summary:     "Leaderboard Batch Tag Assignment Failed",
			Description: "Batch tag assignment failed.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "leaderboard"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceDiscord, Module: "leaderboard"}},
		},

		LeaderboardTagUpdatedV1: {
			Payload:     &LeaderboardTagUpdatedPayloadV1{},
			Summary:     "Leaderboard Tag Updated",
			Description: "User tag updated in the leaderboard.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "leaderboard"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceBackend, Module: "round"}, {Service: sharedevents.ServiceDiscord, Module: "leaderboard"}}},

		TagUpdateForScheduledRoundsV1: {
			Payload:     &TagUpdateForScheduledRoundsPayloadV1{},
			Summary:     "Tag Update For Scheduled Rounds",
			Description: "Update tags for scheduled rounds after leaderboard changes.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "leaderboard"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceBackend, Module: "round"}},
		},

		TagSwapRequestedV1: {
			Payload:     &TagSwapRequestedPayloadV1{},
			Summary:     "Tag Swap Requested",
			Description: "Request to swap tags.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceDiscord, Module: "leaderboard"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceBackend, Module: "leaderboard"}},
		},
		TagSwapInitiatedV1: {
			Payload:     &TagSwapInitiatedPayloadV1{},
			Summary:     "Tag Swap Initiated",
			Description: "Tag swap initiated.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "leaderboard"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceBackend, Module: "leaderboard"}},
		},
		TagSwapProcessedV1: {
			Payload:     &TagSwapProcessedPayloadV1{},
			Summary:     "Tag Swap Processed",
			Description: "Tag swap processed.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "leaderboard"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceDiscord, Module: "leaderboard"}},
		},
		TagSwapFailedV1: {
			Payload:     &TagSwapFailedPayloadV1{},
			Summary:     "Tag Swap Failed",
			Description: "Tag swap failed.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "leaderboard"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceDiscord, Module: "leaderboard"}},
		},

		LeaderboardTraceEventV1: {Payload: &sharedevents.TracePayloadV1{}, Summary: "Leaderboard Trace Event", Description: "Tracing/observability event for leaderboard.", Producer: sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "leaderboard"}},

		// Admin Events
		LeaderboardPointHistoryRequestedV1: {
			Payload:     &PointHistoryRequestedPayloadV1{},
			Summary:     "Point History Requested",
			Description: "Request point history for a member.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceDiscord, Module: "admin"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceBackend, Module: "leaderboard"}},
		},
		LeaderboardPointHistoryResponseV1: {
			Payload:     &PointHistoryResponsePayloadV1{},
			Summary:     "Point History Response",
			Description: "Point history response.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "leaderboard"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceDiscord, Module: "admin"}},
		},
		LeaderboardPointHistoryFailedV1: {
			Payload:     &AdminFailedPayloadV1{},
			Summary:     "Point History Failed",
			Description: "Point history request failed.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "leaderboard"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceDiscord, Module: "admin"}},
		},

		LeaderboardManualPointAdjustmentV1: {
			Payload:     &ManualPointAdjustmentPayloadV1{},
			Summary:     "Manual Point Adjustment",
			Description: "Request manual point adjustment.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceDiscord, Module: "admin"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceBackend, Module: "leaderboard"}},
		},
		LeaderboardManualPointAdjustmentSuccessV1: {
			Payload:     &ManualPointAdjustmentSuccessPayloadV1{},
			Summary:     "Manual Point Adjustment Success",
			Description: "Manual point adjustment succeeded.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "leaderboard"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceDiscord, Module: "admin"}},
		},
		LeaderboardManualPointAdjustmentFailedV1: {
			Payload:     &AdminFailedPayloadV1{},
			Summary:     "Manual Point Adjustment Failed",
			Description: "Manual point adjustment failed.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "leaderboard"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceDiscord, Module: "admin"}},
		},

		LeaderboardRecalculateRoundV1: {
			Payload:     &RecalculateRoundPayloadV1{},
			Summary:     "Recalculate Round",
			Description: "Request round recalculation.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceDiscord, Module: "admin"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceBackend, Module: "leaderboard"}},
		},
		LeaderboardRecalculateRoundSuccessV1: {
			Payload:     &RecalculateRoundSuccessPayloadV1{},
			Summary:     "Recalculate Round Success",
			Description: "Round recalculation succeeded.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "leaderboard"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceDiscord, Module: "admin"}},
		},
		LeaderboardRecalculateRoundFailedV1: {
			Payload:     &AdminFailedPayloadV1{},
			Summary:     "Recalculate Round Failed",
			Description: "Round recalculation failed.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "leaderboard"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceDiscord, Module: "admin"}},
		},

		LeaderboardStartNewSeasonV1: {
			Payload:     &StartNewSeasonPayloadV1{},
			Summary:     "Start New Season",
			Description: "Request to start a new season.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceDiscord, Module: "admin"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceBackend, Module: "leaderboard"}},
		},
		LeaderboardStartNewSeasonSuccessV1: {
			Payload:     &StartNewSeasonSuccessPayloadV1{},
			Summary:     "Start New Season Success",
			Description: "New season started successfully.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "leaderboard"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceDiscord, Module: "admin"}},
		},
		LeaderboardStartNewSeasonFailedV1: {
			Payload:     &AdminFailedPayloadV1{},
			Summary:     "Start New Season Failed",
			Description: "Failed to start new season.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "leaderboard"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceDiscord, Module: "admin"}},
		},

		LeaderboardGetSeasonStandingsV1: {
			Payload:     &GetSeasonStandingsPayloadV1{},
			Summary:     "Get Season Standings",
			Description: "Request season standings.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceDiscord, Module: "admin"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceBackend, Module: "leaderboard"}},
		},
		LeaderboardGetSeasonStandingsResponseV1: {
			Payload:     &GetSeasonStandingsResponsePayloadV1{},
			Summary:     "Get Season Standings Response",
			Description: "Season standings response.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "leaderboard"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceDiscord, Module: "admin"}},
		},
		LeaderboardGetSeasonStandingsFailedV1: {
			Payload:     &AdminFailedPayloadV1{},
			Summary:     "Get Season Standings Failed",
			Description: "Failed to get season standings.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "leaderboard"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceDiscord, Module: "admin"}},
		},

		// List Seasons Request-Reply Events
		LeaderboardListSeasonsRequestedV1: {
			Payload:     &ListSeasonsRequestPayloadV1{},
			Summary:     "List Seasons Requested",
			Description: "Request to list all seasons for a guild (request-reply).",
			Producer:    sharedevents.Actor{Service: sharedevents.ServicePWA, Module: "leaderboard"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceBackend, Module: "leaderboard"}},
		},
		LeaderboardListSeasonsResponseV1: {
			Payload:     &ListSeasonsResponsePayloadV1{},
			Summary:     "List Seasons Response",
			Description: "Seasons list response (request-reply).",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "leaderboard"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServicePWA, Module: "leaderboard"}},
		},
		LeaderboardListSeasonsFailedV1: {
			Payload:     &AdminFailedPayloadV1{},
			Summary:     "List Seasons Failed",
			Description: "Failed to list seasons (request-reply).",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "leaderboard"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServicePWA, Module: "leaderboard"}},
		},

		// Season Standings Request-Reply Events
		LeaderboardSeasonStandingsRequestV1: {
			Payload:     &SeasonStandingsRequestPayloadV1{},
			Summary:     "Season Standings Request",
			Description: "Request season standings (request-reply).",
			Producer:    sharedevents.Actor{Service: sharedevents.ServicePWA, Module: "leaderboard"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceBackend, Module: "leaderboard"}},
		},
		LeaderboardSeasonStandingsResponseV1: {
			Payload:     &SeasonStandingsResponsePayloadV1{},
			Summary:     "Season Standings Response",
			Description: "Season standings response (request-reply).",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "leaderboard"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServicePWA, Module: "leaderboard"}},
		},
		LeaderboardSeasonStandingsFailedV1: {
			Payload:     &AdminFailedPayloadV1{},
			Summary:     "Season Standings Failed",
			Description: "Failed to get season standings (request-reply).",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "leaderboard"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServicePWA, Module: "leaderboard"}},
		},

		// Tag History Request-Reply Events
		LeaderboardTagHistoryRequestedV1: {
			Payload:     &TagHistoryRequestedPayloadV1{},
			Summary:     "Tag History Requested",
			Description: "Request tag history for a member or tag (request-reply).",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceDiscord, Module: "leaderboard"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceBackend, Module: "leaderboard"}},
		},
		LeaderboardTagHistoryResponseV1: {
			Payload:     &TagHistoryResponsePayloadV1{},
			Summary:     "Tag History Response",
			Description: "Tag history response (request-reply).",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "leaderboard"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceDiscord, Module: "leaderboard"}, {Service: sharedevents.ServicePWA, Module: "leaderboard"}},
		},
		LeaderboardTagHistoryFailedV1: {
			Payload:     &TagHistoryFailedPayloadV1{},
			Summary:     "Tag History Failed",
			Description: "Tag history request failed (request-reply).",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "leaderboard"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceDiscord, Module: "leaderboard"}, {Service: sharedevents.ServicePWA, Module: "leaderboard"}},
		},

		// Tag Graph Request-Reply Events
		LeaderboardTagGraphRequestedV1: {
			Payload:     &TagGraphRequestedPayloadV1{},
			Summary:     "Tag Graph Requested",
			Description: "Request PNG tag history chart (request-reply).",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceDiscord, Module: "leaderboard"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceBackend, Module: "leaderboard"}},
		},
		LeaderboardTagGraphResponseV1: {
			Payload:     &TagGraphResponsePayloadV1{},
			Summary:     "Tag Graph Response",
			Description: "PNG tag history chart response (request-reply).",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "leaderboard"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceDiscord, Module: "leaderboard"}},
		},
		LeaderboardTagGraphFailedV1: {
			Payload:     &TagGraphFailedPayloadV1{},
			Summary:     "Tag Graph Failed",
			Description: "Tag graph generation failed (request-reply).",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "leaderboard"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceDiscord, Module: "leaderboard"}},
		},

		// Tag List Request-Reply Events
		LeaderboardTagListRequestedV1: {
			Payload:     &TagListRequestedPayloadV1{},
			Summary:     "Tag List Requested",
			Description: "Request master tag list (request-reply).",
			Producer:    sharedevents.Actor{Service: sharedevents.ServicePWA, Module: "leaderboard"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceBackend, Module: "leaderboard"}},
		},
		LeaderboardTagListResponseV1: {
			Payload:     &TagListResponsePayloadV1{},
			Summary:     "Tag List Response",
			Description: "Master tag list response (request-reply).",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "leaderboard"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServicePWA, Module: "leaderboard"}},
		},
		LeaderboardTagListFailedV1: {
			Payload:     &TagListFailedPayloadV1{},
			Summary:     "Tag List Failed",
			Description: "Tag list request failed (request-reply).",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "leaderboard"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServicePWA, Module: "leaderboard"}},
		},
	}
}
