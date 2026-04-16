package watch_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/your-org/vaultpulse/internal/expiry"
	"github.com/your-org/vaultpulse/internal/vault"
	"github.com/your-org/vaultpulse/internal/watch"
)

func newMockVault(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":{"ttl":3600,"version":1}}`))
	}))
}

func newTestChecker(t *testing.T, srv *httptest.Server) *expiry.Checker {
	t.Helper()
	client, err := vault.NewClient(srv.URL, "test-token")
	if err != nil {
		t.Fatal(err)
	}
	return expiry.NewChecker(client, 300)
}

func TestNew_InvalidInterval(t *testing.T) {
	_, err := watch.New(nil, 0, func(watch.Event) {})
	if err == nil {
		t.Fatal("expected error for zero interval")
	}
}

func TestNew_NilHandler(t *testing.T) {
	_, err := watch.New(nil, time.Second, nil)
	if err == nil {
		t.Fatal("expected error for nil handler")
	}
}

func TestRun_CallsHandlerOnTick(t *testing.T) {
	srv := newMockVault(t)
	defer srv.Close()
	checker := newTestChecker(t, srv)

	var calls atomic.Int32
	w, err := watch.New(checker, 50*time.Millisecond, func(e watch.Event) {
		calls.Add(1)
	})
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Millisecond)
	defer cancel()
	w.Run(ctx, []string{"secret/data/test"})

	if calls.Load() < 2 {
		t.Errorf("expected at least 2 handler calls, got %d", calls.Load())
	}
}

func TestRun_StopsOnContextCancel(t *testing.T) {
	srv := newMockVault(t)
	defer srv.Close()
	checker := newTestChecker(t, srv)

	var calls atomic.Int32
	w, _ := watch.New(checker, 10*time.Millisecond, func(e watch.Event) {
		calls.Add(1)
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately
	w.Run(ctx, []string{"secret/data/test"})

	if calls.Load() != 0 {
		t.Errorf("expected 0 handler calls after immediate cancel, got %d", calls.Load())
	}
}
