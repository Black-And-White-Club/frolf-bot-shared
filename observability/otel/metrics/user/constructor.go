package usermetrics

import "go.opentelemetry.io/otel/metric"

// NewUserMetrics creates a new UserMetrics implementation using OpenTelemetry.
// It requires an OTEL Meter instance and a prefix for metric names.
func NewUserMetrics(meter metric.Meter, prefix string) (UserMetrics, error) {
	// Helper function to create metric names with prefix and subsystem
	metricName := func(name string) string {
		if prefix != "" {
			return prefix + "_user_" + name
		}
		return "user_" + name
	}

	var err error
	m := &userMetrics{meter: meter}

	// --- User Creation Metrics ---
	m.userCreationAttemptCounter, err = meter.Int64Counter(
		metricName("creation_attempts_total"),
		metric.WithDescription("Number of user creation attempts"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.userCreationSuccessCounter, err = meter.Int64Counter(
		metricName("creation_success_total"),
		metric.WithDescription("Number of successful user creations"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.userCreationFailureCounter, err = meter.Int64Counter(
		metricName("creation_failure_total"),
		metric.WithDescription("Number of failed user creations"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.userCreationDuration, err = meter.Float64Histogram(
		metricName("creation_duration_seconds"),
		metric.WithDescription("Duration of user creation in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, err
	}

	// --- User Retrieval Metrics ---
	m.userRetrievalAttemptCounter, err = meter.Int64Counter(
		metricName("retrieval_attempts_total"),
		metric.WithDescription("Number of user retrieval attempts"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.userRetrievalSuccessCounter, err = meter.Int64Counter(
		metricName("retrieval_success_total"),
		metric.WithDescription("Number of successful user retrievals"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.userRetrievalFailureCounter, err = meter.Int64Counter(
		metricName("retrieval_failure_total"),
		metric.WithDescription("Number of failed user retrievals"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.userRetrievalDuration, err = meter.Float64Histogram(
		metricName("retrieval_duration_seconds"),
		metric.WithDescription("Duration of user retrieval in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, err
	}

	// --- User Role Retrieval Metrics ---
	m.userRoleRetrievalAttemptCounter, err = meter.Int64Counter(
		metricName("role_retrieval_attempts_total"),
		metric.WithDescription("Number of user role retrieval attempts"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.userRoleRetrievalSuccessCounter, err = meter.Int64Counter(
		metricName("role_retrieval_success_total"),
		metric.WithDescription("Number of successful user role retrievals"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.userRoleRetrievalFailureCounter, err = meter.Int64Counter(
		metricName("role_retrieval_failure_total"),
		metric.WithDescription("Number of failed user role retrievals"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.userRoleRetrievalDuration, err = meter.Float64Histogram(
		metricName("role_retrieval_duration_seconds"),
		metric.WithDescription("Duration of user role retrieval in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, err
	}

	// --- Tag Availability Check ---
	m.tagAvailableCounter, err = meter.Int64Counter(
		metricName("tag_available_events_total"),
		metric.WithDescription("Number of times a tag became available"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.tagUnavailableCounter, err = meter.Int64Counter(
		metricName("tag_unavailable_events_total"),
		metric.WithDescription("Number of times a tag became unavailable"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	// --- Permission Check Metrics ---
	m.permissionCheckAttemptCounter, err = meter.Int64Counter(
		metricName("permission_check_attempts_total"),
		metric.WithDescription("Number of permission check attempts"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.permissionCheckSuccessCounter, err = meter.Int64Counter(
		metricName("permission_check_success_total"),
		metric.WithDescription("Number of successful permission checks"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.permissionCheckFailureCounter, err = meter.Int64Counter(
		metricName("permission_check_failure_total"),
		metric.WithDescription("Number of failed permission checks"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.permissionCheckDuration, err = meter.Float64Histogram(
		metricName("permission_check_duration_seconds"),
		metric.WithDescription("Duration of permission checks in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, err
	}

	// --- Role Changed Metric ---
	m.roleChangedTotal, err = meter.Int64Counter(
		metricName("role_changed_total"),
		metric.WithDescription("Total number of successful user role changes"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	// --- Role Update Metrics ---
	m.roleUpdateAttemptCounter, err = meter.Int64Counter(
		metricName("role_update_attempts_total"),
		metric.WithDescription("Number of role update attempts"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.roleUpdateSuccessCounter, err = meter.Int64Counter(
		metricName("role_update_success_total"),
		metric.WithDescription("Number of successful role updates"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.roleUpdateFailureCounter, err = meter.Int64Counter(
		metricName("role_update_failure_total"),
		metric.WithDescription("Number of failed role updates"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.roleUpdateDuration, err = meter.Float64Histogram(
		metricName("role_update_duration_seconds"),
		metric.WithDescription("Duration of role updates in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, err
	}

	// --- Database Metrics ---
	m.dbQueryDuration, err = meter.Float64Histogram(
		metricName("db_query_duration_seconds"),
		metric.WithDescription("Duration of database queries related to users in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, err
	}

	// --- User Creation by Tag ---
	m.userCreationByTagCounter, err = meter.Int64Counter(
		metricName("creation_by_tag_total"),
		metric.WithDescription("Number of users created with a specific tag"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	// --- Generic Operation Metrics ---
	m.userOperationAttemptCounter, err = meter.Int64Counter(
		metricName("operation_attempts_total"),
		metric.WithDescription("Number of user-related operation attempts"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.userOperationSuccessCounter, err = meter.Int64Counter(
		metricName("operation_success_total"),
		metric.WithDescription("Number of successful user-related operations"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.userOperationFailureCounter, err = meter.Int64Counter(
		metricName("operation_failure_total"),
		metric.WithDescription("Number of failed user-related operations"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.userOperationDuration, err = meter.Float64Histogram(
		metricName("operation_duration_seconds"),
		metric.WithDescription("Duration of user-related operations in seconds"),
		metric.WithUnit("s"),
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
		metric.WithDescription("Number of successful handler attempts"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.handlerFailureCounter, err = meter.Int64Counter(
		metricName("handler_failure_total"),
		metric.WithDescription("Number of failed handler attempts"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.handlerDuration, err = meter.Float64Histogram(
		metricName("handler_duration_seconds"),
		metric.WithDescription("Duration of handler execution in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, err
	}

	return m, nil
}
