// Package output provides utilities for rendering env var maps and diffs
// in multiple formats suitable for shell consumption, inspection, or export.
//
// Supported formats:
//
//   - env:    plain KEY=VALUE lines (default, suitable for sourcing)
//   - export: export KEY=VALUE lines (suitable for eval in bash/zsh)
//   - dotenv: KEY="VALUE" with Go-style quoting
//   - json:   { "KEY": "VALUE" } JSON object
//
// Example usage:
//
//	vars := map[string]string{"FOO": "bar"}
//	output.WriteEnv(os.Stdout, vars, output.FormatExport)
package output
