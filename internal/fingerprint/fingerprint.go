// Package fingerprint provides secret fingerprinting to detect
// value changes without storing sensitive plaintext.
package fingerprint

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strings"
)

// ErrEmptyInput is returned when the input map is empty or nil.
var ErrEmptyInput = errors.New("fingerprint: input map is empty")

// Result holds the computed fingerprint for a secret path.
type Result struct {
	Path        string
	Fingerprint string
	KeyCount    int
}

// Compute produces a deterministic SHA-256 fingerprint of the
// key-value pairs at the given path. Keys are sorted before hashing
// so insertion order does not affect the result.
func Compute(path string, data map[string]string) (Result, error) {
	if len(data) == 0 {
		return Result{}, ErrEmptyInput
	}

	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		fmt.Fprintf(&sb, "%s=%s;", k, data[k])
	}

	sum := sha256.Sum256([]byte(sb.String()))
	return Result{
		Path:        path,
		Fingerprint: hex.EncodeToString(sum[:]),
		KeyCount:    len(keys),
	}, nil
}

// Changed returns true when the fingerprint of data differs from prev.
func Changed(path, prev string, data map[string]string) (bool, Result, error) {
	r, err := Compute(path, data)
	if err != nil {
		return false, Result{}, err
	}
	return r.Fingerprint != prev, r, nil
}
