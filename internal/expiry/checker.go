package expiry

import (
	"fmt"
	"time"

	"github.com/vaultpulse/internal/vault"
)

// SecretStatus holds expiry information for a single secret.
type SecretStatus struct {
	Path      string
	ExpiresAt time.Time
	TTL       time.Duration
	IsExpired bool
	Warning   bool // true when TTL is within the warning threshold
}

// Checker evaluates secret expiry against configurable thresholds.
type Checker struct {
	client          *vault.Client
	warningThreshold time.Duration
}

// NewChecker creates a Checker with the given Vault client and warning threshold.
func NewChecker(client *vault.Client, warningThreshold time.Duration) *Checker {
	return &Checker{
		client:          client,
		warningThreshold: warningThreshold,
	}
}

// Check retrieves metadata for the secret at path and returns its status.
func (c *Checker) Check(path string) (*SecretStatus, error) {
	meta, err := c.client.GetSecretMeta(path)
	if err != nil {
		return nil, fmt.Errorf("checker: failed to get secret meta for %q: %w", path, err)
	}

	now := time.Now()
	expiresAt, ok := meta["expiration"].(time.Time)
	if !ok {
		// Attempt string parse fallback
		if s, ok2 := meta["expiration"].(string); ok2 {
			parsed, parseErr := time.Parse(time.RFC3339, s)
			if parseErr != nil {
				return nil, fmt.Errorf("checker: cannot parse expiration for %q: %w", path, parseErr)
			}
			expiresAt = parsed
		} else {
			return nil, fmt.Errorf("checker: no expiration field for secret %q", path)
		}
	}

	ttl := expiresAt.Sub(now)
	return &SecretStatus{
		Path:      path,
		ExpiresAt: expiresAt,
		TTL:       ttl,
		IsExpired: ttl <= 0,
		Warning:   ttl > 0 && ttl <= c.warningThreshold,
	}, nil
}

// CheckAll checks multiple secret paths and returns all statuses.
func (c *Checker) CheckAll(paths []string) ([]*SecretStatus, error) {
	statuses := make([]*SecretStatus, 0, len(paths))
	for _, p := range paths {
		status, err := c.Check(p)
		if err != nil {
			return nil, err
		}
		statuses = append(statuses, status)
	}
	return statuses, nil
}
