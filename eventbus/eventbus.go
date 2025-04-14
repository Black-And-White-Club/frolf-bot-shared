package eventbus

import (
	"context"
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
	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-nats/v2/pkg/nats"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	nc "github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel/trace"
)

type eventBus struct {
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

const (
	DelayedMessagesStream  = "delayed"
	DelayedMessagesSubject = "delayed.message"
	DiscordEventsSubject   = "discord.round.event"
)

// EventBus interface
type EventBus interface {
	Publish(topic string, messages ...*message.Message) error
	Subscribe(ctx context.Context, topic string) (<-chan *message.Message, error)
	Close() error
	ProcessDelayedMessages(ctx context.Context, roundID sharedtypes.RoundID, scheduledTime sharedtypes.StartTime) error
	CancelScheduledMessage(ctx context.Context, roundID sharedtypes.RoundID) error
	RecoverScheduledRounds(ctx context.Context)
	ScheduleDelayedMessage(ctx context.Context, originalSubject string, roundID sharedtypes.RoundID, scheduledTime sharedtypes.StartTime, payload []byte) error
	GetNATSConnection() *nc.Conn
	GetJetStream() jetstream.JetStream
	GetHealthCheckers() []HealthChecker
}

// HealthChecker interface for components that can be health-checked
type HealthChecker interface {
	Check(ctx context.Context) error
	Name() string
}

func NewEventBus(ctx context.Context, natsURL string, logger *slog.Logger, appType string, metrics eventbusmetrics.EventBusMetrics, tracer trace.Tracer) (EventBus, error) {
	// Create a contextual logger with the base attributes
	ctxLogger := logger.With(
		"operation", "new_event_bus",
		"nats_url", natsURL,
		"app_type", appType,
	)

	// Log the creation of the EventBus
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

// Publish publishes a message using Watermill's automatic marshaling.
func (eb *eventBus) Publish(topic string, messages ...*message.Message) error {
	// Create a contextual logger with the base attributes
	ctxLogger := eb.logger.With(
		"operation", "publish",
		"topic", topic,
		"message_count", len(messages),
	)

	// Ensure `marshaler` is set
	if eb.marshaler == nil {
		ctxLogger.Error("Failed to publish message", "error", "marshaler not set")
		return fmt.Errorf("eventBus marshaler is not set")
	}

	ctxLogger.Debug("Publishing messages")

	// Record metrics for message publishing
	if len(messages) > 0 {
		eb.metrics.RecordMessagePublish(messages[0].Context(), topic)
	}

	// Let Watermill's `NATSMarshaler` handle the serialization
	if err := eb.publisher.Publish(topic, messages...); err != nil {
		ctxLogger.Error("Failed to publish message", "error", err)
		if len(messages) > 0 {
			eb.metrics.RecordMessagePublishError(messages[0].Context(), topic)
		}
		return fmt.Errorf("failed to publish message to topic %s: %w", topic, err)
	}

	// Log successful publish using the context from each message
	for _, msg := range messages {
		correlationID := middleware.MessageCorrelationID(msg)
		msgContext := msg.Context() // Get the context from the message

		// Create a contextual logger for the message
		msgLogger := ctxLogger.With(
			"message_id", msg.UUID,
			"correlation_id", correlationID,
		)

		// Use InfoContext to log with the message's context
		msgLogger.InfoContext(msgContext, "Message published successfully")
	}

	return nil
}

// Subscribe subscribes to a topic. Handles unmarshaling.
func (eb *eventBus) Subscribe(ctx context.Context, topic string) (<-chan *message.Message, error) {
	// Create a contextual logger with the base attributes
	ctxLogger := eb.logger.With(
		"operation", "subscribe",
		"topic", topic,
	)

	ctxLogger.Info("Entering Subscribe")

	// Record metrics for message subscription using the context
	eb.metrics.RecordMessageSubscribe(ctx, topic)

	var appType string
	if strings.HasPrefix(topic, "discord.") {
		appType = "discord"
	} else {
		appType = "backend"
	}
	ctxLogger = ctxLogger.With("app_type", appType)

	consumerName := fmt.Sprintf("%s-consumer-%s", appType, sanitizeForNATS(topic))
	ctxLogger = ctxLogger.With("consumer_name", consumerName)

	// Determine stream name based on topic
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
	case strings.HasPrefix(topic, "discord."):
		streamName = "discord"
	case strings.HasPrefix(topic, "delayed."):
		streamName = "delayed"
	default:
		ctxLogger.Error("Failed to subscribe to topic", "error", "unknown topic prefix")
		return nil, fmt.Errorf("unknown topic: %s", topic)
	}
	ctxLogger = ctxLogger.With("stream", streamName)

	// Create or get consumer
	consumerConfig := jetstream.ConsumerConfig{
		Durable:       consumerName,
		FilterSubject: topic,
		AckPolicy:     jetstream.AckExplicitPolicy,
		MaxAckPending: 2048,
		DeliverPolicy: jetstream.DeliverAllPolicy,
	}

	cons, err := eb.js.Consumer(ctx, streamName, consumerName)
	if err != nil {
		if errors.Is(err, jetstream.ErrConsumerNotFound) {
			ctxLogger.Info("Consumer not found, creating new one")
			cons, err = eb.js.CreateConsumer(ctx, streamName, consumerConfig)
			if err != nil {
				ctxLogger.Error("Failed to create consumer", "error", err)
				return nil, fmt.Errorf("failed to create consumer: %w", err)
			}
		} else {
			ctxLogger.Error("Unexpected error fetching consumer", "error", err)
			return nil, err
		}
	}

	consInfoAttrs := ctxLogger.With(
		"consumer_durable_name", cons.CachedInfo().Config.Durable,
		"consumer_filter_subject", cons.CachedInfo().Config.FilterSubject,
	)
	consInfoAttrs.Info("Consumer created")

	messages := make(chan *message.Message)
	sub, err := cons.Messages()
	if err != nil {
		ctxLogger.Error("Failed to start consumer", "error", err)
		return nil, fmt.Errorf("failed to start consumer for topic %s: %w", topic, err)
	}

	go func() {
		routineLogger := ctxLogger.With("goroutine", "consumer_handler")
		routineLogger.Info("Starting consumer goroutine")

		defer func() {
			routineLogger.Info("Closing consumer goroutine")
			close(messages)
			sub.Stop()
		}()

		for {
			jetStreamMsg, err := sub.Next()
			if err != nil {
				if errors.Is(err, jetstream.ErrMsgIteratorClosed) {
					routineLogger.Info("Consumer iterator closed", "error", err)
				} else {
					routineLogger.Error("Error receiving message", "error", err)
				}
				return
			}

			wmMsg, err := eb.toWatermillMessage(ctx, jetStreamMsg)
			if err != nil {
				routineLogger.Error("Failed to convert message", "error", err)
				continue
			}

			msgAttrs := routineLogger.With(
				"message_id", wmMsg.UUID,
				"metadata", wmMsg.Metadata,
			)
			msgAttrs.Info("Received message")

			select {
			case messages <- wmMsg:
				select {
				case <-wmMsg.Acked():
					if err := jetStreamMsg.Ack(); err != nil {
						msgAttrs.Error("Failed to ack message", "error", err)
					} else {
						msgAttrs.Info("Message acknowledged")
					}
				case <-wmMsg.Nacked():
					if err := jetStreamMsg.Nak(); err != nil {
						msgAttrs.Error("Failed to nack message", "error", err)
					} else {
						msgAttrs.Info("Message nacked")
					}
				case <-ctx.Done():
					msgAttrs.Warn("Context cancelled (ack/nack)")
					if err := jetStreamMsg.Term(); err != nil {
						msgAttrs.Error("Failed to term message", "error", err)
					}
					return
				}
			case <-ctx.Done():
				msgAttrs.Warn("Context cancelled (sending)", "error", ctx.Err())
				if err := jetStreamMsg.Term(); err != nil {
					msgAttrs.Error("Failed to term message", "error", err)
				}
				return
			}
		}
	}()

	ctxLogger.Info("Subscribed to topic successfully")
	return messages, nil
}

// toWatermillMessage converts a JetStream message to a Watermill message.
func (eb *eventBus) toWatermillMessage(ctx context.Context, jetStreamMsg jetstream.Msg) (*message.Message, error) {
	msgID := string(jetStreamMsg.Headers().Get("Nats-Msg-Id"))
	watermillMsg := message.NewMessage(msgID, jetStreamMsg.Data())

	// Create a contextual logger for this operation
	ctxLogger := eb.logger.With(
		attr.String("operation", "to_watermill_message"),
		attr.String("message_id", watermillMsg.UUID),
	)

	meta, err := jetStreamMsg.Metadata()
	if err == nil {
		// Add metadata to log attributes
		ctxLogger = ctxLogger.With(
			attr.String("stream", meta.Stream),
			attr.String("consumer", meta.Consumer),
			attr.Uint64("stream_seq", meta.Sequence.Stream),
			attr.Uint64("consumer_seq", meta.Sequence.Consumer),
			attr.Uint64("num_delivered", meta.NumDelivered),
		)

		// Set metadata in message
		watermillMsg.Metadata.Set("Stream", meta.Stream)
		watermillMsg.Metadata.Set("Consumer", meta.Consumer)
		watermillMsg.Metadata.Set("Delivered", strconv.FormatInt(int64(meta.NumDelivered), 10))
		watermillMsg.Metadata.Set("StreamSeq", strconv.FormatUint(meta.Sequence.Stream, 10))
		watermillMsg.Metadata.Set("ConsumerSeq", strconv.FormatUint(meta.Sequence.Consumer, 10))
		watermillMsg.Metadata.Set("Timestamp", meta.Timestamp.String())
	} else {
		ctxLogger.Warn("Failed to get message metadata", "error", err)
	}

	// Copy headers to metadata
	for k, v := range jetStreamMsg.Headers() {
		watermillMsg.Metadata.Set(k, v[0])
	}

	// Ensure correlation ID
	if watermillMsg.Metadata.Get(middleware.CorrelationIDMetadataKey) == "" {
		correlationID := watermill.NewUUID()
		watermillMsg.Metadata.Set(middleware.CorrelationIDMetadataKey, correlationID)
		ctxLogger = ctxLogger.With(attr.String("correlation_id", correlationID))
		ctxLogger.Debug("Generated new correlation ID")
	} else {
		ctxLogger = ctxLogger.With(attr.String("correlation_id",
			watermillMsg.Metadata.Get(middleware.CorrelationIDMetadataKey)))
	}

	ctxLogger.Debug("Created Watermill message from JetStream message")

	watermillMsg.SetContext(ctx)
	return watermillMsg, nil
}

// sanitizeForNATS sanitizes a string for use in NATS topics.
func sanitizeForNATS(s string) string {
	reg := regexp.MustCompile(`[^a-zA-Z0-9-]+`) // Allow only alphanumeric and dashes
	return reg.ReplaceAllString(strings.ReplaceAll(s, ".", "-"), "")
}

// Close closes the EventBus and its components.
func (eb *eventBus) Close() error {
	// Create a contextual logger for this operation
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
		eb.natsConn.Close() // Call Close() without checking for a return value
		ctxLogger.Info("NATS connection closed", "component", "nats_connection")
	}

	ctxLogger.Info("EventBus closed successfully")
	return nil
}

