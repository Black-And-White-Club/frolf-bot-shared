// observability/otel/eventbusmetrics/noop.go
package eventbusmetrics

import (
	"context"
	"time"
)

// NoOpMetrics is a metrics collector that does nothing. Useful for unit tests.
type NoOpMetrics struct{}

// NewNoop returns a no-operation implementation of EventBusMetrics
func NewNoop() EventBusMetrics {
	return &NoOpMetrics{}
}

// RecordMessagePublish does nothing
func (n *NoOpMetrics) RecordMessagePublish(ctx context.Context, topic string) {
}

// RecordMessagePublishError does nothing
func (n *NoOpMetrics) RecordMessagePublishError(ctx context.Context, topic string) {
}

// RecordMessageProcess does nothing
func (n *NoOpMetrics) RecordMessageProcess(ctx context.Context, topic string, success bool) {
}

// RecordMessageProcessingTime does nothing
func (n *NoOpMetrics) RecordMessageProcessingTime(ctx context.Context, topic string, duration time.Duration) {
}

// RecordMessageSubscribe does nothing
func (n *NoOpMetrics) RecordMessageSubscribe(ctx context.Context, topic string) {
}

// RecordMessageSubscribeError does nothing
func (n *NoOpMetrics) RecordMessageSubscribeError(ctx context.Context, topic string) {
}
