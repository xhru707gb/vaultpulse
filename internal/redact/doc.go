// Package redact provides a Redactor type that masks sensitive secret
// metadata values — such as tokens, passwords, and keys — before they
// are written to logs, rendered in tables, or forwarded via webhook
// alert payloads.
//
// Patterns are matched case-insensitively against map keys or field
// names. Any matching value is replaced with the string "[REDACTED]".
package redact
