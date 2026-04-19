package secretscore_test

import (
	"testing"
	"time"

	"github.com/yourusername/vaultpulse/internal/secretscore"
)

func defaultCfg() secretscore.Config {
	return secretscore.DefaultConfig()
}

func TestNew_InvalidThresholds(t *testing.T) {
	cfg := secretscore.DefaultConfig()
	cfg.CriticalTTL = cfg.HighTTL // equal — invalid
	_, err := secretscore.New(cfg)
	if err == nil {
		t.Fatal("expected error for invalid thresholds")
	}
}

func TestScore_CriticalTTL(t *testing.T) {
	s, _ := secretscore.New(defaultCfg())
	r := s.Score(secretscore.Input{Path: "secret/a", TTL: 1 * time.Hour})
	if r.Level != secretscore.RiskCritical {
		t.Fatalf("expected critical, got %s", r.Level)
	}
	if r.Score < 75 {
		t.Fatalf("expected score >= 75, got %d", r.Score)
	}
}

func TestScore_HighTTL(t *testing.T) {
	s, _ := secretscore.New(defaultCfg())
	r := s.Score(secretscore.Input{Path: "secret/b", TTL: 48 * time.Hour})
	if r.Level != secretscore.RiskHigh {
		t.Fatalf("expected high, got %s", r.Level)
	}
}

func TestScore_LowTTL(t *testing.T) {
	s, _ := secretscore.New(defaultCfg())
	r := s.Score(secretscore.Input{Path: "secret/c", TTL: 720 * time.Hour})
	if r.Level != secretscore.RiskLow {
		t.Fatalf("expected low, got %s", r.Level)
	}
}

func TestScore_RotationBonusAdded(t *testing.T) {
	s, _ := secretscore.New(defaultCfg())
	without := s.Score(secretscore.Input{Path: "p", TTL: 720 * time.Hour})
	with := s.Score(secretscore.Input{Path: "p", TTL: 720 * time.Hour, RotationOverdue: true})
	if with.Score != without.Score+defaultCfg().RotationBonus {
		t.Fatalf("expected bonus of %d, diff=%d", defaultCfg().RotationBonus, with.Score-without.Score)
	}
}

func TestScore_ViolationBonusAdded(t *testing.T) {
	s, _ := secretscore.New(defaultCfg())
	r := s.Score(secretscore.Input{Path: "p", TTL: 720 * time.Hour, ViolationCount: 2})
	expected := 2 * defaultCfg().ViolationBonus
	if r.Score != expected {
		t.Fatalf("expected score %d, got %d", expected, r.Score)
	}
}

func TestScoreAll_ReturnsAllResults(t *testing.T) {
	s, _ := secretscore.New(defaultCfg())
	inputs := []secretscore.Input{
		{Path: "a", TTL: 1 * time.Hour},
		{Path: "b", TTL: 200 * time.Hour},
	}
	out := s.ScoreAll(inputs)
	if len(out) != 2 {
		t.Fatalf("expected 2 results, got %d", len(out))
	}
}
