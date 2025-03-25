package prometheusfrolfbot

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	databasemetrics "github.com/Black-And-White-Club/frolf-bot-shared/observability/prometheus/database"
	eventbusmetrics "github.com/Black-And-White-Club/frolf-bot-shared/observability/prometheus/eventbus"
	usermetrics "github.com/Black-And-White-Club/frolf-bot-shared/observability/prometheus/user"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// HealthStatus defines the possible health states
type HealthStatus string

const (
	StatusHealthy   HealthStatus = "healthy"
	StatusUnhealthy HealthStatus = "unhealthy"
)

// HealthChecker interface for components that can be health-checked
type HealthChecker interface {
	Check(ctx context.Context) error
	Name() string
}

// PrometheusMetrics provides the Prometheus metrics implementation
type PrometheusMetrics struct {
	registry    *prometheus.Registry
	server      *http.Server
	prefix      string
	environment string

	// Cached metrics instances
	userMetrics     usermetrics.UserMetrics
	databaseMetrics databasemetrics.DatabaseMetrics
	eventBusMetrics eventbusmetrics.EventBusMetrics

	// Health checking
	healthMu       sync.RWMutex
	healthCheckers map[string]HealthChecker

	// Health metrics
	serviceHealthy *prometheus.GaugeVec
	lastCheckTime  *prometheus.GaugeVec
}

// Metrics defines the interface for the metrics component.
type Metrics interface {
	UserMetrics() usermetrics.UserMetrics
	DatabaseMetrics() databasemetrics.DatabaseMetrics
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	GetMetrics() (usermetrics.UserMetrics, databasemetrics.DatabaseMetrics, eventbusmetrics.EventBusMetrics)
	HealthCheck(ctx context.Context) error
	RegisterHealthChecker(checker HealthChecker)
	PerformHealthChecks(ctx context.Context) map[string]error
	EventBusMetrics() eventbusmetrics.EventBusMetrics
}

// RegisterHealthChecker adds a health checker to the metrics
func (p *PrometheusMetrics) RegisterHealthChecker(checker HealthChecker) {
	p.healthMu.Lock()
	defer p.healthMu.Unlock()
	p.healthCheckers[checker.Name()] = checker
}

// NewPrometheusMetrics creates a new Prometheus metrics provider
func NewPrometheusMetrics(
	metricsPort string,
	servicePrefix string,
	environment string,
) (*PrometheusMetrics, error) {
	// Create a new registry
	registry := prometheus.NewRegistry()

	// Initialize health metrics
	serviceHealthy := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: servicePrefix,
			Name:      "health_status",
			Help:      "Health status of service components (1=healthy, 0=unhealthy)",
		},
		[]string{"component", "environment"},
	)

	lastCheckTime := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: servicePrefix,
			Name:      "health_check_timestamp",
			Help:      "Timestamp of the last health check",
		},
		[]string{"component", "environment"},
	)

	registry.MustRegister(serviceHealthy, lastCheckTime)

	// Create Prometheus HTTP server (for metrics only)
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	p := &PrometheusMetrics{
		server: &http.Server{
			Addr:         ":" + metricsPort,
			Handler:      mux,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		registry:       registry,
		prefix:         servicePrefix,
		environment:    environment,
		healthCheckers: make(map[string]HealthChecker),
		serviceHealthy: serviceHealthy,
		lastCheckTime:  lastCheckTime,
	}

	return p, nil
}

// Start initiates the metrics server
func (p *PrometheusMetrics) Start(ctx context.Context) error {
	go func() {
		if err := p.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// Can't log here since we don't have the logger
		}
	}()
	return nil
}

// Stop gracefully shuts down the metrics server
func (p *PrometheusMetrics) Stop(ctx context.Context) error {
	return p.server.Shutdown(ctx)
}

// UserMetrics returns the user metrics
func (p *PrometheusMetrics) UserMetrics() usermetrics.UserMetrics {
	if p.userMetrics == nil {
		p.userMetrics = usermetrics.NewUserMetrics(p.registry, p.prefix)
	}
	return p.userMetrics
}

// DatabaseMetrics returns the database metrics interface
func (p *PrometheusMetrics) DatabaseMetrics() databasemetrics.DatabaseMetrics {
	if p.databaseMetrics == nil {
		p.databaseMetrics = databasemetrics.NewDatabaseMetrics(p.registry, p.prefix)
	}
	return p.databaseMetrics
}

// EventBusMetrics returns the eventbus metrics interface
func (p *PrometheusMetrics) EventBusMetrics() eventbusmetrics.EventBusMetrics {
	if p.eventBusMetrics == nil {
		p.eventBusMetrics = eventbusmetrics.NewEventBusMetrics(p.registry, p.prefix)
	}
	return p.eventBusMetrics
}

// GetMetrics returns the user and database metrics
func (p *PrometheusMetrics) GetMetrics() (usermetrics.UserMetrics, databasemetrics.DatabaseMetrics, eventbusmetrics.EventBusMetrics) {
	return p.UserMetrics(), p.DatabaseMetrics(), p.EventBusMetrics()
}

// HealthCheck performs a health check on the metrics server and all registered checkers
func (p *PrometheusMetrics) HealthCheck(ctx context.Context) error {
	if p.server == nil {
		return fmt.Errorf("metrics server is not initialized")
	}

	// Run all health checks
	results := p.PerformHealthChecks(ctx)

	// Check if any failed
	for name, err := range results {
		if err != nil {
			return fmt.Errorf("health check failed for %s: %w", name, err)
		}
	}

	return nil
}

// PerformHealthChecks runs all registered health checks and updates metrics
func (p *PrometheusMetrics) PerformHealthChecks(ctx context.Context) map[string]error {
	results := make(map[string]error)

	p.healthMu.RLock()
	checkers := make(map[string]HealthChecker, len(p.healthCheckers))
	for name, checker := range p.healthCheckers {
		checkers[name] = checker
	}
	p.healthMu.RUnlock()

	// Run checks without holding the lock
	for name, checker := range checkers {
		err := checker.Check(ctx)
		results[name] = err

		// Update metrics
		if err != nil {
			p.serviceHealthy.WithLabelValues(name, p.environment).Set(0)
		} else {
			p.serviceHealthy.WithLabelValues(name, p.environment).Set(1)
		}
		p.lastCheckTime.WithLabelValues(name, p.environment).Set(float64(time.Now().Unix()))
	}

	return results
}
