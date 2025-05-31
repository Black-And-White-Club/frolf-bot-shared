// observability/opentelemetry/round/constructor.go
package roundmetrics

import (
	"go.opentelemetry.io/otel/metric"
)

// NewRoundMetrics creates a new RoundMetrics implementation using OpenTelemetry.
// It requires an OTEL Meter instance and a prefix for metric names.
func NewRoundMetrics(meter metric.Meter, prefix string) (RoundMetrics, error) {
	// Helper function to create metric names with prefix and subsystem
	metricName := func(name string) string {
		if prefix != "" {
			return prefix + "_round_" + name
		}
		return "round_" + name
	}

	var err error
	m := &roundMetrics{meter: meter}

	// --- Operation Metrics ---
	m.operationAttemptCounter, err = meter.Int64Counter(
		metricName("operation_attempts_total"),
		metric.WithDescription("Total number of round operations attempted"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.operationSuccessCounter, err = meter.Int64Counter(
		metricName("operation_successes_total"),
		metric.WithDescription("Total number of successful round operations"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.operationFailureCounter, err = meter.Int64Counter(
		metricName("operation_failures_total"),
		metric.WithDescription("Total number of failed round operations"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.operationDuration, err = meter.Float64Histogram(
		metricName("operation_duration_seconds"),
		metric.WithDescription("Duration of round operations in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, err
	}

	// --- Round Specific Metrics ---
	m.roundCreatedCounter, err = meter.Int64Counter(
		metricName("created_total"),
		metric.WithDescription("Total number of rounds created"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.roundParticipantAddedCounter, err = meter.Int64Counter(
		metricName("participants_added_total"),
		metric.WithDescription("Total number of participants added to rounds"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.roundFinalizedCounter, err = meter.Int64Counter(
		metricName("finalized_total"),
		metric.WithDescription("Total number of rounds finalized"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.roundCancelledCounter, err = meter.Int64Counter(
		metricName("cancelled_total"),
		metric.WithDescription("Total number of rounds cancelled"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.timeParsingErrorCounter, err = meter.Int64Counter(
		metricName("time_parsing_errors_total"),
		metric.WithDescription("Total number of time parsing errors"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.validationErrorCounter, err = meter.Int64Counter(
		metricName("validation_errors_total"),
		metric.WithDescription("Total number of validation errors"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.dbOperationErrorCounter, err = meter.Int64Counter(
		metricName("db_operation_errors_total"),
		metric.WithDescription("Total number of database operation errors"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	// --- Delete Specific Metrics ---
	m.roundDeleteAttemptCounter, err = meter.Int64Counter(
		metricName("delete_attempts_total"),
		metric.WithDescription("Total number of delete attempts for rounds"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.roundDeleteSuccessCounter, err = meter.Int64Counter(
		metricName("delete_successes_total"),
		metric.WithDescription("Total number of successful round deletions"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.roundDeleteFailureCounter, err = meter.Int64Counter(
		metricName("delete_failures_total"),
		metric.WithDescription("Total number of failed round deletions"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.dbOperationDurationHistogram, err = meter.Float64Histogram(
		metricName("db_operation_duration_seconds"),
		metric.WithDescription("Duration of database operations in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, err
	}

	// --- Validation and Time Parsing Successes ---
	m.validationSuccessCounter, err = meter.Int64Counter(
		metricName("validation_successes_total"),
		metric.WithDescription("Total number of successful validations"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.timeParsingSuccessCounter, err = meter.Int64Counter(
		metricName("time_parsing_successes_total"),
		metric.WithDescription("Total number of successful time parsing operations"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.dbOperationSuccessCounter, err = meter.Int64Counter(
		metricName("db_operation_successes_total"),
		metric.WithDescription("Total number of successful database operations"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	// --- Handler Metrics ---
	m.handlerAttemptCounter, err = meter.Int64Counter(
		metricName("handler_attempts_total"),
		metric.WithDescription("Total number of handler attempts"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.handlerDurationHistogram, err = meter.Float64Histogram(
		metricName("handler_duration_seconds"),
		metric.WithDescription("Duration of handlers in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, err
	}
	m.handlerFailureCounter, err = meter.Int64Counter(
		metricName("handler_failures_total"),
		metric.WithDescription("Total number of handler failures"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.handlerSuccessCounter, err = meter.Int64Counter(
		metricName("handler_successes_total"),
		metric.WithDescription("Total number of handler successes"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	return m, nil
}
