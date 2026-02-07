package eventbus

import (
	"sync"
	"time"

	"github.com/nats-io/nats.go/jetstream"
)

// ConsumerConfig holds configuration for JetStream pull consumers.
// These settings control message delivery, acknowledgment, and retry behavior.
type ConsumerConfig struct {
	// AckWait is the duration the server waits for an acknowledgment before redelivering.
	AckWait time.Duration

	// MaxDeliver is the maximum number of delivery attempts for a message.
	// After this many attempts, the message will be terminated.
	MaxDeliver int

	// BackOff defines the delay between redelivery attempts.
	// Each entry corresponds to a redelivery attempt (1st retry, 2nd retry, etc.).
	BackOff []time.Duration

	// MaxAckPending limits the number of unacknowledged messages the consumer will allow.
	// This provides natural backpressure.
	MaxAckPending int

	// DeliverPolicy determines where in the stream to start delivering messages.
	DeliverPolicy jetstream.DeliverPolicy

	// InactiveThreshold is the duration after which an inactive consumer may be cleaned up.
	// Set to 0 to disable automatic cleanup.
	InactiveThreshold time.Duration
}

// DefaultConsumerConfig returns production-safe default configuration.
// These defaults balance reliability with reasonable resource usage.
func DefaultConsumerConfig() ConsumerConfig {
	return ConsumerConfig{
		AckWait:    60 * time.Second,
		MaxDeliver: 5,
		BackOff: []time.Duration{
			5 * time.Second,
			15 * time.Second,
			30 * time.Second,
			60 * time.Second,
		},
		MaxAckPending:     100,
		DeliverPolicy:     jetstream.DeliverNewPolicy,
		InactiveThreshold: 5 * time.Minute,
	}
}

// ConsumerConfigRegistry manages consumer configurations with support for
// defaults, app-specific overrides, and topic-specific overrides.
// Configuration resolution follows the priority: topic-specific > app-specific > defaults.
type ConsumerConfigRegistry struct {
	mu           sync.RWMutex
	defaults     ConsumerConfig
	appConfigs   map[string]ConsumerConfig
	topicConfigs map[string]ConsumerConfig
}

// NewConsumerConfigRegistry creates a new registry with production-safe defaults.
func NewConsumerConfigRegistry() *ConsumerConfigRegistry {
	return &ConsumerConfigRegistry{
		defaults:     DefaultConsumerConfig(),
		appConfigs:   make(map[string]ConsumerConfig),
		topicConfigs: make(map[string]ConsumerConfig),
	}
}

// SetDefault overrides the default configuration used when no specific config is found.
func (r *ConsumerConfigRegistry) SetDefault(cfg ConsumerConfig) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.defaults = cfg
}

// SetForApp sets a configuration override for a specific application type.
// This configuration takes precedence over defaults but not over topic-specific configs.
func (r *ConsumerConfigRegistry) SetForApp(appType string, cfg ConsumerConfig) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.appConfigs[appType] = cfg
}

// SetForTopic sets a configuration override for a specific topic.
// Topic-specific configurations take the highest precedence.
func (r *ConsumerConfigRegistry) SetForTopic(topic string, cfg ConsumerConfig) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.topicConfigs[topic] = cfg
}

// Resolve returns the appropriate configuration for the given app and topic.
// Resolution priority: topic-specific > app-specific > defaults.
func (r *ConsumerConfigRegistry) Resolve(appType, topic string) ConsumerConfig {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Highest priority: topic-specific config
	if cfg, ok := r.topicConfigs[topic]; ok {
		return cfg
	}

	// Second priority: app-specific config
	if cfg, ok := r.appConfigs[appType]; ok {
		return cfg
	}

	// Fallback: defaults
	return r.defaults
}

// GetDefault returns the current default configuration.
func (r *ConsumerConfigRegistry) GetDefault() ConsumerConfig {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.defaults
}
