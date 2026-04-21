// Package secretlookup provides a reverse-lookup index from secret value
// fingerprints to the paths that reference them, enabling detection of
// shared or duplicated secrets across multiple paths.
package secretlookup

import (
	"errors"
	"sync"
)

// ErrEmptyPath is returned when an empty path is provided.
var ErrEmptyPath = errors.New("secretlookup: path must not be empty")

// ErrEmptyFingerprint is returned when an empty fingerprint is provided.
var ErrEmptyFingerprint = errors.New("secretlookup: fingerprint must not be empty")

// Index maps fingerprints to the set of paths that share that fingerprint.
type Index struct {
	mu      sync.RWMutex
	entries map[string]map[string]struct{} // fingerprint -> set of paths
}

// New creates a new empty Index.
func New() *Index {
	return &Index{
		entries: make(map[string]map[string]struct{}),
	}
}

// Add records that the given path has the given fingerprint.
func (idx *Index) Add(path, fingerprint string) error {
	if path == "" {
		return ErrEmptyPath
	}
	if fingerprint == "" {
		return ErrEmptyFingerprint
	}
	idx.mu.Lock()
	defer idx.mu.Unlock()
	if _, ok := idx.entries[fingerprint]; !ok {
		idx.entries[fingerprint] = make(map[string]struct{})
	}
	idx.entries[fingerprint][path] = struct{}{}
	return nil
}

// Lookup returns all paths that share the given fingerprint.
// Returns nil if the fingerprint is not known.
func (idx *Index) Lookup(fingerprint string) []string {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	paths, ok := idx.entries[fingerprint]
	if !ok {
		return nil
	}
	out := make([]string, 0, len(paths))
	for p := range paths {
		out = append(out, p)
	}
	return out
}

// Remove removes the given path from the index entry for the fingerprint.
// If no paths remain for that fingerprint, the fingerprint is deleted.
func (idx *Index) Remove(path, fingerprint string) error {
	if path == "" {
		return ErrEmptyPath
	}
	if fingerprint == "" {
		return ErrEmptyFingerprint
	}
	idx.mu.Lock()
	defer idx.mu.Unlock()
	paths, ok := idx.entries[fingerprint]
	if !ok {
		return nil
	}
	delete(paths, path)
	if len(paths) == 0 {
		delete(idx.entries, fingerprint)
	}
	return nil
}

// Duplicates returns all fingerprints that are shared by more than one path,
// together with the list of paths sharing them.
func (idx *Index) Duplicates() map[string][]string {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	out := make(map[string][]string)
	for fp, paths := range idx.entries {
		if len(paths) > 1 {
			list := make([]string, 0, len(paths))
			for p := range paths {
				list = append(list, p)
			}
			out[fp] = list
		}
	}
	return out
}
