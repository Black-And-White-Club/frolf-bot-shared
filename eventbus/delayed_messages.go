package eventbus

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/Black-And-White-Club/frolf-bot-shared/observability/attr"
	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/nats-io/nats.go/jetstream"
)

// ProcessDelayedMessage handles delayed message processing with NakWithDelay.
func (eb *eventBus) ProcessDelayedMessage(ctx context.Context, jetStreamMsg jetstream.Msg, scheduledTime time.Time, roundID sharedtypes.RoundID) {
	logAttrs := []attr.LogAttr{
		attr.String("operation", "process_delayed_message"),
		attr.Int("round_id", int(roundID)),
		attr.Time("scheduled_time", scheduledTime),
	}

	// Get current time in UTC
	now := time.Now().UTC()
	logAttrs = append(logAttrs, attr.Time("current_time", now))

	eb.logger.Debug("Processing delayed message", logAttrs...)

	if now.Before(scheduledTime) {
		// Not yet time to process, requeue with delay
		delay := scheduledTime.Sub(now)
		requeueAttrs := append(logAttrs, attr.Duration("delay", delay))
		eb.logger.Info("Message not yet due, requeuing with delay", requeueAttrs...)

		if err := jetStreamMsg.NakWithDelay(delay); err != nil {
			errorAttrs := append(requeueAttrs, attr.Error(err))
			eb.logger.Error("Failed to requeue message with delay", errorAttrs...)
		} else {
			eb.logger.Info("Message successfully requeued with delay", requeueAttrs...)
		}
		return
	}

	processingAttrs := append(logAttrs, attr.String("status", "processing"))
	eb.logger.Info("Time to process delayed message", processingAttrs...)

	// Convert to Watermill message
	wmMsg, err := eb.toWatermillMessage(ctx, jetStreamMsg)
	if err != nil {
		errorAttrs := append(logAttrs,
			attr.Error(err),
			attr.String("action", "terminating_message"),
		)
		eb.logger.Error("Error converting JetStream message to Watermill", errorAttrs...)
		if err := jetStreamMsg.Term(); err != nil {
			termAttrs := append(errorAttrs, attr.Error(err))
			eb.logger.Error("Failed to terminate message", termAttrs...)
		}
		return
	}

	// Inject trace context into the message
	eb.tracer.InjectTraceContext(ctx, wmMsg)

	msgAttrs := append(logAttrs, attr.String("message_id", wmMsg.UUID))

	// Get original subject
	originalSubject := wmMsg.Metadata.Get("Original-Subject")
	if originalSubject == "" {
		errorAttrs := append(msgAttrs, attr.String("error", "missing_original_subject"))
		eb.logger.Error("Delayed message missing Original-Subject", errorAttrs...)
		if err := jetStreamMsg.Term(); err != nil {
			termAttrs := append(errorAttrs, attr.Error(err))
			eb.logger.Error("Failed to terminate message", termAttrs...)
		}
		return
	}

	publishAttrs := append(msgAttrs, attr.String("original_subject", originalSubject))
	eb.logger.Info("Publishing delayed message to original subject", publishAttrs...)

	// Add correlation ID if missing
	if middleware.MessageCorrelationID(wmMsg) == "" {
		correlationID := wmMsg.UUID
		wmMsg.Metadata.Set(middleware.CorrelationIDMetadataKey, correlationID)
		publishAttrs = append(publishAttrs, attr.String("correlation_id", correlationID))
		eb.logger.Debug("Added correlation ID to message", publishAttrs...)
	} else {
		publishAttrs = append(publishAttrs,
			attr.String("correlation_id", middleware.MessageCorrelationID(wmMsg)),
		)
	}

	// Publish the delayed message to the original subject
	if err := eb.publisher.Publish(originalSubject, wmMsg); err != nil {
		errorAttrs := append(publishAttrs, attr.Error(err))
		eb.logger.Error("Failed to republish delayed message", errorAttrs...)
		if err := jetStreamMsg.Nak(); err != nil {
			nakAttrs := append(errorAttrs, attr.Error(err))
			eb.logger.Error("Failed to nack message for retry", nakAttrs...)
		}
		return
	}

	eb.logger.Info("Successfully published delayed message to original subject", publishAttrs...)

	// Acknowledge the message
	if err := jetStreamMsg.Ack(); err != nil {
		ackAttrs := append(publishAttrs, attr.Error(err))
		eb.logger.Error("Failed to ack message", ackAttrs...)
	} else {
		eb.logger.Info("Message acknowledged successfully", publishAttrs...)
	}
}

