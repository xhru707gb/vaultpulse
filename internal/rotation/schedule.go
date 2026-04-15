// Package rotation provides functionality for tracking and evaluating
// secret rotation schedules defined in the vaultpulse configuration.
package rotation

import (
	"fmt"
	"time"
)

// Schedule represents a rotation policy for a single secret path.
type Schedule struct {
	Path     string
	Interval time.Duration
	LastRotated time.Time
}

// Status holds the evaluated rotation state for a secret.
type Status struct {
	Path        string
	DueIn       time.Duration
	Overdue     bool
	LastRotated time.Time
	NextDue     time.Time
}

// Evaluator checks rotation schedules against current time.
type Evaluator struct {
	now func() time.Time
}

// NewEvaluator creates an Evaluator. Pass nil to use real time.
func NewEvaluator(now func() time.Time) *Evaluator {
	if now == nil {
		now = time.Now
	}
	return &Evaluator{now: now}
}

// Evaluate returns a Status for the given Schedule.
func (e *Evaluator) Evaluate(s Schedule) (Status, error) {
	if s.Interval <= 0 {
		return Status{}, fmt.Errorf("rotation: interval must be positive for path %q", s.Path)
	}
	now := e.now()
	nextDue := s.LastRotated.Add(s.Interval)
	dueIn := nextDue.Sub(now)
	return Status{
		Path:        s.Path,
		DueIn:       dueIn,
		Overdue:     dueIn < 0,
		LastRotated: s.LastRotated,
		NextDue:     nextDue,
	}, nil
}

// EvaluateAll evaluates a slice of schedules, returning all statuses.
// Non-fatal errors per schedule are collected and returned together.
func (e *Evaluator) EvaluateAll(schedules []Schedule) ([]Status, []error) {
	var statuses []Status
	var errs []error
	for _, s := range schedules {
		st, err := e.Evaluate(s)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		statuses = append(statuses, st)
	}
	return statuses, errs
}
