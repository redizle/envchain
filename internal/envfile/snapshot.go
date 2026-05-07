package envfile

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot captures the state of an env map at a point in time.
type Snapshot struct {
	CreatedAt time.Time         `json:"created_at"`
	Label     string            `json:"label"`
	Env       map[string]string `json:"env"`
}

// NewSnapshot creates a snapshot from the given env map and label.
func NewSnapshot(label string, env map[string]string) *Snapshot {
	copy := make(map[string]string, len(env))
	for k, v := range env {
		copy[k] = v
	}
	return &Snapshot{
		CreatedAt: time.Now().UTC(),
		Label:     label,
		Env:       copy,
	}
}

// SaveSnapshot writes a snapshot to a JSON file at the given path.
func SaveSnapshot(path string, s *Snapshot) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("snapshot: write %s: %w", path, err)
	}
	return nil
}

// LoadSnapshot reads a snapshot from a JSON file at the given path.
func LoadSnapshot(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: read %s: %w", path, err)
	}
	var s Snapshot
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("snapshot: unmarshal: %w", err)
	}
	return &s, nil
}

// DiffSnapshot compares a snapshot against a current env map and returns
// the Diff between them (snapshot as base, current as overlay).
func DiffSnapshot(s *Snapshot, current map[string]string) map[string]DiffEntry {
	return Diff(s.Env, current)
}
