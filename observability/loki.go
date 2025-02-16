package observability

import (
	"context"
	"log/slog"
	"time"

	"github.com/grafana/loki-client-go/loki"
	slogloki "github.com/samber/slog-loki/v3"
)

// Logger defines the interface for structured logging.
type Logger interface {
	Debug(ctx context.Context, msg string, attrs ...slog.Attr)
	Info(ctx context.Context, msg string, attrs ...slog.Attr)
	Warn(ctx context.Context, msg string, attrs ...slog.Attr)
	Error(ctx context.Context, msg string, attrs ...slog.Attr)
	Shutdown()
}

// LokiLogger wraps a Loki handler and implements the Logger interface.
type LokiLogger struct {
	handler *slogloki.LokiHandler
	client  *loki.Client
}

// NewLokiLogger initializes a Loki logger instance. (Unchanged)
func NewLokiLogger(url, tenantID, serviceName string) (*LokiLogger, error) {
	config, err := loki.NewDefaultConfig(url)
	if err != nil {
		return nil, err
	}
	config.TenantID = tenantID
	// config.ClientTimeout = 30 * time.Second  // Good practice

	client, err := loki.New(config)
	if err != nil {
		return nil, err
	}

	// Define Loki handler with structured attributes
	handler := slogloki.Option{
		Level:  slog.LevelDebug,
		Client: client,
	}.NewLokiHandler().(*slogloki.LokiHandler)

	// Attach a "service" label
	handler = handler.WithAttrs([]slog.Attr{
		slog.String("service", serviceName),
	}).(*slogloki.LokiHandler)

	return &LokiLogger{
		handler: handler,
		client:  client, // Store Loki client for shutdown
	}, nil
}

// Shutdown ensures all logs are flushed.
func (l *LokiLogger) Shutdown() {
	if l.client != nil {
		l.client.Stop()
	}
}

// --- Implement the Logger interface methods ---

func (l *LokiLogger) Debug(ctx context.Context, msg string, attrs ...slog.Attr) {
	l.log(ctx, slog.LevelDebug, msg, attrs...)
}

func (l *LokiLogger) Info(ctx context.Context, msg string, attrs ...slog.Attr) {
	l.log(ctx, slog.LevelInfo, msg, attrs...)
}

func (l *LokiLogger) Warn(ctx context.Context, msg string, attrs ...slog.Attr) {
	l.log(ctx, slog.LevelWarn, msg, attrs...)
}

func (l *LokiLogger) Error(ctx context.Context, msg string, attrs ...slog.Attr) {
	l.log(ctx, slog.LevelError, msg, attrs...)
}

// log is a private helper function to avoid repetition.
func (l *LokiLogger) log(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr) {
	if ctx == nil {
		ctx = context.Background()
	}

	record := slog.Record{
		Time:    time.Now(),
		Level:   level,
		Message: msg,
	}
	record.AddAttrs(attrs...)

	if err := l.handler.Handle(ctx, record); err != nil {
		slog.Error("Failed to send log to Loki", slog.Any("error", err), slog.String("message", msg))
	}
}
