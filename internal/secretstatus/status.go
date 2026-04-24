// Package secretstatus provides aggregated health status evaluation
// across multiple secret dimensions (expiry, rotation, policy, score).
package secretstatus

import (
	"errors"
	"sync"
	"time"
)

// Level represents the overall status level of a secret.
type Level int

const (
	LevelOK Level = iota
	LevelWarning
	LevelCritical
)

// Entry holds the aggregated status for a single secret path.
type Entry struct {
	Path        string
	Level       Level
	Reasons     []string
	EvaluatedAt time.Time
}

// Evaluator aggregates status from registered dimension providers.
type Evaluator struct {
	mu        sync.RWMutex
	providers []Provider
}

// Provider is implemented by any dimension that can contribute a status level.
type Provider interface {
	Name() string
	Evaluate(path string) (Level, string, error)
}

// New returns a new Evaluator with the given providers.
func New(providers ...Provider) (*Evaluator, error) {
	if len(providers) == 0 {
		return nil, errors.New("secretstatus: at least one provider is required")
	}
	return &Evaluator{providers: providers}, nil
}

// Evaluate returns an aggregated Entry for the given secret path.
func (e *Evaluator) Evaluate(path string) (*Entry, error) {
	if path == "" {
		return nil, errors.New("secretstatus: path must not be empty")
	}
	e.mu.RLock()
	defer e.mu.RUnlock()

	entry := &Entry{
		Path:        path,
		Level:       LevelOK,
		EvaluatedAt: time.Now().UTC(),
	}
	for _, p := range e.providers {
		lvl, reason, err := p.Evaluate(path)
		if err != nil {
			continue
		}
		if reason != "" {
			entry.Reasons = append(entry.Reasons, p.Name()+": "+reason)
		}
		if lvl > entry.Level {
			entry.Level = lvl
		}
	}
	return entry, nil
}

// EvaluateAll evaluates a list of paths and returns all entries.
func (e *Evaluator) EvaluateAll(paths []string) ([]*Entry, error) {
	if len(paths) == 0 {
		return nil, errors.New("secretstatus: paths must not be empty")
	}
	entries := make([]*Entry, 0, len(paths))
	for _, p := range paths {
		en, err := e.Evaluate(p)
		if err != nil {
			continue
		}
		entries = append(entries, en)
	}
	return entries, nil
}
