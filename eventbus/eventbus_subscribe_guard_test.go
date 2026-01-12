package eventbus

import (
	"context"
	"io"
	"log/slog"
	"strings"
	"testing"

	discordleaderboard "github.com/Black-And-White-Club/frolf-bot-shared/events/discord/leaderboard"
)

func TestSubscribe_DiscordTopicForbiddenForBackend(t *testing.T) {
	eb := &eventBus{
		appType: "backend",
		logger:  slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
	ctx := context.Background()
	ch, err := eb.Subscribe(ctx, discordleaderboard.LeaderboardRetrieveRequestV1)
	if err == nil {
		t.Fatalf("expected error when subscribing to discord topic as backend")
	}
	if ch != nil {
		t.Fatalf("expected nil channel on forbidden subscribe")
	}
	if !strings.Contains(err.Error(), "forbidden") {
		t.Fatalf("unexpected error message: %v", err)
	}
}

func TestSubscribeForTest_DiscordTopicForbiddenForBackend(t *testing.T) {
	eb := &eventBus{
		appType: "backend",
		logger:  slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
	ctx := context.Background()
	ch, err := eb.SubscribeForTest(ctx, discordleaderboard.LeaderboardRetrieveRequestV1)
	if err == nil {
		t.Fatalf("expected error when subscribing for test to discord topic as backend")
	}
	if ch != nil {
		t.Fatalf("expected nil channel on forbidden subscribe for test")
	}
	if !strings.Contains(err.Error(), "forbidden") {
		t.Fatalf("unexpected error message: %v", err)
	}
}

func TestSubscribe_DiscordLookupResponseForbiddenForBackend(t *testing.T) {
	eb := &eventBus{
		appType: "backend",
		logger:  slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
	ctx := context.Background()
	ch, err := eb.Subscribe(ctx, discordleaderboard.LeaderboardTagLookupSucceededV1)
	if err == nil {
		t.Fatalf("expected error when subscribing to discord lookup response topic as backend")
	}
	if ch != nil {
		t.Fatalf("expected nil channel on forbidden subscribe")
	}
	if !strings.Contains(err.Error(), "forbidden") {
		t.Fatalf("unexpected error message: %v", err)
	}
}

func TestSubscribeForTest_DiscordLookupResponseForbiddenForBackend(t *testing.T) {
	eb := &eventBus{
		appType: "backend",
		logger:  slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
	ctx := context.Background()
	ch, err := eb.SubscribeForTest(ctx, discordleaderboard.LeaderboardTagLookupSucceededV1)
	if err == nil {
		t.Fatalf("expected error when subscribing for test to discord lookup response topic as backend")
	}
	if ch != nil {
		t.Fatalf("expected nil channel on forbidden subscribe for test")
	}
	if !strings.Contains(err.Error(), "forbidden") {
		t.Fatalf("unexpected error message: %v", err)
	}
}
