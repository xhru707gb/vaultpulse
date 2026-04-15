package audit

import (
	"strings"
	"testing"
)

const sampleLog = `{"timestamp":"2024-01-01T00:00:00Z","event":"expiry_check","path":"secret/a","outcome":"ok"}
{"timestamp":"2024-01-01T00:01:00Z","event":"expiry_check","path":"secret/b","outcome":"warning"}
{"timestamp":"2024-01-01T00:02:00Z","event":"expiry_check","path":"secret/c","outcome":"ok"}
{"timestamp":"2024-01-01T00:03:00Z","event":"rotation_check","path":"secret/d","outcome":"overdue"}
{"timestamp":"2024-01-01T00:04:00Z","event":"expiry_check","path":"secret/e","outcome":"warning"}
`

func TestSummary_CountsAreCorrect(t *testing.T) {
	counts, err := summariseReader(strings.NewReader(sampleLog))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(counts) != 3 {
		t.Fatalf("expected 3 distinct event+outcome pairs, got %d", len(counts))
	}

	want := map[string]int{
		"expiry_check:ok":       2,
		"expiry_check:warning":  2,
		"rotation_check:overdue": 1,
	}
	for _, ec := range counts {
		key := ec.Event + ":" + ec.Outcome
		if want[key] != ec.Count {
			t.Errorf("key %q: want count %d, got %d", key, want[key], ec.Count)
		}
	}
}

func TestSummary_SortedOrder(t *testing.T) {
	counts, err := summariseReader(strings.NewReader(sampleLog))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// First two entries should be expiry_check (alphabetically before rotation_check)
	for i := 0; i < 2; i++ {
		if counts[i].Event != "expiry_check" {
			t.Errorf("position %d: expected expiry_check, got %s", i, counts[i].Event)
		}
	}
	if counts[2].Event != "rotation_check" {
		t.Errorf("position 2: expected rotation_check, got %s", counts[2].Event)
	}
}

func TestSummary_EmptyInput(t *testing.T) {
	counts, err := summariseReader(strings.NewReader(""))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(counts) != 0 {
		t.Errorf("expected empty result, got %v", counts)
	}
}

func TestSummary_InvalidJSON(t *testing.T) {
	_, err := summariseReader(strings.NewReader("not-json\n"))
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}
