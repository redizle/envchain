package chain_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/envchain/envchain/internal/chain"
)

func TestLoadConfig_Valid(t *testing.T) {
	dir := t.TempDir()
	content := `{"layers":[{"name":"base","path":".env"},{"name":"local","path":".env.local"}]}`
	p := filepath.Join(dir, "envchain.json")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	cfg, err := chain.LoadConfig(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Layers) != 2 {
		t.Fatalf("expected 2 layers, got %d", len(cfg.Layers))
	}
	if cfg.Layers[0].Name != "base" {
		t.Errorf("layer[0].Name: want %q, got %q", "base", cfg.Layers[0].Name)
	}
}

func TestLoadConfig_NoLayers(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "envchain.json")
	if err := os.WriteFile(p, []byte(`{"layers":[]}`), 0600); err != nil {
		t.Fatal(err)
	}

	_, err := chain.LoadConfig(p)
	if err == nil {
		t.Fatal("expected error for empty layers, got nil")
	}
}

func TestLoadConfig_MissingFile(t *testing.T) {
	_, err := chain.LoadConfig("/nonexistent/envchain.json")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestConfig_ToLayers(t *testing.T) {
	cfg := &chain.Config{
		Layers: []chain.LayerConfig{
			{Name: "base", Path: ".env"},
			{Name: "prod", Path: ".env.prod"},
		},
	}
	layers := cfg.ToLayers()
	if len(layers) != 2 {
		t.Fatalf("expected 2, got %d", len(layers))
	}
	if layers[1].Path != ".env.prod" {
		t.Errorf("layers[1].Path: want %q, got %q", ".env.prod", layers[1].Path)
	}
}
