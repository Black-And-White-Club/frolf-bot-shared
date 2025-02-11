package roundtypes

import "time"

// RoundState represents the state of a round.
type RoundState string

// Enum constants for RoundState
const (
	RoundStateUpcoming   RoundState = "UPCOMING"
	RoundStateInProgress RoundState = "IN_PROGRESS"
	RoundStateFinalized  RoundState = "FINALIZED"
	RoundStateDeleted    RoundState = "DELETED" // Consider removing if not needed
)

// Response represents the possible responses for a participant.
type Response string

// Define the possible response values as constants.
const (
	ResponseAccept    Response = "ACCEPT"
	ResponseTentative Response = "TENTATIVE"
	ResponseDecline   Response = "DECLINE"
)

// Round represents a single round in the tournament.
type Round struct {
	ID               string             `json:"id"`
	Title            string             `json:"title"`
	Description      *string            `json:"description"`
	Location         *string            `json:"location"`
	EventType        *string            `json:"event_type"`
	StartTime        *time.Time         `json:"start_time"`
	EndTime          *time.Time         `json:"end_time"`
	Finalized        bool               `json:"finalized"`
	CreatedBy        string             `json:"created_by"`
	State            RoundState         `json:"state"`
	Participants     []RoundParticipant `json:"participants" bun:",type:jsonb"`
	DiscordEventID   *string            `json:"discord_event_id"`
	DiscordChannelID *string            `json:"discord_channel_id"`
	DiscordGuildID   *string            `json:"discord_guild_id"`
}

// RoundParticipant represents a participant in a round, including their score.
type RoundParticipant struct {
	DiscordID string   `json:"discord_id" validate:"required"`
	TagNumber int      `json:"tag_number"` // Use a special value (-1 or 0) to represent a missing tag
	Response  Response `json:"response" validate:"required"`
	Score     *int     `json:"score,omitempty"`     // Allow score to be nil (not yet submitted)
	IsActive  bool     `json:"is_active,omitempty"` // Optional: Use this to mark participants inactive after the round starts
}

// IsUpcoming checks if the round is in the upcoming state.
func (r *Round) IsUpcoming() bool {
	return r.State == RoundStateUpcoming
}

// IsInProgress checks if the round is in the in-progress state.
func (r *Round) IsInProgress() bool {
	return r.State == RoundStateInProgress
}

// IsFinalized checks if the round is in the finalized state.
func (r *Round) IsFinalized() bool {
	return r.State == RoundStateFinalized
}

// AddParticipant adds a participant to the round.
func (r *Round) AddParticipant(participant RoundParticipant) {
	r.Participants = append(r.Participants, participant)
}

// CreateRoundInput represents the input for creating a new round.
type CreateRoundInput struct {
	Title        string             `json:"title" validate:"required"`
	Description  *string            `json:"description"`
	Location     *string            `json:"location"`
	EventType    *string            `json:"event_type"`
	StartTime    *time.Time         `json:"start_time" validate:"required"`
	EndTime      *time.Time         `json:"end_time,omitempty"`
	Participants []ParticipantInput `json:"participants"`
}

// ParticipantInput represents the input for a participant in a round.
type ParticipantInput struct {
	DiscordID string   `json:"discord_id" validate:"required"`
	TagNumber *int     `json:"tag_number"` // Keep as pointer in input if it's optional
	Response  Response `json:"response" validate:"required"`
}

// UpdateRoundInput represents the input for updating a round.
type UpdateRoundInput struct {
	RoundID     string     `json:"round_id" validate:"required"`
	Title       *string    `json:"title,omitempty"`
	Description *string    `json:"description,omitempty"`
	Location    *string    `json:"location,omitempty"`
	EventType   *string    `json:"event_type,omitempty"`
	StartTime   *time.Time `json:"start_time,omitempty"`
	EndTime     *time.Time `json:"end_time,omitempty"`
}

// UpdateParticipantResponseInput represents the input for updating a participant's response.
type UpdateParticipantResponseInput struct {
	RoundID     string   `json:"round_id" validate:"required"`
	Participant string   `json:"participant" validate:"required"` // Discord ID
	Response    Response `json:"response" validate:"required"`
	TagNumber   *int     `json:"tag_number"` // Allow updating the tag number
}