// ProcessDelayedMessages processes delayed messages with NakWithDelay.
func (eb *eventBus) ProcessDelayedMessages(ctx context.Context, roundID sharedtypes.RoundID, scheduledTime time.Time) {
	roundIDStr := strconv.FormatInt(int64(roundID), 10)
	consumerName := fmt.Sprintf("delayed_processor_%s", roundIDStr)

	logAttrs := []attr.LogAttr{
		attr.String("operation", "process_delayed_messages"),
		attr.Int("round_id", int(roundID)),
		attr.String("round_id_str", roundIDStr),
		attr.Time("scheduled_time", scheduledTime),
		attr.String("consumer_name", consumerName),
		attr.String("stream", DelayedMessagesStream),
	}

	eb.logger.Info("Setting up delayed message processor", logAttrs...)

	// Subject filter for this specific round ID
	filterSubject := fmt.Sprintf("%s.%s", DelayedMessagesSubject, roundIDStr)
	logAttrs = append(logAttrs, attr.String("filter_subject", filterSubject))

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
		errorAttrs := append(logAttrs, attr.Error(err))
		eb.logger.Error("Failed to create/update delayed message consumer", errorAttrs...)
		return
	}

	consAttrs := append(logAttrs, attr.String("status", "consumer_created"))
	eb.logger.Info("Delayed message consumer created/updated", consAttrs...)

	// Subscribe to messages
	sub, err := cons.Messages()
	if err != nil {
		errorAttrs := append(consAttrs, attr.Error(err))
		eb.logger.Error("Failed to subscribe to messages", errorAttrs...)
		return
	}

	subAttrs := append(consAttrs, attr.String("status", "subscription_created"))
	eb.logger.Info("Subscription created for delayed messages", subAttrs...)

	// Process messages in a goroutine
	go func() {
		routineAttrs := append(logAttrs, attr.String("goroutine", "delayed_processor"))
		eb.logger.Info("Starting delayed message processor goroutine", routineAttrs...)

		defer func() {
			eb.logger.Info("Stopping delayed message processor", routineAttrs...)
			sub.Stop()
		}()

		for {
			select {
			case <-ctx.Done():
				ctxAttrs := append(routineAttrs,
					attr.String("reason", "context_canceled"),
					attr.Error(ctx.Err()),
				)
				eb.logger.Warn("Stopping delayed message processor due to context cancellation", ctxAttrs...)
				return
			default:
				jetStreamMsg, err := sub.Next()
				if err != nil {
					if errors.Is(err, jetstream.ErrMsgIteratorClosed) {
						iteratorAttrs := append(routineAttrs, attr.String("reason", "iterator_closed"))
						eb.logger.Info("Message iterator closed", iteratorAttrs...)
						return
					}

					if errors.Is(err, context.Canceled) {
						cancelAttrs := append(routineAttrs, attr.Error(err))
						eb.logger.Info("Context canceled while waiting for message", cancelAttrs...)
						return
					}

					errorAttrs := append(routineAttrs, attr.Error(err))
					eb.logger.Error("Error receiving message", errorAttrs...)
					continue
				}

				// Process the message
				msgMeta, _ := jetStreamMsg.Metadata()
				msgAttrs := append(routineAttrs,
					attr.String("message_subject", jetStreamMsg.Subject()),
				)

				if msgMeta != nil {
					msgAttrs = append(msgAttrs,
						attr.Uint64("stream_seq", msgMeta.Sequence.Stream),
						attr.Uint64("delivery_count", msgMeta.NumDelivered),
					)
				}

				eb.logger.Debug("Processing delayed message", msgAttrs...)
				eb.ProcessDelayedMessage(ctx, jetStreamMsg, scheduledTime, roundID)
			}
		}
	}()

	eb.logger.Info("Delayed message processor started", append(logAttrs, attr.String("status", "started"))...)
}

