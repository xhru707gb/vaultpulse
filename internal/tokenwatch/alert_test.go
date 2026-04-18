package tokenwatch_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/your-org/vaultpulse/internal/tokenwatch"
)

func makeTokenStatuses() []tokenwatch.Status {
	now := time.Now()
	return []tokenwatch.Status{
		{Path: "auth/token/ok", State: tokenwatch.StateOK, Remaining: 48 * time.Hour},
		{Path: "auth/token/warn", State: tokenwatch.StateWarning, Remaining: 2 * time.Hour},
		{Path: "auth/token/crit", State: tokenwatch.StateCritical, Remaining: 10 * time.Minute},
		{Path: "auth/token/exp", State: tokenwatch.StateExpired, Remaining: 0},
		_ = now
	}
}

func TestBuildAlertPayload_FiltersOK(t *testing.T) {
	statuses := makeTokenStatuses()
	payload := tokenwatch.BuildAlertPayload(statuses, time.Now())
	if len(payload.Alerts) != 3 {
		t.Fatalf("expected 3 alerts, got %d", len(payload.Alerts))
	}
	for _, a := range payload.Alerts {
		if a.Path == "auth/token/ok" {
			t.Errorf("OK token should not appear in alerts")
		}
	}
}

func TestBuildAlertPayload_AllOK(t *testing.T) {
	statuses := []tokenwatch.Status{
		{Path: "auth/token/a", State: tokenwatch.StateOK, Remaining: 72 * time.Hour},
	}
	payload := tokenwatch.BuildAlertPayload(statuses, time.Now())
	if len(payload.Alerts) != 0 {
		t.Errorf("expected no alerts, got %d", len(payload.Alerts))
	}
}

func TestBuildAlertPayload_TimestampSet(t *testing.T) {
	fixed := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	payload := tokenwatch.BuildAlertPayload(makeTokenStatuses(), fixed)
	if !payload.Timestamp.Equal(fixed) {
		t.Errorf("expected timestamp %v, got %v", fixed, payload.Timestamp)
	}
}

func TestSendAlert_Success(t *testing.T) {
	var received tokenwatch.AlertPayload
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	payload := tokenwatch.BuildAlertPayload(makeTokenStatuses(), time.Now())
	if err := tokenwatch.SendAlert(srv.URL, payload); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(received.Alerts) != 3 {
		t.Errorf("expected 3 alerts delivered, got %d", len(received.Alerts))
	}
}

func TestSendAlert_Non2xx(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	payload := tokenwatch.BuildAlertPayload(makeTokenStatuses(), time.Now())
	if err := tokenwatch.SendAlert(srv.URL, payload); err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}
