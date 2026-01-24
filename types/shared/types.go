package sharedtypes

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/google/uuid"
)

// DiscordID defines a custom type for Discord IDs.
type DiscordID string

type GuildID string

var userIDRegex = regexp.MustCompile(`^[0-9]+$`) // Matches one or more digits

// IsValid checks if the DiscordID is valid (contains only numbers).
func (id DiscordID) IsValid() bool {
	return userIDRegex.MatchString(string(id))
}

// RoundID defines a custom type for round identifiers.
type RoundID uuid.UUID

func (r RoundID) String() string {
	return uuid.UUID(r).String()
}

func (r RoundID) UUID() uuid.UUID {
	return uuid.UUID(r)
}

func (r *RoundID) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		id, err := uuid.ParseBytes(v)
		if err != nil {
			return fmt.Errorf("invalid UUID []byte: %w", err)
		}
		*r = RoundID(id)
		return nil
	case string:
		id, err := uuid.Parse(v)
		if err != nil {
			return fmt.Errorf("invalid UUID string: %w", err)
		}
		*r = RoundID(id)
		return nil
	default:
		return fmt.Errorf("unsupported Scan type for RoundID: %T", value)
	}
}

func (r RoundID) Value() (driver.Value, error) {
	return uuid.UUID(r).String(), nil
}

// MarshalJSON marshals the RoundID to JSON as a hyphenated UUID string.
// This ensures River queue jobs serialize round_id correctly for JSONB queries.
func (r RoundID) MarshalJSON() ([]byte, error) {
	return json.Marshal(uuid.UUID(r).String())
}

// UnmarshalJSON unmarshals the RoundID from a JSON string.
func (r *RoundID) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("invalid RoundID JSON: expected string, got %s", string(data))
	}
	id, err := uuid.Parse(str)
	if err != nil {
		return fmt.Errorf("invalid UUID string: %w", err)
	}
	*r = RoundID(id)
	return nil
}

type EventMessageID RoundID

// String returns the string representation of EventMessageID.
func (e EventMessageID) String() string {
	return uuid.UUID(e).String()
}

// MarshalJSON marshals the EventMessageID to JSON as a hyphenated UUID string.
func (e EventMessageID) MarshalJSON() ([]byte, error) {
	return json.Marshal(uuid.UUID(e).String())
}

// UnmarshalJSON unmarshals the EventMessageID from JSON.
func (e *EventMessageID) UnmarshalJSON(data []byte) error {
	var r RoundID
	if err := r.UnmarshalJSON(data); err != nil {
		return err
	}
	*e = EventMessageID(r)
	return nil
}

// Score defines a custom type for scores (can be negative or positive).
type Score int

// TagNumber defines a custom type for tag numbers.
type TagNumber int

// IsValid checks if the tag number is valid (e.g., within a certain range).
func (t TagNumber) IsValid() bool {
	return t > 0 && t <= 200
}

// String returns the string representation of the tag.
func (t TagNumber) String() string {
	return fmt.Sprintf("%d", t)
}

// UserRoleEnum represents the role of a user.
type UserRoleEnum string

// Constants for user roles
const (
	UserRoleUnknown UserRoleEnum = ""
	UserRoleUser    UserRoleEnum = "User"
	UserRoleEditor  UserRoleEnum = "Editor"
	UserRoleAdmin   UserRoleEnum = "Admin"
)

// IsValid checks if the given role is valid.
func (ur UserRoleEnum) IsValid() bool {
	switch ur {
	case UserRoleUser, UserRoleEditor, UserRoleAdmin:
		return true
	default:
		return false
	}
}

// String returns the string representation of the UserRoleEnum.
func (ur UserRoleEnum) String() string {
	return string(ur)
}

// StartTime defines a custom type for start times.
type StartTime time.Time

// MarshalJSON marshals the StartTime to JSON.
func (t StartTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(t).Format(time.RFC3339))
}

// UnmarshalJSON unmarshals the StartTime from JSON.
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

// Add delegates to time.Time.Add.
func (st StartTime) Add(d time.Duration) StartTime {
	return StartTime(time.Time(st).Add(d))
}

// UTC delegates to time.Time.UTC.
func (st StartTime) UTC() StartTime {
	return StartTime(time.Time(st).UTC())
}

// AsTime converts StartTime to time.Time.
func (st StartTime) AsTime() time.Time {
	return time.Time(st)
}

// ValidationError for payload/input validation across modules
type ValidationError struct {
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}

// Common error payload for API responses
type BaseErrorPayload struct {
	Error string `json:"error"`
}

// ScoreInfo represents a score with UserID, Score, and TagNumber.
type ScoreInfo struct {
	UserID    DiscordID  `json:"user_id"`
	Score     Score      `json:"score"`
	TagNumber *TagNumber `json:"tag_number,omitempty"`
	TeamID    uuid.UUID  `json:"team_id,omitempty"`
	RawName   string     `json:"raw_name,omitempty"` // For guest/unmatched users
}

// ScoreProcessingResult represents the result of processing scores for a round
type ScoreProcessingResult struct {
	RoundID     RoundID
	TagMappings map[DiscordID]TagNumber
}

type TagMapping struct {
	DiscordID DiscordID `json:"discord_id"`
	TagNumber TagNumber `json:"tag_number"`
}

// TagUpdateSource defines the source of tag updates
type ServiceUpdateSource string

const (
	ServiceUpdateSourceProcessScores ServiceUpdateSource = "process_scores"
	ServiceUpdateSourceCreateUser    ServiceUpdateSource = "create_user"
	ServiceUpdateSourceAdminBatch    ServiceUpdateSource = "admin_batch"
	ServiceUpdateSourceManual        ServiceUpdateSource = "manual"
	ServiceUpdateSourceTagSwap       ServiceUpdateSource = "tag_swap"
)

// TagUpdateMetadata contains data for tag update operations
type TagUpdateContext struct {
	RequestingUserID DiscordID
	BatchID          string
	RoundID          *RoundID
}

// TagAssignmentRequest represents a request to assign a tag to a user
type TagAssignmentRequest struct {
	UserID    DiscordID `json:"user_id"`
	TagNumber TagNumber `json:"tag_number"`
}

type RoundMode string

const (
	RoundModeSingles RoundMode = "SINGLES"
	RoundModeDoubles RoundMode = "DOUBLES"
	RoundModeTriples RoundMode = "TRIPLES"
	RoundModeTeams   RoundMode = "TEAMS"
	RoundModeQuads   RoundMode = "QUADS"
)
