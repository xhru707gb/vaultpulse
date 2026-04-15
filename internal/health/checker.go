// Package health provides Vault cluster health-check utilities.
package health

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/vault/api"
)

// Status holds the result of a single Vault health probe.
type Status struct {
	Initialized bool
	Sealed      bool
	Standby     bool
	CheckedAt   time.Time
	Latency     time.Duration
	Error       error
}

// Checker probes Vault cluster health.
type Checker struct {
	client *api.Client
	now    func() time.Time
}

// NewChecker constructs a Checker from an existing Vault API client.
func NewChecker(client *api.Client) *Checker {
	return &Checker{client: client, now: time.Now}
}

// Check performs a health probe against Vault and returns a Status.
func (c *Checker) Check(ctx context.Context) Status {
	start := c.now()
	s := Status{CheckedAt: start}

	health, err := c.client.Sys().HealthWithContext(ctx)
	s.Latency = c.now().Sub(start)

	if err != nil {
		s.Error = fmt.Errorf("health probe failed: %w", err)
		return s
	}

	s.Initialized = health.Initialized
	s.Sealed = health.Sealed
	s.Standby = health.Standby
	return s
}

// Healthy returns true when Vault is initialised, unsealed, and active.
func (s Status) Healthy() bool {
	return s.Error == nil && s.Initialized && !s.Sealed && !s.Standby
}
