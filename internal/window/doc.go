// Package window provides a generic sliding time window for collecting and
// pruning timestamped entries. It is safe for concurrent use.
//
// A sliding window retains only the entries whose timestamps fall within a
// specified duration relative to the current time. Older entries are pruned
// automatically when new entries are added or when an explicit prune is
// triggered. This makes the package suitable for rate tracking, metrics
// aggregation, and anomaly detection over rolling time intervals.
package window
