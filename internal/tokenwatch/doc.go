// Package tokenwatch provides utilities for monitoring Vault token
// expiration. It fetches token metadata via a Fetcher interface,
// evaluates each token against a configurable warning threshold,
// and returns structured Status values that can be rendered or
// forwarded to alert hooks.
package tokenwatch
