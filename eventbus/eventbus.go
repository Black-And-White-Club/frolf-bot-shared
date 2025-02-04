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
	publisher         message.Publisher
	subscriber        message.Subscriber
	js                jetstream.JetStream
	natsConn          *nc.Conn
	logger            *slog.Logger
	createdStreams    map[string]bool
	streamMutex       sync.Mutex
	processedMessages map[string]bool
}

const (
	delayedMessagesStream  = "delayed"
	delayedMessagesSubject = "delayed.message"
	discordEventsSubject   = "discord.round.event"
)

// EventBus defines the interface for an event bus that can publish and subscribe to messages.
type EventBus interface {
	Publish(topic string, messages ...*message.Message) error
	Subscribe(ctx context.Context, topic string) (<-chan *message.Message, error)
	Close() error
	CreateOrUpdateStream(ctx context.Context, streamCfg jetstream.StreamConfig) (jetstream.Stream, error)
	ProcessDelayedMessages(ctx context.Context)
}

func NewEventBus(ctx context.Context, natsURL string, logger *slog.Logger) (EventBus, error) {
	// Connect to NATS
	natsConn, err := nc.Connect(natsURL,
		nc.RetryOnFailedConnect(true),
		nc.Timeout(30*time.Second),
		nc.ReconnectWait(1*time.Second),
		nc.MaxReconnects(-1),
	)
	if err != nil {
		logger.Error("Failed to connect to NATS", slog.Any("error", err))
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	// Initialize JetStream
	js, err := jetstream.New(natsConn)
	if err != nil {
		natsConn.Close()
		logger.Error("Failed to initialize JetStream", slog.Any("error", err))
		return nil, fmt.Errorf("failed to initialize JetStream: %w", err)
	}

	// Create a Watermill logger that wraps slog
	watermillLogger := watermill.NewSlogLogger(logger)

	// Create a Marshaller for the publisher
	marshaller := &nats.NATSMarshaler{}

	// Initialize the publisher with JetStream and asynchronous publishing
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
				// **Enable JetStream and don't disable it**
				Disabled:      false,
				AutoProvision: true, // You might want to set this to false in production and manage streams manually
				TrackMsgId:    true,
				DurablePrefix: "durable",
				// Set this to function to get more flexibility
				DurableCalculator: func(durablePrefix, topic string) string {
					// Create a sanitized version of the topic by replacing dots with underscores
					sanitizedTopic := strings.ReplaceAll(topic, ".", "_")
					return fmt.Sprintf("%s-%s", durablePrefix, sanitizedTopic)
				},
			},
		},
		watermillLogger,
	)
	if err != nil {
		natsConn.Close()
		logger.Error("Failed to create Watermill publisher", slog.Any("error", err))
		return nil, fmt.Errorf("failed to create Watermill publisher: %w", err)
	}

	// Create the subscriber
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
		logger.Error("Failed to create Watermill subscriber", slog.Any("error", err))
		return nil, fmt.Errorf("failed to create Watermill subscriber: %w", err)
	}

	eventBus := &eventBus{
		publisher:         publisher,
		js:                js,
		natsConn:          natsConn,
		logger:            logger,
		createdStreams:    make(map[string]bool),
		subscriber:        subscriber,
		processedMessages: make(map[string]bool),
	}

	// Create streams and DLQ topics
	if err := eventBus.createStreamsAndDLQs(ctx); err != nil {
		natsConn.Close()
		return nil, fmt.Errorf("failed to create streams and DLQs: %w", err)
	}

	go eventBus.ProcessDelayedMessages(ctx)

	return eventBus, nil
}

