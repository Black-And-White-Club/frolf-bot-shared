package roundtypes

import (
	"encoding/json"
	"time"

	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
)

// Core Type Definitions
type Title string
type Description string
type Location string
type EventType string
type StartTime time.Time
type Finalized bool
type CreatedBy sharedtypes.DiscordID
type Timezone string
type EventMessageID string

func (t StartTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(t).Format(time.RFC3339))
}

func (t *StartTime) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	parsedTime, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return err
	}

	*t = StartTime(parsedTime)
	return nil
}

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

type Participant struct {
	UserID    sharedtypes.DiscordID  `json:"user_id"`
	TagNumber *sharedtypes.TagNumber `json:"tag_number,omitempty"`
	Response  Response               `json:"response"`
	Score     *sharedtypes.Score     `json:"score"`
}

type Round struct {
	ID             sharedtypes.RoundID `json:"id"`
	Title          Title               `json:"title"`
	Description    *Description        `json:"description"`
	Location       *Location           `json:"location"`
	EventType      *EventType          `json:"event_type"`
	StartTime      *StartTime          `json:"start_time"`
	Finalized      Finalized           `json:"finalized"`
	CreatedBy      CreatedBy           `json:"created_by"`
	State          RoundState          `json:"state"`
	Participants   []Participant       `json:"participants"`
	EventMessageID EventMessageID      `json:"event_message_id"`
}

func (r *Round) AddParticipant(participant Participant) {
	r.Participants = append(r.Participants, participant)
}

type BaseRoundPayload struct {
	RoundID     sharedtypes.RoundID   `json:"round_id,omitempty"`
	Title       Title                 `json:"title,omitempty"`
	Description *Description          `json:"description,omitempty"`
	Location    *Location             `json:"location,omitempty"`
	StartTime   *StartTime            `json:"start_time,omitempty"`
	UserID      sharedtypes.DiscordID `json:"user_id,omitempty"`
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
	StartTime   StartTime             `json:"start_time"`
	UserID      sharedtypes.DiscordID `json:"user_id"`
}
