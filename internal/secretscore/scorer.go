// Package secretscore provides risk scoring for Vault secrets based on
// expiry proximity, rotation overdue status, and policy violations.
package secretscore

import (
	"errors"
	"time"
)

// Risk levels.
const (
	RiskLow      = "low"
	RiskMedium   = "medium"
	RiskHigh     = "high"
	RiskCritical = "critical"
)

// Config holds thresholds used when computing scores.
type Config struct {
	CriticalTTL   time.Duration // TTL below which score is critical
	HighTTL       time.Duration // TTL below which score is high
	MediumTTL     time.Duration // TTL below which score is medium
	RotationBonus int           // extra points added when rotation is overdue
	ViolationBonus int          // extra points added per policy violation
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		CriticalTTL:    24 * time.Hour,
		HighTTL:        72 * time.Hour,
		MediumTTL:      168 * time.Hour,
		RotationBonus:  20,
		ViolationBonus: 15,
	}
}

// Input holds the data needed to score a single secret.
type Input struct {
	Path            string
	TTL             time.Duration
	RotationOverdue bool
	ViolationCount  int
}

// Result holds the computed score and risk level for a secret.
type Result struct {
	Path   string
	Score  int
	Level  string
	Reason string
}

// Scorer computes risk scores for secrets.
type Scorer struct {
	cfg Config
}

// New creates a Scorer with the given config.
func New(cfg Config) (*Scorer, error) {
	if cfg.CriticalTTL <= 0 || cfg.HighTTL <= cfg.CriticalTTL || cfg.MediumTTL <= cfg.HighTTL {
		return nil, errors.New("secretscore: invalid TTL thresholds")
	}
	return &Scorer{cfg: cfg}, nil
}

// Score computes a risk result for the given input.
func (s *Scorer) Score(in Input) Result {
	var score int
	var reason string

	switch {
	case in.TTL <= s.cfg.CriticalTTL:
		score += 75
		reason = "TTL critical"
	case in.TTL <= s.cfg.HighTTL:
		score += 50
		reason = "TTL high"
	case in.TTL <= s.cfg.MediumTTL:
		score += 25
		reason = "TTL medium"
	default:
		reason = "TTL ok"
	}

	if in.RotationOverdue {
		score += s.cfg.RotationBonus
		reason += ", rotation overdue"
	}
	score += in.ViolationCount * s.cfg.ViolationBonus
	if in.ViolationCount > 0 {
		reason += ", policy violations"
	}

	return Result{Path: in.Path, Score: score, Level: levelFor(score), Reason: reason}
}

// ScoreAll scores a slice of inputs.
func (s *Scorer) ScoreAll(inputs []Input) []Result {
	out := make([]Result, len(inputs))
	for i, in := range inputs {
		out[i] = s.Score(in)
	}
	return out
}

func levelFor(score int) string {
	switch {
	case score >= 75:
		return RiskCritical
	case score >= 50:
		return RiskHigh
	case score >= 25:
		return RiskMedium
	default:
		return RiskLow
	}
}
