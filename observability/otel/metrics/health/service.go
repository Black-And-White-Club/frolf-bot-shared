package health

import (
	"context"
	"fmt"
	"sync"
)

// Service manages health checkers and performs health checks
type Service struct {
	mu       sync.RWMutex
	checkers map[string]Checker
}

// NewService creates a new health service
func NewService() *Service {
	return &Service{
		checkers: make(map[string]Checker),
	}
}

// RegisterChecker adds a health checker
func (s *Service) RegisterChecker(checker Checker) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.checkers[checker.Name()] = checker
}

// GetChecker retrieves a health checker by name
func (s *Service) GetChecker(name string) (Checker, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	checker, exists := s.checkers[name]
	return checker, exists
}

// CheckAll performs health checks on all registered checkers
func (s *Service) CheckAll(ctx context.Context) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for name, checker := range s.checkers {
		if err := checker.Check(ctx); err != nil {
			return fmt.Errorf("health check failed for %s: %w", name, err)
		}
	}
	return nil
}
