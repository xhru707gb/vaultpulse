// Package secretclassify provides classification of secrets into sensitivity tiers
// based on configurable path patterns and metadata rules.
package secretclassify

import (
	"errors"
	"strings"
	"sync"
)

// Level represents a classification sensitivity tier.
type Level string

const (
	LevelPublic       Level = "public"
	LevelInternal     Level = "internal"
	LevelConfidential Level = "confidential"
	LevelSecret       Level = "secret"
)

// Rule maps a path prefix or substring to a classification level.
type Rule struct {
	Pattern string
	Level   Level
}

// Result holds the classification outcome for a single secret path.
type Result struct {
	Path  string
	Level Level
}

// Classifier assigns sensitivity levels to secret paths.
type Classifier struct {
	mu      sync.RWMutex
	rules   []Rule
	fallback Level
}

// New creates a Classifier with the given rules and a fallback level
// applied when no rule matches.
func New(rules []Rule, fallback Level) (*Classifier, error) {
	if len(rules) == 0 {
		return nil, errors.New("secretclassify: at least one rule is required")
	}
	for _, r := range rules {
		if strings.TrimSpace(r.Pattern) == "" {
			return nil, errors.New("secretclassify: rule pattern must not be empty")
		}
	}
	return &Classifier{rules: rules, fallback: fallback}, nil
}

// Classify returns the classification Result for the given path.
func (c *Classifier) Classify(path string) Result {
	c.mu.RLock()
	defer c.mu.RUnlock()
	for _, r := range c.rules {
		if strings.Contains(path, r.Pattern) {
			return Result{Path: path, Level: r.Level}
		}
	}
	return Result{Path: path, Level: c.fallback}
}

// ClassifyAll classifies a slice of paths and returns results in order.
func (c *Classifier) ClassifyAll(paths []string) []Result {
	out := make([]Result, len(paths))
	for i, p := range paths {
		out[i] = c.Classify(p)
	}
	return out
}
