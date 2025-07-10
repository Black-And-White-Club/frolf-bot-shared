// observability/otel/metrics/registry/constructor.go
package registrymetrics

import "go.opentelemetry.io/otel/metric"

// New is an alias for NewRegistryMetrics for consistency with other metrics packages.
func New(meter metric.Meter, cacheSizeFunc func() int64) (RegistryMetrics, error) {
	return NewRegistryMetrics(meter, cacheSizeFunc)
}
