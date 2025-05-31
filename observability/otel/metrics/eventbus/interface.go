// observability/otel/eventbusmetrics/interface.go
package eventbusmetrics

import (
	"context"
	"time"
)

// EventBusMetrics defines metrics specific to event bus operations using OpenTelemetry
type EventBusMetrics interface {
	// Message publishing metrics
	RecordMessagePublish(ctx context.Context, topic string)
	RecordMessagePublishError(ctx context.Context, topic string)

	// Message processing metrics
	RecordMessageProcess(ctx context.Context, topic string, success bool)
	RecordMessageProcessingTime(ctx context.Context, topic string, duration time.Duration)

	// Subscription metrics
	RecordMessageSubscribe(ctx context.Context, topic string)
	RecordMessageSubscribeError(ctx context.Context, topic string)
}
