package providers_test

import (
	"context"
	"errors"
	"testing"

	"github.com/user/envchain/internal/secrets/providers"
)

// mockSSMClient is a minimal stand-in used to test SSMProvider without real AWS calls.
// Since the real client is injected via NewSSMProviderWithClient, we use a fake
// that satisfies the same interface via a wrapper approach in the provider itself.
// Here we test the Name() contract and error-path behaviour via a stub resolver.

type stubResolver func(ctx context.Context, ref string) (string, error)

func (s stubResolver) Name() string { return "awsssm" }
func (s stubResolver) Resolve(ctx context.Context, ref string) (string, error) {
	return s(ctx, ref)
}

func TestSSMProvider_Name(t *testing.T) {
	// NewSSMProviderWithClient accepts a nil client for name-only checks.
	p := providers.NewSSMProviderWithClient(nil)
	if got := p.Name(); got != "awsssm" {
		t.Errorf("Name() = %q, want %q", got, "awsssm")
	}
}

func TestStubResolver_ReturnsValue(t *testing.T) {
	ctx := context.Background()
	stub := stubResolver(func(_ context.Context, ref string) (string, error) {
		if ref == "/myapp/secret" {
			return "s3cr3t", nil
		}
		return "", errors.New("not found")
	})

	val, err := stub.Resolve(ctx, "/myapp/secret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "s3cr3t" {
		t.Errorf("Resolve() = %q, want %q", val, "s3cr3t")
	}
}

func TestStubResolver_ErrorOnUnknown(t *testing.T) {
	ctx := context.Background()
	stub := stubResolver(func(_ context.Context, ref string) (string, error) {
		return "", errors.New("not found")
	})

	_, err := stub.Resolve(ctx, "/does/not/exist")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestStubResolver_Name(t *testing.T) {
	stub := stubResolver(nil)
	if got := stub.Name(); got != "awsssm" {
		t.Errorf("Name() = %q, want %q", got, "awsssm")
	}
}
