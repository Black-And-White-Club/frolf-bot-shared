package observability

import (
	"context"
	"log/slog"
	"time"

	"github.com/grafana/loki-client-go/loki"
	slogloki "github.com/samber/slog-loki/v3"
)

// Logger defines the interface for structured logging.  This is now *in* loki.go
type Logger interface {
	Info(ctx context.Context, msg string, fields map[string]interface{})
	Error(ctx context.Context, msg string, fields map[string]interface{})
	Debug(ctx context.Context, msg string, fields map[string]interface{})
	Warn(ctx context.Context, msg string, fields map[string]interface{})
	Shutdown() // Add to the interface, for consistency
}

// LokiLogger wraps a Loki handler for structured logging and implements the Logger interface.
type LokiLogger struct {
	handler *slogloki.LokiHandler
	client  *loki.Client // Store the Loki client
}

// NewLokiLogger initializes a Loki logger instance.
func NewLokiLogger(url, tenantID, serviceName string) (*LokiLogger, error) {
	// Create a new Loki client
	config, err := loki.NewDefaultConfig(url)
	if err != nil {
		return nil, err
	}
	config.TenantID = tenantID
	// config.ClientTimeout = 30 * time.Second // Example: Set a reasonable timeout. GOOD PRACTICE

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

// Shutdown ensures all logs are flushed before the app exits.  Good practice!
func (l *LokiLogger) Shutdown() {
	if l.client != nil {
		l.client.Stop() // Flush logs before exit
	}
}

// --- Implement the Logger interface methods ---

func (l *LokiLogger) Info(ctx context.Context, message string, fields map[string]interface{}) {
	l.logWithLoki(ctx, slog.LevelInfo, message, fields)
}

func (l *LokiLogger) Error(ctx context.Context, message string, fields map[string]interface{}) {
	l.logWithLoki(ctx, slog.LevelError, message, fields)
}

func (l *LokiLogger) Debug(ctx context.Context, message string, fields map[string]interface{}) {
	l.logWithLoki(ctx, slog.LevelDebug, message, fields)
}

func (l *LokiLogger) Warn(ctx context.Context, message string, fields map[string]interface{}) {
	l.logWithLoki(ctx, slog.LevelWarn, message, fields)
}

// logWithLoki sends structured logs to Loki.  Internal helper function.
func (l *LokiLogger) logWithLoki(ctx context.Context, level slog.Level, message string, fields map[string]interface{}) {
	// Ensure ctx is not nil (good practice, you already had this)
	if ctx == nil {
		ctx = context.Background()
	}

	// Convert fields to structured attributes (you had this correctly)
	attrs := make([]slog.Attr, 0, len(fields))
	for key, value := range fields {
		attrs = append(attrs, slog.Any(key, value))
	}

	// Create log record (you had this correctly)
	record := slog.Record{
		Time:    time.Now(),
		Level:   level,
		Message: message,
	}
	record.AddAttrs(attrs...) // Keep AddAttrs

	// Send log to Loki
	if err := l.handler.Handle(ctx, record); err != nil {
		// CRITICAL: Handle errors from the handler!  Don't just ignore them.
		slog.Error("Failed to send log to Loki", "error", err, "message", message, "fields", fields) //log it to default.
	}
}