func (eb *eventBus) CreateStream(ctx context.Context, streamName string) error {
	eb.logger.Info("Creating stream", slog.String("stream_name", streamName))
	eb.streamMutex.Lock()
	defer eb.streamMutex.Unlock()

	if eb.createdStreams[streamName] {
		eb.logger.Info("Stream already created in this process", slog.String("stream_name", streamName))
		return nil
	}

	// Define subjects for the stream
	var subjects []string
	switch streamName {
	case "user":
		subjects = []string{
			"user.>", // This captures all events prefixed with "user."
		}
	case "leaderboard":
		subjects = []string{
			"leaderboard.>", // This captures all events prefixed with "leaderboard."
		}
	case "round":
		subjects = []string{
			"round.>", // This captures all events prefixed with "round."
		}
	case "score":
		subjects = []string{
			"score.>", // This captures all events prefixed with "score."
		}
	case "discord":
		subjects = []string{
			"discord.>", // This captures all events prefixed with "discord."
		}
	case delayedMessagesStream:
		subjects = []string{
			delayedMessagesSubject,
		}
	default:
		return fmt.Errorf("unknown stream name: %s", streamName)
	}

	// Define stream configuration with appropriate defaults
	streamCfg := jetstream.StreamConfig{
		Name:      streamName,
		Subjects:  subjects,
		Retention: jetstream.LimitsPolicy, // Default retention policy
		Storage:   jetstream.FileStorage,  // Default storage type
		Replicas:  3,                      // Default number of replicas
	}

	// Customize stream configuration based on stream name
	switch streamName {
	case delayedMessagesStream:
		streamCfg.MaxAge = 24 * time.Hour            // Auto-delete old messages
		streamCfg.MaxMsgs = -1                       // Unlimited messages
		streamCfg.Replicas = 3                       // High availability
		streamCfg.Retention = jetstream.LimitsPolicy // Retain messages until consumed
	default:
		streamCfg.Duplicates = 5 * time.Minute // Default deduplication window for other streams
	}

	// Check and create or update the stream
	stream, err := eb.js.Stream(ctx, streamName)
	if err == jetstream.ErrStreamNotFound {
		_, err = eb.js.CreateStream(ctx, streamCfg)
		if err != nil {
			return fmt.Errorf("failed to create stream %s: %w", streamName, err)
		}
		eb.logger.Info("Stream created", slog.String("stream_name", streamName), slog.Any("subjects", subjects))
	} else if err == nil {
		currentCfg := stream.CachedInfo().Config
		if !streamSubjectsMatch(currentCfg.Subjects, subjects) ||
			currentCfg.Retention != streamCfg.Retention ||
			currentCfg.Storage != streamCfg.Storage ||
			currentCfg.Replicas != streamCfg.Replicas ||
			currentCfg.MaxAge != streamCfg.MaxAge {
			_, err = eb.js.UpdateStream(ctx, streamCfg)
			if err != nil {
				return fmt.Errorf("failed to update stream %s: %w", streamName, err)
			}
			eb.logger.Info("Stream updated", slog.String("stream_name", streamName), slog.Any("subjects", subjects))
		} else {
			eb.logger.Info("Stream config is up-to-date", slog.String("stream_name", streamName))
		}
	} else {
		return fmt.Errorf("failed to retrieve or create stream: %w", err)
	}

	eb.createdStreams[streamName] = true
	return nil
}

func (eb *eventBus) Subscribe(ctx context.Context, topic string) (<-chan *message.Message, error) {
	eb.logger.Info("Entering Subscribe", slog.String("topic", topic))

	// Determine the stream name based on the topic
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
	default:
		return nil, fmt.Errorf("unknown topic: %s", topic)
	}

	// Use a durable consumer name specific to the topic
	consumerName := fmt.Sprintf("consumer-%s", sanitize(topic))

	// Check if the stream exists
	_, err := eb.js.Stream(ctx, streamName)
	if err == jetstream.ErrStreamNotFound {
		eb.logger.Error("Stream not found", slog.String("stream", streamName), slog.Any("error", err))
		return nil, fmt.Errorf("stream not found: %w", err)
	}

	// Define consumer configuration
	consumerConfig := jetstream.ConsumerConfig{
		Durable:       consumerName,
		FilterSubject: topic,
		AckPolicy:     jetstream.AckExplicitPolicy,
		MaxAckPending: 2048,
		DeliverPolicy: jetstream.DeliverAllPolicy,
	}

	// Special handling for delayed messages
	if topic == delayedMessagesSubject {
		consumerConfig.DeliverPolicy = jetstream.DeliverByStartTimePolicy
		startTime := time.Now().Add(-1 * time.Minute) // Start slightly in the past to catch any immediate messages
		consumerConfig.OptStartTime = &startTime
	}

	// Create or update a consumer for the stream, subscribing to the specific topic
	cons, err := eb.js.CreateOrUpdateConsumer(ctx, streamName, consumerConfig)
	if err != nil {
		eb.logger.Error("Failed to create or update consumer",
			slog.String("stream", streamName),
			slog.String("consumer", consumerName),
			slog.Any("error", err),
		)
		return nil, fmt.Errorf("failed to create or update consumer: %w", err)
	}

	// Log consumer details for debugging
	eb.logger.Info("Consumer details",
		slog.String("stream", streamName),
		slog.String("consumer_name", cons.CachedInfo().Name),
		slog.String("consumer_durable_name", cons.CachedInfo().Config.Durable),
		slog.String("consumer_filter_subject", cons.CachedInfo().Config.FilterSubject),
		slog.Any("consumer_ack_policy", cons.CachedInfo().Config.AckPolicy),
	)

	// Create a channel for received messages
	messages := make(chan *message.Message)

	// Subscribe using a JetStream consumer
	sub, err := cons.Messages()
	if err != nil {
		eb.logger.Error("Failed to start consumer", slog.String("consumer", consumerName), slog.Any("error", err))
		return nil, fmt.Errorf("failed to start consumer for topic %s: %w", topic, err)
	}

	// Start a goroutine to handle JetStream messages
	go func() {
		eb.logger.Info("Starting consumer goroutine", slog.String("consumer_name", consumerName), slog.String("stream", streamName), slog.String("topic", topic))
		defer close(messages)

		for {
			jetStreamMsg, err := sub.Next()
			if err != nil {
				if errors.Is(err, jetstream.ErrMsgIteratorClosed) {
					eb.logger.Info("Consumer messages iterator closed", slog.String("consumer", consumerName), slog.Any("error", err))
				} else {
					eb.logger.Error("Error receiving message from consumer", slog.String("consumer", consumerName), slog.Any("error", err))
				}
				return
			}

			// Convert JetStream message to Watermill message
			watermillMsg := message.NewMessage(string(jetStreamMsg.Headers().Get("Nats-Msg-Id")), jetStreamMsg.Data())

			// Add JetStream metadata to the Watermill message
			meta, err := jetStreamMsg.Metadata()
			if err == nil {
				watermillMsg.Metadata.Set("Stream", meta.Stream)
				watermillMsg.Metadata.Set("Consumer", meta.Consumer)
				watermillMsg.Metadata.Set("Delivered", strconv.FormatInt(int64(meta.NumDelivered), 10))
				watermillMsg.Metadata.Set("StreamSeq", strconv.FormatUint(meta.Sequence.Stream, 10))
				watermillMsg.Metadata.Set("ConsumerSeq", strconv.FormatUint(meta.Sequence.Consumer, 10))
				watermillMsg.Metadata.Set("Timestamp", meta.Timestamp.String())
			}

			// Copy metadata from the original message
			for k, v := range jetStreamMsg.Headers() {
				watermillMsg.Metadata.Set(k, v[0])
			}

			// Set correlation ID if not already set
			if watermillMsg.Metadata.Get(middleware.CorrelationIDMetadataKey) == "" {
				watermillMsg.Metadata.Set(middleware.CorrelationIDMetadataKey, watermill.NewUUID())
			}

			eb.logger.Info("Received message from JetStream consumer",
				slog.String("stream", streamName),
				slog.String("consumer", consumerName),
				slog.String("topic", topic),
				slog.String("message_id", watermillMsg.UUID),
				slog.Any("metadata", watermillMsg.Metadata),
			)

			// Send the message to the messages channel and handle acknowledgements
			select {
			case messages <- watermillMsg:
				if err := jetStreamMsg.Ack(); err != nil {
					eb.logger.Error("Failed to acknowledge message", slog.String("message_id", watermillMsg.UUID), slog.Any("error", err))
				}
				eb.logger.Info("Message acknowledged", slog.String("message_id", watermillMsg.UUID))
			case <-ctx.Done():
				eb.logger.Warn("Context cancelled, not processing message", slog.String("message_id", watermillMsg.UUID))
				if err := jetStreamMsg.Nak(); err != nil {
					eb.logger.Error("Failed to Nack message", slog.String("message_id", watermillMsg.UUID), slog.Any("error", err))
				}
				return
			}
		}
	}()

	eb.logger.Info("Successfully subscribed to topic using consumer", slog.String("topic", topic), slog.String("consumer", consumerName))
	return messages, nil
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

