package eventbus

import (
	"testing"

	"github.com/nats-io/nats.go/jetstream"
)

func TestDefaultConsumerConfig_UsesDurableSafeDefaults(t *testing.T) {
	cfg := DefaultConsumerConfig()

	if cfg.DeliverPolicy != jetstream.DeliverAllPolicy {
		t.Fatalf("expected DeliverAllPolicy, got %v", cfg.DeliverPolicy)
	}
	if cfg.InactiveThreshold != 0 {
		t.Fatalf("expected InactiveThreshold=0 for durables, got %s", cfg.InactiveThreshold)
	}
}
