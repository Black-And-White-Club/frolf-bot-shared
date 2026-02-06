package roundtypes

import (
	"fmt"
	"time"

	guildtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/guild"
	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
	"github.com/google/uuid"
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

func (t Title) String() string       { return string(t) }
func (d Description) String() string { return string(d) }
func (l Location) String() string    { return string(l) }
func (e EventType) String() string   { return string(e) }

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
	RoundID                  sharedtypes.RoundID `json:"round_id"`
	EventMessageID           string              `json:"event_message_id"`
	Participants             []Participant       `json:"participants"`
	ParticipantsChangedCount int                 `json:"participants_changed_count"`
	Round                    *Round              `json:"round,omitempty"`
}

type Participant struct {
	UserID    sharedtypes.DiscordID  `json:"user_id"`
	TagNumber *sharedtypes.TagNumber `json:"tag_number,omitempty"`
	Response  Response               `json:"response"`
	Score     *sharedtypes.Score     `json:"score"`
	TeamID    uuid.UUID              `json:"team_id,omitempty"`
	RawName   string                 `json:"raw_name,omitempty"` // For guest/unmatched users
}

type Round struct {
	ID             sharedtypes.RoundID    `json:"id"`
	Title          Title                  `json:"title"`
	Description    Description            `json:"description"`
	Location       Location               `json:"location"`
	EventType      *EventType             `json:"event_type"`
	StartTime      *sharedtypes.StartTime `json:"start_time"`
	Finalized      Finalized              `json:"finalized"`
	CreatedBy      sharedtypes.DiscordID  `json:"created_by"`
	State          RoundState             `json:"state"`
	Participants   []Participant          `json:"participants"`
	EventMessageID string                 `json:"event_message_id"`
	DiscordEventID string                 `json:"discord_event_id,omitempty"`
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
	Teams           []NormalizedTeam      `json:"teams,omitempty"`
}

const DefaultEventType = EventType("casual")

func (r *Round) AddParticipant(participant Participant) {
	r.Participants = append(r.Participants, participant)
}

type BaseRoundPayload struct {
	RoundID     sharedtypes.RoundID    `json:"round_id,omitempty"`
	Title       Title                  `json:"title,omitempty"`
	Description Description            `json:"description,omitempty"`
	Location    Location               `json:"location,omitempty"`
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
	Description Description           `json:"description,omitempty"`
	Location    Location              `json:"location,omitempty"`
	StartTime   string                `json:"start_time"` // StartTime comes in as a string before it's processed
	UserID      sharedtypes.DiscordID `json:"user_id"`
	GuildID     sharedtypes.GuildID   `json:"guild_id"`
	Timezone    string                `json:"timezone"`
	ChannelID   string                `json:"channel_id"`
}

type CreateRoundResult struct {
	Round       *Round                  `json:"round"`
	ChannelID   string                  `json:"channel_id"`
	GuildConfig *guildtypes.GuildConfig `json:"guild_config,omitempty"`
}

type UpdateRoundResult struct {
	Round       *Round                  `json:"round"`
	GuildConfig *guildtypes.GuildConfig `json:"guild_config,omitempty"`
}

type UpdateRoundRequest struct {
	GuildID         sharedtypes.GuildID    `json:"guild_id"`
	RoundID         sharedtypes.RoundID    `json:"round_id"`
	Title           *string                `json:"title,omitempty"`
	Description     *string                `json:"description,omitempty"`
	Location        *string                `json:"location,omitempty"`
	StartTime       *string                `json:"start_time,omitempty"`
	Timezone        *string                `json:"timezone,omitempty"`
	UserID          sharedtypes.DiscordID  `json:"user_id"`
	ParsedStartTime *sharedtypes.StartTime `json:"parsed_start_time,omitempty"`
	EventType       *EventType             `json:"event_type,omitempty"`
}

type UpdateScheduledRoundEventsRequest struct {
	GuildID   sharedtypes.GuildID    `json:"guild_id"`
	RoundID   sharedtypes.RoundID    `json:"round_id"`
	Title     *string                `json:"title,omitempty"`
	StartTime *sharedtypes.StartTime `json:"start_time,omitempty"`
	Location  *string                `json:"location,omitempty"`
}

type JoinRoundRequest struct {
	GuildID    sharedtypes.GuildID    `json:"guild_id"`
	RoundID    sharedtypes.RoundID    `json:"round_id"`
	UserID     sharedtypes.DiscordID  `json:"user_id"`
	Response   Response               `json:"response"`
	TagNumber  *sharedtypes.TagNumber `json:"tag_number,omitempty"`
	JoinedLate *bool                  `json:"joined_late,omitempty"`
}