// Helper function to sanitize a string to make it a valid NATS subject part
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

// Publish publishes messages to the specified topic.
func (eb *eventBus) Publish(topic string, messages ...*message.Message) error {
	for _, msg := range messages {
		// Ensure the message has a unique UUID for Watermill deduplication
		if msg.UUID == "" {
			msg.UUID = watermill.NewUUID()
		}

		// Set correlation ID for traceability
		if msg.Metadata.Get(middleware.CorrelationIDMetadataKey) == "" {
			msg.Metadata.Set(middleware.CorrelationIDMetadataKey, watermill.NewUUID())
		}

		// Set Nats-Msg-Id for JetStream deduplication
		dedupKey := fmt.Sprintf("%s-%s", msg.UUID, sanitize(topic))
		msg.Metadata.Set("topic", topic)
		msg.Metadata.Set("Nats-Msg-Id", dedupKey)

		// Check if this message is for delayed publishing
		if executeAtStr := msg.Metadata.Get("Execute-At"); executeAtStr != "" {
			eb.logger.Info("Publishing delayed message", "topic", topic, "message_id", msg.UUID, "execute_at", executeAtStr)

			// Store original subject in metadata
			msg.Metadata.Set("Original-Subject", topic)

			// Publish to delayed messages stream
			err := eb.publisher.Publish(delayedMessagesSubject, msg)
			if err != nil {
				return fmt.Errorf("failed to publish delayed message to topic %s: %w", delayedMessagesSubject, err)
			}

			eb.logger.Info("Delayed message published", "original_topic", topic, "message_id", msg.UUID, "execute_at", executeAtStr)
			continue // Skip immediate publishing for this message
		}

		// Use JetStream publishing with specific options (e.g., MsgId)
		jsMsg := &nc.Msg{
			Subject: topic,
			Data:    msg.Payload,
			Header:  nc.Header{},
		}

		// Add Watermill message ID as Nats-Msg-Id header
		jsMsg.Header.Set("Nats-Msg-Id", dedupKey)

		// Add other metadata as headers
		for k, v := range msg.Metadata {
			jsMsg.Header.Add(k, v)
		}

		// Use JetStream to publish
		_, err := eb.js.PublishMsg(context.Background(), jsMsg)
		if err != nil {
			return fmt.Errorf("failed to publish message to topic %s: %w", topic, err)
		}

		eb.logger.Info("Message published successfully",
			slog.String("topic", topic),
			slog.String("message_id", msg.UUID),
			slog.String("dedup_key", dedupKey),
		)
	}
	return nil
}
