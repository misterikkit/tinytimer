// +build !wasm

package hack

import "time"

// TimeScale is used to compensate for the misconfigured clock.
// TODO: Fix clock and ws2812 drivers instead of doing this.
const TimeScale = float32(1.75)

// ScaleDuration reduces a duration by a constant ratio to accommodate for a
// misconfigured clock. The problem with this approach is that I need to
// remember to call it everywhere I calculate a duration.
func ScaleDuration(d time.Duration) time.Duration {
	fd := float32(d)
	fd /= TimeScale
	return time.Duration(fd)
}
