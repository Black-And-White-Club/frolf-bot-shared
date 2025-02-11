package roundevents

import (
	"time"

	roundtypes "github.com/Black-And-White-Club/frolf-bot/app/modules/round/domain/types"
)

// Stream names
const (
	RoundStreamName       = "round"
	UserStreamName        = "user"
	LeaderboardStreamName = "leaderboard"
	ScoreStreamName       = "score"
)

// Round-related events
const (
	// --- Create Round ---
	RoundCreateRequest           = "round.create.request"
	RoundValidated               = "round.validated"
	RoundDateTimeParsed          = "round.datetime.parsed"
	RoundEntityCreated           = "round.entity.created"
	RoundStored                  = "round.stored"
	RoundScheduled               = "round.scheduled"
	RoundCreated                 = "round.created"
	RoundError                   = "round.error"
	RoundUpdateError             = "round.update.error"
	RoundFinalizationError       = "round.finalization.error"
	RoundDiscordEventIDUpdated   = "round.discord.event.id.updated"
	ScoreModuleNotificationError = "score.module.notification.error"
	RoundDiscordEventIDUpdate    = "round.discord.event.id.update"

	// --- Update Round ---
	RoundUpdateRequest   = "round.update.request"
	RoundUpdateValidated = "round.update.validated"
	RoundFetched         = "round.fetched"
	RoundEntityUpdated   = "round.entity.updated"
	RoundUpdated         = "round.updated"
	RoundUpdateSuccess   = "round.update.success"

	// --- Delete Round ---
	RoundDeleteRequest      = "round.delete.request"
	RoundDeleteValidated    = "round.delete.validated"
	RoundToDeleteFetched    = "round.to.delete.fetched"
	RoundDeleteAuthorized   = "round.delete.authorized"
	RoundDeleteUnauthorized = "round.delete.unauthorized"
	RoundDeleted            = "round.deleted"
	RoundDeleteError        = "round.delete.error"

	// --- Join Round ---
	RoundParticipantJoinRequest   = "round.participant.join.request"
	RoundParticipantJoinValidated = "round.participant.join.validated"
	RoundParticipantJoinError     = "round.participant.join.error"
	ParticipantJoined             = "round.participant.joined"

	// --- Score Round ---
	RoundScoreUpdateRequest      = "round.score.update.request"
	RoundScoreUpdateValidated    = "round.score.update.validated"
	RoundParticipantScoreUpdated = "round.participant.score.updated"
	RoundAllScoresSubmitted      = "round.all.scores.submitted"
	RoundNotAllScoresSubmitted   = "round.not.all.scores.submitted"
	RoundScoreUpdateError        = "round.score.update.error"

	// --- Finalize Round ---
	RoundFinalized          = "round.finalized"
	RoundScoresNotification = "round.scores.notification"

	// --- Round Reminders ---
	RoundReminder = "round.reminder"

	// --- Round State ---
	RoundStateUpdated = "round.state.updated"

	// --- Start Round ---
	RoundStarted = "round.started"

	// --- Tag Retrieval ---
	RoundTagNumberRequest  = "round.tag.number.request"
	RoundTagNumberFound    = "round.tag.number.found"
	RoundTagNumberNotFound = "round.tag.number.notfound"

	// --- Communication with Other Modules ---
	LeaderboardGetTagNumberRequest  = "leaderboard.get.tag.number.request"
	LeaderboardGetTagNumberResponse = "leaderboard.get.tag.number.response"
	ProcessRoundScoresRequest       = "score.process.scores.request"

	// --- User Authorization ---
	RoundUserRoleCheckRequest = "round.user.role.check.request"
	RoundUserRoleCheckResult  = "round.user.role.check.result"
	RoundUserRoleCheckError   = "round.user.role.check.error"

	// --- Rounds Updated ---
	RoundsUpdated       = "round.rounds.updated"
	RoundScheduleUpdate = "round.schedule.update"
	// --- Delayed Round Messages ---
	DelayedMessagesSubject = "delayed.messages"
	DiscordEventsSubject   = "discord.round.event"
)

// Round Events Payloads

// --- Create Round ---
type RoundCreateRequestPayload struct {
	Title            string                        `json:"title"`
	Location         string                        `json:"location"`
	EventType        *string                       `json:"event_type"`
	DateTime         *time.Time                    `json:"date_time"`
	Participants     []roundtypes.ParticipantInput `json:"participants"`
	DiscordChannelID string                        `json:"discord_channel_id"`
	DiscordGuildID   string                        `json:"discord_guild_id"`
	EndTime          *time.Time                    `json:"end_time,omitempty"` // Correctly a pointer
}

type RoundValidatedPayload struct {
	RoundCreateRequestPayload RoundCreateRequestPayload `json:"round_create_request_payload"`
	EndTime                   *time.Time                `json:"end_time"` // Parsed end time
}

type RoundDateTimeParsedPayload struct {
	RoundCreateRequestPayload RoundCreateRequestPayload `json:"round_create_request_payload"`
	StartTime                 *time.Time                `json:"start_time"`
	EndTime                   *time.Time                `json:"end_time"`
}

