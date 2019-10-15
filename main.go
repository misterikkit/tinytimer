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
)

func main() {
	g := NewGame()
	setup(&g)
	nextTick := time.Now().Add(time.Second / FrameRate)
	tick := func() time.Duration {
		left := nextTick.Sub(time.Now())
		nextTick = nextTick.Add(time.Second / FrameRate)
		return left
	}
	for {
		g.update(time.Now())
		DisplayLEDs(*g.animation.f)
		// timeLeft := nextTick.Sub(time.Now())
		// tick()
		time.Sleep(tick())
	}
}
