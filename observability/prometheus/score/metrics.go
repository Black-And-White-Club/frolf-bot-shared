package scoremetrics

import (
	"strconv"

	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
	"github.com/prometheus/client_golang/prometheus"
)

func roundIDToString(roundID sharedtypes.RoundID) string {
	return strconv.FormatInt(int64(roundID), 10)
}

// ScoreMetrics defines metrics specific to score operations
type ScoreMetrics interface {
	RecordScoreProcessingAttempt(roundID sharedtypes.RoundID)
	RecordScoreProcessingSuccess(roundID sharedtypes.RoundID)
	RecordScoreProcessingFailure(roundID sharedtypes.RoundID)
	RecordScoreProcessingDuration(roundID sharedtypes.RoundID, duration float64)
	RecordScoreCorrectionAttempt(roundID sharedtypes.RoundID)
	RecordScoreCorrectionSuccess(roundID sharedtypes.RoundID)
	RecordScoreCorrectionFailure(roundID sharedtypes.RoundID)
	RecordScoreCorrectionDuration(roundID sharedtypes.RoundID, duration float64)
	RecordLeaderboardUpdateAttempt(roundID sharedtypes.RoundID)
	RecordLeaderboardUpdateSuccess(roundID sharedtypes.RoundID)
	RecordLeaderboardUpdateFailure(roundID sharedtypes.RoundID)
	RecordLeaderboardUpdateDuration(roundID sharedtypes.RoundID, duration float64)
	RecordDBQueryDuration(duration float64)
}

// scoreMetrics implements ScoreMetrics
type scoreMetrics struct {
	scoreProcessingAttemptCounter   *prometheus.CounterVec
	scoreProcessingSuccessCounter   *prometheus.CounterVec
	scoreProcessingFailureCounter   *prometheus.CounterVec
	scoreProcessingDuration         *prometheus.HistogramVec
	scoreCorrectionAttemptCounter   *prometheus.CounterVec
	scoreCorrectionSuccessCounter   *prometheus.CounterVec
	scoreCorrectionFailureCounter   *prometheus.CounterVec
	scoreCorrectionDuration         *prometheus.HistogramVec
	leaderboardUpdateAttemptCounter *prometheus.CounterVec
	leaderboardUpdateSuccessCounter *prometheus.CounterVec
	leaderboardUpdateFailureCounter *prometheus.CounterVec
	leaderboardUpdateDuration       *prometheus.HistogramVec
	dbQueryDuration                 *prometheus.HistogramVec
}

// NewScoreMetrics creates a new ScoreMetrics implementation
func NewScoreMetrics(registry *prometheus.Registry, prefix string) ScoreMetrics {
	scoreProcessingAttemptCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: prefix,
			Subsystem: "score",
			Name:      "processing_attempts_total",
			Help:      "Number of score processing attempts, partitioned by round ID",
		},
		[]string{"round_id"},
	)

	scoreProcessingSuccessCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: prefix,
			Subsystem: "score",
			Name:      "processing_success_total",
			Help:      "Number of successful score processing attempts, partitioned by round ID",
		},
		[]string{"round_id"},
	)

	scoreProcessingFailureCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: prefix,
			Subsystem: "score",
			Name:      "processing_failure_total",
			Help:      "Number of failed score processing attempts, partitioned by round ID",
		},
		[]string{"round_id"},
	)

	scoreProcessingDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: prefix,
			Subsystem: "score",
			Name:      "processing_duration_seconds",
			Help:      "Duration of score processing in seconds, partitioned by round ID",
		},
		[]string{"round_id"},
	)

	scoreCorrectionAttemptCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: prefix,
			Subsystem: "score",
			Name:      "correction_attempts_total",
			Help:      "Number of score correction attempts, partitioned by round ID",
		},
		[]string{"round_id"},
	)

	scoreCorrectionSuccessCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: prefix,
			Subsystem: "score",
			Name:      "correction_success_total",
			Help:      "Number of successful score corrections, partitioned by round ID",
		},
		[]string{"round_id"},
	)

	scoreCorrectionFailureCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: prefix,
			Subsystem: "score",
			Name:      "correction_failure_total",
			Help:      "Number of failed score corrections, partitioned by round ID",
		},
		[]string{"round_id"},
	)

	scoreCorrectionDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: prefix,
			Subsystem: "score",
			Name:      "correction_duration_seconds",
			Help:      "Duration of score corrections in seconds, partitioned by round ID",
		},
		[]string{"round_id"},
	)

	leaderboardUpdateAttemptCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: prefix,
			Subsystem: "score",
			Name:      "leaderboard_update_attempts_total",
			Help:      "Number of leaderboard update attempts, partitioned by round ID",
		},
		[]string{"round_id"},
	)

	leaderboardUpdateSuccessCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: prefix,
			Subsystem: "score",
			Name:      "leaderboard_update_success_total",
			Help:      "Number of successful leaderboard updates, partitioned by round ID",
		},
		[]string{"round_id"},
	)

	leaderboardUpdateFailureCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: prefix,
			Subsystem: "score",
			Name:      "leaderboard_update_failure_total",
			Help:      "Number of failed leaderboard updates, partitioned by round ID",
		},
		[]string{"round_id"},
	)

	leaderboardUpdateDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: prefix,
			Subsystem: "score",
			Name:      "leaderboard_update_duration_seconds",
			Help:      "Duration of leaderboard updates in seconds, partitioned by round ID",
		},
		[]string{"round_id"},
	)

	dbQueryDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: prefix,
			Subsystem: "score",
			Name:      "db_query_duration_seconds",
			Help:      "Duration of database queries in seconds",
		},
		[]string{},
	)

	registry.MustRegister(
		scoreProcessingAttemptCounter,
		scoreProcessingSuccessCounter,
		scoreProcessingFailureCounter,
		scoreProcessingDuration,
		scoreCorrectionAttemptCounter,
		scoreCorrectionSuccessCounter,
		scoreCorrectionFailureCounter,
		scoreCorrectionDuration,
		leaderboardUpdateAttemptCounter,
		leaderboardUpdateSuccessCounter,
		leaderboardUpdateFailureCounter,
		leaderboardUpdateDuration,
		dbQueryDuration,
	)

	return &scoreMetrics{
		scoreProcessingAttemptCounter:   scoreProcessingAttemptCounter,
		scoreProcessingSuccessCounter:   scoreProcessingSuccessCounter,
		scoreProcessingFailureCounter:   scoreProcessingFailureCounter,
		scoreProcessingDuration:         scoreProcessingDuration,
		scoreCorrectionAttemptCounter:   scoreCorrectionAttemptCounter,
		scoreCorrectionSuccessCounter:   scoreCorrectionSuccessCounter,
		scoreCorrectionFailureCounter:   scoreCorrectionFailureCounter,
		scoreCorrectionDuration:         scoreCorrectionDuration,
		leaderboardUpdateAttemptCounter: leaderboardUpdateAttemptCounter,
		leaderboardUpdateSuccessCounter: leaderboardUpdateSuccessCounter,
		leaderboardUpdateFailureCounter: leaderboardUpdateFailureCounter,
		leaderboardUpdateDuration:       leaderboardUpdateDuration,
		dbQueryDuration:                 dbQueryDuration,
	}
}

