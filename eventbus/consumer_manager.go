package eventbus

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"
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

	// Double-check under write lock before issuing a network call.
	cm.mu.Lock()
	if cons, ok := cm.consumers[consumerName]; ok {
		cm.mu.Unlock()
		return cons, nil
	}
	cm.mu.Unlock()

	ctxLogger := cm.logger.With(
		"operation", "ensure_consumer",
		"stream", streamName,
		"topic", topic,
		"app_type", appType,
		"consumer_name", consumerName,
	)

	// Resolve configuration for this app/topic combination
	cfg := cm.registry.Resolve(appType, topic)

	consumerConfig := buildDurableConsumerConfig(consumerName, topic, cfg)

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

	info, err := cons.Info(ctx)
	if err != nil {
		ctxLogger.ErrorContext(ctx, "Failed to fetch consumer info after create/update", "error", err)
		return nil, fmt.Errorf("failed to fetch consumer info %s on stream %s: %w", consumerName, streamName, err)
	}
	if err := validateConsumerInfo(info, consumerConfig); err != nil {
		ctxLogger.ErrorContext(ctx, "Consumer configuration validation failed", "error", err)
		return nil, fmt.Errorf("consumer %s validation failed: %w", consumerName, err)
	}

	// Cache the consumer (in case another goroutine created one while we were in-flight,
	// prefer the first cached instance to keep shared behavior consistent).
	cm.mu.Lock()
	if existing, ok := cm.consumers[consumerName]; ok {
		cm.mu.Unlock()
		return existing, nil
	}
	cm.consumers[consumerName] = cons
	cm.mu.Unlock()

	ctxLogger.InfoContext(ctx, "Consumer ensured successfully")
	return cons, nil
}

func buildDurableConsumerConfig(consumerName, topic string, cfg ConsumerConfig) jetstream.ConsumerConfig {
	return jetstream.ConsumerConfig{
		Durable:       consumerName,
		FilterSubject: topic,
		AckPolicy:     jetstream.AckExplicitPolicy,
		AckWait:       cfg.AckWait,
		MaxDeliver:    cfg.MaxDeliver,
		BackOff:       cfg.BackOff,
		MaxAckPending: cfg.MaxAckPending,
		DeliverPolicy: cfg.DeliverPolicy,
		ReplayPolicy:  jetstream.ReplayInstantPolicy,
		// Keep durables replay-safe by default; non-zero enables server-side cleanup.
		InactiveThreshold: cfg.InactiveThreshold,
	}
}

func validateConsumerInfo(info *jetstream.ConsumerInfo, expected jetstream.ConsumerConfig) error {
	if info == nil {
		return fmt.Errorf("consumer info is nil")
	}
	actual := info.Config

	if actual.Durable != expected.Durable {
		return fmt.Errorf("durable mismatch: expected %q, got %q", expected.Durable, actual.Durable)
	}
	if actual.FilterSubject != expected.FilterSubject {
		return fmt.Errorf("filter subject mismatch: expected %q, got %q", expected.FilterSubject, actual.FilterSubject)
	}
	if actual.DeliverPolicy != expected.DeliverPolicy {
		return fmt.Errorf("deliver policy mismatch: expected %v, got %v", expected.DeliverPolicy, actual.DeliverPolicy)
	}
	if actual.InactiveThreshold != expected.InactiveThreshold {
		return fmt.Errorf("inactive threshold mismatch: expected %s, got %s", expected.InactiveThreshold, actual.InactiveThreshold)
	}
	if actual.AckPolicy != expected.AckPolicy {
		return fmt.Errorf("ack policy mismatch: expected %v, got %v", expected.AckPolicy, actual.AckPolicy)
	}
	if actual.AckWait != expected.AckWait {
		return fmt.Errorf("ack wait mismatch: expected %s, got %s", expected.AckWait, actual.AckWait)
	}
	if actual.MaxDeliver != expected.MaxDeliver {
		return fmt.Errorf("max deliver mismatch: expected %d, got %d", expected.MaxDeliver, actual.MaxDeliver)
	}
	if actual.MaxAckPending != expected.MaxAckPending {
		return fmt.Errorf("max ack pending mismatch: expected %d, got %d", expected.MaxAckPending, actual.MaxAckPending)
	}
	if !reflect.DeepEqual(actual.BackOff, expected.BackOff) {
		return fmt.Errorf("backoff mismatch: expected %v, got %v", expected.BackOff, actual.BackOff)
	}

	return nil
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
