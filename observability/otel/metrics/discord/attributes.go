// observability/otel/discordmetrics/attributes.go
package discordmetrics

import (
	"context"
	"fmt"

	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type key string

const (
	GuildIDKey       key = "guild_id"
	UserIDKey        key = "user_id"
	CommandNameKey   key = "command_name"
	InteractionType  key = "interaction_type"
	CorrelationIDKey key = "correlation_id"
	MessageIDKey     key = "discord_message_id"
)

func WithValue[T any](ctx context.Context, k key, val T) context.Context {
	return context.WithValue(ctx, k, val)
}

func Get(ctx context.Context, k key) (string, bool) {
	switch val := ctx.Value(k).(type) {
	case string:
		return val, true
	case fmt.Stringer:
		return val.String(), true
	case sharedtypes.DiscordID:
		return string(val), true
	default:
		return "", false
	}
}

// Common attributes
func endpointAttr(endpoint string) attribute.KeyValue {
	return attribute.String("endpoint", endpoint)
}

// func errorTypeAttr(errorType string) attribute.KeyValue {
// 	return attribute.String("error_type", errorType)
// }

func eventTypeAttr(eventType string) attribute.KeyValue {
	return attribute.String("event_type", eventType)
}

func disconnectReasonAttr(reason string) attribute.KeyValue {
	return attribute.String("reason", reason)
}

// // Combined attribute sets
// func endpointAttrs(endpoint string) attribute.KeyValue {
// 	return endpointAttr(endpoint)
// }

// func endpointErrorAttrs(endpoint string, errorType string) []attribute.KeyValue {
// 	return []attribute.KeyValue{
// 		endpointAttr(endpoint),
// 		errorTypeAttr(errorType),
// 	}
// }

// attributes.go
func enrichAddAttrs(ctx context.Context, endpoint string, extra ...attribute.KeyValue) []metric.AddOption {
	attrs := gatherCommonAttrs(ctx, endpoint, extra...)
	return []metric.AddOption{metric.WithAttributes(attrs...)}
}

func enrichRecordAttrs(ctx context.Context, endpoint string, extra ...attribute.KeyValue) []metric.RecordOption {
	attrs := gatherCommonAttrs(ctx, endpoint, extra...)
	return []metric.RecordOption{metric.WithAttributes(attrs...)}
}

func gatherCommonAttrs(ctx context.Context, endpoint string, extra ...attribute.KeyValue) []attribute.KeyValue {
	attrs := []attribute.KeyValue{
		attribute.String("endpoint", endpoint),
	}
	attrs = append(attrs, extra...)

	if v, ok := Get(ctx, GuildIDKey); ok {
		attrs = append(attrs, attribute.String("guild_id", v))
	}
	if v, ok := Get(ctx, UserIDKey); ok {
		attrs = append(attrs, attribute.String("user_id", v))
	}
	if v, ok := Get(ctx, CommandNameKey); ok {
		attrs = append(attrs, attribute.String("command", v))
	}
	if v, ok := Get(ctx, InteractionType); ok {
		attrs = append(attrs, attribute.String("interaction_type", v))
	}

	return attrs
}

func enrichMeasureAttrs(ctx context.Context, endpoint string) []metric.MeasurementOption {
	attrs := []attribute.KeyValue{
		attribute.String("endpoint", endpoint),
	}

	if v, ok := Get(ctx, GuildIDKey); ok {
		attrs = append(attrs, attribute.String("guild_id", v))
	}
	if v, ok := Get(ctx, UserIDKey); ok {
		attrs = append(attrs, attribute.String("user_id", v))
	}
	if v, ok := Get(ctx, CommandNameKey); ok {
		attrs = append(attrs, attribute.String("command", v))
	}
	if v, ok := Get(ctx, InteractionType); ok {
		attrs = append(attrs, attribute.String("interaction_type", v))
	}

	return []metric.MeasurementOption{metric.WithAttributes(attrs...)}
}

// func eventTypeAttrs(eventType string) attribute.KeyValue {
// 	return eventTypeAttr(eventType)
// }

// func disconnectAttrs(reason string) attribute.KeyValue {
// 	return disconnectReasonAttr(reason)
// }

func handlerAttr(handlerName string) attribute.KeyValue {
	return attribute.String("handler", handlerName)
}

func handlerAttrs(handlerName string) attribute.KeyValue {
	return handlerAttr(handlerName)
}
