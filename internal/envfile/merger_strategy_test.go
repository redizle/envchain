package envfile

import (
	"testing"
)

func TestMergeWithStrategy_Override(t *testing.T) {
	a := map[string]string{"FOO": "a", "BAR": "1"}
	b := map[string]string{"FOO": "b", "BAZ": "2"}

	got, err := MergeWithStrategy(StrategyOverride, a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["FOO"] != "b" {
		t.Errorf("expected FOO=b, got %q", got["FOO"])
	}
	if got["BAR"] != "1" {
		t.Errorf("expected BAR=1, got %q", got["BAR"])
	}
	if got["BAZ"] != "2" {
		t.Errorf("expected BAZ=2, got %q", got["BAZ"])
	}
}

func TestMergeWithStrategy_KeepFirst(t *testing.T) {
	a := map[string]string{"FOO": "original"}
	b := map[string]string{"FOO": "override", "NEW": "yes"}

	got, err := MergeWithStrategy(StrategyKeepFirst, a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["FOO"] != "original" {
		t.Errorf("expected FOO=original, got %q", got["FOO"])
	}
	if got["NEW"] != "yes" {
		t.Errorf("expected NEW=yes, got %q", got["NEW"])
	}
}

func TestMergeWithStrategy_ErrorOnConflict_NoConflict(t *testing.T) {
	a := map[string]string{"FOO": "same"}
	b := map[string]string{"FOO": "same", "BAR": "new"}

	got, err := MergeWithStrategy(StrategyErrorOnConflict, a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["BAR"] != "new" {
		t.Errorf("expected BAR=new, got %q", got["BAR"])
	}
}

func TestMergeWithStrategy_ErrorOnConflict_Conflict(t *testing.T) {
	a := map[string]string{"FOO": "aaa"}
	b := map[string]string{"FOO": "bbb"}

	_, err := MergeWithStrategy(StrategyErrorOnConflict, a, b)
	if err == nil {
		t.Fatal("expected conflict error, got nil")
	}
	ce, ok := err.(*ConflictError)
	if !ok {
		t.Fatalf("expected *ConflictError, got %T", err)
	}
	if ce.Key != "FOO" {
		t.Errorf("expected conflict key FOO, got %q", ce.Key)
	}
}

func TestMergeWithStrategy_NilLayerSkipped(t *testing.T) {
	a := map[string]string{"X": "1"}

	got, err := MergeWithStrategy(StrategyOverride, a, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["X"] != "1" {
		t.Errorf("expected X=1, got %q", got["X"])
	}
}

func TestConflictError_Message(t *testing.T) {
	err := &ConflictError{Key: "SECRET", Existing: "old", Incoming: "new"}
	msg := err.Error()
	if msg == "" {
		t.Error("expected non-empty error message")
	}
}
