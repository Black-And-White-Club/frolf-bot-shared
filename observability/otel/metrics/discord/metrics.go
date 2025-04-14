// observability/otel/discordmetrics/metrics.go
package discordmetrics

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/metric"
)

// RecordAPIRequestDuration records the time taken for a Discord API request
func (m *discordMetrics) RecordAPIRequestDuration(ctx context.Context, endpoint string, duration time.Duration) {
	m.apiRequestDurationHistogram.Record(ctx, duration.Seconds(), metric.WithAttributes(endpointAttr(endpoint)))
}

// RecordAPIRequest records a Discord API request completion
func (m *discordMetrics) RecordAPIRequest(ctx context.Context, endpoint string) {
	m.apiRequestCounter.Add(ctx, 1, metric.WithAttributes(endpointAttr(endpoint)))
}

// RecordAPIError records a Discord API request error
func (m *discordMetrics) RecordAPIError(ctx context.Context, endpoint string, errorType string) {
	m.apiErrorCounter.Add(ctx, 1, metric.WithAttributes(endpointErrorAttrs(endpoint, errorType)...))
}

// RecordRateLimit records encountering a Discord API rate limit
func (m *discordMetrics) RecordRateLimit(ctx context.Context, endpoint string, resetTime time.Duration) {
	m.rateLimitCounter.Add(ctx, 1, metric.WithAttributes(endpointAttr(endpoint)))
	m.rateLimitDurationHistogram.Record(ctx, resetTime.Seconds())
}

// RecordWebsocketEvent records receiving a Discord websocket event
func (m *discordMetrics) RecordWebsocketEvent(ctx context.Context, eventType string) {
	m.websocketEventCounter.Add(ctx, 1, metric.WithAttributes(eventTypeAttr(eventType)))
}

// RecordWebsocketReconnect records a websocket reconnection
func (m *discordMetrics) RecordWebsocketReconnect(ctx context.Context) {
	m.websocketReconnectCounter.Add(ctx, 1)
}

// RecordWebsocketDisconnect records a websocket disconnection
func (m *discordMetrics) RecordWebsocketDisconnect(ctx context.Context, reason string) {
	m.websocketDisconnectCounter.Add(ctx, 1, metric.WithAttributes(disconnectReasonAttr(reason)))
}
