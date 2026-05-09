package envfile

import (
	"fmt"
	"sort"
	"strings"
)

// ExportFormat defines the output format for exported env vars.
type ExportFormat int

const (
	FormatDotenv ExportFormat = iota
	FormatExport
	FormatDocker
)

// ExportOptions controls how the env map is serialized.
type ExportOptions struct {
	Format  ExportFormat
	Sorted  bool
	OmitEmpty bool
}

// Export serializes an env map to a string in the requested format.
func Export(env map[string]string, opts ExportOptions) string {
	keys := make([]string, 0, len(env))
	for k := range env {
		if opts.OmitEmpty && env[k] == "" {
			continue
		}
		keys = append(keys, k)
	}

	if opts.Sorted {
		sort.Strings(keys)
	}

	var sb strings.Builder
	for _, k := range keys {
		v := env[k]
		switch opts.Format {
		case FormatExport:
			fmt.Fprintf(&sb, "export %s=%s\n", k, quoteValue(v))
		case FormatDocker:
			// Docker --env-file format: no quoting, plain KEY=VALUE
			fmt.Fprintf(&sb, "%s=%s\n", k, v)
		default: // FormatDotenv
			fmt.Fprintf(&sb, "%s=%s\n", k, quoteValue(v))
		}
	}
	return sb.String()
}

// quoteValue wraps a value in double quotes if it contains
// spaces, special characters, or is empty.
func quoteValue(v string) string {
	if v == "" {
		return `""`
	}
	needsQuote := strings.ContainsAny(v, " \t\n\r\"'\\$`!#&;|<>(){}")
	if !needsQuote {
		return v
	}
	escaped := strings.ReplaceAll(v, `\`, `\\`)
	escaped = strings.ReplaceAll(escaped, `"`, `\"`)
	return `"` + escaped + `"`
}
