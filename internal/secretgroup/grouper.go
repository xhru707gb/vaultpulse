// Package secretgroup provides grouping and aggregation of secrets by label,
// prefix, or owner for bulk reporting and alerting.
package secretgroup

import (
	"errors"
	"sort"
	"strings"
	"sync"
)

// Group holds a named collection of secret paths.
type Group struct {
	Name  string
	Paths []string
}

// Grouper manages secret path groupings.
type Grouper struct {
	mu     sync.RWMutex
	groups map[string]*Group
}

// New returns an initialised Grouper.
func New() *Grouper {
	return &Grouper{groups: make(map[string]*Group)}
}

// Add registers path under the named group, creating it if absent.
func (g *Grouper) Add(name, path string) error {
	if name == "" {
		return errors.New("group name must not be empty")
	}
	if path == "" {
		return errors.New("path must not be empty")
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	grp, ok := g.groups[name]
	if !ok {
		grp = &Group{Name: name}
		g.groups[name] = grp
	}
	for _, p := range grp.Paths {
		if p == path {
			return nil
		}
	}
	grp.Paths = append(grp.Paths, path)
	return nil
}

// Remove deletes path from the named group.
func (g *Grouper) Remove(name, path string) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	grp, ok := g.groups[name]
	if !ok {
		return errors.New("group not found: " + name)
	}
	for i, p := range grp.Paths {
		if p == path {
			grp.Paths = append(grp.Paths[:i], grp.Paths[i+1:]...)
			return nil
		}
	}
	return errors.New("path not found in group")
}

// Get returns the Group for name, or false if absent.
func (g *Grouper) Get(name string) (*Group, bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	grp, ok := g.groups[name]
	if !ok {
		return nil, false
	}
	copy := &Group{Name: grp.Name, Paths: append([]string(nil), grp.Paths...)}
	return copy, true
}

// All returns all groups sorted by name.
func (g *Grouper) All() []*Group {
	g.mu.RLock()
	defer g.mu.RUnlock()
	out := make([]*Group, 0, len(g.groups))
	for _, grp := range g.groups {
		out = append(out, &Group{Name: grp.Name, Paths: append([]string(nil), grp.Paths...)})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

// FindByPrefix returns all groups that contain at least one path with the given prefix.
func (g *Grouper) FindByPrefix(prefix string) []*Group {
	g.mu.RLock()
	defer g.mu.RUnlock()
	var out []*Group
	for _, grp := range g.groups {
		for _, p := range grp.Paths {
			if strings.HasPrefix(p, prefix) {
				out = append(out, &Group{Name: grp.Name, Paths: append([]string(nil), grp.Paths...)})
				break
			}
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}
