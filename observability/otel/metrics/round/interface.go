package roundmetrics

import (
	"context"
	"time"
)

// RoundMetrics defines metrics for round operations
type RoundMetrics interface {
	// Operation metrics (used by serviceWrapper)
	RecordOperationAttempt(ctx context.Context, operation, service string)
	RecordOperationDuration(ctx context.Context, operation, service string, duration time.Duration)
	RecordOperationFailure(ctx context.Context, operation, service string)
	RecordOperationSuccess(ctx context.Context, operation, service string)

	// Round specific metrics
	RecordRoundCreated(ctx context.Context, location string)
	RecordRoundParticipantAdded(ctx context.Context, location string)
	RecordRoundFinalized(ctx context.Context, location string)
	RecordRoundCancelled(ctx context.Context, location string)
	RecordTimeParsingError(ctx context.Context)
	RecordValidationError(ctx context.Context)
	RecordDBOperationError(ctx context.Context, operation string)

	// Delete specific metrics
	RecordRoundDeleteAttempt(ctx context.Context)
	RecordRoundDeleteSuccess(ctx context.Context)
	RecordRoundDeleteFailure(ctx context.Context)
	RecordDBOperationDuration(ctx context.Context, operation string, duration time.Duration)

	RecordValidationSuccess(ctx context.Context)
	RecordTimeParsingSuccess(ctx context.Context)
	RecordDBOperationSuccess(ctx context.Context, operation string)

	RecordHandlerAttempt(ctx context.Context, handlerName string)
	RecordHandlerDuration(ctx context.Context, handlerName string, duration time.Duration)
	RecordHandlerFailure(ctx context.Context, handlerName string)
	RecordHandlerSuccess(ctx context.Context, handlerName string)
}
