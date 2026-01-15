package discordroundevents

import sharedevents "github.com/Black-And-White-Club/frolf-bot-shared/events/shared"

// GetV1Registry returns all modern events for the discord round functional area
// Exposes only Discord-prefixed subjects and uses shared EventInfo + Actor metadata.
func GetV1Registry() map[string]sharedevents.EventInfo {
	var actor = sharedevents.Actor{Service: sharedevents.ServiceDiscord, Module: "round"}

	return map[string]sharedevents.EventInfo{
		// Creation flow
		RoundCreateModalSubmittedV1:    {Payload: &CreateRoundModalPayloadV1{}, Summary: "Round Create Modal Submitted", Description: "User submitted a round creation modal via Discord.", Producer: actor, Consumers: []sharedevents.Actor{actor}},
		RoundCreatedDiscordV1:          {Payload: &RoundCreatedDiscordPayloadV1{}, Summary: "Round Created (Discord)", Description: "Notification that a round was created; used to post embeds in Discord.", Producer: actor, Consumers: []sharedevents.Actor{actor}},
		RoundCreationFailedDiscordV1:   {Payload: &RoundCreationFailedDiscordPayloadV1{}, Summary: "Round Creation Failed (Discord)", Description: "Notification that round creation failed, sent to the requester.", Producer: actor, Consumers: []sharedevents.Actor{actor}},
		RoundValidationFailedDiscordV1: {Payload: &RoundCreationFailedDiscordPayloadV1{}, Summary: "Round Validation Failed (Discord)", Description: "Validation failure details for a Discord round creation attempt.", Producer: actor, Consumers: []sharedevents.Actor{actor}},
		RoundCreatedTraceV1:            {Payload: &RoundCreatedTracePayloadV1{}, Summary: "Round Created Trace", Description: "Tracing information emitted when a round is created (Discord).", Producer: actor, Consumers: []sharedevents.Actor{actor}},

		// Lifecycle
		RoundStartedDiscordV1:         {Payload: &RoundStartedDiscordPayloadV1{}, Summary: "Round Started (Discord)", Description: "Notify Discord that the round has started so embeds can be updated.", Producer: actor, Consumers: []sharedevents.Actor{actor}},
		RoundFinalizedDiscordNotifyV1: {Payload: &RoundFinalizedDiscordPayloadV1{}, Summary: "Round Finalized (Discord)", Description: "Notify Discord that the round has been finalized and results are available.", Producer: actor, Consumers: []sharedevents.Actor{actor}},
		RoundReminderDiscordV1:        {Payload: &RoundReminderDiscordPayloadV1{}, Summary: "Round Reminder (Discord)", Description: "Send a reminder to round participants via Discord.", Producer: actor, Consumers: []sharedevents.Actor{actor}},

		// Participants
		RoundParticipantJoinRequestDiscordV1: {Payload: &RoundParticipantJoinRequestDiscordPayloadV1{}, Summary: "Round Participant Join Request (Discord)", Description: "User requested to join a round via Discord interaction.", Producer: actor, Consumers: []sharedevents.Actor{actor}},
		RoundParticipantJoinedDiscordV1:      {Payload: &RoundParticipantJoinedDiscordPayloadV1{}, Summary: "Round Participant Joined (Discord)", Description: "Notification that a participant joined the round, for Discord embed updates.", Producer: actor, Consumers: []sharedevents.Actor{actor}},

		// Scoring
		RoundScoreUpdateRequestDiscordV1:      {Payload: &RoundScoreUpdateRequestDiscordPayloadV1{}, Summary: "Round Score Update Request (Discord)", Description: "User submitted a score update via Discord.", Producer: actor, Consumers: []sharedevents.Actor{actor}},
		RoundParticipantScoreUpdatedDiscordV1: {Payload: &RoundParticipantScoreUpdatedDiscordPayloadV1{}, Summary: "Round Participant Score Updated (Discord)", Description: "Notification that a participant's score was updated, for Discord embed updates.", Producer: actor, Consumers: []sharedevents.Actor{actor}},

		// Update / Delete
		RoundUpdateModalSubmittedV1: {Payload: &RoundUpdateModalSubmittedPayloadV1{}, Summary: "Round Update Modal Submitted (Discord)", Description: "User submitted a round update modal via Discord.", Producer: actor, Consumers: []sharedevents.Actor{actor}},
		RoundUpdateRequestDiscordV1: {Payload: &RoundUpdateRequestDiscordPayloadV1{}, Summary: "Round Update Request (Discord)", Description: "Request to update a round originating from Discord.", Producer: actor, Consumers: []sharedevents.Actor{actor}},
		RoundUpdatedDiscordV1:       {Payload: &RoundUpdatedDiscordPayloadV1{}, Summary: "Round Updated (Discord)", Description: "Notification that a round was updated; used to refresh Discord embeds.", Producer: actor, Consumers: []sharedevents.Actor{actor}},
		RoundDeleteRequestDiscordV1: {Payload: &RoundDeleteRequestDiscordPayloadV1{}, Summary: "Round Delete Request (Discord)", Description: "Request to delete a round originating from Discord.", Producer: actor, Consumers: []sharedevents.Actor{actor}},
		RoundDeletedDiscordV1:       {Payload: &RoundDeletedDiscordPayloadV1{}, Summary: "Round Deleted (Discord)", Description: "Notification that a round was deleted; used to remove or mark Discord embeds.", Producer: actor, Consumers: []sharedevents.Actor{actor}},
	}
}
