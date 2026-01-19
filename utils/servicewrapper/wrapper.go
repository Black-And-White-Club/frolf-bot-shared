// Package servicewrapper provides a generic telemetry wrapper for service operations.
//
// This package standardizes the observability pattern across all service layers:
//   - Automatic span creation and management
//   - Metrics recording (attempts, success, failure, duration)
//   - Structured logging with correlation IDs
//   - Panic recovery with proper error reporting
//
// Usage:
//
//	wrapper := servicewrapper.New(tracer, logger, metrics)
//	result, err := wrapper.Execute(ctx, "CreateUser", userID, func(ctx context.Context) (results.OperationResult, error) {
//	    // Your business logic here
//	})
package servicewrapper

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Black-And-White-Club/frolf-bot-shared/observability/attr"
	"github.com/Black-And-White-Club/frolf-bot-shared/utils/results"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// ServiceMetrics defines the metrics interface for service operations.
// Each domain module should implement this interface with domain-specific metrics.
type ServiceMetrics interface {
	RecordOperationAttempt(ctx context.Context, operation string, resourceID string)
	RecordOperationSuccess(ctx context.Context, operation string, resourceID string)
	RecordOperationFailure(ctx context.Context, operation string, resourceID string)
	RecordOperationDuration(ctx context.Context, operation string, duration time.Duration, resourceID string)
}

// OperationFunc is the signature for wrapped service operations.
type OperationFunc func(ctx context.Context) (results.OperationResult, error)

// Wrapper provides telemetry wrapping for service operations.
type Wrapper struct {
	tracer  trace.Tracer
	logger  *slog.Logger
	metrics ServiceMetrics
}

// New creates a new service wrapper with the provided observability components.
// All parameters are required; use no-op implementations for testing.
func New(tracer trace.Tracer, logger *slog.Logger, metrics ServiceMetrics) *Wrapper {
	return &Wrapper{
		tracer:  tracer,
		logger:  logger,
		metrics: metrics,
	}
}

// Execute wraps a service operation with telemetry.
//
// It provides:
//   - Automatic span creation with operation name and resource ID
//   - Metrics recording for attempts, success, failure, and duration
//   - Structured logging with correlation ID extraction
//   - Panic recovery that converts panics to errors
//
// The resourceID is typically the primary identifier for the resource being operated on
// (e.g., guildID, userID, roundID).
func (w *Wrapper) Execute(
	ctx context.Context,
	operationName string,
	resourceID string,
	op OperationFunc,
) (result results.OperationResult, err error) {
	// Start span
	ctx, span := w.tracer.Start(ctx, operationName, trace.WithAttributes(
		attribute.String("operation", operationName),
		attribute.String("resource_id", resourceID),
	))
	defer span.End()

	// Record attempt
	w.metrics.RecordOperationAttempt(ctx, operationName, resourceID)

	// Track duration
	startTime := time.Now()
	defer func() {
		w.metrics.RecordOperationDuration(ctx, operationName, time.Since(startTime), resourceID)
	}()

	// Log operation start
	w.logger.InfoContext(ctx, "Operation started",
		attr.ExtractCorrelationID(ctx),
		attr.String("operation", operationName),
		attr.String("resource_id", resourceID),
	)

	// Panic recovery
	defer func() {
		if r := recover(); r != nil {
			errorMsg := fmt.Sprintf("panic in %s: %v", operationName, r)
			w.logger.ErrorContext(ctx, errorMsg,
				attr.ExtractCorrelationID(ctx),
				attr.String("resource_id", resourceID),
				attr.Any("panic", r),
			)
			w.metrics.RecordOperationFailure(ctx, operationName, resourceID)
			span.RecordError(fmt.Errorf("%s", errorMsg))

			result = results.OperationResult{}
			err = fmt.Errorf("%s", errorMsg)
		}
	}()

	// Execute the operation
	result, err = op(ctx)
	if err != nil {
		wrappedErr := fmt.Errorf("%s failed: %w", operationName, err)
		w.logger.ErrorContext(ctx, "Operation failed",
			attr.ExtractCorrelationID(ctx),
			attr.String("operation", operationName),
			attr.String("resource_id", resourceID),
			attr.Error(wrappedErr),
		)
		w.metrics.RecordOperationFailure(ctx, operationName, resourceID)
		span.RecordError(wrappedErr)
		return result, wrappedErr
	}

	// Log success
	w.logger.InfoContext(ctx, "Operation completed",
		attr.ExtractCorrelationID(ctx),
		attr.String("operation", operationName),
		attr.String("resource_id", resourceID),
	)
	w.metrics.RecordOperationSuccess(ctx, operationName, resourceID)

	return result, nil
}

// NoOpMetrics provides a no-operation metrics implementation for testing.
type NoOpMetrics struct{}

func (NoOpMetrics) RecordOperationAttempt(ctx context.Context, operation string, resourceID string) {}
func (NoOpMetrics) RecordOperationSuccess(ctx context.Context, operation string, resourceID string) {}
func (NoOpMetrics) RecordOperationFailure(ctx context.Context, operation string, resourceID string) {}
func (NoOpMetrics) RecordOperationDuration(ctx context.Context, operation string, duration time.Duration, resourceID string) {
}
