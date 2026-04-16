package cache_test

import (
	"errors"
	"testing"
	"time"

	"github.com/vaultpulse/vaultpulse/internal/cache"
)

func TestGetOrFetch_CacheMiss_CallsFetcher(t *testing.T) {
	c := cache.New(time.Minute)
	calls := 0
	fetcher := func(key string) (string, error) {
		calls++
		return "value-" + key, nil
	}

	val, err := cache.GetOrFetch(c, "foo", time.Minute, fetcher)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "value-foo" {
		t.Errorf("expected value-foo, got %q", val)
	}
	if calls != 1 {
		t.Errorf("expected 1 fetch call, got %d", calls)
	}
}

func TestGetOrFetch_CacheHit_DoesNotCallFetcher(t *testing.T) {
	c := cache.New(time.Minute)
	c.Set("bar", "cached-value")
	calls := 0
	fetcher := func(key string) (string, error) {
		calls++
		return "fresh", nil
	}

	val, err := cache.GetOrFetch(c, "bar", time.Minute, fetcher)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "cached-value" {
		t.Errorf("expected cached-value, got %q", val)
	}
	if calls != 0 {
		t.Errorf("expected 0 fetch calls, got %d", calls)
	}
}

func TestGetOrFetch_FetcherError_NotCached(t *testing.T) {
	c := cache.New(time.Minute)
	fetchErr := errors.New("upstream down")
	fetcher := func(key string) (string, error) {
		return "", fetchErr
	}

	_, err := cache.GetOrFetch(c, "baz", time.Minute, fetcher)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, fetchErr) {
		t.Errorf("expected fetchErr, got %v", err)
	}
	if _, ok := c.Get("baz"); ok {
		t.Error("expected key not to be cached after fetch error")
	}
}

func TestInvalidate_RemovesMatchingKeys(t *testing.T) {
	c := cache.New(time.Minute)
	c.Set("secret/a", "v1")
	c.Set("secret/b", "v2")
	c.Set("health/x", "v3")

	removed := cache.Invalidate(c, "secret/")
	if removed != 2 {
		t.Errorf("expected 2 removed, got %d", removed)
	}
	if _, ok := c.Get("health/x"); !ok {
		t.Error("expected health/x to remain in cache")
	}
}

func TestInvalidate_NoMatch_ReturnsZero(t *testing.T) {
	c := cache.New(time.Minute)
	c.Set("foo", "bar")

	removed := cache.Invalidate(c, "nonexistent/")
	if removed != 0 {
		t.Errorf("expected 0 removed, got %d", removed)
	}
}
