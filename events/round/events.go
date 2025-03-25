package roundevents

import (
	roundtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/round"
	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
)

// Stream names
const (
	RoundStreamName       = "round"
	UserStreamName        = "user"
	LeaderboardStreamName = "leaderboard"
	ScoreStreamName       = "score"
)

// Event names - grouped by functionality
const (
	// Round Creation Events
	RoundCreateRequest         = "round.create.request"
	RoundValidated             = "round.validated"
	RoundValidationFailed      = "discord.round.validation.failed"
	RoundDateTimeParsed        = "round.datetime.parsed"
	RoundEntityCreated         = "round.entity.created"
	RoundStored                = "round.stored"
	RoundScheduled             = "round.scheduled"
	RoundCreated               = "discord.round.created"
	RoundCreationFailed        = "discord.round.creation.failed"
	RoundError                 = "round.error"
	RoundCreatedEvent          = "discord.round.created"
	RoundEventMessageIDUpdated = "round.discord.event.id.updated"
	RoundEventMessageIDUpdate  = "round.discord.event.id.update"
	RoundTraceEvent            = "round.trace.event"

	// Round Update Events
	RoundUpdateRequest    = "round.update.request"
	RoundUpdateValidated  = "round.update.validated"
	RoundFetched          = "round.fetched"
	RoundEntityUpdated    = "round.entity.updated"
	RoundUpdateSuccess    = "round.update.success"
	RoundUpdateError      = "round.update.error"
	RoundUpdated          = "round.updated"
	RoundScheduleUpdate   = "round.schedule.update"
	RoundStateUpdated     = "round.state.updated"
	RoundsUpdated         = "round.rounds.updated"
	RoundUpdateReschedule = "round.update.reschedule"

	// Round Delete Events
	RoundDeleteRequest      = "round.delete.request"
	RoundDeleted            = "discord.round.deleted"
	RoundDeleteValidated    = "round.delete.validated"
	RoundToDeleteFetched    = "round.to.delete.fetched"
	RoundDeleteAuthorized   = "round.delete.authorized"
	RoundDeleteUnauthorized = "round.delete.unauthorized"
	RoundDeleteError        = "round.delete.error"

	// Participant Events
	RoundParticipantJoinRequest           = "round.participant.join.request"
	RoundParticipantJoinValidated         = "round.participant.join.validated"
	RoundParticipantJoinError             = "round.participant.join.error"
	RoundParticipantDeclined              = "round.participant.declined"
	RoundParticipantDeclinedResponse      = "discord.round.participant.declined"
	RoundParticipantJoined                = "discord.round.participant.joined"
	RoundParticipantRemovalRequest        = "round.participant.removal.request"
	RoundParticipantRemoved               = "discord.round.participant.removed"
	RoundParticipantJoinValidationRequest = "round.participant.join.validation.request"
	RoundParticipantStatusError           = "round.participant.error"
	RoundParticipantStatusFound           = "round.participant.found"

	// Score Events
	RoundScoreUpdateRequest      = "round.score.update.request"
	RoundScoreUpdateValidated    = "round.score.update.validated"
	RoundParticipantScoreUpdated = "round.participant.score.updated"
	RoundAllScoresSubmitted      = "round.all.scores.submitted"
	RoundNotAllScoresSubmitted   = "discord.round.not.all.scores.submitted"
	RoundScoreUpdateError        = "discord.round.score.update.error"
	ProcessRoundScoresRequest    = "score.process.scores.request"
	ScoreModuleNotificationError = "score.module.notification.error"

	// Round Lifecycle Events
	RoundStarted            = "round.started"
	DiscordRoundStarted     = "discord.round.started"
	RoundFinalized          = "round.finalized"
	RoundFinalizationError  = "round.finalization.error"
	RoundScoresNotification = "round.scores.notification"
	RoundReminder           = "round.reminder"
	DiscordRoundReminder    = "discord.round.reminder"
	DiscordRoundFinalized   = "discord.round.finalized"

	// Tag Events
	RoundTagNumberRequest           = "round.tag.number.request"
	RoundTagNumberFound             = "round.tag.number.found"
	RoundTagNumberNotFound          = "round.tag.number.not.found"
	LeaderboardGetTagNumberRequest  = "leaderboard.get.tag.number.request"
	LeaderboardGetTagNumberResponse = "round.get.tag.number.response"

	// Authorization Events
	RoundUserRoleCheckRequest = "round.user.role.check.request"
	RoundUserRoleCheckResult  = "round.user.role.check.result"
	RoundUserRoleCheckError   = "round.user.role.check.error"

	// Discord Events
	DiscordEventsSubject   = "discord.round.event"
	DelayedMessagesSubject = "delayed.messages"
)

