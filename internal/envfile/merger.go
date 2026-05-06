package envfile

// ChangeType describes what happened to a key between two env maps.
type ChangeType string

const (
	ChangeAdded    ChangeType = "added"
	ChangeRemoved  ChangeType = "removed"
	ChangeModified ChangeType = "modified"
)

// Change represents a single key-level difference between two env maps.
type Change struct {
	Type ChangeType
	Old  string
	New  string
}

// Merge combines multiple layers of env maps, with later layers taking precedence.
func Merge(layers ...map[string]string) map[string]string {
	result := make(map[string]string)
	for _, layer := range layers {
		for k, v := range layer {
			result[k] = v
		}
	}
	return result
}

// Diff computes the changes between a base env map and a new env map.
// Returns a map of key -> Change for any keys that were added, removed, or modified.
func Diff(base, next map[string]string) map[string]Change {
	changes := make(map[string]Change)

	for k, v := range next {
		if oldVal, ok := base[k]; !ok {
			changes[k] = Change{Type: ChangeAdded, New: v}
		} else if oldVal != v {
			changes[k] = Change{Type: ChangeModified, Old: oldVal, New: v}
		}
	}

	for k, v := range base {
		if _, ok := next[k]; !ok {
			changes[k] = Change{Type: ChangeRemoved, Old: v}
		}
	}

	return changes
}

// HasChanges returns true if the diff map contains any entries.
func HasChanges(diff map[string]Change) bool {
	return len(diff) > 0
}
