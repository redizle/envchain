package envfile

import (
	"strings"
	"testing"
)

func TestLint_CleanFile(t *testing.T) {
	lines := []string{
		"APP_ENV=production",
		"PORT=8080",
		"SECRET_KEY=ref:ssm:/prod/secret",
	}
	env := map[string]string{
		"APP_ENV":    "production",
		"PORT":       "8080",
		"SECRET_KEY": "ref:ssm:/prod/secret",
	}
	issues := Lint(env, lines)
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %d: %v", len(issues), issues)
	}
}

func TestLint_LowercaseKey(t *testing.T) {
	lines := []string{"app_env=production"}
	env := map[string]string{"app_env": "production"}
	issues := Lint(env, lines)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if !strings.Contains(issues[0].Message, "upper-case") {
		t.Errorf("unexpected message: %s", issues[0].Message)
	}
	if issues[0].Severity != "warn" {
		t.Errorf("expected warn severity, got %s", issues[0].Severity)
	}
}

func TestLint_EmptyValue(t *testing.T) {
	lines := []string{"DATABASE_URL="}
	env := map[string]string{"DATABASE_URL": ""}
	issues := Lint(env, lines)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue for empty value, got %d", len(issues))
	}
	if !strings.Contains(issues[0].Message, "empty value") {
		t.Errorf("unexpected message: %s", issues[0].Message)
	}
}

func TestLint_LeadingWhitespace(t *testing.T) {
	lines := []string{"  APP_NAME=envchain"}
	env := map[string]string{"APP_NAME": "envchain"}
	issues := Lint(env, lines)
	// Should have at least the whitespace warning
	found := false
	for _, iss := range issues {
		if strings.Contains(iss.Message, "whitespace") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected whitespace warning, got: %v", issues)
	}
}

func TestLint_InlineComment(t *testing.T) {
	lines := []string{"LOG_LEVEL=debug #this might be a comment"}
	env := map[string]string{"LOG_LEVEL": "debug #this might be a comment"}
	issues := Lint(env, lines)
	found := false
	for _, iss := range issues {
		if strings.Contains(iss.Message, "inline comment") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected inline comment warning, got: %v", issues)
	}
}

func TestLint_QuotedInlineCommentOk(t *testing.T) {
	lines := []string{`GREETING="hello #world"`}
	env := map[string]string{"GREETING": `"hello #world"`}
	issues := Lint(env, lines)
	for _, iss := range issues {
		if strings.Contains(iss.Message, "inline comment") {
			t.Errorf("should not warn about inline comment inside quoted value")
		}
	}
}

func TestLint_SkipsCommentsAndBlanks(t *testing.T) {
	lines := []string{
		"# this is a comment",
		"",
		"APP=ok",
	}
	env := map[string]string{"APP": "ok"}
	issues := Lint(env, lines)
	if len(issues) != 0 {
		t.Errorf("expected no issues, got %v", issues)
	}
}

func TestLintIssue_String(t *testing.T) {
	issue := LintIssue{Line: 3, Key: "foo", Message: "bad key", Severity: "warn"}
	s := issue.String()
	if !strings.Contains(s, "warn") || !strings.Contains(s, "foo") || !strings.Contains(s, "bad key") {
		t.Errorf("unexpected String() output: %s", s)
	}
}
