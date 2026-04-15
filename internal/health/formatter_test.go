package health_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/user/vaultpulse/internal/health"
)

func makeStatus(sealed, standby bool, err error) health.Status {
	return health.Status{
		Initialized: true,
		Sealed:      sealed,
		Standby:     standby,
		CheckedAt:   time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
		Latency:     5 * time.Millisecond,
		Error:       err,
	}
}

func TestFormatTable_ContainsHeaders(t *testing.T) {
	var buf bytes.Buffer
	health.FormatTable(&buf, makeStatus(false, false, nil))
	out := buf.String()
	for _, hdr := range []string{"STATUS", "INITIALIZED", "SEALED", "STANDBY", "LATENCY", "CHECKED AT"} {
		if !strings.Contains(out, hdr) {
			t.Errorf("missing header %q in output", hdr)
		}
	}
}

func TestFormatTable_StatusLabels(t *testing.T) {
	cases := []struct {
		name   string
		s      health.Status
		wantLbl string
	}{
		{"ok", makeStatus(false, false, nil), "OK"},
		{"sealed", makeStatus(true, false, nil), "SEALED"},
		{"standby", makeStatus(false, true, nil), "STANDBY"},
		{"error", makeStatus(false, false, errors.New("boom")), "ERROR"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			health.FormatTable(&buf, tc.s)
			if !strings.Contains(buf.String(), tc.wantLbl) {
				t.Errorf("expected label %q in:\n%s", tc.wantLbl, buf.String())
			}
		})
	}
}

func TestFormatTable_LatencyShown(t *testing.T) {
	var buf bytes.Buffer
	health.FormatTable(&buf, makeStatus(false, false, nil))
	if !strings.Contains(buf.String(), "5ms") {
		t.Errorf("expected latency '5ms' in output: %s", buf.String())
	}
}
