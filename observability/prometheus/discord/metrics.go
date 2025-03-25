// observability/prometheus/discord/metrics.go
package discordmetrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// DiscordMetrics defines metrics specific to Discord API operations
type DiscordMetrics interface {
	// Record API request duration
	RecordAPIRequestDuration(seconds float64)

	// Record API request completion
	RecordAPIRequest(endpoint string)

	// Record API request error
	RecordAPIError(errorType string)

	// Record rate limit encountered
	RecordRateLimit(endpoint string, resetTimeSeconds float64)

	// Record websocket event received
	RecordWebsocketEvent(eventType string)
}

// discordMetrics implements DiscordMetrics
type discordMetrics struct {
	apiRequestDurationHistogram *prometheus.Histogram
	apiRequestCounter           *prometheus.CounterVec
	apiErrorCounter             *prometheus.CounterVec
	rateLimitCounter            *prometheus.CounterVec
	rateLimitDurationHistogram  *prometheus.Histogram
	websocketEventCounter       *prometheus.CounterVec
}

// NewDiscordMetrics creates a new DiscordMetrics implementation
func NewDiscordMetrics(registry *prometheus.Registry, prefix string) DiscordMetrics {
	apiRequestDurationHistogram := prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: prefix,
			Subsystem: "discord",
			Name:      "api_request_duration_seconds",
			Help:      "Time taken to send API requests to Discord.",
			Buckets:   prometheus.DefBuckets,
		},
	)

	apiRequestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: prefix,
			Subsystem: "discord",
			Name:      "api_requests_total",
			Help:      "Total number of API requests sent to Discord, partitioned by endpoint.",
		},
		[]string{"endpoint"},
	)

	apiErrorCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: prefix,
			Subsystem: "discord",
			Name:      "api_errors_total",
			Help:      "Total number of Discord API request failures, partitioned by error type.",
		},
		[]string{"error_type"},
	)

	rateLimitCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: prefix,
			Subsystem: "discord",
			Name:      "rate_limits_total",
			Help:      "Total number of Discord API rate limits encountered, partitioned by endpoint.",
		},
		[]string{"endpoint"},
	)

	rateLimitDurationHistogram := prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: prefix,
			Subsystem: "discord",
			Name:      "rate_limit_reset_seconds",
			Help:      "Time until rate limit reset as reported by Discord API.",
			Buckets:   []float64{1, 2, 5, 10, 30, 60, 300, 600},
		},
	)

	websocketEventCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: prefix,
			Subsystem: "discord",
			Name:      "websocket_events_total",
			Help:      "Total number of Discord websocket events received, partitioned by event type.",
		},
		[]string{"event_type"},
	)

	registry.MustRegister(
		apiRequestDurationHistogram,
		apiRequestCounter,
		apiErrorCounter,
		rateLimitCounter,
		rateLimitDurationHistogram,
		websocketEventCounter,
	)

	return &discordMetrics{
		apiRequestDurationHistogram: &apiRequestDurationHistogram,
		apiRequestCounter:           apiRequestCounter,
		apiErrorCounter:             apiErrorCounter,
		rateLimitCounter:            rateLimitCounter,
		rateLimitDurationHistogram:  &rateLimitDurationHistogram,
		websocketEventCounter:       websocketEventCounter,
	}
}

// RecordAPIRequestDuration records the time taken for a Discord API request
func (m *discordMetrics) RecordAPIRequestDuration(seconds float64) {
	(*m.apiRequestDurationHistogram).Observe(seconds)
}

// RecordAPIRequest records a Discord API request completion
func (m *discordMetrics) RecordAPIRequest(endpoint string) {
	m.apiRequestCounter.WithLabelValues(endpoint).Inc()
}

// RecordAPIError records a Discord API request error
func (m *discordMetrics) RecordAPIError(errorType string) {
	m.apiErrorCounter.WithLabelValues(errorType).Inc()
}

// RecordRateLimit records encountering a Discord API rate limit
func (m *discordMetrics) RecordRateLimit(endpoint string, resetTimeSeconds float64) {
	m.rateLimitCounter.WithLabelValues(endpoint).Inc()
	(*m.rateLimitDurationHistogram).Observe(resetTimeSeconds)
}

// RecordWebsocketEvent records receiving a Discord websocket event
func (m *discordMetrics) RecordWebsocketEvent(eventType string) {
	m.websocketEventCounter.WithLabelValues(eventType).Inc()
}
