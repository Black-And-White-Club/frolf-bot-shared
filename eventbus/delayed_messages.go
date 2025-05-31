package eventbus

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/Black-And-White-Club/frolf-bot-shared/observability/attr"
	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go/jetstream"
)

// CancelScheduledMessage cancels scheduled messages for a given roundID.
func (eb *eventBus) CancelScheduledMessage(ctx context.Context, roundID sharedtypes.RoundID) error {
	// Create a contextual logger for this operation
	ctxLogger := eb.logger.With(
		attr.String("operation", "cancel_scheduled_message"),
		attr.RoundID("round_id", roundID),
		attr.String("stream", DelayedMessagesStream),
	)

	ctxLogger.Info("Cancelling scheduled message")

	stream, err := eb.js.Stream(ctx, DelayedMessagesStream)
	if err != nil {
		ctxLogger.Error("Failed to access stream", attr.Error(err))
		return fmt.Errorf("failed to access stream %s: %w", DelayedMessagesStream, err)
	}

	subjectToDelete := fmt.Sprintf("%s.%v", DelayedMessagesSubject, roundID)
	ctxLogger = ctxLogger.With(attr.String("subject_to_delete", subjectToDelete))

	if err := stream.Purge(ctx, jetstream.WithPurgeSubject(subjectToDelete)); err != nil {
		ctxLogger.Error("Failed to purge messages", attr.Error(err))
		return fmt.Errorf("failed to purge messages for round %v: %w", roundID, err)
	}

	ctxLogger.Info("Successfully cancelled scheduled message")
	return nil
}

// RecoverScheduledRounds recovers scheduled rounds after restart.
func (eb *eventBus) RecoverScheduledRounds(ctx context.Context) {
	// Create a contextual logger for this operation
	ctxLogger := eb.logger.With(
		attr.String("operation", "recover_scheduled_rounds"),
		attr.String("stream", DelayedMessagesStream),
	)

	ctxLogger.Info("Recovering scheduled rounds after restart")

	consumerName := fmt.Sprintf("delayed_recovery_consumer_%d", time.Now().Unix())
	ctxLogger = ctxLogger.With(attr.String("consumer_name", consumerName))

	consumerConfig := jetstream.ConsumerConfig{
		Durable:       consumerName,
		FilterSubject: DelayedMessagesSubject + ".>",
		AckPolicy:     jetstream.AckExplicitPolicy,
		MaxAckPending: 100,
		DeliverPolicy: jetstream.DeliverAllPolicy,
		AckWait:       30 * time.Second,
	}

	consumer, err := eb.js.CreateOrUpdateConsumer(ctx, DelayedMessagesStream, consumerConfig)
	if err != nil {
		ctxLogger.Error("Failed to create recovery consumer", attr.Error(err))
		return
	}

	ctxLogger.Info("Recovery consumer created", attr.String("status", "consumer_created"))

	sub, err := consumer.Messages()
	if err != nil {
		ctxLogger.Error("Failed to fetch scheduled messages", attr.Error(err))
		return
	}

	ctxLogger.Info("Recovery subscription created", attr.String("status", "subscription_created"))

	defer func() {
		sub.Stop()
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := eb.js.DeleteConsumer(cleanupCtx, DelayedMessagesStream, consumerName); err != nil {
			ctxLogger.Warn("Failed to delete temporary recovery consumer", attr.Error(err))
		} else {
			ctxLogger.Info("Temporary recovery consumer deleted", attr.String("status", "consumer_deleted"))
		}
	}()

	recoveredRounds := make(map[sharedtypes.RoundID]bool)

	for {
		select {
		case <-ctx.Done():
			ctxLogger.Warn("Context canceled, stopping recovery", attr.String("reason", "context_canceled"), attr.Error(ctx.Err()))
			return
		default:
			fetchCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			jetStreamMsg, err := sub.Next()
			cancel()

			if err != nil {
				if errors.Is(err, jetstream.ErrMsgIteratorClosed) ||
					errors.Is(err, context.DeadlineExceeded) ||
					errors.Is(err, context.Canceled) {
					ctxLogger.Info("No more scheduled messages to recover or context canceled", attr.String("status", "recovery_complete"))
					return
				}
				ctxLogger.Error("Error fetching delayed message", attr.Error(err))
				continue
			}

			wmMsg, err := eb.toWatermillMessage(fetchCtx, jetStreamMsg)
			if err != nil {
				ctxLogger.Error("Error converting JetStream message to Watermill", attr.Error(err))
				jetStreamMsg.Term()
				continue
			}

			roundID := wmMsg.Metadata.Get("Round-ID")
			executeAtStr := wmMsg.Metadata.Get("Execute-At")

			if roundID == "" || executeAtStr == "" {
				ctxLogger.Error("Missing required metadata in delayed message", attr.String("error", "missing_metadata"))
				jetStreamMsg.Term()
				continue
			}

			executeAt, err := time.Parse(time.RFC3339, executeAtStr)
			if err != nil {
				ctxLogger.Error("Invalid Execute-At timestamp", attr.Error(err))
				jetStreamMsg.Term()
				continue
			}

			// Convert roundID string to sharedtypes.RoundID
			parsedRoundID, err := uuid.Parse(roundID) // Parse the string to a UUID
			if err != nil {
				ctxLogger.Error("Invalid roundID string in recovery", attr.String("error", "invalid_round_id"))
				jetStreamMsg.Term()
				continue
			}
			sharedRoundID := sharedtypes.RoundID(parsedRoundID) // Now convert to sharedtypes.RoundID

			if recoveredRounds[sharedRoundID] {
				ctxLogger.Debug("Round already recovered, skipping duplicate", attr.RoundID("round_id", sharedRoundID))
				jetStreamMsg.Ack()
				continue
			}

			if time.Now().Before(executeAt) {
				ctxLogger.Info("Recreating consumer for scheduled round", attr.RoundID("round_id", sharedRoundID))

				if err := eb.ProcessDelayedMessages(ctx, sharedRoundID, sharedtypes.StartTime(executeAt)); err != nil {
					ctxLogger.Error("Failed to process delayed messages", attr.Error(err))
				}

				recoveredRounds[sharedRoundID] = true

				ctxLogger.Info("Successfully recovered delayed round", attr.RoundID("round_id", sharedRoundID))
			} else {
				ctxLogger.Warn("Skipping recovery of past round",
					attr.RoundID("round_id", sharedRoundID),
					attr.String("status", "skipped_past_round"),
					attr.String("reason", "scheduled_time_in_past"),
				)
			}

			if err := jetStreamMsg.Ack(); err != nil {
				ctxLogger.Error("Failed to ack recovery message", attr.Error(err))
			}
		}
	}
}

