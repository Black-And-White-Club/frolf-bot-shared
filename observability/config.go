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
	OTLPEndpoint    string
	OTLPTransport   string // grpc|http

	// OTEL log batching (optional; zeros use sensible defaults)
	LogBatchMaxQueueSize       int // e.g., 256 dev, 2048 prod
	LogBatchMaxExportBatchSize int // e.g., 64 dev, 512 prod
	LogBatchTimeoutSeconds     int // e.g., 2 dev, 5 prod
	LogExportTimeoutSeconds    int // e.g., 3 dev, 10 prod
}

func (c Config) LokiEnabled() bool {
	return c.LokiURL != ""
}

func (c Config) TracingEnabled() bool {
	return c.TempoEndpoint != "" || c.OTLPEndpoint != ""
}

func (c Config) MetricsEnabled() bool {
	return c.MetricsAddress != "" || c.OTLPEndpoint != ""
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
