package events

import (
	"reflect"
	"time"
)

// MetadataCarrier is an interface for payloads that provide their own metadata.
type MetadataCarrier interface {
	GetEventName() string
	GetDomain() string
}

// CommonMetadata defines common fields for all events.
type CommonMetadata struct {
	Domain    string    `json:"domain"`
	EventName string    `json:"event_name"`
	Timestamp time.Time `json:"timestamp"`
}

// GetEventName implements the MetadataCarrier interface.
func (c CommonMetadata) GetEventName() string {
	return c.EventName
}

// GetDomain implements the MetadataCarrier interface.
func (c CommonMetadata) GetDomain() string {
	return c.Domain
}

// getEventName attempts to extract the EventName from a struct.
func getEventName(payload interface{}) (string, bool) {
	val := reflect.ValueOf(payload)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return "", false
	}
	field := val.FieldByName("EventName")
	if !field.IsValid() {
		return "", false
	}
	if field.Kind() != reflect.String {
		return "", false
	}

	return field.String(), true
}
