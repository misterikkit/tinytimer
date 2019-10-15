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
	t := time.NewTicker(time.Second / FrameRate)
	for now := range t.C {
		g.update(now)
		DisplayLEDs(*g.animation.f)
	}
}
