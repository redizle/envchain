package providers_test

import (
	"context"
	"errors"
	"testing"

	"github.com/user/envchain/internal/secrets/providers"
)

// stubVaultClient implements VaultClient for testing.
type stubVaultClient struct {
	values map[string]string
}

func (s *stubVaultClient) GetSecret(_ context.Context, path string) (string, error) {
	v, ok := s.values[path]
	if !ok {
		return "", errors.New("vault: secret not found: " + path)
	}
	return v, nil
}

func TestVaultProvider_Name(t *testing.T) {
	p := providers.NewVaultProviderWithClient(&stubVaultClient{})
	if p.Name() != "vault" {
		t.Errorf("expected name 'vault', got %q", p.Name())
	}
}

func TestVaultProvider_Resolve_Found(t *testing.T) {
	stub := &stubVaultClient{
		values: map[string]string{
			"secret/myapp/db_password": "supersecret",
		},
	}
	p := providers.NewVaultProviderWithClient(stub)

	val, err := p.Resolve(context.Background(), "secret/myapp/db_password")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "supersecret" {
		t.Errorf("expected 'supersecret', got %q", val)
	}
}

func TestVaultProvider_Resolve_NotFound(t *testing.T) {
	stub := &stubVaultClient{values: map[string]string{}}
	p := providers.NewVaultProviderWithClient(stub)

	_, err := p.Resolve(context.Background(), "secret/missing")
	if err == nil {
		t.Fatal("expected error for missing secret, got nil")
	}
}

func TestVaultProvider_Resolve_MultipleSecrets(t *testing.T) {
	stub := &stubVaultClient{
		values: map[string]string{
			"secret/a": "alpha",
			"secret/b": "beta",
		},
	}
	p := providers.NewVaultProviderWithClient(stub)

	for path, expected := range stub.values {
		val, err := p.Resolve(context.Background(), path)
		if err != nil {
			t.Errorf("path %q: unexpected error: %v", path, err)
		}
		if val != expected {
			t.Errorf("path %q: expected %q, got %q", path, expected, val)
		}
	}
}