// ProcessDelayedMessage handles delayed message processing with NakWithDelay.
func (eb *eventBus) ProcessDelayedMessage(ctx context.Context, jetStreamMsg jetstream.Msg, scheduledTime time.Time, roundID sharedtypes.RoundID) {
	// Create a contextual logger for this operation
	ctxLogger := eb.logger.With(
		attr.String("operation", "process_delayed_message"),
		attr.RoundID("round_id", roundID),
		attr.Time("scheduled_time", scheduledTime),
	)

	// Get current time in UTC
	now := time.Now().UTC()
	ctxLogger = ctxLogger.With(attr.Time("current_time", now))

	ctxLogger.Debug("Processing delayed message")

	if now.Before(scheduledTime) {
		// Not yet time to process, requeue with delay
		delay := scheduledTime.Sub(now)
		ctxLogger.Info("Message not yet due, requeuing with delay", attr.Duration("delay", delay))

		if err := jetStreamMsg.NakWithDelay(delay); err != nil {
			ctxLogger.Error("Failed to requeue message with delay", attr.Error(err))
		} else {
			ctxLogger.Info("Message successfully requeued with delay")
		}
		return
	}

	ctxLogger.Info("Time to process delayed message", attr.String("status", "processing"))

	// Convert to Watermill message
	wmMsg, err := eb.toWatermillMessage(ctx, jetStreamMsg)
	if err != nil {
		ctxLogger.Error("Error converting JetStream message to Watermill", attr.Error(err), attr.String("action", "terminating_message"))
		if termErr := jetStreamMsg.Term(); termErr != nil {
			ctxLogger.Error("Failed to terminate message", attr.Error(termErr))
		}
		return
	}

	// Get original subject
	originalSubject := wmMsg.Metadata.Get("Original-Subject")
	if originalSubject == "" {
		ctxLogger.Error("Delayed message missing Original-Subject", attr.String("error", "missing_original_subject"))
		if termErr := jetStreamMsg.Term(); termErr != nil {
			ctxLogger.Error("Failed to terminate message", attr.Error(termErr))
		}
		return
	}

	ctxLogger.Info("Publishing delayed message to original subject", attr.String("original_subject", originalSubject))

	// Add correlation ID if missing
	if middleware.MessageCorrelationID(wmMsg) == "" {
		correlationID := wmMsg.UUID
		wmMsg.Metadata.Set(middleware.CorrelationIDMetadataKey, correlationID)
		ctxLogger.Debug("Added correlation ID to message", attr.String("correlation_id", correlationID))
	} else {
		ctxLogger.Debug("Correlation ID already present", attr.String("correlation_id", middleware.MessageCorrelationID(wmMsg)))
	}

	// Publish the delayed message to the original subject
	if err := eb.publisher.Publish(originalSubject, wmMsg); err != nil {
		ctxLogger.Error("Failed to republish delayed message", attr.Error(err))
		if nakErr := jetStreamMsg.Nak(); nakErr != nil {
			ctxLogger.Error("Failed to nack message for retry", attr.Error(nakErr))
		}
		return
	}

	ctxLogger.Info("Successfully published delayed message to original subject", attr.String("original_subject", originalSubject))

	// Acknowledge the message
	if ackErr := jetStreamMsg.Ack(); ackErr != nil {
		ctxLogger.Error("Failed to ack message", attr.Error(ackErr))
	} else {
		ctxLogger.Info("Message acknowledged successfully", attr.String("original_subject", originalSubject))
	}
}

