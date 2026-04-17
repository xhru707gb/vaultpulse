// Package digest provides periodic secret-state summarisation,
// rolling up expiry, rotation and health signals into a single report.
package digest

import (
	"fmt"
	"io"
	"time"
)

// Entry holds a single secret's aggregated state.
type Entry struct {
	Path        string
	Expired     bool
	ExpiresSoon bool
	Overdue     bool
	Unhealthy   bool
	TTL         time.Duration
}

// Report is the output of a digest run.
type Report struct {
	GeneratedAt  time.Time
	Entries      []Entry
	TotalSecrets int
	AlertCount   int
}

// Builder assembles a digest Report from individual signal slices.
type Builder struct {
	now func() time.Time
}

// NewBuilder returns a Builder. If now is nil, time.Now is used.
func NewBuilder(now func() time.Time) *Builder {
	if now == nil {
		now = time.Now
	}
	return &Builder{now: now}
}

// Build creates a Report from the supplied entries.
func (b *Builder) Build(entries []Entry) Report {
	alerts := 0
	for _, e := range entries {
		if e.Expired || e.ExpiresSoon || e.Overdue || e.Unhealthy {
			alerts++
		}
	}
	return Report{
		GeneratedAt:  b.now().UTC(),
		Entries:      entries,
		TotalSecrets: len(entries),
		AlertCount:   alerts,
	}
}

// WriteTo writes a human-readable digest to w.
func (r Report) WriteTo(w io.Writer) (int64, error) {
	n, err := fmt.Fprintf(w,
		"Digest generated at %s\nTotal secrets: %d  Alerts: %d\n",
		r.GeneratedAt.Format(time.RFC3339), r.TotalSecrets, r.AlertCount,
	)
	return int64(n), err
}
