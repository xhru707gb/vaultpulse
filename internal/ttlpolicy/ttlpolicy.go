// Package ttlpolicy enforces minimum and maximum TTL bounds on Vault secrets.
package ttlpolicy

import (
	"errors"
	"fmt"
	"time"
)

// Rule defines TTL bounds for a path prefix.
type Rule struct {
	Prefix     string
	MinTTL     time.Duration
	MaxTTL     time.Duration
}

// Result holds the outcome of evaluating a secret against TTL policy rules.
type Result struct {
	Path      string
	TTL       time.Duration
	Compliant bool
	Violation string
}

// Enforcer evaluates secret TTLs against a set of rules.
type Enforcer struct {
	rules []Rule
}

// New creates a new Enforcer. Returns an error if no rules are provided.
func New(rules []Rule) (*Enforcer, error) {
	if len(rules) == 0 {
		return nil, errors.New("ttlpolicy: at least one rule is required")
	}
	for _, r := range rules {
		if r.Prefix == "" {
			return nil, errors.New("ttlpolicy: rule prefix must not be empty")
		}
		if r.MaxTTL > 0 && r.MinTTL > r.MaxTTL {
			return nil, fmt.Errorf("ttlpolicy: minTTL exceeds maxTTL for prefix %q", r.Prefix)
		}
	}
	return &Enforcer{rules: rules}, nil
}

// Evaluate checks a single secret path and TTL against matching rules.
// The first matching rule (by prefix) is applied.
func (e *Enforcer) Evaluate(path string, ttl time.Duration) Result {
	for _, r := range e.rules {
		if len(path) >= len(r.Prefix) && path[:len(r.Prefix)] == r.Prefix {
			return applyRule(path, ttl, r)
		}
	}
	return Result{Path: path, TTL: ttl, Compliant: true}
}

// EvaluateAll evaluates a map of path→TTL pairs and returns all results.
func (e *Enforcer) EvaluateAll(secrets map[string]time.Duration) []Result {
	results := make([]Result, 0, len(secrets))
	for path, ttl := range secrets {
		results = append(results, e.Evaluate(path, ttl))
	}
	return results
}

func applyRule(path string, ttl time.Duration, r Rule) Result {
	if r.MinTTL > 0 && ttl < r.MinTTL {
		return Result{
			Path:      path,
			TTL:       ttl,
			Compliant: false,
			Violation: fmt.Sprintf("TTL %s is below minimum %s", ttl, r.MinTTL),
		}
	}
	if r.MaxTTL > 0 && ttl > r.MaxTTL {
		return Result{
			Path:      path,
			TTL:       ttl,
			Compliant: false,
			Violation: fmt.Sprintf("TTL %s exceeds maximum %s", ttl, r.MaxTTL),
		}
	}
	return Result{Path: path, TTL: ttl, Compliant: true}
}
