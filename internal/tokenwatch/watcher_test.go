package tokenwatch_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/yourusername/vaultpulse/internal/tokenwatch"
)

type mockFetcher struct {
	tokens []tokenwatch.TokenInfo
	err    error
}

func (m *mockFetcher) LookupTokens(_ context.Context) ([]tokenwatch.TokenInfo, error) {
	return m.tokens, m.err
}

func fixedNow() time.Time {
	return time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
}

func newTestWatcher(tokens []tokenwatch.TokenInfo, err error) *tokenwatch.Watcher {
	w, _ := tokenwatch.NewWithClock(&mockFetcher{tokens: tokens, err: err}, 30*time.Minute, fixedNow)
	return w
}

func TestEvaluate_OK(t *testing.T) {
	now := fixedNow()
	tokens := []tokenwatch.TokenInfo{
		{Accessor: "abc", DisplayName: "root", ExpireTime: now.Add(2 * time.Hour)},
	}
	w := newTestWatcher(tokens, nil)
	statuses, err := w.Evaluate(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if statuses[0].State != "ok" {
		t.Errorf("expected ok, got %s", statuses[0].State)
	}
}

func TestEvaluate_Warning(t *testing.T) {
	now := fixedNow()
	tokens := []tokenwatch.TokenInfo{
		{Accessor: "abc", DisplayName: "ci", ExpireTime: now.Add(10 * time.Minute)},
	}
	w := newTestWatcher(tokens, nil)
	statuses, _ := w.Evaluate(context.Background())
	if statuses[0].State != "warning" {
		t.Errorf("expected warning, got %s", statuses[0].State)
	}
}

func TestEvaluate_Expired(t *testing.T) {
	now := fixedNow()
	tokens := []tokenwatch.TokenInfo{
		{Accessor: "abc", DisplayName: "old", ExpireTime: now.Add(-1 * time.Minute)},
	}
	w := newTestWatcher(tokens, nil)
	statuses, _ := w.Evaluate(context.Background())
	if statuses[0].State != "expired" {
		t.Errorf("expected expired, got %s", statuses[0].State)
	}
}

func TestEvaluate_FetcherError(t *testing.T) {
	w := newTestWatcher(nil, errors.New("vault down"))
	_, err := w.Evaluate(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestNew_NilFetcher(t *testing.T) {
	_, err := tokenwatch.New(nil, 30*time.Minute)
	if err == nil {
		t.Fatal("expected error for nil fetcher")
	}
}
