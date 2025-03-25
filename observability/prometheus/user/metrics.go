// observability/prometheus/user/metrics.go
package usermetrics

import (
	"fmt"

	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
	"github.com/prometheus/client_golang/prometheus"
)

// UserMetrics defines metrics specific to user operations
type UserMetrics interface {
	RecordUserCreation(userType string, source string, status string)
	RecordUserRetrieval(success bool, userID sharedtypes.DiscordID)
	RecordUserRoleRetrieval(success bool, userID sharedtypes.DiscordID)
	RecordUserRetrievalDuration(duration float64)
	RecordTagAvailabilityCheck(available bool, tag int)
	RecordPermissionCheck(role sharedtypes.UserRoleEnum, allowed bool, action string, resource string)
	RecordRoleUpdate(oldRole, newRole sharedtypes.UserRoleEnum, context string, userID sharedtypes.DiscordID)
	UserCreationDuration(duration float64)
	DBQueryDuration(duration float64)
	UserCreationByTag(tag int)
	RecordOperationAttempt(operationName string, userID sharedtypes.DiscordID)
	RecordOperationSuccess(operationName string, userID sharedtypes.DiscordID)
	RecordOperationFailure(operationName string, userID sharedtypes.DiscordID)
	RecordOperationDuration(operationName string, duration float64)
	RecordUserRoleUpdateAttempt(userID sharedtypes.DiscordID, newRole string)
	RecordUserRoleUpdateSuccess(userID sharedtypes.DiscordID, newRole string)
	RecordUserRoleUpdateFailure(userID sharedtypes.DiscordID, newRole string)
	RecordHandlerAttempt(handlerName string)
	RecordHandlerDuration(handlerName string, duration float64)
	RecordHandlerFailure(handlerName string)
	RecordHandlerSuccess(handlerName string)
}

// userMetrics implements UserMetrics
type userMetrics struct {
	userCreationCounter          *prometheus.CounterVec
	userRetrievalCounter         *prometheus.CounterVec
	userRoleRetrievalCounter     *prometheus.CounterVec
	tagAvailabilityCounter       *prometheus.CounterVec
	permissionCheckCounter       *prometheus.CounterVec
	roleUpdateCounter            *prometheus.CounterVec
	userCreationDuration         *prometheus.HistogramVec
	userRetrievalDuration        *prometheus.HistogramVec
	dbQueryDuration              *prometheus.HistogramVec
	userCreationByTagCounter     *prometheus.CounterVec
	userOperationAttemptCounter  *prometheus.CounterVec
	userOperationSuccessCounter  *prometheus.CounterVec
	userOperationFailureCounter  *prometheus.CounterVec
	userOperationDuration        *prometheus.HistogramVec
	userRoleUpdateAttemptCounter *prometheus.CounterVec
	userRoleUpdateSuccessCounter *prometheus.CounterVec
	userRoleUpdateFailureCounter *prometheus.CounterVec
	handlerAttemptCounter        *prometheus.CounterVec
	handlerSuccessCounter        *prometheus.CounterVec
	handlerFailureCounter        *prometheus.CounterVec
	handlerDuration              *prometheus.HistogramVec
}

