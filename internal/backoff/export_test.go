package backoff

// AttemptOf exposes the internal attempt counter for white-box tests.
func AttemptOf(b *Backoff) int { return b.attempt }
