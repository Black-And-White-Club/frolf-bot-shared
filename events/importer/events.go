package importerevents

import (
	roundtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/round"
	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
)

// Stream name
const (
	ImporterStreamName = "importer"
)

// Event names
const (
	ScorecardUploaded    = "udisc.scorecard.uploaded"
	ScorecardParsed      = "udisc.scorecard.parsed"
	ImportCompleted      = "udisc.import.completed"
	ImportFailed         = "udisc.import.failed"
	ParticipantAutoAdded = "round.participant.auto_added"
	ScoresImported       = "round.scores.imported"
)

// ScorecardUploadedPayload is published when a user uploads a scorecard file.
type ScorecardUploadedPayload struct {
	ImportID   string                `json:"import_id"`
	RoundID    sharedtypes.RoundID   `json:"round_id"`
	GuildID    sharedtypes.GuildID   `json:"guild_id"`
	UploaderID sharedtypes.DiscordID `json:"uploader_id"`
	ChannelID  string                `json:"channel_id"`
	MessageID  string                `json:"message_id"`
	FileData   []byte                `json:"file_data"`
	FileName   string                `json:"file_name"`
	FileType   string                `json:"file_type"` // csv or xlsx
	Timestamp  int64                 `json:"timestamp"`
}

// ScorecardParsedPayload is published when scorecard parsing succeeds.
type ScorecardParsedPayload struct {
	ImportID   string                      `json:"import_id"`
	RoundID    sharedtypes.RoundID         `json:"round_id"`
	GuildID    sharedtypes.GuildID         `json:"guild_id"`
	UploaderID sharedtypes.DiscordID       `json:"uploader_id"`
	ParsedData *roundtypes.ParsedScorecard `json:"parsed_data"`
	Timestamp  int64                       `json:"timestamp"`
}

// ImportCompletedPayload is published when import completes successfully.
type ImportCompletedPayload struct {
	ImportID       string                `json:"import_id"`
	RoundID        sharedtypes.RoundID   `json:"round_id"`
	GuildID        sharedtypes.GuildID   `json:"guild_id"`
	UploaderID     sharedtypes.DiscordID `json:"uploader_id"`
	MatchedPlayers []MatchedPlayer       `json:"matched_players"`
	SkippedPlayers []string              `json:"skipped_players"`
	ScoresImported int                   `json:"scores_imported"`
	Timestamp      int64                 `json:"timestamp"`
}

// MatchedPlayer represents a player successfully matched and imported.
type MatchedPlayer struct {
	DiscordID sharedtypes.DiscordID `json:"discord_id"`
	UDiscName string                `json:"udisc_name"`
	Score     int                   `json:"score"`
}

// ImportFailedPayload is published when import fails.
type ImportFailedPayload struct {
	ImportID   string                `json:"import_id"`
	RoundID    sharedtypes.RoundID   `json:"round_id"`
	GuildID    sharedtypes.GuildID   `json:"guild_id"`
	UploaderID sharedtypes.DiscordID `json:"uploader_id"`
	Error      string                `json:"error"`
	Timestamp  int64                 `json:"timestamp"`
}

// ParticipantAutoAddedPayload is published when a player is auto-added to a round.
type ParticipantAutoAddedPayload struct {
	RoundID   sharedtypes.RoundID   `json:"round_id"`
	GuildID   sharedtypes.GuildID   `json:"guild_id"`
	UserID    sharedtypes.DiscordID `json:"user_id"`
	ImportID  string                `json:"import_id"`
	Reason    string                `json:"reason"` // e.g., "imported_from_scorecard"
	Timestamp int64                 `json:"timestamp"`
}

// ScoresImportedPayload is published when scores are inserted into the round.
type ScoresImportedPayload struct {
	RoundID     sharedtypes.RoundID `json:"round_id"`
	GuildID     sharedtypes.GuildID `json:"guild_id"`
	ImportID    string              `json:"import_id"`
	ScoresCount int                 `json:"scores_count"`
	Timestamp   int64               `json:"timestamp"`
}
