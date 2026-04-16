// Package notify provides a unified dispatch layer that fans out alerts
// to multiple configured channels (webhook, log, stdout).
package notify

import (
	"context"
	"fmt"
	"io"
	"os"
)

// Channel represents a named output target for alert dispatch.
type Channel string

const (
	ChannelWebhook Channel = "webhook"
	ChannelStdout  Channel = "stdout"
	ChannelLog     Channel = "log"
)

// Message is the normalised payload sent to every channel.
type Message struct {
	Level   string // "info", "warning", "critical"
	Subject string
	Body    string
}

// Sender is implemented by any channel that can deliver a Message.
type Sender interface {
	Send(ctx context.Context, msg Message) error
	Name() Channel
}

// Dispatcher fans a Message out to all registered Senders.
type Dispatcher struct {
	senders []Sender
	out     io.Writer // fallback writer for errors
}

// New returns a Dispatcher with the given senders.
func New(senders ...Sender) *Dispatcher {
	return &Dispatcher{senders: senders, out: os.Stderr}
}

// Dispatch sends msg to every registered Sender.
// Errors from individual senders are collected and returned as a combined error.
func (d *Dispatcher) Dispatch(ctx context.Context, msg Message) error {
	var errs []error
	for _, s := range d.senders {
		if err := s.Send(ctx, msg); err != nil {
			fmt.Fprintf(d.out, "notify: sender %q error: %v\n", s.Name(), err)
			errs = append(errs, fmt.Errorf("%s: %w", s.Name(), err))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("dispatch errors: %v", errs)
	}
	return nil
}

// RegisterSender appends a Sender at runtime.
func (d *Dispatcher) RegisterSender(s Sender) {
	d.senders = append(d.senders, s)
}
