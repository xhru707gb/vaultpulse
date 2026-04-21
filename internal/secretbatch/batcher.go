// Package secretbatch provides batched evaluation of multiple secret paths,
// collecting results and errors in a single pass.
package secretbatch

import (
	"errors"
	"fmt"
	"sync"
)

// Result holds the outcome of evaluating a single secret path.
type Result struct {
	Path    string
	OK      bool
	Message string
	Err     error
}

// Evaluator is a function that evaluates a single secret path.
type Evaluator func(path string) (ok bool, message string, err error)

// Batcher runs an Evaluator over a set of paths concurrently and collects results.
type Batcher struct {
	mu          sync.Mutex
	concurrency int
	evaluator   Evaluator
}

// New creates a new Batcher with the given concurrency limit and evaluator.
// concurrency must be >= 1.
func New(concurrency int, evaluator Evaluator) (*Batcher, error) {
	if concurrency < 1 {
		return nil, errors.New("secretbatch: concurrency must be at least 1")
	}
	if evaluator == nil {
		return nil, errors.New("secretbatch: evaluator must not be nil")
	}
	return &Batcher{concurrency: concurrency, evaluator: evaluator}, nil
}

// Run evaluates all paths concurrently (bounded by concurrency) and returns
// one Result per path in the same order as the input slice.
func (b *Batcher) Run(paths []string) []Result {
	if len(paths) == 0 {
		return nil
	}

	results := make([]Result, len(paths))
	sem := make(chan struct{}, b.concurrency)
	var wg sync.WaitGroup

	for i, p := range paths {
		wg.Add(1)
		sem <- struct{}{}
		go func(idx int, path string) {
			defer wg.Done()
			defer func() { <-sem }()

			ok, msg, err := b.evaluator(path)
			results[idx] = Result{
				Path:    path,
				OK:      ok,
				Message: msg,
				Err:     err,
			}
		}(i, p)
	}

	wg.Wait()
	return results
}

// Failures returns only the results where OK is false or Err is non-nil.
func Failures(results []Result) []Result {
	var out []Result
	for _, r := range results {
		if !r.OK || r.Err != nil {
			out = append(out, r)
		}
	}
	return out
}

// Summary returns a human-readable summary string.
func Summary(results []Result) string {
	total := len(results)
	failed := len(Failures(results))
	return fmt.Sprintf("%d evaluated, %d failed, %d ok", total, failed, total-failed)
}
