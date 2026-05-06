package providers

import (
	"context"
	"fmt"
	"strings"
)

// Provider resolves a secret reference to a plaintext value.
type Provider interface {
	Name() string
	Resolve(ctx context.Context, ref string) (string, error)
}

// Registry holds named secret providers and dispatches resolution.
type Registry struct {
	providers map[string]Provider
}

// NewRegistry returns an empty Registry.
func NewRegistry() *Registry {
	return &Registry{providers: make(map[string]Provider)}
}

// Register adds a provider under its Name().
func (r *Registry) Register(p Provider) {
	r.providers[p.Name()] = p
}

// Get returns the provider for a given name, or an error if not found.
func (r *Registry) Get(name string) (Provider, error) {
	p, ok := r.providers[name]
	if !ok {
		return nil, fmt.Errorf("unknown secret provider %q", name)
	}
	return p, nil
}

// Resolve parses a reference of the form "provider:ref" and delegates
// to the matching registered provider.
func (r *Registry) Resolve(ctx context.Context, raw string) (string, error) {
	idx := strings.Index(raw, ":")
	if idx < 0 {
		return "", fmt.Errorf("secret ref %q missing provider prefix (expected provider:ref)", raw)
	}
	providerName := raw[:idx]
	ref := raw[idx+1:]

	p, err := r.Get(providerName)
	if err != nil {
		return "", err
	}
	return p.Resolve(ctx, ref)
}

// Names returns all registered provider names.
func (r *Registry) Names() []string {
	names := make([]string, 0, len(r.providers))
	for n := range r.providers {
		names = append(names, n)
	}
	return names
}
