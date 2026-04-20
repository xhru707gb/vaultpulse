// Package secretsink provides a lightweight fan-out router that delivers
// secret lifecycle events (expiry, rotation, warnings) to one or more
// pluggable output sinks.
//
// Usage:
//
//	router := secretsink.New()
//	router.Register(mySink)
//	errs := router.Dispatch(secretsink.Event{
//		Path:    "secret/api-key",
//		Kind:    "expired",
//		Message: "secret has expired",
//	})
package secretsink
