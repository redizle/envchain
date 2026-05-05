package envfile

// Merge combines multiple env maps in order, with later maps taking precedence.
// Keys from later layers override keys from earlier layers.
func Merge(layers ...map[string]string) map[string]string {
	result := make(map[string]string)
	for _, layer := range layers {
		for k, v := range layer {
			result[k] = v
		}
	}
	return result
}

// DiffResult holds the changes between two env maps.
type DiffResult struct {
	Added   map[string]string
	Removed map[string]string
	Changed map[string][2]string // key -> [old, new]
}

// Diff computes the difference between a base env map and a new env map.
func Diff(base, next map[string]string) DiffResult {
	result := DiffResult{
		Added:   make(map[string]string),
		Removed: make(map[string]string),
		Changed: make(map[string][2]string),
	}

	for k, v := range next {
		oldVal, exists := base[k]
		if !exists {
			result.Added[k] = v
		} else if oldVal != v {
			result.Changed[k] = [2]string{oldVal, v}
		}
	}

	for k, v := range base {
		if _, exists := next[k]; !exists {
			result.Removed[k] = v
		}
	}

	return result
}

// HasChanges returns true if the DiffResult contains any additions, removals, or changes.
func (d DiffResult) HasChanges() bool {
	return len(d.Added) > 0 || len(d.Removed) > 0 || len(d.Changed) > 0
}
