package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeValidateFixture(t *testing.T, dir, cfg, envContent string) string {
	t.Helper()
	envPath := filepath.Join(dir, "base.env")
	if err := os.WriteFile(envPath, []byte(envContent), 0600); err != nil {
		t.Fatal(err)
	}
	cfgPath := filepath.Join(dir, ".envchain.yaml")
	if err := os.WriteFile(cfgPath, []byte(cfg), 0600); err != nil {
		t.Fatal(err)
	}
	return cfgPath
}

func TestValidateCmd_FlagsExist(t *testing.T) {
	cmd := newValidateCmd()
	if cmd.Flags().Lookup("config") == nil {
		t.Error("expected --config flag")
	}
}

func TestValidateCmd_ValidEnv(t *testing.T) {
	dir := t.TempDir()
	cfgYAML := "layers:\n  - " + filepath.Join(dir, "base.env") + "\n"
	cfgPath := writeValidateFixture(t, dir, cfgYAML, "APP_NAME=envchain\nPORT=8080\n")

	root := NewRootCmd()
	root.SetArgs([]string{"validate", "--config", cfgPath})
	var out strings.Builder
	root.SetOut(&out)
	if err := root.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "OK") {
		t.Errorf("expected OK in output, got: %s", out.String())
	}
}

func TestValidateCmd_MissingConfig(t *testing.T) {
	root := NewRootCmd()
	root.SetArgs([]string{"validate", "--config", "/nonexistent/.envchain.yaml"})
	if err := root.Execute(); err == nil {
		t.Fatal("expected error for missing config")
	}
}
