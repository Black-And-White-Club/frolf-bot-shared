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
	logger         *slog.Logger
	createdStreams map[string]bool
	streamMutex    sync.Mutex
	marshaler      nats.Marshaler
}

const (
	delayedMessagesStream  = "delayed"
	delayedMessagesSubject = "delayed.message"
	discordEventsSubject   = "discord.round.event"
)

// EventBus interface
type EventBus interface {
	Publish(topic string, messages ...*message.Message) error
	Subscribe(ctx context.Context, topic string) (<-chan *message.Message, error)
	Close() error
	ProcessDelayedMessages(ctx context.Context)
	CancelScheduledMessage(ctx context.Context, roundID string) error
}

// NewEventBus creates a new JetStream-backed EventBus.
func NewEventBus(ctx context.Context, natsURL string, logger *slog.Logger) (EventBus, error) {
	natsConn, err := nc.Connect(natsURL,
		nc.RetryOnFailedConnect(true),
		nc.Timeout(30*time.Second),
		nc.ReconnectWait(1*time.Second),
		nc.MaxReconnects(-1),
	)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to connect to NATS", slog.Any("error", err))
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	js, err := jetstream.New(natsConn)
	if err != nil {
		natsConn.Close()
		logger.ErrorContext(ctx, "Failed to initialize JetStream", slog.Any("error", err))
		return nil, fmt.Errorf("failed to initialize JetStream: %w", err)
	}

	watermillLogger := watermill.NewSlogLogger(logger)

	// Use nats.NATSMarshaler for automatic JSON marshaling/unmarshaling.
	marshaller := &nats.NATSMarshaler{}

	publisher, err := nats.NewPublisher(
		nats.PublisherConfig{
			URL:       natsURL,
			Marshaler: marshaller, // <--- Use the NATSMarshaler
			NatsOptions: []nc.Option{
				nc.RetryOnFailedConnect(true),
				nc.Timeout(30 * time.Second),
				nc.ReconnectWait(1 * time.Second),
				nc.MaxReconnects(-1),
			},
			JetStream: nats.JetStreamConfig{
				Disabled:      false,
				AutoProvision: true, // Consider setting to false in production
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
		logger.ErrorContext(ctx, "Failed to create Watermill publisher", slog.Any("error", err))
		return nil, fmt.Errorf("failed to create Watermill publisher: %w", err)
	}

	subscriber, err := nats.NewSubscriber(
		nats.SubscriberConfig{
			URL:            natsURL,
			AckWaitTimeout: 30 * time.Second,
			Unmarshaler:    marshaller, // <--- Use the NATSMarshaler
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
		logger.ErrorContext(ctx, "Failed to create Watermill subscriber", slog.Any("error", err))
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
	}

	if err := eventBus.createStreamsAndDLQs(ctx); err != nil {
		natsConn.Close()
		return nil, fmt.Errorf("failed to create streams and DLQs: %w", err)
	}

	go eventBus.ProcessDelayedMessages(ctx)

	return eventBus, nil
}

// Publish publishes a message using Watermill's automatic marshaling.
func (eb *eventBus) Publish(topic string, messages ...*message.Message) error {
	// Ensure `marshaler` is set
	ctx := messages[0].Context()
	if eb.marshaler == nil {
		return fmt.Errorf("eventBus marshaler is not set")
	}

	// Convert the payload into a Watermill message
	wmMsg := message.NewMessage(watermill.NewUUID(), nil)
	wmMsg.SetContext(ctx)
	wmMsg.Metadata.Set("topic", topic)

	// Attempt to retrieve Correlation ID from metadata
	correlationID := middleware.MessageCorrelationID(wmMsg)
	if correlationID == "" {
		correlationID = watermill.NewUUID() // Generate new one if missing
	}
	wmMsg.Metadata.Set(middleware.CorrelationIDMetadataKey, correlationID)

	// Let Watermill's `NATSMarshaler` handle the serialization
	if err := eb.publisher.Publish(topic, wmMsg); err != nil {
		return fmt.Errorf("failed to publish message to topic %s: %w", topic, err)
	}

	eb.logger.InfoContext(ctx, "Message published successfully",
		slog.String("topic", topic),
		slog.String("message_id", wmMsg.UUID),
		slog.String("correlation_id", correlationID),
	)
	return nil
}

// Subscribe subscribes to a topic.  Handles unmarshaling.
func (eb *eventBus) Subscribe(ctx context.Context, topic string) (<-chan *message.Message, error) {
	eb.logger.InfoContext(ctx, "Entering Subscribe", slog.String("topic", topic))

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
	case topic == delayedMessagesSubject:
		streamName = delayedMessagesStream //subscribe to the delayed message stream.
	default:
		return nil, fmt.Errorf("unknown topic: %s", topic)
	}

	consumerName := fmt.Sprintf("consumer-%s", sanitize(topic))

	consumerConfig := jetstream.ConsumerConfig{
		Durable:       consumerName,
		FilterSubject: topic,
		AckPolicy:     jetstream.AckExplicitPolicy,
		MaxAckPending: 2048,
		DeliverPolicy: jetstream.DeliverAllPolicy, // All available messages
	}

	if topic == delayedMessagesSubject {
		consumerConfig.DeliverPolicy = jetstream.DeliverByStartTimePolicy // Start from now
		startTime := time.Now().Add(-1 * time.Minute)                     // Slightly in the past
		consumerConfig.OptStartTime = &startTime
	}

	cons, err := eb.js.CreateOrUpdateConsumer(ctx, streamName, consumerConfig)
	if err != nil {
		eb.logger.ErrorContext(ctx, "Failed to create/update consumer", slog.String("stream", streamName), slog.String("consumer", consumerName), slog.Any("error", err))
		return nil, fmt.Errorf("failed to create/update consumer: %w", err)
	}

	eb.logger.InfoContext(ctx, "Consumer details",
		slog.String("stream", streamName),
		slog.String("consumer_name", cons.CachedInfo().Name),
		slog.String("consumer_durable_name", cons.CachedInfo().Config.Durable),
		slog.String("consumer_filter_subject", cons.CachedInfo().Config.FilterSubject),
		slog.Any("consumer_ack_policy", cons.CachedInfo().Config.AckPolicy),
	)

	messages := make(chan *message.Message)

	sub, err := cons.Messages() // Use cons.Messages() for auto unmarshaling
	if err != nil {
		eb.logger.ErrorContext(ctx, "Failed to start consumer", slog.String("consumer", consumerName), slog.Any("error", err))
		return nil, fmt.Errorf("failed to start consumer for topic %s: %w", topic, err)
	}

	go func() {
		eb.logger.InfoContext(ctx, "Starting consumer goroutine", slog.String("consumer_name", consumerName), slog.String("stream", streamName), slog.String("topic", topic))
		defer func() {
			close(messages)
			sub.Stop()
		}()

		for {
			jetStreamMsg, err := sub.Next()
			if err != nil {
				if errors.Is(err, jetstream.ErrMsgIteratorClosed) {
					eb.logger.InfoContext(ctx, "Consumer iterator closed", slog.String("consumer", consumerName), slog.Any("error", err))
				} else {
					eb.logger.ErrorContext(ctx, "Error receiving message", slog.String("consumer", consumerName), slog.Any("error", err))
				}
				return
			}

			wmMsg, err := eb.toWatermillMessage(ctx, jetStreamMsg) // Convert to Watermill message
			if err != nil {
				eb.logger.ErrorContext(ctx, "Failed to convert message", slog.Any("err", err))
				continue
			}

			eb.logger.InfoContext(ctx, "Received message", slog.String("stream", streamName), slog.String("consumer", consumerName),
				slog.String("topic", topic), slog.String("message_id", wmMsg.UUID), slog.Any("metadata", wmMsg.Metadata))

			select {
			case messages <- wmMsg:
				select {
				case <-wmMsg.Acked():
					if err := jetStreamMsg.Ack(); err != nil {
						eb.logger.ErrorContext(ctx, "Failed to ack message", slog.String("message_id", wmMsg.UUID), slog.Any("error", err))
					}
					eb.logger.InfoContext(ctx, "Message acknowledged", slog.String("message_id", wmMsg.UUID))
				case <-wmMsg.Nacked():
					if err := jetStreamMsg.Nak(); err != nil {
						eb.logger.ErrorContext(ctx, "Failed to nack message", slog.String("message_id", wmMsg.UUID), slog.Any("error", err))
					}
					eb.logger.InfoContext(ctx, "Message nacked", slog.String("message_id", wmMsg.UUID))
				case <-ctx.Done():
					eb.logger.WarnContext(ctx, "Context cancelled (ack/nack)", slog.String("message_id", wmMsg.UUID))
					if err := jetStreamMsg.Term(); err != nil {
						eb.logger.ErrorContext(ctx, "Failed to term message", slog.String("message_id", wmMsg.UUID), slog.Any("error", err))
					}
					return
				}
			case <-ctx.Done():
				eb.logger.WarnContext(ctx, "Context cancelled (sending)", slog.String("message_id", wmMsg.UUID))
				if err := jetStreamMsg.Term(); err != nil {
					eb.logger.ErrorContext(ctx, "Failed to term message", slog.String("message_id", wmMsg.UUID), slog.Any("error", err))
				}
				return
			}
		}
	}()

	eb.logger.InfoContext(ctx, "Subscribed to topic", slog.String("topic", topic), slog.String("consumer", consumerName))
	return messages, nil
}

func (eb *eventBus) toWatermillMessage(ctx context.Context, jetStreamMsg jetstream.Msg) (*message.Message, error) {
	watermillMsg := message.NewMessage(string(jetStreamMsg.Headers().Get("Nats-Msg-Id")), jetStreamMsg.Data())

	meta, err := jetStreamMsg.Metadata()
	if err == nil {
		watermillMsg.Metadata.Set("Stream", meta.Stream)
		watermillMsg.Metadata.Set("Consumer", meta.Consumer)
		watermillMsg.Metadata.Set("Delivered", strconv.FormatInt(int64(meta.NumDelivered), 10))
		watermillMsg.Metadata.Set("StreamSeq", strconv.FormatUint(meta.Sequence.Stream, 10))
		watermillMsg.Metadata.Set("ConsumerSeq", strconv.FormatUint(meta.Sequence.Consumer, 10))
		watermillMsg.Metadata.Set("Timestamp", meta.Timestamp.String())
	}

	for k, v := range jetStreamMsg.Headers() {
		watermillMsg.Metadata.Set(k, v[0])
	}

	if watermillMsg.Metadata.Get(middleware.CorrelationIDMetadataKey) == "" {
		watermillMsg.Metadata.Set(middleware.CorrelationIDMetadataKey, watermill.NewUUID())
	}
	watermillMsg.SetContext(ctx)
	return watermillMsg, nil
}

func streamSubjectsMatch(existing, new []string) bool {
	if len(existing) != len(new) {
		return false
	}
	seen := make(map[string]struct{})
	for _, sub := range existing {
		seen[sub] = struct{}{}
	}
	for _, sub := range new {
		if _, found := seen[sub]; !found {
			return false
		}
	}
	return true
}

func sanitize(s string) string {
	reg := regexp.MustCompile("[^a-zA-Z0-9-]+")
	return reg.ReplaceAllString(s, "")
}

func (eb *eventBus) createStreamsAndDLQs(ctx context.Context) error {
	streams := []string{"user", "leaderboard", "round", "score", "discord", delayedMessagesStream}
	for _, stream := range streams {
		if err := eb.CreateStream(ctx, stream); err != nil {
			return fmt.Errorf("failed to create stream %s: %w", stream, err)
		}
	}
	return nil
}

func (eb *eventBus) Close() error {
	if eb.publisher != nil {
		if err := eb.publisher.Close(); err != nil {
			eb.logger.Error("Error closing publisher", slog.Any("error", err))
		}
	}
	if eb.subscriber != nil {
		if err := eb.subscriber.Close(); err != nil {
			eb.logger.Error("Error closing subscriber", slog.Any("error", err))
		}
	}

	if eb.natsConn != nil {
		eb.natsConn.Close()
	}

	return nil
}

// CreateStream creates or updates a single JetStream stream (helper function).
func (eb *eventBus) CreateStream(ctx context.Context, streamName string) error {
	eb.logger.InfoContext(ctx, "Creating/Updating stream", slog.String("stream_name", streamName))
	eb.streamMutex.Lock() // Protect against concurrent stream creation
	defer eb.streamMutex.Unlock()

	if eb.createdStreams[streamName] {
		eb.logger.InfoContext(ctx, "Stream already created in this process", slog.String("stream_name", streamName))
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
	case delayedMessagesStream:
		subjects = []string{delayedMessagesSubject + ".>"}
	default:
		return fmt.Errorf("unknown stream name: %s", streamName)
	}

	// Define default stream configuration
	streamCfg := jetstream.StreamConfig{
		Name:      streamName,
		Subjects:  subjects,
		Retention: jetstream.LimitsPolicy, // Retain based on limits (age, size)
		Storage:   jetstream.FileStorage,  // Store on disk
	}

	// Customize for delayed messages stream
	switch streamName {
	case delayedMessagesStream:
		streamCfg.MaxAge = 24 * time.Hour            // Auto-delete after 24 hours
		streamCfg.MaxMsgs = -1                       // Unlimited messages (rely on MaxAge)
		streamCfg.Retention = jetstream.LimitsPolicy // Use limits policy
	default:
		streamCfg.Duplicates = 5 * time.Minute // 5-minute deduplication window
	}

	// Create or update the stream (idempotent)
	_, err := eb.js.CreateOrUpdateStream(ctx, streamCfg)
	if err != nil {
		return fmt.Errorf("failed to create or update stream %s: %w", streamName, err)
	}

	eb.logger.InfoContext(ctx, "Stream created or updated", slog.String("stream_name", streamName), slog.Any("subjects", subjects))
	eb.createdStreams[streamName] = true // Mark as created
	return nil
}
