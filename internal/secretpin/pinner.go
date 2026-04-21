// Package secretpin provides version pinning for secrets, allowing operators
// to lock a secret path to a specific version and detect drift from that pin.
package secretpin

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// ErrAlreadyPinned is returned when a path is already pinned.
var ErrAlreadyPinned = errors.New("path already pinned")

// ErrNotPinned is returned when a path has no active pin.
var ErrNotPinned = errors.New("path not pinned")

// Pin holds the pinned version metadata for a secret path.
type Pin struct {
	Path      string
	Version   int
	PinnedAt  time.Time
	PinnedBy  string
}

// DriftResult describes whether the current version matches the pin.
type DriftResult struct {
	Path           string
	PinnedVersion  int
	CurrentVersion int
	Drifted        bool
}

// Pinner manages version pins for secret paths.
type Pinner struct {
	mu   sync.RWMutex
	pins map[string]Pin
	now  func() time.Time
}

// New creates a new Pinner.
func New() *Pinner {
	return &Pinner{
		pins: make(map[string]Pin),
		now:  time.Now,
	}
}

// Pin registers a version pin for the given path.
func (p *Pinner) Pin(path string, version int, pinnedBy string) error {
	if path == "" {
		return errors.New("path must not be empty")
	}
	if version < 1 {
		return fmt.Errorf("version must be >= 1, got %d", version)
	}
	if pinnedBy == "" {
		return errors.New("pinnedBy must not be empty")
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, exists := p.pins[path]; exists {
		return fmt.Errorf("%w: %s", ErrAlreadyPinned, path)
	}
	p.pins[path] = Pin{
		Path:     path,
		Version:  version,
		PinnedAt: p.now(),
		PinnedBy: pinnedBy,
	}
	return nil
}

// Unpin removes the pin for the given path.
func (p *Pinner) Unpin(path string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, exists := p.pins[path]; !exists {
		return fmt.Errorf("%w: %s", ErrNotPinned, path)
	}
	delete(p.pins, path)
	return nil
}

// Check compares the current version against the pinned version.
func (p *Pinner) Check(path string, currentVersion int) (DriftResult, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	pin, exists := p.pins[path]
	if !exists {
		return DriftResult{}, fmt.Errorf("%w: %s", ErrNotPinned, path)
	}
	return DriftResult{
		Path:           path,
		PinnedVersion:  pin.Version,
		CurrentVersion: currentVersion,
		Drifted:        currentVersion != pin.Version,
	}, nil
}

// All returns all active pins.
func (p *Pinner) All() []Pin {
	p.mu.RLock()
	defer p.mu.RUnlock()
	out := make([]Pin, 0, len(p.pins))
	for _, pin := range p.pins {
		out = append(out, pin)
	}
	return out
}
