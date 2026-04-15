package alert_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/your-org/vaultpulse/internal/alert"
	"github.com/your-org/vaultpulse/internal/expiry"
)

func makeStatuses() []expiry.Status {
	return []expiry.Status{
		{Path: "secret/ok", TTL: 72 * time.Hour, State: expiry.StateOK},
		{Path: "secret/warn", TTL: 20 * time.Hour, State: expiry.StateWarning},
		{Path: "secret/gone", TTL: 0, State: expiry.StateExpired},
	}
}

func TestNotify_SendsOnlyAlerts(t *testing.T) {
	var received alert.WebhookPayload

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := alert.NewNotifier(ts.URL)
	if err := n.Notify(makeStatuses()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(received.Alerts) != 2 {
		t.Errorf("expected 2 alerts, got %d", len(received.Alerts))
	}
	if received.Summary == "" {
		t.Error("expected non-empty summary")
	}
}

func TestNotify_SendsOnlyAlerts_AlertPaths(t *testing.T) {
	var received alert.WebhookPayload

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := alert.NewNotifier(ts.URL)
	if err := n.Notify(makeStatuses()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Only warning and expired secrets should appear in the payload.
	paths := make(map[string]bool, len(received.Alerts))
	for _, a := range received.Alerts {
		paths[a.Path] = true
	}
	if !paths["secret/warn"] {
		t.Error("expected secret/warn in alerts")
	}
	if !paths["secret/gone"] {
		t.Error("expected secret/gone in alerts")
	}
	if paths["secret/ok"] {
		t.Error("secret/ok should not appear in alerts")
	}
}

func TestNotify_NoAlertsSkipsWebhook(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := alert.NewNotifier(ts.URL)
	statuses := []expiry.Status{
		{Path: "secret/ok", TTL: 48 * time.Hour, State: expiry.StateOK},
	}
	if err := n.Notify(statuses); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("webhook should not be called when there are no alerts")
	}
}

func TestNotify_Non2xxReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := alert.NewNotifier(ts.URL)
	statuses := []expiry.Status{
		{Path: "secret/expired", TTL: 0, State: expiry.StateExpired},
	}
	if err := n.Notify(statuses); err == nil {
		t.Error("expected error for non-2xx response")
	}
}
