package vault

import (
	"context"
	"fmt"
	"time"

	vaultapi "github.com/hashicorp/vault/api"
)

// SecretMeta holds metadata about a Vault secret relevant to expiration tracking.
type SecretMeta struct {
	Path       string
	Expiration time.Time
	TTL        time.Duration
	Renewable  bool
}

// Client wraps the Vault API client with vaultpulse-specific operations.
type Client struct {
	vc      *vaultapi.Client
	Address string
}

// NewClient creates a new Vault client using the provided address and token.
func NewClient(address, token string) (*Client, error) {
	cfg := vaultapi.DefaultConfig()
	cfg.Address = address

	vc, err := vaultapi.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create vault client: %w", err)
	}

	vc.SetToken(token)

	return &Client{
		vc:      vc,
		Address: address,
	}, nil
}

// GetSecretMeta reads a secret at the given path and returns its expiration metadata.
func (c *Client) GetSecretMeta(ctx context.Context, path string) (*SecretMeta, error) {
	secret, err := c.vc.Logical().ReadWithContext(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to read secret at %q: %w", path, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("secret not found at path %q", path)
	}

	ttl := time.Duration(secret.LeaseDuration) * time.Second
	expiration := time.Now().Add(ttl)

	return &SecretMeta{
		Path:       path,
		Expiration: expiration,
		TTL:        ttl,
		Renewable:  secret.Renewable,
	}, nil
}

// Ping checks connectivity to the Vault server by fetching its health status.
func (c *Client) Ping(ctx context.Context) error {
	_, err := c.vc.Sys().HealthWithContext(ctx)
	if err != nil {
		return fmt.Errorf("vault unreachable at %s: %w", c.Address, err)
	}
	return nil
}
