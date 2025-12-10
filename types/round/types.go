package roundtypes

import (
	"time"

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
	GuildID        sharedtypes.GuildID    `json:"guild_id"`
	// Import/scorecard fields
	ImportID        string     `json:"import_id,omitempty"`
	ImportStatus    string     `json:"import_status,omitempty"`
	ImportType      string     `json:"import_type,omitempty"`
	FileData        []byte     `json:"file_data,omitempty"`
	FileName        string     `json:"file_name,omitempty"`
	UDiscURL        string     `json:"udisc_url,omitempty"`
	ImportNotes     string     `json:"import_notes,omitempty"`
	ImportError     string     `json:"import_error,omitempty"`
	ImportErrorCode string     `json:"import_error_code,omitempty"`
	ImportedAt      *time.Time `json:"imported_at,omitempty"`
	// Import context - who initiated and where to respond
	ImportUserID    sharedtypes.DiscordID `json:"import_user_id,omitempty"`
	ImportChannelID string                `json:"import_channel_id,omitempty"`
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

// ParsedScorecard represents the result of parsing a scorecard file
type ParsedScorecard struct {
	ImportID     string              `json:"import_id"`
	RoundID      sharedtypes.RoundID `json:"round_id"`
	GuildID      sharedtypes.GuildID `json:"guild_id"`
	ParScores    []int               `json:"par_scores"`
	PlayerScores []PlayerScoreRow    `json:"player_scores"`
	StartTime    *time.Time          `json:"start_time,omitempty"`
	EndTime      *time.Time          `json:"end_time,omitempty"`
}

// PlayerScoreRow represents a single player's scores from a parsed scorecard
type PlayerScoreRow struct {
	PlayerName string `json:"player_name"`
	HoleScores []int  `json:"hole_scores"`
	Total      int    `json:"total"`
}
