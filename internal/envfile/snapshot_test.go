package envfile_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nicholasgasior/envchain/internal/envfile"
)

func TestNewSnapshot_CopiesEnv(t *testing.T) {
	original := map[string]string{"FOO": "bar", "BAZ": "qux"}
	s := envfile.NewSnapshot("test", original)

	if s.Label != "test" {
		t.Fatalf("expected label 'test', got %q", s.Label)
	}
	if s.Env["FOO"] != "bar" {
		t.Fatalf("expected FOO=bar")
	}

	// Mutating original should not affect snapshot
	original["FOO"] = "mutated"
	if s.Env["FOO"] != "bar" {
		t.Fatal("snapshot env was mutated by original map change")
	}
}

func TestSaveAndLoadSnapshot(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	env := map[string]string{"KEY": "value", "NUM": "42"}
	s := envfile.NewSnapshot("prod", env)

	if err := envfile.SaveSnapshot(path, s); err != nil {
		t.Fatalf("SaveSnapshot error: %v", err)
	}

	loaded, err := envfile.LoadSnapshot(path)
	if err != nil {
		t.Fatalf("LoadSnapshot error: %v", err)
	}

	if loaded.Label != "prod" {
		t.Fatalf("expected label 'prod', got %q", loaded.Label)
	}
	if loaded.Env["KEY"] != "value" {
		t.Fatalf("expected KEY=value")
	}
	if loaded.Env["NUM"] != "42" {
		t.Fatalf("expected NUM=42")
	}
}

func TestLoadSnapshot_MissingFile(t *testing.T) {
	_, err := envfile.LoadSnapshot("/nonexistent/snap.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestDiffSnapshot_DetectsChanges(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2", "C": "3"}
	s := envfile.NewSnapshot("base", base)

	current := map[string]string{"A": "1", "B": "changed", "D": "new"}
	diff := envfile.DiffSnapshot(s, current)

	if diff["B"].Old != "2" || diff["B"].New != "changed" {
		t.Fatalf("expected B to be changed, got %+v", diff["B"])
	}
	if diff["C"].Old != "3" || diff["C"].New != "" {
		t.Fatalf("expected C to be removed, got %+v", diff["C"])
	}
	if diff["D"].Old != "" || diff["D"].New != "new" {
		t.Fatalf("expected D to be added, got %+v", diff["D"])
	}
	if _, ok := diff["A"]; ok {
		t.Fatal("expected A to be unchanged and absent from diff")
	}
}

func TestSaveSnapshot_UnwritablePath(t *testing.T) {
	s := envfile.NewSnapshot("x", map[string]string{})
	err := envfile.SaveSnapshot("/no/such/dir/snap.json", s)
	if err == nil {
		t.Fatal("expected error writing to unwritable path")
	}
	_ = os.Remove("/no/such/dir/snap.json")
}
