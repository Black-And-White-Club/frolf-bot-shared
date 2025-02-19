// eventbus/delayed.go
package eventbus

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/nats-io/nats.go/jetstream"
)

func (eb *eventBus) ProcessDelayedMessages(ctx context.Context) {
	// Ensure a durable consumer exists, subscribing to ALL delayed messages.
	consumer, err := eb.js.CreateOrUpdateConsumer(ctx, delayedMessagesStream, jetstream.ConsumerConfig{
		Durable:       "delayed_processor",
		FilterSubject: delayedMessagesSubject + ".>", // Subscribe to ALL delayed messages
		AckPolicy:     jetstream.AckExplicitPolicy,
		DeliverPolicy: jetstream.DeliverByStartTimePolicy,
		OptStartTime:  &[]time.Time{time.Now()}[0],
	})
	if err != nil {
		eb.logger.ErrorContext(ctx, "Failed to create/update consumer", slog.Any("error", err))
		return
	}

	sub, err := consumer.Consume(func(jetStreamMsg jetstream.Msg) {
		executeAtStr := jetStreamMsg.Headers().Get("Execute-At")
		executeAt, err := time.Parse(time.RFC3339, executeAtStr)
		if err != nil {
			eb.logger.ErrorContext(ctx, "Invalid Execute-At timestamp, Terminating message", slog.Any("error", err))
			jetStreamMsg.Term()
			return
		}

		if time.Now().After(executeAt) {
			wmMsg, err := eb.toWatermillMessage(ctx, jetStreamMsg)
			if err != nil {
				eb.logger.ErrorContext(ctx, "Error converting jetstream message to watermill", slog.Any("error", err))
				jetStreamMsg.Term()
				return
			}

			originalSubject := wmMsg.Metadata.Get("Original-Subject")
			if originalSubject == "" {
				eb.logger.ErrorContext(ctx, "Delayed message missing Original-Subject", slog.String("message_id", wmMsg.UUID))
				jetStreamMsg.Term()
				return
			}

			// Clear out the delay-related metadata
			wmMsg.Metadata.Set("Execute-At", "")
			wmMsg.Metadata.Set("Round-ID", "") // Clear Round-ID
			wmMsg.Metadata.Set("Nats-Delay", "")

			if err := eb.publisher.Publish(originalSubject, wmMsg); err != nil {
				eb.logger.ErrorContext(ctx, "Failed to republish delayed message", slog.String("original_topic", originalSubject), slog.String("message_id", wmMsg.UUID), slog.Any("error", err))
				jetStreamMsg.Nak() // Retry
				return
			}

			jetStreamMsg.Ack()
			eb.logger.InfoContext(ctx, "Processed delayed message", slog.String("subject", originalSubject), slog.String("message_id", wmMsg.UUID))

		} else {
			eb.logger.InfoContext(ctx, "Message not ready, Nacking", slog.String("delay", time.Until(executeAt).String()))
			jetStreamMsg.Nak()
		}
	})
	if err != nil {
		eb.logger.ErrorContext(ctx, "Failed to start consumer", slog.Any("error", err))
		return
	}

	<-ctx.Done()
	eb.logger.WarnContext(ctx, "Stopping delayed message processor due to context cancellation")
	sub.Stop()
}

func (eb *eventBus) CancelScheduledMessage(ctx context.Context, roundID string) error {
	// Access the stream
	stream, err := eb.js.Stream(ctx, delayedMessagesStream)
	if err != nil {
		return fmt.Errorf("failed to access stream %s: %w", delayedMessagesStream, err)
	}

	// Construct the subject for the specific round.  This is now correct.
	subjectToDelete := fmt.Sprintf("%s.%s", delayedMessagesSubject, roundID)

	// Purge messages on the specific subject
	if err := stream.Purge(ctx, jetstream.WithPurgeSubject(subjectToDelete)); err != nil {
		return fmt.Errorf("failed to purge messages for round %s: %w", roundID, err)
	}

	eb.logger.InfoContext(ctx, "Cancelled scheduled message", slog.String("subject", subjectToDelete), slog.String("round_id", roundID))
	return nil
}
