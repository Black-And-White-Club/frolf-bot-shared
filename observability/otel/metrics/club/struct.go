package clubmetrics

import (
	"go.opentelemetry.io/otel/metric"
)

// clubMetrics implements ClubMetrics using OpenTelemetry.
type clubMetrics struct {
	meter metric.Meter

	operationAttemptCounter metric.Int64Counter
	operationSuccessCounter metric.Int64Counter
	operationFailureCounter metric.Int64Counter
	operationDuration       metric.Float64Histogram
}
