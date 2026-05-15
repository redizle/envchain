package envfile

import (
	"testing"
)

func baseEnv() map[string]string {
	return map[string]string{
		"FOO": "foo",
		"BAR": "bar",
		"BAZ": "baz",
	}
}

func TestPatch_SetNew(t *testing.T) {
	out, res, err := Patch(baseEnv(), []PatchOp{{Action: "set", Key: "NEW", Value: "hello"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["NEW"] != "hello" {
		t.Errorf("expected NEW=hello, got %q", out["NEW"])
	}
	if len(res.Applied) != 1 {
		t.Errorf("expected 1 applied, got %d", len(res.Applied))
	}
}

func TestPatch_SetOverwrite(t *testing.T) {
	out, _, err := Patch(baseEnv(), []PatchOp{{Action: "set", Key: "FOO", Value: "updated"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["FOO"] != "updated" {
		t.Errorf("expected FOO=updated, got %q", out["FOO"])
	}
}

func TestPatch_Delete(t *testing.T) {
	out, res, err := Patch(baseEnv(), []PatchOp{{Action: "delete", Key: "BAR"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, exists := out["BAR"]; exists {
		t.Error("expected BAR to be deleted")
	}
	if len(res.Applied) != 1 {
		t.Errorf("expected 1 applied, got %d", len(res.Applied))
	}
}

func TestPatch_DeleteMissing(t *testing.T) {
	_, res, err := Patch(baseEnv(), []PatchOp{{Action: "delete", Key: "MISSING"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 {
		t.Errorf("expected 1 skipped, got %d", len(res.Skipped))
	}
}

func TestPatch_Rename(t *testing.T) {
	out, res, err := Patch(baseEnv(), []PatchOp{{Action: "rename", Key: "FOO", NewKey: "FOO_RENAMED"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, exists := out["FOO"]; exists {
		t.Error("expected FOO to be removed after rename")
	}
	if out["FOO_RENAMED"] != "foo" {
		t.Errorf("expected FOO_RENAMED=foo, got %q", out["FOO_RENAMED"])
	}
	if len(res.Applied) != 1 {
		t.Errorf("expected 1 applied, got %d", len(res.Applied))
	}
}

func TestPatch_RenameConflictWarning(t *testing.T) {
	_, res, err := Patch(baseEnv(), []PatchOp{{Action: "rename", Key: "FOO", NewKey: "BAR"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Warnings) != 1 {
		t.Errorf("expected 1 warning for overwrite, got %d", len(res.Warnings))
	}
}

func TestPatch_UnknownAction(t *testing.T) {
	_, _, err := Patch(baseEnv(), []PatchOp{{Action: "upsert", Key: "X"}})
	if err == nil {
		t.Error("expected error for unknown action")
	}
}

func TestPatch_DoesNotMutateInput(t *testing.T) {
	env := baseEnv()
	Patch(env, []PatchOp{{Action: "delete", Key: "FOO"}})
	if _, ok := env["FOO"]; !ok {
		t.Error("input env was mutated")
	}
}