// Event Payloads - structured by extending base payloads where possible

type CreateRoundRequestedPayload struct {
	Title       roundtypes.Title       `json:"title"`
	Description roundtypes.Description `json:"description"`
	StartTime   string                 `json:"start_time"`
	Location    roundtypes.Location    `json:"location"`
	UserID      sharedtypes.DiscordID  `json:"user_id"`
	ChannelID   string                 `json:"channel_id"`
	Timezone    roundtypes.Timezone    `json:"timezone"`
}

type RoundCreateRequestPayload struct {
	roundtypes.BaseRoundPayload
	Timezone roundtypes.Timezone `json:"timezone"`
}

type RoundDateTimeParsedPayload struct {
	RoundCreateRequestPayload RoundCreateRequestPayload `json:"round_create_request_payload"`
	StartTime                 *roundtypes.StartTime     `json:"start_time"`
}

type RoundScheduledPayload struct {
	roundtypes.BaseRoundPayload
	EventMessageID *roundtypes.EventMessageID `json:"discord_message_id"`
}

type RoundStartedPayload struct {
	RoundID   sharedtypes.RoundID   `json:"round_id"`
	Title     roundtypes.Title      `json:"title"`
	Location  *roundtypes.Location  `json:"location"`
	StartTime *roundtypes.StartTime `json:"start_time"`
	ChannelID string                `json:"channel_id"`
}

type RoundFinalizedEmbedUpdatePayload struct {
	RoundID          sharedtypes.RoundID
	Title            roundtypes.Title
	StartTime        *roundtypes.StartTime
	Location         *roundtypes.Location
	Participants     []roundtypes.Participant
	EventMessageID   *roundtypes.EventMessageID
	DiscordChannelID string `json:"discord_channel_id,omitempty"`
}

// ---- Round Creation Payloads ----

type RoundValidatedPayload struct {
	CreateRoundRequestedPayload CreateRoundRequestedPayload `json:"round_create_request_payload"`
}

type RoundValidationFailedPayload struct {
	UserID       sharedtypes.DiscordID `json:"user_id"`
	ErrorMessage []string              `json:"error_messages"`
}
type RoundEntityCreatedPayload struct {
	Round            roundtypes.Round `json:"round"`
	DiscordChannelID string           `json:"discord_channel_id"`
	DiscordGuildID   string           `json:"discord_guild_id"`
}

type RoundStoredPayload struct {
	Round roundtypes.Round `json:"round"`
}

type RoundCreatedPayload struct {
	roundtypes.BaseRoundPayload
	ChannelID string `json:"channel_id"`
}

type RoundCreationFailedPayload struct {
	UserID       sharedtypes.DiscordID `json:"user_id"`
	ErrorMessage string                `json:"error_message"`
	ChannelID    string                `json:"channel_id"`
	GuildID      string                `json:"guild_id"`
}

type RoundEventCreatedPayload struct {
	RoundID        sharedtypes.RoundID `json:"round_id"`
	EventMessageID string              `json:"discord_message_id"`
}

type RoundEventMessageIDUpdatedPayload struct {
	RoundID        sharedtypes.RoundID `json:"round_id"`
	EventMessageID string              `json:"discord_message_id"`
}

