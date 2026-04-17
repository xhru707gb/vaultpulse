// Package sampling provides a thread-safe probabilistic sampler used to
// reduce the volume of audit events forwarded to downstream sinks.
//
// Usage:
//
//	s, err := sampling.New(sampling.Config{Rate: 0.25})
//	if err != nil { ... }
//	if s.Sample() {
//		// forward the event
//	}
package sampling