// CancelScheduledMessage cancels scheduled messages for a given roundID.
func (eb *eventBus) CancelScheduledMessage(ctx context.Context, roundID sharedtypes.RoundID) error {
	logAttrs := []attr.LogAttr{
		attr.String("operation", "cancel_scheduled_message"),
		attr.Int("round_id", int(roundID)),
		attr.String("stream", DelayedMessagesStream),
	}

	eb.logger.Info("Cancelling scheduled message", logAttrs...)

	stream, err := eb.js.Stream(ctx, DelayedMessagesStream)
	if err != nil {
		errorAttrs := append(logAttrs, attr.Error(err))
		eb.logger.Error("Failed to access stream", errorAttrs...)
		return fmt.Errorf("failed to access stream %s: %w", DelayedMessagesStream, err)
	}

	subjectToDelete := fmt.Sprintf("%s.%d ", DelayedMessagesSubject, roundID)
	logAttrs = append(logAttrs, attr.String("subject_to_delete", subjectToDelete))

	if err := stream.Purge(ctx, jetstream.WithPurgeSubject(subjectToDelete)); err != nil {
		errorAttrs := append(logAttrs, attr.Error(err))
		eb.logger.Error("Failed to purge messages", errorAttrs...)
		return fmt.Errorf("failed to purge messages for round %v: %w", roundID, err)
	}

	eb.logger.Info("Successfully cancelled scheduled message", logAttrs...)
	return nil
}

// RecoverScheduledRounds recovers scheduled rounds after restart.
func (eb *eventBus) RecoverScheduledRounds(ctx context.Context) {
	logAttrs := []attr.LogAttr{
		attr.String("operation", "recover_scheduled_rounds"),
		attr.String("stream", DelayedMessagesStream),
	}

	eb.logger.Info("Recovering scheduled rounds after restart", logAttrs...)

	consumerName := fmt.Sprintf("delayed_recovery_consumer_%d", time.Now().Unix())
	logAttrs = append(logAttrs, attr.String("consumer_name", consumerName))

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
		errorAttrs := append(logAttrs, attr.Error(err))
		eb.logger.Error("Failed to create recovery consumer", errorAttrs...)
		return
	}

	consAttrs := append(logAttrs, attr.String("status", "consumer_created"))
	eb.logger.Info("Recovery consumer created", consAttrs...)

	sub, err := consumer.Messages()
	if err != nil {
		errorAttrs := append(consAttrs, attr.Error(err))
		eb.logger.Error("Failed to fetch scheduled messages", errorAttrs...)
		return
	}

	subAttrs := append(consAttrs, attr.String("status", "subscription_created"))
	eb.logger.Info("Recovery subscription created", subAttrs...)

	defer func() {
		sub.Stop()
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := eb.js.DeleteConsumer(cleanupCtx, DelayedMessagesStream, consumerName); err != nil {
			eb.logger.Warn("Failed to delete temporary recovery consumer",
				append(subAttrs, attr.Error(err))...)
		} else {
			eb.logger.Info("Temporary recovery consumer deleted",
				append(subAttrs, attr.String("status", "consumer_deleted"))...)
		}
	}()

	recoveredRounds := make(map[sharedtypes.RoundID]bool)

	for {
		select {
		case <-ctx.Done():
			cancelAttrs := append(subAttrs,
				attr.String("reason", "context_canceled"),
				attr.Error(ctx.Err()),
			)
			eb.logger.Warn("Context canceled, stopping recovery", cancelAttrs...)
			return
		default:
			fetchCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			jetStreamMsg, err := sub.Next()
			cancel()

			if err != nil {
				if errors.Is(err, jetstream.ErrMsgIteratorClosed) ||
					errors.Is(err, context.DeadlineExceeded) ||
					errors.Is(err, context.Canceled) {
					completeAttrs := append(subAttrs, attr.String("status", "recovery_complete"))
					eb.logger.Info("No more scheduled messages to recover or context canceled", completeAttrs...)
					return
				}
				errorAttrs := append(subAttrs, attr.Error(err))
				eb.logger.Error("Error fetching delayed message", errorAttrs...)
				continue
			}

			wmMsg, err := eb.toWatermillMessage(fetchCtx, jetStreamMsg)
			if err != nil {
				errorAttrs := append(subAttrs, attr.Error(err))
				eb.logger.Error("Error converting JetStream message to Watermill", errorAttrs...)
				jetStreamMsg.Term()
				continue
			}

			// Inject trace context into the message
			eb.tracer.InjectTraceContext(ctx, wmMsg)

			roundIDStr := wmMsg.Metadata.Get("Round-ID")
			executeAtStr := wmMsg.Metadata.Get("Execute-At")

			msgAttrs := append(subAttrs,
				attr.String("message_id", wmMsg.UUID),
				attr.String("round_id_str", roundIDStr),
				attr.String("execute_at", executeAtStr),
			)

			if roundIDStr == "" || executeAtStr == "" {
				errorAttrs := append(msgAttrs, attr.String("error", "missing_metadata"))
				eb.logger.Error("Missing required metadata in delayed message", errorAttrs...)
				jetStreamMsg.Term()
				continue
			}

			executeAt, err := time.Parse(time.RFC3339, executeAtStr)
			if err != nil {
				errorAttrs := append(msgAttrs, attr.Error(err))
				eb.logger.Error("Invalid Execute-At timestamp", errorAttrs...)
				jetStreamMsg.Term()
				continue
			}

			msgAttrs = append(msgAttrs, attr.Time("execute_at_parsed", executeAt))

			roundIDInt, err := strconv.ParseInt(roundIDStr, 10, 64)
			if err != nil {
				errorAttrs := append(msgAttrs, attr.Error(err))
				eb.logger.Error("Invalid roundID string in recovery", errorAttrs...)
				jetStreamMsg.Term()
				continue
			}

			roundID := sharedtypes.RoundID(roundIDInt)
			msgAttrs = append(msgAttrs, attr.Int("round_id", int(roundID)))

			if recoveredRounds[roundID] {
				eb.logger.Debug("Round already recovered, skipping duplicate", msgAttrs...)
				jetStreamMsg.Ack()
				continue
			}

			if time.Now().Before(executeAt) {
				msgAttrs = append(msgAttrs, attr.String("status", "recreating_consumer"))
				eb.logger.Info("Recreating consumer for scheduled round", msgAttrs...)

				eb.ProcessDelayedMessages(ctx, roundID, executeAt)

				recoveredRounds[roundID] = true

				msgAttrs = append(msgAttrs, attr.String("status", "recovery_successful"))
				eb.logger.Info("Successfully recovered delayed round", msgAttrs...)
			} else {
				msgAttrs = append(msgAttrs,
					attr.String("status", "skipped_past_round"),
					attr.String("reason", "scheduled_time_in_past"),
				)
				eb.logger.Warn("Skipping recovery of past round", msgAttrs...)
			}

			if err := jetStreamMsg.Ack(); err != nil {
				ackAttrs := append(msgAttrs, attr.Error(err))
				eb.logger.Error("Failed to ack recovery message", ackAttrs...)
			}
		}
	}
}

