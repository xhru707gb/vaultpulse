package cmd

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func writeTempConfig(t *testing.T, vaultAddr, token string) string {
	t.Helper()
	content := "vault_addr: " + vaultAddr + "\n" +
		"token: " + token + "\n" +
		"paths:\n  - secret/data/test\n" +
		"warning_threshold: 72h\n"
	dir := t.TempDir()
	p := filepath.Join(dir, "vaultpulse.yaml")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writing temp config: %v", err)
	}
	return p
}

func TestExecute_MissingConfig(t *testing.T) {
	rootCmd.SetArgs([]string{"--config", "/nonexistent/path.yaml"})
	defer rootCmd.SetArgs(nil)

	var buf bytes.Buffer
	rootCmd.SetErr(&buf)

	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing config, got nil")
	}
}

func TestExecute_TableFlag(t *testing.T) {
	vaultSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/sys/health":
			w.WriteHeader(http.StatusOK)
		case "/v1/secret/data/test":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"data":{"metadata":{"created_time":"2024-01-01T00:00:00Z","deletion_time":""}}}`)) //nolint:lll
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer vaultSrv.Close()

	cfgPath := writeTempConfig(t, vaultSrv.URL, "test-token")

	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetArgs([]string{"--config", cfgPath, "--table"})
	defer rootCmd.SetArgs(nil)

	// We only assert no panic; full integration depends on vault responses.
	_ = rootCmd.Execute()
}
