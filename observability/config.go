// shared/observability/config.go
package observability

import (
	"log/slog"

	"go.opentelemetry.io/otel/attribute"
)

type Config struct {
	ServiceName string
	Environment string
	Version     string

	LokiURL         string
	MetricsAddress  string
	TempoEndpoint   string
	TempoInsecure   bool
	TempoSampleRate float64
}

func (c Config) TracingEnabled() bool {
	return c.TempoEndpoint != ""
}

func (c Config) MetricsEnabled() bool {
	return c.MetricsAddress != ""
}

func parseLogLevel(c Config) slog.Level {
	switch c.Environment {
	case "prod":
		return slog.LevelInfo
	default:
		return slog.LevelDebug
	}
}

func (c Config) ResourceAttributes() []attribute.KeyValue {
	return []attribute.KeyValue{
		attribute.String("service.name", c.ServiceName),
		attribute.String("service.version", c.Version),
		attribute.String("deployment.environment", c.Environment),
	}
}
