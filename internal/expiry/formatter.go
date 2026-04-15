package expiry

import (
	"fmt"
	"io"
	"text/tabwriter"
	"time"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorGreen  = "\033[32m"
)

// FormatTable writes a human-readable table of secret statuses to w.
func FormatTable(w io.Writer, statuses []*SecretStatus, useColor bool) error {
	tw := tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)
	fmt.Fprintln(tw, "PATH\tEXPIRES AT\tTTL\tSTATUS")
	fmt.Fprintln(tw, "----\t----------\t---\t------")

	for _, s := range statuses {
		label, color := statusLabel(s)
		ttlStr := formatTTL(s.TTL)
		expStr := s.ExpiresAt.UTC().Format(time.RFC3339)

		if useColor {
			fmt.Fprintf(tw, "%s\t%s\t%s\t%s%s%s\n",
				s.Path, expStr, ttlStr, color, label, colorReset)
		} else {
			fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n",
				s.Path, expStr, ttlStr, label)
		}
	}
	return tw.Flush()
}

func statusLabel(s *SecretStatus) (string, string) {
	switch {
	case s.IsExpired:
		return "EXPIRED", colorRed
	case s.Warning:
		return "WARNING", colorYellow
	default:
		return "OK", colorGreen
	}
}

func formatTTL(d time.Duration) string {
	if d <= 0 {
		return "expired"
	}
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	if h > 0 {
		return fmt.Sprintf("%dh%dm", h, m)
	}
	return fmt.Sprintf("%dm", m)
}
