package guildevents

import sharedevents "github.com/Black-And-White-Club/frolf-bot-shared/events/shared"

// GetV1Registry returns all modern events for the guild functional area
func GetV1Registry() map[string]sharedevents.EventInfo {
	return map[string]sharedevents.EventInfo{
		GuildSetupRequestedV1: {
			Payload:     &GuildConfigCreationRequestedPayloadV1{},
			Summary:     "Guild Setup Requested",
			Description: "Initial guild setup requested.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceDiscord, Module: "guild"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceBackend, Module: "guild"}},
		},
		GuildConfigCreationRequestedV1: {
			Payload:     &GuildConfigCreationRequestedPayloadV1{},
			Summary:     "Guild Config Creation Requested",
			Description: "Request to create guild config.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceDiscord, Module: "guild"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceBackend, Module: "guild"}},
		},
		GuildConfigCreatedV1: {
			Payload:     &GuildConfigCreatedPayloadV1{},
			Summary:     "Guild Config Created",
			Description: "Guild config created.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "guild"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceDiscord, Module: "guild"}, {Service: sharedevents.ServiceBackend, Module: "leaderboard"}}},
		GuildConfigCreationFailedV1: {
			Payload:     &GuildConfigCreationFailedPayloadV1{},
			Summary:     "Guild Config Creation Failed",
			Description: "Guild config creation failed.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "guild"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceDiscord, Module: "guild"}},
		},

		GuildConfigUpdateRequestedV1: {
			Payload:     &GuildConfigUpdateRequestedPayloadV1{},
			Summary:     "Guild Config Update Requested",
			Description: "Request to update guild config.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceDiscord, Module: "guild"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceBackend, Module: "guild"}},
		},
		GuildConfigUpdatedV1: {
			Payload:     &GuildConfigUpdatedPayloadV1{},
			Summary:     "Guild Config Updated",
			Description: "Guild config updated.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "guild"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceDiscord, Module: "guild"}},
		},
		GuildConfigUpdateFailedV1: {
			Payload:     &GuildConfigUpdateFailedPayloadV1{},
			Summary:     "Guild Config Update Failed",
			Description: "Guild config update failed.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "guild"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceDiscord, Module: "guild"}},
		},

		GuildConfigRetrievalRequestedV1: {
			Payload:     &GuildConfigRetrievalRequestedPayloadV1{},
			Summary:     "Guild Config Retrieval Requested",
			Description: "Request to retrieve guild config.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceDiscord, Module: "guild"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceBackend, Module: "guild"}},
		},
		GuildConfigRetrievedV1: {
			Payload:     &GuildConfigRetrievedPayloadV1{},
			Summary:     "Guild Config Retrieved",
			Description: "Guild config retrieved.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "guild"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceDiscord, Module: "guild"}},
		},
		GuildConfigRetrievalFailedV1: {
			Payload:     &GuildConfigRetrievalFailedPayloadV1{},
			Summary:     "Guild Config Retrieval Failed",
			Description: "Guild config retrieval failed.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "guild"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceDiscord, Module: "guild"}},
		},

		GuildConfigDeletionRequestedV1: {
			Payload:     &GuildConfigDeletionRequestedPayloadV1{},
			Summary:     "Guild Config Deletion Requested",
			Description: "Request to delete guild config.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceDiscord, Module: "guild"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceBackend, Module: "guild"}},
		},
		GuildConfigDeletedV1: {
			Payload:     &GuildConfigDeletedPayloadV1{},
			Summary:     "Guild Config Deleted",
			Description: "Guild config deleted.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "guild"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceDiscord, Module: "guild"}},
		},
		GuildConfigDeletionFailedV1: {
			Payload:     &GuildConfigDeletionFailedPayloadV1{},
			Summary:     "Guild Config Deletion Failed",
			Description: "Guild config deletion failed.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "guild"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceDiscord, Module: "guild"}},
		},
		GuildConfigDeletionResultsV1: {Payload: &GuildConfigDeletionResultsPayloadV1{}, Summary: "Guild Config Deletion Results", Description: "Per-resource deletion results.", Producer: sharedevents.Actor{Service: sharedevents.ServiceDiscord, Module: "guild"}},
	}
}
