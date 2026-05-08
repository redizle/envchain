package envfile

import (
	"testing"
)

func TestValidateSchema_AllGood(t *testing.T) {
	env := map[string]string{
		"APP_ENV": "production",
		"PORT":    "8080",
	}
	schema := Schema{
		Rules: []SchemaRule{
			{Key: "APP_ENV", Required: true},
			{Key: "PORT", Required: true, Pattern: `^\d+$`},
		},
	}
	violations := ValidateSchema(env, schema)
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %v", violations)
	}
}

func TestValidateSchema_MissingRequired(t *testing.T) {
	env := map[string]string{}
	schema := Schema{
		Rules: []SchemaRule{
			{Key: "DATABASE_URL", Required: true},
		},
	}
	violations := ValidateSchema(env, schema)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Key != "DATABASE_URL" {
		t.Errorf("expected key DATABASE_URL, got %s", violations[0].Key)
	}
}

func TestValidateSchema_PatternMismatch(t *testing.T) {
	env := map[string]string{"PORT": "not-a-number"}
	schema := Schema{
		Rules: []SchemaRule{
			{Key: "PORT", Pattern: `^\d+$`},
		},
	}
	violations := ValidateSchema(env, schema)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
}

func TestValidateSchema_EmptyValueRequired(t *testing.T) {
	env := map[string]string{"SECRET": "   "}
	schema := Schema{
		Rules: []SchemaRule{
			{Key: "SECRET", Required: true},
		},
	}
	violations := ValidateSchema(env, schema)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation for blank value, got %d", len(violations))
	}
}

func TestValidateSchema_OptionalMissingKeySkipsPattern(t *testing.T) {
	env := map[string]string{}
	schema := Schema{
		Rules: []SchemaRule{
			{Key: "LOG_LEVEL", Required: false, Pattern: `^(debug|info|warn|error)$`},
		},
	}
	violations := ValidateSchema(env, schema)
	if len(violations) != 0 {
		t.Fatalf("expected no violations for missing optional key, got %v", violations)
	}
}

func TestValidateSchema_InvalidPatternInRule(t *testing.T) {
	env := map[string]string{"FOO": "bar"}
	schema := Schema{
		Rules: []SchemaRule{
			{Key: "FOO", Pattern: `[invalid(`},
		},
	}
	violations := ValidateSchema(env, schema)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation for bad pattern, got %d", len(violations))
	}
}
