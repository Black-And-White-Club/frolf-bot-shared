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
	DiscordStreamName     = "discord"
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
	RoundScheduleUpdate   = "discord.round.schedule.update"
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
	RoundParticipantDeclined              = "discord.round.participant.declined"
	RoundParticipantDeclinedResponse      = "discord.round.participant.declined"
	RoundParticipantJoined                = "discord.round.participant.joined"
	RoundParticipantRemovalRequest        = "round.participant.removal.request"
	RoundParticipantRemoved               = "discord.round.participant.removed"
	RoundParticipantJoinValidationRequest = "round.participant.join.validation.request"
	RoundParticipantStatusError           = "round.participant.error"
	RoundParticipantStatusFound           = "round.participant.found"
	RoundParticipantStatusCheckError      = "discord.round.participant.status.check.error"
	RoundParticipantRemovalError          = "round.participant.removal.error"
	RoundParticipantUpdateError           = "round.participant.update.error"
	RoundParticipantStatusUpdateRequest   = "round.participant.status.update.request"

	// Score Events
	RoundScoreUpdateRequest        = "round.score.update.request"
	RoundScoreUpdateValidated      = "round.score.update.validated"
	RoundParticipantScoreUpdated   = "round.participant.score.updated"
	DiscordParticipantScoreUpdated = "discord.participant.score.updated"
	RoundAllScoresSubmitted        = "round.all.scores.submitted"
	RoundNotAllScoresSubmitted     = "discord.round.not.all.scores.submitted"
	RoundScoreUpdateError          = "discord.round.score.update.error"
	ProcessRoundScoresRequest      = "score.process.round.scores.request"
	ScoreModuleNotificationError   = "score.module.notification.error"

	// Round Lifecycle Events
	RoundStarted            = "round.started"
	DiscordRoundStarted     = "discord.round.started"
	RoundFinalized          = "round.finalized"
	RoundFinalizationError  = "round.finalization.error"
	RoundScoresNotification = "round.scores.notification"
	RoundReminder           = "round.reminder"
	DiscordRoundReminder    = "discord.reminder"
	DiscordRoundFinalized   = "discord.round.finalized"

	// Tag Events
	RoundTagNumberRequest          = "round.tag.number.request"
	RoundTagNumberFound            = "round.leaderboard.tag.found"
	RoundTagNumberNotFound         = "round.leaderboard.tag.not.found"
	LeaderboardGetTagNumberRequest = "leaderboard.round.tag.get.by.user.id.request"

	LeaderboardGetTagNumberResponse = "round.get.tag.number.response"
	TagUpdateForScheduledRounds     = "round.tag.update.for.scheduled.rounds"
	TagsUpdatedForScheduledRounds   = "discord.tags.updated.for.scheduled.rounds"

	// Authorization Events
	RoundUserRoleCheckRequest = "round.user.role.check.request"
	RoundUserRoleCheckResult  = "round.user.role.check.result"
	RoundUserRoleCheckError   = "round.user.role.check.error"

	// Discord Events
	DiscordEventsSubject   = "discord.round.event"
	DelayedMessagesSubject = "delayed.messages"

	// Round Retrieval Events
	GetRoundRequest = "round.get.request"
	RoundRetrieved  = "discord.round.retrieved"
)

// Event Payloads - structured by extending base payloads where possible

type GetRoundRequestPayload struct {
	RoundID        sharedtypes.RoundID   `json:"round_id"`
	EventMessageID string                `json:"event_message_id"`
	UserID         sharedtypes.DiscordID `json:"user_id"`
}

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
	StartTime                 *sharedtypes.StartTime    `json:"start_time"`
}

type RoundScheduledPayload struct {
	roundtypes.BaseRoundPayload
	EventMessageID string `json:"discord_message_id"`
}

type RoundStartedPayload struct {
	RoundID   sharedtypes.RoundID    `json:"round_id"`
	Title     roundtypes.Title       `json:"title"`
	Location  *roundtypes.Location   `json:"location"`
	StartTime *sharedtypes.StartTime `json:"start_time"`
	ChannelID string                 `json:"channel_id"`
}

