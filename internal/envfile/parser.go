package envfile

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Entry represents a single key-value pair in an env file.
type Entry struct {
	Key     string
	Value   string
	Comment string // inline or preceding comment
}

// EnvFile holds all parsed entries from a .env file.
type EnvFile struct {
	Entries []Entry
}

// Parse reads an env file from r and returns a structured EnvFile.
// Supports KEY=VALUE, # comments, and quoted values.
func Parse(r io.Reader) (*EnvFile, error) {
	ef := &EnvFile{}
	scanner := bufio.NewScanner(r)
	lineNum := 0
	var pendingComment string

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			pendingComment = ""
			continue
		}

		if strings.HasPrefix(line, "#") {
			pendingComment = strings.TrimPrefix(line, "#")
			pendingComment = strings.TrimSpace(pendingComment)
			continue
		}

		idx := strings.IndexByte(line, '=')
		if idx < 0 {
			return nil, fmt.Errorf("line %d: missing '=' in %q", lineNum, line)
		}

		key := strings.TrimSpace(line[:idx])
		raw := strings.TrimSpace(line[idx+1:])

		if key == "" {
			return nil, fmt.Errorf("line %d: empty key", lineNum)
		}

		value := unquote(raw)

		ef.Entries = append(ef.Entries, Entry{
			Key:     key,
			Value:   value,
			Comment: pendingComment,
		})
		pendingComment = ""
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %w", err)
	}

	return ef, nil
}

// ToMap returns a flat key→value map from the parsed entries.
func (ef *EnvFile) ToMap() map[string]string {
	m := make(map[string]string, len(ef.Entries))
	for _, e := range ef.Entries {
		m[e.Key] = e.Value
	}
	return m
}

// unquote strips surrounding single or double quotes from a value.
func unquote(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
