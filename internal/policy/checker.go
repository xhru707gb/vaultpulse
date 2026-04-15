// Package policy evaluates Vault secrets against defined compliance policies.
package policy

import (
	"fmt"
	"time"
)

// Policy defines rules a secret must satisfy.
type Policy struct {
	Name           string
	MaxTTLDays     int
	RequireRotation bool
	RotationDays   int
}

// Status holds the result of evaluating a secret against a policy.
type Status struct {
	Path       string
	Policy     string
	Compliant  bool
	Violations []string
	CheckedAt  time.Time
}

// Checker evaluates secrets against policies.
type Checker struct {
	policies []Policy
	now      func() time.Time
}

// NewChecker returns a Checker with the provided policies.
func NewChecker(policies []Policy, now func() time.Time) *Checker {
	if now == nil {
		now = time.Now
	}
	return &Checker{policies: policies, now: now}
}

// Evaluate checks a single secret path against all matching policies.
func (c *Checker) Evaluate(path string, ttl time.Duration, lastRotated time.Time) []Status {
	results := make([]Status, 0, len(c.policies))
	for _, p := range c.policies {
		s := Status{
			Path:      path,
			Policy:    p.Name,
			Compliant: true,
			CheckedAt: c.now(),
		}
		if p.MaxTTLDays > 0 {
			maxTTL := time.Duration(p.MaxTTLDays) * 24 * time.Hour
			if ttl > maxTTL {
				s.Compliant = false
				s.Violations = append(s.Violations,
					fmt.Sprintf("TTL %s exceeds max %dd", ttl.Round(time.Hour), p.MaxTTLDays))
			}
		}
		if p.RequireRotation && p.RotationDays > 0 && !lastRotated.IsZero() {
			age := c.now().Sub(lastRotated)
			max := time.Duration(p.RotationDays) * 24 * time.Hour
			if age > max {
				s.Compliant = false
				s.Violations = append(s.Violations,
					fmt.Sprintf("last rotated %dd ago, max %dd", int(age.Hours()/24), p.RotationDays))
			}
		}
		results = append(results, s)
	}
	return results
}

// EvaluateAll evaluates multiple secrets and returns all statuses.
func (c *Checker) EvaluateAll(secrets map[string]SecretMeta) []Status {
	var all []Status
	for path, meta := range secrets {
		all = append(all, c.Evaluate(path, meta.TTL, meta.LastRotated)...)
	}
	return all
}

// SecretMeta holds metadata for a secret to be evaluated.
type SecretMeta struct {
	TTL         time.Duration
	LastRotated time.Time
}
