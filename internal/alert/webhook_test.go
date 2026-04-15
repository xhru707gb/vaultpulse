package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yourusername/vaultpulse/internal/expiry"
)

func makeStatuses2() []expiry.SecretStatus {
	return []expiry.SecretStatus{
		{Path: "secret/ok", State: expiry.StateOK, TTL: 86400},
		{Path: "secret/warn", State: expiry.StateWarning, TTL: 3600},
		{Path: "secret/expired", State: expiry.StateExpired, TTL: 0},
	}
}

func TestBuildPayload_FiltersOK(t *testing.T) {
	statuses := makeStatuses2()
	payload := buildPayload(statuses)

	if len(payload.Alerts) != 2 {
		t.Fatalf("expected 2 alerts, got %d", len(payload.Alerts))
	}
	if payload.Alerts[0].Path != "secret/warn" {
		t.Errorf("expected first alert path 'secret/warn', got %q", payload.Alerts[0].Path)
	}
	if payload.Alerts[1].Status != string(expiry.StateExpired) {
		t.Errorf("expected second alert status 'expired', got %q", payload.Alerts[1].Status)
	}
}

func TestBuildPayload_TimestampSet(t *testing.T) {
	payload := buildPayload(makeStatuses2())
	if payload.Timestamp == "" {
		t.Error("expected non-empty timestamp")
	}
}

func TestSendWebhook_Success(t *testing.T) {
	var received WebhookPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	payload := buildPayload(makeStatuses2())
	if err := sendWebhook(ts.URL, payload); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(received.Alerts) != 2 {
		t.Errorf("expected 2 alerts in received payload, got %d", len(received.Alerts))
	}
}

func TestSendWebhook_Non2xxError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	payload := buildPayload(makeStatuses2())
	err := sendWebhook(ts.URL, payload)
	if err == nil {
		t.Fatal("expected error for non-2xx response, got nil")
	}
}

func TestFormatTTLDuration(t *testing.T) {
	cases := []struct {
		seconds int64
		want    string
	}{
		{0, "expired"},
		{-10, "expired"},
		{1800, "30m"},
		{3661, "1h1m"},
	}
	for _, c := range cases {
		got := formatTTLDuration(c.seconds)
		if got != c.want {
			t.Errorf("formatTTLDuration(%d) = %q, want %q", c.seconds, got, c.want)
		}
	}
}
