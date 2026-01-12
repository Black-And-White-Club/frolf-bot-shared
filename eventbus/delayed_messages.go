package eventbus

// import (
// 	"context"
// 	"errors"
// 	"fmt"
// 	"log/slog"
// 	"time"

// 	"github.com/Black-And-White-Club/frolf-bot-shared/observability/attr"
// 	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
// 	"github.com/ThreeDotsLabs/watermill"
// 	"github.com/ThreeDotsLabs/watermill/message"
// 	"github.com/nats-io/nats.go/jetstream"
// )

// const (
// 	GlobalDelayedConsumerName = "global_delayed_processor"
// )

// // Just delete and recreate the consumer, not the stream
// func (eb *eventBus) StartGlobalDelayedMessageProcessor(ctx context.Context) error {
// 	fmt.Printf("[DELAYED DEBUG] Starting efficient delayed message processor\n")

// 	ctxLogger := eb.logger.With(
// 		attr.String("operation", "start_global_delayed_processor"),
// 		attr.String("consumer_name", GlobalDelayedConsumerName),
// 	)

// 	// Delete just the consumer (not the stream!)
// 	if err := eb.js.DeleteConsumer(ctx, DelayedMessagesStream, GlobalDelayedConsumerName); err != nil {
// 		fmt.Printf("[DELAYED DEBUG] Delete consumer result (expected if not exists): %v\n", err)
// 	}
// 	time.Sleep(1 * time.Second)

// 	// Corrected consumer config with only valid fields
// 	consumerConfig := jetstream.ConsumerConfig{
// 		Durable:       GlobalDelayedConsumerName,
// 		FilterSubject: DelayedMessagesSubject + ".>",
// 		AckPolicy:     jetstream.AckExplicitPolicy,
// 		DeliverPolicy: jetstream.DeliverAllPolicy,
// 		MaxAckPending: 1,               // Process one at a time for timing accuracy
// 		AckWait:       2 * time.Minute, // Reasonable timeout for processing
// 		MaxDeliver:    3,               // Limited retries
// 		BackOff: []time.Duration{ // Reasonable backoff
// 			30 * time.Second,
// 			2 * time.Minute,
// 			5 * time.Minute,
// 		},
// 		InactiveThreshold: 5 * time.Minute,               // Clean up inactive consumers
// 		ReplayPolicy:      jetstream.ReplayInstantPolicy, // Send messages as fast as possible
// 		HeadersOnly:       false,                         // We need the full message
// 	}

// 	consumer, err := eb.js.CreateConsumer(ctx, DelayedMessagesStream, consumerConfig)
// 	if err != nil {
// 		ctxLogger.Error("Failed to create consumer", attr.Error(err))
// 		return fmt.Errorf("failed to create consumer: %w", err)
// 	}

// 	// Start subscription
// 	sub, err := consumer.Messages()
// 	if err != nil {
// 		ctxLogger.Error("Failed to create subscription", attr.Error(err))
// 		return fmt.Errorf("failed to create subscription: %w", err)
// 	}

// 	// Start processor
// 	go eb.runGlobalDelayedProcessor(ctx, sub, ctxLogger)

// 	ctxLogger.Info("Efficient delayed message processor started")
// 	fmt.Printf("[DELAYED DEBUG] Efficient processor started\n")
// 	return nil
// }

// // Add this method to reset pending messages
// func (eb *eventBus) ResetDelayedConsumer(ctx context.Context) error {
// 	ctxLogger := eb.logger.With(attr.String("operation", "reset_delayed_consumer"))

// 	fmt.Printf("[RESET] Resetting delayed consumer to clear pending messages\n")
// 	ctxLogger.Info("Resetting delayed consumer")

// 	// Delete the existing consumer
// 	if err := eb.js.DeleteConsumer(ctx, DelayedMessagesStream, GlobalDelayedConsumerName); err != nil {
// 		// It's okay if the consumer doesn't exist
// 		fmt.Printf("[RESET] Consumer deletion result: %v\n", err)
// 	}

// 	// Wait a moment
// 	time.Sleep(2 * time.Second)

// 	// Restart the processor (which will recreate the consumer)
// 	return eb.StartGlobalDelayedMessageProcessor(ctx)
// }

