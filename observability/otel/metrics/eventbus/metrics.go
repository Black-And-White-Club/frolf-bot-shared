// observability/otel/eventbusmetrics/metrics.go
package eventbusmetrics

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/metric"
)

// RecordMessagePublish records a message publish event
func (m *eventBusMetrics) RecordMessagePublish(ctx context.Context, topic string) {
	m.messagePublishCounter.Add(ctx, 1, metric.WithAttributes(topicAttrs(topic)))
}

// RecordMessagePublishError records a message publish error event
func (m *eventBusMetrics) RecordMessagePublishError(ctx context.Context, topic string) {
	m.messagePublishErrorCounter.Add(ctx, 1, metric.WithAttributes(topicAttrs(topic)))
}

// RecordMessageProcess records a message processing event
func (m *eventBusMetrics) RecordMessageProcess(ctx context.Context, topic string, success bool) {
	m.messageProcessCounter.Add(ctx, 1, metric.WithAttributes(topicSuccessAttrs(topic, success)...))
}

// RecordMessageProcessingTime records the duration of message processing
func (m *eventBusMetrics) RecordMessageProcessingTime(ctx context.Context, topic string, duration time.Duration) {
	m.messageProcessingTimeHistogram.Record(ctx, duration.Seconds(), metric.WithAttributes(topicAttrs(topic)))
}

// RecordMessageSubscribe records a message subscription event
func (m *eventBusMetrics) RecordMessageSubscribe(ctx context.Context, topic string) {
	m.messageSubscribeCounter.Add(ctx, 1, metric.WithAttributes(topicAttrs(topic)))
}

// RecordMessageSubscribeError records a message subscription error event
func (m *eventBusMetrics) RecordMessageSubscribeError(ctx context.Context, topic string) {
	m.messageSubscribeErrorCounter.Add(ctx, 1, metric.WithAttributes(topicAttrs(topic)))
}
