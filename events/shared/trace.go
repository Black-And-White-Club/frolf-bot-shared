package sharedevents

// TracePayloadV1 contains generic trace event data for observability.
//
// Schema History:
//   - v1.0 (January 2026): Initial version
type TracePayloadV1 struct {
	Message string `json:"message"`
	GuildID string `json:"guild_id"`
	EventID string `json:"event_id,omitempty"`
	Module  string `json:"module,omitempty"`
}