// RecordScoreProcessingAttempt records a score processing attempt
func (m *scoreMetrics) RecordScoreProcessingAttempt(roundID sharedtypes.RoundID) {
	m.scoreProcessingAttemptCounter.WithLabelValues(roundIDToString(roundID)).Inc()
}

// RecordScoreProcessingSuccess records a successful score processing attempt
func (m *scoreMetrics) RecordScoreProcessingSuccess(roundID sharedtypes.RoundID) {
	m.scoreProcessingSuccessCounter.WithLabelValues(roundIDToString(roundID)).Inc()
}

// RecordScoreProcessingFailure records a failed score processing attempt
func (m *scoreMetrics) RecordScoreProcessingFailure(roundID sharedtypes.RoundID) {
	m.scoreProcessingFailureCounter.WithLabelValues(roundIDToString(roundID)).Inc()
}

// RecordScoreProcessingDuration records the duration of score processing
func (m *scoreMetrics) RecordScoreProcessingDuration(roundID sharedtypes.RoundID, duration float64) {
	m.scoreProcessingDuration.WithLabelValues(roundIDToString(roundID)).Observe(duration)
}

// RecordScoreCorrectionAttempt records a score correction attempt
func (m *scoreMetrics) RecordScoreCorrectionAttempt(roundID sharedtypes.RoundID) {
	m.scoreCorrectionAttemptCounter.WithLabelValues(roundIDToString(roundID)).Inc()
}

// RecordScoreCorrectionSuccess records a successful score correction
func (m *scoreMetrics) RecordScoreCorrectionSuccess(roundID sharedtypes.RoundID) {
	m.scoreCorrectionSuccessCounter.WithLabelValues(roundIDToString(roundID)).Inc()
}

// RecordScoreCorrectionFailure records a failed score correction
func (m *scoreMetrics) RecordScoreCorrectionFailure(roundID sharedtypes.RoundID) {
	m.scoreCorrectionFailureCounter.WithLabelValues(roundIDToString(roundID)).Inc()
}

// RecordScoreCorrectionDuration records the duration of score correction
func (m *scoreMetrics) RecordScoreCorrectionDuration(roundID sharedtypes.RoundID, duration float64) {
	m.scoreCorrectionDuration.WithLabelValues(roundIDToString(roundID)).Observe(duration)
}

// RecordLeaderboardUpdateAttempt records a leaderboard update attempt
func (m *scoreMetrics) RecordLeaderboardUpdateAttempt(roundID sharedtypes.RoundID) {
	m.leaderboardUpdateAttemptCounter.WithLabelValues(roundIDToString(roundID)).Inc()
}

// RecordLeaderboardUpdateSuccess records a successful leaderboard update
func (m *scoreMetrics) RecordLeaderboardUpdateSuccess(roundID sharedtypes.RoundID) {
	m.leaderboardUpdateSuccessCounter.WithLabelValues(roundIDToString(roundID)).Inc()
}

// RecordLeaderboardUpdateFailure records a failed leaderboard update
func (m *scoreMetrics) RecordLeaderboardUpdateFailure(roundID sharedtypes.RoundID) {
	m.leaderboardUpdateFailureCounter.WithLabelValues(roundIDToString(roundID)).Inc()
}

// RecordLeaderboardUpdateDuration records the duration of a leaderboard update
func (m *scoreMetrics) RecordLeaderboardUpdateDuration(roundID sharedtypes.RoundID, duration float64) {
	m.leaderboardUpdateDuration.WithLabelValues(roundIDToString(roundID)).Observe(duration)
}

// RecordDBQueryDuration records the duration of a database query
func (m *scoreMetrics) RecordDBQueryDuration(duration float64) {
	m.dbQueryDuration.WithLabelValues().Observe(duration)
}
