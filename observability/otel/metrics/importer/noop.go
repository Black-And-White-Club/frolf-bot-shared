package importermetrics

import (
	"context"
	"time"
)

// NoOpMetrics is a no-op implementation of ImporterMetrics
type NoOpMetrics struct{}

// NewNoOpMetrics creates a new no-op importer metrics instance
func NewNoOpMetrics() ImporterMetrics {
	return &NoOpMetrics{}
}

// RecordImportAttempt is a no-op implementation
func (n *NoOpMetrics) RecordImportAttempt(ctx context.Context) {}

// RecordImportSuccess is a no-op implementation
func (n *NoOpMetrics) RecordImportSuccess(ctx context.Context) {}

// RecordImportFailure is a no-op implementation
func (n *NoOpMetrics) RecordImportFailure(ctx context.Context) {}

// RecordParseDuration is a no-op implementation
func (n *NoOpMetrics) RecordParseDuration(ctx context.Context, duration time.Duration) {}

// RecordParseSuccess is a no-op implementation
func (n *NoOpMetrics) RecordParseSuccess(ctx context.Context, playerCount int) {}

// RecordParseFailure is a no-op implementation
func (n *NoOpMetrics) RecordParseFailure(ctx context.Context) {}

// RecordMatchDuration is a no-op implementation
func (n *NoOpMetrics) RecordMatchDuration(ctx context.Context, duration time.Duration) {}

// RecordMatchSuccess is a no-op implementation
func (n *NoOpMetrics) RecordMatchSuccess(ctx context.Context, matched int, unmatched int) {}

// RecordMatchFailure is a no-op implementation
func (n *NoOpMetrics) RecordMatchFailure(ctx context.Context) {}

// RecordScoreIngestionDuration is a no-op implementation
func (n *NoOpMetrics) RecordScoreIngestionDuration(ctx context.Context, duration time.Duration) {}

// RecordScoresIngested is a no-op implementation
func (n *NoOpMetrics) RecordScoresIngested(ctx context.Context, count int) {}

// RecordScoreIngestionFailure is a no-op implementation
func (n *NoOpMetrics) RecordScoreIngestionFailure(ctx context.Context) {}
