package envfile

import (
	"fmt"
	"regexp"
	"strings"
)

// SchemaRule defines a validation rule for a specific key.
type SchemaRule struct {
	Key      string
	Required bool
	Pattern  string // optional regex pattern the value must match
	pattern  *regexp.Regexp
}

// Schema holds a set of rules for validating an env map.
type Schema struct {
	Rules []SchemaRule
}

// SchemaViolation describes a single schema violation.
type SchemaViolation struct {
	Key     string
	Message string
}

func (v SchemaViolation) Error() string {
	return fmt.Sprintf("%s: %s", v.Key, v.Message)
}

// ValidateSchema checks env against the schema and returns any violations.
func ValidateSchema(env map[string]string, schema Schema) []SchemaViolation {
	var violations []SchemaViolation

	for i := range schema.Rules {
		rule := &schema.Rules[i]

		// compile pattern once
		if rule.Pattern != "" && rule.pattern == nil {
			compiled, err := regexp.Compile(rule.Pattern)
			if err != nil {
				violations = append(violations, SchemaViolation{
					Key:     rule.Key,
					Message: fmt.Sprintf("invalid pattern in schema rule: %v", err),
				})
				continue
			}
			rule.pattern = compiled
		}

		val, exists := env[rule.Key]

		if rule.Required && (!exists || strings.TrimSpace(val) == "") {
			violations = append(violations, SchemaViolation{
				Key:     rule.Key,
				Message: "required key is missing or empty",
			})
			continue
		}

		if exists && rule.pattern != nil && !rule.pattern.MatchString(val) {
			violations = append(violations, SchemaViolation{
				Key:     rule.Key,
				Message: fmt.Sprintf("value %q does not match pattern %q", val, rule.Pattern),
			})
		}
	}

	return violations
}
