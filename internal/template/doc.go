// Package template provides lightweight text template rendering for
// VaultPulse alert and notification messages.
//
// A Renderer is created via New and exposes a Render method that accepts
// a Go text/template string and a Data struct. Built-in helper functions
// (fmtTime, upper) are available inside every template.
//
// Example:
//
//	r := template.New()
//	msg, err := r.Render("[{{ .Status }}] {{ .Path }} expires {{ fmtTime .ExpiresAt }}", data)
package template
