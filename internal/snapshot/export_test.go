// export_test.go exposes internal fields for white-box testing.
package snapshot

import "time"

// SetNow overrides the clock used by Writer, enabling deterministic tests.
func (w *Writer) SetNow(fn func() time.Time) {
	w.now = fn
}
