package envfile

import (
	"testing"
)

func TestMerge_BasicOverride(t *testing.T) {
	base := map[string]string{"FOO": "bar", "BAZ": "qux"}
	override := map[string]string{"FOO": "overridden", "NEW": "value"}

	result := Merge(base, override)

	if result["FOO"] != "overridden" {
		t.Errorf("expected FOO=overridden, got %s", result["FOO"])
	}
	if result["BAZ"] != "qux" {
		t.Errorf("expected BAZ=qux, got %s", result["BAZ"])
	}
	if result["NEW"] != "value" {
		t.Errorf("expected NEW=value, got %s", result["NEW"])
	}
}

func TestMerge_EmptyLayers(t *testing.T) {
	result := Merge(map[string]string{}, map[string]string{})
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d keys", len(result))
	}
}

func TestMerge_SingleLayer(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2"}
	result := Merge(base)
	if result["A"] != "1" || result["B"] != "2" {
		t.Errorf("single layer merge failed")
	}
}

func TestDiff_AddedRemovedChanged(t *testing.T) {
	base := map[string]string{"FOO": "old", "KEEP": "same", "GONE": "bye"}
	next := map[string]string{"FOO": "new", "KEEP": "same", "FRESH": "hi"}

	d := Diff(base, next)

	if _, ok := d.Added["FRESH"]; !ok {
		t.Error("expected FRESH to be in Added")
	}
	if _, ok := d.Removed["GONE"]; !ok {
		t.Error("expected GONE to be in Removed")
	}
	if pair, ok := d.Changed["FOO"]; !ok || pair[0] != "old" || pair[1] != "new" {
		t.Errorf("expected FOO changed from old to new, got %v", pair)
	}
	if _, ok := d.Changed["KEEP"]; ok {
		t.Error("KEEP should not appear in Changed")
	}
}

func TestDiff_HasChanges(t *testing.T) {
	base := map[string]string{"A": "1"}
	next := map[string]string{"A": "1"}

	d := Diff(base, next)
	if d.HasChanges() {
		t.Error("expected no changes")
	}

	next["B"] = "2"
	d = Diff(base, next)
	if !d.HasChanges() {
		t.Error("expected changes")
	}
}
