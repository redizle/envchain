package chain_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/envchain/envchain/internal/chain"
)

func writeTempEnv(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("writeTempEnv: %v", err)
	}
	return p
}

func TestLoad_BasicMerge(t *testing.T) {
	dir := t.TempDir()
	base := writeTempEnv(t, dir, ".env", "APP=base\nDEBUG=false\n")
	local := writeTempEnv(t, dir, ".env.local", "DEBUG=true\nEXTRA=yes\n")

	c, err := chain.Load([]chain.Layer{
		{Name: "base", Path: base},
		{Name: "local", Path: local},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if c.Env["APP"] != "base" {
		t.Errorf("APP: want %q, got %q", "base", c.Env["APP"])
	}
	if c.Env["DEBUG"] != "true" {
		t.Errorf("DEBUG: want %q, got %q", "true", c.Env["DEBUG"])
	}
	if c.Env["EXTRA"] != "yes" {
		t.Errorf("EXTRA: want %q, got %q", "yes", c.Env["EXTRA"])
	}
}

func TestLoad_MissingLayerSkipped(t *testing.T) {
	dir := t.TempDir()
	base := writeTempEnv(t, dir, ".env", "KEY=value\n")

	c, err := chain.Load([]chain.Layer{
		{Name: "base", Path: base},
		{Name: "missing", Path: filepath.Join(dir, ".env.missing")},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Env["KEY"] != "value" {
		t.Errorf("KEY: want %q, got %q", "value", c.Env["KEY"])
	}
}

func TestLoad_EmptyLayers(t *testing.T) {
	c, err := chain.Load([]chain.Layer{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(c.Env) != 0 {
		t.Errorf("expected empty env, got %v", c.Env)
	}
}
