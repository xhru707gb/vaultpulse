package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the top-level vaultpulse configuration.
type Config struct {
	Vault   VaultConfig   `yaml:"vault"`
	Alerting AlertingConfig `yaml:"alerting"`
	Schedule ScheduleConfig `yaml:"schedule"`
}

// VaultConfig contains connection details for HashiCorp Vault.
type VaultConfig struct {
	Address   string `yaml:"address"`
	Token     string `yaml:"token"`
	Namespace string `yaml:"namespace"`
	SkipTLS   bool   `yaml:"skip_tls_verify"`
}

// AlertingConfig defines alerting hook configuration.
type AlertingConfig struct {
	WebhookURL      string        `yaml:"webhook_url"`
	WarnThreshold   time.Duration `yaml:"warn_threshold"`
	CriticalThreshold time.Duration `yaml:"critical_threshold"`
}

// ScheduleConfig controls how often vaultpulse polls Vault.
type ScheduleConfig struct {
	Interval time.Duration `yaml:"interval"`
}

// Load reads a YAML config file from the given path and returns a Config.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file %q: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file %q: %w", path, err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

// validate enforces required fields and sensible defaults.
func (c *Config) validate() error {
	if c.Vault.Address == "" {
		c.Vault.Address = "http://127.0.0.1:8200"
	}
	if c.Vault.Token == "" {
		if tok := os.Getenv("VAULT_TOKEN"); tok != "" {
			c.Vault.Token = tok
		} else {
			return fmt.Errorf("vault.token is required (or set VAULT_TOKEN env var)")
		}
	}
	if c.Alerting.WarnThreshold == 0 {
		c.Alerting.WarnThreshold = 7 * 24 * time.Hour
	}
	if c.Alerting.CriticalThreshold == 0 {
		c.Alerting.CriticalThreshold = 24 * time.Hour
	}
	if c.Schedule.Interval == 0 {
		c.Schedule.Interval = 5 * time.Minute
	}
	return nil
}
