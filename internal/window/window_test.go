package window_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/yourusername/vaultpulse/internal/window"
)

var epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func fixedNow(t time.Time) func() time.Time { return func() time.Time { return t } }

func TestNew_InvalidDuration(t *testing.T) {
	_, err := window.New[string](0)
	if err == nil {
		t.Fatal("expected error for zero duration")
	}
}

func TestAdd_And_Len(t *testing.T) {
	w, _ := window.New[string](time.Minute)
	w.SetNow(fixedNow(epoch))
	w.Add("a")
	w.Add("b")
	if w.Len() != 2 {
		t.Fatalf("expected 2, got %d", w.Len())
	}
}

func TestPrune_RemovesStaleEntries(t *testing.T) {
	w, _ := window.New[string](time.Minute)
	w.SetNow(fixedNow(epoch))
	w.Add("old")
	w.SetNow(fixedNow(epoch.Add(2 * time.Minute)))
	w.Add("new")
	if w.Len() != 1 {
		t.Fatalf("expected 1 after prune, got %d", w.Len())
	}
	if w.Entries()[0].Value != "new" {
		t.Fatal("expected remaining entry to be 'new'")
	}
}

func TestReset_ClearsEntries(t *testing.T) {
	w, _ := window.New[int](time.Minute)
	w.SetNow(fixedNow(epoch))
	w.Add(1)
	w.Add(2)
	w.Reset()
	if w.Len() != 0 {
		t.Fatal("expected 0 after reset")
	}
}

func TestFormatTable_ContainsHeaders(t *testing.T) {
	w, _ := window.New[string](time.Minute)
	w.SetNow(fixedNow(epoch))
	w.Add("secret/foo")
	out := window.FormatTable(w.Entries(), func(s string) string { return s })
	for _, hdr := range []string{"VALUE", "RECORDED AT"} {
		if !strings.Contains(out, hdr) {
			t.Errorf("missing header %q", hdr)
		}
	}
}

func TestFormatSummary_ShowsCount(t *testing.T) {
	w, _ := window.New[int](30 * time.Second)
	w.SetNow(fixedNow(epoch))
	for i := 0; i < 5; i++ {
		w.Add(i)
	}
	summary := window.FormatSummary(w.Entries(), 30*time.Second)
	if !strings.Contains(summary, fmt.Sprintf("%d entries", 5)) {
		t.Errorf("unexpected summary: %s", summary)
	}
}
