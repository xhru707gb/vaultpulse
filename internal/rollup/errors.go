package rollup

import "errors"

// ErrInvalidWindow is returned when a non-positive window duration is supplied.
var ErrInvalidWindow = errors.New("rollup: window must be greater than zero")
