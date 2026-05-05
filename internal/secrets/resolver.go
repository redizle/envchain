package secrets

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Provider is a source of secret values keyed by name.
type Provider interface {
	Get(key string) (string, error)
}

// Resolver interpolates secret references in env values.
// References use the syntax: ${secret:KEY}
type Resolver struct {
	providers map[string]Provider
}

var refPattern = regexp.MustCompile(`\$\{([a-zA-Z0-9_]+):([^}]+)\}`)

// NewResolver creates a Resolver with the given named providers.
func NewResolver(providers map[string]Provider) *Resolver {
	return &Resolver{providers: providers}
}

// Resolve replaces all secret references in the given env map.
// Returns a new map with resolved values; original is not modified.
func (r *Resolver) Resolve(env map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		resolved, err := r.resolveValue(v)
		if err != nil {
			return nil, fmt.Errorf("resolving %q: %w", k, err)
		}
		out[k] = resolved
	}
	return out, nil
}

func (r *Resolver) resolveValue(value string) (string, error) {
	var resolveErr error
	result := refPattern.ReplaceAllStringFunc(value, func(match string) string {
		if resolveErr != nil {
			return match
		}
		parts := refPattern.FindStringSubmatch(match)
		if len(parts) != 3 {
			return match
		}
		providerName, key := parts[1], parts[2]
		p, ok := r.providers[providerName]
		if !ok {
			resolveErr = fmt.Errorf("unknown provider %q", providerName)
			return match
		}
		val, err := p.Get(key)
		if err != nil {
			resolveErr = fmt.Errorf("provider %q: %w", providerName, err)
			return match
		}
		return val
	})
	return result, resolveErr
}

// EnvProvider reads secrets from the process environment.
type EnvProvider struct{}

func (e *EnvProvider) Get(key string) (string, error) {
	v, ok := os.LookupEnv(key)
	if !ok {
		return "", fmt.Errorf("env var %q not set", key)
	}
	return v, nil
}

// StaticProvider holds a fixed map of secrets (useful for testing).
type StaticProvider struct {
	Secrets map[string]string
}

func (s *StaticProvider) Get(key string) (string, error) {
	v, ok := s.Secrets[key]
	if !ok {
		return "", fmt.Errorf("secret %q not found", key)
	}
	return strings.TrimSpace(v), nil
}