type ParticipantStatusCheckResult struct {
	Action        string                `json:"action"` // "VALIDATE", "REMOVE"
	CurrentStatus string                `json:"current_status"`
	RoundID       sharedtypes.RoundID   `json:"round_id"`
	UserID        sharedtypes.DiscordID `json:"user_id"`
	Response      Response              `json:"response"`
	GuildID       sharedtypes.GuildID   `json:"guild_id"`
}

type DeleteRoundInput struct {
	GuildID sharedtypes.GuildID   `json:"guild_id"`
	RoundID sharedtypes.RoundID   `json:"round_id"`
	UserID  sharedtypes.DiscordID `json:"user_id"`
}

type StartRoundRequest struct {
	GuildID sharedtypes.GuildID `json:"guild_id"`
	RoundID sharedtypes.RoundID `json:"round_id"`
}

type FinalizeRoundInput struct {
	GuildID sharedtypes.GuildID `json:"guild_id"`
	RoundID sharedtypes.RoundID `json:"round_id"`
}

type FinalizeRoundResult struct {
	Round        *Round           `json:"round"`
	Participants []Participant    `json:"participants"`
	Teams        []NormalizedTeam `json:"teams"`
}

// ParsedScorecard represents the result of parsing a scorecard file
type ParsedScorecard struct {
	ImportID     string                `json:"import_id"`
	RoundID      sharedtypes.RoundID   `json:"round_id"`
	GuildID      sharedtypes.GuildID   `json:"guild_id"`
	ParScores    []int                 `json:"par_scores"`
	PlayerScores []PlayerScoreRow      `json:"player_scores"`
	StartTime    *time.Time            `json:"start_time,omitempty"`
	EndTime      *time.Time            `json:"end_time,omitempty"`
	Mode         sharedtypes.RoundMode `json:"mode"`
}

type TeamMember struct {
	UserID  *sharedtypes.DiscordID `json:"user_id,omitempty"` // Nil if guest/unmatched
	RawName string                 `json:"raw_name"`          // Always present (from UDisc)
}

// PlayerScoreRow represents a single player's (or team's) scores from a parsed scorecard
type PlayerScoreRow struct {
	PlayerName string   `json:"player_name"`
	HoleScores []int    `json:"hole_scores"`
	Total      int      `json:"total"`
	IsTeam     bool     `json:"is_team,omitempty"`
	TeamNames  []string `json:"team_names,omitempty"`
}

// MatchedPlayer represents a player successfully matched and imported.
type MatchedPlayer struct {
	DiscordID sharedtypes.DiscordID `json:"discord_id"`
	UDiscName string                `json:"udisc_name"`
	Score     int                   `json:"score"`
}

type NormalizedScorecard struct {
	ID        string                `json:"id"`
	RoundID   sharedtypes.RoundID   `json:"round_id"`
	GuildID   sharedtypes.GuildID   `json:"guild_id"`
	ImportID  string                `json:"import_id"`
	Mode      sharedtypes.RoundMode `json:"mode"`              // SINGLES or DOUBLES
	Players   []NormalizedPlayer    `json:"players,omitempty"` // For singles
	Teams     []NormalizedTeam      `json:"teams,omitempty"`   // For doubles
	ParScores []int                 `json:"par_scores,omitempty"`
	CreatedAt time.Time             `json:"created_at"`
}

type NormalizedTeam struct {
	TeamID     uuid.UUID    `json:"team_id"`
	Members    []TeamMember `json:"members"`
	Total      int          `json:"total"`
	HoleScores []int        `json:"hole_scores,omitempty"`
}

type NormalizedPlayer struct {
	DisplayName string `json:"display_name"`
	Total       int    `json:"total"`
	HoleScores  []int  `json:"hole_scores"`
}

// Displayable represents any type that can provide a Discord ID and optional raw name.
type Displayable interface {
	UserIDPointer() *sharedtypes.DiscordID
	RawNameString() string
}

// DisplayName returns the human-readable display name for any Displayable.
func DisplayName(userID *sharedtypes.DiscordID, rawName string) string {
	if userID != nil {
		return fmt.Sprintf("<@%s>", *userID)
	}
	if rawName != "" {
		return rawName
	}
	return "Unknown Player"
}

// Participant implements Displayable
func (p Participant) UserIDPointer() *sharedtypes.DiscordID {
	if p.UserID == "" {
		return nil
	}
	return &p.UserID
}
func (p Participant) RawNameString() string { return p.RawName }

