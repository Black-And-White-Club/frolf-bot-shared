package importermetrics

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/metric"
)

// importerMetrics implements the ImporterMetrics interface using OpenTelemetry.
type importerMetrics struct {
	// Import operation counters
	importAttempts metric.Int64Counter
	importSuccess  metric.Int64Counter
	importFailure  metric.Int64Counter

	// Parse phase metrics
	parseDuration metric.Float64Histogram
	parseSuccess  metric.Int64Counter
	parseFailure  metric.Int64Counter

	// Match phase metrics
	matchDuration metric.Float64Histogram
	matchSuccess  metric.Int64Counter
	matchFailure  metric.Int64Counter

	// Score ingestion metrics
	scoreIngestionDuration metric.Float64Histogram
	scoresIngested         metric.Int64Counter
	scoreIngestionFailure  metric.Int64Counter
}

// NewImporterMetrics creates a new ImporterMetrics instance using OpenTelemetry.
func NewImporterMetrics(meter metric.Meter) (ImporterMetrics, error) {
	m := &importerMetrics{}
	var err error

	// Import operation counters
	m.importAttempts, err = meter.Int64Counter("importer.import.attempts.total")
	if err != nil {
		return nil, fmt.Errorf("failed to create import attempts counter: %w", err)
	}

	m.importSuccess, err = meter.Int64Counter("importer.import.success.total")
	if err != nil {
		return nil, fmt.Errorf("failed to create import success counter: %w", err)
	}

	m.importFailure, err = meter.Int64Counter("importer.import.failure.total")
	if err != nil {
		return nil, fmt.Errorf("failed to create import failure counter: %w", err)
	}

	// Parse phase metrics
	m.parseDuration, err = meter.Float64Histogram("importer.parse.duration.ms")
	if err != nil {
		return nil, fmt.Errorf("failed to create parse duration histogram: %w", err)
	}

	m.parseSuccess, err = meter.Int64Counter("importer.parse.success.total")
	if err != nil {
		return nil, fmt.Errorf("failed to create parse success counter: %w", err)
	}

	m.parseFailure, err = meter.Int64Counter("importer.parse.failure.total")
	if err != nil {
		return nil, fmt.Errorf("failed to create parse failure counter: %w", err)
	}

	// Match phase metrics
	m.matchDuration, err = meter.Float64Histogram("importer.match.duration.ms")
	if err != nil {
		return nil, fmt.Errorf("failed to create match duration histogram: %w", err)
	}

	m.matchSuccess, err = meter.Int64Counter("importer.match.success.total")
	if err != nil {
		return nil, fmt.Errorf("failed to create match success counter: %w", err)
	}

	m.matchFailure, err = meter.Int64Counter("importer.match.failure.total")
	if err != nil {
		return nil, fmt.Errorf("failed to create match failure counter: %w", err)
	}

	// Score ingestion metrics
	m.scoreIngestionDuration, err = meter.Float64Histogram("importer.score_ingestion.duration.ms")
	if err != nil {
		return nil, fmt.Errorf("failed to create score ingestion duration histogram: %w", err)
	}

	m.scoresIngested, err = meter.Int64Counter("importer.scores_ingested.total")
	if err != nil {
		return nil, fmt.Errorf("failed to create scores ingested counter: %w", err)
	}

	m.scoreIngestionFailure, err = meter.Int64Counter("importer.score_ingestion.failure.total")
	if err != nil {
		return nil, fmt.Errorf("failed to create score ingestion failure counter: %w", err)
	}

	return m, nil
}

// RecordImportAttempt records an import attempt.
func (m *importerMetrics) RecordImportAttempt(ctx context.Context) {
	m.importAttempts.Add(ctx, 1)
}

// RecordImportSuccess records a successful import.
func (m *importerMetrics) RecordImportSuccess(ctx context.Context) {
	m.importSuccess.Add(ctx, 1)
}

// RecordImportFailure records a failed import.
func (m *importerMetrics) RecordImportFailure(ctx context.Context) {
	m.importFailure.Add(ctx, 1)
}

// RecordParseDuration records the duration of the parse phase in milliseconds.
func (m *importerMetrics) RecordParseDuration(ctx context.Context, duration time.Duration) {
	m.parseDuration.Record(ctx, float64(duration.Milliseconds()))
}

// RecordParseSuccess records a successful parse with the player count.
func (m *importerMetrics) RecordParseSuccess(ctx context.Context, playerCount int) {
	m.parseSuccess.Add(ctx, int64(playerCount))
}

// RecordParseFailure records a failed parse.
func (m *importerMetrics) RecordParseFailure(ctx context.Context) {
	m.parseFailure.Add(ctx, 1)
}

// RecordMatchDuration records the duration of the match phase in milliseconds.
func (m *importerMetrics) RecordMatchDuration(ctx context.Context, duration time.Duration) {
	m.matchDuration.Record(ctx, float64(duration.Milliseconds()))
}

// RecordMatchSuccess records successful player matching.
func (m *importerMetrics) RecordMatchSuccess(ctx context.Context, matched int, unmatched int) {
	m.matchSuccess.Add(ctx, int64(matched))
}

// RecordMatchFailure records a failed match operation.
func (m *importerMetrics) RecordMatchFailure(ctx context.Context) {
	m.matchFailure.Add(ctx, 1)
}

// RecordScoreIngestionDuration records the duration of score ingestion in milliseconds.
func (m *importerMetrics) RecordScoreIngestionDuration(ctx context.Context, duration time.Duration) {
	m.scoreIngestionDuration.Record(ctx, float64(duration.Milliseconds()))
}

// RecordScoresIngested records the number of scores ingested.
func (m *importerMetrics) RecordScoresIngested(ctx context.Context, count int) {
	m.scoresIngested.Add(ctx, int64(count))
}

// RecordScoreIngestionFailure records a failed score ingestion.
func (m *importerMetrics) RecordScoreIngestionFailure(ctx context.Context) {
	m.scoreIngestionFailure.Add(ctx, 1)
}
