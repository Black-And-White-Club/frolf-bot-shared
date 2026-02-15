package eventbus

import (
	"testing"
	"time"
)

func TestAckHeartbeatAndDeadline_UsesConfiguredAckWait(t *testing.T) {
	heartbeat, deadline := ackHeartbeatAndDeadline(30 * time.Second)

	if heartbeat != 10*time.Second {
		t.Fatalf("expected heartbeat interval 10s, got %s", heartbeat)
	}
	if deadline != 90*time.Second {
		t.Fatalf("expected max processing duration 90s, got %s", deadline)
	}
}

func TestAckHeartbeatAndDeadline_EnforcesMinimumHeartbeat(t *testing.T) {
	heartbeat, deadline := ackHeartbeatAndDeadline(2 * time.Second)

	if heartbeat != time.Second {
		t.Fatalf("expected heartbeat interval floor of 1s, got %s", heartbeat)
	}
	if deadline != 6*time.Second {
		t.Fatalf("expected max processing duration 6s, got %s", deadline)
	}
}

func TestAckHeartbeatAndDeadline_FallsBackToDefaultAckWait(t *testing.T) {
	defaultAckWait := DefaultConsumerConfig().AckWait

	heartbeat, deadline := ackHeartbeatAndDeadline(0)

	if heartbeat != defaultAckWait/3 {
		t.Fatalf("expected heartbeat interval %s, got %s", defaultAckWait/3, heartbeat)
	}
	if deadline != defaultAckWait*maxAckWaitExtensions {
		t.Fatalf("expected max processing duration %s, got %s", defaultAckWait*maxAckWaitExtensions, deadline)
	}
}

func TestFetchErrorBackoff(t *testing.T) {
	tests := []struct {
		name     string
		attempts int
		want     time.Duration
	}{
		{name: "normal first attempt", attempts: 1, want: 100 * time.Millisecond},
		{name: "grows exponentially", attempts: 4, want: 800 * time.Millisecond},
		{name: "clamps at max", attempts: 7, want: 5 * time.Second},
		{name: "negative attempts fallback", attempts: -3, want: 100 * time.Millisecond},
		{name: "very large attempts avoid overflow", attempts: 1_000_000, want: 5 * time.Second},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := fetchErrorBackoff(tc.attempts)
			if got != tc.want {
				t.Fatalf("attempts=%d expected %s, got %s", tc.attempts, tc.want, got)
			}
		})
	}
}