// TeamMember implements Displayable
func (m TeamMember) UserIDPointer() *sharedtypes.DiscordID { return m.UserID }
func (m TeamMember) RawNameString() string                 { return m.RawName }

type NormalizedSinglesEntry struct {
	PlayerName string
	Score      int
}

type NormalizedTeamEntry struct {
	TeamID  uuid.UUID
	Members []string // raw names only
	Score   int
}

type Metadata struct {
	ImportID       string
	GuildID        sharedtypes.GuildID
	RoundID        sharedtypes.RoundID
	UserID         sharedtypes.DiscordID
	ChannelID      string
	EventMessageID string
}

// ReminderData contains the data needed to schedule a round reminder.
// Moved to shared types to avoid circular dependencies between Service and Queue.
type ReminderData struct {
	GuildID          sharedtypes.GuildID
	RoundID          sharedtypes.RoundID
	ReminderType     string
	RoundTitle       string
	Location         string
	StartTime        string
	EventMessageID   string
	DiscordChannelID string
}

type ProcessRoundReminderRequest struct {
	GuildID          sharedtypes.GuildID `json:"guild_id"`
	RoundID          sharedtypes.RoundID `json:"round_id"`
	ReminderType     string              `json:"reminder_type"`
	RoundTitle       string              `json:"round_title"`
	Location         string              `json:"location"`
	StartTime        string              `json:"start_time"`
	EventMessageID   string              `json:"event_message_id"`
	DiscordChannelID string              `json:"discord_channel_id"`
	DiscordGuildID   string              `json:"discord_guild_id"`
}

type ProcessRoundReminderResult struct {
	GuildID          sharedtypes.GuildID     `json:"guild_id"`
	RoundID          sharedtypes.RoundID     `json:"round_id"`
	RoundTitle       string                  `json:"round_title"`
	StartTime        string                  `json:"start_time"`
	Location         string                  `json:"location"`
	UserIDs          []sharedtypes.DiscordID `json:"user_ids"`
	ReminderType     string                  `json:"reminder_type"`
	EventMessageID   string                  `json:"event_message_id"`
	DiscordChannelID string                  `json:"discord_channel_id"`
	DiscordGuildID   string                  `json:"discord_guild_id"`
}

type ScheduleRoundEventsRequest struct {
	GuildID        sharedtypes.GuildID     `json:"guild_id"`
	RoundID        sharedtypes.RoundID     `json:"round_id"`
	Title          string                  `json:"title"`
	Description    string                  `json:"description"`
	Location       string                  `json:"location"`
	StartTime      sharedtypes.StartTime   `json:"start_time"`
	UserID         sharedtypes.DiscordID   `json:"user_id"`
	EventMessageID string                  `json:"event_message_id"`
	ChannelID      string                  `json:"channel_id,omitempty"`
	Config         *guildtypes.GuildConfig `json:"config,omitempty"`
}

type ScheduleRoundEventsResult struct {
	RoundID        sharedtypes.RoundID   `json:"round_id"`
	GuildID        sharedtypes.GuildID   `json:"guild_id"`
	Title          string                `json:"title"`
	Description    string                `json:"description"`
	Location       string                `json:"location"`
	StartTime      sharedtypes.StartTime `json:"start_time"`
	EventMessageID string                `json:"event_message_id"`
	ChannelID      string                `json:"channel_id"`
}

type ScoreUpdateRequest struct {
	GuildID   sharedtypes.GuildID    `json:"guild_id"`
	RoundID   sharedtypes.RoundID    `json:"round_id"`
	UserID    sharedtypes.DiscordID  `json:"user_id"`
	Score     *sharedtypes.Score     `json:"score"`
	TagNumber *sharedtypes.TagNumber `json:"tag_number,omitempty"`
}

type ScoreUpdateResult struct {
	RoundID             sharedtypes.RoundID `json:"round_id"`
	GuildID             sharedtypes.GuildID `json:"guild_id"`
	EventMessageID      string              `json:"event_message_id"`
	UpdatedParticipants []Participant       `json:"updated_participants"`
}

type BulkScoreUpdateRequest struct {
	GuildID sharedtypes.GuildID  `json:"guild_id"`
	RoundID sharedtypes.RoundID  `json:"round_id"`
	Updates []ScoreUpdateRequest `json:"updates"`
}

type BulkScoreUpdateResult struct {
	GuildID sharedtypes.GuildID  `json:"guild_id"`
	RoundID sharedtypes.RoundID  `json:"round_id"`
	Updates []ScoreUpdateRequest `json:"updates"`
}

