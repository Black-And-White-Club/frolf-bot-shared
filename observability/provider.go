package observability

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"

	// OTEL exporters (gRPC only)
	otlploggrpc "go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	otlpmetricgrpc "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	otlptracegrpc "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"

	// OTEL SDK
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type Provider struct {
	TracerProvider trace.TracerProvider
	MeterProvider  metric.MeterProvider
	Logger         *slog.Logger
	Shutdown       func(ctx context.Context) error
}

func Setup(ctx context.Context, cfg Config) (*Provider, error) {
	// Create resource with attributes
	res, err := resource.New(ctx, resource.WithAttributes(cfg.ResourceAttributes()...))
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	shutdownFuncs := []func(context.Context) error{}

	// ========== Tracer ==========
	var tracerProvider trace.TracerProvider = trace.NewNoopTracerProvider()

	if cfg.TracingEnabled() {
		tp, shutdown, err := setupTracing(ctx, cfg, res)
		if err != nil {
			return nil, fmt.Errorf("failed to setup tracing: %w", err)
		}
		tracerProvider = tp
		shutdownFuncs = append(shutdownFuncs, shutdown)
	}

	otel.SetTracerProvider(tracerProvider)

	// ========== Metrics ==========
	var meterProvider metric.MeterProvider = noop.NewMeterProvider()

	if cfg.MetricsEnabled() {
		mp, shutdown, err := setupMetrics(ctx, cfg, res)
		if err != nil {
			return nil, fmt.Errorf("failed to setup metrics: %w", err)
		}
		meterProvider = mp
		shutdownFuncs = append(shutdownFuncs, shutdown)
	}

	otel.SetMeterProvider(meterProvider)

	// ========== Logging ==========
	logger, logShutdown, err := setupLogging(ctx, cfg, res)
	if err != nil {
		return nil, fmt.Errorf("failed to setup logging: %w", err)
	}
	shutdownFuncs = append(shutdownFuncs, logShutdown)

	// Combined shutdown function
	shutdown := func(ctx context.Context) error {
		for _, fn := range shutdownFuncs {
			if err := fn(ctx); err != nil {
				return err
			}
		}
		return nil
	}

	return &Provider{
		TracerProvider: tracerProvider,
		MeterProvider:  meterProvider,
		Logger:         logger,
		Shutdown:       shutdown,
	}, nil
}

func setupLogging(ctx context.Context, cfg Config, res *resource.Resource) (*slog.Logger, func(context.Context) error, error) {
	if cfg.LokiEnabled() {
		// Use OTLP log exporter to send to Alloy, which forwards to Loki
		return setupOTLPLogging(ctx, cfg, res)
	}

	// Fallback to stdout logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: parseLogLevel(cfg),
	}))
	logger = logger.With("service", cfg.ServiceName, "version", cfg.Version, "environment", cfg.Environment)

	return logger, func(context.Context) error { return nil }, nil
}

func setupOTLPLogging(ctx context.Context, cfg Config, res *resource.Resource) (*slog.Logger, func(context.Context) error, error) {
	endpoint := firstNonEmpty(cfg.OTLPEndpoint, cfg.MetricsAddress)
	exporter, err := otlploggrpc.New(ctx, otlploggrpc.WithEndpoint(endpoint), otlploggrpc.WithInsecure())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create OTLP log exporter: %w", err)
	}

	// Create log processor
	processor := sdklog.NewBatchProcessor(exporter)

	// Create logger provider
	loggerProvider := sdklog.NewLoggerProvider(
		sdklog.WithResource(res),
		sdklog.WithProcessor(processor),
	)

	// Create slog handler that uses OTEL
	handler := NewOTELHandler(loggerProvider, cfg)
	logger := slog.New(handler)

	shutdown := func(ctx context.Context) error {
		return loggerProvider.Shutdown(ctx)
	}

	return logger, shutdown, nil
}

// Custom slog handler that sends to OTEL
type otelHandler struct {
	loggerProvider *sdklog.LoggerProvider
	level          slog.Level
	attrs          []slog.Attr
}

func NewOTELHandler(provider *sdklog.LoggerProvider, cfg Config) slog.Handler {
	return &otelHandler{
		loggerProvider: provider,
		level:          parseLogLevel(cfg),
	}
}

func (h *otelHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *otelHandler) Handle(ctx context.Context, record slog.Record) error {
	// Get logger from provider
	logger := h.loggerProvider.Logger("frolf-bot")

	// Convert slog to OTEL log record
	var logRecord log.Record
	logRecord.SetTimestamp(record.Time)
	logRecord.SetBody(log.StringValue(record.Message))
	logRecord.SetSeverity(convertLevel(record.Level))

	// Add attributes
	record.Attrs(func(attr slog.Attr) bool {
		logRecord.AddAttributes(log.String(attr.Key, attr.Value.String()))
		return true
	})

	// Add handler attrs
	for _, attr := range h.attrs {
		logRecord.AddAttributes(log.String(attr.Key, attr.Value.String()))
	}

	// Emit the log (no return value)
	logger.Emit(ctx, logRecord)
	return nil
}

func (h *otelHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newAttrs := make([]slog.Attr, len(h.attrs)+len(attrs))
	copy(newAttrs, h.attrs)
	copy(newAttrs[len(h.attrs):], attrs)

	return &otelHandler{
		loggerProvider: h.loggerProvider,
		level:          h.level,
		attrs:          newAttrs,
	}
}

func (h *otelHandler) WithGroup(name string) slog.Handler {
	// For simplicity, ignore groups for now
	return h
}

func convertLevel(level slog.Level) log.Severity {
	switch level {
	case slog.LevelDebug:
		return log.SeverityDebug
	case slog.LevelInfo:
		return log.SeverityInfo
	case slog.LevelWarn:
		return log.SeverityWarn
	case slog.LevelError:
		return log.SeverityError
	default:
		return log.SeverityInfo
	}
}

func setupTracing(ctx context.Context, cfg Config, res *resource.Resource) (trace.TracerProvider, func(context.Context) error, error) {
	endpoint := firstNonEmpty(cfg.OTLPEndpoint, cfg.TempoEndpoint)
	if endpoint == "" {
		// tracing disabled
		return trace.NewNoopTracerProvider(), func(context.Context) error { return nil }, nil
	}
	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithEndpoint(endpoint), otlptracegrpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	return tp, tp.Shutdown, nil
}

func setupMetrics(ctx context.Context, cfg Config, res *resource.Resource) (metric.MeterProvider, func(context.Context) error, error) {
	endpoint := firstNonEmpty(cfg.OTLPEndpoint, cfg.MetricsAddress)
	if endpoint == "" {
		// metrics disabled
		return noop.NewMeterProvider(), func(context.Context) error { return nil }, nil
	}
	exporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithEndpoint(endpoint), otlpmetricgrpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter)),
		sdkmetric.WithResource(res),
	)

	return mp, mp.Shutdown, nil
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}
