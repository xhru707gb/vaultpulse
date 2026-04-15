package health

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func makeHealthStatuses() []Status {
	return []Status{
		{Node: "vault-1", Healthy: true, Sealed: false, LatencyMs: 12},
		{Node: "vault-2", Healthy: false, Sealed: false, LatencyMs: 350},
		{Node: "vault-3", Healthy: false, Sealed: true, LatencyMs: 0},
	}
}

func TestBuildAlertPayload_FiltersHealthy(t *testing.T) {
	statuses := makeHealthStatuses()
	payload := BuildAlertPayload(statuses)
	if payload == nil {
		t.Fatal("expected non-nil payload")
	}
	if len(payload.Alerts) != 2 {
		t.Fatalf("expected 2 alerts, got %d", len(payload.Alerts))
	}
	for _, a := range payload.Alerts {
		if a.Node == "vault-1" {
			t.Errorf("healthy node vault-1 should not appear in alerts")
		}
	}
}

func TestBuildAlertPayload_AllHealthy(t *testing.T) {
	statuses := []Status{
		{Node: "vault-1", Healthy: true, Sealed: false, LatencyMs: 5},
	}
	payload := BuildAlertPayload(statuses)
	if payload != nil {
		t.Errorf("expected nil payload when all nodes healthy, got %+v", payload)
	}
}

func TestBuildAlertPayload_TimestampSet(t *testing.T) {
	statuses := makeHealthStatuses()
	payload := BuildAlertPayload(statuses)
	if payload == nil {
		t.Fatal("expected non-nil payload")
	}
	if payload.Timestamp == "" {
		t.Error("expected Timestamp to be set")
	}
}

func TestSendAlert_Success(t *testing.T) {
	var received AlertPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	payload := BuildAlertPayload(makeHealthStatuses())
	if err := SendAlert(ts.URL, payload); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(received.Alerts) != 2 {
		t.Errorf("expected 2 alerts in received payload, got %d", len(received.Alerts))
	}
}

func TestSendAlert_Non2xxReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	payload := BuildAlertPayload(makeHealthStatuses())
	if err := SendAlert(ts.URL, payload); err == nil {
		t.Error("expected error for non-2xx response")
	}
}

func TestSendAlert_NilPayloadNoOp(t *testing.T) {
	if err := SendAlert("http://unused", nil); err != nil {
		t.Errorf("expected no error for nil payload, got %v", err)
	}
}