// ProcessDelayedMessages processes delayed messages with NakWithDelay.
func (eb *eventBus) ProcessDelayedMessages(ctx context.Context, roundID sharedtypes.RoundID, scheduledTime sharedtypes.StartTime) error {
	consumerName := fmt.Sprintf("delayed_processor_%s", roundID)

	// Create a contextual logger for this operation
	ctxLogger := eb.logger.With(
		attr.String("operation", "process_delayed_messages"),
		attr.RoundID("round_id", roundID),
		attr.Time("scheduled_time", scheduledTime.AsTime()),
		attr.String("consumer_name", consumerName),
		attr.String("stream", DelayedMessagesStream),
	)

	ctxLogger.Info("Setting up delayed message processor")

	// Subject filter for this specific round ID
	filterSubject := fmt.Sprintf("%s.%s", DelayedMessagesSubject, roundID)
	ctxLogger = ctxLogger.With(attr.String("filter_subject", filterSubject))

	// Create or update the consumer
	consumerConfig := jetstream.ConsumerConfig{
		Durable:       consumerName,
		FilterSubject: filterSubject,
		AckPolicy:     jetstream.AckExplicitPolicy,
		DeliverPolicy: jetstream.DeliverAllPolicy,
		MaxAckPending: 1000,
		AckWait:       30 * time.Second,
	}

	cons, err := eb.js.CreateOrUpdateConsumer(ctx, DelayedMessagesStream, consumerConfig)
	if err != nil {
		ctxLogger.Error("Failed to create/update delayed message consumer", attr.Error(err))
		return err
	}

	ctxLogger.Info("Delayed message consumer created/updated", attr.String("status", "consumer_created"))

	// Subscribe to messages
	sub, err := cons.Messages()
	if err != nil {
		ctxLogger.Error("Failed to subscribe to messages", attr.Error(err))
		return err
	}

	ctxLogger.Info("Subscription created for delayed messages", attr.String("status", "subscription_created"))

	// Process messages in a goroutine
	go func() {
		routineAttrs := ctxLogger.With(attr.String("goroutine", "delayed_processor"))
		routineAttrs.Info("Starting delayed message processor goroutine")

		defer func() {
			// Recover from panics
			if r := recover(); r != nil {
				routineAttrs.Error("Recovered from panic in delayed message processor",
					attr.String("panic", fmt.Sprintf("%v", r)),
					attr.String("stack", string(debug.Stack())),
				)
			}
			routineAttrs.Info("Stopping delayed message processor")
			sub.Stop()
		}()

		for {
			select {
			case <-ctx.Done():
				ctxAttrs := routineAttrs.With(
					attr.String("reason", "context_canceled"),
					attr.Error(ctx.Err()),
				)
				ctxAttrs.Warn("Stopping delayed message processor due to context cancellation")
				return
			default:
				func() {
					// Use a separate recover function to handle panics within message processing
					defer func() {
						if r := recover(); r != nil {
							routineAttrs.Error("Recovered from panic while processing message",
								attr.String("panic", fmt.Sprintf("%v", r)),
								attr.String("stack", string(debug.Stack())),
							)
						}
					}()

					jetStreamMsg, err := sub.Next()
					if err != nil {
						if errors.Is(err, jetstream.ErrMsgIteratorClosed) {
							routineAttrs.Info("Message iterator closed", attr.String("reason", "iterator_closed"))
							return
						}

						if errors.Is(err, context.Canceled) {
							routineAttrs.Info("Context canceled while waiting for message", attr.Error(err))
							return
						}

						routineAttrs.Error("Error receiving message", attr.Error(err))
						return
					}

					// Process the message
					msgMeta, _ := jetStreamMsg.Metadata()
					msgAttrs := routineAttrs.With(attr.String("message_subject", jetStreamMsg.Subject()))

					if msgMeta != nil {
						msgAttrs = msgAttrs.With(
							attr.Uint64("stream_seq", msgMeta.Sequence.Stream),
							attr.Uint64("delivery_count", msgMeta.NumDelivered),
						)
					}

					msgAttrs.Debug("Processing delayed message")
					eb.ProcessDelayedMessage(ctx, jetStreamMsg, scheduledTime.AsTime(), roundID)
				}()
			}
		}
	}()

	ctxLogger.Info("Delayed message processor started", attr.String("status", "started"))
	return nil
}

