// Package audit provides structured audit logging for secret check
// and rotation events, writing entries to a configurable sink.
package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// Entry represents a single audit log record.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Event     string    `json:"event"`
	Path      string    `json:"path"`
	Status    string    `json:"status"`
	Detail    string    `json:"detail,omitempty"`
}

// Logger writes audit entries as newline-delimited JSON.
type Logger struct {
	out io.Writer
	now func() time.Time
}

// NewLogger returns a Logger that writes to the given file path.
// Pass an empty path to write to stdout.
func NewLogger(path string) (*Logger, error) {
	var w io.Writer
	if path == "" {
		w = os.Stdout
	} else {
		f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o640)
		if err != nil {
			return nil, fmt.Errorf("audit: open log file: %w", err)
		}
		w = f
	}
	return &Logger{out: w, now: time.Now}, nil
}

// Log writes a single audit entry.
func (l *Logger) Log(event, path, status, detail string) error {
	e := Entry{
		Timestamp: l.now().UTC(),
		Event:     event,
		Path:      path,
		Status:    status,
		Detail:    detail,
	}
	b, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("audit: marshal entry: %w", err)
	}
	_, err = fmt.Fprintf(l.out, "%s\n", b)
	return err
}
