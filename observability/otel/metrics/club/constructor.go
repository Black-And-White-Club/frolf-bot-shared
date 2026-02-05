package clubmetrics

import (
	"go.opentelemetry.io/otel/metric"
)

// NewClubMetrics creates a new ClubMetrics implementation using OpenTelemetry.
func NewClubMetrics(meter metric.Meter, prefix string) (ClubMetrics, error) {
	metricName := func(name string) string {
		if prefix != "" {
			return prefix + "_club_" + name
		}
		return "club_" + name
	}

	var err error
	m := &clubMetrics{meter: meter}

	m.operationAttemptCounter, err = meter.Int64Counter(
		metricName("operation_attempts_total"),
		metric.WithDescription("Number of club operation attempts"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	m.operationSuccessCounter, err = meter.Int64Counter(
		metricName("operation_success_total"),
		metric.WithDescription("Number of successful club operations"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	m.operationFailureCounter, err = meter.Int64Counter(
		metricName("operation_failure_total"),
		metric.WithDescription("Number of failed club operations"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	m.operationDuration, err = meter.Float64Histogram(
		metricName("operation_duration_seconds"),
		metric.WithDescription("Duration of club operations in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, err
	}

	return m, nil
}
