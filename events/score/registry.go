package scoreevents

import sharedevents "github.com/Black-And-White-Club/frolf-bot-shared/events/shared"

// GetV1Registry returns all modern events for the score functional area
func GetV1Registry() map[string]sharedevents.EventInfo {
	// NOTE: Score events are shared across modules and registered in events/shared/registry.go.
	return map[string]sharedevents.EventInfo{}
}
