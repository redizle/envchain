package envfile

import (
	"strings"
	"testing"
)

func TestAudit_CleanEnv(t *testing.T) {
	env := map[string]string{
		"APP_NAME": "myapp",
		"PORT":     "8080",
	}
	result := Audit(env)
	if result.HasIssues() {
		t.Errorf("expected no issues, got: %s", result)
	}
}

func TestAudit_PlainTextSecret(t *testing.T) {
	env := map[string]string{
		"DB_PASSWORD": "hunter2",
	}
	result := Audit(env)
	if len(result.Warnings) == 0 {
		t.Error("expected a warning for plain-text password")
	}
	if !strings.Contains(result.Warnings[0], "DB_PASSWORD") {
		t.Errorf("warning should mention key name, got: %s", result.Warnings[0])
	}
}

func TestAudit_SecretRefIsOk(t *testing.T) {
	env := map[string]string{
		"API_KEY": "${{ssm:/prod/api_key}}",
	}
	result := Audit(env)
	if len(result.Warnings) > 0 {
		t.Errorf("secret ref should not trigger warning, got: %v", result.Warnings)
	}
}

func TestAudit_LiteralNewline(t *testing.T) {
	env := map[string]string{
		"SOME_VAR": "line1\nline2",
	}
	result := Audit(env)
	if len(result.Errors) == 0 {
		t.Error("expected an error for literal newline in value")
	}
}

func TestAudit_LongValue(t *testing.T) {
	env := map[string]string{
		"BIG_DATA": strings.Repeat("x", 600),
	}
	result := Audit(env)
	if len(result.Warnings) == 0 {
		t.Error("expected a warning for unusually long value")
	}
}

func TestAudit_MultipleIssues(t *testing.T) {
	env := map[string]string{
		"SECRET_TOKEN": "abc123",
		"BAD_VAR":      "foo\nbar",
	}
	result := Audit(env)
	if len(result.Warnings) == 0 {
		t.Error("expected warnings")
	}
	if len(result.Errors) == 0 {
		t.Error("expected errors")
	}
}

func TestAuditResult_String(t *testing.T) {
	r := &AuditResult{
		Warnings: []string{"something suspicious"},
		Errors:   []string{"something bad"},
	}
	s := r.String()
	if !strings.Contains(s, "[warn]") {
		t.Error("expected [warn] in output")
	}
	if !strings.Contains(s, "[error]") {
		t.Error("expected [error] in output")
	}
}
