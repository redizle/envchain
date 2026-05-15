package envfile

import (
	"fmt"
	"strings"
)

// PatchOp represents a single patch operation on an env map.
type PatchOp struct {
	Action string // "set", "delete", "rename"
	Key    string
	Value  string // used by "set"
	NewKey string // used by "rename"
}

// PatchResult holds the outcome of applying a patch.
type PatchResult struct {
	Applied  []string
	Skipped  []string
	Warnings []string
}

// Patch applies a slice of PatchOps to a copy of the provided env map.
// It returns the patched map and a PatchResult describing what happened.
func Patch(env map[string]string, ops []PatchOp) (map[string]string, PatchResult, error) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}

	var result PatchResult

	for _, op := range ops {
		switch strings.ToLower(op.Action) {
		case "set":
			if op.Key == "" {
				return nil, result, fmt.Errorf("set op missing key")
			}
			out[op.Key] = op.Value
			result.Applied = append(result.Applied, fmt.Sprintf("set %s", op.Key))

		case "delete":
			if _, exists := out[op.Key]; !exists {
				result.Skipped = append(result.Skipped, fmt.Sprintf("delete %s (not found)", op.Key))
				continue
			}
			delete(out, op.Key)
			result.Applied = append(result.Applied, fmt.Sprintf("delete %s", op.Key))

		case "rename":
			if op.NewKey == "" {
				return nil, result, fmt.Errorf("rename op missing new_key for key %q", op.Key)
			}
			v, exists := out[op.Key]
			if !exists {
				result.Skipped = append(result.Skipped, fmt.Sprintf("rename %s (not found)", op.Key))
				continue
			}
			if _, conflict := out[op.NewKey]; conflict {
				result.Warnings = append(result.Warnings, fmt.Sprintf("rename %s -> %s overwrites existing key", op.Key, op.NewKey))
			}
			out[op.NewKey] = v
			delete(out, op.Key)
			result.Applied = append(result.Applied, fmt.Sprintf("rename %s -> %s", op.Key, op.NewKey))

		default:
			return nil, result, fmt.Errorf("unknown patch action %q", op.Action)
		}
	}

	return out, result, nil
}
