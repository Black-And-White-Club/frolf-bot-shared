package lokifrolfbot

import (
	"context"

	"github.com/Black-And-White-Club/frolf-bot-shared/observability/attr"
)

// NoOpLogger is a logger that does nothing. Useful for unit tests.
type NoOpLogger struct{}

func (n *NoOpLogger) Debug(msg string, attrs ...attr.LogAttr) {}
func (n *NoOpLogger) Info(msg string, attrs ...attr.LogAttr)  {}
func (n *NoOpLogger) Warn(msg string, attrs ...attr.LogAttr)  {}
func (n *NoOpLogger) Error(msg string, attrs ...attr.LogAttr) {}
func (n *NoOpLogger) Close()                                  {}
func (n *NoOpLogger) WithContext(ctx context.Context) Logger  { return n }
