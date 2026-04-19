// Package secretversion tracks version history for Vault secret paths,
// recording when each secret was first observed and when it was last updated.
// It is used to detect unexpected version rollbacks or stale secrets that
// have not been rotated within an expected window.
package secretversion