// // Much simpler and more efficient delayed message processor
// func (eb *eventBus) runGlobalDelayedProcessor(ctx context.Context, sub jetstream.MessagesContext, ctxLogger *slog.Logger) {
// 	fmt.Printf("[DELAYED GOROUTINE] Efficient delayed processor started\n")

// 	routineLogger := ctxLogger.With(attr.String("goroutine", "global_delayed_processor"))
// 	routineLogger.Info("Efficient delayed processor started - will block until messages are ready")

// 	defer func() {
// 		if r := recover(); r != nil {
// 			fmt.Printf("[DELAYED GOROUTINE] PANIC recovered: %v\n", r)
// 			routineLogger.Error("Recovered from panic in delayed processor",
// 				attr.String("panic", fmt.Sprintf("%v", r)),
// 			)
// 		}
// 		fmt.Printf("[DELAYED GOROUTINE] Processor stopping\n")
// 		sub.Stop()
// 	}()

// 	messageCount := 0
// 	consecutiveErrors := 0

// 	for {
// 		select {
// 		case <-ctx.Done():
// 			routineLogger.Info("Context canceled, stopping delayed processor")
// 			return
// 		default:
// 			messageCount++

// 			// EFFICIENT: Block until JetStream delivers a ready message
// 			// No polling, no timeouts, no waste!
// 			fmt.Printf("[DELAYED GOROUTINE] Waiting for ready message #%d...\n", messageCount)
// 			jetStreamMsg, err := sub.Next()
// 			if err != nil {
// 				consecutiveErrors++
// 				fmt.Printf("[DELAYED GOROUTINE] Error: %v (consecutive: %d)\n", err, consecutiveErrors)

// 				// Handle iterator closure - restart processor
// 				if errors.Is(err, jetstream.ErrMsgIteratorClosed) {
// 					fmt.Printf("[DELAYED GOROUTINE] Iterator closed - restarting\n")
// 					go func() {
// 						time.Sleep(2 * time.Second)
// 						eb.StartGlobalDelayedMessageProcessor(ctx)
// 					}()
// 					return
// 				}

// 				if errors.Is(err, context.Canceled) {
// 					return
// 				}

// 				// Too many errors - restart
// 				if consecutiveErrors >= 5 {
// 					fmt.Printf("[DELAYED GOROUTINE] Too many errors - restarting\n")
// 					go func() {
// 						time.Sleep(5 * time.Second)
// 						eb.StartGlobalDelayedMessageProcessor(ctx)
// 					}()
// 					return
// 				}

// 				// Simple backoff and continue
// 				time.Sleep(time.Duration(consecutiveErrors) * time.Second)
// 				continue
// 			}

// 			// Success - reset error counter
// 			consecutiveErrors = 0

// 			// Process the message (JetStream already handled the timing!)
// 			fmt.Printf("[DELAYED GOROUTINE] Processing ready message #%d\n", messageCount)
// 			if err := eb.processDelayedMessage(ctx, jetStreamMsg, routineLogger); err != nil {
// 				routineLogger.Error("Failed to process message", attr.Error(err))
// 			}
// 		}
// 	}
// }

// // Update processDelayedMessage to validate round still exists and is current
// func (eb *eventBus) processDelayedMessage(ctx context.Context, jetStreamMsg jetstream.Msg, routineLogger *slog.Logger) error {
// 	msgMetadata, err := jetStreamMsg.Metadata()
// 	if err != nil {
// 		jetStreamMsg.Term()
// 		return fmt.Errorf("metadata error: %w", err)
// 	}

// 	fmt.Printf("[DELAYED MESSAGE] Processing delivery #%d of %s\n",
// 		msgMetadata.NumDelivered, jetStreamMsg.Subject())

// 	// Convert to Watermill message
// 	wmMsg, err := eb.toWatermillMessage(ctx, jetStreamMsg)
// 	if err != nil {
// 		jetStreamMsg.Term()
// 		return fmt.Errorf("conversion error: %w", err)
// 	}

// 	// Extract metadata
// 	roundID := wmMsg.Metadata.Get("Round-ID")
// 	executeAtStr := wmMsg.Metadata.Get("Execute-At")
// 	originalSubject := wmMsg.Metadata.Get("Original-Subject")

// 	if roundID == "" || executeAtStr == "" || originalSubject == "" {
// 		fmt.Printf("[DELAYED MESSAGE] Missing required metadata - terminating\n")
// 		jetStreamMsg.Term()
// 		return fmt.Errorf("missing metadata")
// 	}

