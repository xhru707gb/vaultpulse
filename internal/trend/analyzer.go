// Package trend analyses historical audit log entries to detect
// secret expiry and rotation trends over time.
package trend

import (
	"sort"
	"time"
)

// Point represents a single data point in a trend series.
type Point struct {
	Timestamp time.Time
	Count     int
}

// Report holds trend analysis results.
type Report struct {
	Path       string
	EventType  string
	Points     []Point
	TotalCount int
	First      time.Time
	Last       time.Time
}

// Analyzer computes trends from a slice of audit entries.
type Analyzer struct {
	bucketSize time.Duration
}

// NewAnalyzer returns an Analyzer that buckets events by bucketSize.
func NewAnalyzer(bucketSize time.Duration) (*Analyzer, error) {
	if bucketSize <= 0 {
		return nil, ErrInvalidBucket
	}
	return &Analyzer{bucketSize: bucketSize}, nil
}

// Entry is a minimal audit record consumed by the analyzer.
type Entry struct {
	Path      string
	EventType string
	Timestamp time.Time
}

// Analyse groups entries by path+eventType and builds time-bucketed reports.
func (a *Analyzer) Analyse(entries []Entry) []Report {
	type key struct{ path, event string }
	buckets := make(map[key]map[int64]int)

	for _, e := range entries {
		k := key{e.Path, e.EventType}
		if buckets[k] == nil {
			buckets[k] = make(map[int64]int)
		}
		b := e.Timestamp.Truncate(a.bucketSize).Unix()
		buckets[k][b]++
	}

	var reports []Report
	for k, bmap := range buckets {
		var points []Point
		for ts, cnt := range bmap {
			points = append(points, Point{Timestamp: time.Unix(ts, 0).UTC(), Count: cnt})
		}
		sort.Slice(points, func(i, j int) bool {
			return points[i].Timestamp.Before(points[j].Timestamp)
		})
		total := 0
		for _, p := range points {
			total += p.Count
		}
		reports = append(reports, Report{
			Path:       k.path,
			EventType:  k.event,
			Points:     points,
			TotalCount: total,
			First:      points[0].Timestamp,
			Last:       points[len(points)-1].Timestamp,
		})
	}
	sort.Slice(reports, func(i, j int) bool {
		if reports[i].Path != reports[j].Path {
			return reports[i].Path < reports[j].Path
		}
		return reports[i].EventType < reports[j].EventType
	})
	return reports
}
