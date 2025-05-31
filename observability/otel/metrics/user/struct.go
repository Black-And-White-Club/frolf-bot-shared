package usermetrics

import "go.opentelemetry.io/otel/metric"

// userMetrics implements UserMetrics using OpenTelemetry
type userMetrics struct {
	meter metric.Meter // OTEL Meter

	// --- User Creation Metrics ---
	userCreationAttemptCounter metric.Int64Counter
	userCreationSuccessCounter metric.Int64Counter
	userCreationFailureCounter metric.Int64Counter
	userCreationDuration       metric.Float64Histogram

	// --- User Retrieval Metrics ---
	userRetrievalAttemptCounter metric.Int64Counter
	userRetrievalSuccessCounter metric.Int64Counter
	userRetrievalFailureCounter metric.Int64Counter
	userRetrievalDuration       metric.Float64Histogram

	// --- User Role Retrieval Metrics ---
	userRoleRetrievalAttemptCounter metric.Int64Counter
	userRoleRetrievalSuccessCounter metric.Int64Counter
	userRoleRetrievalFailureCounter metric.Int64Counter
	userRoleRetrievalDuration       metric.Float64Histogram

	// --- Tag Availability Check ---
	tagAvailableCounter   metric.Int64Counter
	tagUnavailableCounter metric.Int64Counter

	// --- Permission Check Metrics ---
	permissionCheckAttemptCounter metric.Int64Counter
	permissionCheckSuccessCounter metric.Int64Counter
	permissionCheckFailureCounter metric.Int64Counter
	permissionCheckDuration       metric.Float64Histogram

	// --- Role Update Metrics ---
	roleUpdateAttemptCounter metric.Int64Counter
	roleUpdateSuccessCounter metric.Int64Counter
	roleUpdateFailureCounter metric.Int64Counter
	roleUpdateDuration       metric.Float64Histogram

	// --- Database Metrics ---
	dbQueryDuration metric.Float64Histogram

	// --- User Creation by Tag ---
	userCreationByTagCounter metric.Int64Counter

	// --- Generic Operation Metrics ---
	userOperationAttemptCounter metric.Int64Counter
	userOperationSuccessCounter metric.Int64Counter
	userOperationFailureCounter metric.Int64Counter
	userOperationDuration       metric.Float64Histogram

	// --- Role Changed Metric ---
	roleChangedTotal metric.Int64Counter

	// --- Handler Metrics ---
	handlerAttemptCounter metric.Int64Counter
	handlerSuccessCounter metric.Int64Counter
	handlerFailureCounter metric.Int64Counter
	handlerDuration       metric.Float64Histogram
}