// 	// Parse execution time
// 	executeAt, err := time.Parse(time.RFC3339, executeAtStr)
// 	if err != nil {
// 		fmt.Printf("[TIME DEBUG] Failed to parse time '%s': %v\n", executeAtStr, err)
// 		jetStreamMsg.Term()
// 		return fmt.Errorf("invalid timestamp: %w", err)
// 	}

// 	now := time.Now().UTC()
// 	executeAtUTC := executeAt.UTC()
// 	timeSinceExecute := now.Sub(executeAtUTC)
// 	messageAge := now.Sub(msgMetadata.Timestamp)

// 	// Debug logging
// 	fmt.Printf("[TIME DEBUG] Raw Execute-At: '%s'\n", executeAtStr)
// 	fmt.Printf("[TIME DEBUG] Current time: %s\n", now.Format(time.RFC3339))
// 	fmt.Printf("[TIME DEBUG] Execute time: %s\n", executeAtUTC.Format(time.RFC3339))
// 	fmt.Printf("[TIME DEBUG] Message age: %v\n", messageAge)
// 	fmt.Printf("[TIME DEBUG] Time since execute: %v\n", timeSinceExecute)

// 	// VALIDATION: Check if the round still exists (for very overdue messages)
// 	if timeSinceExecute > 30*time.Minute {
// 		fmt.Printf("[DELAYED MESSAGE] Message is %v overdue - checking if round still exists\n", timeSinceExecute)

// 		// You'd need to inject your round service/DB here to check if round exists
// 		// For now, we'll just log a warning and continue processing
// 		routineLogger.Warn("Processing significantly overdue delayed message",
// 			attr.String("round_id", roundID),
// 			attr.Duration("overdue_by", timeSinceExecute),
// 			attr.Duration("message_age", messageAge),
// 			attr.String("original_subject", originalSubject))

// 		// TODO: Add actual round existence check here if needed:
// 		// if eb.roundDB != nil {
// 		//     round, err := eb.roundDB.GetRound(ctx, roundID)
// 		//     if err != nil {
// 		//         fmt.Printf("[DELAYED MESSAGE] Round %s no longer exists - terminating message\n", roundID)
// 		//         jetStreamMsg.Term()
// 		//         return nil
// 		//     }
// 		//
// 		//     // Check if start time changed significantly (indicating an update occurred)
// 		//     if !round.StartTime.AsTime().Equal(executeAtUTC) {
// 		//         fmt.Printf("[DELAYED MESSAGE] Round start time changed - terminating outdated message\n")
// 		//         jetStreamMsg.Term()
// 		//         return nil
// 		//     }
// 		// }
// 	}

// 	// Check if ready to execute
// 	if now.Before(executeAtUTC) {
// 		delay := executeAtUTC.Sub(now)
// 		fmt.Printf("[DELAYED MESSAGE] Not ready - delaying for: %v\n", delay)

// 		// Sanity check for unreasonably long delays
// 		if delay > 365*24*time.Hour {
// 			fmt.Printf("[DELAYED MESSAGE] WARNING: Extremely long delay (%v) - possible time zone issue\n", delay)
// 			routineLogger.Warn("Extremely long delay detected",
// 				attr.String("round_id", roundID),
// 				attr.Duration("delay", delay),
// 				attr.String("execute_at", executeAtStr))
// 		}

// 		if err := jetStreamMsg.NakWithDelay(delay); err != nil {
// 			return fmt.Errorf("delay failed: %w", err)
// 		}

// 		fmt.Printf("[DELAYED MESSAGE] Delayed successfully - JetStream will redeliver in %v\n", delay)
// 		return nil
// 	}

// 	// Ready to execute
// 	fmt.Printf("[DELAYED MESSAGE] Executing: %s -> %s", roundID, originalSubject)
// 	if timeSinceExecute > 0 {
// 		fmt.Printf(" (overdue by %v)", timeSinceExecute)
// 	}
// 	fmt.Printf("\n")

// 	msgLogger := routineLogger.With(
// 		attr.String("round_id", roundID),
// 		attr.String("original_subject", originalSubject),
// 		attr.Time("execute_at", executeAtUTC),
// 		attr.Duration("overdue_by", timeSinceExecute),
// 	)

