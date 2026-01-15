package discordscoreevents

import sharedevents "github.com/Black-And-White-Club/frolf-bot-shared/events/shared"

// GetV1Registry returns all modern events for the discord score functional area
// Exposes only Discord-prefixed subjects and uses shared EventInfo + Actor metadata.
func GetV1Registry() map[string]sharedevents.EventInfo {
	var actor = sharedevents.Actor{Service: sharedevents.ServiceDiscord, Module: "score"}
	return map[string]sharedevents.EventInfo{

		ScoreUpdateRequestDiscordV1:  {Payload: &ScoreUpdateRequestDiscordPayloadV1{}, Summary: "Score Update Request (Discord)", Description: "Request to update a score originating from Discord.", Producer: actor, Consumers: []sharedevents.Actor{actor}},
		ScoreUpdateResponseDiscordV1: {Payload: &ScoreUpdateResponseDiscordPayloadV1{}, Summary: "Score Update Response (Discord)", Description: "Confirmation that a score update succeeded (Discord).", Producer: actor, Consumers: []sharedevents.Actor{actor}},
		ScoreUpdateFailedDiscordV1:   {Payload: &ScoreUpdateFailedDiscordPayloadV1{}, Summary: "Score Update Failed (Discord)", Description: "Notification that a score update failed (Discord).", Producer: actor, Consumers: []sharedevents.Actor{actor}},

		ScoreBulkUpdateRequestDiscordV1:  {Payload: &ScoreBulkUpdateRequestDiscordPayloadV1{}, Summary: "Bulk Score Update Request (Discord)", Description: "Request to apply multiple score updates from Discord.", Producer: actor, Consumers: []sharedevents.Actor{actor}},
		ScoreBulkUpdateResponseDiscordV1: {Payload: &ScoreBulkUpdateResponseDiscordPayloadV1{}, Summary: "Bulk Score Update Response (Discord)", Description: "Result of a bulk score update operation for Discord.", Producer: actor, Consumers: []sharedevents.Actor{actor}},
	}
}
