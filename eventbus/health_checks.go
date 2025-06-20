package eventbus

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	nc "github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// GetHealthCheckers returns health checkers for the EventBus components
func (eb *eventBus) GetHealthCheckers() []HealthChecker {
	return []HealthChecker{
		&natsHealthChecker{
			natsConn: eb.natsConn,
			js:       eb.js,
			logger:   eb.logger,
		},
	}
}

// natsHealthChecker implements health checking for NATS/JetStream
type natsHealthChecker struct {
	natsConn *nc.Conn
	js       jetstream.JetStream
	logger   *slog.Logger
}

func (n *natsHealthChecker) Name() string {
	return "NATS/JetStream"
}

func (n *natsHealthChecker) Check(ctx context.Context) error {
	// Check NATS connection
	if n.natsConn == nil {
		return errors.New("NATS connection is nil")
	}

	if !n.natsConn.IsConnected() {
		return errors.New("NATS connection is not connected")
	}

	// Check JetStream availability
	if n.js == nil {
		return errors.New("JetStream is nil")
	}

	// Try a simple JetStream operation with timeout
	checkCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := n.js.AccountInfo(checkCtx)
	if err != nil {
		return fmt.Errorf("JetStream health check failed: %w", err)
	}

	return nil
}

// isRetryableError determines if an error should trigger a retry
func (eb *eventBus) isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	errStr := strings.ToLower(err.Error())

	// Network and timeout related errors that should be retried
	retryableErrors := []string{
		"timeout",
		"connection",
		"network",
		"temporary",
		"i/o timeout",
		"context deadline exceeded",
		"no servers available",
		"connection refused",
		"connection reset",
	}

	for _, retryable := range retryableErrors {
		if strings.Contains(errStr, retryable) {
			return true
		}
	}

	// NATS specific errors that should be retried
	if errors.Is(err, nc.ErrConnectionClosed) ||
		errors.Is(err, nc.ErrNoServers) ||
		errors.Is(err, nc.ErrTimeout) {
		return true
	}

	return false
}
