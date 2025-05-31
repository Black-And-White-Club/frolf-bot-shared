// observability/opentelemetry/leaderboard/constructor.go
package leaderboardmetrics

import (
	"go.opentelemetry.io/otel/metric"
)

// NewLeaderboardMetrics creates a new LeaderboardMetrics implementation using OpenTelemetry.
// It requires an OTEL Meter instance and a prefix for metric names.
func NewLeaderboardMetrics(meter metric.Meter, prefix string) (LeaderboardMetrics, error) {
	// Helper function to create metric names with prefix and subsystem
	metricName := func(name string) string {
		if prefix != "" {
			return prefix + "_leaderboard_" + name
		}
		return "leaderboard_" + name
	}

	var err error
	m := &leaderboardMetrics{meter: meter}

	// --- Leaderboard Update ---
	m.leaderboardUpdateCounter, err = meter.Int64Counter(
		metricName("update_total"),
		metric.WithDescription("Number of leaderboard updates"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	// --- Tag Assignment ---
	m.tagAssignmentCounter, err = meter.Int64Counter(
		metricName("tag_assignment_total"),
		metric.WithDescription("Number of tag assignments"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	// --- Tag Availability Check ---
	m.tagAvailabilityCounter, err = meter.Int64Counter(
		metricName("tag_availability_checks_total"),
		metric.WithDescription("Number of tag availability checks"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	// --- Operation Metrics ---
	m.operationAttemptCounter, err = meter.Int64Counter(
		metricName("operation_attempts_total"),
		metric.WithDescription("Number of operation attempts"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.operationSuccessCounter, err = meter.Int64Counter(
		metricName("operation_success_total"),
		metric.WithDescription("Number of successful operations"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.operationFailureCounter, err = meter.Int64Counter(
		metricName("operation_failure_total"),
		metric.WithDescription("Number of failed operations"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.operationDuration, err = meter.Float64Histogram(
		metricName("operation_duration_seconds"),
		metric.WithDescription("Duration of operations in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, err
	}

	// --- Service Metrics ---
	m.serviceAttemptCounter, err = meter.Int64Counter(
		metricName("service_attempts_total"),
		metric.WithDescription("Number of service attempts"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.serviceSuccessCounter, err = meter.Int64Counter(
		metricName("service_success_total"),
		metric.WithDescription("Number of successful service executions"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.serviceFailureCounter, err = meter.Int64Counter(
		metricName("service_failure_total"),
		metric.WithDescription("Number of failed service executions"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.serviceDuration, err = meter.Float64Histogram(
		metricName("service_duration_seconds"),
		metric.WithDescription("Duration of service executions in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, err
	}

	// --- Leaderboard Update Attempts ---
	m.leaderboardUpdateAttemptCounter, err = meter.Int64Counter(
		metricName("update_attempts_total"),
		metric.WithDescription("Number of leaderboard update attempts"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.leaderboardUpdateSuccessCounter, err = meter.Int64Counter(
		metricName("update_success_total"),
		metric.WithDescription("Number of successful leaderboard updates"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.leaderboardUpdateFailureCounter, err = meter.Int64Counter(
		metricName("update_failure_total"),
		metric.WithDescription("Number of failed leaderboard updates"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.leaderboardUpdateDuration, err = meter.Float64Histogram(
		metricName("update_duration_seconds"),
		metric.WithDescription("Duration of leaderboard updates in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, err
	}

	// --- Leaderboard Get Attempts ---
	m.leaderboardGetAttemptCounter, err = meter.Int64Counter(
		metricName("get_attempts_total"),
		metric.WithDescription("Number of leaderboard retrieval attempts"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.leaderboardGetSuccessCounter, err = meter.Int64Counter(
		metricName("get_success_total"),
		metric.WithDescription("Number of successful leaderboard retrievals"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.leaderboardGetFailureCounter, err = meter.Int64Counter(
		metricName("get_failure_total"),
		metric.WithDescription("Number of failed leaderboard retrievals"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.leaderboardGetDuration, err = meter.Float64Histogram(
		metricName("get_duration_seconds"),
		metric.WithDescription("Duration of leaderboard retrievals in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, err
	}

	// --- Tag Get Attempts ---
	m.tagGetAttemptCounter, err = meter.Int64Counter(
		metricName("tag_get_attempts_total"),
		metric.WithDescription("Number of tag retrieval attempts"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.tagGetSuccessCounter, err = meter.Int64Counter(
		metricName("tag_get_success_total"),
		metric.WithDescription("Number of successful tag retrievals"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.tagGetFailureCounter, err = meter.Int64Counter(
		metricName("tag_get_failure_total"),
		metric.WithDescription("Number of failed tag retrievals"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.tagGetDuration, err = meter.Float64Histogram(
		metricName("tag_get_duration_seconds"),
		metric.WithDescription("Duration of tag retrievals in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, err
	}

	// --- Tag Assignment Attempts ---
	m.tagAssignmentAttemptCounter, err = meter.Int64Counter(
		metricName("tag_assignment_attempts_total"),
		metric.WithDescription("Number of tag assignment attempts"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.tagAssignmentSuccessCounter, err = meter.Int64Counter(
		metricName("tag_assignment_success_total"),
		metric.WithDescription("Number of successful tag assignments"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.tagAssignmentFailureCounter, err = meter.Int64Counter(
		metricName("tag_assignment_failure_total"),
		metric.WithDescription("Number of failed tag assignments"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.tagAssignmentDuration, err = meter.Float64Histogram(
		metricName("tag_assignment_duration_seconds"),
		metric.WithDescription("Duration of tag assignments in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, err
	}

	// --- Tag Swap Attempts ---
	m.tagSwapAttemptCounter, err = meter.Int64Counter(
		metricName("tag_swap_attempts_total"),
		metric.WithDescription("Number of tag swap attempts"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.tagSwapSuccessCounter, err = meter.Int64Counter(
		metricName("tag_swap_success_total"),
		metric.WithDescription("Number of successful tag swaps"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.tagSwapFailureCounter, err = meter.Int64Counter(
		metricName("tag_swap_failure_total"),
		metric.WithDescription("Number of failed tag swaps"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	// --- Tag Assignment Update ---
	m.tagAssignmentUpdateCounter, err = meter.Int64Counter(
		metricName("tag_assignment_updates_total"),
		metric.WithDescription("Number of tag assignment updates"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	// --- New Tag Assignment ---
	m.newTagAssignmentCounter, err = meter.Int64Counter(
		metricName("new_tag_assignments_total"),
		metric.WithDescription("Number of new tag assignments"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	// --- Tag Removal ---
	m.tagRemovalCounter, err = meter.Int64Counter(
		metricName("tag_removals_total"),
		metric.WithDescription("Number of tag removals"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	// --- Handler Metrics ---
	m.handlerAttemptCounter, err = meter.Int64Counter(
		metricName("handler_attempts_total"),
		metric.WithDescription("Number of handler attempts"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.handlerSuccessCounter, err = meter.Int64Counter(
		metricName("handler_success_total"),
		metric.WithDescription("Number of successful handler executions"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.handlerFailureCounter, err = meter.Int64Counter(
		metricName("handler_failure_total"),
		metric.WithDescription("Number of failed handler executions"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.handlerDuration, err = meter.Float64Histogram(
		metricName("handler_duration_seconds"),
		metric.WithDescription("Duration of handler executions in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, err
	}

	return m, nil
}
