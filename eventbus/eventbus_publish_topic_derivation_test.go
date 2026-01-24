package eventbus

import (
	"io"
	"log/slog"
	"testing"

	"github.com/ThreeDotsLabs/watermill-nats/v2/pkg/nats"
	"github.com/ThreeDotsLabs/watermill/message"
)

type fakePublisher struct {
	calls []publishCall
}

type publishCall struct {
	topic string
	msgs  []*message.Message
}

func (f *fakePublisher) Publish(topic string, messages ...*message.Message) error {
	f.calls = append(f.calls, publishCall{topic: topic, msgs: messages})
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

	if len(fp.calls) != 1 {
		t.Fatalf("expected 1 publish call, got %d", len(fp.calls))
	}
	if fp.calls[0].topic != "leaderboard.tag.lookup.by.user.id.success.v1" {
		t.Fatalf("expected publisher topic to be derived topic, got %s", fp.calls[0].topic)
	}
}

func TestPublish_GroupsMessagesByTopic(t *testing.T) {
	fp := &fakePublisher{}
	eb := &eventBus{
		appType:   "backend",
		publisher: fp,
		logger:    slog.New(slog.NewTextHandler(io.Discard, nil)),
		marshaler: &nats.NATSMarshaler{},
	}

	// Create messages with different topics (simulating handler returning multiple results)
	msg1 := message.NewMessage("id1", []byte("payload1"))
	msg1.Metadata.Set("topic", "round.finalized.discord.v1")

	msg2 := message.NewMessage("id2", []byte("payload2"))
	msg2.Metadata.Set("topic", "round.finalized.v1")

	msg3 := message.NewMessage("id3", []byte("payload3"))
	msg3.Metadata.Set("topic", "round.finalized.discord.v1") // Same as msg1

	// Publish all messages with empty topic (dynamic routing)
	if err := eb.Publish("", msg1, msg2, msg3); err != nil {
		t.Fatalf("expected publish to succeed, got error: %v", err)
	}

	// Should have 2 separate publish calls (one per unique topic)
	if len(fp.calls) != 2 {
		t.Fatalf("expected 2 publish calls (one per topic), got %d", len(fp.calls))
	}

	// Verify topics were separated correctly
	topicCounts := make(map[string]int)
	for _, call := range fp.calls {
		topicCounts[call.topic] += len(call.msgs)
	}

	if topicCounts["round.finalized.discord.v1"] != 2 {
		t.Errorf("expected 2 messages to discord topic, got %d", topicCounts["round.finalized.discord.v1"])
	}
	if topicCounts["round.finalized.v1"] != 1 {
		t.Errorf("expected 1 message to backend topic, got %d", topicCounts["round.finalized.v1"])
	}
}

func TestPublish_FailsOnMissingTopic(t *testing.T) {
	fp := &fakePublisher{}
	eb := &eventBus{
		appType:   "backend",
		publisher: fp,
		logger:    slog.New(slog.NewTextHandler(io.Discard, nil)),
		marshaler: &nats.NATSMarshaler{},
	}

	msg := message.NewMessage("id", []byte("payload"))
	// No topic set

	err := eb.Publish("", msg)
	if err == nil {
		t.Fatal("expected error when message has no topic metadata")
	}
}
