package clubmetrics

import (
	"context"
	"time"
)

// NoOpMetrics is a metrics collector that does nothing. Useful for unit tests.
type NoOpMetrics struct{}

func NewNoop() ClubMetrics {
	return &NoOpMetrics{}
}

func (n *NoOpMetrics) RecordOperationAttempt(ctx context.Context, operationName, serviceName string) {
}

func (n *NoOpMetrics) RecordOperationSuccess(ctx context.Context, operationName, serviceName string) {
}

func (n *NoOpMetrics) RecordOperationFailure(ctx context.Context, operationName, serviceName string) {
}

func (n *NoOpMetrics) RecordOperationDuration(ctx context.Context, operationName, serviceName string, duration time.Duration) {
}
