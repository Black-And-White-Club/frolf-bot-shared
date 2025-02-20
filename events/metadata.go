package events

import "github.com/ThreeDotsLabs/watermill/message"

// MetadataCarrier is an interface for payloads that provide their own metadata.
type MetadataCarrier interface {
	GetEventName() string
	GetDomain() string
}

// CommonMetadata holds metadata fields common to many events.
type CommonMetadata struct {
	EventName string `json:"event_name"`
	Domain    string `json:"domain"`
}

// GetEventName implements the MetadataCarrier interface.
func (c CommonMetadata) GetEventName() string {
	return c.EventName
}

// GetDomain implements the MetadataCarrier interface.
func (c CommonMetadata) GetDomain() string {
	return c.Domain
}

// WithMetadata sets common metadata on a Watermill message.
func WithMetadata(msg *message.Message, payload interface{}) {

	// Copy existing metadata (especially correlation ID)
	for k, v := range msg.Metadata {
		if _, ok := msg.Metadata[k]; !ok { // Prevent overwriting.
			msg.Metadata.Set(k, v)
		}
	}

	// Handle the common metadata, checking if there is an override
	if p, ok := payload.(MetadataCarrier); ok {
		msg.Metadata.Set("event", p.GetEventName())
		msg.Metadata.Set("domain", p.GetDomain())
	} else if p, ok := payload.(CommonMetadata); ok {
		msg.Metadata.Set("event", p.GetEventName())
		msg.Metadata.Set("domain", p.GetDomain())
	} else if p, ok := payload.(*CommonMetadata); ok { //Pointer receiver, just in case
		msg.Metadata.Set("event", p.GetEventName())
		msg.Metadata.Set("domain", p.GetDomain())
	}
}
