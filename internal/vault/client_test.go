package vault

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func newMockVaultServer(t *testing.T) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()

	// Health endpoint
	mux.HandleFunc("/v1/sys/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"initialized": true, "sealed": false})
	})

	// Secret endpoint
	mux.HandleFunc("/v1/secret/data/myapp/db", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"lease_duration": 3600,
			"renewable":      true,
			"data":           map[string]interface{}{"password": "s3cr3t"},
		})
	})

	return httptest.NewServer(mux)
}

func TestNewClient(t *testing.T) {
	_, err := NewClient("http://127.0.0.1:8200", "test-token")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestPing_Success(t *testing.T) {
	srv := newMockVaultServer(t)
	defer srv.Close()

	client, err := NewClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	if err := client.Ping(context.Background()); err != nil {
		t.Errorf("expected ping to succeed, got: %v", err)
	}
}

func TestGetSecretMeta_Success(t *testing.T) {
	srv := newMockVaultServer(t)
	defer srv.Close()

	client, err := NewClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	meta, err := client.GetSecretMeta(context.Background(), "secret/data/myapp/db")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if meta.Path != "secret/data/myapp/db" {
		t.Errorf("unexpected path: %s", meta.Path)
	}
	if meta.TTL != 3600*time.Second {
		t.Errorf("unexpected TTL: %v", meta.TTL)
	}
	if !meta.Renewable {
		t.Error("expected renewable to be true")
	}
	if meta.Expiration.Before(time.Now()) {
		t.Error("expiration should be in the future")
	}
}
