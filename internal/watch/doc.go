// Package watch implements a periodic polling loop for Vault secret expiry
// status. A Watcher runs at a configurable interval, invokes the expiry
// checker for a set of paths, and dispatches the results to a caller-supplied
// Handler function.
//
// Typical usage:
//
//	w, _ := watch.New(checker, 30*time.Second, func(e watch.Event) {
//	    fmt.Println(expiry.FormatTable(e.Statuses))
//	})
//	w.Run(ctx, []string{"secret/db", "secret/api"})
package watch
