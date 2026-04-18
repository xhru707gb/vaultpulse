// Package escalation provides multi-level alert escalation based on severity thresholds.
package escalation

import (
	"errors"
	"time"
)

// Level represents an escalation severity level.
type Level string

const (
	LevelInfo     Level = "info"
	LevelWarning  Level = "warning"
	LevelCritical Level = "critical"
)

// Rule defines when escalation to a given level should occur.
type Rule struct {
	Level     Level
	Threshold time.Duration // TTL below which this rule triggers
}

// Event is the result of evaluating a secret against escalation rules.
type Event struct {
	Path      string
	Level     Level
	TTL       time.Duration
	Evaluated time.Time
}

// Escalator evaluates secrets against escalation rules.
type Escalator struct {
	rules []Rule
	now   func() time.Time
}

// ErrNoRules is returned when no rules are configured.
var ErrNoRules = errors.New("escalation: no rules configured")

// New creates an Escalator with the given rules.
// Rules should be ordered from most to least severe.
func New(rules []Rule, now func() time.Time) (*Escalator, error) {
	if len(rules) == 0 {
		return nil, ErrNoRules
	}
	if now == nil {
		now = time.Now
	}
	return &Escalator{rules: rules, now: now}, nil
}

// Evaluate returns an Event for the given path and TTL.
// Returns nil if no rule matches.
func (e *Escalator) Evaluate(path string, ttl time.Duration) *Event {
	for _, r := range e.rules {
		if ttl <= r.Threshold {
			return &Event{
				Path:      path,
				Level:     r.Level,
				TTL:       ttl,
				Evaluated: e.now(),
			}
		}
	}
	return nil
}

// EvaluateAll evaluates a map of path→TTL and returns all matching events.
func (e *Escalator) EvaluateAll(secrets map[string]time.Duration) []Event {
	var events []Event
	for path, ttl := range secrets {
		if ev := e.Evaluate(path, ttl); ev != nil {
			events = append(events, *ev)
		}
	}
	return events
}
