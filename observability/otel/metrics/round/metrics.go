package roundmetrics

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/metric"
)

// Implementations of the RoundMetrics interface
func (m *roundMetrics) RecordRoundDeleteAttempt(ctx context.Context) {
	m.roundDeleteAttemptCounter.Add(ctx, 1)
}

func (m *roundMetrics) RecordRoundDeleteSuccess(ctx context.Context) {
	m.roundDeleteSuccessCounter.Add(ctx, 1)
}

func (m *roundMetrics) RecordRoundDeleteFailure(ctx context.Context) {
	m.roundDeleteFailureCounter.Add(ctx, 1)
}

func (m *roundMetrics) RecordDBOperationDuration(ctx context.Context, operation string, duration time.Duration) {
	attrs := dbOperationAttrs(operation)
	m.dbOperationDurationHistogram.Record(ctx, duration.Seconds(), metric.WithAttributes(attrs...))
}

func (m *roundMetrics) RecordOperationAttempt(ctx context.Context, operation, service string) {
	attrs := operationAttrs(operation, service)
	m.operationAttemptCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
}

func (m *roundMetrics) RecordOperationDuration(ctx context.Context, operation, service string, duration time.Duration) {
	attrs := operationAttrs(operation, service)
	m.operationDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(attrs...))
}

func (m *roundMetrics) RecordOperationFailure(ctx context.Context, operation, service string) {
	attrs := operationAttrs(operation, service)
	m.operationFailureCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
}

func (m *roundMetrics) RecordOperationSuccess(ctx context.Context, operation, service string) {
	attrs := operationAttrs(operation, service)
	m.operationSuccessCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
}

func (m *roundMetrics) RecordRoundCreated(ctx context.Context, location string) {
	attrs := locationAttrs(location)
	m.roundCreatedCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
}

func (m *roundMetrics) RecordRoundParticipantAdded(ctx context.Context, location string) {
	attrs := locationAttrs(location)
	m.roundParticipantAddedCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
}

func (m *roundMetrics) RecordRoundFinalized(ctx context.Context, location string) {
	attrs := locationAttrs(location)
	m.roundFinalizedCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
}

func (m *roundMetrics) RecordRoundCancelled(ctx context.Context, location string) {
	attrs := locationAttrs(location)
	m.roundCancelledCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
}

func (m *roundMetrics) RecordTimeParsingError(ctx context.Context) {
	m.timeParsingErrorCounter.Add(ctx, 1)
}

func (m *roundMetrics) RecordValidationError(ctx context.Context) {
	m.validationErrorCounter.Add(ctx, 1)
}

func (m *roundMetrics) RecordDBOperationError(ctx context.Context, operation string) {
	attrs := dbOperationAttrs(operation)
	m.dbOperationErrorCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
}

func (m *roundMetrics) RecordValidationSuccess(ctx context.Context) {
	m.validationSuccessCounter.Add(ctx, 1)
}

func (m *roundMetrics) RecordTimeParsingSuccess(ctx context.Context) {
	m.timeParsingSuccessCounter.Add(ctx, 1)
}

func (m *roundMetrics) RecordDBOperationSuccess(ctx context.Context, operation string) {
	attrs := dbOperationAttrs(operation)
	m.dbOperationSuccessCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
}

func (m *roundMetrics) RecordHandlerAttempt(ctx context.Context, handlerName string) {
	attrs := handlerNameAttrs(handlerName)
	m.handlerAttemptCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
}

func (m *roundMetrics) RecordHandlerDuration(ctx context.Context, handlerName string, duration time.Duration) {
	attrs := handlerNameAttrs(handlerName)
	m.handlerDurationHistogram.Record(ctx, duration.Seconds(), metric.WithAttributes(attrs...))
}

func (m *roundMetrics) RecordHandlerFailure(ctx context.Context, handlerName string) {
	attrs := handlerNameAttrs(handlerName)
	m.handlerFailureCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
}

func (m *roundMetrics) RecordHandlerSuccess(ctx context.Context, handlerName string) {
	attrs := handlerNameAttrs(handlerName)
	m.handlerSuccessCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
}
