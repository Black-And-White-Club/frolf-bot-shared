package guildmetrics

import (
	"context"
	"time"

	sharedtypes "github.com/Black-And-White-Club/frolf-bot-shared/types/shared"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// Common attributes for guild metrics

func handlerAttr(handlerName string) attribute.KeyValue {
	return attribute.String("handler", handlerName)
}

func (m *guildMetrics) RecordGuildCreated(ctx context.Context, success bool, guildID sharedtypes.GuildID, source string) {
	m.guildCreatedCounter.Add(ctx, 1, metric.WithAttributes(
		successAttr(success),
		guildIDAttr(guildID),
		sourceAttr(source),
	))
}

func (m *guildMetrics) RecordGuildDeleted(ctx context.Context, success bool, guildID sharedtypes.GuildID, source string) {
	m.guildDeletedCounter.Add(ctx, 1, metric.WithAttributes(
		successAttr(success),
		guildIDAttr(guildID),
		sourceAttr(source),
	))
}

func (m *guildMetrics) RecordOperationAttempt(ctx context.Context, operationName string, guildID sharedtypes.GuildID, serviceName string) {
	m.operationAttemptCounter.Add(ctx, 1, metric.WithAttributes(
		operationAttr(operationName),
		guildIDAttr(guildID),
		serviceAttr(serviceName),
	))
}

func (m *guildMetrics) RecordOperationSuccess(ctx context.Context, operationName string, guildID sharedtypes.GuildID, serviceName string) {
	m.operationSuccessCounter.Add(ctx, 1, metric.WithAttributes(
		operationAttr(operationName),
		guildIDAttr(guildID),
		serviceAttr(serviceName),
	))
}

func (m *guildMetrics) RecordOperationFailure(ctx context.Context, operationName string, guildID sharedtypes.GuildID, serviceName string) {
	m.operationFailureCounter.Add(ctx, 1, metric.WithAttributes(
		operationAttr(operationName),
		guildIDAttr(guildID),
		serviceAttr(serviceName),
	))
}

func (m *guildMetrics) RecordOperationDuration(ctx context.Context, operationName string, guildID sharedtypes.GuildID, serviceName string, duration time.Duration) {
	m.operationDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(
		operationAttr(operationName),
		guildIDAttr(guildID),
		serviceAttr(serviceName),
	))
}

// RecordHandlerAttempt records a handler attempt for a guild event handler.
func (m *guildMetrics) RecordHandlerAttempt(ctx context.Context, handlerName string) {
	m.handlerAttemptCounter.Add(ctx, 1, metric.WithAttributes(handlerAttr(handlerName)))
}

// RecordHandlerSuccess records a handler success for a guild event handler.
func (m *guildMetrics) RecordHandlerSuccess(ctx context.Context, handlerName string) {
	m.handlerSuccessCounter.Add(ctx, 1, metric.WithAttributes(handlerAttr(handlerName)))
}

// RecordHandlerFailure records a handler failure for a guild event handler.
func (m *guildMetrics) RecordHandlerFailure(ctx context.Context, handlerName string) {
	m.handlerFailureCounter.Add(ctx, 1, metric.WithAttributes(handlerAttr(handlerName)))
}

// RecordHandlerDuration records the duration of a handler for a guild event handler.
func (m *guildMetrics) RecordHandlerDuration(ctx context.Context, handlerName string, duration time.Duration) {
	m.handlerDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(handlerAttr(handlerName)))
}