type RoundErrorPayload struct {
	Round *RoundCreateRequestPayload `json:"round"`
	Error string                     `json:"error"`
}

// ---- Round Update Payloads ----

type RoundUpdateRequestPayload struct {
	roundtypes.BaseRoundPayload
	EventType *roundtypes.EventType `json:"event_type,omitempty"`
}

type RoundUpdateValidatedPayload struct {
	RoundUpdateRequestPayload RoundUpdateRequestPayload `json:"round_update_request_payload"`
}

type RoundFetchedPayload struct {
	Round                     roundtypes.Round          `json:"round"`
	RoundUpdateRequestPayload RoundUpdateRequestPayload `json:"round_update_request_payload"`
}

type RoundEntityUpdatedPayload struct {
	Round roundtypes.Round `json:"round"`
}

type RoundUpdateSuccessPayload struct {
	RoundID sharedtypes.RoundID `json:"round_id"`
}

type RoundUpdateErrorPayload struct {
	RoundUpdateRequest *RoundUpdateRequestPayload `json:"round_update_request"`
	Error              string                     `json:"error"`
}

type RoundScheduleUpdatePayload struct {
	RoundID   sharedtypes.RoundID   `json:"round_id"`
	Title     roundtypes.Title      `json:"title"`
	StartTime *roundtypes.StartTime `json:"start_time"`
	Location  *roundtypes.Location  `json:"location"`
}

type RoundStateUpdatedPayload struct {
	RoundID sharedtypes.RoundID   `json:"round_id"`
	State   roundtypes.RoundState `json:"state"`
}

// ---- Round Delete Payloads ----

type RoundDeleteRequestPayload struct {
	RoundID              sharedtypes.RoundID   `json:"round_id" validate:"required"`
	RequestingUserUserID sharedtypes.DiscordID `json:"requesting_user_user_id" validate:"required"`
}

type RoundDeleteValidatedPayload struct {
	RoundDeleteRequestPayload RoundDeleteRequestPayload `json:"round_delete_request_payload"`
}

type RoundToDeleteFetchedPayload struct {
	Round                     roundtypes.Round          `json:"round"`
	RoundDeleteRequestPayload RoundDeleteRequestPayload `json:"round_delete_request_payload"`
}

type RoundDeleteAuthorizedPayload struct {
	RoundID sharedtypes.RoundID `json:"round_id"`
}

type RoundDeletedPayload struct {
	RoundID        sharedtypes.RoundID       `json:"round_id"`
	EventMessageID roundtypes.EventMessageID `json:"discord_message_id"`
}

type RoundDeleteErrorPayload struct {
	RoundDeleteRequest *RoundDeleteRequestPayload `json:"round_delete_request"`
	Error              string                     `json:"error"`
}

// ---- Participant Payloads ----

type RoundParticipant struct {
	UserID    sharedtypes.DiscordID  `json:"user_id"`
	TagNumber *sharedtypes.TagNumber `json:"tag_number"`
	Response  roundtypes.Response    `json:"response"`
	Score     *sharedtypes.Score     `json:"score"`
}

type ParticipantJoinRequestPayload struct {
	RoundID    sharedtypes.RoundID    `json:"round_id"`
	UserID     sharedtypes.DiscordID  `json:"user_id"`
	Response   roundtypes.Response    `json:"response"`
	TagNumber  *sharedtypes.TagNumber `json:"tag_number"`
	JoinedLate *bool                  `json:"joined_late,omitempty"`
}

type ParticipantJoinValidatedPayload struct {
	ParticipantJoinRequestPayload ParticipantJoinRequestPayload `json:"participant_join_request_payload"`
}

type RoundParticipantJoinErrorPayload struct {
	ParticipantJoinRequest *ParticipantJoinRequestPayload `json:"participant_join_request"`
	Error                  string                         `json:"error"`
	EventMessageID         roundtypes.EventMessageID      `json:"discord_message_id"`
}

