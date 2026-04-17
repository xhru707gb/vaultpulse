package schedule

import (
	"context"
	"fmt"
	"sync"
)

// Runner manages multiple named schedulers.
type Runner struct {
	mu         sync.Mutex
	schedulers map[string]*Scheduler
}

// NewRunner creates an empty Runner.
func NewRunner() *Runner {
	return &Runner{schedulers: make(map[string]*Scheduler)}
}

// Register adds a named scheduler to the runner.
func (r *Runner) Register(name string, s *Scheduler) error {
	if name == "" {
		return fmt.Errorf("schedule: name must not be empty")
	}
	if s == nil {
		return fmt.Errorf("schedule: scheduler must not be nil")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.schedulers[name]; exists {
		return fmt.Errorf("schedule: %q already registered", name)
	}
	r.schedulers[name] = s
	return nil
}

// RunAll starts all registered schedulers concurrently.
// It returns the first error encountered, cancelling the remaining.
func (r *Runner) RunAll(ctx context.Context) error {
	r.mu.Lock()
	names := make([]string, 0, len(r.schedulers))
	for n := range r.schedulers {
		names = append(names, n)
	}
	r.mu.Unlock()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errCh := make(chan error, len(names))
	for _, n := range names {
		s := r.schedulers[n]
		go func() {
			errCh <- s.Run(ctx)
		}()
	}
	return <-errCh
}
