package eventbus

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Black-And-White-Club/frolf-bot-shared/observability/attr"
	lokifrolfbot "github.com/Black-And-White-Club/frolf-bot-shared/observability/loki"
	eventbusmetrics "github.com/Black-And-White-Club/frolf-bot-shared/observability/prometheus/eventbus"
	tempofrolfbot "github.com/Black-And-White-Club/frolf-bot-shared/observability/tempo"
	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-nats/v2/pkg/nats"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	nc "github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type eventBus struct {
	publisher      message.Publisher
	subscriber     message.Subscriber
	js             jetstream.JetStream
	natsConn       *nc.Conn
	logger         lokifrolfbot.Logger
	createdStreams map[string]bool
	streamMutex    sync.Mutex
	marshaler      nats.Marshaler
	metrics        eventbusmetrics.EventBusMetrics
	tracer         tempofrolfbot.Tracer
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
	ProcessDelayedMessages(ctx context.Context, roundID sharedtypes.RoundID, scheduledTime time.Time)
	CancelScheduledMessage(ctx context.Context, roundID sharedtypes.RoundID) error
	ScheduleRoundProcessing(ctx context.Context, roundID sharedtypes.RoundID, scheduledTime time.Time)
	RecoverScheduledRounds(ctx context.Context)

	GetNATSConnection() *nc.Conn
	GetJetStream() jetstream.JetStream
	GetHealthCheckers() []HealthChecker
}

// HealthChecker interface for components that can be health-checked
type HealthChecker interface {
	Check(ctx context.Context) error
	Name() string
}

