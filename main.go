package main

import (
	"math"
	"time"
)

// Constants used in rendering
const (
	// Tau is better than Pi.
	Tau = 2 * math.Pi

	FrameRate = 60 // per second

	FrameSize = 24 // pixels

	PixelWidth = Tau / FrameSize // radians

	// TimeScale is used to compensate for the misconfigured clock.
	// TODO: Fix clock and ws2812 drivers instead of doing this.
	TimeScale = 1.75
)

func main() {
	g := NewGame()
	setup(&g)
	DisplayLEDs(newFrame()) // blank the LEDs

	tickInterval := scaleDuration(time.Second / FrameRate)
	nextTick := time.Now().Add(tickInterval)
	// return the time until next timer tick, and update `nextTick`
	tick := func() time.Duration {
		// TODO: skip a tick if needed
		left := nextTick.Sub(time.Now())
		nextTick = nextTick.Add(tickInterval)
		return left
	}

	for {
		g.update(time.Now())
		DisplayLEDs(*g.animation.f)
		time.Sleep(tick())
	}
}

// scaleDuration reduces a duration by a constant ratio to accommodate for a
// misconfigured clock. The problem with this approach is that I need to
// remember to call it everywhere I calculate a duration.
func scaleDuration(d time.Duration) time.Duration {
	fd := float32(d)
	fd /= TimeScale
	return time.Duration(fd)
}
