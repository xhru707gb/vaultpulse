// Package secretpriority assigns and tracks priority levels to secrets
// based on configurable rules, enabling triage of high-impact secrets first.
package secretpriority

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

// Level represents the priority of a secret.
type Level int

const (
	LevelLow      Level = iota // default
	LevelMedium               // moderate importance
	LevelHigh                 // elevated importance
	LevelCritical             // must be addressed immediately
)

// Rule maps a path prefix to a priority level.
type Rule struct {
	Prefix string
	Level  Level
}

// Result holds the evaluated priority for a single secret path.
type Result struct {
	Path  string
	Level Level
	Rule  string // matched rule prefix, or "default"
}

// Evaluator assigns priority levels to secret paths.
type Evaluator struct {
	mu      sync.RWMutex
	rules   []Rule
	default_ Level
}

// New creates an Evaluator with the provided rules and a fallback default level.
// Returns an error if no rules are provided.
func New(rules []Rule, defaultLevel Level) (*Evaluator, error) {
	if len(rules) == 0 {
		return nil, errors.New("secretpriority: at least one rule is required")
	}
	return &Evaluator{rules: rules, default_: defaultLevel}, nil
}

// Evaluate returns the priority Result for the given secret path.
func (e *Evaluator) Evaluate(path string) Result {
	e.mu.RLock()
	defer e.mu.RUnlock()

	for _, r := range e.rules {
		if strings.HasPrefix(path, r.Prefix) {
			return Result{Path: path, Level: r.Level, Rule: r.Prefix}
		}
	}
	return Result{Path: path, Level: e.default_, Rule: "default"}
}

// EvaluateAll returns priority Results for all provided paths.
func (e *Evaluator) EvaluateAll(paths []string) []Result {
	out := make([]Result, 0, len(paths))
	for _, p := range paths {
		out = append(out, e.Evaluate(p))
	}
	return out
}

// LevelLabel returns a human-readable label for a Level.
func LevelLabel(l Level) string {
	switch l {
	case LevelCritical:
		return "CRITICAL"
	case LevelHigh:
		return "HIGH"
	case LevelMedium:
		return "MEDIUM"
	default:
		return fmt.Sprintf("LOW(%d)", l)
	}
}
