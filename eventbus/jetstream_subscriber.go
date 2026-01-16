package eventbus

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	nc "github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

const (
	// defaultMaxConcurrentAcks is the default limit for concurrent ACK handling goroutines.
	defaultMaxConcurrentAcks = 50

	// pullBatchSize is the number of messages to fetch in each pull request.
	pullBatchSize = 10

	// pullExpiry is how long to wait for messages in each pull request.
	pullExpiry = 5 * time.Second
)

// JetStreamSubscriberAdapter implements message.Subscriber using native JetStream pull consumers.
// It provides bounded concurrency for ACK handling and graceful shutdown support.
type JetStreamSubscriberAdapter struct {
	js              jetstream.JetStream
	consumerManager *ConsumerManager
	appType         string
	logger          *slog.Logger

	maxConcurrentAcks int

	mu            sync.Mutex
	subscriptions []*subscription
	closed        atomic.Bool

	// termination tracks topics that should terminate messages instead of retrying.
	terminationMu sync.RWMutex
	termination   map[string]bool
}

// subscription tracks an active subscription for cleanup during Close.
type subscription struct {
	topic       string
	cancel      context.CancelFunc
	msgContext  jetstream.MessagesContext
	wg          *sync.WaitGroup
	outputCh    chan *message.Message
	closedOnce  sync.Once
	closedDone  chan struct{}
}

// JetStreamSubscriberOption configures the JetStreamSubscriberAdapter.
type JetStreamSubscriberOption func(*JetStreamSubscriberAdapter)

// WithMaxConcurrentAcks sets the maximum number of concurrent ACK handling goroutines.
// This bounds memory usage and prevents overwhelming the NATS server with ACK requests.
func WithMaxConcurrentAcks(max int) JetStreamSubscriberOption {
	return func(s *JetStreamSubscriberAdapter) {
		if max > 0 {
			s.maxConcurrentAcks = max
		}
	}
}

