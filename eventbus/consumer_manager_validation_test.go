package eventbus

import (
	"strings"
	"testing"
	"time"

	"github.com/nats-io/nats.go/jetstream"
)

func TestValidateConsumerInfo_SucceedsForExpectedConfig(t *testing.T) {
	cfg := ConsumerConfig{
		AckWait:           45 * time.Second,
		MaxDeliver:        7,
		BackOff:           []time.Duration{time.Second, 2 * time.Second},
		MaxAckPending:     200,
		DeliverPolicy:     jetstream.DeliverAllPolicy,
		InactiveThreshold: 0,
	}
	expected := buildDurableConsumerConfig("backend-round-finalized-v1", "round.finalized.v1", cfg)
	info := &jetstream.ConsumerInfo{Config: expected}

	if err := validateConsumerInfo(info, expected); err != nil {
		t.Fatalf("expected validation success, got error: %v", err)
	}
}

func TestValidateConsumerInfo_FailsOnDeliverPolicyMismatch(t *testing.T) {
	cfg := ConsumerConfig{
		AckWait:           45 * time.Second,
		MaxDeliver:        7,
		BackOff:           []time.Duration{time.Second, 2 * time.Second},
		MaxAckPending:     200,
		DeliverPolicy:     jetstream.DeliverAllPolicy,
		InactiveThreshold: 0,
	}
	expected := buildDurableConsumerConfig("backend-round-finalized-v1", "round.finalized.v1", cfg)
	actual := expected
	actual.DeliverPolicy = jetstream.DeliverNewPolicy

	info := &jetstream.ConsumerInfo{Config: actual}
	err := validateConsumerInfo(info, expected)
	if err == nil {
		t.Fatal("expected validation error for deliver policy mismatch")
	}
	if !strings.Contains(err.Error(), "deliver policy mismatch") {
		t.Fatalf("expected deliver policy mismatch error, got: %v", err)
	}
}
