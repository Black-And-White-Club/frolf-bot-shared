package clubevents

import sharedevents "github.com/Black-And-White-Club/frolf-bot-shared/events/shared"

// GetV1Registry returns all modern events for the club functional area
func GetV1Registry() map[string]sharedevents.EventInfo {
	return map[string]sharedevents.EventInfo{
		ClubInfoRequestV1: {
			Payload:     &ClubInfoRequestPayloadV1{},
			Summary:     "Club Info Requested",
			Description: "Request to retrieve club information.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServicePWA, Module: "club"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServiceBackend, Module: "club"}},
		},
		ClubInfoResponseV1: {
			Payload:     &ClubInfoResponsePayloadV1{},
			Summary:     "Club Info Response",
			Description: "Response containing club information.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "club"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServicePWA, Module: "club"}},
		},
		ClubUpdatedV1: {
			Payload:     &ClubUpdatedPayloadV1{},
			Summary:     "Club Updated",
			Description: "Club name or icon was updated.",
			Producer:    sharedevents.Actor{Service: sharedevents.ServiceBackend, Module: "club"},
			Consumers:   []sharedevents.Actor{{Service: sharedevents.ServicePWA, Module: "club"}},
		},
	}
}
