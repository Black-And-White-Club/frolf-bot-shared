package discordleaderboardevents

import sharedevents "github.com/Black-And-White-Club/frolf-bot-shared/events/shared"

// GetV1Registry returns all modern events for the discord leaderboard functional area
func GetV1Registry() map[string]sharedevents.EventInfo {

	return map[string]sharedevents.EventInfo{
		// Retrieval
		LeaderboardRetrieveRequestV1: {Payload: &LeaderboardRetrieveRequestPayloadV1{}, Summary: "Leaderboard Retrieve Request", Description: "Discord request to retrieve the leaderboard.", Producer: sharedevents.Actor{Service: sharedevents.ServiceDiscord, Module: "leaderboard"}, Consumers: []sharedevents.Actor{{Service: sharedevents.ServiceDiscord, Module: "leaderboard"}}},
		LeaderboardRetrievedV1:       {Payload: &LeaderboardRetrievedPayloadV1{}, Summary: "Leaderboard Retrieved", Description: "Leaderboard data ready for Discord display.", Producer: sharedevents.Actor{Service: sharedevents.ServiceDiscord, Module: "leaderboard"}, Consumers: []sharedevents.Actor{{Service: sharedevents.ServiceDiscord, Module: "leaderboard"}}},

		// Tag assignment
		LeaderboardTagAssignRequestV1: {Payload: &LeaderboardTagAssignRequestPayloadV1{}, Summary: "Leaderboard Tag Assign Request", Description: "Admin requested a tag assignment via Discord.", Producer: sharedevents.Actor{Service: sharedevents.ServiceDiscord, Module: "leaderboard"}, Consumers: []sharedevents.Actor{{Service: sharedevents.ServiceDiscord, Module: "leaderboard"}}},
		LeaderboardTagAssignedV1:      {Payload: &LeaderboardTagAssignedPayloadV1{}, Summary: "Leaderboard Tag Assigned", Description: "Tag assignment succeeded (Discord).", Producer: sharedevents.Actor{Service: sharedevents.ServiceDiscord, Module: "leaderboard"}, Consumers: []sharedevents.Actor{{Service: sharedevents.ServiceDiscord, Module: "leaderboard"}}},
		LeaderboardTagAssignFailedV1:  {Payload: &LeaderboardTagAssignFailedPayloadV1{}, Summary: "Leaderboard Tag Assign Failed", Description: "Tag assignment failed (Discord).", Producer: sharedevents.Actor{Service: sharedevents.ServiceDiscord, Module: "leaderboard"}, Consumers: []sharedevents.Actor{{Service: sharedevents.ServiceDiscord, Module: "leaderboard"}}},
		LeaderboardBatchTagAssignedV1: {Payload: &BatchTagAssignedPayloadV1{}, Summary: "Leaderboard Batch Tag Assigned", Description: "Batch tag assignment result for Discord.", Producer: sharedevents.Actor{Service: sharedevents.ServiceDiscord, Module: "leaderboard"}, Consumers: []sharedevents.Actor{{Service: sharedevents.ServiceDiscord, Module: "leaderboard"}}},

		// Availability
		LeaderboardTagAvailabilityRequestV1:  {Payload: &LeaderboardTagAvailabilityRequestPayloadV1{}, Summary: "Leaderboard Tag Availability Request", Description: "Discord request to check tag availability.", Producer: sharedevents.Actor{Service: sharedevents.ServiceDiscord, Module: "leaderboard"}, Consumers: []sharedevents.Actor{{Service: sharedevents.ServiceDiscord, Module: "leaderboard"}}},
		LeaderboardTagAvailabilityResponseV1: {Payload: &LeaderboardTagAvailabilityResponsePayloadV1{}, Summary: "Leaderboard Tag Availability Response", Description: "Discord response with tag availability.", Producer: sharedevents.Actor{Service: sharedevents.ServiceDiscord, Module: "leaderboard"}, Consumers: []sharedevents.Actor{{Service: sharedevents.ServiceDiscord, Module: "leaderboard"}}},

		// Tag swap
		LeaderboardTagSwapRequestV1: {Payload: &LeaderboardTagSwapRequestPayloadV1{}, Summary: "Leaderboard Tag Swap Request", Description: "Discord admin requested a tag swap.", Producer: sharedevents.Actor{Service: sharedevents.ServiceDiscord, Module: "leaderboard"}, Consumers: []sharedevents.Actor{{Service: sharedevents.ServiceDiscord, Module: "leaderboard"}}},
		LeaderboardTagSwappedV1:     {Payload: &LeaderboardTagSwappedPayloadV1{}, Summary: "Leaderboard Tag Swapped", Description: "Tag swap succeeded (Discord).", Producer: sharedevents.Actor{Service: sharedevents.ServiceDiscord, Module: "leaderboard"}, Consumers: []sharedevents.Actor{{Service: sharedevents.ServiceDiscord, Module: "leaderboard"}}},
		LeaderboardTagSwapFailedV1:  {Payload: &LeaderboardTagSwapFailedPayloadV1{}, Summary: "Leaderboard Tag Swap Failed", Description: "Tag swap failed (Discord).", Producer: sharedevents.Actor{Service: sharedevents.ServiceDiscord, Module: "leaderboard"}, Consumers: []sharedevents.Actor{{Service: sharedevents.ServiceDiscord, Module: "leaderboard"}}},
	}
}