type CheckAllScoresSubmittedRequest struct {
	GuildID sharedtypes.GuildID   `json:"guild_id"`
	RoundID sharedtypes.RoundID   `json:"round_id"`
	UserID  sharedtypes.DiscordID `json:"user_id"`
}

type AllScoresSubmittedResult struct {
	IsComplete   bool             `json:"is_complete"`
	Round        *Round           `json:"round,omitempty"`
	Participants []Participant    `json:"participants,omitempty"`
	Teams        []NormalizedTeam `json:"teams,omitempty"`
}

type ImportScoreData struct {
	UserID  sharedtypes.DiscordID `json:"user_id"`
	RawName string                `json:"raw_name"`
	Score   int                   `json:"score"`
	TeamID  uuid.UUID             `json:"team_id,omitempty"`
}

type ImportApplyScoresInput struct {
	GuildID  sharedtypes.GuildID `json:"guild_id"`
	RoundID  sharedtypes.RoundID `json:"round_id"`
	ImportID string              `json:"import_id"`
	Scores   []ImportScoreData   `json:"scores"`
}

type ImportApplyScoresResult struct {
	GuildID        sharedtypes.GuildID `json:"guild_id"`
	RoundID        sharedtypes.RoundID `json:"round_id"`
	ImportID       string              `json:"import_id"`
	Participants   []Participant       `json:"participants"`
	Teams          []NormalizedTeam    `json:"teams,omitempty"`
	RoundData      *Round              `json:"round_data,omitempty"`
	EventMessageID string              `json:"event_message_id"`
	Timestamp      time.Time           `json:"timestamp"`
}

type ImportIngestScorecardInput struct {
	ImportID       string                `json:"import_id"`
	GuildID        sharedtypes.GuildID   `json:"guild_id"`
	RoundID        sharedtypes.RoundID   `json:"round_id"`
	UserID         sharedtypes.DiscordID `json:"user_id"`
	ChannelID      string                `json:"channel_id"`
	EventMessageID string                `json:"event_message_id"`
	NormalizedData NormalizedScorecard   `json:"normalized_data"`
}

type IngestScorecardResult struct {
	ImportID         string                  `json:"import_id"`
	GuildID          sharedtypes.GuildID     `json:"guild_id"`
	RoundID          sharedtypes.RoundID     `json:"round_id"`
	UserID           sharedtypes.DiscordID   `json:"user_id"`
	ChannelID        string                  `json:"channel_id"`
	ScoresIngested   int                     `json:"scores_ingested"`
	MatchedPlayers   int                     `json:"matched_players"`
	UnmatchedPlayers int                     `json:"unmatched_players"`
	SkippedPlayers   []string                `json:"skipped_players"`
	Scores           []sharedtypes.ScoreInfo `json:"scores"`
	RoundMode        sharedtypes.RoundMode   `json:"round_mode"`
	EventMessageID   string                  `json:"event_message_id"`
	Timestamp        time.Time               `json:"timestamp"`
}

type ImportCreateJobInput struct {
	ImportID  string                `json:"import_id"`
	GuildID   sharedtypes.GuildID   `json:"guild_id"`
	RoundID   sharedtypes.RoundID   `json:"round_id"`
	UserID    sharedtypes.DiscordID `json:"user_id"`
	ChannelID string                `json:"channel_id"`
	FileName  string                `json:"file_name"`
	FileData  []byte                `json:"file_data"`
	UDiscURL  string                `json:"udisc_url"`
	Notes     string                `json:"notes"`
}

type ImportParseScorecardInput struct {
	ImportID string              `json:"import_id"`
	GuildID  sharedtypes.GuildID `json:"guild_id"`
	RoundID  sharedtypes.RoundID `json:"round_id"`
	JobID    string              `json:"job_id"`
	FileName string              `json:"file_name"`
	FileData []byte              `json:"file_data"`
	FileURL  string              `json:"file_url"`
}

type CreateImportJobResult struct {
	Job *ImportCreateJobInput `json:"job"`
}

type UpdateScheduledRoundsWithNewTagsRequest struct {
	GuildID     sharedtypes.GuildID                             `json:"guild_id"`
	ChangedTags map[sharedtypes.DiscordID]sharedtypes.TagNumber `json:"changed_tags"`
}

type ScheduledRoundsSyncResult struct {
	GuildID      sharedtypes.GuildID `json:"guild_id"`
	Updates      []RoundUpdate       `json:"updates"`
	TotalChecked int                 `json:"total_checked"`
}
