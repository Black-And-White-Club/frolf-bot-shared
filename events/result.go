package events

import "time"

// ResultPayload is a generic payload for success/failure events.
type ResultPayload struct {
	CommonMetadata           // Embed the common metadata
	Status         string    `json:"status"`
	ErrorDetail    string    `json:"error_detail,omitempty"`
	Timestamp      time.Time `json:"timestamp"`
}
