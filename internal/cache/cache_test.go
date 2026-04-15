package cache_test

import (
	"testing"
	"time"

	"github.com/yourusername/vaultpulse/internal/cache"
)

func TestSet_And_Get_Hit(t *testing.T) {
	c := cache.New(5 * time.Second)
	c.Set("foo", "bar")

	v, ok := c.Get("foo")
	if !ok {
		t.Fatal("expected cache hit")
	}
	if v.(string) != "bar" {
		t.Fatalf("expected bar, got %v", v)
	}
}

func TestGet_Miss_UnknownKey(t *testing.T) {
	c := cache.New(5 * time.Second)

	_, ok := c.Get("missing")
	if ok {
		t.Fatal("expected cache miss for unknown key")
	}
}

func TestGet_Miss_AfterExpiry(t *testing.T) {
	c := cache.New(1 * time.Millisecond)
	c.Set("key", 42)

	time.Sleep(5 * time.Millisecond)

	_, ok := c.Get("key")
	if ok {
		t.Fatal("expected cache miss after TTL expiry")
	}
}

func TestSetWithTTL_OverridesDefault(t *testing.T) {
	c := cache.New(1 * time.Millisecond)
	c.SetWithTTL("long", "value", 10*time.Second)

	time.Sleep(5 * time.Millisecond)

	_, ok := c.Get("long")
	if !ok {
		t.Fatal("expected cache hit with extended TTL")
	}
}

func TestDelete_RemovesEntry(t *testing.T) {
	c := cache.New(5 * time.Second)
	c.Set("del", true)
	c.Delete("del")

	_, ok := c.Get("del")
	if ok {
		t.Fatal("expected entry to be deleted")
	}
}

func TestFlush_ClearsAll(t *testing.T) {
	c := cache.New(5 * time.Second)
	c.Set("a", 1)
	c.Set("b", 2)
	c.Flush()

	if c.Len() != 0 {
		t.Fatalf("expected 0 entries after flush, got %d", c.Len())
	}
}

func TestLen_ReturnsCount(t *testing.T) {
	c := cache.New(5 * time.Second)
	c.Set("x", 1)
	c.Set("y", 2)

	if c.Len() != 2 {
		t.Fatalf("expected 2, got %d", c.Len())
	}
}
