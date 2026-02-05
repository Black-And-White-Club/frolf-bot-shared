package clubmetrics

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/metric"
)

func (m *clubMetrics) RecordOperationAttempt(ctx context.Context, operationName, serviceName string) {
	m.operationAttemptCounter.Add(ctx, 1, metric.WithAttributes(
		operationAttr(operationName),
		serviceAttr(serviceName),
	))
}

func (m *clubMetrics) RecordOperationSuccess(ctx context.Context, operationName, serviceName string) {
	m.operationSuccessCounter.Add(ctx, 1, metric.WithAttributes(
		operationAttr(operationName),
		serviceAttr(serviceName),
	))
}

func (m *clubMetrics) RecordOperationFailure(ctx context.Context, operationName, serviceName string) {
	m.operationFailureCounter.Add(ctx, 1, metric.WithAttributes(
		operationAttr(operationName),
		serviceAttr(serviceName),
	))
}

func (m *clubMetrics) RecordOperationDuration(ctx context.Context, operationName, serviceName string, duration time.Duration) {
	m.operationDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(
		operationAttr(operationName),
		serviceAttr(serviceName),
	))
}
