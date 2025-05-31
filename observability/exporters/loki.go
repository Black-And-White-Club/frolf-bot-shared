// exporters/logger.go
package exporters

import (
	"log/slog"
	"os"
)

// NewJSONLogger returns a JSON slog.Logger to stdout.
func NewJSONLogger(level slog.Level) *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))
}