type RoundEntityCreatedPayload struct {
	Round            roundtypes.Round `json:"round"`
	DiscordChannelID string           `json:"discord_channel_id"`
	DiscordGuildID   string           `json:"discord_guild_id"`
}

type RoundStoredPayload struct {
	Round roundtypes.Round `json:"round"`
}

type RoundScheduledPayload struct {
	RoundID          string     `json:"round_id"`
	StartTime        *time.Time `json:"start_time"`
	EndTime          *time.Time `json:"end_time"`
	Title            string     `json:"title"`
	Description      string     `json:"description"`
	Location         string     `json:"location"`
	DiscordEventID   string     `json:"discord_event_id"`
	DiscordChannelID string     `json:"discord_channel_id"`
	DiscordGuildID   string     `json:"discord_guild_id"`
}

type RoundCreatedPayload struct {
	RoundID   string     `json:"round_id"`
	Name      string     `json:"name"`
	StartTime *time.Time `json:"start_time"`
}

type RoundEventCreatedPayload struct {
	RoundID        string `json:"round_id"`
	DiscordEventID string `json:"discord_event_id"`
}

type RoundDiscordEventIDUpdatedPayload struct {
	RoundID        string `json:"round_id"`
	DiscordEventID string `json:"discord_event_id"`
}

type RoundErrorPayload struct {
	CorrelationID string                     `json:"correlation_id"`
	Round         *RoundCreateRequestPayload `json:"round"`
	Error         string                     `json:"error"`
}

