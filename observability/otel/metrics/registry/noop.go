// observability/otel/metrics/registry/noop.go
package registrymetrics

import (
	"context"
	"time"
)

type noopRegistryMetrics struct{}

func NewNoop() RegistryMetrics { return &noopRegistryMetrics{} }

func (n *noopRegistryMetrics) RecordConfigRequest(ctx context.Context, guildID string, success bool, duration time.Duration) {
}
func (n *noopRegistryMetrics) RecordCacheHit(ctx context.Context, guildID string)  {}
func (n *noopRegistryMetrics) RecordCacheMiss(ctx context.Context, guildID string) {}
func (n *noopRegistryMetrics) RecordInflightRequest(ctx context.Context, guildID string, delta int64) {
}
func (n *noopRegistryMetrics) RecordError(ctx context.Context, guildID, errorType string) {}
func (n *noopRegistryMetrics) GetErrorMetrics() map[string]map[string]int64 {
	return map[string]map[string]int64{}
}
func (n *noopRegistryMetrics) ResetErrorMetrics() {}
func (n *noopRegistryMetrics) RecordCacheOperation(ctx context.Context, operation string, duration time.Duration) {
}

func (n *noopRegistryMetrics) RecordBackendLatency(ctx context.Context, guildID string, duration time.Duration) {
}
