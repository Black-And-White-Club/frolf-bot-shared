// observability/opentelemetry/leaderboard/struct.go
package leaderboardmetrics

import (
	"go.opentelemetry.io/otel/metric"
)

// leaderboardMetrics implements LeaderboardMetrics using OpenTelemetry
type leaderboardMetrics struct {
	meter metric.Meter // OTEL Meter

	leaderboardUpdateCounter        metric.Int64Counter
	tagAssignmentCounter            metric.Int64Counter
	tagAvailabilityCounter          metric.Int64Counter
	operationAttemptCounter         metric.Int64Counter
	operationSuccessCounter         metric.Int64Counter
	operationFailureCounter         metric.Int64Counter
	operationDuration               metric.Float64Histogram
	serviceAttemptCounter           metric.Int64Counter
	serviceSuccessCounter           metric.Int64Counter
	serviceFailureCounter           metric.Int64Counter
	serviceDuration                 metric.Float64Histogram
	leaderboardUpdateAttemptCounter metric.Int64Counter
	leaderboardUpdateSuccessCounter metric.Int64Counter
	leaderboardUpdateFailureCounter metric.Int64Counter
	leaderboardUpdateDuration       metric.Float64Histogram
	leaderboardGetAttemptCounter    metric.Int64Counter
	leaderboardGetSuccessCounter    metric.Int64Counter
	leaderboardGetFailureCounter    metric.Int64Counter
	leaderboardGetDuration          metric.Float64Histogram
	tagGetAttemptCounter            metric.Int64Counter
	tagGetSuccessCounter            metric.Int64Counter
	tagGetFailureCounter            metric.Int64Counter
	tagGetDuration                  metric.Float64Histogram
	tagAssignmentAttemptCounter     metric.Int64Counter
	tagAssignmentSuccessCounter     metric.Int64Counter
	tagAssignmentFailureCounter     metric.Int64Counter
	tagAssignmentDuration           metric.Float64Histogram
	tagSwapAttemptCounter           metric.Int64Counter
	tagSwapSuccessCounter           metric.Int64Counter
	tagSwapFailureCounter           metric.Int64Counter
	tagAssignmentUpdateCounter      metric.Int64Counter
	newTagAssignmentCounter         metric.Int64Counter
	tagRemovalCounter               metric.Int64Counter
	handlerAttemptCounter           metric.Int64Counter
	handlerSuccessCounter           metric.Int64Counter
	handlerFailureCounter           metric.Int64Counter
	handlerDuration                 metric.Float64Histogram
}
