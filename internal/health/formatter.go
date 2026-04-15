package health

import (
	"fmt"
	"io"
	"text/tabwriter"
)

const (
	labelOK      = "OK"
	labelSealed  = "SEALED"
	labelStandby = "STANDBY"
	labelError   = "ERROR"
)

// FormatTable writes a human-readable health report to w.
func FormatTable(w io.Writer, s Status) {
	tw := tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)
	fmt.Fprintln(tw, "STATUS\tINITIALIZED\tSEALED\tSTANDBY\tLATENCY\tCHECKED AT")
	fmt.Fprintf(tw, "%s\t%v\t%v\t%v\t%s\t%s\n",
		statusLabel(s),
		s.Initialized,
		s.Sealed,
		s.Standby,
		fmt.Sprintf("%dms", s.Latency.Milliseconds()),
		s.CheckedAt.UTC().Format("2006-01-02 15:04:05"),
	)
	_ = tw.Flush()
}

func statusLabel(s Status) string {
	switch {
	case s.Error != nil:
		return labelError
	case s.Sealed:
		return labelSealed
	case s.Standby:
		return labelStandby
	default:
		return labelOK
	}
}
