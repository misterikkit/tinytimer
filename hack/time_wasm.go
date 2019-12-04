// +build wasm

package hack

import "time"

// TimeScale is not needed in WASM land
const TimeScale = float32(1.0)

// ScaleDuration is a noop.
func ScaleDuration(d time.Duration) time.Duration {
	return d
}
