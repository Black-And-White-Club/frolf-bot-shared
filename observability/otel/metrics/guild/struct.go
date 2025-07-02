// observability/opentelemetry/guild/struct.go
package guildmetrics

import (
	"go.opentelemetry.io/otel/metric"
)

// guildMetrics implements GuildMetrics using OpenTelemetry
type guildMetrics struct {
	meter metric.Meter // OTEL Meter

	guildCreatedCounter     metric.Int64Counter
	guildDeletedCounter     metric.Int64Counter
	operationAttemptCounter metric.Int64Counter
	operationSuccessCounter metric.Int64Counter
	operationFailureCounter metric.Int64Counter
	operationDuration       metric.Float64Histogram
}
