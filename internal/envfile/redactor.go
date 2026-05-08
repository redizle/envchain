package envfile

import (
	"regexp"
	"strings"
)

// RedactedValue is the placeholder used in place of sensitive values.
const RedactedValue = "***REDACTED***"

// RedactOptions controls redaction behaviour.
type RedactOptions struct {
	// PreserveRefs keeps secret reference syntax (e.g. ${ssm:/path}) unredacted.
	PreserveRefs bool
}

var secretRefPattern = regexp.MustCompile(`^\$\{[a-zA-Z0-9_]+:.+\}$`)

// Redact returns a copy of env where all sensitive values are replaced with
// RedactedValue. Secret references are optionally preserved based on opts.
func Redact(env map[string]string, opts RedactOptions) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		if opts.PreserveRefs && secretRefPattern.MatchString(strings.TrimSpace(v)) {
			out[k] = v
			continue
		}
		if isSensitiveKey(k) {
			out[k] = RedactedValue
			continue
		}
		out[k] = v
	}
	return out
}

// RedactSlice returns a []string slice of KEY=VALUE pairs with sensitive
// values redacted, suitable for display or logging.
func RedactSlice(env map[string]string, opts RedactOptions) []string {
	redacted := Redact(env, opts)
	lines := make([]string, 0, len(redacted))
	for k, v := range redacted {
		lines = append(lines, k+"="+v)
	}
	return lines
}