// ScheduleDelayedMessage schedules a message to be delivered at a future time
func (eb *eventBus) ScheduleDelayedMessage(ctx context.Context, originalSubject string, roundID sharedtypes.RoundID, scheduledTime time.Time, payload []byte) error {
	logAttrs := []attr.LogAttr{
		attr.String("operation", "schedule_delayed_message"),
		attr.String("original_subject", originalSubject),
		attr.Int("round_id", int(roundID)),
		attr.Time("scheduled_time", scheduledTime),
	}

	eb.logger.Info("Scheduling delayed message", logAttrs...)

	msg := message.NewMessage(watermill.NewUUID(), payload)

	msg.Metadata.Set("Original-Subject", originalSubject)
	msg.Metadata.Set("Round-ID", strconv.FormatInt(int64(roundID), 10))
	msg.Metadata.Set("Execute-At", scheduledTime.Format(time.RFC3339))

	if middleware.MessageCorrelationID(msg) == "" {
		correlationID := watermill.NewUUID()
		msg.Metadata.Set(middleware.CorrelationIDMetadataKey, correlationID)
	}

	// Inject trace context into the message
	eb.tracer.InjectTraceContext(ctx, msg)

	msgAttrs := append(logAttrs,
		attr.String("message_id", msg.UUID),
		attr.String("correlation_id", middleware.MessageCorrelationID(msg)),
	)

	delayedSubject := fmt.Sprintf("%s.%d", DelayedMessagesSubject, roundID)
	msgAttrs = append(msgAttrs, attr.String("delayed_subject", delayedSubject))

	eb.logger.Info("Publishing delayed message", msgAttrs...)

	err := eb.Publish(delayedSubject, msg)
	if err != nil {
		errorAttrs := append(msgAttrs, attr.Error(err))
		eb.logger.Error("Failed to publish delayed message", errorAttrs...)
		return fmt.Errorf("failed to publish delayed message: %w", err)
	}

	eb.logger.Info("Successfully scheduled delayed message", msgAttrs...)

	eb.ProcessDelayedMessages(ctx, roundID, scheduledTime)

	return nil
}
