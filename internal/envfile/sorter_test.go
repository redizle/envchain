package envfile

import (
	"testing"
)

func TestSort_Alpha(t *testing.T) {
	env := map[string]string{"ZEBRA": "1", "APPLE": "2", "MANGO": "3"}
	_, keys := Sort(env, SortOptions{Order: SortAlpha})
	want := []string{"APPLE", "MANGO", "ZEBRA"}
	for i, k := range keys {
		if k != want[i] {
			t.Errorf("index %d: got %q, want %q", i, k, want[i])
		}
	}
}

func TestSort_AlphaDesc(t *testing.T) {
	env := map[string]string{"ZEBRA": "1", "APPLE": "2", "MANGO": "3"}
	_, keys := Sort(env, SortOptions{Order: SortAlphaDesc})
	want := []string{"ZEBRA", "MANGO", "APPLE"}
	for i, k := range keys {
		if k != want[i] {
			t.Errorf("index %d: got %q, want %q", i, k, want[i])
		}
	}
}

func TestSort_ByGroup(t *testing.T) {
	env := map[string]string{
		"DB_HOST": "localhost",
		"AWS_KEY":  "key",
		"DB_PORT": "5432",
		"AWS_SECRET": "secret",
		"APP_ENV": "prod",
	}
	_, keys := Sort(env, SortOptions{Order: SortByGroup})

	// Verify groups are contiguous
	groups := make([]string, len(keys))
	for i, k := range keys {
		groups[i] = groupPrefix(k)
	}
	for i := 1; i < len(groups); i++ {
		if groups[i] < groups[i-1] {
			t.Errorf("groups out of order at %d: %q before %q", i, groups[i-1], groups[i])
		}
	}
}

func TestSort_IgnoreCase(t *testing.T) {
	env := map[string]string{"zebra": "1", "APPLE": "2", "Mango": "3"}
	_, keys := Sort(env, SortOptions{Order: SortAlpha, IgnoreCase: true})
	want := []string{"APPLE", "Mango", "zebra"}
	for i, k := range keys {
		if k != want[i] {
			t.Errorf("index %d: got %q, want %q", i, k, want[i])
		}
	}
}

func TestSort_DoesNotMutateInput(t *testing.T) {
	env := map[string]string{"B": "2", "A": "1"}
	original := map[string]string{"B": "2", "A": "1"}
	Sort(env, SortOptions{Order: SortAlpha})
	for k, v := range original {
		if env[k] != v {
			t.Errorf("input mutated: key %q", k)
		}
	}
}

func TestGroupPrefix_NoUnderscore(t *testing.T) {
	if got := groupPrefix("NOPREFIX"); got != "NOPREFIX" {
		t.Errorf("got %q, want NOPREFIX", got)
	}
}

func TestGroupPrefix_WithUnderscore(t *testing.T) {
	if got := groupPrefix("DB_HOST"); got != "DB" {
		t.Errorf("got %q, want DB", got)
	}
}
