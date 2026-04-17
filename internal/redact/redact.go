// Package redact provides utilities for masking sensitive secret values
// in logs, output, and alert payloads.
package redact

import "strings"

const defaultMask = "[REDACTED]"

// Redactor masks sensitive values based on configured key patterns.
type Redactor struct {
	patterns []string
	mask     string
}

// New returns a Redactor that masks values whose keys contain any of the
// given patterns (case-insensitive).
func New(patterns []string) *Redactor {
	return &Redactor{patterns: patterns, mask: defaultMask}
}

// ShouldRedact reports whether the given key matches a sensitive pattern.
func (r *Redactor) ShouldRedact(key string) bool {
	lower := strings.ToLower(key)
	for _, p := range r.patterns {
		if strings.Contains(lower, strings.ToLower(p)) {
			return true
		}
	}
	return false
}

// Value returns the masked string if the key is sensitive, otherwise the
// original value.
func (r *Redactor) Value(key, value string) string {
	if r.ShouldRedact(key) {
		return r.mask
	}
	return value
}

// Map returns a copy of m with sensitive values replaced by the mask.
func (r *Redactor) Map(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = r.Value(k, v)
	}
	return out
}
