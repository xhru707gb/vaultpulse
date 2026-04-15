package expiry

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func makeStatus(path string, ttl time.Duration) *SecretStatus {
	now := time.Now()
	return &SecretStatus{
		Path:      path,
		ExpiresAt: now.Add(ttl),
		TTL:       ttl,
		IsExpired: ttl <= 0,
		Warning:   ttl > 0 && ttl <= 24*time.Hour,
	}
}

func TestFormatTable_ContainsHeaders(t *testing.T) {
	var buf bytes.Buffer
	statuses := []*SecretStatus{makeStatus("secret/db", 48*time.Hour)}

	if err := FormatTable(&buf, statuses, false); err != nil {
		t.Fatalf("FormatTable error: %v", err)
	}

	out := buf.String()
	for _, header := range []string{"PATH", "EXPIRES AT", "TTL", "STATUS"} {
		if !strings.Contains(out, header) {
			t.Errorf("expected header %q in output", header)
		}
	}
}

func TestFormatTable_StatusLabels(t *testing.T) {
	tests := []struct {
		name   string
		ttl    time.Duration
		wantLabel string
	}{
		{"ok", 72 * time.Hour, "OK"},
		{"warning", 12 * time.Hour, "WARNING"},
		{"expired", -1 * time.Hour, "EXPIRED"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			statuses := []*SecretStatus{makeStatus("secret/test", tc.ttl)}
			if err := FormatTable(&buf, statuses, false); err != nil {
				t.Fatalf("FormatTable error: %v", err)
			}
			if !strings.Contains(buf.String(), tc.wantLabel) {
				t.Errorf("expected label %q in output:\n%s", tc.wantLabel, buf.String())
			}
		})
	}
}

func TestFormatTTL(t *testing.T) {
	tests := []struct {
		d    time.Duration
		want string
	}{
		{0, "expired"},
		{-5 * time.Minute, "expired"},
		{90 * time.Minute, "1h30m"},
		{45 * time.Minute, "45m"},
		{25 * time.Hour, "25h0m"},
	}
	for _, tc := range tests {
		got := formatTTL(tc.d)
		if got != tc.want {
			t.Errorf("formatTTL(%v) = %q, want %q", tc.d, got, tc.want)
		}
	}
}
