package audit

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func fixedNow() time.Time {
	return time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
}

func newTestLogger(buf *bytes.Buffer) *Logger {
	return &Logger{out: buf, now: fixedNow}
}

func TestLog_WritesValidJSON(t *testing.T) {
	var buf bytes.Buffer
	l := newTestLogger(&buf)

	if err := l.Log("check", "secret/db", "warning", "expires in 3d"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var e Entry
	if err := json.Unmarshal(buf.Bytes(), &e); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if e.Event != "check" {
		t.Errorf("event: want %q, got %q", "check", e.Event)
	}
	if e.Path != "secret/db" {
		t.Errorf("path: want %q, got %q", "secret/db", e.Path)
	}
	if e.Status != "warning" {
		t.Errorf("status: want %q, got %q", "warning", e.Status)
	}
	if e.Detail != "expires in 3d" {
		t.Errorf("detail: want %q, got %q", "expires in 3d", e.Detail)
	}
}

func TestLog_TimestampUTC(t *testing.T) {
	var buf bytes.Buffer
	l := newTestLogger(&buf)
	_ = l.Log("rotation", "secret/api", "ok", "")

	var e Entry
	_ = json.Unmarshal(buf.Bytes(), &e)
	if e.Timestamp != fixedNow() {
		t.Errorf("timestamp: want %v, got %v", fixedNow(), e.Timestamp)
	}
}

func TestLog_NewlineDelimited(t *testing.T) {
	var buf bytes.Buffer
	l := newTestLogger(&buf)
	_ = l.Log("check", "secret/a", "ok", "")
	_ = l.Log("check", "secret/b", "expired", "ttl elapsed")

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	for i, line := range lines {
		var e Entry
		if err := json.Unmarshal([]byte(line), &e); err != nil {
			t.Errorf("line %d is not valid JSON: %v", i, err)
		}
	}
}

func TestNewLogger_InvalidPath(t *testing.T) {
	_, err := NewLogger("/nonexistent/dir/audit.log")
	if err == nil {
		t.Fatal("expected error for invalid path, got nil")
	}
}
