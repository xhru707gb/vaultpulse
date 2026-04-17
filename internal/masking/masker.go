// Package masking provides utilities for masking sensitive secret values
// before they are displayed or transmitted.
package masking

import "strings"

// Level controls how aggressively values are masked.
type Level int

const (
	// LevelFull replaces the entire value with asterisks.
	LevelFull Level = iota
	// LevelPartial reveals the first and last two characters.
	LevelPartial
	// LevelNone performs no masking.
	LevelNone
)

// Masker masks secret values according to a configured level.
type Masker struct {
	level Level
	keys  []string // key substrings that trigger masking
}

// New returns a Masker with the given level and sensitive key patterns.
func New(level Level, sensitiveKeys []string) *Masker {
	return &Masker{level: level, keys: sensitiveKeys}
}

// ShouldMask reports whether the given key matches a sensitive pattern.
func (m *Masker) ShouldMask(key string) bool {
	lower := strings.ToLower(key)
	for _, k := range m.keys {
		if strings.Contains(lower, strings.ToLower(k)) {
			return true
		}
	}
	return false
}

// Mask returns the masked representation of value for the given key.
func (m *Masker) Mask(key, value string) string {
	if !m.ShouldMask(key) {
		return value
	}
	switch m.level {
	case LevelNone:
		return value
	case LevelPartial:
		return partial(value)
	default:
		return strings.Repeat("*", 8)
	}
}

// MaskMap masks all sensitive keys in a map, returning a new map.
func (m *Masker) MaskMap(data map[string]string) map[string]string {
	out := make(map[string]string, len(data))
	for k, v := range data {
		out[k] = m.Mask(k, v)
	}
	return out
}

func partial(v string) string {
	if len(v) <= 4 {
		return strings.Repeat("*", len(v))
	}
	return string(v[:2]) + strings.Repeat("*", len(v)-4) + string(v[len(v)-2:])
}