// NewUserMetrics creates a new UserMetrics implementation
func NewUserMetrics(registry *prometheus.Registry, prefix string) UserMetrics {
	userCreationCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: prefix,
			Subsystem: "user",
			Name:      "creation_total",
			Help:      "Number of users created, partitioned by user type, source, and status",
		},
		[]string{"user_type", "source", "status"},
	)

	userRetrievalCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: prefix,
			Subsystem: "user",
			Name:      "retrieval_total",
			Help:      "Number of user retrievals, partitioned by success and user ID",
		},
		[]string{"success", "user_id"},
	)

	userRoleRetrievalCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: prefix,
			Subsystem: "user",
			Name:      "role_retrieval_total",
			Help:      "Number of user role retrievals, partitioned by success and user ID",
		},
		[]string{"success", "user_id"},
	)

	tagAvailabilityCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: prefix,
			Subsystem: "user",
			Name:      "tag_availability_checks_total",
			Help:      "Number of tag availability checks, partitioned by availability result and tag type",
		},
		[]string{"available", "tag_type"},
	)

	permissionCheckCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: prefix,
			Subsystem: "user",
			Name:      "permission_checks_total",
			Help:      "Number of permission checks, partitioned by role, result, action, and resource",
		},
		[]string{"role", "allowed", "action", "resource"},
	)

	roleUpdateCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: prefix,
			Subsystem: "user",
			Name:      "role_updates_total",
			Help:      "Number of role updates, partitioned by old role, new role, context, and user ID",
		},
		[]string{"old_role", "new_role", "context", "user_id"},
	)

	userCreationDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: prefix,
			Subsystem: "user",
			Name:      "creation_duration_seconds",
			Help:      "Duration of user creation in seconds",
		},
		[]string{},
	)

	userRetrievalDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: prefix,
			Subsystem: "user",
			Name:      "retrieval_duration_seconds",
			Help:      "Duration of user retrieval in seconds",
		},
		[]string{},
	)

	dbQueryDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: prefix,
			Subsystem: "user",
			Name:      "db_query_duration_seconds",
			Help:      "Duration of database queries in seconds",
		},
		[]string{},
	)

	userCreationByTagCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: prefix,
			Subsystem: "user",
			Name:      "creation_by_tag_total",
			Help:      "Number of users created, partitioned by tag",
		},
		[]string{"tag"},
	)

	userOperationAttemptCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: prefix,
			Subsystem: "user",
			Name:      "operation_attempts_total",
			Help:      "Number of operation attempts, partitioned by operation name and user ID",
		},
		[]string{"operation", "user_id"},
	)

	userOperationSuccessCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: prefix,
			Subsystem: "user",
			Name:      "operation_success_total",
			Help:      "Number of successful operations, partitioned by operation name and user ID",
		},
		[]string{"operation", "user_id"},
	)

	userOperationFailureCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: prefix,
			Subsystem: "user",
			Name:      "operation_failure_total",
			Help:      "Number of failed operations, partitioned by operation name and user ID",
		},
		[]string{"operation", "user_id"},
	)

	userOperationDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: prefix,
			Subsystem: "user",
			Name:      "operation_duration_seconds",
			Help:      "Duration of operations in seconds",
		},
		[]string{"operation"},
	)

	userRoleUpdateAttemptCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: prefix,
			Subsystem: "user",
			Name:      "role_update_attempts_total",
			Help:      "Number of user role update attempts",
		},
		[]string{"user_id", "new_role"},
	)

	userRoleUpdateSuccessCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: prefix,
			Subsystem: "user",
			Name:      "role_update_success_total",
			Help:      "Number of successful user role updates",
		},
		[]string{"user_id", "new_role"},
	)

	userRoleUpdateFailureCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: prefix,
			Subsystem: "user",
			Name:      "role_update_failure_total",
			Help:      "Number of failed user role updates",
		},
		[]string{"user_id", "new_role"},
	)
	handlerAttemptCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: prefix,
			Subsystem: "user",
			Name:      "handler_attempts_total",
			Help:      "Number of handler attempts, partitioned by handler name",
		},
		[]string{"handler"},
	)

	handlerSuccessCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: prefix,
			Subsystem: "user",
			Name:      "handler_success_total",
			Help:      "Number of successful handler executions, partitioned by handler name",
		},
		[]string{"handler"},
	)

	handlerFailureCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: prefix,
			Subsystem: "user",
			Name:      "handler_failure_total",
			Help:      "Number of failed handler executions, partitioned by handler name",
		},
		[]string{"handler"},
	)

	handlerDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: prefix,
			Subsystem: "user",
			Name:      "handler_duration_seconds",
			Help:      "Duration of handler executions in seconds",
		},
		[]string{"handler"},
	)

	registry.MustRegister(
		userCreationCounter,
		userRetrievalCounter,
		userRoleRetrievalCounter,
		tagAvailabilityCounter,
		permissionCheckCounter,
		roleUpdateCounter,
		userCreationDuration,
		userRetrievalDuration,
		dbQueryDuration,
		userCreationByTagCounter,
		userOperationAttemptCounter,
		userOperationSuccessCounter,
		userOperationFailureCounter,
		userOperationDuration,
		userRoleUpdateAttemptCounter,
		userRoleUpdateSuccessCounter,
		userRoleUpdateFailureCounter,
		handlerAttemptCounter,
		handlerSuccessCounter,
		handlerFailureCounter,
		handlerDuration,
	)

	return &userMetrics{
		userCreationCounter:          userCreationCounter,
		userRetrievalCounter:         userRetrievalCounter,
		userRoleRetrievalCounter:     userRoleRetrievalCounter,
		tagAvailabilityCounter:       tagAvailabilityCounter,
		permissionCheckCounter:       permissionCheckCounter,
		roleUpdateCounter:            roleUpdateCounter,
		userCreationDuration:         userCreationDuration,
		userRetrievalDuration:        userRetrievalDuration,
		dbQueryDuration:              dbQueryDuration,
		userCreationByTagCounter:     userCreationByTagCounter,
		userOperationAttemptCounter:  userOperationAttemptCounter,
		userOperationSuccessCounter:  userOperationSuccessCounter,
		userOperationFailureCounter:  userOperationFailureCounter,
		userOperationDuration:        userOperationDuration,
		userRoleUpdateAttemptCounter: userRoleUpdateAttemptCounter,
		userRoleUpdateSuccessCounter: userRoleUpdateSuccessCounter,
		userRoleUpdateFailureCounter: userRoleUpdateFailureCounter,
		handlerAttemptCounter:        handlerAttemptCounter,
		handlerSuccessCounter:        handlerSuccessCounter,
		handlerFailureCounter:        handlerFailureCounter,
		handlerDuration:              handlerDuration,
	}

}

// RecordUserRetrieval records a user retrieval attempt
func (m *userMetrics) RecordUserRetrieval(success bool, userID sharedtypes.DiscordID) {
	m.userRetrievalCounter.WithLabelValues(fmt.Sprintf("%t", success), string(userID)).Inc()
}

