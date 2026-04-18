package renew

import (
	"fmt"
	"strings"
	"time"
)

const (
	maxPath = 40
)

// FormatTable renders renewal entries as a plain-text table.
func FormatTable(entries []*Entry) string {
	if len(entries) == 0 {
		return "No secrets currently tracked for renewal.\n"
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%-42s %-12s %-10s %s\n", "PATH", "LEASE TTL", "RENEWALS", "NEXT RENEW AT")
	fmt.Fprintf(&sb, "%s\n", strings.Repeat("-", 82))
	for _, e := range entries {
		fmt.Fprintf(&sb, "%-42s %-12s %-10d %s\n",
			truncate(e.Path, maxPath),
			formatDuration(e.LeaseTTL),
			e.RenewCount,
			e.RenewAt.UTC().Format(time.RFC3339),
		)
	}
	return sb.String()
}

func formatDuration(d time.Duration) string {
	if d == 0 {
		return "unknown"
	}
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	if h > 0 {
		return fmt.Sprintf("%dh%dm", h, m)
	}
	if m > 0 {
		return fmt.Sprintf("%dm%ds", m, s)
	}
	return fmt.Sprintf("%ds", s)
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return "..." + s[len(s)-n+3:]
}
