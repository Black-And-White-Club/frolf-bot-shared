package sharedevents_test

import (
	"reflect"
	"testing"

	leaderboardevents "github.com/Black-And-White-Club/frolf-bot-shared/events/leaderboard"
	sharedevents "github.com/Black-And-White-Club/frolf-bot-shared/events/shared"
)

func TestLeaderboardBatchTagAssignmentRequestedPayloadTypeIsCanonical(t *testing.T) {
	subject := sharedevents.LeaderboardBatchTagAssignmentRequestedV1

	sharedRegistry := sharedevents.GetV1Registry()
	leaderboardRegistry := leaderboardevents.GetV1Registry()

	sharedInfo, ok := sharedRegistry[subject]
	if !ok {
		t.Fatalf("shared registry missing subject %q", subject)
	}
	leaderboardInfo, ok := leaderboardRegistry[subject]
	if !ok {
		t.Fatalf("leaderboard registry missing subject %q", subject)
	}

	expectedType := reflect.TypeOf(&sharedevents.BatchTagAssignmentRequestedPayloadV1{})
	sharedType := reflect.TypeOf(sharedInfo.Payload)
	leaderboardType := reflect.TypeOf(leaderboardInfo.Payload)

	if sharedType != expectedType {
		t.Fatalf("shared registry payload type mismatch for %q: expected %v, got %v", subject, expectedType, sharedType)
	}
	if leaderboardType != expectedType {
		t.Fatalf("leaderboard registry payload type mismatch for %q: expected %v, got %v", subject, expectedType, leaderboardType)
	}
}
