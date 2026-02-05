package eventbus

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Black-And-White-Club/frolf-bot-shared/observability/attr"
	lokifrolfbot "github.com/Black-And-White-Club/frolf-bot-shared/observability/otel/logging"
	eventbusmetrics "github.com/Black-And-White-Club/frolf-bot-shared/observability/otel/metrics/eventbus"
	"github.com/ThreeDotsLabs/watermill-nats/v2/pkg/nats"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	nc "github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel/trace"
)

type eventBus struct {
	appType           string
	publisher         message.Publisher
	subscriber        message.Subscriber // Watermill subscriber (kept for compatibility)
	subscriberAdapter *JetStreamSubscriberAdapter
	consumerManager   *ConsumerManager
	js                jetstream.JetStream
	natsConn          *nc.Conn
	logger            *slog.Logger
	createdStreams    map[string]bool
	streamMutex       sync.Mutex
	marshaler         nats.Marshaler
	metrics           eventbusmetrics.EventBusMetrics
	tracer            trace.Tracer
}

// EventBus interface
type EventBus interface {
	Publish(topic string, messages ...*message.Message) error
	Subscribe(ctx context.Context, topic string) (<-chan *message.Message, error)
	Close() error
	GetNATSConnection() *nc.Conn
	GetJetStream() jetstream.JetStream
	GetHealthCheckers() []HealthChecker
	CreateStream(ctx context.Context, streamName string) error
	SubscribeForTest(ctx context.Context, topic string) (<-chan *message.Message, error)
}

// HealthChecker interface for components that can be health-checked
type HealthChecker interface {
	Check(ctx context.Context) error
	Name() string
}