// CreateStream creates or updates a single JetStream stream (helper function).
func (eb *eventBus) CreateStream(ctx context.Context, streamName string) error {
	// Create a contextual logger for this operation
	ctxLogger := eb.logger.With(
		attr.String("operation", "create_stream"),
		attr.String("stream_name", streamName),
	)

	ctxLogger.Info("Creating/Updating stream")
	eb.streamMutex.Lock() // Protect against concurrent stream creation
	defer eb.streamMutex.Unlock()

	if eb.createdStreams[streamName] {
		ctxLogger.Info("Stream already created in this process")
		return nil // Already created, nothing to do
	}

	var subjects []string
	switch streamName {
	case "user":
		subjects = []string{"user.>"} // e.g., user.created, user.updated
	case "leaderboard":
		subjects = []string{"leaderboard.>"}
	case "round":
		subjects = []string{"round.>"}
	case "score":
		subjects = []string{"score.>"}
	case "discord":
		subjects = []string{"discord.>"}
	case DelayedMessagesStream:
		subjects = []string{"delayed.>"}
	default:
		ctxLogger.Error("Failed to create stream", "error", "unknown stream name")
		return fmt.Errorf("unknown stream name: %s", streamName)
	}

	ctxLogger = ctxLogger.With(attr.Any("subjects", subjects))

	// Define default stream configuration
	streamCfg := jetstream.StreamConfig{
		Name:      streamName,
		Subjects:  subjects,
		Retention: jetstream.LimitsPolicy, // Retain based on limits (age, size)
	}

	// Customize for delayed messages stream
	switch streamName {
	case DelayedMessagesStream:
		streamCfg.MaxAge = 24 * time.Hour            // Auto-delete after 24 hours
		streamCfg.MaxMsgs = -1                       // Unlimited messages (rely on MaxAge)
		streamCfg.Retention = jetstream.LimitsPolicy // Use limits policy
		ctxLogger = ctxLogger.With(
			attr.Duration("max_age", streamCfg.MaxAge),
			attr.Int64("max_msgs", streamCfg.MaxMsgs),
		)
	default:
		streamCfg.Duplicates = 5 * time.Minute // 5-minute deduplication window
		ctxLogger = ctxLogger.With(attr.Duration("duplicates_window", streamCfg.Duplicates))
	}

	// Create or update the stream (idempotent)
	stream, err := eb.js.CreateOrUpdateStream(ctx, streamCfg)
	if err != nil {
		ctxLogger.Error("Failed to create or update stream", "error", err)
		return fmt.Errorf("failed to create or update stream %s: %w", streamName, err)
	}

	ctxLogger.Info("Stream created or updated successfully", "stream_name", stream.CachedInfo().Config.Name)
	eb.createdStreams[streamName] = true // Mark as created
	return nil
}

// createStreamsForApp creates streams for a specific application type.
func (eb *eventBus) createStreamsForApp(ctx context.Context, appType string) error {
	// Create a contextual logger for this operation
	ctxLogger := eb.logger.With(
		attr.String("operation", "create_streams_for_app"),
		attr.String("app_type", appType),
	)

	ctxLogger.Info("Creating streams for application type")

	var streams []string
	switch appType {
	case "backend":
		streams = []string{"user", "leaderboard", "round", "score", "delayed"}
	case "discord":
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
	return eb.natsConn // Assuming conn is your N ATS connection
}

// GetJetStream returns the JetStream context.
func (eb *eventBus) GetJetStream() jetstream.JetStream {
	return eb.js // Assuming js is your JetStream context
}
