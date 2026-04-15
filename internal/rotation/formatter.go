package rotation

import (
	"fmt"
	"io"
	"text/tabwriter"
	"time"
)

const (
	labelOverdue = "OVERDUE"
	labelDueSoon = "DUE SOON"
	labelOK      = "OK"
	dueSoonThreshold = 24 * time.Hour
)

// FormatTable writes a human-readable rotation status table to w.
func FormatTable(w io.Writer, statuses []Status) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "PATH\tSTATUS\tLAST ROTATED\tNEXT DUE\tDUE IN")
	fmt.Fprintln(tw, "----\t------\t------------\t--------\t------")
	for _, s := range statuses {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\n",
			s.Path,
			rotationLabel(s),
			s.LastRotated.UTC().Format(time.RFC3339),
			s.NextDue.UTC().Format(time.RFC3339),
			formatDueIn(s.DueIn),
		)
	}
	tw.Flush()
}

func rotationLabel(s Status) string {
	switch {
	case s.Overdue:
		return labelOverdue
	case s.DueIn <= dueSoonThreshold:
		return labelDueSoon
	default:
		return labelOK
	}
}

func formatDueIn(d time.Duration) string {
	if d < 0 {
		return fmt.Sprintf("-%s ago", d.Abs().Round(time.Minute))
	}
	return d.Round(time.Minute).String()
}
