// Package secretscore computes a numeric risk score for each monitored
// Vault secret. The score is derived from remaining TTL, rotation overdue
// status, and the number of active policy violations. Scores are bucketed
// into four risk levels: low, medium, high, and critical.
package secretscore