// NewEventBus creates a new JetStream-backed EventBus.
func NewEventBus(ctx context.Context, natsURL string, logger lokifrolfbot.Logger, appType string, metrics eventbusmetrics.EventBusMetrics, tracer tempofrolfbot.Tracer) (EventBus, error) {

	logAttrs := []attr.LogAttr{
		attr.String("operation", "new_event_bus"),
		attr.String("nats_url", natsURL),
		attr.String("app_type", appType),
	}

	logger.Info("Creating new EventBus", logAttrs...)

	natsConn, err := nc.Connect(natsURL,
		nc.RetryOnFailedConnect(true),
		nc.Timeout(30*time.Second),
		nc.ReconnectWait(1*time.Second),
		nc.MaxReconnects(-1),
	)
	if err != nil {
		errorAttrs := append(logAttrs, attr.Error(err))
		logger.Error("Failed to connect to NATS", errorAttrs...)
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	js, err := jetstream.New(natsConn)
	if err != nil {
		natsConn.Close()
		errorAttrs := append(logAttrs, attr.Error(err))
		logger.Error("Failed to initialize JetStream", errorAttrs...)
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
				AutoProvision: false, // Now controlled manually by appType
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
		errorAttrs := append(logAttrs, attr.Error(err))
		logger.Error("Failed to create Watermill publisher", errorAttrs...)
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
		errorAttrs := append(logAttrs, attr.Error(err))
		logger.Error("Failed to create Watermill subscriber", errorAttrs...)
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

	// Only create the necessary streams for the given app type
	if err := eventBus.createStreamsForApp(ctx, appType); err != nil {
		natsConn.Close()
		errorAttrs := append(logAttrs, attr.Error(err))
		logger.Error("Failed to create streams for app", errorAttrs...)
		return nil, fmt.Errorf("failed to create streams: %w", err)
	}

	logger.Info("EventBus created successfully", logAttrs...)
	return eventBus, nil
}

func (eb *eventBus) ScheduleRoundProcessing(ctx context.Context, roundID sharedtypes.RoundID, scheduledTime time.Time) {
	logAttrs := []attr.LogAttr{
		attr.String("operation", "schedule_round_processing"),
		attr.Int("round_id", int(roundID)),
		attr.Time("scheduled_time", scheduledTime),
	}

	eb.logger.Info("Scheduling round processing", logAttrs...)
	go eb.ProcessDelayedMessages(ctx, roundID, scheduledTime)
}

// Publish publishes a message using Watermill's automatic marshaling.
func (eb *eventBus) Publish(topic string, messages ...*message.Message) error {
	logAttrs := []attr.LogAttr{
		attr.String("operation", "publish"),
		attr.String("topic", topic),
		attr.Int("message_count", len(messages)),
	}

	// Ensure `marshaler` is set
	if eb.marshaler == nil {
		errorAttrs := append(logAttrs, attr.String("error", "marshaler not set"))
		eb.logger.Error("Failed to publish message", errorAttrs...)
		return fmt.Errorf("eventBus marshaler is not set")
	}

	eb.logger.Debug("Publishing messages", logAttrs...)

	// Record metrics for message publishing
	eb.metrics.RecordMessagePublish(topic)

	// Inject trace context into messages
	for _, msg := range messages {
		eb.tracer.InjectTraceContext(msg.Context(), msg) // Use the new method
	}

	// Let Watermill's `NATSMarshaler` handle the serialization
	if err := eb.publisher.Publish(topic, messages...); err != nil {
		errorAttrs := append(logAttrs, attr.Error(err))
		eb.logger.Error("Failed to publish message", errorAttrs...)
		eb.metrics.RecordMessagePublishError(topic) // Record error metric
		return fmt.Errorf("failed to publish message to topic %s: %w", topic, err)
	}

	// Log successful publish
	for _, msg := range messages {
		correlationID := middleware.MessageCorrelationID(msg)
		msgAttrs := append(logAttrs,
			attr.String("message_id", msg.UUID),
			attr.String("correlation_id", correlationID),
		)
		eb.logger.Info("Message published successfully", msgAttrs...)
	}

	return nil
}

// Subscribe subscribes to a topic. Handles unmarshaling.
func (eb *eventBus) Subscribe(ctx context.Context, topic string) (<-chan *message.Message, error) {
	logAttrs := []attr.LogAttr{
		attr.String("operation", "subscribe"),
		attr.String("topic", topic),
	}

	eb.logger.Info("Entering Subscribe", logAttrs...)

	// Record metrics for message subscription
	eb.metrics.RecordMessageSubscribe(topic)

	var appType string
	if strings.HasPrefix(topic, "discord.") {
		appType = "discord"
	} else {
		appType = "backend"
	}
	logAttrs = append(logAttrs, attr.String("app_type", appType))

	consumerName := fmt.Sprintf("%s-consumer-%s", appType, sanitizeForNATS(topic))
	logAttrs = append(logAttrs, attr.String("consumer_name", consumerName))

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
		errorAttrs := append(logAttrs, attr.String("error", "unknown topic prefix"))
		eb.logger.Error("Failed to subscribe to topic", errorAttrs...)
		return nil, fmt.Errorf("unknown topic: %s", topic)
	}
	logAttrs = append(logAttrs, attr.String("stream", streamName))

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
			eb.logger.Info("Consumer not found, creating new one", logAttrs...)
			cons, err = eb.js.CreateConsumer(ctx, streamName, consumerConfig)
			if err != nil {
				errorAttrs := append(logAttrs, attr.Error(err))
				eb.logger.Error("Failed to create consumer", errorAttrs...)
				return nil, fmt.Errorf("failed to create consumer: %w", err)
			}
		} else {
			errorAttrs := append(logAttrs, attr.Error(err))
			eb.logger.Error("Unexpected error fetching consumer", errorAttrs...)
			return nil, err
		}
	}

	consInfoAttrs := append(logAttrs,
		attr.String("consumer_durable_name", cons.CachedInfo().Config.Durable),
		attr.String("consumer_filter_subject", cons.CachedInfo().Config.FilterSubject),
	)
	eb.logger.Info("Consumer created", consInfoAttrs...)

	messages := make(chan *message.Message)
	sub, err := cons.Messages()
	if err != nil {
		errorAttrs := append(logAttrs, attr.Error(err))
		eb.logger.Error("Failed to start consumer", errorAttrs...)
		return nil, fmt.Errorf("failed to start consumer for topic %s: %w", topic, err)
	}

	go func() {
		routineAttrs := append(logAttrs, attr.String("goroutine", "consumer_handler"))
		eb.logger.Info("Starting consumer goroutine", routineAttrs...)

		defer func() {
			eb.logger.Info("Closing consumer goroutine", routineAttrs...)
			close(messages)
			sub.Stop()
		}()

		for {
			jetStreamMsg, err := sub.Next()
			if err != nil {
				if errors.Is(err, jetstream.ErrMsgIteratorClosed) {
					iteratorAttrs := append(routineAttrs, attr.Error(err))
					eb.logger.Info("Consumer iterator closed", iteratorAttrs...)
				} else {
					errorAttrs := append(routineAttrs, attr.Error(err))
					eb.logger.Error("Error receiving message", errorAttrs...)
				}
				return
			}

			wmMsg, err := eb.toWatermillMessage(ctx, jetStreamMsg)
			if err != nil {
				errorAttrs := append(routineAttrs, attr.Error(err))
				eb.logger.Error("Failed to convert message", errorAttrs...)
				continue
			}

			// Inject trace context into the message
			eb.tracer.InjectTraceContext(ctx, wmMsg)

			msgAttrs := append(routineAttrs,
				attr.String("message_id", wmMsg.UUID),
				attr.Any("metadata", wmMsg.Metadata),
			)
			eb.logger.Info("Received message", msgAttrs...)

			select {
			case messages <- wmMsg:
				select {
				case <-wmMsg.Acked():
					if err := jetStreamMsg.Ack(); err != nil {
						ackAttrs := append(msgAttrs, attr.Error(err))
						eb.logger.Error("Failed to ack message", ackAttrs...)
					} else {
						eb.logger.Info("Message acknowledged", msgAttrs...)
					}
				case <-wmMsg.Nacked():
					if err := jetStreamMsg.Nak(); err != nil {
						nackAttrs := append(msgAttrs, attr.Error(err))
						eb.logger.Error("Failed to nack message", nackAttrs...)
					} else {
						eb.logger.Info("Message nacked", msgAttrs...)
					}
				case <-ctx.Done():
					eb.logger.Warn("Context cancelled (ack/nack)", msgAttrs...)
					if err := jetStreamMsg.Term(); err != nil {
						termAttrs := append(msgAttrs, attr.Error(err))
						eb.logger.Error("Failed to term message", termAttrs...)
					}
					return
				}
			case <-ctx.Done():
				ctxAttrs := append(msgAttrs, attr.Error(ctx.Err()))
				eb.logger.Warn("Context cancelled (sending)", ctxAttrs...)
				if err := jetStreamMsg.Term(); err != nil {
					termAttrs := append(ctxAttrs, attr.Error(err))
					eb.logger.Error("Failed to term message", termAttrs...)
				}
				return
			}
		}
	}()

	eb.logger.Info("Subscribed to topic successfully", logAttrs...)
	return messages, nil
}

