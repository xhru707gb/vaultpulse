package ownership

import (
	"fmt"
	"strings"
)

const (
	colPath    = 36
	colOwner   = 20
	colTeam    = 20
	colContact = 28
)

// FormatTable renders ownership entries as a plain-text table.
func FormatTable(entries []Entry) string {
	if len(entries) == 0 {
		return "no ownership records found\n"
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%-*s %-*s %-*s %-*s\n",
		colPath, "PATH",
		colOwner, "OWNER",
		colTeam, "TEAM",
		colContact, "CONTACT",
	)
	sb.WriteString(strings.Repeat("-", colPath+colOwner+colTeam+colContact+3) + "\n")
	for _, e := range entries {
		fmt.Fprintf(&sb, "%-*s %-*s %-*s %-*s\n",
			colPath, truncate(e.Path, colPath),
			colOwner, truncate(e.Owner, colOwner),
			colTeam, truncate(e.Team, colTeam),
			colContact, truncate(e.Contact, colContact),
		)
	}
	return sb.String()
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}
