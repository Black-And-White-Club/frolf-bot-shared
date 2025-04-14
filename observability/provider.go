package observability

import (
	"context"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
)

type Provider struct {
	TracerProvider trace.TracerProvider
	MeterProvider  metric.MeterProvider
	Logger         *slog.Logger
	Shutdown       func(ctx context.Context) error
}

func Setup(ctx context.Context, cfg Config) (*Provider, error) {
	var shutdownFns []func(ctx context.Context) error

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: parseLogLevel(cfg), // helper function
	}))

	res, err := resource.New(ctx,
		resource.WithAttributes(cfg.ResourceAttributes()...),
		resource.WithFromEnv(),
		resource.WithTelemetrySDK(),
	)
	if err != nil {
		return nil, err
	}

	// ========== Tracer ==========
	var tracerProvider trace.TracerProvider = trace.NewNoopTracerProvider()
	if cfg.TracingEnabled() {
		tp, shutdown, err := setupTracing(ctx, cfg, res)
		if err != nil {
			return nil, err
		}
		tracerProvider = tp
		shutdownFns = append(shutdownFns, shutdown)
	}

	// ========== Metrics ==========
	var meterProvider metric.MeterProvider = noop.NewMeterProvider()
	if cfg.MetricsEnabled() {
		mp, shutdown, err := setupMetrics(ctx, cfg, res)
		if err != nil {
			return nil, err
		}
		meterProvider = mp
		shutdownFns = append(shutdownFns, shutdown)
	}

	// You can forward logs to OTEL here in future via slog OTEL handler or bridge

	return &Provider{
		Logger:         logger,
		TracerProvider: tracerProvider,
		MeterProvider:  meterProvider,
		Shutdown: func(ctx context.Context) error {
			for _, fn := range shutdownFns {
				_ = fn(ctx) // log errors if needed
			}
			return nil
		},
	}, nil
}

func setupTracing(ctx context.Context, cfg Config, res *resource.Resource) (trace.TracerProvider, func(context.Context) error, error) {
	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(cfg.TempoEndpoint),
		otlptracegrpc.WithInsecure(), // optional: toggle for prod
	)
	if err != nil {
		return nil, nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tp)
	return tp, tp.Shutdown, nil
}

func setupMetrics(ctx context.Context, cfg Config, res *resource.Resource) (metric.MeterProvider, func(context.Context) error, error) {
	exporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpoint(cfg.MetricsAddress),
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		return nil, nil, err
	}

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter)),
		sdkmetric.WithResource(res),
	)

	otel.SetMeterProvider(mp)
	return mp, func(ctx context.Context) error {
		return exporter.Shutdown(ctx)
	}, nil
}
