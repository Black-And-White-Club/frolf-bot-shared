package eventbus

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Black-And-White-Club/frolf-bot-shared/observability/attr"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// NATSHealthChecker is a health checker for NATS connections
type NATSHealthChecker struct {
	conn    *nats.Conn
	timeout time.Duration
	logger  *slog.Logger // Updated to use slog.Logger
}

// NewNATSHealthChecker creates a new health checker for NATS
func NewNATSHealthChecker(conn *nats.Conn, logger *slog.Logger) *NATSHealthChecker {
	return &NATSHealthChecker{
		conn:    conn,
		timeout: 5 * time.Second, // Default timeout
		logger:  logger,
	}
}

// Name returns the name of the health checker
func (n *NATSHealthChecker) Name() string {
	return "NATS"
}

// Check performs the health check
func (n *NATSHealthChecker) Check(ctx context.Context) error {
	if n.conn == nil {
		err := fmt.Errorf("NATS connection is nil")
		n.logger.Warn("NATS health check failed", attr.Error(err)) // Directly use attr.LogAttr
		return err
	}

	if !n.conn.IsConnected() {
		err := fmt.Errorf("NATS connection is not established")
		n.logger.Warn("NATS health check failed", attr.Error(err)) // Directly use attr.LogAttr
		return err
	}

	// Test connectivity with ping/flush
	if err := n.conn.FlushTimeout(n.timeout); err != nil {
		n.logger.Warn("NATS ping failed", attr.Error(err), attr.Duration("timeout", n.timeout)) // Directly use attr.LogAttr
		return fmt.Errorf("NATS server ping failed: %w", err)
	}

	n.logger.Debug("NATS health check passed") // No attributes needed
	return nil
}

// JetStreamHealthChecker is a health checker for JetStream
type JetStreamHealthChecker struct {
	js     jetstream.JetStream
	logger *slog.Logger // Updated to use slog.Logger
}

// NewJetStreamHealthChecker creates a new health checker for JetStream
func NewJetStreamHealthChecker(js jetstream.JetStream, logger *slog.Logger) *JetStreamHealthChecker {
	return &JetStreamHealthChecker{
		js:     js,
		logger: logger,
	}
}

// Name returns the name of the health checker
func (j *JetStreamHealthChecker) Name() string {
	return "JetStream"
}

// Check performs the health check
func (j *JetStreamHealthChecker) Check(ctx context.Context) error {
	if j.js == nil {
		err := fmt.Errorf("JetStream context is nil")
		j.logger.Warn("JetStream health check failed", attr.Error(err)) // Directly use attr.LogAttr
		return err
	}

	// Attempt to fetch account info to verify JetStream is operational
	accountInfo, err := j.js.AccountInfo(ctx)
	if err != nil {
		j.logger.Warn("JetStream health check failed", attr.Error(err)) // Directly use attr.LogAttr
		return fmt.Errorf("JetStream account info check failed: %w", err)
	}

	if accountInfo == nil {
		err := fmt.Errorf("JetStream account info is nil")
		j.logger.Warn("JetStream health check failed", attr.Error(err)) // Directly use attr.LogAttr
		return err
	}

	j.logger.Debug("JetStream health check passed") // No attributes needed
	return nil
}

// GetHealthCheckers returns health checkers for NATS and JetStream connections
func (eb *eventBus) GetHealthCheckers() []HealthChecker {
	var checkers []HealthChecker

	// Create NATS health checker
	checkers = append(checkers, &NATSHealthChecker{
		conn:    eb.natsConn,
		timeout: 5 * time.Second,
		logger:  eb.logger, // Ensure this is updated to use *slog.Logger
	})

	// Create JetStream health checker
	checkers = append(checkers, &JetStreamHealthChecker{
		js:     eb.js,
		logger: eb.logger, // Ensure this is updated to use *slog.Logger
	})

	return checkers
}
