// Package registrymetrics provides OpenTelemetry metrics for the registry cache.
package registrymetrics

import (
	"context"
	"sync"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// Ensure registryMetrics implements RegistryMetrics
var _ RegistryMetrics = (*registryMetrics)(nil)

// ErrorStats tracks error counts by type and guild in a thread-safe way.
type ErrorStats struct {
	mu    sync.RWMutex
	stats map[string]map[string]int64 // errorType -> guildID -> count
}

func NewErrorStats() *ErrorStats {
	return &ErrorStats{stats: make(map[string]map[string]int64)}
}

func (e *ErrorStats) Record(guildID, errorType string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.stats[errorType] == nil {
		e.stats[errorType] = make(map[string]int64)
	}
	e.stats[errorType][guildID]++
}

func (e *ErrorStats) Get() map[string]map[string]int64 {
	e.mu.RLock()
	defer e.mu.RUnlock()
	copyStats := make(map[string]map[string]int64)
	for et, guilds := range e.stats {
		copyStats[et] = make(map[string]int64)
		for gid, cnt := range guilds {
			copyStats[et][gid] = cnt
		}
	}
	return copyStats
}

func (e *ErrorStats) Reset() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.stats = make(map[string]map[string]int64)
}

// NewRegistryMetrics creates and registers all registry metrics instruments.
func NewRegistryMetrics(meter metric.Meter, cacheSizeFunc func() int64) (RegistryMetrics, error) {
	configRequests, err := meter.Int64Counter(
		"registry_config_requests_total",
		metric.WithDescription("Total registry config requests (success/failure)"),
	)
	if err != nil {
		return nil, err
	}
	cacheHits, err := meter.Int64Counter(
		"registry_cache_hits_total",
		metric.WithDescription("Total registry cache hits"),
	)
	if err != nil {
		return nil, err
	}
	cacheMisses, err := meter.Int64Counter(
		"registry_cache_misses_total",
		metric.WithDescription("Total registry cache misses"),
	)
	if err != nil {
		return nil, err
	}
	errors, err := meter.Int64Counter(
		"registry_errors_total",
		metric.WithDescription("Total registry errors (by type/guild)"),
	)
	if err != nil {
		return nil, err
	}
	requestDuration, err := meter.Float64Histogram(
		"registry_config_request_duration_seconds",
		metric.WithDescription("Config request duration (seconds)"),
	)
	if err != nil {
		return nil, err
	}
	cacheOpDuration, err := meter.Float64Histogram(
		"registry_cache_operation_duration_seconds",
		metric.WithDescription("Cache operation duration (seconds)"),
	)
	if err != nil {
		return nil, err
	}
	inflightRequests, err := meter.Int64UpDownCounter(
		"registry_inflight_requests",
		metric.WithDescription("Current in-flight registry requests"),
	)
	if err != nil {
		return nil, err
	}
	var cacheSize metric.Int64ObservableGauge
	cacheSize, err = meter.Int64ObservableGauge(
		"registry_cache_size",
		metric.WithDescription("Current registry cache size"),
	)
	if err != nil {
		return nil, err
	}
	// Register callback for cache size gauge
	_, err = meter.RegisterCallback(
		func(ctx context.Context, o metric.Observer) error {
			o.ObserveInt64(cacheSize, cacheSizeFunc())
			return nil
		},
		cacheSize,
	)
	if err != nil {
		return nil, err
	}

	return &registryMetrics{
		meter:            meter,
		configRequests:   configRequests,
		cacheHits:        cacheHits,
		cacheMisses:      cacheMisses,
		errors:           errors,
		requestDuration:  requestDuration,
		cacheOpDuration:  cacheOpDuration,
		inflightRequests: inflightRequests,
		cacheSize:        cacheSize,
		errorStats:       NewErrorStats(),
	}, nil
}

// RecordConfigRequest records a config request metric.
func (m *registryMetrics) RecordConfigRequest(ctx context.Context, guildID string, success bool, duration time.Duration) {
	attrs := []attribute.KeyValue{guildIDAttr(guildID), successAttr(success)}
	m.configRequests.Add(ctx, 1, enrichAddAttrs(ctx, "config_request", attrs...)...)
	m.requestDuration.Record(ctx, duration.Seconds(), enrichRecordAttrs(ctx, "config_request", attrs...)...)
}

// RecordCacheHit records a cache hit metric.
func (m *registryMetrics) RecordCacheHit(ctx context.Context, guildID string) {
	m.cacheHits.Add(ctx, 1, enrichAddAttrs(ctx, "cache_hit", guildIDAttr(guildID))...)
}

// RecordCacheMiss records a cache miss metric.
func (m *registryMetrics) RecordCacheMiss(ctx context.Context, guildID string) {
	m.cacheMisses.Add(ctx, 1, enrichAddAttrs(ctx, "cache_miss", guildIDAttr(guildID))...)
}

// RecordInflightRequest adjusts the in-flight request gauge.
func (m *registryMetrics) RecordInflightRequest(ctx context.Context, guildID string, delta int64) {
	m.inflightRequests.Add(ctx, delta, enrichAddAttrs(ctx, "inflight_request", guildIDAttr(guildID))...)
}

// RecordError records an error metric and updates local error stats.
func (m *registryMetrics) RecordError(ctx context.Context, guildID, errorType string) {
	m.errors.Add(ctx, 1, enrichAddAttrs(ctx, "error", guildIDAttr(guildID), errorTypeAttr(errorType))...)
	m.errorStats.Record(guildID, errorType)
}

// GetErrorMetrics returns a thread-safe copy of error stats.
func (m *registryMetrics) GetErrorMetrics() map[string]map[string]int64 {
	return m.errorStats.Get()
}

// ResetErrorMetrics resets local error stats (for testing).
func (m *registryMetrics) ResetErrorMetrics() {
	m.errorStats.Reset()
}

// RecordCacheOperation records a cache operation duration.
func (m *registryMetrics) RecordCacheOperation(ctx context.Context, operation string, duration time.Duration) {
	m.cacheOpDuration.Record(ctx, duration.Seconds(), enrichRecordAttrs(ctx, operation)...)
}

// RecordBackendLatency is a placeholder for backend latency metrics.
func (m *registryMetrics) RecordBackendLatency(ctx context.Context, guildID string, duration time.Duration) {
	// Implement if needed for backend calls
}
