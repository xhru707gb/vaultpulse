// Package schedule provides periodic job scheduling for VaultPulse checks.
package schedule

import (
	"context"
	"errors"
	"time"
)

// ErrInvalidInterval is returned when the interval is zero or negative.
var ErrInvalidInterval = errors.New("schedule: interval must be positive")

// Job is a function executed on each tick.
type Job func(ctx context.Context) error

// Scheduler runs a Job at a fixed interval.
type Scheduler struct {
	interval time.Duration
	job      Job
	now      func() time.Time
}

// New creates a Scheduler with the given interval and job.
func New(interval time.Duration, job Job) (*Scheduler, error) {
	if interval <= 0 {
		return nil, ErrInvalidInterval
	}
	if job == nil {
		return nil, errors.New("schedule: job must not be nil")
	}
	return &Scheduler{interval: interval, job: job, now: time.Now}, nil
}

// Run blocks and executes the job on every tick until ctx is cancelled.
// The first execution happens after the first interval elapses.
func (s *Scheduler) Run(ctx context.Context) error {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := s.job(ctx); err != nil {
				return err
			}
		}
	}
}
