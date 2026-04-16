package watch

import "errors"

// Sentinel errors returned by the watch package.
var (
	ErrInvalidInterval = errors.New("watch: interval must be greater than zero")
	ErrNilHandler      = errors.New("watch: handler must not be nil")
)
