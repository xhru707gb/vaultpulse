// Package tokenwatch monitors Vault token TTLs and emits alerts
// when tokens are approaching expiration.
package tokenwatch

import (
	"context"
	"fmt"
	"time"
)

// TokenInfo holds metadata about a Vault token.
type TokenInfo struct {
	Accessor   string
	DisplayName string
	TTL        time.Duration
	ExpireTime time.Time
}

// Status represents the evaluated state of a token.
type Status struct {
	Token     TokenInfo
	State     string // "ok", "warning", "expired"
	Remaining time.Duration
}

// Fetcher retrieves token info from Vault.
type Fetcher interface {
	LookupTokens(ctx context.Context) ([]TokenInfo, error)
}

// Watcher evaluates token expiry states.
type Watcher struct {
	fetcher       Fetcher
	warnThreshold time.Duration
	now           func() time.Time
}

// New creates a Watcher with the given fetcher and warning threshold.
func New(f Fetcher, warnThreshold time.Duration) (*Watcher, error) {
	if f == nil {
		return nil, fmt.Errorf("tokenwatch: fetcher must not be nil")
	}
	if warnThreshold <= 0 {
		return nil, fmt.Errorf("tokenwatch: warnThreshold must be positive")
	}
	return &Watcher{fetcher: f, warnThreshold: warnThreshold, now: time.Now}, nil
}

// Evaluate fetches tokens and returns their expiry statuses.
func (w *Watcher) Evaluate(ctx context.Context) ([]Status, error) {
	tokens, err := w.fetcher.LookupTokens(ctx)
	if err != nil {
		return nil, fmt.Errorf("tokenwatch: lookup failed: %w", err)
	}
	now := w.now()
	statuses := make([]Status, 0, len(tokens))
	for _, t := range tokens {
		remaining := t.ExpireTime.Sub(now)
		state := "ok"
		if remaining <= 0 {
			state = "expired"
		} else if remaining <= w.warnThreshold {
			state = "warning"
		}
		statuses = append(statuses, Status{Token: t, State: state, Remaining: remaining})
	}
	return statuses, nil
}
