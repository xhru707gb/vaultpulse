// Package diff provides snapshot diffing to detect secret changes between runs.
package diff

import (
	"fmt"
	"time"
)

// ChangeKind describes the type of change detected between two snapshots.
type ChangeKind string

const (
	ChangeAdded   ChangeKind = "added"
	ChangeRemoved ChangeKind = "removed"
	ChangeModified ChangeKind = "modified"
)

// SecretEntry represents a single secret's state at a point in time.
type SecretEntry struct {
	Path      string    `json:"path"`
	ExpiresAt time.Time `json:"expires_at"`
	Version   int       `json:"version"`
}

// Change represents a detected difference between two snapshots.
type Change struct {
	Path string
	Kind ChangeKind
	Prev *SecretEntry
	Curr *SecretEntry
}

// String returns a human-readable description of the change.
func (c Change) String() string {
	switch c.Kind {
	case ChangeAdded:
		return fmt.Sprintf("[added]    %s (version %d)", c.Path, c.Curr.Version)
	case ChangeRemoved:
		return fmt.Sprintf("[removed]  %s (was version %d)", c.Path, c.Prev.Version)
	case ChangeModified:
		return fmt.Sprintf("[modified] %s (v%d -> v%d)", c.Path, c.Prev.Version, c.Curr.Version)
	default:
		return fmt.Sprintf("[unknown]  %s", c.Path)
	}
}

// Compute compares two slices of SecretEntry (previous and current) and returns
// the list of changes. Entries are matched by Path.
func Compute(prev, curr []SecretEntry) []Change {
	prevMap := make(map[string]SecretEntry, len(prev))
	for _, e := range prev {
		prevMap[e.Path] = e
	}

	currMap := make(map[string]SecretEntry, len(curr))
	for _, e := range curr {
		currMap[e.Path] = e
	}

	var changes []Change

	for path, currEntry := range currMap {
		prevEntry, existed := prevMap[path]
		if !existed {
			c := currEntry
			changes = append(changes, Change{Path: path, Kind: ChangeAdded, Curr: &c})
		} else if prevEntry.Version != currEntry.Version || !prevEntry.ExpiresAt.Equal(currEntry.ExpiresAt) {
			p, c := prevEntry, currEntry
			changes = append(changes, Change{Path: path, Kind: ChangeModified, Prev: &p, Curr: &c})
		}
	}

	for path, prevEntry := range prevMap {
		if _, stillExists := currMap[path]; !stillExists {
			p := prevEntry
			changes = append(changes, Change{Path: path, Kind: ChangeRemoved, Prev: &p})
		}
	}

	return changes
}
