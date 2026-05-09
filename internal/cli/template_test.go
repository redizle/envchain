package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTemplateFixture(t *testing.T, dir string) {
	t.Helper()

	cfg := `layers:
  - name: base
    path: base.env
`
	base := `BASE_URL=https://example.com
API_URL={{ .BASE_URL }}/api
PLAIN=static
`
	if err := os.WriteFile(filepath.Join(dir, ".envchain.yaml"), []byte(cfg), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "base.env"), []byte(base), 0644); err != nil {
		t.Fatal(err)
	}
}

func TestTemplateCmd_FlagsExist(t *testing.T) {
	cmd := newTemplateCmd()
	flags := []string{"config", "env", "missing-key", "format"}
	for _, f := range flags {
		if cmd.Flags().Lookup(f) == nil {
			t.Errorf("expected flag --%s", f)
		}
	}
}

func TestTemplateCmd_MissingConfig(t *testing.T) {
	dir := t.TempDir()
	old, _ := os.Getwd()
	defer os.Chdir(old)
	os.Chdir(dir)

	cmd := newTemplateCmd()
	cmd.SetArgs([]string{"--config", "nonexistent.yaml"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for missing config")
	}
}

func TestTemplateCmd_RendersTemplates(t *testing.T) {
	dir := t.TempDir()
	old, _ := os.Getwd()
	defer os.Chdir(old)
	os.Chdir(dir)

	writeTemplateFixture(t, dir)

	root := NewRootCmd()
	root.SetArgs([]string{"template", "--config", ".envchain.yaml", "--format", "dotenv", "--missing-key", "zero"})
	if err := root.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTemplateCmd_ErrorOnMissingKeyDefault(t *testing.T) {
	dir := t.TempDir()
	old, _ := os.Getwd()
	defer os.Chdir(old)
	os.Chdir(dir)

	cfg := `layers:
  - name: base
    path: base.env
`
	base := `X={{ .DOES_NOT_EXIST }}
`
	os.WriteFile(filepath.Join(dir, ".envchain.yaml"), []byte(cfg), 0644)
	os.WriteFile(filepath.Join(dir, "base.env"), []byte(base), 0644)

	cmd := newTemplateCmd()
	cmd.SetArgs([]string{"--config", ".envchain.yaml", "--missing-key", "error"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for missing key")
	}
}
