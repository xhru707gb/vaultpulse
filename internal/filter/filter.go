// Package filter provides path-based filtering for secret status results.
package filter

import (
	"strings"
)

// Options holds filtering criteria.
type Options struct {
	Prefix  string
	Contains string
	Exclude  string
}

// Filter applies Options to a slice of paths, returning only matching ones.
func Filter(paths []string, opts Options) []string {
	out := make([]string, 0, len(paths))
	for _, p := range paths {
		if opts.Prefix != "" && !strings.HasPrefix(p, opts.Prefix) {
			continue
		}
		if opts.Contains != "" && !strings.Contains(p, opts.Contains) {
			continue
		}
		if opts.Exclude != "" && strings.Contains(p, opts.Exclude) {
			continue
		}
		out = append(out, p)
	}
	return out
}

// MatchesAny returns true if the path satisfies any of the given prefixes.
func MatchesAny(path string, prefixes []string) bool {
	for _, pfx := range prefixes {
		if strings.HasPrefix(path, pfx) {
			return true
		}
	}
	return false
}
