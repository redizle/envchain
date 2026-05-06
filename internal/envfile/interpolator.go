package envfile

import (
	"fmt"
	"regexp"
	"strings"
)

// refPattern matches ${VAR_NAME} style references within values.
var refPattern = regexp.MustCompile(`\$\{([^}]+)\}`)

// InterpolationError describes a missing variable reference.
type InterpolationError struct {
	Key string
	Ref string
}

func (e *InterpolationError) Error() string {
	return fmt.Sprintf("key %q references undefined variable %q", e.Key, e.Ref)
}

// Interpolate resolves ${VAR} references inside env values using the provided
// env map as the source of truth. References to unknown variables are returned
// as errors. The input map is not mutated; a new map is returned.
func Interpolate(env map[string]string) (map[string]string, []error) {
	result := make(map[string]string, len(env))
	var errs []error

	for key, value := range env {
		resolved, err := interpolateValue(key, value, env)
		if err != nil {
			errs = append(errs, err)
			result[key] = value // keep original on error
			continue
		}
		result[key] = resolved
	}

	return result, errs
}

// interpolateValue replaces all ${REF} tokens in value with their resolved
// counterparts from env.
func interpolateValue(key, value string, env map[string]string) (string, error) {
	var firstErr error

	result := refPattern.ReplaceAllStringFunc(value, func(match string) string {
		// Extract name from ${NAME}
		name := strings.TrimSuffix(strings.TrimPrefix(match, "${"), "}")
		if v, ok := env[name]; ok {
			return v
		}
		if firstErr == nil {
			firstErr = &InterpolationError{Key: key, Ref: name}
		}
		return match // leave token intact
	})

	return result, firstErr
}