// RecordUserRoleRetrieval records a user role retrieval attempt
func (m *userMetrics) RecordUserRoleRetrieval(success bool, userID sharedtypes.DiscordID) {
	m.userRoleRetrievalCounter.WithLabelValues(fmt.Sprintf("%t", success), string(userID)).Inc()
}

// RecordUserRetrievalDuration records the duration of user retrieval
func (m *userMetrics) RecordUserRetrievalDuration(duration float64) {
	m.userRetrievalDuration.WithLabelValues().Observe(duration)
}

// RecordOperationAttempt records an operation attempt
func (m *userMetrics) RecordOperationAttempt(operationName string, userID sharedtypes.DiscordID) {
	m.userOperationAttemptCounter.WithLabelValues(operationName, string(userID)).Inc()
}

// RecordOperationSuccess records a successful operation
func (m *userMetrics) RecordOperationSuccess(operationName string, userID sharedtypes.DiscordID) {
	m.userOperationSuccessCounter.WithLabelValues(operationName, string(userID)).Inc()
}

// RecordOperationFailure records a failed operation
func (m *userMetrics) RecordOperationFailure(operationName string, userID sharedtypes.DiscordID) {
	m.userOperationFailureCounter.WithLabelValues(operationName, string(userID)).Inc()
}

// RecordOperationDuration records the duration of an operation
func (m *userMetrics) RecordOperationDuration(operationName string, duration float64) {
	m.userOperationDuration.WithLabelValues(operationName).Observe(duration)
}

// RecordUserCreation records a user creation event
func (m *userMetrics) RecordUserCreation(userType string, source string, status string) {
	m.userCreationCounter.WithLabelValues(userType, source, status).Inc()
}

// RecordDBQueryDuration records the duration of a DB query
func (m *userMetrics) DBQueryDuration(duration float64) {
	m.dbQueryDuration.WithLabelValues().Observe(duration)
}

// RecordPermissionCheck records a permission check
func (m *userMetrics) RecordPermissionCheck(role sharedtypes.UserRoleEnum, allowed bool, action string, resource string) {
	m.permissionCheckCounter.WithLabelValues(fmt.Sprintf("%v", role), fmt.Sprintf("%t", allowed), action, resource).Inc()
}

// RecordRoleUpdate records a role update event
func (m *userMetrics) RecordRoleUpdate(oldRole, newRole sharedtypes.UserRoleEnum, context string, userID sharedtypes.DiscordID) {
	m.roleUpdateCounter.WithLabelValues(fmt.Sprintf("%v", oldRole), fmt.Sprintf("%v", newRole), context, string(userID)).Inc()
}

// RecordTagAvailabilityCheck records a tag availability check
func (m *userMetrics) RecordTagAvailabilityCheck(available bool, tag int) {
	m.tagAvailabilityCounter.WithLabelValues(fmt.Sprintf("%t", available), fmt.Sprintf("%d", tag)).Inc()
}

// UserCreationByTag records user creation by tag
func (m *userMetrics) UserCreationByTag(tag int) {
	m.userCreationByTagCounter.WithLabelValues(fmt.Sprintf("%d", tag)).Inc()
}

// UserCreationDuration records the duration of user creation
func (m *userMetrics) UserCreationDuration(duration float64) {
	m.userCreationDuration.WithLabelValues().Observe(duration)
}

// RecordUserRoleUpdateAttempt records an attempt to update a user’s role
func (m *userMetrics) RecordUserRoleUpdateAttempt(userID sharedtypes.DiscordID, newRole string) {
	m.userRoleUpdateAttemptCounter.WithLabelValues(string(userID), newRole).Inc()
}

// RecordUserRoleUpdateSuccess records a successful user role update
func (m *userMetrics) RecordUserRoleUpdateSuccess(userID sharedtypes.DiscordID, newRole string) {
	m.userRoleUpdateSuccessCounter.WithLabelValues(string(userID), newRole).Inc()
}

// RecordUserRoleUpdateFailure records a failed user role update
func (m *userMetrics) RecordUserRoleUpdateFailure(userID sharedtypes.DiscordID, newRole string) {
	m.userRoleUpdateFailureCounter.WithLabelValues(string(userID), newRole).Inc()
}

// RecordHandlerAttempt records a handler attempt
func (m *userMetrics) RecordHandlerAttempt(handlerName string) {
	m.handlerAttemptCounter.WithLabelValues(handlerName).Inc()
}

// RecordHandlerSuccess records a successful handler execution
func (m *userMetrics) RecordHandlerSuccess(handlerName string) {
	m.handlerSuccessCounter.WithLabelValues(handlerName).Inc()
}

// RecordHandlerFailure records a failed handler execution
func (m *userMetrics) RecordHandlerFailure(handlerName string) {
	m.handlerFailureCounter.WithLabelValues(handlerName).Inc()
}

// RecordHandlerDuration records the duration of a handler execution
func (m *userMetrics) RecordHandlerDuration(handlerName string, duration float64) {
	m.handlerDuration.WithLabelValues(handlerName).Observe(duration)
}
