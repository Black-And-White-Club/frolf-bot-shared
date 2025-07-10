// Package registrymetrics defines the public interface for registry metrics.
package registrymetrics

import (
	"context"
	"time"
)

type RegistryMetrics interface {
	RecordConfigRequest(ctx context.Context, guildID string, success bool, duration time.Duration)
	RecordCacheHit(ctx context.Context, guildID string)
	RecordCacheMiss(ctx context.Context, guildID string)
	RecordInflightRequest(ctx context.Context, guildID string, delta int64)
	RecordError(ctx context.Context, guildID, errorType string)
	GetErrorMetrics() map[string]map[string]int64
	ResetErrorMetrics()
	RecordCacheOperation(ctx context.Context, operation string, duration time.Duration)
	RecordBackendLatency(ctx context.Context, guildID string, duration time.Duration)
}
