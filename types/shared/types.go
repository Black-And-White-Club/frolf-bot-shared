package sharedtypes

import (
	"encoding/json"
	"fmt"
	"regexp"
	"time"
)

// DiscordID defines a custom type for Discord IDs.
type DiscordID string

var userIDRegex = regexp.MustCompile(`^[0-9]+$`) // Matches one or more digits

// IsValid checks if the DiscordID is valid (contains only numbers).
func (id DiscordID) IsValid() bool {
	return userIDRegex.MatchString(string(id))
}

// RoundID defines a custom type for round identifiers.
type RoundID int64

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
	UserRoleRattler UserRoleEnum = "Rattler"
	UserRoleEditor  UserRoleEnum = "Editor"
	UserRoleAdmin   UserRoleEnum = "Admin"
)

// IsValid checks if the given role is valid.
func (ur UserRoleEnum) IsValid() bool {
	switch ur {
	case UserRoleRattler, UserRoleEditor, UserRoleAdmin:
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

// Common error payload for API responses
type BaseErrorPayload struct {
	Error string `json:"error"`
}
