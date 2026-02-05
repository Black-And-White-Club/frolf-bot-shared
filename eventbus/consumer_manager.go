package eventbus

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"

	"github.com/nats-io/nats.go/jetstream"
)

// ConsumerManager handles the creation and caching of JetStream pull consumers.
// It ensures idempotent consumer creation using CreateOrUpdateConsumer.
type ConsumerManager struct {
	js       jetstream.JetStream
	logger   *slog.Logger
	registry *ConsumerConfigRegistry

	mu        sync.RWMutex
	consumers map[string]jetstream.Consumer
}

// NewConsumerManager creates a new ConsumerManager.
func NewConsumerManager(js jetstream.JetStream, logger *slog.Logger, registry *ConsumerConfigRegistry) *ConsumerManager {
	if registry == nil {
		registry = NewConsumerConfigRegistry()
	}
	return &ConsumerManager{
		js:        js,
		logger:    logger,
		registry:  registry,
		consumers: make(map[string]jetstream.Consumer),
	}
}

// EnsureConsumer creates or updates a durable consumer for the given stream and topic.
// It uses CreateOrUpdateConsumer for idempotent operations and caches the result.
// Returns the consumer, which can be used to create message iterators.
func (cm *ConsumerManager) EnsureConsumer(ctx context.Context, streamName, topic, appType string) (jetstream.Consumer, error) {
	consumerName := buildConsumerName(appType, topic)

	// Check cache first
	cm.mu.RLock()
	if cons, ok := cm.consumers[consumerName]; ok {
		cm.mu.RUnlock()
		return cons, nil
	}
	cm.mu.RUnlock()

	// Need to create/update - acquire write lock
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// Double-check after acquiring write lock
	if cons, ok := cm.consumers[consumerName]; ok {
		return cons, nil
	}

	ctxLogger := cm.logger.With(
		"operation", "ensure_consumer",
		"stream", streamName,
		"topic", topic,
		"app_type", appType,
		"consumer_name", consumerName,
	)

	// Resolve configuration for this app/topic combination
	cfg := cm.registry.Resolve(appType, topic)

	consumerConfig := jetstream.ConsumerConfig{
		Durable:       consumerName,
		FilterSubject: topic,

		AckPolicy:     jetstream.AckExplicitPolicy,
		AckWait:       cfg.AckWait,
		MaxDeliver:    cfg.MaxDeliver,
		BackOff:       cfg.BackOff,
		MaxAckPending: cfg.MaxAckPending,

		DeliverPolicy: cfg.DeliverPolicy,
		ReplayPolicy:  jetstream.ReplayInstantPolicy,

		InactiveThreshold: cfg.InactiveThreshold,
	}

	ctxLogger.DebugContext(ctx, "Creating or updating consumer",
		"ack_wait", cfg.AckWait,
		"max_deliver", cfg.MaxDeliver,
		"max_ack_pending", cfg.MaxAckPending,
	)

	// CreateOrUpdateConsumer is idempotent - safe to call multiple times
	cons, err := cm.js.CreateOrUpdateConsumer(ctx, streamName, consumerConfig)
	if err != nil {
		ctxLogger.ErrorContext(ctx, "Failed to create or update consumer", "error", err)
		return nil, fmt.Errorf("failed to create or update consumer %s on stream %s: %w", consumerName, streamName, err)
	}

	// Cache the consumer
	cm.consumers[consumerName] = cons

	ctxLogger.InfoContext(ctx, "Consumer ensured successfully")
	return cons, nil
}

// GetConsumerInfo returns information about a cached consumer.
// Returns nil if the consumer is not cached.
func (cm *ConsumerManager) GetConsumerInfo(ctx context.Context, appType, topic string) (*jetstream.ConsumerInfo, error) {
	consumerName := buildConsumerName(appType, topic)

	cm.mu.RLock()
	cons, ok := cm.consumers[consumerName]
	cm.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("consumer %s not found in cache", consumerName)
	}

	return cons.Info(ctx)
}

// GetRegistry returns the config registry for external configuration.
func (cm *ConsumerManager) GetRegistry() *ConsumerConfigRegistry {
	return cm.registry
}

// buildConsumerName creates a standardized consumer name from app type and topic.
// Format: {appType}-{sanitized-topic}
// Sanitization replaces: dots with dashes, * with star, > with all
func buildConsumerName(appType, topic string) string {
	sanitized := topic
	sanitized = strings.ReplaceAll(sanitized, ".", "-")
	sanitized = strings.ReplaceAll(sanitized, "*", "star")
	sanitized = strings.ReplaceAll(sanitized, ">", "all")

	return fmt.Sprintf("%s-%s", appType, sanitized)
}

// ResolveStreamFromTopic determines the stream name from a topic.
// This centralizes stream resolution logic.
func ResolveStreamFromTopic(topic string) (string, error) {
	switch {
	case strings.HasPrefix(topic, "user."):
		return "user", nil
	case strings.HasPrefix(topic, "leaderboard."):
		return "leaderboard", nil
	case strings.HasPrefix(topic, "round."):
		return "round", nil
	case strings.HasPrefix(topic, "score."):
		return "score", nil
	case strings.HasPrefix(topic, "guild."):
		return "guild", nil
	case strings.HasPrefix(topic, "discord."):
		return "discord", nil
	case strings.HasPrefix(topic, "auth."):
		return "auth", nil
	case strings.HasPrefix(topic, "club."):
		return "club", nil
	default:
		return "", fmt.Errorf("unknown topic prefix: %s", topic)
	}
}
