// Package registrymetrics provides attribute helpers for registry metrics.
package registrymetrics

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

const (
	AttrGuildID     = "guild_id"
	AttrErrorType   = "error_type"
	AttrRequestType = "request_type"
	AttrOperation   = "operation"
	AttrSuccess     = "success"
)

func guildIDAttr(guildID string) attribute.KeyValue {
	return attribute.String(AttrGuildID, guildID)
}

func errorTypeAttr(errorType string) attribute.KeyValue {
	return attribute.String(AttrErrorType, errorType)
}

func requestTypeAttr(requestType string) attribute.KeyValue {
	return attribute.String(AttrRequestType, requestType)
}

func cacheOperationAttr(operation string) attribute.KeyValue {
	return attribute.String(AttrOperation, operation)
}

func successAttr(success bool) attribute.KeyValue {
	return attribute.Bool(AttrSuccess, success)
}

// Attribute helpers for context enrichment, following discordmetrics pattern
// These helpers are used to build []metric.AddOption and []metric.RecordOption
func enrichAddAttrs(ctx context.Context, operation string, extra ...attribute.KeyValue) []metric.AddOption {
	attrs := gatherCommonAttrs(ctx, operation, extra...)
	return []metric.AddOption{metric.WithAttributes(attrs...)}
}

func enrichRecordAttrs(ctx context.Context, operation string, extra ...attribute.KeyValue) []metric.RecordOption {
	attrs := gatherCommonAttrs(ctx, operation, extra...)
	return []metric.RecordOption{metric.WithAttributes(attrs...)}
}

func gatherCommonAttrs(ctx context.Context, operation string, extra ...attribute.KeyValue) []attribute.KeyValue {
	attrs := []attribute.KeyValue{
		attribute.String("operation", operation),
	}
	attrs = append(attrs, extra...)
	// Add more context-based attributes as needed (guild_id, error_type, etc.)
	return attrs
}
