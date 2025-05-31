package scoremetrics

import "go.opentelemetry.io/otel/metric"

// NewScoreMetrics creates a new ScoreMetrics implementation using OpenTelemetry.
// It requires an OTEL Meter instance.
func NewScoreMetrics(meter metric.Meter, prefix string) (ScoreMetrics, error) {
	// Helper function to create metric names with prefix and subsystem
	metricName := func(name string) string {
		if prefix != "" {
			return prefix + "_score_" + name
		}
		return "score_" + name
	}

	var err error
	m := &scoreMetrics{meter: meter}

	// --- Score Processing Metrics ---
	m.scoreProcessingAttemptCounter, err = meter.Int64Counter(
		metricName("processing_attempts_total"),
		metric.WithDescription("Number of score processing attempts"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.scoreProcessingSuccessCounter, err = meter.Int64Counter(
		metricName("processing_success_total"),
		metric.WithDescription("Number of successful score processing attempts"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.scoreProcessingFailureCounter, err = meter.Int64Counter(
		metricName("processing_failure_total"),
		metric.WithDescription("Number of failed score processing attempts"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.scoreProcessingDuration, err = meter.Float64Histogram(
		metricName("processing_duration_seconds"),
		metric.WithDescription("Duration of score processing in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, err
	}

	// --- Score Correction Metrics ---
	m.scoreCorrectionAttemptCounter, err = meter.Int64Counter(
		metricName("correction_attempts_total"),
		metric.WithDescription("Number of score correction attempts"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.scoreCorrectionSuccessCounter, err = meter.Int64Counter(
		metricName("correction_success_total"),
		metric.WithDescription("Number of successful score corrections"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.scoreCorrectionFailureCounter, err = meter.Int64Counter(
		metricName("correction_failure_total"),
		metric.WithDescription("Number of failed score corrections"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.scoreCorrectionDuration, err = meter.Float64Histogram(
		metricName("correction_duration_seconds"),
		metric.WithDescription("Duration of score corrections in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, err
	}

	// --- Generic Operation Metrics ---
	m.scoreOperationAttemptCounter, err = meter.Int64Counter(
		metricName("operation_attempts_total"),
		metric.WithDescription("Number of operation attempts"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.scoreOperationSuccessCounter, err = meter.Int64Counter(
		metricName("operation_success_total"),
		metric.WithDescription("Number of successful operations"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.scoreOperationFailureCounter, err = meter.Int64Counter(
		metricName("operation_failure_total"),
		metric.WithDescription("Number of failed operations"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.scoreOperationDuration, err = meter.Float64Histogram(
		metricName("operation_duration_seconds"),
		metric.WithDescription("Duration of operations in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, err
	}

	// --- Leaderboard Metrics ---
	m.leaderboardUpdateAttemptCounter, err = meter.Int64Counter(
		metricName("leaderboard_update_attempts_total"),
		metric.WithDescription("Number of leaderboard update attempts"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.leaderboardUpdateSuccessCounter, err = meter.Int64Counter(
		metricName("leaderboard_update_success_total"),
		metric.WithDescription("Number of successful leaderboard updates"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.leaderboardUpdateFailureCounter, err = meter.Int64Counter(
		metricName("leaderboard_update_failure_total"),
		metric.WithDescription("Number of failed leaderboard updates"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.leaderboardUpdateDuration, err = meter.Float64Histogram(
		metricName("leaderboard_update_duration_seconds"),
		metric.WithDescription("Duration of leaderboard updates in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, err
	}

	// --- Database Metrics ---
	m.dbQueryDuration, err = meter.Float64Histogram(
		metricName("db_query_duration_seconds"),
		metric.WithDescription("Duration of database queries in seconds"),
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

	// --- Detailed Score Processing Metrics ---
	m.scoreUpdateAttemptCounter, err = meter.Int64Counter(
		metricName("score_update_attempts_total"),
		metric.WithDescription("Number of score update attempts"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.roundScoresProcessingAttemptCounter, err = meter.Int64Counter(
		metricName("round_scores_processing_attempts_total"),
		metric.WithDescription("Number of round scores processing attempts"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.scoresProcessedCounter, err = meter.Int64Counter(
		metricName("scores_processed_total"),
		metric.WithDescription("Number of scores processed per round"),
		metric.WithUnit("1"), // Unit is "scores" but Dimensionless is common
	)
	if err != nil {
		return nil, err
	}
	m.taggedPlayersProcessedCounter, err = meter.Int64Counter(
		metricName("tagged_players_processed_total"),
		metric.WithDescription("Number of tagged players processed per round"),
		metric.WithUnit("1"), // Unit is "players"
	)
	if err != nil {
		return nil, err
	}
	m.untaggedPlayersProcessedCounter, err = meter.Int64Counter(
		metricName("untagged_players_processed_total"),
		metric.WithDescription("Number of untagged players processed per round"),
		metric.WithUnit("1"), // Unit is "players"
	)
	if err != nil {
		return nil, err
	}
	m.scoreSortingDurationHistogram, err = meter.Float64Histogram(
		metricName("score_sorting_duration_seconds"),
		metric.WithDescription("Duration of score sorting operations in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, err
	}
	m.tagExtractionDurationHistogram, err = meter.Float64Histogram(
		metricName("tag_extraction_duration_seconds"),
		metric.WithDescription("Duration of tag extraction operations in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, err
	}

	// --- Player/Tag Specific Metrics ---
	m.playerScoreGauge, err = meter.Int64Histogram( // Changed to Int64Histogram
		metricName("player_score"),
		metric.WithDescription("Score achieved by a player in a round"),
		metric.WithUnit("1"), // Or appropriate score unit if defined in UCUM
	)
	if err != nil {
		return nil, err
	}
	m.playerTagGauge, err = meter.Int64Histogram(
		metricName("player_tag"),
		metric.WithDescription("Tag number associated with a player in a round"),
		metric.WithUnit("1"), // Tag numbers are unitless
	)
	if err != nil {
		return nil, err
	}
	m.tagPerformanceGauge, err = meter.Int64Histogram( // Changed to Int64Histogram
		metricName("tag_performance"),
		metric.WithDescription("Score associated with a specific tag in a round"),
		metric.WithUnit("1"), // Or appropriate score unit if defined in UCUM
	)
	if err != nil {
		return nil, err
	}
	m.tagMovementCounter, err = meter.Int64Counter(
		metricName("tag_movement_total"),
		metric.WithDescription("Number of times a tag has moved between players"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}
	m.untaggedPlayerCounter, err = meter.Int64Counter(
		metricName("untagged_player_total"),
		metric.WithDescription("Count of players without tags participating in rounds"),
		metric.WithUnit("1"), // Unit is "players"
	)
	if err != nil {
		return nil, err
	}

	// OTEL metrics don't need explicit registration like Prometheus client_golang.
	// The MeterProvider handles the export pipeline.

	return m, nil
}
