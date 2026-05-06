package providers

import (
	"context"
	"fmt"
)

// Provider is the interface all secret providers must implement.
type Provider interface {
	Name() string
	Resolve(ctx context.Context, ref string) (string, error)
}

// Registry holds named providers and dispatches resolution by name.
type Registry struct {
	providers map[string]Provider
}

// NewRegistry creates an empty provider registry.
func NewRegistry() *Registry {
	return &Registry{providers: make(map[string]Provider)}
}

// Register adds a provider to the registry. Panics on duplicate names.
func (r *Registry) Register(p Provider) {
	name := p.Name()
	if _, exists := r.providers[name]; exists {
		panic(fmt.Sprintf("providers: duplicate provider name %q", name))
	}
	r.providers[name] = p
}

// Get returns the provider for the given name, or an error if not found.
func (r *Registry) Get(name string) (Provider, error) {
	p, ok := r.providers[name]
	if !ok {
		return nil, fmt.Errorf("providers: unknown provider %q", name)
	}
	return p, nil
}

// Names returns all registered provider names.
func (r *Registry) Names() []string {
	names := make([]string, 0, len(r.providers))
	for n := range r.providers {
		names = append(names, n)
	}
	return names
}

// Resolve dispatches a resolution call to the named provider.
func (r *Registry) Resolve(ctx context.Context, providerName, ref string) (string, error) {
	p, err := r.Get(providerName)
	if err != nil {
		return "", err
	}
	return p.Resolve(ctx, ref)
}
