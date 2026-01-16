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
	appType        string
	publisher      message.Publisher
	subscriber     message.Subscriber
	js             jetstream.JetStream
	natsConn       *nc.Conn
	logger         *slog.Logger
	createdStreams map[string]bool
	streamMutex    sync.Mutex
	marshaler      nats.Marshaler
	metrics        eventbusmetrics.EventBusMetrics
	tracer         trace.Tracer
}

// EventBus interface - REMOVED all delayed message methods
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

	subscriber, err := nats.NewSubscriber(
		nats.SubscriberConfig{
			URL:            natsURL,
			AckWaitTimeout: 30 * time.Second,
			Unmarshaler:    marshaller,
			NatsOptions: []nc.Option{
				nc.RetryOnFailedConnect(true),
				nc.Timeout(30 * time.Second),
				nc.ReconnectWait(1 * time.Second),
				nc.MaxReconnects(-1),
			},
			CloseTimeout:     5 * time.Second,
			SubscribeTimeout: 5 * time.Second,
		},
		watermillLogger,
	)
	if err != nil {
		natsConn.Close()
		ctxLogger.ErrorContext(ctx, "Failed to create Watermill subscriber", "error", err)
		return nil, fmt.Errorf("failed to create Watermill subscriber: %w", err)
	}

	eventBus := &eventBus{
		appType:        appType,
		publisher:      publisher,
		subscriber:     subscriber,
		js:             js,
		natsConn:       natsConn,
		logger:         logger,
		createdStreams: make(map[string]bool),
		streamMutex:    sync.Mutex{},
		marshaler:      marshaller,
		metrics:        metrics,
		tracer:         tracer,
	}

	if err := eventBus.createStreamsForApp(ctx, appType); err != nil {
		natsConn.Close()
		ctxLogger.ErrorContext(ctx, "Failed to create streams for app", "error", err)
		return nil, fmt.Errorf("failed to create streams for app: %w", err)
	}

	ctxLogger.InfoContext(ctx, "EventBus created successfully")
	return eventBus, nil
}

