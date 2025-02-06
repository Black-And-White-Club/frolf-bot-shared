package eventbus

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/nats-io/nats.go/jetstream"
)

func (eb *eventBus) ProcessDelayedMessages(ctx context.Context) {
	// Ensure a durable consumer exists
	consumer, err := eb.js.CreateOrUpdateConsumer(ctx, "DELAYED_MESSAGES", jetstream.ConsumerConfig{
		Durable:       "delayed_processor",
		FilterSubject: "delayed.stream",
		AckPolicy:     jetstream.AckExplicitPolicy, // Ensure explicit ACKs
	})
	if err != nil {
		eb.logger.Error("Failed to create/update consumer", "error", err)
		return
	}

	// Subscribe with a callback function (JetStream will call this function per message)
	sub, err := consumer.Consume(func(jetStreamMsg jetstream.Msg) {
		msgId := jetStreamMsg.Headers().Get("Nats-Msg-Id")

		// **Check if message was already processed**
		if _, exists := eb.processedMessages[msgId]; exists {
			eb.logger.Warn("Skipping duplicate execution", "msg_id", msgId)
			jetStreamMsg.Ack()
			return
		}

		// **Extract Execution Time**
		executeAtStr := jetStreamMsg.Headers().Get("Execute-At")
		executeAt, err := time.Parse(time.RFC3339, executeAtStr)
		if err != nil {
			eb.logger.Error("Invalid Execute-At timestamp, NAKing message", "msg_id", msgId, "error", err)
			jetStreamMsg.Nak()
			return
		}

		// **Reschedule message if execution time hasn't arrived**
		delay := time.Until(executeAt)
		if delay > 0 {
			eb.logger.Info("Message not ready, rescheduling", "msg_id", msgId, "delay", delay)
			jetStreamMsg.NakWithDelay(delay)
			return
		}

		// **Mark message as processed**
		eb.processedMessages[msgId] = true

		// **Republish message to correct subject**
		originalSubject := jetStreamMsg.Headers().Get("Original-Subject")
		_, err = eb.js.Publish(ctx, originalSubject, jetStreamMsg.Data())
		if err != nil {
			eb.logger.Error("Failed to republish delayed message", "msg_id", msgId, "error", err)
			jetStreamMsg.Nak() // Retry later
			return
		}

		// **Acknowledge message**
		jetStreamMsg.Ack()
		eb.logger.Info("Processed delayed message", "msg_id", msgId, "subject", originalSubject)
	})
	if err != nil {
		eb.logger.Error("Failed to start consumer", "error", err)
		return
	}
	defer sub.Stop()

	// Wait until context is done before stopping the processor
	<-ctx.Done()
	eb.logger.Warn("Stopping delayed message processor due to context cancellation")
}

func (eb *eventBus) CancelScheduledMessage(ctx context.Context, roundID string) error {
	// Access the stream
	stream, err := eb.js.Stream(ctx, delayedMessagesStream)
	if err != nil {
		return fmt.Errorf("failed to access stream %s: %w", delayedMessagesStream, err)
	}

	// Construct the subject for the specific round
	subjectToDelete := fmt.Sprintf("%s.%s", delayedMessagesSubject, roundID)

	// Purge messages on the specific subject
	if err := stream.Purge(ctx, jetstream.WithPurgeSubject(subjectToDelete)); err != nil {
		return fmt.Errorf("failed to purge messages for round %s: %w", roundID, err)
	}

	eb.logger.Info("Purged scheduled messages from delayed stream", "subject", subjectToDelete, "round_id", roundID)
	return nil
}

// CreateOrUpdateStream creates or updates a JetStream stream with the given configuration.
func (eb *eventBus) CreateOrUpdateStream(ctx context.Context, streamCfg jetstream.StreamConfig) (jetstream.Stream, error) {
	eb.logger.Info("Creating/updating stream", slog.String("stream_name", streamCfg.Name))
	eb.streamMutex.Lock()
	defer eb.streamMutex.Unlock()

	js, err := jetstream.New(eb.natsConn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to jetstream: %w", err)
	}

	stream, err := js.Stream(ctx, streamCfg.Name)
	if err != nil {
		if errors.Is(err, jetstream.ErrStreamNotFound) {
			stream, err = js.CreateStream(ctx, streamCfg)
			if err != nil {
				return nil, fmt.Errorf("failed to create stream %s: %w", streamCfg.Name, err)
			}
			eb.logger.Info("Stream created", slog.String("stream_name", streamCfg.Name), slog.Any("subjects", streamCfg.Subjects))
		} else {
			return nil, fmt.Errorf("failed to get stream %s: %w", streamCfg.Name, err)
		}
	} else {
		eb.logger.Info("Stream found, checking config", slog.String("stream_name", streamCfg.Name))
		currentCfg := stream.CachedInfo().Config
		if !streamSubjectsMatch(currentCfg.Subjects, streamCfg.Subjects) ||
			currentCfg.Retention != streamCfg.Retention ||
			currentCfg.Storage != streamCfg.Storage ||
			currentCfg.Replicas != streamCfg.Replicas {
			_, err = js.UpdateStream(ctx, streamCfg)
			if err != nil {
				return nil, fmt.Errorf("failed to update stream %s: %w", streamCfg.Name, err)
			}
			eb.logger.Info("Stream updated", slog.String("stream_name", streamCfg.Name), slog.Any("subjects", streamCfg.Subjects))
		} else {
			eb.logger.Info("Stream config is up-to-date", slog.String("stream_name", streamCfg.Name))
		}
	}

	eb.createdStreams[streamCfg.Name] = true
	return stream, nil
}
