// Package envelope provides secret envelope encryption helpers,
// wrapping a secret value with metadata about its encryption key version.
package envelope

import (
	"errors"
	"fmt"
	"time"
)

// ErrMissingCiphertext is returned when an Envelope has no ciphertext.
var ErrMissingCiphertext = errors.New("envelope: ciphertext is required")

// ErrMissingKeyVersion is returned when an Envelope has no key version.
var ErrMissingKeyVersion = errors.New("envelope: key version is required")

// Envelope holds an encrypted secret value together with its key metadata.
type Envelope struct {
	Path        string    `json:"path"`
	KeyVersion  string    `json:"key_version"`
	Ciphertext  string    `json:"ciphertext"`
	EncryptedAt time.Time `json:"encrypted_at"`
}

// New creates a new Envelope, stamping EncryptedAt to now (UTC).
func New(path, keyVersion, ciphertext string) (*Envelope, error) {
	if ciphertext == "" {
		return nil, ErrMissingCiphertext
	}
	if keyVersion == "" {
		return nil, ErrMissingKeyVersion
	}
	return &Envelope{
		Path:        path,
		KeyVersion:  keyVersion,
		Ciphertext:  ciphertext,
		EncryptedAt: time.Now().UTC(),
	}, nil
}

// Age returns how long ago the envelope was encrypted.
func (e *Envelope) Age() time.Duration {
	return time.Since(e.EncryptedAt)
}

// String returns a human-readable summary of the envelope.
func (e *Envelope) String() string {
	return fmt.Sprintf("path=%s key=%s age=%s", e.Path, e.KeyVersion, e.Age().Round(time.Second))
}
