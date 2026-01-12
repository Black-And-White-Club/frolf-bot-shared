package roundevents

// This file adds the GetEventMessageID method implementations for V1 payloads
// that carry an EventMessageID so the handler wrapper can promote the value
// into outgoing message metadata automatically.

// GetEventMessageID returns the embedded EventMessageID for scheduled rounds.
func (p RoundScheduledPayloadV1) GetEventMessageID() string { return p.EventMessageID }

// GetEventMessageID returns the embedded EventMessageID for discord round start payloads.
func (p DiscordRoundStartPayloadV1) GetEventMessageID() string { return p.EventMessageID }

// GetEventMessageID returns the embedded EventMessageID for finalized discord payloads.
func (p RoundFinalizedDiscordPayloadV1) GetEventMessageID() string { return p.EventMessageID }

// GetEventMessageID returns the embedded EventMessageID for finalized embed update payloads.
func (p RoundFinalizedEmbedUpdatePayloadV1) GetEventMessageID() string { return p.EventMessageID }

// GetEventMessageID returns the embedded EventMessageID for reminder payloads.
func (p DiscordReminderPayloadV1) GetEventMessageID() string { return p.EventMessageID }

// GetEventMessageID returns the embedded EventMessageID for event created payloads.
func (p RoundEventCreatedPayloadV1) GetEventMessageID() string { return p.EventMessageID }

// GetEventMessageID returns the embedded EventMessageID for round update info (tags updates).
func (p RoundUpdateInfoV1) GetEventMessageID() string { return p.EventMessageID }
