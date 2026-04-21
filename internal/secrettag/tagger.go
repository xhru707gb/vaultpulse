// Package secrettag provides tagging and retrieval of arbitrary string labels
// against secret paths, enabling grouping and filtering by user-defined tags.
package secrettag

import (
	"errors"
	"fmt"
	"sort"
	"sync"
)

// ErrEmptyPath is returned when an empty path is provided.
var ErrEmptyPath = errors.New("secrettag: path must not be empty")

// ErrEmptyTag is returned when an empty tag is provided.
var ErrEmptyTag = errors.New("secrettag: tag must not be empty")

// ErrNotFound is returned when no tags exist for the given path.
var ErrNotFound = errors.New("secrettag: path not found")

// Tagger stores tags associated with secret paths.
type Tagger struct {
	mu   sync.RWMutex
	data map[string]map[string]struct{} // path -> set of tags
}

// New returns an initialised Tagger.
func New() *Tagger {
	return &Tagger{data: make(map[string]map[string]struct{})}
}

// Add associates tag with path. Duplicates are silently ignored.
func (t *Tagger) Add(path, tag string) error {
	if path == "" {
		return ErrEmptyPath
	}
	if tag == "" {
		return ErrEmptyTag
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.data[path] == nil {
		t.data[path] = make(map[string]struct{})
	}
	t.data[path][tag] = struct{}{}
	return nil
}

// Remove removes tag from path. Returns ErrNotFound if the path has no tags.
func (t *Tagger) Remove(path, tag string) error {
	if path == "" {
		return ErrEmptyPath
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	set, ok := t.data[path]
	if !ok {
		return fmt.Errorf("%w: %s", ErrNotFound, path)
	}
	delete(set, tag)
	if len(set) == 0 {
		delete(t.data, path)
	}
	return nil
}

// Tags returns a sorted slice of tags for path.
func (t *Tagger) Tags(path string) ([]string, error) {
	if path == "" {
		return nil, ErrEmptyPath
	}
	t.mu.RLock()
	defer t.mu.RUnlock()
	set, ok := t.data[path]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrNotFound, path)
	}
	out := make([]string, 0, len(set))
	for tag := range set {
		out = append(out, tag)
	}
	sort.Strings(out)
	return out, nil
}

// PathsWithTag returns all paths that carry the given tag, sorted.
func (t *Tagger) PathsWithTag(tag string) []string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	var out []string
	for path, set := range t.data {
		if _, ok := set[tag]; ok {
			out = append(out, path)
		}
	}
	sort.Strings(out)
	return out
}

// Len returns the number of paths that have at least one tag.
func (t *Tagger) Len() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return len(t.data)
}