type ParticipantStatusRequestPayload struct {
	RoundID sharedtypes.RoundID   `json:"round_id"`
	UserID  sharedtypes.DiscordID `json:"user_id"`
}

type ParticipantStatusFoundPayload struct {
	RoundID sharedtypes.RoundID   `json:"round_id"`
	UserID  sharedtypes.DiscordID `json:"user_id"`
	Status  string                `json:"status"`
}

type ParticipantRemovalRequestPayload struct {
	RoundID sharedtypes.RoundID   `json:"round_id"`
	UserID  sharedtypes.DiscordID `json:"user_id"`
}

type ParticipantRemovedPayload struct {
	RoundID        sharedtypes.RoundID       `json:"round_id"`
	UserID         sharedtypes.DiscordID     `json:"user_id"`
	Response       roundtypes.Response       `json:"response"`
	EventMessageID roundtypes.EventMessageID `json:"discord_message_id"`
}

type ParticipantJoinValidationRequestPayload struct {
	RoundID  sharedtypes.RoundID   `json:"round_id"`
	UserID   sharedtypes.DiscordID `json:"user_id"`
	Response roundtypes.Response   `json:"response"`
}

type ParticipantJoinedPayload struct {
	RoundID               sharedtypes.RoundID       `json:"round_id"`
	AcceptedParticipants  []roundtypes.Participant  `jsonb:"accepted_participants"`
	DeclinedParticipants  []roundtypes.Participant  `jsonb:"declined_participants"`
	TentativeParticipants []roundtypes.Participant  `jsonb:"tentative_participants"`
	EventMessageID        roundtypes.EventMessageID `json:"discord_message_id"`
	JoinedLate            *bool                     `json:"joined_late,omitempty"`
}

type ParticipantDeclinedPayload struct {
	RoundID        sharedtypes.RoundID       `json:"round_id"`
	UserID         sharedtypes.DiscordID     `json:"user_id"`
	EventMessageID roundtypes.EventMessageID `json:"discord_message_id"`
}

// ---- Score Payloads ----

type ScoreUpdateRequestPayload struct {
	RoundID     sharedtypes.RoundID   `json:"round_id"`
	Participant sharedtypes.DiscordID `json:"participant"`
	Score       *sharedtypes.Score    `json:"score"`
}

type ScoreUpdateValidatedPayload struct {
	ScoreUpdateRequestPayload ScoreUpdateRequestPayload `json:"score_update_request_payload"`
}

type ParticipantScoreUpdatedPayload struct {
	RoundID        sharedtypes.RoundID        `json:"round_id"`
	Participant    sharedtypes.DiscordID      `json:"participant"`
	Score          sharedtypes.Score          `json:"score"`
	ChannelID      string                     `json:"channel_id"`
	EventMessageID *roundtypes.EventMessageID `json:"discord_message_id"`
}

type AllScoresSubmittedPayload struct {
	RoundID        sharedtypes.RoundID        `json:"round_id"`
	EventMessageID *roundtypes.EventMessageID `json:"discord_message_id"`
}

type RoundScoreUpdateErrorPayload struct {
	ScoreUpdateRequest *ScoreUpdateRequestPayload `json:"score_update_request"`
	Error              string                     `json:"error"`
}

type ParticipantScore struct {
	UserID    sharedtypes.DiscordID  `json:"user_id"`
	TagNumber *sharedtypes.TagNumber `json:"tag_number"`
	Score     sharedtypes.Score      `json:"score"`
}

type RoundScoresNotificationPayload struct {
	RoundID sharedtypes.RoundID `json:"round_id"`
	Scores  []ParticipantScore  `json:"scores"`
}

type ProcessRoundScoresRequestPayload struct {
	RoundID sharedtypes.RoundID `json:"round_id"`
	Scores  []ParticipantScore  `json:"scores"`
}

// ---- Round Lifecycle Payloads ----

