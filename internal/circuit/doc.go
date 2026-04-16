// Package circuit provides a circuit breaker implementation for vaultpulse.
//
// The circuit breaker transitions between three states:
//   - Closed: normal operation, requests pass through.
//   - Open: failures exceeded threshold; requests are rejected with ErrOpen.
//   - Half-Open: after OpenTimeout elapses, one probe request is allowed.
//
// Usage:
//
//	br, _ := circuit.New(circuit.Config{MaxFailures: 3, OpenTimeout: 10 * time.Second})
//	if err := br.Allow(); err != nil {
//	    // circuit is open, skip call
//	}
package circuit
