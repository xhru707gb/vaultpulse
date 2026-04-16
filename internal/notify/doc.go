// Package notify implements a multi-channel alert dispatcher for VaultPulse.
//
// A Dispatcher holds one or more Sender implementations and fans out a
// normalised Message to all of them in order.  Built-in senders:
//
//   - StdoutSender  – writes to any io.Writer (default os.Stdout)
//
// Additional senders (webhook, log file) can be registered via
// RegisterSender at runtime.
//
// Usage:
//
//	d := notify.New(notify.NewStdoutSender(nil))
//	d.Dispatch(ctx, notify.Message{Level: "warning", Subject: "secret/db", Body: "expires in 2h"})
package notify
