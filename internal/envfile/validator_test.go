package envfile

import (
	"strings"
	"testing"
)

func TestValidate_ValidEnv(t *testing.T) {
	env := map[string]string{
		"APP_NAME": "envchain",
		"PORT":     "8080",
		"_SECRET":  "abc123",
	}
	if err := Validate(env); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidate_KeyStartsWithDigit(t *testing.T) {
	env := map[string]string{"1INVALID": "value"}
	err := Validate(env)
	if err == nil {
		t.Fatal("expected error for digit-leading key")
	}
	if !strings.Contains(err.Error(), "must not start with a digit") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestValidate_KeyWithInvalidChar(t *testing.T) {
	env := map[string]string{"BAD-KEY": "value"}
	err := Validate(env)
	if err == nil {
		t.Fatal("expected error for key with hyphen")
	}
	if !strings.Contains(err.Error(), "invalid character") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestValidate_ValueWithNullByte(t *testing.T) {
	env := map[string]string{"KEY": "val\x00ue"}
	err := Validate(env)
	if err == nil {
		t.Fatal("expected error for null byte in value")
	}
	if !strings.Contains(err.Error(), "null byte") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestValidate_MultipleIssues(t *testing.T) {
	env := map[string]string{
		"1BAD":    "ok",
		"ALSO-BAD": "ok",
	}
	err := Validate(env)
	if err == nil {
		t.Fatal("expected error for multiple bad keys")
	}
	ve, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if len(ve.Issues) < 2 {
		t.Errorf("expected at least 2 issues, got %d", len(ve.Issues))
	}
}

func TestValidate_EmptyMap(t *testing.T) {
	if err := Validate(map[string]string{}); err != nil {
		t.Fatalf("expected no error for empty map, got %v", err)
	}
}
