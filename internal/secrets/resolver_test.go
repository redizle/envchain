package secrets_test

import (
	"testing"

	"github.com/user/envchain/internal/secrets"
)

func staticProvider(m map[string]string) *secrets.StaticProvider {
	return &secrets.StaticProvider{Secrets: m}
}

func TestResolve_NoRefs(t *testing.T) {
	r := secrets.NewResolver(nil)
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	out, err := r.Resolve(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["FOO"] != "bar" || out["BAZ"] != "qux" {
		t.Errorf("expected unchanged values, got %v", out)
	}
}

func TestResolve_StaticProvider(t *testing.T) {
	p := staticProvider(map[string]string{"DB_PASS": "s3cr3t"})
	r := secrets.NewResolver(map[string]secrets.Provider{"static": p})
	env := map[string]string{"DATABASE_URL": "postgres://user:${static:DB_PASS}@localhost/db"}
	out, err := r.Resolve(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "postgres://user:s3cr3t@localhost/db"
	if out["DATABASE_URL"] != want {
		t.Errorf("got %q, want %q", out["DATABASE_URL"], want)
	}
}

func TestResolve_MultipleRefs(t *testing.T) {
	p := staticProvider(map[string]string{"USER": "admin", "PASS": "hunter2"})
	r := secrets.NewResolver(map[string]secrets.Provider{"static": p})
	env := map[string]string{"CREDS": "${static:USER}:${static:PASS}"}
	out, err := r.Resolve(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["CREDS"] != "admin:hunter2" {
		t.Errorf("got %q", out["CREDS"])
	}
}

func TestResolve_UnknownProvider(t *testing.T) {
	r := secrets.NewResolver(map[string]secrets.Provider{})
	env := map[string]string{"TOKEN": "${vault:MY_TOKEN}"}
	_, err := r.Resolve(env)
	if err == nil {
		t.Fatal("expected error for unknown provider")
	}
}

func TestResolve_MissingSecret(t *testing.T) {
	p := staticProvider(map[string]string{})
	r := secrets.NewResolver(map[string]secrets.Provider{"static": p})
	env := map[string]string{"API_KEY": "${static:MISSING_KEY}"}
	_, err := r.Resolve(env)
	if err == nil {
		t.Fatal("expected error for missing secret")
	}
}

func TestResolve_OriginalUnmodified(t *testing.T) {
	p := staticProvider(map[string]string{"X": "resolved"})
	r := secrets.NewResolver(map[string]secrets.Provider{"static": p})
	env := map[string]string{"VAL": "${static:X}"}
	_, err := r.Resolve(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["VAL"] != "${static:X}" {
		t.Error("original env map was mutated")
	}
}
