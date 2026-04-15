package audit

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
)

// EventCount holds the number of occurrences for a single event+outcome pair.
type EventCount struct {
	Event   string
	Outcome string
	Count   int
}

// Summary reads a JSONL audit log file and returns aggregated event counts
// sorted by event name then outcome.
func Summary(path string) ([]EventCount, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("audit summary: open %s: %w", path, err)
	}
	defer f.Close()
	return summariseReader(f)
}

func summariseReader(r io.Reader) ([]EventCount, error) {
	type key struct{ event, outcome string }
	counts := make(map[key]int)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var entry Entry
		if err := json.Unmarshal(line, &entry); err != nil {
			return nil, fmt.Errorf("audit summary: parse line: %w", err)
		}
		counts[key{entry.Event, entry.Outcome}]++
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("audit summary: scan: %w", err)
	}

	result := make([]EventCount, 0, len(counts))
	for k, n := range counts {
		result = append(result, EventCount{Event: k.event, Outcome: k.outcome, Count: n})
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Event != result[j].Event {
			return result[i].Event < result[j].Event
		}
		return result[i].Outcome < result[j].Outcome
	})
	return result, nil
}