// Publish publishes a message using Watermill's automatic marshaling with idempotency.
func (eb *eventBus) Publish(topic string, messages ...*message.Message) error {
	ctxLogger := eb.logger.With(
		"operation", "publish",
		"topic", topic,
		"message_count", len(messages),
	)

	if eb.marshaler == nil {
		ctxLogger.Error("Failed to publish message", "error", "marshaler not set")
		return fmt.Errorf("eventBus marshaler is not set")
	}

	// Guard: prevent non-discord apps from publishing discord.* topics per boundary rules
	if strings.HasPrefix(topic, "discord.") && eb.appType != "discord" {
		ctxLogger.Error("Publish forbidden: app not allowed to publish discord topics", "app_type", eb.appType, "topic", topic)
		return fmt.Errorf("publishing to discord topics is forbidden for app type %q", eb.appType)
	}

	ctxLogger.Debug("Publishing messages")

	// If router passed empty topic, try to derive the subject from the message metadata.
	if topic == "" {
		if len(messages) == 0 {
			return errors.New("no topic provided and no messages to derive topic from")
		}
		// Prefer metadata keys (case-sensitive) set by helpers.CreateResultMessage
		// Try common keys in order: "topic", "Topic", "event_name", "topic_hint"
		meta := messages[0].Metadata
		derived := meta.Get("topic")
		if derived == "" {
			derived = meta.Get("Topic")
		}
		if derived == "" {
			derived = meta.Get("event_name")
		}
		if derived == "" {
			derived = meta.Get("topic_hint")
		}
		if derived == "" {
			return fmt.Errorf("no publish topic provided and message metadata missing topic")
		}
		topic = derived
	}

	if len(messages) > 0 && eb.metrics != nil {
		eb.metrics.RecordMessagePublish(messages[0].Context(), topic)
	}

	// Ensure each message has a unique ID for deduplication
	for _, msg := range messages {
		// Set Nats-Msg-Id header for deduplication (JetStream uses this)
		if msg.Metadata.Get("Nats-Msg-Id") == "" {
			if idempotencyKey := msg.Metadata.Get("idempotency_key"); idempotencyKey != "" {
				msg.Metadata.Set("Nats-Msg-Id", dedupeMsgID(idempotencyKey, topic))
			} else {
				msg.Metadata.Set("Nats-Msg-Id", msg.UUID)
			}
		}

		ctxLogger.Debug("Attempting to publish message",
			attr.String("message_uuid", msg.UUID),
			attr.String("nats_msg_id", msg.Metadata.Get("Nats-Msg-Id")),
			attr.String("topic_name", topic),
		)
	}

	// Publish with retry logic for network issues
	maxRetries := 3
	for attempt := 1; attempt <= maxRetries; attempt++ {
		if err := eb.publisher.Publish(topic, messages...); err != nil {
			ctxLogger.Warn("Publish attempt failed",
				"attempt", attempt,
				"max_retries", maxRetries,
				"error", err,
			)

			// Check if it's a network/timeout error that should be retried
			if attempt < maxRetries && eb.isRetryableError(err) {
				backoffDuration := time.Duration(attempt) * 100 * time.Millisecond
				ctxLogger.Info("Retrying publish after backoff", "backoff", backoffDuration)
				time.Sleep(backoffDuration)
				continue
			}

			ctxLogger.Error("Failed to publish message after retries", "error", err)
			if len(messages) > 0 && eb.metrics != nil {
				eb.metrics.RecordMessagePublishError(messages[0].Context(), topic)
			}

			return fmt.Errorf("failed to publish message to topic %s after %d attempts: %w", topic, maxRetries, err)
		}

		// Success
		break
	}

	for _, msg := range messages {
		correlationID := middleware.MessageCorrelationID(msg)
		msgContext := msg.Context()

		msgLogger := ctxLogger.With(
			"message_id", msg.UUID,
			"nats_msg_id", msg.Metadata.Get("Nats-Msg-Id"),
			"correlation_id", correlationID,
		)

		msgLogger.InfoContext(msgContext, "Message published successfully")
	}

	return nil
}

func dedupeMsgID(idempotencyKey, topic string) string {
	// Include topic to avoid collisions across subjects within the same stream.
	seed := topic + "|" + idempotencyKey
	sum := sha256.Sum256([]byte(seed))
	return "idem-" + hex.EncodeToString(sum[:])
}

