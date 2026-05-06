package providers

import (
	"context"
	"fmt"
	"testing"
)

type stubDopplerClient struct {
	secrets map[string]string
}

func (s *stubDopplerClient) GetSecret(_ context.Context, _, project, config, name string) (string, error) {
	key := fmt.Sprintf("%s/%s/%s", project, config, name)
	val, ok := s.secrets[key]
	if !ok {
		return "", fmt.Errorf("secret %q not found in doppler", name)
	}
	return val, nil
}

func TestDopplerProvider_Name(t *testing.T) {
	p := NewDopplerProvider("tok")
	if p.Name() != "doppler" {
		t.Errorf("expected 'doppler', got %q", p.Name())
	}
}

func TestDopplerProvider_Resolve_Found(t *testing.T) {
	client := &stubDopplerClient{secrets: map[string]string{
		"myapp/production/DB_PASSWORD": "supersecret",
	}}
	p := NewDopplerProviderWithClient("tok", client)

	val, err := p.Resolve(context.Background(), "myapp/production/DB_PASSWORD")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "supersecret" {
		t.Errorf("expected 'supersecret', got %q", val)
	}
}

func TestDopplerProvider_Resolve_NotFound(t *testing.T) {
	client := &stubDopplerClient{secrets: map[string]string{}}
	p := NewDopplerProviderWithClient("tok", client)

	_, err := p.Resolve(context.Background(), "myapp/production/MISSING")
	if err == nil {
		t.Fatal("expected error for missing secret")
	}
}

func TestDopplerProvider_Resolve_BadRef(t *testing.T) {
	p := NewDopplerProvider("tok")

	_, err := p.Resolve(context.Background(), "bad-ref")
	if err == nil {
		t.Fatal("expected error for malformed ref")
	}
}

func TestDopplerProvider_Resolve_MultipleSecrets(t *testing.T) {
	client := &stubDopplerClient{secrets: map[string]string{
		"proj/staging/API_KEY":  "key-abc",
		"proj/staging/API_SECRET": "sec-xyz",
	}}
	p := NewDopplerProviderWithClient("tok", client)

	for ref, want := range map[string]string{
		"proj/staging/API_KEY":    "key-abc",
		"proj/staging/API_SECRET": "sec-xyz",
	} {
		got, err := p.Resolve(context.Background(), ref)
		if err != nil {
			t.Fatalf("ref %q: unexpected error: %v", ref, err)
		}
		if got != want {
			t.Errorf("ref %q: expected %q, got %q", ref, want, got)
		}
	}
}
