package envfile

import (
	"testing"
)

func TestInterpolate_NoRefs(t *testing.T) {
	env := map[string]string{
		"HOST": "localhost",
		"PORT": "5432",
	}
	result, errs := Interpolate(env)
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
	if result["HOST"] != "localhost" || result["PORT"] != "5432" {
		t.Errorf("unexpected result: %v", result)
	}
}

func TestInterpolate_BasicRef(t *testing.T) {
	env := map[string]string{
		"BASE": "postgres",
		"URL":  "${BASE}://localhost",
	}
	result, errs := Interpolate(env)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if result["URL"] != "postgres://localhost" {
		t.Errorf("got %q, want %q", result["URL"], "postgres://localhost")
	}
}

func TestInterpolate_MultipleRefs(t *testing.T) {
	env := map[string]string{
		"SCHEME": "https",
		"HOST":   "example.com",
		"URL":    "${SCHEME}://${HOST}/path",
	}
	result, errs := Interpolate(env)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if result["URL"] != "https://example.com/path" {
		t.Errorf("got %q", result["URL"])
	}
}

func TestInterpolate_UnknownRef(t *testing.T) {
	env := map[string]string{
		"URL": "${UNDEFINED}://host",
	}
	result, errs := Interpolate(env)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
	ie, ok := errs[0].(*InterpolationError)
	if !ok {
		t.Fatalf("expected *InterpolationError, got %T", errs[0])
	}
	if ie.Key != "URL" || ie.Ref != "UNDEFINED" {
		t.Errorf("unexpected error fields: %+v", ie)
	}
	// Original value preserved
	if result["URL"] != "${UNDEFINED}://host" {
		t.Errorf("expected original value preserved, got %q", result["URL"])
	}
}

func TestInterpolate_EmptyMap(t *testing.T) {
	result, errs := Interpolate(map[string]string{})
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(result) != 0 {
		t.Errorf("expected empty result, got %v", result)
	}
}

func TestInterpolate_OriginalNotMutated(t *testing.T) {
	env := map[string]string{
		"A": "hello",
		"B": "${A} world",
	}
	_, _ = Interpolate(env)
	if env["B"] != "${A} world" {
		t.Errorf("original map was mutated")
	}
}
