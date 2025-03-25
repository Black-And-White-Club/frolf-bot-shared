package lokifrolfbot

import (
	"context"
	"log/slog"
	"time"

	"github.com/Black-And-White-Club/frolf-bot-shared/observability/attr"
	"github.com/grafana/loki-client-go/loki"
	slogloki "github.com/samber/slog-loki/v3"
)

// Logger defines the interface for structured logging.
type Logger interface {
	Debug(msg string, attrs ...attr.LogAttr)
	Info(msg string, attrs ...attr.LogAttr)
	Warn(msg string, attrs ...attr.LogAttr)
	Error(msg string, attrs ...attr.LogAttr)
	Close()
	WithContext(ctx context.Context) Logger
}

// LokiLogger wraps a Loki handler and implements the Logger interface.
type LokiLogger struct {
	handler *slogloki.LokiHandler
	client  *loki.Client
}

// NewLokiLogger initializes a Loki logger instance.
func NewLokiLogger(url, tenantID, serviceName, environment string) (*LokiLogger, error) {
	config, err := loki.NewDefaultConfig(url)
	if err != nil {
		return nil, err
	}
	config.TenantID = tenantID

	client, err := loki.New(config)
	if err != nil {
		return nil, err
	}

	handler := slogloki.Option{
		Level:  slog.LevelDebug,
		Client: client,
	}.NewLokiHandler().(*slogloki.LokiHandler)

	handler = handler.WithAttrs([]slog.Attr{
		slog.String("service", serviceName),
	}).(*slogloki.LokiHandler)

	return &LokiLogger{
		handler: handler,
		client:  client,
	}, nil
}

// Close ensures all logs are flushed before shutting down.
func (l *LokiLogger) Close() {
	if l.client != nil {
		l.client.Stop()
	}
}

// WithContext injects correlation ID into logs.
func (l *LokiLogger) WithContext(ctx context.Context) Logger {
	if ctx == nil {
		ctx = context.Background()
	}

	correlationID := attr.ExtractCorrelationID(ctx)

	newHandler := l.handler.WithAttrs([]slog.Attr{
		correlationID.ToSlogAttr(),
	}).(*slogloki.LokiHandler)

	return &LokiLogger{
		handler: newHandler,
		client:  l.client,
	}
}

// Debug, Info, Warn, Error are implemented for structured logging.
func (l *LokiLogger) Debug(msg string, attrs ...attr.LogAttr) { l.log(slog.LevelDebug, msg, attrs...) }
func (l *LokiLogger) Info(msg string, attrs ...attr.LogAttr)  { l.log(slog.LevelInfo, msg, attrs...) }
func (l *LokiLogger) Warn(msg string, attrs ...attr.LogAttr)  { l.log(slog.LevelWarn, msg, attrs...) }
func (l *LokiLogger) Error(msg string, attrs ...attr.LogAttr) { l.log(slog.LevelError, msg, attrs...) }

// log is a helper function for sending logs to Loki.
func (l *LokiLogger) log(level slog.Level, msg string, attrs ...attr.LogAttr) {
	slogAttrs := make([]slog.Attr, len(attrs))
	for i, a := range attrs {
		slogAttrs[i] = a.ToSlogAttr()
	}

	record := slog.Record{
		Time:    time.Now(),
		Level:   level,
		Message: msg,
	}
	record.AddAttrs(slogAttrs...)

	if err := l.handler.Handle(context.Background(), record); err != nil {
		slog.Error("Failed to send log to Loki", slog.Any("error", err), slog.String("message", msg))
	}
}
