package importermetrics

import (
	"context"
	"time"
)

// ImporterMetrics defines metrics for importer operations.
type ImporterMetrics interface {
	// Import operation metrics
	RecordImportAttempt(ctx context.Context)
	RecordImportSuccess(ctx context.Context)
	RecordImportFailure(ctx context.Context)

	// Parse phase metrics
	RecordParseDuration(ctx context.Context, duration time.Duration)
	RecordParseSuccess(ctx context.Context, playerCount int)
	RecordParseFailure(ctx context.Context)

	// Match phase metrics
	RecordMatchDuration(ctx context.Context, duration time.Duration)
	RecordMatchSuccess(ctx context.Context, matched int, unmatched int)
	RecordMatchFailure(ctx context.Context)

	// Score ingestion metrics
	RecordScoreIngestionDuration(ctx context.Context, duration time.Duration)
	RecordScoresIngested(ctx context.Context, count int)
	RecordScoreIngestionFailure(ctx context.Context)
}
