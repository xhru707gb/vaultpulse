package trend

import "errors"

// ErrInvalidBucket is returned when a non-positive bucket size is provided.
var ErrInvalidBucket = errors.New("trend: bucket size must be greater than zero")
