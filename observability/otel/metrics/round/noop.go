package roundmetrics

import (
	"context"
	"time"
)

// NoOpMetrics is a metrics collector that does nothing. Useful for unit tests.
type NoOpMetrics struct{}

func NewNoop() RoundMetrics {
	return &NoOpMetrics{}
}

func (n *NoOpMetrics) RecordOperationAttempt(ctx context.Context, operation, service string) {}
func (n *NoOpMetrics) RecordOperationDuration(ctx context.Context, operation, service string, duration time.Duration) {
}
func (n *NoOpMetrics) RecordOperationFailure(ctx context.Context, operation, service string) {}
func (n *NoOpMetrics) RecordOperationSuccess(ctx context.Context, operation, service string) {}

// Round specific metrics
func (n *NoOpMetrics) RecordRoundCreated(ctx context.Context, location string)          {}
func (n *NoOpMetrics) RecordRoundParticipantAdded(ctx context.Context, location string) {}
func (n *NoOpMetrics) RecordRoundFinalized(ctx context.Context, location string)        {}
func (n *NoOpMetrics) RecordRoundCancelled(ctx context.Context, location string)        {}
func (n *NoOpMetrics) RecordTimeParsingError(ctx context.Context)                       {}
func (n *NoOpMetrics) RecordValidationError(ctx context.Context)                        {}
func (n *NoOpMetrics) RecordDBOperationError(ctx context.Context, operation string)     {}

// Delete specific metrics
func (n *NoOpMetrics) RecordRoundDeleteAttempt(ctx context.Context) {}
func (n *NoOpMetrics) RecordRoundDeleteSuccess(ctx context.Context) {}
func (n *NoOpMetrics) RecordRoundDeleteFailure(ctx context.Context) {}
func (n *NoOpMetrics) RecordDBOperationDuration(ctx context.Context, operation string, duration time.Duration) {
}

func (n *NoOpMetrics) RecordValidationSuccess(ctx context.Context)                    {}
func (n *NoOpMetrics) RecordTimeParsingSuccess(ctx context.Context)                   {}
func (n *NoOpMetrics) RecordDBOperationSuccess(ctx context.Context, operation string) {}

func (n *NoOpMetrics) RecordHandlerAttempt(ctx context.Context, handlerName string) {}
func (n *NoOpMetrics) RecordHandlerDuration(ctx context.Context, handlerName string, duration time.Duration) {
}
func (n *NoOpMetrics) RecordHandlerFailure(ctx context.Context, handlerName string) {}
func (n *NoOpMetrics) RecordHandlerSuccess(ctx context.Context, handlerName string) {}
