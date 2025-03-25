package eventbus

import (
	"context"
	"fmt"
	"time"

	"github.com/Black-And-White-Club/frolf-bot-shared/observability/attr"
	lokifrolfbot "github.com/Black-And-White-Club/frolf-bot-shared/observability/loki"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// NATSHealthChecker is a health checker for NATS connections
type NATSHealthChecker struct {
	conn    *nats.Conn
	timeout time.Duration
	logger  lokifrolfbot.Logger
}

// NewNATSHealthChecker creates a new health checker for NATS
func NewNATSHealthChecker(conn *nats.Conn, logger lokifrolfbot.Logger) *NATSHealthChecker {
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
	logAttrs := []attr.LogAttr{
		attr.String("component", "nats"),
		attr.String("operation", "health_check"),
	}

	if n.conn == nil {
		err := fmt.Errorf("NATS connection is nil")
		if n.logger != nil {
			n.logger.Warn("NATS health check failed", append(logAttrs, attr.Error(err))...)
		}
		return err
	}

	if !n.conn.IsConnected() {
		err := fmt.Errorf("NATS connection is not established")
		if n.logger != nil {
			n.logger.Warn("NATS health check failed", append(logAttrs, attr.Error(err))...)
		}
		return err
	}

	// Test connectivity with ping/flush
	if err := n.conn.FlushTimeout(n.timeout); err != nil {
		if n.logger != nil {
			n.logger.Warn("NATS ping failed", append(logAttrs,
				attr.Error(err),
				attr.Duration("timeout", n.timeout))...)
		}
		return fmt.Errorf("NATS server ping failed: %w", err)
	}

	if n.logger != nil {
		n.logger.Debug("NATS health check passed", logAttrs...)
	}
	return nil
}

// JetStreamHealthChecker is a health checker for JetStream
type JetStreamHealthChecker struct {
	js     jetstream.JetStream
	logger lokifrolfbot.Logger
}

// NewJetStreamHealthChecker creates a new health checker for JetStream
func NewJetStreamHealthChecker(js jetstream.JetStream, logger lokifrolfbot.Logger) *JetStreamHealthChecker {
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
	logAttrs := []attr.LogAttr{
		attr.String("component", "jetstream"),
		attr.String("operation", "health_check"),
	}

	if j.js == nil {
		err := fmt.Errorf("JetStream context is nil")
		if j.logger != nil {
			j.logger.Warn("JetStream health check failed", append(logAttrs, attr.Error(err))...)
		}
		return err
	}

	// Attempt to fetch account info to verify JetStream is operational
	accountInfo, err := j.js.AccountInfo(ctx)
	if err != nil {
		if j.logger != nil {
			j.logger.Warn("JetStream health check failed", append(logAttrs, attr.Error(err))...)
		}
		return fmt.Errorf("JetStream account info check failed: %w", err)
	}

	if accountInfo == nil {
		err := fmt.Errorf("JetStream account info is nil")
		if j.logger != nil {
			j.logger.Warn("JetStream health check failed", append(logAttrs, attr.Error(err))...)
		}
		return err
	}

	if j.logger != nil {
		j.logger.Debug("JetStream health check passed", logAttrs...)
	}
	return nil
}

// GetHealthCheckers returns health checkers for NATS and JetStream connections
func (eb *eventBus) GetHealthCheckers() []HealthChecker {
	var checkers []HealthChecker

	// Create NATS health checker
	checkers = append(checkers, &NATSHealthChecker{
		conn:    eb.natsConn,
		timeout: 5 * time.Second,
		logger:  eb.logger,
	})

	// Create JetStream health checker
	checkers = append(checkers, &JetStreamHealthChecker{
		js:     eb.js,
		logger: eb.logger,
	})

	return checkers
}
