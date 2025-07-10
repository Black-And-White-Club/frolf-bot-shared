// observability/opentelemetry/guild/constructor.go
package guildmetrics

import (
	"go.opentelemetry.io/otel/metric"
)

// NewGuildMetrics creates a new GuildMetrics implementation using OpenTelemetry.
// It requires an OTEL Meter instance and a prefix for metric names.
func NewGuildMetrics(meter metric.Meter, prefix string) (GuildMetrics, error) {
	metricName := func(name string) string {
		if prefix != "" {
			return prefix + "_guild_" + name
		}
		return "guild_" + name
	}

	var err error
	m := &guildMetrics{meter: meter}

	// Guild Created Counter
	m.guildCreatedCounter, err = meter.Int64Counter(
		metricName("created_total"),
		metric.WithDescription("Number of guilds created"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	// Guild Deleted Counter
	m.guildDeletedCounter, err = meter.Int64Counter(
		metricName("deleted_total"),
		metric.WithDescription("Number of guilds deleted"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	// Operation metrics
	m.operationAttemptCounter, err = meter.Int64Counter(
		metricName("operation_attempts_total"),
		metric.WithDescription("Number of guild operation attempts"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	m.operationSuccessCounter, err = meter.Int64Counter(
		metricName("operation_success_total"),
		metric.WithDescription("Number of successful guild operations"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	m.operationFailureCounter, err = meter.Int64Counter(
		metricName("operation_failure_total"),
		metric.WithDescription("Number of failed guild operations"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	m.operationDuration, err = meter.Float64Histogram(
		metricName("operation_duration_seconds"),
		metric.WithDescription("Duration of guild operations in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, err
	}

	// Handler metrics
	m.handlerAttemptCounter, err = meter.Int64Counter(
		metricName("handler_attempts_total"),
		metric.WithDescription("Number of guild handler attempts"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	m.handlerSuccessCounter, err = meter.Int64Counter(
		metricName("handler_success_total"),
		metric.WithDescription("Number of successful guild handler executions"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	m.handlerFailureCounter, err = meter.Int64Counter(
		metricName("handler_failure_total"),
		metric.WithDescription("Number of failed guild handler executions"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	m.handlerDuration, err = meter.Float64Histogram(
		metricName("handler_duration_seconds"),
		metric.WithDescription("Duration of guild handler executions in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, err
	}

	return m, nil
}
