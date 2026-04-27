package secretexpiry_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/vaultpulse/vaultpulse/internal/secretexpiry"
)

func makeExpiryStatuses() []secretexpiry.Status {
	now := time.Now().UTC()
	return []secretexpiry.Status{
		{Path: "secret/ok", State: secretexpiry.StateOK, Remaining: 72 * time.Hour},
		{Path: "secret/warn", State: secretexpiry.StateWarning, Remaining: 6 * time.Hour},
		{Path: "secret/expired", State: secretexpiry.StateExpired, Remaining: 0},
		{Path: "secret/ok2", State: secretexpiry.StateOK, Remaining: 48 * time.Hour},
		_ = now
	}
}

func TestBuildExpiryAlertPayload_FiltersOK(t *testing.T) {
	statuses := makeExpiryStatuses()
	payload := secretexpiry.BuildAlertPayload(statuses)
	if payload.AlertCount != 2 {
		t.Fatalf("expected 2 alerts, got %d", payload.AlertCount)
	}
	for _, e := range payload.Entries {
		if e.State == "OK" {
			t.Errorf("OK entry should not appear in alert payload: %s", e.Path)
		}
	}
}

func TestBuildExpiryAlertPayload_AllOK(t *testing.T) {
	statuses := []secretexpiry.Status{
		{Path: "secret/a", State: secretexpiry.StateOK, Remaining: 24 * time.Hour},
		{Path: "secret/b", State: secretexpiry.StateOK, Remaining: 48 * time.Hour},
	}
	payload := secretexpiry.BuildAlertPayload(statuses)
	if payload.AlertCount != 0 {
		t.Fatalf("expected 0 alerts, got %d", payload.AlertCount)
	}
	if len(payload.Entries) != 0 {
		t.Errorf("expected empty entries slice")
	}
}

func TestBuildExpiryAlertPayload_TimestampSet(t *testing.T) {
	before := time.Now().UTC().Add(-time.Second)
	payload := secretexpiry.BuildAlertPayload(makeExpiryStatuses())
	if payload.GeneratedAt.Before(before) {
		t.Errorf("GeneratedAt should be recent, got %v", payload.GeneratedAt)
	}
}

func TestSendExpiryAlert_Success(t *testing.T) {
	var received secretexpiry.AlertPayload
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	payload := secretexpiry.BuildAlertPayload(makeExpiryStatuses())
	if err := secretexpiry.SendAlert(srv.URL, payload); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.AlertCount != payload.AlertCount {
		t.Errorf("expected alert_count %d, got %d", payload.AlertCount, received.AlertCount)
	}
}

func TestSendExpiryAlert_Non2xxReturnsError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	payload := secretexpiry.BuildAlertPayload(makeExpiryStatuses())
	if err := secretexpiry.SendAlert(srv.URL, payload); err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}
