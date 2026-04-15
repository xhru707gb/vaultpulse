package rotation_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/vaultpulse/internal/rotation"
)

func makeRotationStatuses() []rotation.Status {
	return []rotation.Status{
		{Path: "secret/ok", LastRotated: time.Now().Add(-24 * time.Hour), Interval: 72 * time.Hour, Overdue: false},
		{Path: "secret/overdue", LastRotated: time.Now().Add(-100 * time.Hour), Interval: 72 * time.Hour, Overdue: true},
		{Path: "secret/also-overdue", LastRotated: time.Now().Add(-200 * time.Hour), Interval: 72 * time.Hour, Overdue: true},
	}
}

func TestBuildRotationAlertPayload_FiltersOK(t *testing.T) {
	statuses := makeRotationStatuses()
	payload, err := rotation.BuildAlertPayload(statuses)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(payload.Alerts) != 2 {
		t.Errorf("expected 2 alerts, got %d", len(payload.Alerts))
	}
	for _, a := range payload.Alerts {
		if !a.Overdue {
			t.Errorf("expected only overdue alerts, got path=%s", a.Path)
		}
	}
}

func TestBuildRotationAlertPayload_AllOK(t *testing.T) {
	statuses := []rotation.Status{
		{Path: "secret/fine", Overdue: false},
	}
	payload, err := rotation.BuildAlertPayload(statuses)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(payload.Alerts) != 0 {
		t.Errorf("expected 0 alerts, got %d", len(payload.Alerts))
	}
}

func TestBuildRotationAlertPayload_TimestampSet(t *testing.T) {
	before := time.Now().UTC()
	payload, _ := rotation.BuildAlertPayload(makeRotationStatuses())
	after := time.Now().UTC()
	if payload.GeneratedAt.Before(before) || payload.GeneratedAt.After(after) {
		t.Errorf("timestamp out of expected range: %v", payload.GeneratedAt)
	}
}

func TestSendRotationAlert_Success(t *testing.T) {
	var received map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	statuses := makeRotationStatuses()
	err := rotation.SendAlert(server.URL, statuses)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["generated_at"] == nil {
		t.Error("expected generated_at in payload")
	}
}

func TestSendRotationAlert_Non2xxReturnsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	err := rotation.SendAlert(server.URL, makeRotationStatuses())
	if err == nil {
		t.Error("expected error for non-2xx response")
	}
}
