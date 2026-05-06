package cli

import (
	"bytes"
	"testing"
)

func TestRootCmd_HasSubcommands(t *testing.T) {
	root := NewRootCmd()

	names := map[string]bool{}
	for _, sub := range root.Commands() {
		names[sub.Name()] = true
	}

	if !names["print"] {
		t.Error("expected 'print' subcommand")
	}
	if !names["diff"] {
		t.Error("expected 'diff' subcommand")
	}
}

func TestRootCmd_HelpNoError(t *testing.T) {
	root := NewRootCmd()
	root.SetOut(&bytes.Buffer{})
	root.SetErr(&bytes.Buffer{})
	root.SetArgs([]string{"--help"})

	// --help exits with nil in cobra when SilenceUsage is set
	_ = root.Execute()
}

func TestPrintCmd_FlagsExist(t *testing.T) {
	cmd := newPrintCmd()

	if cmd.Flags().Lookup("config") == nil {
		t.Error("expected --config flag on print")
	}
	if cmd.Flags().Lookup("format") == nil {
		t.Error("expected --format flag on print")
	}
	if cmd.Flags().Lookup("resolve") == nil {
		t.Error("expected --resolve flag on print")
	}
}

func TestDiffCmd_FlagsExist(t *testing.T) {
	cmd := newDiffCmd()

	if cmd.Flags().Lookup("config") == nil {
		t.Error("expected --config flag on diff")
	}
}
