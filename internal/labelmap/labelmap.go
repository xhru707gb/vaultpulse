// Package labelmap provides key-value label management for secrets,
// enabling grouping, filtering, and annotation of vault paths.
package labelmap

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

// ErrEmptyKey is returned when an empty label key is provided.
var ErrEmptyKey = errors.New("labelmap: key must not be empty")

// Map holds labels associated with vault secret paths.
type Map struct {
	entries map[string]map[string]string // path -> labels
}

// New creates an empty Map.
func New() *Map {
	return &Map{entries: make(map[string]map[string]string)}
}

// Set assigns a label key=value to the given path.
func (m *Map) Set(path, key, value string) error {
	if strings.TrimSpace(key) == "" {
		return ErrEmptyKey
	}
	if _, ok := m.entries[path]; !ok {
		m.entries[path] = make(map[string]string)
	}
	m.entries[path][key] = value
	return nil
}

// Get returns all labels for a path, or nil if none exist.
func (m *Map) Get(path string) map[string]string {
	labels, ok := m.entries[path]
	if !ok {
		return nil
	}
	out := make(map[string]string, len(labels))
	for k, v := range labels {
		out[k] = v
	}
	return out
}

// Delete removes a specific label key from a path.
func (m *Map) Delete(path, key string) {
	if labels, ok := m.entries[path]; ok {
		delete(labels, key)
		if len(labels) == 0 {
			delete(m.entries, path)
		}
	}
}

// Filter returns paths whose labels match all provided key=value pairs.
func (m *Map) Filter(selector map[string]string) []string {
	var matched []string
	for path, labels := range m.entries {
		if matchesAll(labels, selector) {
			matched = append(matched, path)
		}
	}
	sort.Strings(matched)
	return matched
}

// FormatLabels returns a compact string representation of a path's labels.
func FormatLabels(labels map[string]string) string {
	if len(labels) == 0 {
		return "<none>"
	}
	keys := make([]string, 0, len(labels))
	for k := range labels {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%s", k, labels[k]))
	}
	return strings.Join(parts, ", ")
}

func matchesAll(labels, selector map[string]string) bool {
	for k, v := range selector {
		if labels[k] != v {
			return false
		}
	}
	return true
}
