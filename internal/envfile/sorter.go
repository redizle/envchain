package envfile

import (
	"sort"
	"strings"
)

// SortOrder defines how env keys should be sorted.
type SortOrder int

const (
	// SortAlpha sorts keys alphabetically (A-Z).
	SortAlpha SortOrder = iota
	// SortAlphaDesc sorts keys reverse alphabetically (Z-A).
	SortAlphaDesc
	// SortByGroup sorts keys grouped by common prefix (e.g. DB_, AWS_).
	SortByGroup
)

// SortOptions controls sorting behaviour.
type SortOptions struct {
	Order      SortOrder
	IgnoreCase bool
}

// Sort returns a new map with keys ordered according to opts and a slice of
// keys in that order. The original map is not mutated.
func Sort(env map[string]string, opts SortOptions) (map[string]string, []string) {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}

	cmp := func(a, b string) bool {
		if opts.IgnoreCase {
			return strings.ToLower(a) < strings.ToLower(b)
		}
		return a < b
	}

	switch opts.Order {
	case SortAlphaDesc:
		sort.Slice(keys, func(i, j int) bool { return !cmp(keys[i], keys[j]) })
	case SortByGroup:
		sort.Slice(keys, func(i, j int) bool {
			gi := groupPrefix(keys[i])
			gj := groupPrefix(keys[j])
			if gi != gj {
				if opts.IgnoreCase {
					return strings.ToLower(gi) < strings.ToLower(gj)
				}
				return gi < gj
			}
			return cmp(keys[i], keys[j])
		})
	default: // SortAlpha
		sort.Slice(keys, func(i, j int) bool { return cmp(keys[i], keys[j]) })
	}

	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}
	return out, keys
}

// groupPrefix returns the prefix before the first underscore, or the whole key
// if no underscore is present.
func groupPrefix(key string) string {
	if idx := strings.Index(key, "_"); idx > 0 {
		return key[:idx]
	}
	return key
}
