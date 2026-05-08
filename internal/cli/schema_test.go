package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func writeSchemaFixture(t *testing.T, dir string) (cfgPath, schemaPath string) {
	t.Helper()

	envPath := filepath.Join(dir, "base.env")
	_ = os.WriteFile(envPath, []byte("PORT=8080\nAPP_ENV=production\n"), 0644)

	cfgPath = filepath.Join(dir, ".envchain.yaml")
	cfgContent := "layers:\n  - " + envPath + "\n"
	_ = os.WriteFile(cfgPath, []byte(cfgContent), 0644)

	schemaPath = filepath.Join(dir, ".envschema.yaml")
	schemaContent := "rules:\n  - key: PORT\n    required: true\n    pattern: '^\\d+$'\n  - key: APP_ENV\n    required: true\n"
	_ = os.WriteFile(schemaPath, []byte(schemaContent), 0644)

	return cfgPath, schemaPath
}

func TestSchemaCmd_FlagsExist(t *testing.T) {
	cmd := newSchemaCmd()
	if cmd.Flags().Lookup("config") == nil {
		t.Error("expected --config flag")
	}
	if cmd.Flags().Lookup("schema") == nil {
		t.Error("expected --schema flag")
	}
}

func TestSchemaCmd_PassesValidEnv(t *testing.T) {
	dir := t.TempDir()
	cfgPath, schemaPath := writeSchemaFixture(t, dir)

	cmd := newSchemaCmd()
	cmd.SetArgs([]string{"--config", cfgPath, "--schema", schemaPath})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestSchemaCmd_MissingConfig(t *testing.T) {
	cmd := newSchemaCmd()
	cmd.SetArgs([]string{"--config", "/nonexistent/.envchain.yaml", "--schema", "/nonexistent/.envschema.yaml"})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing config")
	}
}

func TestSchemaCmd_MissingSchemaFile(t *testing.T) {
	dir := t.TempDir()
	cfgPath, _ := writeSchemaFixture(t, dir)

	cmd := newSchemaCmd()
	cmd.SetArgs([]string{"--config", cfgPath, "--schema", filepath.Join(dir, "noschema.yaml")})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing schema file")
	}
}