// 	// Add execution metadata
// 	wmMsg.Metadata.Set("Delayed-Executed-At", now.Format(time.RFC3339))
// 	if timeSinceExecute > 0 {
// 		wmMsg.Metadata.Set("Delayed-Overdue-By", timeSinceExecute.String())
// 	}

// 	// Ensure correlation ID exists
// 	if wmMsg.Metadata.Get("correlation_id") == "" {
// 		wmMsg.Metadata.Set("correlation_id", wmMsg.UUID)
// 	}

// 	// Republish to original subject
// 	if err := eb.publisher.Publish(originalSubject, wmMsg); err != nil {
// 		msgLogger.Error("Failed to republish delayed message", attr.Error(err))
// 		fmt.Printf("[DELAYED MESSAGE] Failed to republish: %v\n", err)

// 		// Nack for retry
// 		if nakErr := jetStreamMsg.Nak(); nakErr != nil {
// 			msgLogger.Error("Failed to nack message for retry", attr.Error(nakErr))
// 		}
// 		return fmt.Errorf("republish failed: %w", err)
// 	}

// 	// Success - ACK the message
// 	if err := jetStreamMsg.Ack(); err != nil {
// 		msgLogger.Error("Failed to ack processed message", attr.Error(err))
// 		return fmt.Errorf("ack failed: %w", err)
// 	}

// 	msgLogger.Info("Delayed message executed successfully")
// 	fmt.Printf("[DELAYED MESSAGE] Successfully executed and ACKed\n")
// 	return nil
// }

// // RecoverScheduledRounds handles recovery after restart - now just ensures global processor is running
// func (eb *eventBus) RecoverScheduledRounds(ctx context.Context) {
// 	ctxLogger := eb.logger.With(attr.String("operation", "recover_scheduled_rounds"))
// 	ctxLogger.Info("Recovery called - global processor handles all delayed messages automatically")

// 	// The global processor with DeliverAllPolicy automatically handles recovery
// 	// by processing all existing delayed messages in the stream
// 	// No additional recovery logic needed
// }

// // Add this method for debugging
// func (eb *eventBus) DebugDelayedMessages(ctx context.Context) error {
// 	stream, err := eb.js.Stream(ctx, DelayedMessagesStream)
// 	if err != nil {
// 		return fmt.Errorf("failed to get stream: %w", err)
// 	}

// 	info, err := stream.Info(ctx)
// 	if err != nil {
// 		return fmt.Errorf("failed to get stream info: %w", err)
// 	}

// 	fmt.Printf("=== DELAYED STREAM DEBUG ===\n")
// 	fmt.Printf("Total messages: %d\n", info.State.Msgs)
// 	fmt.Printf("Total bytes: %d\n", info.State.Bytes)
// 	fmt.Printf("First sequence: %d\n", info.State.FirstSeq)
// 	fmt.Printf("Last sequence: %d\n", info.State.LastSeq)

// 	// Check the global consumer
// 	cons, err := eb.js.Consumer(ctx, DelayedMessagesStream, GlobalDelayedConsumerName)
// 	if err != nil {
// 		fmt.Printf("Global consumer not found: %v\n", err)
// 		return err
// 	}

// 	consInfo, err := cons.Info(ctx)
// 	if err != nil {
// 		fmt.Printf("Failed to get consumer info: %v\n", err)
// 		return err
// 	}

// 	fmt.Printf("\n=== GLOBAL CONSUMER DEBUG ===\n")
// 	fmt.Printf("Consumer name: %s\n", consInfo.Name)
// 	fmt.Printf("Delivered messages: %d\n", consInfo.Delivered.Consumer)
// 	fmt.Printf("Acknowledged: %d\n", consInfo.AckFloor.Consumer)
// 	fmt.Printf("Pending: %d\n", consInfo.NumPending)
// 	fmt.Printf("Redelivered: %d\n", consInfo.NumRedelivered)
// 	fmt.Printf("Waiting: %d\n", consInfo.NumWaiting)
// 	fmt.Printf("Active: %v\n", consInfo.PushBound)

// 	return nil
// }

