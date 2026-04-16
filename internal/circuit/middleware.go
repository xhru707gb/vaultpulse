package circuit

import "context"

// Do executes fn through the circuit breaker.
// It records success or failure automatically.
func Do(ctx context.Context, br *Breaker, fn func(ctx context.Context) error) error {
	if err := br.Allow(); err != nil {
		return err
	}
	if err := fn(ctx); err != nil {
		br.RecordFailure()
		return err
	}
	br.RecordSuccess()
	return nil
}
