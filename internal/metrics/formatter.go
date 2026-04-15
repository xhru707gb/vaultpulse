package metrics

import (
	"fmt"
	"strings"
	"time"
)

const (
	headerLine = "%-20s %s"
	rowLine    = "%-20s %v"
)

// FormatTable renders the snapshot as a human-readable table.
func FormatTable(s Snapshot) string {
	var sb strings.Builder

	sb.WriteString("=== VaultPulse Metrics Snapshot ===\n")
	sb.WriteString(fmt.Sprintf("%-20s %s\n", "Collected At", formatTime(s.CollectedAt)))
	sb.WriteString(fmt.Sprintf("%-20s %s\n", "Last Check Duration", formatDuration(s.LastCheckDur)))
	sb.WriteString(strings.Repeat("-", 40) + "\n")
	sb.WriteString(fmt.Sprintf(headerLine+"\n", "Metric", "Value"))
	sb.WriteString(strings.Repeat("-", 40) + "\n")
	sb.WriteString(fmt.Sprintf(rowLine+"\n", "Total Secrets", s.TotalSecrets))
	sb.WriteString(fmt.Sprintf(rowLine+"\n", "Expired", s.Expired))
	sb.WriteString(fmt.Sprintf(rowLine+"\n", "Warning", s.Warning))
	sb.WriteString(fmt.Sprintf(rowLine+"\n", "Healthy", s.Healthy))
	sb.WriteString(fmt.Sprintf(rowLine+"\n", "Overdue Rotation", s.OverdueRotation))
	sb.WriteString(fmt.Sprintf(rowLine+"\n", "Policy Violations", s.PolicyViolation))

	return sb.String()
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return "n/a"
	}
	return t.UTC().Format(time.RFC3339)
}

func formatDuration(d time.Duration) string {
	if d == 0 {
		return "n/a"
	}
	return d.Round(time.Millisecond).String()
}
