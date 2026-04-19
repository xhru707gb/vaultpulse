// Package secretreport aggregates secret health data into a unified report.
package secretreport

import (
	"errors"
	"time"
)

// Severity represents the overall health severity of a secret.
type Severity string

const (
	SeverityOK       Severity = "ok"
	SeverityWarning  Severity = "warning"
	SeverityCritical Severity = "critical"
)

// Entry holds aggregated health information for a single secret path.
type Entry struct {
	Path        string
	Severity    Severity
	TTL         time.Duration
	LastRotated time.Time
	Notes       []string
}

// Report holds all entries produced by the reporter.
type Report struct {
	GeneratedAt time.Time
	Entries     []Entry
	Total       int
	AlertCount  int
}

// Reporter builds unified secret reports.
type Reporter struct {
	now func() time.Time
}

// New returns a new Reporter.
func New() *Reporter {
	return &Reporter{now: time.Now}
}

// Build constructs a Report from the provided entries.
// Returns an error if entries is nil.
func (r *Reporter) Build(entries []Entry) (Report, error) {
	if entries == nil {
		return Report{}, errors.New("secretreport: entries must not be nil")
	}

	alerts := 0
	for _, e := range entries {
		if e.Severity == SeverityWarning || e.Severity == SeverityCritical {
			alerts++
		}
	}

	return Report{
		GeneratedAt: r.now().UTC(),
		Entries:     entries,
		Total:       len(entries),
		AlertCount:  alerts,
	}, nil
}
