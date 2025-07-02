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

	// Example: Guild Created Counter
	m.guildCreatedCounter, err = meter.Int64Counter(
		metricName("created_total"),
		metric.WithDescription("Number of guilds created"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	// Add more metric initializations as needed

	return m, nil
}
