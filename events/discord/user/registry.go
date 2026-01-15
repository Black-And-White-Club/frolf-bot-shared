package discorduserevents

import sharedevents "github.com/Black-And-White-Club/frolf-bot-shared/events/shared"

// GetV1Registry returns all modern events for the discord user functional area
// NOTE: This registry exposes ONLY the Discord-specific subjects (those that
// use the "discord." prefix) and uses shared EventInfo + Actor metadata.
func GetV1Registry() map[string]sharedevents.EventInfo {
	var actor = sharedevents.Actor{Service: sharedevents.ServiceDiscord, Module: "user"}
	return map[string]sharedevents.EventInfo{

		// Roles flow
		RoleUpdateCommandV1:       {Payload: &RoleUpdateCommandPayloadV1{}, Summary: "Role Update Command (Discord)", Description: "Admin requested a role update command via Discord.", Producer: actor, Consumers: []sharedevents.Actor{actor}},
		RoleUpdateButtonPressV1:   {Payload: &RoleUpdateButtonPressPayloadV1{}, Summary: "Role Update Button Press (Discord)", Description: "User/admin pressed a role update button in Discord.", Producer: actor, Consumers: []sharedevents.Actor{actor}},
		RoleUpdateTimeoutV1:       {Payload: &RoleUpdateTimeoutPayloadV1{}, Summary: "Role Update Timeout (Discord)", Description: "Role update flow timed out.", Producer: actor, Consumers: []sharedevents.Actor{actor}},
		RoleOptionsRequestedV1:    {Payload: &struct{}{}, Summary: "Role Options Requested (Discord)", Description: "Request to display role options in Discord.", Producer: actor, Consumers: []sharedevents.Actor{actor}},
		RoleResponseV1:            {Payload: &RoleUpdateResponsePayloadV1{}, Summary: "Role Response (Discord)", Description: "User response to role selection.", Producer: actor, Consumers: []sharedevents.Actor{actor}},
		RoleResponseFailedV1:      {Payload: &struct{}{}, Summary: "Role Response Failed (Discord)", Description: "Role response processing failed.", Producer: actor, Consumers: []sharedevents.Actor{actor}},
		RoleUpdateResponseTraceV1: {Payload: &struct{}{}, Summary: "Role Update Response Trace", Description: "Trace for role update responses.", Producer: actor, Consumers: []sharedevents.Actor{actor}},

		// Signup flow
		SignupStartedV1:             {Payload: &SignupStartedPayloadV1{}, Summary: "Signup Started (Discord)", Description: "User initiated signup via Discord.", Producer: actor, Consumers: []sharedevents.Actor{actor}},
		SignupFormSubmittedV1:       {Payload: &SignupFormSubmittedPayloadV1{}, Summary: "Signup Form Submitted (Discord)", Description: "User submitted the signup form via Discord.", Producer: actor, Consumers: []sharedevents.Actor{actor}},
		SignupSubmissionV1:          {Payload: &SignupFormSubmittedPayloadV1{}, Summary: "Signup Submission (Discord)", Description: "Signup submission for processing.", Producer: actor, Consumers: []sharedevents.Actor{actor}},
		SignupTagAskV1:              {Payload: &struct{}{}, Summary: "Signup Tag Ask (Discord)", Description: "Prompt to ask user for tag number during signup.", Producer: actor, Consumers: []sharedevents.Actor{actor}},
		SignupTagSkipV1:             {Payload: &struct{}{}, Summary: "Signup Tag Skip (Discord)", Description: "User skipped entering a tag during signup.", Producer: actor, Consumers: []sharedevents.Actor{actor}},
		SignupTagIncludeRequestedV1: {Payload: &struct{}{}, Summary: "Signup Tag Include Requested (Discord)", Description: "User requested to include a tag during signup.", Producer: actor, Consumers: []sharedevents.Actor{actor}},
		SignupTagPromptSentV1:       {Payload: &struct{}{}, Summary: "Signup Tag Prompt Sent (Discord)", Description: "Tag prompt was sent to the user.", Producer: actor, Consumers: []sharedevents.Actor{actor}},
		SignupSuccessV1:             {Payload: &SignupSuccessPayloadV1{}, Summary: "Signup Success (Discord)", Description: "Signup completed successfully.", Producer: actor, Consumers: []sharedevents.Actor{actor}},
		SignupFailedV1:              {Payload: &SignupFailedPayloadV1{}, Summary: "Signup Failed (Discord)", Description: "Signup failed with details.", Producer: actor, Consumers: []sharedevents.Actor{actor}},
		SignupCanceledV1:            {Payload: &CancelPayloadV1{}, Summary: "Signup Canceled (Discord)", Description: "User canceled signup.", Producer: actor, Consumers: []sharedevents.Actor{actor}},
		SignupAddRoleV1:             {Payload: &AddRolePayloadV1{}, Summary: "Signup Add Role (Discord)", Description: "Request to add a Discord role after signup.", Producer: actor, Consumers: []sharedevents.Actor{actor}},
		SignupRoleAddedV1:           {Payload: &RoleAddedPayloadV1{}, Summary: "Signup Role Added (Discord)", Description: "Confirmation that role was added after signup.", Producer: actor, Consumers: []sharedevents.Actor{actor}},
		SignupRoleAdditionFailedV1:  {Payload: &RoleAdditionFailedPayloadV1{}, Summary: "Signup Role Addition Failed (Discord)", Description: "Role addition failed after signup.", Producer: actor, Consumers: []sharedevents.Actor{actor}},
		SignupResponseTraceV1:       {Payload: &struct{}{}, Summary: "Signup Response Trace", Description: "Trace event for signup responses.", Producer: actor, Consumers: []sharedevents.Actor{actor}},

		// Tag flow
		TagNumberRequestedV1: {Payload: &TagNumberRequestedPayloadV1{}, Summary: "Tag Number Requested (Discord)", Description: "Request for user to provide tag number.", Producer: actor, Consumers: []sharedevents.Actor{actor}},
		TagNumberResponseV1:  {Payload: &TagNumberResponsePayloadV1{}, Summary: "Tag Number Response (Discord)", Description: "User provided tag number.", Producer: actor, Consumers: []sharedevents.Actor{actor}},
	}
}
