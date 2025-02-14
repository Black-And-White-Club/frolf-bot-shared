package observability

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/ThreeDotsLabs/watermill/components/metrics"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics defines the interface for metrics collection.  Moved here!
type Metrics interface {
	AddPrometheusRouterMetrics(r *message.Router) error
}

// PrometheusMetricsBuilder provides methods to decorate publishers, subscribers, and handlers.
//
//	It now implements the Metrics interface.
type PrometheusMetricsBuilder struct {
	*metrics.PrometheusMetricsBuilder
	Registry *prometheus.Registry
}

// NewPrometheusMetricsBuilder initializes a new PrometheusMetricsBuilder.
func NewPrometheusMetricsBuilder(namespace, subsystem string) *PrometheusMetricsBuilder {
	// Create a *new* registry for your application.  Don't use the global one.
	registry := prometheus.NewRegistry()

	// Create the Watermill metrics builder, using *your* registry.
	watermillMetrics := metrics.NewPrometheusMetricsBuilder(registry, namespace, subsystem)

	builder := &PrometheusMetricsBuilder{PrometheusMetricsBuilder: &watermillMetrics, Registry: registry}

	//Better nil check.
	if builder.PrometheusMetricsBuilder == nil {
		panic("watermill PrometheusMetricsBuilder is nil") //panic as it should not ever happen.
	}
	return builder
}

// AddPrometheusRouterMetrics adds metrics middleware to the router.
func (b *PrometheusMetricsBuilder) AddPrometheusRouterMetrics(r *message.Router) error {
	b.PrometheusMetricsBuilder.AddPrometheusRouterMetrics(r)
	return nil
}

// MetricsServer serves Prometheus metrics on the given address.
type MetricsServer struct {
	server *http.Server
}

// NewMetricsServer creates and starts a new Prometheus metrics server.
func NewMetricsServer(addr string, registry *prometheus.Registry) (*MetricsServer, error) {
	// Create a ServeMux to handle different paths (good practice)
	mux := http.NewServeMux()
	// Use promhttp package, and register the registry FOR the handler,
	mux.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	server := &http.Server{
		Addr:    addr,
		Handler: mux, // Use the ServeMux
		// Good practice to set timeouts:
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second, // Adjust as needed
	}

	return &MetricsServer{server: server}, nil
}

// Start starts the metrics server in a goroutine.
func (m *MetricsServer) Start(ctx context.Context) { // Added context
	go func() {
		slog.Info("Starting Prometheus metrics server on %s", "address", m.server.Addr)
		if err := m.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Prometheus metrics server failed", "error", err)
		}
	}()
}

// Shutdown gracefully shuts down the metrics server.
func (m *MetricsServer) Shutdown(ctx context.Context) error {
	slog.Info("Shutting down Prometheus metrics server")
	if err := m.server.Shutdown(ctx); err != nil {
		// Already good!  You're logging and returning the error.
		return fmt.Errorf("error during prometheus metrics server shutdown: %w", err)
	}
	return nil
}

// IsRunning checks if the metrics server is running.
func (m *MetricsServer) IsRunning() bool {
	return m.server != nil // Correct
}
