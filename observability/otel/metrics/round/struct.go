// observability/opentelemetry/round/struct.go
package roundmetrics

import (
	"go.opentelemetry.io/otel/metric"
)

// roundMetrics implements RoundMetrics using OpenTelemetry
type roundMetrics struct {
	meter metric.Meter // OTEL Meter

	// --- Operation Metrics ---
	operationAttemptCounter metric.Int64Counter
	operationSuccessCounter metric.Int64Counter
	operationFailureCounter metric.Int64Counter
	operationDuration       metric.Float64Histogram

	// --- Round Specific Metrics ---
	roundCreatedCounter          metric.Int64Counter
	roundParticipantAddedCounter metric.Int64Counter
	roundFinalizedCounter        metric.Int64Counter
	roundCancelledCounter        metric.Int64Counter
	timeParsingErrorCounter      metric.Int64Counter
	validationErrorCounter       metric.Int64Counter
	dbOperationErrorCounter      metric.Int64Counter

	// --- Delete Specific Metrics ---
	roundDeleteAttemptCounter    metric.Int64Counter
	roundDeleteSuccessCounter    metric.Int64Counter
	roundDeleteFailureCounter    metric.Int64Counter
	dbOperationDurationHistogram metric.Float64Histogram

	// --- Validation and Time Parsing Successes ---
	validationSuccessCounter  metric.Int64Counter
	timeParsingSuccessCounter metric.Int64Counter
	dbOperationSuccessCounter metric.Int64Counter

	// --- Handler Metrics ---
	handlerAttemptCounter    metric.Int64Counter
	handlerDurationHistogram metric.Float64Histogram
	handlerFailureCounter    metric.Int64Counter
	handlerSuccessCounter    metric.Int64Counter
}
