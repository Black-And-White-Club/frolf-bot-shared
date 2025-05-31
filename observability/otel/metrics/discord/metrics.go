// observability/otel/discordmetrics/metrics.go
package discordmetrics

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// RecordAPIRequestDuration records the time taken for a Discord API request
func (m *discordMetrics) RecordAPIRequestDuration(ctx context.Context, endpoint string, duration time.Duration) {
	m.apiRequestDurationHistogram.Record(ctx, duration.Seconds(), enrichRecordAttrs(ctx, endpoint)...)
}

func (m *discordMetrics) RecordAPIRequest(ctx context.Context, endpoint string) {
	m.apiRequestCounter.Add(ctx, 1, enrichAddAttrs(ctx, endpoint)...)
}

func (m *discordMetrics) RecordAPIError(ctx context.Context, endpoint string, errorType string) {
	extra := []attribute.KeyValue{attribute.String("error_type", errorType)}
	m.apiErrorCounter.Add(ctx, 1, enrichAddAttrs(ctx, endpoint, extra...)...)
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

// RecordHandlerAttempt records a handler attempt
func (m *discordMetrics) RecordHandlerAttempt(ctx context.Context, handlerName string) {
	m.handlerAttemptCounter.Add(ctx, 1, metric.WithAttributes(handlerAttrs(handlerName)))
}

// RecordHandlerSuccess records a successful handler attempt
func (m *discordMetrics) RecordHandlerSuccess(ctx context.Context, handlerName string) {
	m.handlerSuccessCounter.Add(ctx, 1, metric.WithAttributes(handlerAttrs(handlerName)))
}

// RecordHandlerFailure records a failed handler attempt
func (m *discordMetrics) RecordHandlerFailure(ctx context.Context, handlerName string) {
	m.handlerFailureCounter.Add(ctx, 1, metric.WithAttributes(handlerAttrs(handlerName)))
}

// RecordHandlerDuration records the duration of a handler execution
func (m *discordMetrics) RecordHandlerDuration(ctx context.Context, handlerName string, duration time.Duration) {
	m.handlerDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(handlerAttrs(handlerName)))
}
