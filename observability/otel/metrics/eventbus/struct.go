// observability/otel/eventbusmetrics/struct.go
package eventbusmetrics

import "go.opentelemetry.io/otel/metric"

// eventBusMetrics implements EventBusMetrics using OpenTelemetry
type eventBusMetrics struct {
	meter metric.Meter // OTEL Meter

	// Message Publish Metrics
	messagePublishCounter      metric.Int64Counter
	messagePublishErrorCounter metric.Int64Counter

	// Message Process Metrics
	messageProcessCounter          metric.Int64Counter
	messageProcessingTimeHistogram metric.Float64Histogram // Duration in seconds

	// Message Subscribe Metrics
	messageSubscribeCounter      metric.Int64Counter
	messageSubscribeErrorCounter metric.Int64Counter
}
