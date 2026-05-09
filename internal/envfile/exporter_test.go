package envfile

import (
	"strings"
	"testing"
)

func TestExport_DotenvFormat(t *testing.T) {
	env := map[string]string{
		"FOO": "bar",
		"BAZ": "qux",
	}
	out := Export(env, ExportOptions{Format: FormatDotenv, Sorted: true})
	if !strings.Contains(out, "BAZ=qux\n") {
		t.Errorf("expected BAZ=qux, got: %s", out)
	}
	if !strings.Contains(out, "FOO=bar\n") {
		t.Errorf("expected FOO=bar, got: %s", out)
	}
}

func TestExport_ExportFormat(t *testing.T) {
	env := map[string]string{"KEY": "value"}
	out := Export(env, ExportOptions{Format: FormatExport, Sorted: true})
	if !strings.Contains(out, "export KEY=value\n") {
		t.Errorf("expected export prefix, got: %s", out)
	}
}

func TestExport_DockerFormat(t *testing.T) {
	env := map[string]string{"MY_VAR": "hello world"}
	out := Export(env, ExportOptions{Format: FormatDocker, Sorted: true})
	// Docker format does NOT quote values
	if !strings.Contains(out, "MY_VAR=hello world\n") {
		t.Errorf("expected unquoted docker format, got: %s", out)
	}
}

func TestExport_QuotesSpecialChars(t *testing.T) {
	env := map[string]string{"MSG": "hello world"}
	out := Export(env, ExportOptions{Format: FormatDotenv, Sorted: true})
	if !strings.Contains(out, `MSG="hello world"`) {
		t.Errorf("expected quoted value, got: %s", out)
	}
}

func TestExport_OmitEmpty(t *testing.T) {
	env := map[string]string{
		"PRESENT": "yes",
		"EMPTY":   "",
	}
	out := Export(env, ExportOptions{Format: FormatDotenv, Sorted: true, OmitEmpty: true})
	if strings.Contains(out, "EMPTY") {
		t.Errorf("expected EMPTY to be omitted, got: %s", out)
	}
	if !strings.Contains(out, "PRESENT=yes") {
		t.Errorf("expected PRESENT=yes, got: %s", out)
	}
}

func TestExport_EmptyValueQuoted(t *testing.T) {
	env := map[string]string{"BLANK": ""}
	out := Export(env, ExportOptions{Format: FormatDotenv, Sorted: true})
	if !strings.Contains(out, `BLANK=""`) {
		t.Errorf("expected empty string to be quoted, got: %s", out)
	}
}

func TestExport_SortedOutput(t *testing.T) {
	env := map[string]string{
		"ZZZ": "last",
		"AAA": "first",
	}
	out := Export(env, ExportOptions{Format: FormatDotenv, Sorted: true})
	idxA := strings.Index(out, "AAA")
	idxZ := strings.Index(out, "ZZZ")
	if idxA > idxZ {
		t.Errorf("expected AAA before ZZZ in sorted output")
	}
}
