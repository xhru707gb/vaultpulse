package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// ReportEntry represents a single line in a human-readable audit report.
type ReportEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Path      string    `json:"path"`
	Event     string    `json:"event"`
	Status    string    `json:"status"`
	TTL       string    `json:"ttl,omitempty"`
}

// Report reads newline-delimited JSON audit log entries from r and writes a
// formatted plain-text report to w. Lines that cannot be decoded are skipped.
func Report(r io.Reader, w io.Writer) error {
	decoder := json.NewDecoder(r)

	fmt.Fprintf(w, "%-30s %-40s %-12s %-10s %s\n",
		"TIMESTAMP", "PATH", "EVENT", "STATUS", "TTL")
	fmt.Fprintf(w, "%s\n", repeatChar('-', 100))

	var count int
	for decoder.More() {
		var entry ReportEntry
		if err := decoder.Decode(&entry); err != nil {
			continue
		}
		fmt.Fprintf(w, "%-30s %-40s %-12s %-10s %s\n",
			entry.Timestamp.UTC().Format(time.RFC3339),
			truncate(entry.Path, 40),
			entry.Event,
			entry.Status,
			entry.TTL,
		)
		count++
	}

	fmt.Fprintf(w, "\n%d audit record(s) shown.\n", count)
	return nil
}

func repeatChar(ch rune, n int) string {
	buf := make([]rune, n)
	for i := range buf {
		buf[i] = ch
	}
	return string(buf)
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}
