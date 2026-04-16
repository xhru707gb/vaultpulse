package filter_test

import (
	"testing"

	"github.com/vaultpulse/vaultpulse/internal/filter"
)

var paths = []string{
	"secret/prod/db",
	"secret/prod/api",
	"secret/staging/db",
	"kv/prod/token",
}

func TestFilter_Prefix(t *testing.T) {
	got := filter.Filter(paths, filter.Options{Prefix: "secret/prod"})
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilter_Contains(t *testing.T) {
	got := filter.Filter(paths, filter.Options{Contains: "db"})
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilter_Exclude(t *testing.T) {
	got := filter.Filter(paths, filter.Options{Exclude: "staging"})
	if len(got) != 3 {
		t.Fatalf("expected 3, got %d", len(got))
	}
}

func TestFilter_Combined(t *testing.T) {
	got := filter.Filter(paths, filter.Options{Prefix: "secret", Exclude: "staging"})
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilter_Empty(t *testing.T) {
	got := filter.Filter(paths, filter.Options{})
	if len(got) != len(paths) {
		t.Fatalf("expected %d, got %d", len(paths), len(got))
	}
}

func TestMatchesAny_True(t *testing.T) {
	if !filter.MatchesAny("secret/prod/db", []string{"kv/", "secret/prod"}) {
		t.Fatal("expected match")
	}
}

func TestMatchesAny_False(t *testing.T) {
	if filter.MatchesAny("secret/prod/db", []string{"kv/", "pki/"}) {
		t.Fatal("expected no match")
	}
}