// --- Update Round ---
type RoundUpdateRequestPayload struct {
	RoundID        string     `json:"round_id" validate:"required"`
	DiscordEventID string     `json:"discord_event_id" validate:"required"`
	Title          *string    `json:"title,omitempty"`
	Location       *string    `json:"location,omitempty"`
	EventType      *string    `json:"event_type,omitempty"`
	StartTime      *time.Time `json:"start_time,omitempty"` //  <-- CHANGED HERE
	EndTime        *time.Time `json:"end_time,omitempty"`
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

type RoundUpdatedPayload struct {
	RoundID string `json:"round_id"`
}

type RoundUpdateSuccessPayload struct {
	RoundID string `json:"round_id"`
}

type RoundUpdateErrorPayload struct {
	CorrelationID      string                     `json:"correlation_id"`
	RoundUpdateRequest *RoundUpdateRequestPayload `json:"round_update_request"`
	Error              string                     `json:"error"`
}

// --- Delete Round ---
type RoundDeleteRequestPayload struct {
	RoundID                 string `json:"round_id" validate:"required"`
	RequestingUserDiscordID string `json:"requesting_user_discord_id" validate:"required"`
}

type RoundDeleteValidatedPayload struct {
	RoundDeleteRequestPayload RoundDeleteRequestPayload `json:"round_delete_request_payload"`
}

type RoundToDeleteFetchedPayload struct {
	Round                     roundtypes.Round          `json:"round"`
	RoundDeleteRequestPayload RoundDeleteRequestPayload `json:"round_delete_request_payload"`
}

type RoundDeleteAuthorizedPayload struct {
	RoundID string `json:"round_id"`
}

type RoundDeletedPayload struct {
	RoundID string `json:"round_id"`
}

type RoundDeleteErrorPayload struct {
	CorrelationID      string                     `json:"correlation_id"`
	RoundDeleteRequest *RoundDeleteRequestPayload `json:"round_delete_request"`
	Error              string                     `json:"error"`
}

// --- Join Round ---
type ParticipantJoinRequestPayload struct {
	RoundID     string `json:"round_id"`
	Participant string `json:"participant"` // Discord ID
}

type ParticipantJoinValidatedPayload struct {
	ParticipantJoinRequestPayload ParticipantJoinRequestPayload `json:"participant_join_request_payload"`
}

type RoundParticipantJoinErrorPayload struct {
	CorrelationID          string                         `json:"correlation_id"`
	ParticipantJoinRequest *ParticipantJoinRequestPayload `json:"participant_join_request"`
	Error                  string                         `json:"error"`
}

type ParticipantJoinedPayload struct {
	RoundID     string `json:"round_id"`
	Participant string `json:"participant"`
	TagNumber   int    `json:"tag_number,omitempty"`
	Response    string `json:"response"`
}

// --- Score Round ---
type ScoreUpdateRequestPayload struct {
	RoundID     string `json:"round_id"`
	Participant string `json:"participant"` // Discord ID
	Score       *int   `json:"score"`
}

type ScoreUpdateValidatedPayload struct {
	ScoreUpdateRequestPayload ScoreUpdateRequestPayload `json:"score_update_request_payload"`
}

type ParticipantScoreUpdatedPayload struct {
	RoundID     string `json:"round_id"`
	Participant string `json:"participant"` // Discord ID
	Score       int    `json:"score"`
}

type AllScoresSubmittedPayload struct {
	RoundID string `json:"round_id"`
}

type RoundScoreUpdateErrorPayload struct {
	CorrelationID      string                     `json:"correlation_id"`
	ScoreUpdateRequest *ScoreUpdateRequestPayload `json:"score_update_request"`
	Error              string                     `json:"error"`
}

// RoundScheduleUpdatePayload is the payload for the round.schedule.update event.
type RoundScheduleUpdatePayload struct {
	RoundID   string     `json:"round_id"`
	Title     string     `json:"title"`
	StartTime *time.Time `json:"start_time"`
	Location  string     `json:"location"`
}

// --- Finalize Round ---
type RoundFinalizedPayload struct {
	RoundID string `json:"round_id"`
}

type RoundFinalizationErrorPayload struct {
	CorrelationID string `json:"correlation_id"`
	RoundID       string `json:"round_id"`
	Error         string `json:"error"`
}

type ScoreModuleNotificationErrorPayload struct {
	CorrelationID string `json:"correlation_id"`
	RoundID       string `json:"round_id"`
	Error         string `json:"error"`
}

// --- Round Reminders ---
type RoundReminderPayload struct {
	RoundID          string     `json:"round_id"`
	ReminderType     string     `json:"reminder_type"` // e.g., "1h", "30m"
	RoundTitle       string     `json:"round_title"`
	StartTime        *time.Time `json:"start_time"`
	Location         string     `json:"location"`
	DiscordChannelID string     `json:"discord_channel_id,omitempty"`
	DiscordGuildID   string     `json:"discord_guild_id,omitempty"`
}

// --- Round Start ---
type RoundStartedPayload struct {
	RoundID   string     `json:"round_id"`
	Title     string     `json:"title"`
	Location  string     `json:"location"`
	StartTime *time.Time `json:"start_time"`
}

// --- Communication with Discord ---
type DiscordReminderPayload struct {
	RoundID          string     `json:"round_id"`
	ReminderType     string     `json:"reminder_type"`
	RoundTitle       string     `json:"round_title"`
	StartTime        *time.Time `json:"start_time"`
	Location         string     `json:"location"`
	UserIDs          []string   `json:"user_ids"`
	DiscordChannelID string     `json:"discord_channel_id,omitempty"`
	DiscordGuildID   string     `json:"discord_guild_id,omitempty"`
}

type DiscordRoundStartPayload struct {
	RoundID          string                    `json:"round_id"`
	Title            string                    `json:"title"`
	Location         string                    `json:"location"`
	StartTime        *time.Time                `json:"start_time"`
	Participants     []DiscordRoundParticipant `json:"participants"`
	DiscordChannelID string                    `json:"discord_channel_id,omitempty"`
	DiscordGuildID   string                    `json:"discord_guild_id,omitempty"`
}

type DiscordRoundParticipant struct {
	DiscordID string `json:"discord_id"`
	TagNumber int    `json:"tag_number"`
	Score     *int   `json:"score"`
}

// --- Round State ---
type RoundStateUpdatedPayload struct {
	RoundID string                `json:"round_id"`
	State   roundtypes.RoundState `json:"state"`
}

// Participant represents a participant in a round with their tag number.
type Participant struct {
	DiscordID string `json:"discord_id"`
	TagNumber int    `json:"tag_number"`
}

// --- Tag Retrieval ---
type TagNumberRequestPayload struct {
	DiscordID string        `json:"discord_id"`
	Timeout   time.Duration `json:"timeout"`
}

type RoundTagNumberFoundPayload struct {
	RoundID   string `json:"round_id"`
	DiscordID string `json:"discord_id"`
	TagNumber int    `json:"tag_number"`
}

type RoundTagNumberNotFoundPayload struct {
	DiscordID string `json:"discord_id"`
}

// --- Notify Score Module ---
type RoundScoresNotificationPayload struct {
	RoundID string             `json:"round_id"`
	Scores  []ParticipantScore `json:"scores"`
}

type ParticipantScore struct {
	DiscordID string  `json:"discord_id"`
	TagNumber string  `json:"tag_number"` // Assuming you want to keep this as a string
	Score     float64 `json:"score"`
}

// --- Process Round Scores ---
type ProcessRoundScoresRequestPayload struct {
	RoundID string             `json:"round_id"`
	Scores  []ParticipantScore `json:"scores"`
}

// --- User Authorization ---
type UserRoleCheckRequestPayload struct {
	DiscordID     string `json:"discord_id"`
	RoundID       string `json:"round_id"`       // Context for the request
	CorrelationID string `json:"correlation_id"` // To correlate with the response
}

type UserRoleCheckResultPayload struct {
	DiscordID string `json:"discord_id"`
	RoundID   string `json:"round_id"` // Context for the request
	HasRole   bool   `json:"has_role"`
	Error     string `json:"error"` // Error message if the check failed
}

type RoundUserRoleCheckErrorPayload struct {
	CorrelationID string `json:"correlation_id"`
	DiscordID     string `json:"discord_id"`
	RoundID       string `json:"round_id"`
	Error         string `json:"error"`
}

// --- Payloads for Tag Retrieval ---
type GetTagNumberResponsePayload struct {
	DiscordID string `json:"discord_id"`
	TagNumber int    `json:"tag_number"`
	Error     string `json:"error,omitempty"` // Include an error string
}