// Subscribe subscribes to a topic. Handles unmarshaling with retry logic.
func (eb *eventBus) Subscribe(ctx context.Context, topic string) (<-chan *message.Message, error) {
	ctxLogger := eb.logger.With(
		"operation", "subscribe",
		"topic", topic,
	)

	if eb.metrics != nil {
		eb.metrics.RecordMessageSubscribe(ctx, topic)
	}

	// Boundary guard
	if strings.HasPrefix(topic, "discord.") && eb.appType != "discord" {
		return nil, fmt.Errorf("subscription to discord topics forbidden for app %q", eb.appType)
	}

	var appType string
	if strings.HasPrefix(topic, "discord.") {
		appType = "discord"
	} else {
		appType = "backend"
	}

	consumerName := fmt.Sprintf("%s-consumer-%s", appType, sanitizeForNATS(topic))
	if v := ctx.Value("eventbus_consumer_name"); v != nil {
		if s, ok := v.(string); ok && s != "" {
			consumerName = s
		}
	}

	// Stream resolution
	var streamName string
	switch {
	case strings.HasPrefix(topic, "user."):
		streamName = "user"
	case strings.HasPrefix(topic, "leaderboard."):
		streamName = "leaderboard"
	case strings.HasPrefix(topic, "round."):
		streamName = "round"
	case strings.HasPrefix(topic, "score."):
		streamName = "score"
	case strings.HasPrefix(topic, "guild."):
		streamName = "guild"
	case strings.HasPrefix(topic, "discord."):
		streamName = "discord"
	default:
		return nil, fmt.Errorf("unknown topic: %s", topic)
	}

	if err := eb.CreateStream(ctx, streamName); err != nil {
		return nil, err
	}

	ackWait := 30 * time.Second

	consumerConfig := jetstream.ConsumerConfig{
		Durable:       consumerName,
		FilterSubject: topic,

		AckPolicy:     jetstream.AckExplicitPolicy,
		AckWait:       ackWait,
		MaxDeliver:    3,
		BackOff:       []time.Duration{1 * time.Second, 5 * time.Second},
		MaxAckPending: 100,

		DeliverPolicy: jetstream.DeliverNewPolicy,
		ReplayPolicy:  jetstream.ReplayInstantPolicy,
	}

	// Declarative sync
	cons, err := eb.js.CreateOrUpdateConsumer(ctx, streamName, consumerConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to sync consumer: %w", err)
	}

	sub, err := cons.Messages()
	if err != nil {
		return nil, err
	}

	out := make(chan *message.Message, 128)

	// Bounded ACK goroutines
	ackSem := make(chan struct{}, 128)

	go func() {
		defer func() {
			sub.Stop()
			close(out)
		}()

		for {
			jsMsg, err := sub.Next()
			if err != nil {
				if !errors.Is(err, jetstream.ErrMsgIteratorClosed) {
					ctxLogger.Error("consumer error", "error", err)
				}
				return
			}

			wmMsg, err := eb.toWatermillMessage(ctx, jsMsg)
			if err != nil {
				_ = jsMsg.Term()
				continue
			}

			select {
			case out <- wmMsg:
				ackSem <- struct{}{}

				processTimeout := time.Duration(float64(ackWait) * 0.8)
				pCtx, cancel := context.WithTimeout(ctx, processTimeout)

				go func() {
					defer cancel()
					defer func() { <-ackSem }()

					select {
					case <-wmMsg.Acked():
						_ = jsMsg.Ack()
					case <-wmMsg.Nacked():
						_ = jsMsg.Nak()
					case <-pCtx.Done():
						// Let AckWait expire → redelivery
					case <-ctx.Done():
						// graceful shutdown → redelivery
					}
				}()

			case <-ctx.Done():
				return
			}
		}
	}()

	return out, nil
}

func (eb *eventBus) SubscribeForTest(ctx context.Context, topic string) (<-chan *message.Message, error) {
	var streamName string
	switch {
	case strings.HasPrefix(topic, "user."):
		streamName = "user"
	case strings.HasPrefix(topic, "leaderboard."):
		streamName = "leaderboard"
	case strings.HasPrefix(topic, "round."):
		streamName = "round"
	case strings.HasPrefix(topic, "score."):
		streamName = "score"
	case strings.HasPrefix(topic, "guild."):
		streamName = "guild"
	case strings.HasPrefix(topic, "discord."):
		streamName = "discord"
	default:
		return nil, fmt.Errorf("unknown topic: %s", topic)
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
func (eb *eventBus) Close() error {
	ctxLogger := eb.logger.With(attr.String("operation", "close"))

	ctxLogger.Info("Closing EventBus")

	if eb.publisher != nil {
		if err := eb.publisher.Close(); err != nil {
			ctxLogger.Error("Error closing publisher", "component", "publisher", "error", err)
		} else {
			ctxLogger.Info("Publisher closed successfully", "component", "publisher")
		}
	}

	if eb.subscriber != nil {
		if err := eb.subscriber.Close(); err != nil {
			ctxLogger.Error("Error closing subscriber", "component", "subscriber", "error", err)
		} else {
			ctxLogger.Info("Subscriber closed successfully", "component", "subscriber")
		}
	}

	if eb.natsConn != nil {
		eb.natsConn.Close()
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
		streams = []string{"user", "leaderboard", "round", "score", "guild"}
	case "discord":
		// Discord creates its own internal stream and the shared guild stream
		// It will subscribe to backend streams (which backend creates)
		streams = []string{"discord", "guild"}
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
