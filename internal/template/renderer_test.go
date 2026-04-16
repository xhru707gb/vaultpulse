package template_test

import (
	"strings"
	"testing"
	"time"

	"github.com/vaultpulse/internal/template"
)

func makeData() template.Data {
	return template.Data{
		Path:      "secret/db/password",
		Status:    "expired",
		TTL:       "0s",
		ExpiresAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Timestamp: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
		Extra:     map[string]string{"env": "prod"},
	}
}

func TestRender_BasicSubstitution(t *testing.T) {
	r := template.New()
	out, err := r.Render("path={{ .Path }} status={{ .Status }}", makeData())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "secret/db/password") {
		t.Errorf("expected path in output, got: %s", out)
	}
	if !strings.Contains(out, "expired") {
		t.Errorf("expected status in output, got: %s", out)
	}
}

func TestRender_FmtTimeHelper(t *testing.T) {
	r := template.New()
	out, err := r.Render("ts={{ fmtTime .Timestamp }}", makeData())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "2024-06-01T12:00:00Z") {
		t.Errorf("expected formatted timestamp, got: %s", out)
	}
}

func TestRender_ExtraMap(t *testing.T) {
	r := template.New()
	out, err := r.Render("env={{ index .Extra \"env\" }}", makeData())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "prod") {
		t.Errorf("expected extra value, got: %s", out)
	}
}

func TestRender_InvalidTemplate(t *testing.T) {
	r := template.New()
	_, err := r.Render("{{ .Unclosed", makeData())
	if err == nil {
		t.Fatal("expected parse error")
	}
}

func TestRender_EmptyTemplate(t *testing.T) {
	r := template.New()
	out, err := r.Render("", makeData())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "" {
		t.Errorf("expected empty output, got: %s", out)
	}
}
