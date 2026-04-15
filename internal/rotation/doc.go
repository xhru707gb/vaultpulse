// Package rotation implements rotation schedule tracking for vaultpulse.
//
// It provides an Evaluator that compares a secret's last-rotated timestamp
// against a configured rotation interval to determine whether rotation is
// overdue or upcoming. Results are surfaced as Status values and can be
// rendered via FormatTable for CLI output.
//
// Typical usage:
//
//	eval := rotation.NewEvaluator(nil)
//	statuses, errs := eval.EvaluateAll(schedules)
//	rotation.FormatTable(os.Stdout, statuses)
package rotation
