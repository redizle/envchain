package envfile

import (
	"strings"
	"testing"
)

func TestParse_BasicKeyValue(t *testing.T) {
	input := `APP_ENV=production
DATABASE_URL=postgres://localhost/mydb
`
	ef, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ef.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(ef.Entries))
	}
	if ef.Entries[0].Key != "APP_ENV" || ef.Entries[0].Value != "production" {
		t.Errorf("unexpected first entry: %+v", ef.Entries[0])
	}
}

func TestParse_QuotedValues(t *testing.T) {
	input := `GREETING="hello world"
NAME='envchain'
`
	ef, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ef.Entries[0].Value != "hello world" {
		t.Errorf("expected 'hello world', got %q", ef.Entries[0].Value)
	}
	if ef.Entries[1].Value != "envchain" {
		t.Errorf("expected 'envchain', got %q", ef.Entries[1].Value)
	}
}

func TestParse_Comments(t *testing.T) {
	input := `# database config
DB_HOST=localhost
`
	ef, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ef.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(ef.Entries))
	}
	if ef.Entries[0].Comment != "database config" {
		t.Errorf("expected comment 'database config', got %q", ef.Entries[0].Comment)
	}
}

func TestParse_MissingEquals(t *testing.T) {
	input := `BADLINE
`
	_, err := Parse(strings.NewReader(input))
	if err == nil {
		t.Fatal("expected error for missing '=', got nil")
	}
}

func TestParse_EmptyKey(t *testing.T) {
	input := `=value
`
	_, err := Parse(strings.NewReader(input))
	if err == nil {
		t.Fatal("expected error for empty key, got nil")
	}
}

func TestToMap(t *testing.T) {
	input := `FOO=bar
BAZ=qux
`
	ef, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := ef.ToMap()
	if m["FOO"] != "bar" || m["BAZ"] != "qux" {
		t.Errorf("unexpected map contents: %v", m)
	}
}
