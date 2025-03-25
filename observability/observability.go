package observability

import (
	"context"
	"fmt"

	lokifrolfbot "github.com/Black-And-White-Club/frolf-bot-shared/observability/loki"
	prometheusfrolfbot "github.com/Black-And-White-Club/frolf-bot-shared/observability/prometheus"
	tempofrolfbot "github.com/Black-And-White-Club/frolf-bot-shared/observability/tempo"
)

// Config holds all configuration for observability components
type Config struct {
	LokiURL         string
	LokiTenantID    string
	ServiceName     string
	ServiceVersion  string
	Environment     string
	MetricsAddress  string
	TempoEndpoint   string
	TempoInsecure   bool
	TempoSampleRate float64
}

// Observability defines the unified interface for observability components
type Observability interface {
	GetLogger() lokifrolfbot.Logger
	GetTracer() tempofrolfbot.Tracer
	GetMetrics() prometheusfrolfbot.Metrics
	RegisterHealthChecker(checker prometheusfrolfbot.HealthChecker)
	PerformHealthChecks(ctx context.Context) map[string]error
	HealthCheck(ctx context.Context) error
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

// ObservabilityService provides a unified service for all observability needs
type ObservabilityService struct {
	logger         lokifrolfbot.Logger
	tracer         tempofrolfbot.Tracer
	metrics        prometheusfrolfbot.Metrics
	tracerShutdown func() // Store the tracer shutdown function
}

// NewObservability creates and configures all observability components
func NewObservability(ctx context.Context, config Config) (*ObservabilityService, error) {
	// Initialize Loki logger
	logger, err := lokifrolfbot.NewLokiLogger(
		config.LokiURL,
		config.LokiTenantID,
		config.ServiceName,
		config.Environment,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	// Initialize Tempo tracer
	tracer, tracerShutdown, err := tempofrolfbot.NewTracer(
		config.ServiceName,
		config.TempoEndpoint,
		config.TempoInsecure,
		config.TempoSampleRate,
		config.ServiceVersion,
		config.Environment,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize tracer: %w", err)
	}

	// Initialize Prometheus metrics
	metrics, err := prometheusfrolfbot.NewPrometheusMetrics(
		config.MetricsAddress,
		config.ServiceName,
		config.Environment,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize metrics: %w", err)
	}

	return &ObservabilityService{
		logger:         logger,
		tracer:         tracer,
		metrics:        metrics,
		tracerShutdown: tracerShutdown,
	}, nil
}

// Start initializes all observability components
func (o *ObservabilityService) Start(ctx context.Context) error {
	o.logger.Info("Starting observability components")
	if err := o.metrics.Start(ctx); err != nil {
		return err
	}
	// Additional start logic for tracer if needed
	return nil
}

// Stop gracefully shuts down all observability components
func (o *ObservabilityService) Stop(ctx context.Context) error {
	o.logger.Info("Stopping observability components")
	if o.tracerShutdown != nil {
		o.tracerShutdown()
	}
	return o.metrics.Stop(ctx)
}

// GetLogger returns the logger instance
func (o *ObservabilityService) GetLogger() lokifrolfbot.Logger {
	return o.logger
}

// GetTracer returns the tracer instance
func (o *ObservabilityService) GetTracer() tempofrolfbot.Tracer {
	return o.tracer
}

// GetMetrics returns the metrics provider instance
func (o *ObservabilityService) GetMetrics() prometheusfrolfbot.Metrics {
	return o.metrics
}

// RegisterHealthChecker adds a health checker to the metrics
func (o *ObservabilityService) RegisterHealthChecker(checker prometheusfrolfbot.HealthChecker) {
	o.metrics.RegisterHealthChecker(checker)
}

// PerformHealthChecks runs health checks and returns results
func (o *ObservabilityService) PerformHealthChecks(ctx context.Context) map[string]error {
	return o.metrics.PerformHealthChecks(ctx)
}

// HealthCheck performs health checks on all observability components
func (o *ObservabilityService) HealthCheck(ctx context.Context) error {
	// Check metrics server health
	if err := o.metrics.HealthCheck(ctx); err != nil {
		return fmt.Errorf("metrics health check failed: %w", err)
	}
	return nil
}
