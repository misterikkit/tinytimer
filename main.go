package main

import (
	"math"
	"time"
)

// Constants used in rendering
const (
	// Tau is better than Pi.
	Tau = 2 * math.Pi

	FrameRate = 30 // per second

	FrameSize = 24 // pixels

	PixelWidth = Tau / FrameSize // radians
)

func main() {
	g := NewGame()
	setup(&g)
	DisplayLEDs(newFrame()) // blank the LEDs

	nextTick := time.Now().Add(time.Second / FrameRate)
	// return the time until next timer tick, and update `nextTick`
	tick := func() time.Duration {
		// TODO: skip a tick if needed
		left := nextTick.Sub(time.Now())
		nextTick = nextTick.Add(time.Second / FrameRate)
		return left
	}

	for {
		g.update(time.Now())
		DisplayLEDs(*g.animation.f)
		time.Sleep(tick())
	}
}
