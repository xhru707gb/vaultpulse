// Package secrethash provides utilities for hashing secret values
// and detecting changes across evaluations.
package secrethash

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"
)

// ErrEmptyPath is returned when an empty path is provided.
var ErrEmptyPath = errors.New("secrethash: path must not be empty")

// ErrEmptyValue is returned when an empty value is provided.
var ErrEmptyValue = errors.New("secrethash: value must not be empty")

// Entry holds the hash and metadata for a single secret.
type Entry struct {
	Path      string
	Hash      string
	ChangedAt time.Time
	Version   int
}

// Hasher tracks hashed secret values and detects changes.
type Hasher struct {
	mu      sync.RWMutex
	entries map[string]Entry
	now     func() time.Time
}

// New creates a new Hasher.
func New() *Hasher {
	return &Hasher{
		entries: make(map[string]Entry),
		now:     time.Now,
	}
}

// Record hashes the given value for path and stores it.
// Returns true if the hash changed (or is new), false if unchanged.
func (h *Hasher) Record(path, value string) (bool, error) {
	if path == "" {
		return false, ErrEmptyPath
	}
	if value == "" {
		return false, ErrEmptyValue
	}

	hash := hashValue(value)

	h.mu.Lock()
	defer h.mu.Unlock()

	existing, ok := h.entries[path]
	if ok && existing.Hash == hash {
		return false, nil
	}

	version := 1
	if ok {
		version = existing.Version + 1
	}

	h.entries[path] = Entry{
		Path:      path,
		Hash:      hash,
		ChangedAt: h.now(),
		Version:   version,
	}
	return true, nil
}

// Get returns the Entry for a given path.
func (h *Hasher) Get(path string) (Entry, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	e, ok := h.entries[path]
	return e, ok
}

// All returns a snapshot of all tracked entries.
func (h *Hasher) All() []Entry {
	h.mu.RLock()
	defer h.mu.RUnlock()
	out := make([]Entry, 0, len(h.entries))
	for _, e := range h.entries {
		out = append(out, e)
	}
	return out
}

func hashValue(v string) string {
	sum := sha256.Sum256([]byte(v))
	return fmt.Sprintf("%s", hex.EncodeToString(sum[:]))
}
