package config

import (
	"os"
	"testing"
	"time"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "vaultpulse-*.yaml")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_ValidConfig(t *testing.T) {
	path := writeTempConfig(t, `
vault:
  address: "https://vault.example.com"
  token: "s.testtoken"
alerting:
  webhook_url: "https://hooks.example.com/alert"
  warn_threshold: 168h
  critical_threshold: 24h
schedule:
  interval: 10m
`)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Vault.Address != "https://vault.example.com" {
		t.Errorf("expected vault address, got %q", cfg.Vault.Address)
	}
	if cfg.Schedule.Interval != 10*time.Minute {
		t.Errorf("expected 10m interval, got %v", cfg.Schedule.Interval)
	}
}

func TestLoad_DefaultsApplied(t *testing.T) {
	t.Setenv("VAULT_TOKEN", "s.envtoken")
	path := writeTempConfig(t, `vault: {}`)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Vault.Address != "http://127.0.0.1:8200" {
		t.Errorf("expected default address, got %q", cfg.Vault.Address)
	}
	if cfg.Vault.Token != "s.envtoken" {
		t.Errorf("expected token from env, got %q", cfg.Vault.Token)
	}
	if cfg.Alerting.WarnThreshold != 7*24*time.Hour {
		t.Errorf("expected default warn threshold, got %v", cfg.Alerting.WarnThreshold)
	}
	if cfg.Schedule.Interval != 5*time.Minute {
		t.Errorf("expected default interval, got %v", cfg.Schedule.Interval)
	}
}

func TestLoad_MissingToken(t *testing.T) {
	os.Unsetenv("VAULT_TOKEN")
	path := writeTempConfig(t, `vault:
  address: "http://127.0.0.1:8200"
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for missing token, got nil")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := Load("/nonexistent/path/config.yaml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
