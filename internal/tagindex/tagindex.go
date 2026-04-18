// Package tagindex provides a tag-based index for grouping and querying
// Vault secret paths by user-defined labels.
package tagindex

import (
	"errors"
	"sort"
	"sync"
)

// ErrEmptyTag is returned when an empty tag is provided.
var ErrEmptyTag = errors.New("tagindex: tag must not be empty")

// ErrEmptyPath is returned when an empty path is provided.
var ErrEmptyPath = errors.New("tagindex: path must not be empty")

// Index maps tags to sets of secret paths.
type Index struct {
	mu      sync.RWMutex
	entries map[string]map[string]struct{}
}

// New returns an initialised Index.
func New() *Index {
	return &Index{entries: make(map[string]map[string]struct{})}
}

// Add associates path with tag.
func (idx *Index) Add(tag, path string) error {
	if tag == "" {
		return ErrEmptyTag
	}
	if path == "" {
		return ErrEmptyPath
	}
	idx.mu.Lock()
	defer idx.mu.Unlock()
	if idx.entries[tag] == nil {
		idx.entries[tag] = make(map[string]struct{})
	}
	idx.entries[tag][path] = struct{}{}
	return nil
}

// Remove disassociates path from tag.
func (idx *Index) Remove(tag, path string) error {
	if tag == "" {
		return ErrEmptyTag
	}
	if path == "" {
		return ErrEmptyPath
	}
	idx.mu.Lock()
	defer idx.mu.Unlock()
	delete(idx.entries[tag], path)
	if len(idx.entries[tag]) == 0 {
		delete(idx.entries, tag)
	}
	return nil
}

// Paths returns all paths associated with tag, sorted.
func (idx *Index) Paths(tag string) []string {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	set := idx.entries[tag]
	out := make([]string, 0, len(set))
	for p := range set {
		out = append(out, p)
	}
	sort.Strings(out)
	return out
}

// Tags returns all known tags, sorted.
func (idx *Index) Tags() []string {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	out := make([]string, 0, len(idx.entries))
	for t := range idx.entries {
		out = append(out, t)
	}
	sort.Strings(out)
	return out
}

// Len returns the number of paths associated with tag.
func (idx *Index) Len(tag string) int {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	return len(idx.entries[tag])
}
