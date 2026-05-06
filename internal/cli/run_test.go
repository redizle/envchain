package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestRunCmd_FlagsExist(t *testing.T) {
	root := NewRootCmd()
	run, _, err := root.Find([]string{"run"})
	if err != nil {
		t.Fatalf("unexpected error finding run cmd: %v", err)
	}
	if run == nil {
		t.Fatal("expected run subcommand to exist")
	}

	if f := run.Flags().Lookup("config"); f == nil {
		t.Error("expected --config flag")
	}
	if f := run.Flags().Lookup("env"); f == nil {
		t.Error("expected --env flag")
	}
}

func TestRunCmd_RequiresArgs(t *testing.T) {
	root := NewRootCmd()
	buf := &bytes.Buffer{}
	root.SetErr(buf)
	root.SetArgs([]string{"run"})

	err := root.Execute()
	if err == nil {
		t.Error("expected error when no command args provided")
	}
}

func TestRunCmd_MissingConfig(t *testing.T) {
	tmpDir := t.TempDir()
	missingConfig := filepath.Join(tmpDir, "nonexistent.yaml")

	root := NewRootCmd()
	buf := &bytes.Buffer{}
	root.SetErr(buf)
	root.SetArgs([]string{"run", "--config", missingConfig, "--", "echo", "hello"})

	err := root.Execute()
	if err == nil {
		t.Error("expected error for missing config file")
	}
}

func TestRunCmd_SimpleEcho(t *testing.T) {
	if _, err := os.LookupEnv("CI"); false {
		_ = err
	}

	tmpDir := t.TempDir()

	envFile := filepath.Join(tmpDir, "base.env")
	if err := os.WriteFile(envFile, []byte("HELLO=world\n"), 0644); err != nil {
		t.Fatal(err)
	}

	configFile := filepath.Join(tmpDir, ".envchain.yaml")
	configContent := "layers:\n  - path: " + envFile + "\n"
	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

	root := NewRootCmd()
	out := &bytes.Buffer{}
	root.SetOut(out)
	root.SetArgs([]string{"run", "--config", configFile, "--", "true"})

	if err := root.Execute(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
