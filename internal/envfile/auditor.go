package envfile

import (
	"fmt"
	"strings"
)

// AuditResult holds findings from auditing an env map.
type AuditResult struct {
	Warnings []string
	Errors   []string
}

// HasIssues returns true if there are any warnings or errors.
func (a *AuditResult) HasIssues() bool {
	return len(a.Warnings) > 0 || len(a.Errors) > 0
}

// String returns a human-readable summary.
func (a *AuditResult) String() string {
	var sb strings.Builder
	for _, e := range a.Errors {
		fmt.Fprintf(&sb, "[error] %s\n", e)
	}
	for _, w := range a.Warnings {
		fmt.Fprintf(&sb, "[warn]  %s\n", w)
	}
	return sb.String()
}

// sensitivePatterns are key substrings that suggest a value should not be plain.
var sensitivePatterns = []string{
	"PASSWORD", "SECRET", "TOKEN", "API_KEY", "PRIVATE_KEY", "CREDENTIALS",
}

// Audit inspects an env map for common security and quality issues.
func Audit(env map[string]string) *AuditResult {
	result := &AuditResult{}

	for k, v := range env {
		upper := strings.ToUpper(k)

		// Warn if a sensitive-looking key holds a plain (non-ref) non-empty value.
		if isSensitiveKey(upper) && v != "" && !isSecretRef(v) {
			result.Warnings = append(result.Warnings,
				fmt.Sprintf("key %q looks sensitive but has a plain-text value", k))
		}

		// Error if value contains literal newlines (likely a paste mistake).
		if strings.Contains(v, "\n") {
			result.Errors = append(result.Errors,
				fmt.Sprintf("key %q value contains a literal newline", k))
		}

		// Warn on very long values that are not secret refs.
		if len(v) > 512 && !isSecretRef(v) {
			result.Warnings = append(result.Warnings,
				fmt.Sprintf("key %q has an unusually long value (%d chars)", k, len(v)))
		}
	}

	return result
}

func isSensitiveKey(upper string) bool {
	for _, p := range sensitivePatterns {
		if strings.Contains(upper, p) {
			return true
		}
	}
	return false
}

// isSecretRef returns true if the value looks like a secret interpolation ref.
func isSecretRef(v string) bool {
	return strings.HasPrefix(v, "${{") && strings.HasSuffix(v, "}}")
}
