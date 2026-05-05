package chain

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config describes the layer stack for a project, typically loaded from
// envchain.json in the project root.
type Config struct {
	// Layers defines the ordered list of env file paths to merge.
	// Paths are relative to the config file's directory.
	Layers []LayerConfig `json:"layers"`
}

// LayerConfig is the JSON representation of a single layer entry.
type LayerConfig struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

// LoadConfig reads and parses an envchain.json config file.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("chain: reading config %q: %w", path, err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("chain: parsing config %q: %w", path, err)
	}

	if len(cfg.Layers) == 0 {
		return nil, fmt.Errorf("chain: config %q defines no layers", path)
	}

	return &cfg, nil
}

// ToLayers converts config layer entries to Layer values.
func (c *Config) ToLayers() []Layer {
	out := make([]Layer, len(c.Layers))
	for i, l := range c.Layers {
		out[i] = Layer{Name: l.Name, Path: l.Path}
	}
	return out
}
