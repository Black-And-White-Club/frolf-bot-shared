// observability/otel/discordmetrics/constructor.go
package discordmetrics

import "go.opentelemetry.io/otel/metric"

// NewDiscordMetrics creates a new DiscordMetrics implementation using OpenTelemetry.
func NewDiscordMetrics(meter metric.Meter, prefix string) (DiscordMetrics, error) {
	// Helper function to create metric names with prefix
	metricName := func(name string) string {
		if prefix != "" {
			return prefix + "_discord_" + name
		}
		return "discord_" + name
	}

	var err error
	m := &discordMetrics{meter: meter}

	// API Request Metrics
	m.apiRequestDurationHistogram, err = meter.Float64Histogram(
		metricName("api_request_duration_seconds"),
		metric.WithDescription("Time taken to send API requests to Discord."),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, err
	}

	m.apiRequestCounter, err = meter.Int64Counter(
		metricName("api_requests_total"),
		metric.WithDescription("Total number of API requests sent to Discord, partitioned by endpoint."),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	m.apiErrorCounter, err = meter.Int64Counter(
		metricName("api_errors_total"),
		metric.WithDescription("Total number of Discord API request failures, partitioned by error type."),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	// Rate Limit Metrics
	m.rateLimitCounter, err = meter.Int64Counter(
		metricName("rate_limits_total"),
		metric.WithDescription("Total number of Discord API rate limits encountered, partitioned by endpoint."),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	m.rateLimitDurationHistogram, err = meter.Float64Histogram(
		metricName("rate_limit_reset_seconds"),
		metric.WithDescription("Time until rate limit reset as reported by Discord API."),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, err
	}

	// Websocket Metrics
	m.websocketEventCounter, err = meter.Int64Counter(
		metricName("websocket_events_total"),
		metric.WithDescription("Total number of Discord websocket events received, partitioned by event type."),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	m.websocketReconnectCounter, err = meter.Int64Counter(
		metricName("websocket_reconnects_total"),
		metric.WithDescription("Total number of websocket reconnections."),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	m.websocketDisconnectCounter, err = meter.Int64Counter(
		metricName("websocket_disconnects_total"),
		metric.WithDescription("Total number of websocket disconnects, partitioned by reason."),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	return m, nil
}
