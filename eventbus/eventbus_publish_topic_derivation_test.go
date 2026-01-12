package eventbus

import (
	"io"
	"log/slog"
	"testing"

	"github.com/ThreeDotsLabs/watermill-nats/v2/pkg/nats"
	"github.com/ThreeDotsLabs/watermill/message"
)

type fakePublisher struct {
	lastTopic string
	msgs      []*message.Message
}

func (f *fakePublisher) Publish(topic string, messages ...*message.Message) error {
	f.lastTopic = topic
	f.msgs = append(f.msgs, messages...)
	return nil
}

func (f *fakePublisher) Close() error { return nil }

func TestPublish_DerivesTopicFromMessageMetadataEvenWithoutMetrics(t *testing.T) {
	fp := &fakePublisher{}
	eb := &eventBus{
		appType:   "backend",
		publisher: fp,
		logger:    slog.New(slog.NewTextHandler(io.Discard, nil)),
		marshaler: &nats.NATSMarshaler{},
	}

	msg := message.NewMessage("id", []byte("payload"))
	msg.Metadata.Set("topic", "leaderboard.tag.lookup.by.user.id.success.v1")

	if err := eb.Publish("", msg); err != nil {
		t.Fatalf("expected publish to derive topic from message metadata, got error: %v", err)
	}

	if fp.lastTopic != "leaderboard.tag.lookup.by.user.id.success.v1" {
		t.Fatalf("expected publisher.lastTopic to be derived topic, got %s", fp.lastTopic)
	}
}
