package output

import (
	"strings"
	"testing"

	"github.com/envchain/envchain/internal/envfile"
)

func TestWriteEnv_DefaultFormat(t *testing.T) {
	vars := map[string]string{"FOO": "bar", "BAZ": "qux"}
	var buf strings.Builder
	if err := WriteEnv(&buf, vars, FormatEnv); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "FOO=bar") {
		t.Errorf("expected FOO=bar in output, got: %s", out)
	}
	if !strings.Contains(out, "BAZ=qux") {
		t.Errorf("expected BAZ=qux in output, got: %s", out)
	}
}

func TestWriteEnv_ExportFormat(t *testing.T) {
	vars := map[string]string{"KEY": "value with spaces"}
	var buf strings.Builder
	if err := WriteEnv(&buf, vars, FormatExport); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "export KEY=") {
		t.Errorf("expected export prefix, got: %s", out)
	}
}

func TestWriteEnv_JSONFormat(t *testing.T) {
	vars := map[string]string{"A": "1"}
	var buf strings.Builder
	if err := WriteEnv(&buf, vars, FormatJSON); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, `"A"`) || !strings.Contains(out, `"1"`) {
		t.Errorf("unexpected JSON output: %s", out)
	}
}

func TestWriteEnv_DotenvFormat(t *testing.T) {
	vars := map[string]string{"X": "hello world"}
	var buf strings.Builder
	if err := WriteEnv(&buf, vars, FormatDotenv); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, `X="hello world"`) {
		t.Errorf("expected quoted value, got: %s", out)
	}
}

func TestWriteDiff(t *testing.T) {
	diff := map[string]envfile.Change{
		"ADDED": {Type: envfile.ChangeAdded, New: "newval"},
		"REMOVED": {Type: envfile.ChangeRemoved, Old: "oldval"},
		"CHANGED": {Type: envfile.ChangeModified, Old: "v1", New: "v2"},
	}
	var buf strings.Builder
	if err := WriteDiff(&buf, diff); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "+ ADDED=newval") {
		t.Errorf("missing added line: %s", out)
	}
	if !strings.Contains(out, "- REMOVED=oldval") {
		t.Errorf("missing removed line: %s", out)
	}
	if !strings.Contains(out, "~ CHANGED: v1 -> v2") {
		t.Errorf("missing changed line: %s", out)
	}
}

func TestShellQuote(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"simple", "simple"},
		{"with space", "'with space'"},
		{"it's", `'it'\''s'`},
	}
	for _, c := range cases {
		got := shellQuote(c.input)
		if got != c.expected {
			t.Errorf("shellQuote(%q) = %q, want %q", c.input, got, c.expected)
		}
	}
}