type RoundFinalizedPayload struct {
	RoundID sharedtypes.RoundID `json:"round_id"`
}

type RoundFinalizationErrorPayload struct {
	RoundID sharedtypes.RoundID `json:"round_id"`
	Error   string              `json:"error"`
}

type ScoreModuleNotificationErrorPayload struct {
	RoundID sharedtypes.RoundID `json:"round_id"`
	Error   string              `json:"error"`
}

// Discord-related payloads
type DiscordReminderPayload struct {
	RoundID          sharedtypes.RoundID       `json:"round_id"`
	ReminderType     string                    `json:"reminder_type"`
	RoundTitle       roundtypes.Title          `json:"round_title"`
	StartTime        *roundtypes.StartTime     `json:"start_time"`
	Location         *roundtypes.Location      `json:"location"`
	UserIDs          []sharedtypes.DiscordID   `json:"user_ids"`
	DiscordChannelID string                    `json:"discord_channel_id,omitempty"`
	DiscordGuildID   string                    `json:"discord_guild_id,omitempty"`
	EventMessageID   roundtypes.EventMessageID `json:"event_message_id"`
}

type DiscordRoundStartPayload struct {
	RoundID          sharedtypes.RoundID       `json:"round_id"`
	Title            roundtypes.Title          `json:"title"`
	Location         *roundtypes.Location      `json:"location"`
	StartTime        *roundtypes.StartTime     `json:"start_time"`
	Participants     []RoundParticipant        `jsonb:"participants"`
	DiscordChannelID string                    `json:"discord_channel_id,omitempty"`
	DiscordGuildID   string                    `json:"discord_guild_id,omitempty"`
	EventMessageID   roundtypes.EventMessageID `json:"event_message_id"`
}

type DiscordRoundParticipant struct {
	UserID    sharedtypes.DiscordID  `json:"user_id"`
	TagNumber *sharedtypes.TagNumber `json:"tag_number"`
	Score     *sharedtypes.Score     `json:"score"`
}

// ---- Tag Retrieval Payloads ----

type TagNumberRequestPayload struct {
	UserID  sharedtypes.DiscordID `json:"user_id"`
	RoundID sharedtypes.RoundID   `json:"round_id"`
}

type RoundTagNumberFoundPayload struct {
	RoundID   sharedtypes.RoundID    `json:"round_id"`
	UserID    sharedtypes.DiscordID  `json:"user_id"`
	TagNumber *sharedtypes.TagNumber `json:"tag_number"`
}

type RoundTagNumberNotFoundPayload struct {
	UserID  sharedtypes.DiscordID `json:"user_id"`
	RoundID sharedtypes.RoundID   `json:"round_id"`
}

type GetTagNumberResponsePayload struct {
	UserID    sharedtypes.DiscordID  `json:"user_id"`
	TagNumber *sharedtypes.TagNumber `json:"tag_number"`
	Error     string                 `json:"error,omitempty"`
	RoundID   sharedtypes.RoundID    `json:"round_id"`
}

// ---- Authorization Payloads ----

type UserRoleCheckRequestPayload struct {
	UserID  sharedtypes.DiscordID `json:"user_id"`
	RoundID sharedtypes.RoundID   `json:"round_id"`
}

type UserRoleCheckResultPayload struct {
	UserID  sharedtypes.DiscordID `json:"user_id"`
	RoundID sharedtypes.RoundID   `json:"round_id"`
	HasRole bool                  `json:"has_role"`
	Error   string                `json:"error,omitempty"`
}

type RoundUserRoleCheckErrorPayload struct {
	UserID  sharedtypes.DiscordID `json:"user_id"`
	RoundID sharedtypes.RoundID   `json:"round_id"`
	Error   string                `json:"error"`
}

func (p *AllScoresSubmittedPayload) ToRoundFinalizedPayload() RoundFinalizedPayload {
	return RoundFinalizedPayload{
		RoundID: p.RoundID,
	}
}
