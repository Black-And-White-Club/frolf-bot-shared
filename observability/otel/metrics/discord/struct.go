// observability/otel/discordmetrics/struct.go
package discordmetrics

import "go.opentelemetry.io/otel/metric"

// discordMetrics implements DiscordMetrics using OpenTelemetry
type discordMetrics struct {
	meter metric.Meter // OTEL Meter

	// API metrics
	apiRequestDurationHistogram metric.Float64Histogram
	apiRequestCounter           metric.Int64Counter
	apiErrorCounter             metric.Int64Counter

	// Rate limit metrics
	rateLimitCounter           metric.Int64Counter
	rateLimitDurationHistogram metric.Float64Histogram

	// Websocket metrics
	websocketEventCounter      metric.Int64Counter
	websocketReconnectCounter  metric.Int64Counter
	websocketDisconnectCounter metric.Int64Counter

	// --- Handler Metrics ---
	handlerAttemptCounter metric.Int64Counter
	handlerSuccessCounter metric.Int64Counter
	handlerFailureCounter metric.Int64Counter
	handlerDuration       metric.Float64Histogram
}
