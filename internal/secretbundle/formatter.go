package secretbundle

import (
	"fmt"
	"strings"
)

const (
	_labelHealthy  = "OK"
	_labelDegraded = "DEGRADED"
)

func healthLabel(healthy bool) string {
	if healthy {
		return _labelHealthy
	}
	return _labelDegraded
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}

// FormatTable renders a slice of EvalResult as an aligned table.
func FormatTable(results []EvalResult) string {
	if len(results) == 0 {
		return "no bundles registered\n"
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%-30s  %7s  %7s  %8s\n", "BUNDLE", "TOTAL", "EXPIRED", "STATUS")
	fmt.Fprintf(&sb, "%s\n", strings.Repeat("-", 58))
	for _, r := range results {
		fmt.Fprintf(&sb, "%-30s  %7d  %7d  %8s\n",
			truncate(r.Name, 30),
			r.Total,
			r.Expired,
			healthLabel(r.Healthy),
		)
	}
	return sb.String()
}

// FormatSummary returns a one-line summary across all results.
func FormatSummary(results []EvalResult) string {
	total, degraded := 0, 0
	for _, r := range results {
		total++
		if !r.Healthy {
			degraded++
		}
	}
	return fmt.Sprintf("bundles: %d total, %d degraded, %d healthy\n",
		total, degraded, total-degraded)
}
