package roundtypes

import (
	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
)

// Core Type Definitions
type (
	Title       string
	Description string
	Location    string
	EventType   string
	Finalized   bool
	CreatedBy   sharedtypes.DiscordID
	Timezone    string
)

// RoundState represents the state of a round.
type RoundState string

const (
	RoundStateUpcoming   RoundState = "UPCOMING"
	RoundStateInProgress RoundState = "IN_PROGRESS"
	RoundStateFinalized  RoundState = "FINALIZED"
	RoundStateDeleted    RoundState = "DELETED"
)

type Response string

const (
	ResponseAccept    Response = "ACCEPT"
	ResponseTentative Response = "TENTATIVE"
	ResponseDecline   Response = "DECLINE"
)

type RoundUpdate struct {
	RoundID        sharedtypes.RoundID `json:"round_id"`
	EventMessageID string              `json:"event_message_id"`
	Participants   []Participant       `json:"participants"`
	Round          *Round              `json:"round,omitempty"`
}

type Participant struct {
	UserID    sharedtypes.DiscordID  `json:"user_id"`
	TagNumber *sharedtypes.TagNumber `json:"tag_number,omitempty"`
	Response  Response               `json:"response"`
	Score     *sharedtypes.Score     `json:"score"`
}

type Round struct {
	ID             sharedtypes.RoundID    `json:"id"`
	Title          Title                  `json:"title"`
	Description    *Description           `json:"description"`
	Location       *Location              `json:"location"`
	EventType      *EventType             `json:"event_type"`
	StartTime      *sharedtypes.StartTime `json:"start_time"`
	Finalized      Finalized              `json:"finalized"`
	CreatedBy      sharedtypes.DiscordID  `json:"created_by"`
	State          RoundState             `json:"state"`
	Participants   []Participant          `json:"participants"`
	EventMessageID string                 `json:"event_message_id"`
}

const DefaultEventType = EventType("casual")

func (r *Round) AddParticipant(participant Participant) {
	r.Participants = append(r.Participants, participant)
}

type BaseRoundPayload struct {
	RoundID     sharedtypes.RoundID    `json:"round_id,omitempty"`
	Title       Title                  `json:"title,omitempty"`
	Description *Description           `json:"description,omitempty"`
	Location    *Location              `json:"location,omitempty"`
	StartTime   *sharedtypes.StartTime `json:"start_time,omitempty"`
	UserID      sharedtypes.DiscordID  `json:"user_id,omitempty"`
}

type BaseParticipantPayload struct {
	UserID    sharedtypes.DiscordID  `json:"user_id"`
	TagNumber *sharedtypes.TagNumber `json:"tag_number"`
	Response  Response               `json:"response"`
}

type BaseErrorPayload struct {
	Error string `json:"error"`
}

type CreateRoundInput struct {
	Title       Title                 `json:"title"`
	Description *Description          `json:"description,omitempty"`
	Location    *Location             `json:"location,omitempty"`
	StartTime   string                `json:"start_time"` // StartTime comes in as a string before it's processed
	UserID      sharedtypes.DiscordID `json:"user_id"`
}
