// Package secretmeta provides a registry for attaching arbitrary metadata
// key-value pairs to secret paths, enabling richer querying and reporting.
package secretmeta

import (
	"errors"
	"fmt"
	"sync"
)

// ErrEmptyPath is returned when an empty path is provided.
var ErrEmptyPath = errors.New("secretmeta: path must not be empty")

// ErrEmptyKey is returned when an empty metadata key is provided.
var ErrEmptyKey = errors.New("secretmeta: metadata key must not be empty")

// ErrNotFound is returned when a path has no registered metadata.
var ErrNotFound = errors.New("secretmeta: path not found")

// Registry stores metadata for secret paths.
type Registry struct {
	mu   sync.RWMutex
	data map[string]map[string]string
}

// New creates an empty Registry.
func New() *Registry {
	return &Registry{data: make(map[string]map[string]string)}
}

// Set attaches a key-value pair to the given path.
func (r *Registry) Set(path, key, value string) error {
	if path == "" {
		return ErrEmptyPath
	}
	if key == "" {
		return ErrEmptyKey
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.data[path] == nil {
		r.data[path] = make(map[string]string)
	}
	r.data[path][key] = value
	return nil
}

// Get returns the metadata map for a path, or ErrNotFound.
func (r *Registry) Get(path string) (map[string]string, error) {
	if path == "" {
		return nil, ErrEmptyPath
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	m, ok := r.data[path]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrNotFound, path)
	}
	copy := make(map[string]string, len(m))
	for k, v := range m {
		copy[k] = v
	}
	return copy, nil
}

// Delete removes all metadata for a path.
func (r *Registry) Delete(path string) error {
	if path == "" {
		return ErrEmptyPath
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[path]; !ok {
		return fmt.Errorf("%w: %s", ErrNotFound, path)
	}
	delete(r.data, path)
	return nil
}

// Paths returns all registered paths.
func (r *Registry) Paths() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	paths := make([]string, 0, len(r.data))
	for p := range r.data {
		paths = append(paths, p)
	}
	return paths
}
