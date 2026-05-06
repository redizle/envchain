package output

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/envchain/envchain/internal/envfile"
)

// Format controls how env vars are written.
type Format string

const (
	FormatEnv    Format = "env"    // KEY=VALUE
	FormatExport Format = "export" // export KEY=VALUE
	FormatJSON   Format = "json"   // {"KEY": "VALUE"}
	FormatDotenv Format = "dotenv" // same as env but with quotes
)

// WriteEnv writes a map of env vars to w in the given format.
func WriteEnv(w io.Writer, vars map[string]string, format Format) error {
	keys := sortedKeys(vars)

	switch format {
	case FormatExport:
		for _, k := range keys {
			_, err := fmt.Fprintf(w, "export %s=%s\n", k, shellQuote(vars[k]))
			if err != nil {
				return err
			}
		}
	case FormatJSON:
		return writeJSON(w, vars, keys)
	case FormatDotenv:
		for _, k := range keys {
			_, err := fmt.Fprintf(w, "%s=%q\n", k, vars[k])
			if err != nil {
				return err
			}
		}
	default: // FormatEnv
		for _, k := range keys {
			_, err := fmt.Fprintf(w, "%s=%s\n", k, vars[k])
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// WriteDiff writes a human-readable diff to w.
func WriteDiff(w io.Writer, diff map[string]envfile.Change) error {
	keys := make([]string, 0, len(diff))
	for k := range diff {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		c := diff[k]
		switch c.Type {
		case envfile.ChangeAdded:
			_, err := fmt.Fprintf(w, "+ %s=%s\n", k, c.New)
			if err != nil {
				return err
			}
		case envfile.ChangeRemoved:
			_, err := fmt.Fprintf(w, "- %s=%s\n", k, c.Old)
			if err != nil {
				return err
			}
		case envfile.ChangeModified:
			_, err := fmt.Fprintf(w, "~ %s: %s -> %s\n", k, c.Old, c.New)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func shellQuote(s string) string {
	if !strings.ContainsAny(s, " \t\n\r'\"\\$`") {
		return s
	}
	return `'` + strings.ReplaceAll(s, `'`, `'\''`) + `'`
}

func writeJSON(w io.Writer, vars map[string]string, keys []string) error {
	_, err := fmt.Fprintln(w, "{")
	if err != nil {
		return err
	}
	for i, k := range keys {
		comma := ","
		if i == len(keys)-1 {
			comma = ""
		}
		_, err = fmt.Fprintf(w, "  %q: %q%s\n", k, vars[k], comma)
		if err != nil {
			return err
		}
	}
	_, err = fmt.Fprintln(w, "}")
	return err
}
