// Package ttlpolicy provides enforcement of minimum and maximum TTL bounds
// for Vault secret paths. Rules are matched by path prefix; the first
// matching rule is applied. Secrets with no matching rule are considered
// compliant by default.
//
// Usage:
//
//	e, err := ttlpolicy.New(rules)
//	result := e.Evaluate(path, ttl)
package ttlpolicy
