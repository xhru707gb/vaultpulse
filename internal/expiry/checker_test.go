package expiry

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/vaultpulse/internal/vault"
)

func newMockVaultServer(t *testing.T, expiration string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/sys/health":
			w.WriteHeader(http.StatusOK)
		default:
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{
					"expiration": expiration,
				},
			})
		}
	}))
}

func newTestChecker(t *testing.T, expiration string, threshold time.Duration) *Checker {
	t.Helper()
	srv := newMockVaultServer(t, expiration)
	t.Cleanup(srv.Close)
	client, err := vault.NewClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("failed to create vault client: %v", err)
	}
	return NewChecker(client, threshold)
}

func TestCheck_NotExpired(t *testing.T) {
	future := time.Now().Add(48 * time.Hour).UTC().Format(time.RFC3339)
	checker := newTestChecker(t, future, 24*time.Hour)

	status, err := checker.Check("secret/my-app/db")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status.IsExpired {
		t.Error("expected secret to not be expired")
	}
	if status.Warning {
		t.Error("expected no warning for secret with 48h TTL and 24h threshold")
	}
}

func TestCheck_Warning(t *testing.T) {
	future := time.Now().Add(12 * time.Hour).UTC().Format(time.RFC3339)
	checker := newTestChecker(t, future, 24*time.Hour)

	status, err := checker.Check("secret/my-app/db")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status.IsExpired {
		t.Error("expected secret to not be expired")
	}
	if !status.Warning {
		t.Error("expected warning for secret within threshold")
	}
}

func TestCheck_Expired(t *testing.T) {
	past := time.Now().Add(-1 * time.Hour).UTC().Format(time.RFC3339)
	checker := newTestChecker(t, past, 24*time.Hour)

	status, err := checker.Check("secret/my-app/db")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !status.IsExpired {
		t.Error("expected secret to be expired")
	}
}

func TestCheck_InvalidExpiration(t *testing.T) {
	checker := newTestChecker(t, "not-a-valid-timestamp", 24*time.Hour)

	_, err := checker.Check("secret/my-app/db")
	if err == nil {
		t.Error("expected error for invalid expiration format, got nil")
	}
}

func TestCheckAll_MultipleSecrets(t *testing.T) {
	future := time.Now().Add(72 * time.Hour).UTC().Format(time.RFC3339)
	checker := newTestChecker(t, future, 24*time.Hour)

	paths := []string{"secret/a", "secret/b", "secret/c"}
	statuses, err := checker.CheckAll(paths)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(statuses) != len(paths) {
		t.Errorf("expected %d statuses, got %d", len(paths), len(statuses))
	}
}
