package scoremetrics

import "go.opentelemetry.io/otel/metric"

// scoreMetrics implements ScoreMetrics using OpenTelemetry
type scoreMetrics struct {
	meter metric.Meter // OTEL Meter

	// --- Score Processing Metrics ---
	scoreProcessingAttemptCounter metric.Int64Counter
	scoreProcessingSuccessCounter metric.Int64Counter
	scoreProcessingFailureCounter metric.Int64Counter
	scoreProcessingDuration       metric.Float64Histogram // Use Float64 for duration in seconds

	// --- Score Correction Metrics ---
	scoreCorrectionAttemptCounter metric.Int64Counter
	scoreCorrectionSuccessCounter metric.Int64Counter
	scoreCorrectionFailureCounter metric.Int64Counter
	scoreCorrectionDuration       metric.Float64Histogram

	// --- Generic Operation Metrics ---
	scoreOperationAttemptCounter metric.Int64Counter
	scoreOperationSuccessCounter metric.Int64Counter
	scoreOperationFailureCounter metric.Int64Counter
	scoreOperationDuration       metric.Float64Histogram

	// --- Leaderboard Metrics ---
	leaderboardUpdateAttemptCounter metric.Int64Counter
	leaderboardUpdateSuccessCounter metric.Int64Counter
	leaderboardUpdateFailureCounter metric.Int64Counter
	leaderboardUpdateDuration       metric.Float64Histogram

	// --- Database Metrics ---
	dbQueryDuration metric.Float64Histogram

	// --- Handler Metrics ---
	handlerAttemptCounter metric.Int64Counter
	handlerSuccessCounter metric.Int64Counter
	handlerFailureCounter metric.Int64Counter
	handlerDuration       metric.Float64Histogram

	// --- Detailed Score Processing Metrics ---
	scoreUpdateAttemptCounter           metric.Int64Counter // Note: Success/failure could be an attribute
	roundScoresProcessingAttemptCounter metric.Int64Counter // Note: Success/failure could be an attribute
	scoresProcessedCounter              metric.Int64Counter
	taggedPlayersProcessedCounter       metric.Int64Counter
	untaggedPlayersProcessedCounter     metric.Int64Counter
	scoreSortingDurationHistogram       metric.Float64Histogram
	tagExtractionDurationHistogram      metric.Float64Histogram

	// --- Player/Tag Specific Metrics ---
	//  Use Int64Histogram for integer scores and tag numbers
	playerScoreGauge      metric.Int64Histogram
	playerTagGauge        metric.Int64Histogram
	tagPerformanceGauge   metric.Int64Histogram
	tagMovementCounter    metric.Int64Counter
	untaggedPlayerCounter metric.Int64Counter
}
