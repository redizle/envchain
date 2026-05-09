package envfile

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

// TemplateOptions controls rendering behaviour.
type TemplateOptions struct {
	// MissingKey controls how missing keys are handled: "zero", "error", or "default".
	MissingKey string
}

// RenderTemplate renders a Go text/template string against the provided env map.
// Keys in the env map are available as {{ .KEY }}.
func RenderTemplate(tmpl string, env map[string]string, opts TemplateOptions) (string, error) {
	missingKey := opts.MissingKey
	if missingKey == "" {
		missingKey = "error"
	}

	t, err := template.New("envchain").
		Option("missingkey=" + missingKey).
		Funcs(templateFuncs()).
		Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("template parse error: %w", err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, envToData(env)); err != nil {
		return "", fmt.Errorf("template render error: %w", err)
	}
	return buf.String(), nil
}

// RenderTemplateMap renders every value in env that contains "{{" as a template,
// using the full env map as context. Non-template values are returned as-is.
func RenderTemplateMap(env map[string]string, opts TemplateOptions) (map[string]string, error) {
	out := make(map[string]string, len(env))
	data := envToData(env)
	for k, v := range env {
		if !strings.Contains(v, "{{") {
			out[k] = v
			continue
		}
		rendered, err := renderOne(v, data, opts)
		if err != nil {
			return nil, fmt.Errorf("key %q: %w", k, err)
		}
		out[k] = rendered
	}
	return out, nil
}

func renderOne(tmpl string, data map[string]string, opts TemplateOptions) (string, error) {
	missingKey := opts.MissingKey
	if missingKey == "" {
		missingKey = "error"
	}
	t, err := template.New("").
		Option("missingkey=" + missingKey).
		Funcs(templateFuncs()).
		Parse(tmpl)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func envToData(env map[string]string) map[string]string {
	copy := make(map[string]string, len(env))
	for k, v := range env {
		copy[k] = v
	}
	return copy
}

func templateFuncs() template.FuncMap {
	return template.FuncMap{
		"default": func(def, val string) string {
			if val == "" {
				return def
			}
			return val
		},
		"upper": strings.ToUpper,
		"lower": strings.ToLower,
	}
}
