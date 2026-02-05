package clubmetrics

import (
	"context"
	"time"
)

// ClubMetrics defines metrics specific to club operations.
type ClubMetrics interface {
	RecordOperationAttempt(ctx context.Context, operationName, serviceName string)
	RecordOperationSuccess(ctx context.Context, operationName, serviceName string)
	RecordOperationFailure(ctx context.Context, operationName, serviceName string)
	RecordOperationDuration(ctx context.Context, operationName, serviceName string, duration time.Duration)
}