func (eb *eventBus) toWatermillMessage(ctx context.Context, jetStreamMsg jetstream.Msg) (*message.Message, error) {
	msgID := string(jetStreamMsg.Headers().Get("Nats-Msg-Id"))
	watermillMsg := message.NewMessage(msgID, jetStreamMsg.Data())

	logAttrs := []attr.LogAttr{
		attr.String("operation", "to_watermill_message"),
		attr.String("message_id", watermillMsg.UUID),
	}

	meta, err := jetStreamMsg.Metadata()
	if err == nil {
		// Add metadata to log attributes
		logAttrs = append(logAttrs,
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
		logAttrs = append(logAttrs, attr.Error(err))
		eb.logger.Warn("Failed to get message metadata", logAttrs...)
	}

	// Copy headers to metadata
	for k, v := range jetStreamMsg.Headers() {
		watermillMsg.Metadata.Set(k, v[0])
	}

	// Ensure correlation ID
	if watermillMsg.Metadata.Get(middleware.CorrelationIDMetadataKey) == "" {
		correlationID := watermill.NewUUID()
		watermillMsg.Metadata.Set(middleware.CorrelationIDMetadataKey, correlationID)
		logAttrs = append(logAttrs, attr.String("correlation_id", correlationID))
		eb.logger.Debug("Generated new correlation ID", logAttrs...)
	} else {
		logAttrs = append(logAttrs, attr.String("correlation_id",
			watermillMsg.Metadata.Get(middleware.CorrelationIDMetadataKey)))
	}

	eb.logger.Debug("Created Watermill message from JetStream message", logAttrs...)

	watermillMsg.SetContext(ctx)
	return watermillMsg, nil
}

func sanitizeForNATS(s string) string {
	reg := regexp.MustCompile(`[^a-zA-Z0-9-]+`) // Allow only alphanumeric and dashes
	return reg.ReplaceAllString(strings.ReplaceAll(s, ".", "-"), "")
}

func (eb *eventBus) Close() error {
	logAttrs := []attr.LogAttr{
		attr.String("operation", "close"),
	}

	eb.logger.Info("Closing EventBus", logAttrs...)

	if eb.publisher != nil {
		if err := eb.publisher.Close(); err != nil {
			errorAttrs := append(logAttrs,
				attr.String("component", "publisher"),
				attr.Error(err),
			)
			eb.logger.Error("Error closing publisher", errorAttrs...)
		} else {
			eb.logger.Info("Publisher closed successfully",
				append(logAttrs, attr.String("component", "publisher"))...)
		}
	}

	if eb.subscriber != nil {
		if err := eb.subscriber.Close(); err != nil {
			errorAttrs := append(logAttrs,
				attr.String("component", "subscriber"),
				attr.Error(err),
			)
			eb.logger.Error("Error closing subscriber", errorAttrs...)
		} else {
			eb.logger.Info("Subscriber closed successfully",
				append(logAttrs, attr.String("component", "subscriber"))...)
		}
	}

	if eb.natsConn != nil {
		eb.natsConn.Close()
		eb.logger.Info("NATS connection closed",
			append(logAttrs, attr.String("component", "nats_connection"))...)
	}

	eb.logger.Info("EventBus closed successfully", logAttrs...)
	return nil
}

// CreateStream creates or updates a single JetStream stream (helper function).
func (eb *eventBus) CreateStream(ctx context.Context, streamName string) error {
	logAttrs := []attr.LogAttr{
		attr.String("operation", "create_stream"),
		attr.String("stream_name", streamName),
	}

	eb.logger.Info("Creating/Updating stream", logAttrs...)
	eb.streamMutex.Lock() // Protect against concurrent stream creation
	defer eb.streamMutex.Unlock()

	if eb.createdStreams[streamName] {
		eb.logger.Info("Stream already created in this process", logAttrs...)
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
		errorAttrs := append(logAttrs, attr.String("error", "unknown stream name"))
		eb.logger.Error("Failed to create stream", errorAttrs...)
		return fmt.Errorf("unknown stream name: %s", streamName)
	}

	logAttrs = append(logAttrs, attr.Any("subjects", subjects))

	// Define default stream configuration
	streamCfg := jetstream.StreamConfig{
		Name:      streamName,
		Subjects:  subjects,
		Retention: jetstream.LimitsPolicy, // Retain based on limits (age, size)
		// Storage:   jetstream.FileStorage,  // Store on disk
	}

	// Customize for delayed messages stream
	switch streamName {
	case DelayedMessagesStream:
		streamCfg.MaxAge = 24 * time.Hour            // Auto-delete after 24 hours
		streamCfg.MaxMsgs = -1                       // Unlimited messages (rely on MaxAge)
		streamCfg.Retention = jetstream.LimitsPolicy // Use limits policy
		logAttrs = append(logAttrs,
			attr.Duration("max_age", streamCfg.MaxAge),
			attr.Int64("max_msgs", streamCfg.MaxMsgs))
	default:
		streamCfg.Duplicates = 5 * time.Minute // 5-minute deduplication window
		logAttrs = append(logAttrs, attr.Duration("duplicates_window", streamCfg.Duplicates))
	}

	// Create or update the stream (idempotent)
	stream, err := eb.js.CreateOrUpdateStream(ctx, streamCfg)
	if err != nil {
		errorAttrs := append(logAttrs, attr.Error(err))
		eb.logger.Error("Failed to create or update stream", errorAttrs...)
		return fmt.Errorf("failed to create or update stream %s: %w", streamName, err)
	}

	successAttrs := append(logAttrs, attr.String("stream_name", stream.CachedInfo().Config.Name))
	eb.logger.Info("Stream created or updated successfully", successAttrs...)
	eb.createdStreams[streamName] = true // Mark as created
	return nil
}

func (eb *eventBus) createStreamsForApp(ctx context.Context, appType string) error {
	logAttrs := []attr.LogAttr{
		attr.String("operation", "create_streams_for_app"),
		attr.String("app_type", appType),
	}

	eb.logger.Info("Creating streams for application type", logAttrs...)

	var streams []string
	switch appType {
	case "backend":
		streams = []string{"user", "leaderboard", "round", "score", "delayed"}
	case "discord":
		streams = []string{"discord"}
	default:
		errorAttrs := append(logAttrs, attr.String("error", "unknown app type"))
		eb.logger.Error("Failed to create streams for app", errorAttrs...)
		return fmt.Errorf("unknown app type: %s", appType)
	}

	logAttrs = append(logAttrs, attr.Any("streams", streams))
	eb.logger.Info("Streams to create", logAttrs...)

	for _, stream := range streams {
		streamAttrs := append(logAttrs, attr.String("current_stream", stream))
		eb.logger.Debug("Creating stream", streamAttrs...)

		if err := eb.CreateStream(ctx, stream); err != nil {
			errorAttrs := append(streamAttrs, attr.Error(err))
			eb.logger.Error("Failed to create stream", errorAttrs...)
			return fmt.Errorf("failed to create stream %s: %w", stream, err)
		}
	}

	eb.logger.Info("All streams created successfully for app type", logAttrs...)
	return nil
}

// GetNATSConnection returns the underlying NATS connection
func (eb *eventBus) GetNATSConnection() *nc.Conn {
	return eb.natsConn // Assuming conn is your NATS connection
}

// GetJetStream returns the JetStream context
func (eb *eventBus) GetJetStream() jetstream.JetStream {
	return eb.js // Assuming js is your JetStream context
}