// // Updated CancelScheduledMessage with more thorough cleanup
// func (eb *eventBus) CancelScheduledMessage(ctx context.Context, roundID sharedtypes.RoundID) error {
// 	ctxLogger := eb.logger.With(
// 		attr.String("operation", "cancel_scheduled_message"),
// 		attr.RoundID("round_id", roundID),
// 		attr.String("stream", DelayedMessagesStream),
// 	)

// 	fmt.Printf("[CANCEL] Starting cancellation for round: %s\n", roundID)
// 	ctxLogger.Info("Cancelling scheduled messages for round")

// 	// Get the stream
// 	stream, err := eb.js.Stream(ctx, DelayedMessagesStream)
// 	if err != nil {
// 		ctxLogger.Error("Failed to access stream", attr.Error(err))
// 		return fmt.Errorf("failed to access stream %s: %w", DelayedMessagesStream, err)
// 	}

// 	// Subject pattern for this round's messages
// 	subjectPattern := fmt.Sprintf("%s.%v", DelayedMessagesSubject, roundID)
// 	fmt.Printf("[CANCEL] Purging messages with subject pattern: %s\n", subjectPattern)

// 	// Purge messages for this specific round from the stream
// 	err = stream.Purge(ctx, jetstream.WithPurgeSubject(subjectPattern))
// 	if err != nil {
// 		ctxLogger.Error("Failed to purge messages from stream", attr.Error(err))
// 		return fmt.Errorf("failed to purge messages for round %v: %w", roundID, err)
// 	}

// 	fmt.Printf("[CANCEL] Successfully purged messages from stream\n")
// 	ctxLogger.Info("Messages purged from stream",
// 		attr.String("subject_pattern", subjectPattern))

// 	// The consumer will continue processing other messages normally
// 	// Purged messages will no longer be delivered

// 	ctxLogger.Info("Cancellation completed successfully")
// 	return nil
// }

// // ScheduleDelayedMessage schedules a message to be delivered at a future time
// func (eb *eventBus) ScheduleDelayedMessage(ctx context.Context, originalSubject string, roundID sharedtypes.RoundID, scheduledTime sharedtypes.StartTime, payload []byte, additionalMetadata map[string]string) error {
// 	ctxLogger := eb.logger.With(
// 		attr.String("operation", "schedule_delayed_message"),
// 		attr.String("original_subject", originalSubject),
// 		attr.RoundID("round_id", roundID),
// 		attr.Time("scheduled_time", scheduledTime.AsTime()),
// 	)

// 	ctxLogger.Info("Scheduling delayed message")

// 	// Validate timing with buffer
// 	now := time.Now().UTC()
// 	scheduledTimeUTC := scheduledTime.AsTime().UTC()
// 	minBuffer := 5 * time.Second

// 	if !scheduledTimeUTC.After(now.Add(minBuffer)) {
// 		ctxLogger.Warn("Scheduled time too close to current time",
// 			attr.Time("current_time", now),
// 			attr.Time("scheduled_time", scheduledTimeUTC),
// 			attr.Duration("minimum_buffer", minBuffer),
// 		)
// 		return fmt.Errorf("scheduled time must be at least %v in the future", minBuffer)
// 	}

// 	// Create message with metadata
// 	msg := message.NewMessage(watermill.NewUUID(), payload)
// 	msg.Metadata.Set("Original-Subject", originalSubject)
// 	msg.Metadata.Set("Round-ID", roundID.String())
// 	msg.Metadata.Set("Execute-At", scheduledTimeUTC.Format(time.RFC3339))

// 	// Add additional metadata
// 	for key, value := range additionalMetadata {
// 		msg.Metadata.Set(key, value)
// 	}

// 	// Publish to delayed messages subject
// 	delayedSubject := fmt.Sprintf("%s.%v", DelayedMessagesSubject, roundID)

// 	ctxLogger.Info("Publishing delayed message",
// 		attr.String("discord_message_id", msg.UUID),
// 		attr.String("delayed_subject", delayedSubject),
// 	)

// 	if err := eb.Publish(delayedSubject, msg); err != nil {
// 		ctxLogger.Error("Failed to publish delayed message", attr.Error(err))
// 		return fmt.Errorf("failed to publish delayed message: %w", err)
// 	}

// 	ctxLogger.Info("Successfully scheduled delayed message",
// 		attr.String("discord_message_id", msg.UUID),
// 		attr.String("delayed_subject", delayedSubject),
// 	)

// 	return nil
// }
