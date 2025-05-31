package health

import (
	"context"
)

// Status represents the health status of a component
type Status string

const (
	StatusHealthy   Status = "healthy"
	StatusUnhealthy Status = "unhealthy"
	StatusDegraded  Status = "degraded"
)

// Checker defines the interface for health checking components
type Checker interface {
	// Check performs a health check and returns an error if unhealthy
	Check(ctx context.Context) error

	// Name returns the name of the component being checked
	Name() string
}

// HealthCheck is a function type for health check callbacks
type HealthCheck func(ctx context.Context) error

// SimpleChecker is a basic implementation of the Checker interface
type SimpleChecker struct {
	name      string
	checkFunc HealthCheck
}

// NewSimpleChecker creates a new simple health checker
func NewSimpleChecker(name string, checkFunc HealthCheck) *SimpleChecker {
	return &SimpleChecker{
		name:      name,
		checkFunc: checkFunc,
	}
}

// Name returns the name of this checker
func (c *SimpleChecker) Name() string {
	return c.name
}

// Check performs the health check
func (c *SimpleChecker) Check(ctx context.Context) error {
	return c.checkFunc(ctx)
}
