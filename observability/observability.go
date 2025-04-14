package observability

import (
	"context"
)

type Observability struct {
	Provider *Provider
	Registry *Registry
}

// Init sets up the full observability stack: providers and module metrics
func Init(ctx context.Context, cfg Config) (*Observability, error) {
	provider, err := Setup(ctx, cfg)
	if err != nil {
		return nil, err
	}

	registry := NewRegistry(provider, cfg)

	return &Observability{
		Provider: provider,
		Registry: registry,
	}, nil
}
