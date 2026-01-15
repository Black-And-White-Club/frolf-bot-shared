package discordevents

import sharedevents "github.com/Black-And-White-Club/frolf-bot-shared/events/shared"

// GetV1Registry returns all modern events for the discord functional area
func GetV1Registry() map[string]sharedevents.EventInfo {
	return map[string]sharedevents.EventInfo{
		SendDMV1:               {Payload: &SendDMPayloadV1{}, Summary: "Send DM", Description: "Request to send a direct message to a user via Discord.", Producer: sharedevents.Actor{Service: sharedevents.ServiceDiscord, Module: "app"}},
		DMSentV1:               {Payload: &DMSentPayloadV1{}, Summary: "DM Sent", Description: "Confirmation that a DM was sent to a user.", Producer: sharedevents.Actor{Service: sharedevents.ServiceDiscord, Module: "app"}},
		DMErrorV1:              {Payload: &DMErrorPayloadV1{}, Summary: "DM Error", Description: "Details about a failure to send a DM.", Producer: sharedevents.Actor{Service: sharedevents.ServiceDiscord, Module: "app"}},
		InteractionRespondedV1: {Payload: &InteractionRespondedPayloadV1{}, Summary: "Interaction Responded", Description: "Record that a Discord interaction was responded to.", Producer: sharedevents.Actor{Service: sharedevents.ServiceDiscord, Module: "app"}},
		DiscordEventTraceV1:    {Payload: &TracePayloadV1{}, Summary: "Discord Event Trace", Description: "Observability/trace events emitted by the discord service.", Producer: sharedevents.Actor{Service: sharedevents.ServiceDiscord, Module: "app"}},
	}
}
