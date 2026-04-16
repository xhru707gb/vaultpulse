package notify

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// StdoutSender writes formatted alert messages to an io.Writer (default: os.Stdout).
type StdoutSender struct {
	w io.Writer
}

// NewStdoutSender returns a StdoutSender that writes to w.
// If w is nil, os.Stdout is used.
func NewStdoutSender(w io.Writer) *StdoutSender {
	if w == nil {
		w = os.Stdout
	}
	return &StdoutSender{w: w}
}

// Name implements Sender.
func (s *StdoutSender) Name() Channel { return ChannelStdout }

// Send writes a single-line alert to the writer.
func (s *StdoutSender) Send(_ context.Context, msg Message) error {
	level := strings.ToUpper(msg.Level)
	_, err := fmt.Fprintf(s.w, "[%s] %s  %s  %s\n",
		time.Now().UTC().Format(time.RFC3339),
		level,
		msg.Subject,
		msg.Body,
	)
	return err
}
