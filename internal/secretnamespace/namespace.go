// Package secretnamespace provides namespace-based grouping and isolation
// for secrets within a Vault instance.
package secretnamespace

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
)

// ErrEmptyNamespace is returned when an empty namespace is provided.
var ErrEmptyNamespace = errors.New("namespace must not be empty")

// ErrEmptyPath is returned when an empty secret path is provided.
var ErrEmptyPath = errors.New("secret path must not be empty")

// ErrDuplicatePath is returned when a path already exists in the namespace.
var ErrDuplicatePath = errors.New("path already registered in namespace")

// ErrUnknownNamespace is returned when the namespace does not exist.
var ErrUnknownNamespace = errors.New("unknown namespace")

// Registry manages secret paths organised by namespace.
type Registry struct {
	mu   sync.RWMutex
	data map[string][]string // namespace -> paths
}

// New creates a new namespace Registry.
func New() *Registry {
	return &Registry{data: make(map[string][]string)}
}

// Add registers a secret path under the given namespace.
func (r *Registry) Add(namespace, path string) error {
	if strings.TrimSpace(namespace) == "" {
		return ErrEmptyNamespace
	}
	if strings.TrimSpace(path) == "" {
		return ErrEmptyPath
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, p := range r.data[namespace] {
		if p == path {
			return fmt.Errorf("%w: %s", ErrDuplicatePath, path)
		}
	}
	r.data[namespace] = append(r.data[namespace], path)
	return nil
}

// Paths returns all paths registered under the given namespace.
func (r *Registry) Paths(namespace string) ([]string, error) {
	if strings.TrimSpace(namespace) == "" {
		return nil, ErrEmptyNamespace
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	paths, ok := r.data[namespace]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrUnknownNamespace, namespace)
	}
	out := make([]string, len(paths))
	copy(out, paths)
	return out, nil
}

// Namespaces returns all registered namespace names in sorted order.
func (r *Registry) Namespaces() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.data))
	for ns := range r.data {
		names = append(names, ns)
	}
	sort.Strings(names)
	return names
}

// Remove deletes a path from the given namespace.
func (r *Registry) Remove(namespace, path string) error {
	if strings.TrimSpace(namespace) == "" {
		return ErrEmptyNamespace
	}
	if strings.TrimSpace(path) == "" {
		return ErrEmptyPath
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	paths, ok := r.data[namespace]
	if !ok {
		return fmt.Errorf("%w: %s", ErrUnknownNamespace, namespace)
	}
	for i, p := range paths {
		if p == path {
			r.data[namespace] = append(paths[:i], paths[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("path not found in namespace %s: %s", namespace, path)
}