// NewJetStreamSubscriberAdapter creates a new subscriber adapter.
func NewJetStreamSubscriberAdapter(
	js jetstream.JetStream,
	consumerManager *ConsumerManager,
	appType string,
	logger *slog.Logger,
	opts ...JetStreamSubscriberOption,
) *JetStreamSubscriberAdapter {
	s := &JetStreamSubscriberAdapter{
		js:                js,
		consumerManager:   consumerManager,
		appType:           appType,
		logger:            logger,
		maxConcurrentAcks: defaultMaxConcurrentAcks,
		termination:       make(map[string]bool),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// Subscribe creates a subscription to the given topic.
// It returns a channel of Watermill messages that integrates with Watermill routers.
// The subscription uses a native JetStream pull consumer with bounded ACK concurrency.
func (s *JetStreamSubscriberAdapter) Subscribe(ctx context.Context, topic string) (<-chan *message.Message, error) {
	if s.closed.Load() {
		return nil, errors.New("subscriber is closed")
	}

	ctxLogger := s.logger.With(
		"operation", "subscribe",
		"topic", topic,
		"app_type", s.appType,
	)

	// Resolve stream from topic
	streamName, err := ResolveStreamFromTopic(topic)
	if err != nil {
		ctxLogger.ErrorContext(ctx, "Failed to resolve stream from topic", "error", err)
		return nil, err
	}

	// Ensure the consumer exists (idempotent)
	cons, err := s.consumerManager.EnsureConsumer(ctx, streamName, topic, s.appType)
	if err != nil {
		return nil, err
	}

	// Create the message iterator
	msgCtx, err := cons.Messages(
		jetstream.PullMaxMessages(pullBatchSize),
		jetstream.PullExpiry(pullExpiry),
	)
	if err != nil {
		ctxLogger.ErrorContext(ctx, "Failed to create message iterator", "error", err)
		return nil, fmt.Errorf("failed to create message iterator: %w", err)
	}

	// Create subscription context that we control
	subCtx, cancel := context.WithCancel(ctx)

	outputCh := make(chan *message.Message)
	wg := &sync.WaitGroup{}

	sub := &subscription{
		topic:      topic,
		cancel:     cancel,
		msgContext: msgCtx,
		wg:         wg,
		outputCh:   outputCh,
		closedDone: make(chan struct{}),
	}

	// Track subscription for cleanup
	s.mu.Lock()
	s.subscriptions = append(s.subscriptions, sub)
	s.mu.Unlock()

	// Start the message pump
	go s.messagePump(subCtx, sub, ctxLogger)

	ctxLogger.InfoContext(ctx, "Subscription created",
		"stream", streamName,
		"max_concurrent_acks", s.maxConcurrentAcks,
	)

	return outputCh, nil
}

// messagePump reads messages from JetStream and sends them to the output channel.
// It uses a semaphore to bound the number of concurrent ACK handling goroutines.
func (s *JetStreamSubscriberAdapter) messagePump(ctx context.Context, sub *subscription, logger *slog.Logger) {
	defer func() {
		// Stop the message iterator
		sub.msgContext.Stop()

		// Wait for all in-flight ACK goroutines to complete
		sub.wg.Wait()

		// Close the output channel
		sub.closedOnce.Do(func() {
			close(sub.outputCh)
			close(sub.closedDone)
		})

		logger.InfoContext(ctx, "Message pump stopped", "topic", sub.topic)
	}()

	// Semaphore for bounded ACK concurrency
	ackSem := make(chan struct{}, s.maxConcurrentAcks)

	for {
		// Check for shutdown
		select {
		case <-ctx.Done():
			logger.DebugContext(ctx, "Context cancelled, stopping message pump")
			return
		default:
		}

		// Fetch next message (this blocks with PullExpiry timeout)
		jsMsg, err := sub.msgContext.Next()
		if err != nil {
			if errors.Is(err, jetstream.ErrMsgIteratorClosed) {
				logger.DebugContext(ctx, "Message iterator closed")
				return
			}
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return
			}
			// Log unexpected errors but continue
			logger.WarnContext(ctx, "Error fetching message", "error", err)
			continue
		}

		// Convert to Watermill message
		wmMsg, err := s.toWatermillMessage(ctx, jsMsg)
		if err != nil {
			logger.ErrorContext(ctx, "Failed to convert message", "error", err)
			// On conversion error, delay redelivery
			if nakErr := jsMsg.NakWithDelay(5 * time.Second); nakErr != nil {
				logger.ErrorContext(ctx, "Failed to nak message", "error", nakErr)
			}
			continue
		}

		// Acquire semaphore slot before sending to output
		// This bounds memory usage from ACK goroutines
		select {
		case ackSem <- struct{}{}:
			// Got slot
		case <-ctx.Done():
			// Shutdown - nak for redelivery
			if nakErr := jsMsg.Nak(); nakErr != nil {
				logger.ErrorContext(ctx, "Failed to nak message during shutdown", "error", nakErr)
			}
			return
		}

		// Track this ACK goroutine
		sub.wg.Add(1)

		// Send to output channel
		select {
		case sub.outputCh <- wmMsg:
			// Message sent, spawn ACK handler
			go s.waitForAck(ctx, jsMsg, wmMsg, sub, ackSem, logger)

		case <-ctx.Done():
			// Shutdown - release semaphore and nak
			<-ackSem
			sub.wg.Done()
			if nakErr := jsMsg.Nak(); nakErr != nil {
				logger.ErrorContext(ctx, "Failed to nak message during shutdown", "error", nakErr)
			}
			return
		}
	}
}

// waitForAck listens for Watermill ACK/NACK signals and translates them to JetStream.
func (s *JetStreamSubscriberAdapter) waitForAck(
	ctx context.Context,
	jsMsg jetstream.Msg,
	wmMsg *message.Message,
	sub *subscription,
	ackSem chan struct{},
	logger *slog.Logger,
) {
	defer func() {
		// Release semaphore slot
		<-ackSem
		// Mark goroutine as done
		sub.wg.Done()
	}()

	// Get configured AckWait to set a reasonable timeout
	cfg := s.consumerManager.GetRegistry().Resolve(s.appType, sub.topic)
	// Use 80% of AckWait as our timeout to leave room for the actual ACK
	processTimeout := time.Duration(float64(cfg.AckWait) * 0.8)

	ackCtx, cancel := context.WithTimeout(ctx, processTimeout)
	defer cancel()

	select {
	case <-wmMsg.Acked():
		if err := jsMsg.Ack(); err != nil {
			logger.ErrorContext(ctx, "Failed to ack message", "error", err, "message_id", wmMsg.UUID)
		} else {
			logger.DebugContext(ctx, "Message acked", "message_id", wmMsg.UUID)
		}

	case <-wmMsg.Nacked():
		// Check if we should terminate instead of retry
		if s.shouldTerminate(sub.topic) {
			if err := jsMsg.Term(); err != nil {
				logger.ErrorContext(ctx, "Failed to terminate message", "error", err, "message_id", wmMsg.UUID)
			} else {
				logger.InfoContext(ctx, "Message terminated (marked for termination)", "message_id", wmMsg.UUID)
			}
			return
		}

		if err := jsMsg.Nak(); err != nil {
			logger.ErrorContext(ctx, "Failed to nak message", "error", err, "message_id", wmMsg.UUID)
		} else {
			logger.DebugContext(ctx, "Message nacked for redelivery", "message_id", wmMsg.UUID)
		}

	case <-ackCtx.Done():
		// Timeout - let AckWait handle redelivery
		logger.WarnContext(ctx, "ACK timeout, letting AckWait handle redelivery",
			"message_id", wmMsg.UUID,
			"timeout", processTimeout,
		)

	case <-ctx.Done():
		// Graceful shutdown - nak for redelivery
		if err := jsMsg.Nak(); err != nil {
			logger.ErrorContext(ctx, "Failed to nak message during shutdown", "error", err, "message_id", wmMsg.UUID)
		}
	}
}

// toWatermillMessage converts a JetStream message to a Watermill message.
func (s *JetStreamSubscriberAdapter) toWatermillMessage(ctx context.Context, jsMsg jetstream.Msg) (*message.Message, error) {
	// Use Nats-Msg-Id header as message ID, or generate one
	msgID := jsMsg.Headers().Get("Nats-Msg-Id")
	if msgID == "" {
		msgID = nc.NewInbox()
	}

	wmMsg := message.NewMessage(msgID, jsMsg.Data())

	// Copy all headers to metadata
	for k, v := range jsMsg.Headers() {
		if len(v) > 0 {
			wmMsg.Metadata.Set(k, v[0])
		}
	}

	// Add JetStream-specific metadata with _js_ prefix
	if meta, err := jsMsg.Metadata(); err == nil {
		wmMsg.Metadata.Set("_js_stream", meta.Stream)
		wmMsg.Metadata.Set("_js_consumer", meta.Consumer)
		wmMsg.Metadata.Set("_js_num_delivered", strconv.FormatUint(meta.NumDelivered, 10))
		wmMsg.Metadata.Set("_js_stream_seq", strconv.FormatUint(meta.Sequence.Stream, 10))
		wmMsg.Metadata.Set("_js_consumer_seq", strconv.FormatUint(meta.Sequence.Consumer, 10))
		wmMsg.Metadata.Set("_js_timestamp", meta.Timestamp.Format(time.RFC3339Nano))
		wmMsg.Metadata.Set("_js_domain", meta.Domain)

		// Also set legacy keys for backward compatibility
		wmMsg.Metadata.Set("stream", meta.Stream)
		wmMsg.Metadata.Set("consumer", meta.Consumer)
		wmMsg.Metadata.Set("deliveries", strconv.FormatUint(meta.NumDelivered, 10))
		wmMsg.Metadata.Set("stream_seq", strconv.FormatUint(meta.Sequence.Stream, 10))
		wmMsg.Metadata.Set("consumer_seq", strconv.FormatUint(meta.Sequence.Consumer, 10))
	}

	// Set message context derived from parent
	wmMsg.SetContext(context.WithValue(ctx, messageIDKey{}, wmMsg.UUID))

	return wmMsg, nil
}

// messageIDKey is the context key for message ID.
type messageIDKey struct{}

// shouldTerminate checks if a topic is marked for message termination.
func (s *JetStreamSubscriberAdapter) shouldTerminate(topic string) bool {
	s.terminationMu.RLock()
	defer s.terminationMu.RUnlock()
	return s.termination[topic]
}

// MarkForTermination marks a topic so that NACKed messages are terminated instead of retried.
// This is useful for poison messages that should not be redelivered.
func (s *JetStreamSubscriberAdapter) MarkForTermination(topic string, terminate bool) {
	s.terminationMu.Lock()
	defer s.terminationMu.Unlock()
	if terminate {
		s.termination[topic] = true
	} else {
		delete(s.termination, topic)
	}
}

// Close gracefully shuts down the subscriber.
// It cancels all subscription contexts, waits for in-flight ACK goroutines,
// and closes all output channels.
func (s *JetStreamSubscriberAdapter) Close() error {
	if !s.closed.CompareAndSwap(false, true) {
		return nil // Already closed
	}

	s.logger.Info("Closing JetStream subscriber adapter")

	s.mu.Lock()
	subs := s.subscriptions
	s.subscriptions = nil
	s.mu.Unlock()

	// Cancel all subscription contexts
	for _, sub := range subs {
		sub.cancel()
	}

	// Wait for all subscriptions to finish
	for _, sub := range subs {
		<-sub.closedDone
	}

	s.logger.Info("JetStream subscriber adapter closed")
	return nil
}
