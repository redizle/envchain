package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeAuditFixture(t *testing.T, dir string, envContent string) string {
	t.Helper()
	envFile := filepath.Join(dir, "base.env")
	if err := os.WriteFile(envFile, []byte(envContent), 0644); err != nil {
		t.Fatal(err)
	}
	cfg := "layers:\n  - path: base.env\n"
	cfgFile := filepath.Join(dir, ".envchain.yaml")
	if err := os.WriteFile(cfgFile, []byte(cfg), 0644); err != nil {
		t.Fatal(err)
	}
	return cfgFile
}

func TestAuditCmd_FlagsExist(t *testing.T) {
	cmd := newAuditCmd()
	if cmd.Flags().Lookup("config") == nil {
		t.Error("expected --config flag")
	}
	if cmd.Flags().Lookup("env") == nil {
		t.Error("expected --env flag")
	}
}

func TestAuditCmd_CleanEnv(t *testing.T) {
	dir := t.TempDir()
	cfgFile := writeAuditFixture(t, dir, "APP_NAME=myapp\nPORT=8080\n")

	cmd := NewRootCmd()
	cmd.SetArgs([]string{"audit", "--config", cfgFile})
	var out strings.Builder
	cmd.SetOut(&out)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "No issues found") {
		t.Errorf("expected clean message, got: %s", out.String())
	}
}

func TestAuditCmd_DetectsPlainSecret(t *testing.T) {
	dir := t.TempDir()
	cfgFile := writeAuditFixture(t, dir, "DB_PASSWORD=hunter2\n")

	cmd := NewRootCmd()
	cmd.SetArgs([]string{"audit", "--config", cfgFile})
	var out strings.Builder
	cmd.SetOut(&out)

	// We expect a warning (non-zero but not necessarily error exit here in test)
	_ = cmd.Execute()
	if !strings.Contains(out.String(), "[warn]") {
		t.Errorf("expected warning output, got: %s", out.String())
	}
}

func TestAuditCmd_MissingConfig(t *testing.T) {
	cmd := NewRootCmd()
	cmd.SetArgs([]string{"audit", "--config", "/nonexistent/.envchain.yaml"})
	if err := cmd.Execute(); err == nil {
		t.Error("expected error for missing config")
	}
}
