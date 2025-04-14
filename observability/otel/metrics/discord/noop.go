// observability/otel/discordmetrics/noop.go
package discordmetrics

import (
	"context"
	"time"
)

// NoOpMetrics is a metrics collector that does nothing. Useful for unit tests.
type NoOpMetrics struct{}

// NewNoop returns a no-operation implementation of DiscordMetrics
func NewNoop() DiscordMetrics {
	return &NoOpMetrics{}
}

// RecordAPIRequestDuration does nothing
func (n *NoOpMetrics) RecordAPIRequestDuration(ctx context.Context, endpoint string, duration time.Duration) {
}

// RecordAPIRequest does nothing
func (n *NoOpMetrics) RecordAPIRequest(ctx context.Context, endpoint string) {
}

// RecordAPIError does nothing
func (n *NoOpMetrics) RecordAPIError(ctx context.Context, endpoint string, errorType string) {
}

// RecordRateLimit does nothing
func (n *NoOpMetrics) RecordRateLimit(ctx context.Context, endpoint string, resetTime time.Duration) {
}

// RecordWebsocketEvent does nothing
func (n *NoOpMetrics) RecordWebsocketEvent(ctx context.Context, eventType string) {
}

// RecordWebsocketReconnect does nothing
func (n *NoOpMetrics) RecordWebsocketReconnect(ctx context.Context) {
}

// RecordWebsocketDisconnect does nothing
func (n *NoOpMetrics) RecordWebsocketDisconnect(ctx context.Context, reason string) {
}
