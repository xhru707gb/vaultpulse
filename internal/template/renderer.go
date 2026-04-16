// Package template provides message template rendering for alert notifications.
package template

import (
	"bytes"
	"fmt"
	"text/template"
	"time"
)

// Data holds the variables available inside a template.
type Data struct {
	Path      string
	Status    string
	TTL       string
	ExpiresAt time.Time
	Timestamp time.Time
	Extra     map[string]string
}

// Renderer renders text templates with alert data.
type Renderer struct {
	funcs template.FuncMap
}

// New returns a Renderer with a set of built-in helper functions.
func New() *Renderer {
	return &Renderer{
		funcs: template.FuncMap{
			"upper": func(s string) string {
				return fmt.Sprintf("%s", bytes.ToUpper([]byte(s)))
			},
			"fmtTime": func(t time.Time) string {
				return t.UTC().Format(time.RFC3339)
			},
		},
	}
}

// Render executes the given template text with d and returns the result.
func (r *Renderer) Render(tmplText string, d Data) (string, error) {
	t, err := template.New("").Funcs(r.funcs).Parse(tmplText)
	if err != nil {
		return "", fmt.Errorf("template parse: %w", err)
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, d); err != nil {
		return "", fmt.Errorf("template execute: %w", err)
	}
	return buf.String(), nil
}