type TagLookupRequestPayload struct {
	UserID           sharedtypes.DiscordID `json:"user_id"`
	RoundID          sharedtypes.RoundID   `json:"round_id"`
	Response         roundtypes.Response   `json:"response"`          // ← Current/Active response
	OriginalResponse roundtypes.Response   `json:"original_response"` // ← Same as Response (for consistency)
	JoinedLate       *bool                 `json:"joined_late,omitempty"`
}
type RoundFinalizedEmbedUpdatePayload struct {
	RoundID          sharedtypes.RoundID      `json:"round_id"`
	Title            roundtypes.Title         `json:"title"`
	StartTime        *sharedtypes.StartTime   `json:"start_time"`
	Location         *roundtypes.Location     `json:"location"`
	Participants     []roundtypes.Participant `json:"participants"`
	EventMessageID   string                   `json:"discord_message_id"`
	DiscordChannelID string                   `json:"discord_channel_id,omitempty"`
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

type RoundMessageIDUpdatePayload struct {
	RoundID sharedtypes.RoundID `json:"round_id"`
}
type RoundErrorPayload struct {
	RoundID sharedtypes.RoundID `json:"round"`
	Error   string              `json:"error"`
}

// ---- Round Update Payloads ----

type UpdateRoundRequestedPayload struct {
	RoundID     sharedtypes.RoundID     `json:"round_id"`
	UserID      sharedtypes.DiscordID   `json:"user_id"`
	ChannelID   string                  `json:"channel_id"`
	MessageID   string                  `json:"message_id"`
	Title       *roundtypes.Title       `json:"title,omitempty"`
	Description *roundtypes.Description `json:"description,omitempty"`
	StartTime   *string                 `json:"start_time,omitempty"`
	Timezone    *roundtypes.Timezone    `json:"timezone"`
	Location    *roundtypes.Location    `json:"location,omitempty"`
}

type RoundUpdateRequestPayload struct {
	RoundID     sharedtypes.RoundID     `json:"round_id"`
	Title       roundtypes.Title        `json:"title,omitempty"`
	Description *roundtypes.Description `json:"description,omitempty"`
	Location    *roundtypes.Location    `json:"location,omitempty"`
	StartTime   *sharedtypes.StartTime  `json:"start_time,omitempty"`
	EventType   *roundtypes.EventType   `json:"event_type,omitempty"`
	UserID      sharedtypes.DiscordID   `json:"user_id"`
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
	RoundID   sharedtypes.RoundID    `json:"round_id"`
	Title     roundtypes.Title       `json:"title"`
	StartTime *sharedtypes.StartTime `json:"start_time"`
	Location  *roundtypes.Location   `json:"location"`
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
	RoundID        sharedtypes.RoundID `json:"round_id"`
	EventMessageID string              `json:"discord_message_id"`
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
	EventMessageID         string                         `json:"discord_message_id"`
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
	RoundID               sharedtypes.RoundID      `json:"round_id"`
	UserID                sharedtypes.DiscordID    `json:"user_id"`
	AcceptedParticipants  []roundtypes.Participant `json:"accepted_participants"`
	DeclinedParticipants  []roundtypes.Participant `json:"declined_participants"`
	TentativeParticipants []roundtypes.Participant `json:"tentative_participants"`
	EventMessageID        string                   `json:"discord_message_id"`
}

type ParticipantJoinValidationRequestPayload struct {
	RoundID  sharedtypes.RoundID   `json:"round_id"`
	UserID   sharedtypes.DiscordID `json:"user_id"`
	Response roundtypes.Response   `json:"response"`
}

type ParticipantJoinedPayload struct {
	RoundID               sharedtypes.RoundID      `json:"round_id"`
	AcceptedParticipants  []roundtypes.Participant `jsonb:"accepted_participants"`
	DeclinedParticipants  []roundtypes.Participant `jsonb:"declined_participants"`
	TentativeParticipants []roundtypes.Participant `jsonb:"tentative_participants"`
	EventMessageID        string                   `json:"discord_message_id"`
	JoinedLate            *bool                    `json:"joined_late,omitempty"`
}

type ParticipantDeclinedPayload struct {
	RoundID        sharedtypes.RoundID   `json:"round_id"`
	UserID         sharedtypes.DiscordID `json:"user_id"`
	EventMessageID string                `json:"discord_message_id"`
}

// ParticipantStatusCheckErrorPayload holds data for errors during participant status check.
type ParticipantStatusCheckErrorPayload struct {
	RoundID sharedtypes.RoundID   `json:"round_id"`
	UserID  sharedtypes.DiscordID `json:"user_id"`
	Error   string                `json:"error"`
}

// ParticipantUpdateErrorPayload holds data for errors during participant DB updates.
type ParticipantUpdateErrorPayload struct {
	RoundID sharedtypes.RoundID   `json:"round_id"`
	UserID  sharedtypes.DiscordID `json:"user_id"`
	Error   string                `json:"error"`
}

// ParticipantRemovalErrorPayload holds data for errors during participant removal.
type ParticipantRemovalErrorPayload struct {
	RoundID sharedtypes.RoundID   `json:"round_id"`
	UserID  sharedtypes.DiscordID `json:"user_id"`
	Error   string                `json:"error"`
}

// ParticipantDeclineErrorPayload holds data for errors during decline handling.
type ParticipantDeclineErrorPayload struct {
	RoundID sharedtypes.RoundID   `json:"round_id"`
	UserID  sharedtypes.DiscordID `json:"user_id"`
	Error   string                `json:"error"`
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
	RoundID        sharedtypes.RoundID      `json:"round_id"`
	Participant    sharedtypes.DiscordID    `json:"participant"`
	Score          sharedtypes.Score        `json:"score"`
	ChannelID      string                   `json:"channel_id"`
	EventMessageID string                   `json:"discord_message_id"`
	Participants   []roundtypes.Participant `json:"participants,omitempty"`
}

type AllScoresSubmittedPayload struct {
	RoundID        sharedtypes.RoundID      `json:"round_id"`
	EventMessageID string                   `json:"discord_message_id"`
	RoundData      roundtypes.Round         `json:"round_data"`
	Participants   []roundtypes.Participant `json:"participants,omitempty"`
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
	RoundID   sharedtypes.RoundID `json:"round_id"`
	RoundData roundtypes.Round    `json:"round_data"`
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
	RoundID          sharedtypes.RoundID     `json:"round_id"`
	ReminderType     string                  `json:"reminder_type"`
	RoundTitle       roundtypes.Title        `json:"round_title"`
	StartTime        *sharedtypes.StartTime  `json:"start_time"`
	Location         *roundtypes.Location    `json:"location"`
	UserIDs          []sharedtypes.DiscordID `json:"user_ids"`
	DiscordChannelID string                  `json:"discord_channel_id,omitempty"`
	DiscordGuildID   string                  `json:"discord_guild_id,omitempty"`
	EventMessageID   string                  `json:"event_message_id"`
}

type RoundReminderProcessedPayload struct {
	RoundID sharedtypes.RoundID `json:"round_id"`
}

type DiscordRoundStartPayload struct {
	RoundID          sharedtypes.RoundID    `json:"round_id"`
	Title            roundtypes.Title       `json:"title"`
	Location         *roundtypes.Location   `json:"location"`
	StartTime        *sharedtypes.StartTime `json:"start_time"`
	Participants     []RoundParticipant     `jsonb:"participants"`
	DiscordChannelID string                 `json:"discord_channel_id,omitempty"`
	DiscordGuildID   string                 `json:"discord_guild_id,omitempty"`
	EventMessageID   string                 `json:"event_message_id"`
}

type DiscordRoundParticipant struct {
	UserID    sharedtypes.DiscordID  `json:"user_id"`
	TagNumber *sharedtypes.TagNumber `json:"tag_number"`
	Score     *sharedtypes.Score     `json:"score"`
}

type DiscordRoundUpdatePayload struct {
	Participants    []roundtypes.Participant `json:"participants"`
	RoundIDs        []sharedtypes.RoundID    `json:"round_ids"`
	EventMessageIDs []string                 `json:"event_message_ids"`
}

// ---- Tag Retrieval Payloads ----

type TagNumberRequestPayload struct {
	UserID  sharedtypes.DiscordID `json:"user_id"`
	RoundID sharedtypes.RoundID   `json:"round_id"`
}

type RoundTagNumberFoundPayload struct {
	RoundID            sharedtypes.RoundID    `json:"round_id"`
	UserID             sharedtypes.DiscordID  `json:"user_id"`
	TagNumber          *sharedtypes.TagNumber `json:"tag_number"`
	OriginalResponse   roundtypes.Response    `json:"original_response"`
	OriginalJoinedLate *bool                  `json:"original_joined_late,omitempty"`
}

type RoundTagNumberNotFoundPayload struct {
	RoundID            sharedtypes.RoundID   `json:"round_id"`
	UserID             sharedtypes.DiscordID `json:"user_id"`
	OriginalResponse   roundtypes.Response   `json:"original_response"`
	OriginalJoinedLate *bool                 `json:"original_joined_late,omitempty"`
	Reason             string                `json:"reason,omitempty"`
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

func (p *AllScoresSubmittedPayload) ToRoundFinalizedPayload(round roundtypes.Round) RoundFinalizedPayload {
	return RoundFinalizedPayload{
		RoundID:   p.RoundID,
		RoundData: round, // Pass the actual round instance here
	}
}

// ScheduledRoundTagUpdatePayload represents the payload for updating scheduled round tags.
type ScheduledRoundTagUpdatePayload struct {
	ChangedTags map[sharedtypes.DiscordID]*sharedtypes.TagNumber
}

type NotAllScoresSubmittedPayload struct {
	RoundID        sharedtypes.RoundID      `json:"round_id"`
	Participant    sharedtypes.DiscordID    `json:"participant"`
	Score          sharedtypes.Score        `json:"score"`
	EventMessageID string                   `json:"event_message_id"`
	Scores         []ParticipantScore       `json:"scores"`
	Participants   []roundtypes.Participant `json:"participants"`
}

// ParticipantUpdatePayload defines the common interface for participant update payloads.
type ParticipantUpdatePayload interface {
	GetRoundID() sharedtypes.RoundID
	GetUserID() sharedtypes.DiscordID
	GetTagNumber() *sharedtypes.TagNumber
	GetJoinedLate() *bool
}

type RoundUpdateInfo struct {
	RoundID             sharedtypes.RoundID      `json:"round_id"`
	EventMessageID      string                   `json:"event_message_id"`
	Title               roundtypes.Title         `json:"title"`
	StartTime           *sharedtypes.StartTime   `json:"start_time"`
	Location            *roundtypes.Location     `json:"location"`
	UpdatedParticipants []roundtypes.Participant `json:"updated_participants"`
	ParticipantsChanged int                      `json:"participants_changed"`
}

type UpdateSummary struct {
	TotalRoundsProcessed int `json:"total_rounds_processed"`
	RoundsUpdated        int `json:"rounds_updated"`
	ParticipantsUpdated  int `json:"participants_updated"`
}

type TagsUpdatedForScheduledRoundsPayload struct {
	UpdatedRounds []RoundUpdateInfo `json:"updated_rounds"`
	Summary       UpdateSummary     `json:"summary"`
}

// Implement the interface for ParticipantDeclinedPayload
func (p *ParticipantDeclinedPayload) GetRoundID() sharedtypes.RoundID {
	return p.RoundID
}

func (p *ParticipantDeclinedPayload) GetUserID() sharedtypes.DiscordID {
	return p.UserID
}

func (p *ParticipantDeclinedPayload) GetTagNumber() *sharedtypes.TagNumber {
	return nil // No tag number for declined participants
}

func (p *ParticipantDeclinedPayload) GetJoinedLate() *bool {
	return nil
}

// Implement the interface for RoundTagNumberNotFoundPayload
func (p *RoundTagNumberNotFoundPayload) GetRoundID() sharedtypes.RoundID {
	return p.RoundID
}

func (p *RoundTagNumberNotFoundPayload) GetUserID() sharedtypes.DiscordID {
	return p.UserID
}

func (p *RoundTagNumberNotFoundPayload) GetTagNumber() *sharedtypes.TagNumber {
	return nil // No tag number found
}

func (p *RoundTagNumberNotFoundPayload) GetJoinedLate() *bool {
	return p.OriginalJoinedLate
}

// Implement the interface for RoundTagNumberFoundPayload
func (p *RoundTagNumberFoundPayload) GetRoundID() sharedtypes.RoundID {
	return p.RoundID
}

func (p *RoundTagNumberFoundPayload) GetUserID() sharedtypes.DiscordID {
	return p.UserID
}

func (p *RoundTagNumberFoundPayload) GetTagNumber() *sharedtypes.TagNumber {
	return p.TagNumber
}

func (p *RoundTagNumberFoundPayload) GetJoinedLate() *bool {
	return p.OriginalJoinedLate
}

func (p *ParticipantJoinRequestPayload) GetRoundID() sharedtypes.RoundID {
	return p.RoundID
}

func (p *ParticipantJoinRequestPayload) GetUserID() sharedtypes.DiscordID {
	return p.UserID
}

func (p *ParticipantJoinRequestPayload) GetTagNumber() *sharedtypes.TagNumber {
	return p.TagNumber
}

func (p *ParticipantJoinRequestPayload) GetJoinedLate() *bool {
	return p.JoinedLate
}
