// observability/otel/discordmetrics/interface.go
package discordmetrics

import (
	"context"
	"time"
)

// DiscordMetrics defines metrics specific to Discord API operations using OpenTelemetry
type DiscordMetrics interface {
	// Record API request duration with enriched context-based attributes
	RecordAPIRequestDuration(ctx context.Context, endpoint string, duration time.Duration)

	// Record API request completion with enriched context-based attributes
	RecordAPIRequest(ctx context.Context, endpoint string)

	// Record API request error with enriched context-based attributes
	RecordAPIError(ctx context.Context, endpoint string, errorType string)

	// Record rate limit encountered
	RecordRateLimit(ctx context.Context, endpoint string, resetTime time.Duration)

	// Record websocket event received
	RecordWebsocketEvent(ctx context.Context, eventType string)

	// Record websocket connection state changes
	RecordWebsocketReconnect(ctx context.Context)
	RecordWebsocketDisconnect(ctx context.Context, reason string)

	RecordHandlerAttempt(ctx context.Context, handlerName string)
	RecordHandlerSuccess(ctx context.Context, handlerName string)
	RecordHandlerFailure(ctx context.Context, handlerName string)
	RecordHandlerDuration(ctx context.Context, handlerName string, duration time.Duration)
}