// ScheduleDelayedMessage schedules a message to be delivered at a future time
func (eb *eventBus) ScheduleDelayedMessage(ctx context.Context, originalSubject string, roundID sharedtypes.RoundID, scheduledTime sharedtypes.StartTime, payload []byte, additionalMetadata map[string]string) error {
	// Create a contextual logger for this operation
	ctxLogger := eb.logger.With(
		attr.String("operation", "schedule_delayed_message"),
		attr.String("original_subject", originalSubject),
		attr.RoundID("round_id", roundID),
		attr.Time("scheduled_time", scheduledTime.AsTime()),
	)

	ctxLogger.Info("Scheduling delayed message")

	// Validate that scheduled time is in the future
	now := time.Now().UTC()
	if scheduledTime.AsTime().Before(now) {
		ctxLogger.Warn("Scheduled time is in the past, skipping message scheduling",
			attr.Time("current_time", now),
			attr.Time("scheduled_time", scheduledTime.AsTime()),
		)
		return fmt.Errorf("cannot schedule message in the past: %s is before %s", scheduledTime.AsTime().Format(time.RFC3339), now.Format(time.RFC3339))
	}

	msg := message.NewMessage(watermill.NewUUID(), payload)

	msg.Metadata.Set("Original-Subject", originalSubject)
	msg.Metadata.Set("Round-ID", roundID.String())
	msg.Metadata.Set("Execute-At", scheduledTime.AsTime().Format(time.RFC3339))

	// Add any additional metadata
	for key, value := range additionalMetadata {
		msg.Metadata.Set(key, value)
	}

	// Prepare the delayed subject for logging
	delayedSubject := fmt.Sprintf("%s.%v", DelayedMessagesSubject, roundID)

	// Log the message publishing details
	ctxLogger.Info("Publishing delayed message",
		attr.String("message_id", msg.UUID),
		attr.String("correlation_id", middleware.MessageCorrelationID(msg)),
		attr.String("delayed_subject", delayedSubject),
	)

	err := eb.Publish(delayedSubject, msg)
	if err != nil {
		ctxLogger.Error("Failed to publish delayed message", attr.Error(err))
		return fmt.Errorf("failed to publish delayed message: %w", err)
	}

	ctxLogger.Info("Successfully scheduled delayed message",
		attr.String("message_id", msg.UUID),
		attr.String("correlation_id", middleware.MessageCorrelationID(msg)),
		attr.String("delayed_subject", delayedSubject),
	)

	// Optionally process delayed messages immediately after scheduling
	eb.ProcessDelayedMessages(ctx, roundID, scheduledTime)

	return nil
}
