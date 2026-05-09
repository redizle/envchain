package envfile

import (
	"fmt"
	"strings"
)

// LintIssue represents a single linting warning or error.
type LintIssue struct {
	Line    int
	Key     string
	Message string
	Severity string // "warn" or "error"
}

func (i LintIssue) String() string {
	return fmt.Sprintf("[%s] line %d (%s): %s", i.Severity, i.Line, i.Key, i.Message)
}

// Lint inspects a parsed env map and raw lines for common style and correctness issues.
// It returns a slice of LintIssues; an empty slice means the file is clean.
func Lint(env map[string]string, lines []string) []LintIssue {
	var issues []LintIssue

	for lineNum, raw := range lines {
		trimmed := strings.TrimSpace(raw)

		// Skip blank lines and comments
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		// Warn about lines with leading whitespace
		if raw != trimmed {
			issues = append(issues, LintIssue{
				Line:     lineNum + 1,
				Message:  "leading or trailing whitespace around entry",
				Severity: "warn",
			})
		}

		eqIdx := strings.IndexByte(trimmed, '=')
		if eqIdx < 0 {
			continue
		}

		key := trimmed[:eqIdx]
		val := trimmed[eqIdx+1:]

		// Warn about UPPER_SNAKE_CASE convention
		if key != strings.ToUpper(key) {
			issues = append(issues, LintIssue{
				Line:     lineNum + 1,
				Key:      key,
				Message:  "key is not upper-case; convention is UPPER_SNAKE_CASE",
				Severity: "warn",
			})
		}

		// Warn about empty values without a comment explaining intent
		if val == "" {
			issues = append(issues, LintIssue{
				Line:     lineNum + 1,
				Key:      key,
				Message:  "empty value; consider adding a comment or placeholder",
				Severity: "warn",
			})
		}

		// Warn about inline comments that may be accidentally included in value
		if unquotedInlineComment(val) {
			issues = append(issues, LintIssue{
				Line:     lineNum + 1,
				Key:      key,
				Message:  "possible inline comment in unquoted value (use quotes if intentional)",
				Severity: "warn",
			})
		}
	}

	return issues
}

// unquotedInlineComment returns true if val contains ` #` outside of quotes.
func unquotedInlineComment(val string) bool {
	if len(val) == 0 {
		return false
	}
	if val[0] == '"' || val[0] == '\'' {
		return false
	}
	return strings.Contains(val, " #")
}
