package chain

import (
	"fmt"
	"os"

	"github.com/envchain/envchain/internal/envfile"
)

// Layer represents a named .env file layer
type Layer struct {
	Name string
	Path string
}

// Chain holds an ordered list of layers and the merged result
type Chain struct {
	Layers []Layer
	Env    map[string]string
}

// Load reads all layers in order and merges them into a single env map.
// Later layers override earlier ones.
func Load(layers []Layer) (*Chain, error) {
	maps := make([]map[string]string, 0, len(layers))

	for _, layer := range layers {
		data, err := os.ReadFile(layer.Path)
		if err != nil {
			if os.IsNotExist(err) {
				// missing layers are silently skipped
				continue
			}
			return nil, fmt.Errorf("chain: reading layer %q: %w", layer.Name, err)
		}

		parsed, err := envfile.Parse(string(data))
		if err != nil {
			return nil, fmt.Errorf("chain: parsing layer %q: %w", layer.Name, err)
		}

		maps = append(maps, parsed)
	}

	merged := envfile.Merge(maps...)

	return &Chain{
		Layers: layers,
		Env:    merged,
	}, nil
}

// DiffLayers returns the diff between two named layers within the chain.
func DiffLayers(base, override map[string]string) envfile.DiffResult {
	return envfile.Diff(base, override)
}
