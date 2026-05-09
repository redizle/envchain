package envfile

import (
	"testing"
)

func TestRenderTemplate_NoTemplate(t *testing.T) {
	out, err := RenderTemplate("hello world", nil, TemplateOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if out != "hello world" {
		t.Errorf("expected %q, got %q", "hello world", out)
	}
}

func TestRenderTemplate_BasicSubstitution(t *testing.T) {
	env := map[string]string{"NAME": "envchain"}
	out, err := RenderTemplate("hello {{ .NAME }}", env, TemplateOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if out != "hello envchain" {
		t.Errorf("got %q", out)
	}
}

func TestRenderTemplate_MissingKeyError(t *testing.T) {
	_, err := RenderTemplate("{{ .MISSING }}", map[string]string{}, TemplateOptions{MissingKey: "error"})
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestRenderTemplate_MissingKeyZero(t *testing.T) {
	out, err := RenderTemplate("{{ .MISSING }}", map[string]string{}, TemplateOptions{MissingKey: "zero"})
	if err != nil {
		t.Fatal(err)
	}
	if out != "" {
		t.Errorf("expected empty string, got %q", out)
	}
}

func TestRenderTemplate_DefaultFunc(t *testing.T) {
	out, err := RenderTemplate(`{{ default "fallback" .MISSING }}`, map[string]string{}, TemplateOptions{MissingKey: "zero"})
	if err != nil {
		t.Fatal(err)
	}
	if out != "fallback" {
		t.Errorf("got %q", out)
	}
}

func TestRenderTemplate_UpperFunc(t *testing.T) {
	out, err := RenderTemplate(`{{ upper .VAL }}`, map[string]string{"VAL": "hello"}, TemplateOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if out != "HELLO" {
		t.Errorf("got %q", out)
	}
}

func TestRenderTemplateMap_RendersTemplateValues(t *testing.T) {
	env := map[string]string{
		"BASE_URL": "https://example.com",
		"API_URL":  `{{ .BASE_URL }}/api`,
		"PLAIN":    "no-template",
	}
	out, err := RenderTemplateMap(env, TemplateOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if out["API_URL"] != "https://example.com/api" {
		t.Errorf("API_URL = %q", out["API_URL"])
	}
	if out["PLAIN"] != "no-template" {
		t.Errorf("PLAIN = %q", out["PLAIN"])
	}
}

func TestRenderTemplateMap_ErrorOnMissingKey(t *testing.T) {
	env := map[string]string{
		"X": `{{ .DOES_NOT_EXIST }}`,
	}
	_, err := RenderTemplateMap(env, TemplateOptions{MissingKey: "error"})
	if err == nil {
		t.Fatal("expected error")
	}
}
