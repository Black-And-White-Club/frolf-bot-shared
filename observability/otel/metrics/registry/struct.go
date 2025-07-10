// observability/otel/metrics/registry/struct.go
package registrymetrics

import "go.opentelemetry.io/otel/metric"

// registryMetrics implements RegistryMetrics using OpenTelemetry
// Mirrors the discordMetrics struct pattern
// Add local errorStats for error tracking

type registryMetrics struct {
	meter metric.Meter

	configRequests   metric.Int64Counter
	cacheHits        metric.Int64Counter
	cacheMisses      metric.Int64Counter
	errors           metric.Int64Counter
	requestDuration  metric.Float64Histogram
	cacheOpDuration  metric.Float64Histogram
	inflightRequests metric.Int64UpDownCounter
	cacheSize        metric.Int64ObservableGauge

	errorStats *ErrorStats
}
