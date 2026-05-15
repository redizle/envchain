package envfile

// MergeStrategy controls how values are combined when merging layers.
type MergeStrategy int

const (
	// StrategyOverride replaces existing values with new ones (default).
	StrategyOverride MergeStrategy = iota
	// StrategyKeepFirst retains the first-seen value and ignores later ones.
	StrategyKeepFirst
	// StrategyErrorOnConflict returns an error if two layers define the same key with different values.
	StrategyErrorOnConflict
)

// ConflictError is returned by MergeWithStrategy when StrategyErrorOnConflict
// detects a duplicate key with a differing value.
type ConflictError struct {
	Key      string
	Existing string
	Incoming string
}

func (e *ConflictError) Error() string {
	return "envchain: merge conflict on key \"" + e.Key + "\": existing=\"" + e.Existing + "\" incoming=\"" + e.Incoming + "\""
}

// MergeWithStrategy merges layers according to the given strategy.
// Layers are processed left-to-right; later layers are considered "incoming".
func MergeWithStrategy(strategy MergeStrategy, layers ...map[string]string) (map[string]string, error) {
	result := make(map[string]string)

	for _, layer := range layers {
		if layer == nil {
			continue
		}
		for k, v := range layer {
			existing, exists := result[k]
			switch strategy {
			case StrategyKeepFirst:
				if !exists {
					result[k] = v
				}
			case StrategyErrorOnConflict:
				if exists && existing != v {
					return nil, &ConflictError{Key: k, Existing: existing, Incoming: v}
				}
				result[k] = v
			default: // StrategyOverride
				result[k] = v
			}
		}
	}

	return result, nil
}
