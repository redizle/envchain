package envfile

import (
	"sort"
	"testing"
)

func TestRedact_NonSensitiveKeysUnchanged(t *testing.T) {
	env := map[string]string{
		"APP_NAME": "myapp",
		"PORT":     "8080",
	}
	out := Redact(env, RedactOptions{})
	if out["APP_NAME"] != "myapp" {
		t.Errorf("expected APP_NAME=myapp, got %s", out["APP_NAME"])
	}
	if out["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %s", out["PORT"])
	}
}

func TestRedact_SensitiveKeysRedacted(t *testing.T) {
	env := map[string]string{
		"DB_PASSWORD": "supersecret",
		"API_SECRET":  "abc123",
		"AUTH_TOKEN":  "tok_xyz",
	}
	out := Redact(env, RedactOptions{})
	for k, v := range out {
		if v != RedactedValue {
			t.Errorf("expected %s to be redacted, got %s", k, v)
		}
	}
}

func TestRedact_PreserveRefs_KeepsSecretRef(t *testing.T) {
	env := map[string]string{
		"DB_PASSWORD": "${ssm:/prod/db/password}",
		"API_SECRET":  "plaintext",
	}
	out := Redact(env, RedactOptions{PreserveRefs: true})
	if out["DB_PASSWORD"] != "${ssm:/prod/db/password}" {
		t.Errorf("expected ref preserved, got %s", out["DB_PASSWORD"])
	}
	if out["API_SECRET"] != RedactedValue {
		t.Errorf("expected plaintext secret redacted, got %s", out["API_SECRET"])
	}
}

func TestRedact_PreserveRefs_False_RedactsRef(t *testing.T) {
	env := map[string]string{
		"DB_PASSWORD": "${ssm:/prod/db/password}",
	}
	out := Redact(env, RedactOptions{PreserveRefs: false})
	if out["DB_PASSWORD"] != RedactedValue {
		t.Errorf("expected ref to be redacted, got %s", out["DB_PASSWORD"])
	}
}

func TestRedact_DoesNotMutateInput(t *testing.T) {
	env := map[string]string{"DB_PASSWORD": "secret"}
	_ = Redact(env, RedactOptions{})
	if env["DB_PASSWORD"] != "secret" {
		t.Error("input map was mutated")
	}
}

func TestRedactSlice_ReturnsPairs(t *testing.T) {
	env := map[string]string{
		"APP_NAME": "myapp",
		"DB_PASSWORD": "secret",
	}
	lines := RedactSlice(env, RedactOptions{})
	sort.Strings(lines)
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	if lines[0] != "APP_NAME=myapp" {
		t.Errorf("unexpected line: %s", lines[0])
	}
	if lines[1] != "DB_PASSWORD="+RedactedValue {
		t.Errorf("unexpected line: %s", lines[1])
	}
}
