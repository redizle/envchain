package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInitCmd_FlagsExist(t *testing.T) {
	cmd := newInitCmd()
	if cmd.Flags().Lookup("force") == nil {
		t.Error("expected --force flag")
	}
}

// changeDirTemp changes the working directory to dir and restores it on cleanup.
func changeDirTemp(t *testing.T, dir string) {
	t.Helper()
	old, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("failed to chdir to %s: %v", dir, err)
	}
	t.Cleanup(func() { os.Chdir(old) })
}

func TestInitCmd_CreatesFiles(t *testing.T) {
	dir := t.TempDir()
	changeDirTemp(t, dir)

	cmd := newInitCmd()
	buf := &strings.Builder{}
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, name := range []string{"envchain.yaml", "base.env", "local.env"} {
		path := filepath.Join(dir, name)
		if _, err := os.Stat(path); err != nil {
			t.Errorf("expected file %s to exist: %v", name, err)
		}
	}

	if !strings.Contains(buf.String(), "initialized") {
		t.Error("expected initialization message in output")
	}
}

func TestInitCmd_SkipsExistingFiles(t *testing.T) {
	dir := t.TempDir()
	changeDirTemp(t, dir)

	existing := filepath.Join(dir, "envchain.yaml")
	os.WriteFile(existing, []byte("original"), 0644)

	cmd := newInitCmd()
	buf := &strings.Builder{}
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content, _ := os.ReadFile(existing)
	if string(content) != "original" {
		t.Error("expected existing file to be unchanged")
	}
	if !strings.Contains(buf.String(), "skipping envchain.yaml") {
		t.Error("expected skip message for existing file")
	}
}

func TestInitCmd_ForceOverwrites(t *testing.T) {
	dir := t.TempDir()
	changeDirTemp(t, dir)

	existing := filepath.Join(dir, "envchain.yaml")
	os.WriteFile(existing, []byte("original"), 0644)

	cmd := newInitCmd()
	cmd.SetArgs([]string{"--force"})
	buf := &strings.Builder{}
	cmd.SetOut(buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content, _ := os.ReadFile(existing)
	if string(content) == "original" {
		t.Error("expected file to be overwritten with --force")
	}
}
