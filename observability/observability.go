package observability

import (
	"context"
	"fmt"
	"log/slog"
)

// Observability defines a combined interface for logging, metrics, and tracing.
type Observability interface {
	GetLogger() Logger
	GetMetrics() Metrics
	GetTracer() Tracer
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

// Config holds configuration for the observability components.  This is good!
type Config struct {
	LokiURL         string
	LokiTenantID    string
	ServiceName     string
	MetricsAddress  string
	TempoEndpoint   string
	TempoInsecure   bool
	ServiceVersion  string
	TempoSampleRate float64
	// ... other configuration options ...
}

// concreteObservability is the concrete implementation of the Observability interface.
type concreteObservability struct {
	logger        Logger
	metrics       Metrics
	tracer        Tracer
	metricsServer *MetricsServer // Hold the server here
}

// NewObservability creates and initializes the observability components. This is the factory.
func NewObservability(ctx context.Context, config Config) (Observability, error) {
	logger, err := NewLokiLogger(config.LokiURL, config.LokiTenantID, config.ServiceName)
	if err != nil {
		return nil, fmt.Errorf("failed to create Loki logger: %w", err)
	}

	metricsBuilder := NewPrometheusMetricsBuilder(config.ServiceName, "subsystem") // Replace "subsystem" with your desired value

	metricsServer, err := NewMetricsServer(config.MetricsAddress, metricsBuilder.Registry)
	if err != nil {
		return nil, fmt.Errorf("failed to create metrics server: %w", err)
	}

	tracer := &TempoTracer{} // Create an instance of the concrete TempoTracer
	tracingOpts := TracingOptions{
		ServiceName:    config.ServiceName,
		TempoEndpoint:  config.TempoEndpoint,
		Insecure:       config.TempoInsecure,
		ServiceVersion: config.ServiceVersion,
		SampleRate:     config.TempoSampleRate,
	}
	tracerShutdown, err := tracer.InitTracing(ctx, tracingOpts) // Initialize tracing
	if err != nil {
		return nil, fmt.Errorf("failed to initialize tracing: %w", err)
	}
	_ = tracerShutdown //keep reference to the tracer shutdown.

	return &concreteObservability{
		logger:        logger,
		metrics:       metricsBuilder,
		tracer:        tracer,
		metricsServer: metricsServer,
	}, nil
}

// GetLogger returns the Logger instance.
func (o *concreteObservability) GetLogger() Logger {
	return o.logger
}

// GetMetrics returns the Metrics instance.
func (o *concreteObservability) GetMetrics() Metrics {
	return o.metrics
}

// GetTracer returns the Tracer instance.
func (o *concreteObservability) GetTracer() Tracer {
	return o.tracer
}

// Start initializes the observability components.
func (o *concreteObservability) Start(ctx context.Context) error {
	o.metricsServer.Start(ctx)
	return nil // Currently, only the metrics server needs starting.  If you had other startup tasks, add them here.
}

// Stop shuts down the observability components.
func (o *concreteObservability) Stop(ctx context.Context) error {
	o.logger.Shutdown()                                   // Shutdown the logger
	if err := o.metricsServer.Shutdown(ctx); err != nil { // Shutdown the metrics server
		slog.Error("Failed to shut down metrics server", "error", err) //best effort.
	}
	if closer, ok := o.tracer.(interface{ Shutdown(context.Context) error }); ok { //check to make sure interface has shutdown method
		if err := closer.Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shut down tracer: %w", err)
		}
	}

	return nil
}
