// Package schedule provides a lightweight periodic job scheduler for
// VaultPulse. It supports running a single Job at a fixed interval via
// Scheduler, and orchestrating multiple named schedulers concurrently
// via Runner.
//
// Usage:
//
//	s, err := schedule.New(30*time.Second, myJob)
//	if err != nil { ... }
//	s.Run(ctx)
package schedule
