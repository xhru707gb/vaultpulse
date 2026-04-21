// Package secretrelay implements a fan-out relay for secret change events.
//
// A Relay holds a named set of Handler functions. When Dispatch is called with
// a secret path and its payload, every registered handler is invoked in order.
// Errors from individual handlers are collected and returned together so that a
// single failing handler does not block the others.
//
// Typical usage:
//
//	relay := secretrelay.New()
//	_ = relay.Register("webhook", myWebhookHandler)
//	_ = relay.Register("audit",   myAuditHandler)
//
//	errs := relay.Dispatch("secret/db/password", payload)
package secretrelay
