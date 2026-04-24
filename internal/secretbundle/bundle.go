// Package secretbundle groups related secrets into named bundles
// and provides aggregate evaluation across all members.
package secretbundle

import (
	"errors"
	"fmt"
	"sync"
)

// Entry represents a single secret within a bundle.
type Entry struct {
	Path    string
	Version int
	Expired bool
}

// Bundle holds a named collection of secret paths.
type Bundle struct {
	Name    string
	Entries []Entry
}

// EvalResult summarises the health of a bundle.
type EvalResult struct {
	Name       string
	Total      int
	Expired    int
	Healthy    bool
}

// Registry stores named bundles.
type Registry struct {
	mu      sync.RWMutex
	bundles map[string]*Bundle
}

// New creates an empty Registry.
func New() *Registry {
	return &Registry{bundles: make(map[string]*Bundle)}
}

// Add registers a new bundle. Returns an error if the name already exists or is empty.
func (r *Registry) Add(name string) error {
	if name == "" {
		return errors.New("bundle name must not be empty")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.bundles[name]; ok {
		return fmt.Errorf("bundle %q already exists", name)
	}
	r.bundles[name] = &Bundle{Name: name}
	return nil
}

// AddEntry appends a secret entry to an existing bundle.
func (r *Registry) AddEntry(bundleName, path string, version int, expired bool) error {
	if path == "" {
		return errors.New("path must not be empty")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	b, ok := r.bundles[bundleName]
	if !ok {
		return fmt.Errorf("bundle %q not found", bundleName)
	}
	b.Entries = append(b.Entries, Entry{Path: path, Version: version, Expired: expired})
	return nil
}

// Evaluate returns an EvalResult for the named bundle.
func (r *Registry) Evaluate(name string) (EvalResult, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	b, ok := r.bundles[name]
	if !ok {
		return EvalResult{}, fmt.Errorf("bundle %q not found", name)
	}
	expired := 0
	for _, e := range b.Entries {
		if e.Expired {
			expired++
		}
	}
	return EvalResult{
		Name:    name,
		Total:   len(b.Entries),
		Expired: expired,
		Healthy: expired == 0,
	}, nil
}

// EvaluateAll returns results for every registered bundle.
func (r *Registry) EvaluateAll() []EvalResult {
	r.mu.RLock()
	defer r.mu.RUnlock()
	results := make([]EvalResult, 0, len(r.bundles))
	for name, b := range r.bundles {
		expired := 0
		for _, e := range b.Entries {
			if e.Expired {
				expired++
			}
		}
		results = append(results, EvalResult{
			Name:    name,
			Total:   len(b.Entries),
			Expired: expired,
			Healthy: expired == 0,
		})
	}
	return results
}
