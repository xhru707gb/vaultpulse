package audit

import (
	"bytes"
	"strings"
	"testing"
)

const sampleLog = `{"timestamp":"2024-05-01T10:00:00Z","path":"secret/db/password","event":"check","status":"expired","ttl":"0s"}
{"timestamp":"2024-05-01T10:01:00Z","path":"secret/api/key","event":"check","status":"warning","ttl":"2h30m"}
{"timestamp":"2024-05-01T10:02:00Z","path":"secret/tls/cert","event":"rotation","status":"ok","ttl":""}
`

func TestReport_ContainsHeaders(t *testing.T) {
	var buf bytes.Buffer
	err := Report(strings.NewReader(sampleLog), &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, header := range []string{"TIMESTAMP", "PATH", "EVENT", "STATUS", "TTL"} {
		if !strings.Contains(out, header) {
			t.Errorf("expected header %q in output", header)
		}
	}
}

func TestReport_ContainsEntries(t *testing.T) {
	var buf bytes.Buffer
	_ = Report(strings.NewReader(sampleLog), &buf)
	out := buf.String()

	if !strings.Contains(out, "secret/db/password") {
		t.Error("expected path secret/db/password in output")
	}
	if !strings.Contains(out, "expired") {
		t.Error("expected status 'expired' in output")
	}
	if !strings.Contains(out, "warning") {
		t.Error("expected status 'warning' in output")
	}
}

func TestReport_CountLine(t *testing.T) {
	var buf bytes.Buffer
	_ = Report(strings.NewReader(sampleLog), &buf)
	out := buf.String()
	if !strings.Contains(out, "3 audit record(s) shown.") {
		t.Errorf("expected count line, got:\n%s", out)
	}
}

func TestReport_EmptyInput(t *testing.T) {
	var buf bytes.Buffer
	_ = Report(strings.NewReader(""), &buf)
	out := buf.String()
	if !strings.Contains(out, "0 audit record(s) shown.") {
		t.Errorf("expected zero count, got:\n%s", out)
	}
}

func TestReport_SkipsInvalidJSON(t *testing.T) {
	input := "{not-valid}\n" + `{"timestamp":"2024-05-01T10:00:00Z","path":"p","event":"check","status":"ok"}` + "\n"
	var buf bytes.Buffer
	_ = Report(strings.NewReader(input), &buf)
	out := buf.String()
	if !strings.Contains(out, "1 audit record(s) shown.") {
		t.Errorf("expected 1 record after skipping invalid JSON, got:\n%s", out)
	}
}
