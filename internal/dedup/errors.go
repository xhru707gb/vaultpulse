package dedup

import "errors"

// ErrInvalidWindow is returned when the suppression window is non-positive.
var ErrInvalidWindow = errors.New("dedup: window must be greater than zero")
