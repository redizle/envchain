package envfile

import (
	"fmt"
	"strings"
	"unicode"
)

// ValidationError holds a list of issues found during validation.
type ValidationError struct {
	Issues []string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed: %s", strings.Join(e.Issues, "; "))
}

// Validate checks a parsed env map for common issues.
// It returns a *ValidationError if any issues are found, nil otherwise.
func Validate(env map[string]string) error {
	var issues []string

	for key, value := range env {
		if err := validateKey(key); err != nil {
			issues = append(issues, err.Error())
		}
		if err := validateValue(key, value); err != nil {
			issues = append(issues, err.Error())
		}
	}

	if len(issues) > 0 {
		return &ValidationError{Issues: issues}
	}
	return nil
}

func validateKey(key string) error {
	if key == "" {
		return fmt.Errorf("empty key found")
	}
	for i, ch := range key {
		if i == 0 && unicode.IsDigit(ch) {
			return fmt.Errorf("key %q must not start with a digit", key)
		}
		if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) && ch != '_' {
			return fmt.Errorf("key %q contains invalid character %q", key, ch)
		}
	}
	return nil
}

func validateValue(key, value string) error {
	if strings.ContainsAny(value, "\x00") {
		return fmt.Errorf("value for key %q contains null byte", key)
	}
	return nil
}