func NewEventBus(ctx context.Context, natsURL string, logger *slog.Logger, appType string, metrics eventbusmetrics.EventBusMetrics, tracer trace.Tracer) (EventBus, error) {
	ctxLogger := logger.With(
		"operation", "new_event_bus",
		"nats_url", natsURL,
		"app_type", appType,
	)

	ctxLogger.InfoContext(ctx, "Creating new EventBus")

	natsConn, err := nc.Connect(natsURL,
		nc.RetryOnFailedConnect(true),
		nc.Timeout(30*time.Second),
		nc.ReconnectWait(1*time.Second),
		nc.MaxReconnects(-1),
	)
	if err != nil {
		ctxLogger.ErrorContext(ctx, "Failed to connect to NATS", "error", err)
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	js, err := jetstream.New(natsConn)
	if err != nil {
		natsConn.Close()
		ctxLogger.ErrorContext(ctx, "Failed to initialize JetStream", "error", err)
		return nil, fmt.Errorf("failed to initialize JetStream: %w", err)
	}

	marshaller := &nats.NATSMarshaler{}
	watermillLogger := lokifrolfbot.ToWatermillAdapter(logger)

	publisher, err := nats.NewPublisher(
		nats.PublisherConfig{
			URL:       natsURL,
			Marshaler: marshaller,
			NatsOptions: []nc.Option{
				nc.RetryOnFailedConnect(true),
				nc.Timeout(30 * time.Second),
				nc.ReconnectWait(1 * time.Second),
				nc.MaxReconnects(-1),
			},
			JetStream: nats.JetStreamConfig{
				Disabled:      false,
				AutoProvision: false,
				TrackMsgId:    true,
				DurablePrefix: "durable",
				DurableCalculator: func(durablePrefix, topic string) string {
					sanitizedTopic := strings.ReplaceAll(topic, ".", "_")
					return fmt.Sprintf("%s-%s", durablePrefix, sanitizedTopic)
				},
			},
		},
		watermillLogger,
	)
	if err != nil {
		natsConn.Close()
		ctxLogger.ErrorContext(ctx, "Failed to create Watermill publisher", "error", err)
		return nil, fmt.Errorf("failed to create Watermill publisher: %w", err)
	}

	// Create consumer config registry with production defaults
	registry := NewConsumerConfigRegistry()

	// Create consumer manager for managing JetStream consumers
	consumerManager := NewConsumerManager(js, logger, registry)

	// Create the native JetStream subscriber adapter
	subscriberAdapter := NewJetStreamSubscriberAdapter(
		js,
		consumerManager,
		appType,
		logger,
		WithMaxConcurrentAcks(50),
	)

	eventBus := &eventBus{
		appType:           appType,
		publisher:         publisher,
		subscriber:        nil, // No longer using Watermill subscriber
		subscriberAdapter: subscriberAdapter,
		consumerManager:   consumerManager,
		js:                js,
		natsConn:          natsConn,
		logger:            logger,
		createdStreams:    make(map[string]bool),
		streamMutex:       sync.Mutex{},
		marshaler:         marshaller,
		metrics:           metrics,
		tracer:            tracer,
	}

	if err := eventBus.createStreamsForApp(ctx, appType); err != nil {
		natsConn.Close()
		ctxLogger.ErrorContext(ctx, "Failed to create streams for app", "error", err)
		return nil, fmt.Errorf("failed to create streams for app: %w", err)
	}

	ctxLogger.InfoContext(ctx, "EventBus created successfully")
	return eventBus, nil
}

// Publish publishes messages using Watermill with idempotency support.
// When topic is empty, messages are grouped by their metadata topic and published separately.
func (eb *eventBus) Publish(topic string, messages ...*message.Message) error {
	if len(messages) == 0 {
		return nil
	}

	if eb.marshaler == nil {
		return errors.New("eventBus marshaler is not set")
	}

	// When topic is provided, publish all messages to that topic
	if topic != "" {
		return eb.publishToTopic(topic, messages)
	}

	// Group messages by their metadata topic for correct routing
	topicGroups := make(map[string][]*message.Message)
	for _, msg := range messages {
		msgTopic := resolveMessageTopic(msg)
		if msgTopic == "" {
			return fmt.Errorf("message %s has no topic in metadata", msg.UUID)
		}
		topicGroups[msgTopic] = append(topicGroups[msgTopic], msg)
	}

	// Publish each group to its correct topic
	var errs []error
	for t, msgs := range topicGroups {
		if err := eb.publishToTopic(t, msgs); err != nil {
			errs = append(errs, fmt.Errorf("topic %s: %w", t, err))
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

// resolveMessageTopic extracts the topic from message metadata.
func resolveMessageTopic(msg *message.Message) string {
	for _, key := range []string{"topic", "Topic", "event_name", "topic_hint"} {
		if t := msg.Metadata.Get(key); t != "" {
			return t
		}
	}
	return ""
}

// publishToTopic publishes messages to a specific topic with retries.
func (eb *eventBus) publishToTopic(topic string, messages []*message.Message) error {
	if len(messages) == 0 {
		return nil
	}

	// Boundary guard for discord topics
	if strings.HasPrefix(topic, "discord.") && eb.appType != "discord" {
		return fmt.Errorf("publishing to discord topics forbidden for app %q", eb.appType)
	}

	ctxLogger := eb.logger.With(
		"operation", "publish",
		"topic", topic,
		"message_count", len(messages),
	)

	// Check if this is an inbox subject (core NATS request-reply pattern)
	// Inbox subjects start with "_INBOX." and are ephemeral - they don't belong to JetStream
	if strings.HasPrefix(topic, "_INBOX.") {
		return eb.publishToInbox(topic, messages, ctxLogger)
	}

	// Set deduplication IDs
	for _, msg := range messages {
		if msg.Metadata.Get("Nats-Msg-Id") == "" {
			if key := msg.Metadata.Get("idempotency_key"); key != "" {
				msg.Metadata.Set("Nats-Msg-Id", dedupeMsgID(key, topic))
			} else {
				msg.Metadata.Set("Nats-Msg-Id", msg.UUID)
			}
		}
		ctxLogger.Debug("Publishing message",
			attr.String("message_uuid", msg.UUID),
			attr.String("nats_msg_id", msg.Metadata.Get("Nats-Msg-Id")),
		)
	}

	if eb.metrics != nil {
		eb.metrics.RecordMessagePublish(messages[0].Context(), topic)
	}

	// Publish with retry
	const maxRetries = 3
	var lastErr error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		if err := eb.publisher.Publish(topic, messages...); err != nil {
			lastErr = err
			if attempt < maxRetries && eb.isRetryableError(err) {
				time.Sleep(time.Duration(attempt) * 100 * time.Millisecond)
				continue
			}
			if eb.metrics != nil {
				eb.metrics.RecordMessagePublishError(messages[0].Context(), topic)
			}
			ctxLogger.Error("Publish failed", "error", err, "attempts", attempt)
			return fmt.Errorf("publish to %s failed after %d attempts: %w", topic, attempt, lastErr)
		}
		break
	}

	for _, msg := range messages {
		ctxLogger.InfoContext(msg.Context(), "Message published",
			"message_id", msg.UUID,
			"correlation_id", middleware.MessageCorrelationID(msg),
		)
	}

	return nil
}

// publishToInbox publishes messages to a core NATS inbox subject.
// Inbox subjects are used for request-reply patterns and don't belong to JetStream.
func (eb *eventBus) publishToInbox(topic string, messages []*message.Message, ctxLogger *slog.Logger) error {
	for _, msg := range messages {
		ctxLogger.Debug("Publishing to inbox via core NATS",
			attr.String("message_uuid", msg.UUID),
			attr.String("inbox", topic),
		)

		// Publish directly via core NATS connection
		if err := eb.natsConn.Publish(topic, msg.Payload); err != nil {
			ctxLogger.Error("Failed to publish to inbox", "error", err)
			return fmt.Errorf("failed to publish to inbox %s: %w", topic, err)
		}

		ctxLogger.Info("Message published to inbox",
			"message_id", msg.UUID,
			"inbox", topic,
		)
	}

	return nil
}

func dedupeMsgID(idempotencyKey, topic string) string {
	// Include topic to avoid collisions across subjects within the same stream.
	seed := topic + "|" + idempotencyKey
	sum := sha256.Sum256([]byte(seed))
	return "idem-" + hex.EncodeToString(sum[:])
}

// Subscribe subscribes to a topic using the native JetStream subscriber adapter.
// It provides bounded ACK concurrency and graceful shutdown support.
func (eb *eventBus) Subscribe(ctx context.Context, topic string) (<-chan *message.Message, error) {
	ctxLogger := eb.logger.With(
		"operation", "subscribe",
		"topic", topic,
	)

	if eb.metrics != nil {
		eb.metrics.RecordMessageSubscribe(ctx, topic)
	}

	// Boundary guard: prevent non-discord apps from subscribing to discord.* topics
	if strings.HasPrefix(topic, "discord.") && eb.appType != "discord" {
		ctxLogger.Error("Subscribe forbidden: app not allowed to subscribe to discord topics",
			"app_type", eb.appType,
			"topic", topic,
		)
		return nil, fmt.Errorf("subscription to discord topics forbidden for app %q", eb.appType)
	}

	// Resolve stream from topic and ensure it exists
	streamName, err := ResolveStreamFromTopic(topic)
	if err != nil {
		ctxLogger.ErrorContext(ctx, "Failed to resolve stream from topic", "error", err)
		return nil, err
	}

	if err := eb.CreateStream(ctx, streamName); err != nil {
		ctxLogger.ErrorContext(ctx, "Failed to create stream", "error", err, "stream", streamName)
		return nil, err
	}

	// Delegate to the subscriber adapter
	return eb.subscriberAdapter.Subscribe(ctx, topic)
}

// SubscribeForTest creates an ephemeral (non-durable) subscription for testing.
// Unlike Subscribe, this creates a non-durable consumer that is automatically
// cleaned up when the subscription is closed.
func (eb *eventBus) SubscribeForTest(ctx context.Context, topic string) (<-chan *message.Message, error) {
	// Boundary guard: prevent non-discord apps from subscribing to discord.* topics
	if strings.HasPrefix(topic, "discord.") && eb.appType != "discord" {
		return nil, fmt.Errorf("subscription to discord topics forbidden for app %q", eb.appType)
	}

	streamName, err := ResolveStreamFromTopic(topic)
	if err != nil {
		return nil, err
	}

	if err := eb.CreateStream(ctx, streamName); err != nil {
		return nil, err
	}

	cons, err := eb.js.CreateConsumer(ctx, streamName, jetstream.ConsumerConfig{
		Durable:       "",
		FilterSubject: topic,
		AckPolicy:     jetstream.AckExplicitPolicy,
		DeliverPolicy: jetstream.DeliverNewPolicy,
		MaxDeliver:    1,
		AckWait:       30 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	sub, err := cons.Messages()
	if err != nil {
		return nil, err
	}

	out := make(chan *message.Message, 16)

	go func() {
		defer func() {
			sub.Stop()
			close(out)
		}()

		for {
			jsMsg, err := sub.Next()
			if err != nil {
				return
			}

			wmMsg, err := eb.toWatermillMessage(ctx, jsMsg)
			if err != nil {
				_ = jsMsg.Term()
				continue
			}

			select {
			case out <- wmMsg:
				go func() {
					select {
					case <-wmMsg.Acked():
						_ = jsMsg.Ack()
					case <-ctx.Done():
					}
				}()
			case <-ctx.Done():
				return
			}
		}
	}()

	return out, nil
}

// toWatermillMessage converts a JetStream message to a Watermill message.
func (eb *eventBus) toWatermillMessage(parentCtx context.Context, jsMsg jetstream.Msg) (*message.Message, error) {
	msgID := jsMsg.Headers().Get("Nats-Msg-Id")
	if msgID == "" {
		msgID = nc.NewInbox()
	}

	msg := message.NewMessage(msgID, jsMsg.Data())

	// Copy headers
	for k, v := range jsMsg.Headers() {
		msg.Metadata.Set(k, v[0])
	}

	if meta, err := jsMsg.Metadata(); err == nil {
		msg.Metadata.Set("stream", meta.Stream)
		msg.Metadata.Set("consumer", meta.Consumer)
		msg.Metadata.Set("deliveries", strconv.FormatUint(meta.NumDelivered, 10))
		msg.Metadata.Set("stream_seq", strconv.FormatUint(meta.Sequence.Stream, 10))
		msg.Metadata.Set("consumer_seq", strconv.FormatUint(meta.Sequence.Consumer, 10))
	}

	// IMPORTANT: derive context, do not overwrite
	msg.SetContext(context.WithValue(parentCtx, "message_id", msg.UUID))
	return msg, nil
}

// sanitizeForNATS sanitizes a string for use in NATS topics.
func sanitizeForNATS(s string) string {
	reg := regexp.MustCompile(`[^a-zA-Z0-9-]+`)
	return reg.ReplaceAllString(strings.ReplaceAll(s, ".", "-"), "")
}

// Close closes the EventBus and its components.
// It closes the subscriber adapter first to allow graceful drain of in-flight messages,
// then the publisher, and finally drains the NATS connection.
func (eb *eventBus) Close() error {
	ctxLogger := eb.logger.With(attr.String("operation", "close"))

	ctxLogger.Info("Closing EventBus")

	// Close subscriber adapter first to gracefully drain in-flight messages
	if eb.subscriberAdapter != nil {
		if err := eb.subscriberAdapter.Close(); err != nil {
			ctxLogger.Error("Error closing subscriber adapter", "component", "subscriber_adapter", "error", err)
		} else {
			ctxLogger.Info("Subscriber adapter closed successfully", "component", "subscriber_adapter")
		}
	}

	// Close legacy subscriber if present (for backward compatibility)
	if eb.subscriber != nil {
		if err := eb.subscriber.Close(); err != nil {
			ctxLogger.Error("Error closing subscriber", "component", "subscriber", "error", err)
		} else {
			ctxLogger.Info("Subscriber closed successfully", "component", "subscriber")
		}
	}

	// Close publisher
	if eb.publisher != nil {
		if err := eb.publisher.Close(); err != nil {
			ctxLogger.Error("Error closing publisher", "component", "publisher", "error", err)
		} else {
			ctxLogger.Info("Publisher closed successfully", "component", "publisher")
		}
	}

	// Drain and close NATS connection
	if eb.natsConn != nil {
		// Drain allows in-flight messages to complete
		if err := eb.natsConn.Drain(); err != nil {
			ctxLogger.Warn("Error draining NATS connection, forcing close", "error", err)
			eb.natsConn.Close()
		}
		ctxLogger.Info("NATS connection closed", "component", "nats_connection")
	}

	ctxLogger.Info("EventBus closed successfully")
	return nil
}

// CreateStream creates or updates a single JetStream stream (helper function).
func (eb *eventBus) CreateStream(ctx context.Context, streamName string) error {
	ctxLogger := eb.logger.With(
		attr.String("operation", "create_stream"),
		attr.String("stream_name", streamName),
	)

	ctxLogger.Info("Creating/Updating stream")
	eb.streamMutex.Lock()
	defer eb.streamMutex.Unlock()

	if eb.createdStreams[streamName] {
		ctxLogger.Info("Stream already created in this process")
		return nil
	}

	var subjects []string
	switch streamName {
	case "user":
		subjects = []string{"user.>"}
	case "leaderboard":
		subjects = []string{"leaderboard.>"}
	case "round":
		subjects = []string{"round.>"}
	case "score":
		subjects = []string{"score.>"}
	case "discord":
		subjects = []string{"discord.>"}
	case "guild":
		subjects = []string{"guild.>"}
	case "auth":
		subjects = []string{"auth.>"}
	case "club":
		subjects = []string{"club.>"}
	default:
		ctxLogger.Error("Failed to create stream", "error", "unknown stream name")
		return fmt.Errorf("unknown stream name: %s", streamName)
	}

	ctxLogger = ctxLogger.With(attr.Any("subjects", subjects))

	// Define stream configuration with better durability and deduplication
	streamCfg := jetstream.StreamConfig{
		Name:         streamName,
		Subjects:     subjects,
		Retention:    jetstream.LimitsPolicy,
		Duplicates:   10 * time.Minute,      // Extended deduplication window
		MaxAge:       24 * time.Hour,        // Keep messages for 24 hours
		MaxBytes:     1024 * 1024 * 100,     // 100MB max stream size
		MaxMsgs:      100000,                // Max 100k messages
		NoAck:        false,                 // Require acknowledgments
		Discard:      jetstream.DiscardOld,  // Discard old messages when limits hit
		Storage:      jetstream.FileStorage, // Persistent storage
		Replicas:     1,                     // Single replica for now
		MaxConsumers: -1,                    // Unlimited consumers
		MaxMsgSize:   1024 * 1024,           // 1MB max message size
	}

	ctxLogger = ctxLogger.With(attr.Duration("duplicates_window", streamCfg.Duplicates))

	// Create or update the stream (idempotent)
	stream, err := eb.js.CreateOrUpdateStream(ctx, streamCfg)
	if err != nil {
		ctxLogger.Error("Failed to create or update stream", "error", err)
		return fmt.Errorf("failed to create or update stream %s: %w", streamName, err)
	}

	ctxLogger.Info("Stream created or updated successfully", "stream_name", stream.CachedInfo().Config.Name)
	eb.createdStreams[streamName] = true
	return nil
}

// createStreamsForApp creates streams for a specific application type.
func (eb *eventBus) createStreamsForApp(ctx context.Context, appType string) error {
	ctxLogger := eb.logger.With(
		attr.String("operation", "create_streams_for_app"),
		attr.String("app_type", appType),
	)

	ctxLogger.Info("Creating streams for application type")

	var streams []string
	switch appType {
	case "backend":
		streams = []string{"user", "leaderboard", "round", "score", "guild", "auth", "club"}
	case "discord":
		// Discord creates its own internal stream.
		// It will subscribe to backend streams (user, guild, auth, etc) which backend creates.
		// We avoid creating shared streams here to prevent ambiguity and ensure Backend owns the domain streams.
		streams = []string{"discord"}
	default:
		ctxLogger.Error("Failed to create streams for app", "error", "unknown app type")
		return fmt.Errorf("unknown app type: %s", appType)
	}

	ctxLogger = ctxLogger.With(attr.Any("streams", streams))
	ctxLogger.Info("Streams to create")

	for _, stream := range streams {
		streamAttrs := ctxLogger.With(attr.String("current_stream", stream))
		streamAttrs.Debug("Creating stream")

		if err := eb.CreateStream(ctx, stream); err != nil {
			streamAttrs.Error("Failed to create stream", "error", err)
			return fmt.Errorf("failed to create stream %s: %w", stream, err)
		}
	}

	ctxLogger.Info("All streams created successfully for app type")
	return nil
}

// GetNATSConnection returns the underlying NATS connection.
func (eb *eventBus) GetNATSConnection() *nc.Conn {
	return eb.natsConn
}

// GetJetStream returns the JetStream context.
func (eb *eventBus) GetJetStream() jetstream.JetStream {
	return eb.js
}
