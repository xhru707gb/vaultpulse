package health_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/user/vaultpulse/internal/health"
)

func newMockVaultServer(t *testing.T, sealed bool) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/sys/health" {
			code := http.StatusOK
			if sealed {
				code = http.StatusServiceUnavailable
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(code)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"initialized": true,
				"sealed":      sealed,
				"standby":     false,
			})
		}
	}))
}

func newTestChecker(t *testing.T, srv *httptest.Server) *health.Checker {
	t.Helper()
	cfg := api.DefaultConfig()
	cfg.Address = srv.URL
	client, err := api.NewClient(cfg)
	if err != nil {
		t.Fatalf("api.NewClient: %v", err)
	}
	return health.NewChecker(client)
}

func TestCheck_Healthy(t *testing.T) {
	srv := newMockVaultServer(t, false)
	defer srv.Close()

	checker := newTestChecker(t, srv)
	s := checker.Check(context.Background())

	if s.Error != nil {
		t.Fatalf("unexpected error: %v", s.Error)
	}
	if !s.Healthy() {
		t.Errorf("expected healthy status, got sealed=%v standby=%v", s.Sealed, s.Standby)
	}
}

func TestCheck_Sealed(t *testing.T) {
	srv := newMockVaultServer(t, true)
	defer srv.Close()

	checker := newTestChecker(t, srv)
	s := checker.Check(context.Background())

	if s.Healthy() {
		t.Error("expected unhealthy status for sealed vault")
	}
	if !s.Sealed {
		t.Error("expected Sealed=true")
	}
}

func TestCheck_Latency(t *testing.T) {
	srv := newMockVaultServer(t, false)
	defer srv.Close()

	checker := newTestChecker(t, srv)
	s := checker.Check(context.Background())

	if s.Latency < 0 {
		t.Errorf("negative latency: %v", s.Latency)
	}
}

func TestCheck_CancelledContext(t *testing.T) {
	srv := newMockVaultServer(t, false)
	defer srv.Close()

	checker := newTestChecker(t, srv)
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately before making the request

	s := checker.Check(ctx)
	if s.Error == nil {
		t.Error("expected error for cancelled context, got nil")
	}
	if s.Healthy() {
		t.Error("expected unhealthy status for cancelled context")
	}
}
