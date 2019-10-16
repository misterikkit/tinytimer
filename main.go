package main

import (
	"machine"
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
	// g := NewGame()
	setup(nil)
	DisplayLEDs(newFrame()) // blank the LEDs

	spin := newSpinner()
	interval := time.Second / FrameRate

	// blink an LED to indicate that we're not stuck.
	ledOut := true

	for {
		spin.update(time.Now())
		DisplayLEDs(spin.f)

		machine.LED.Set(ledOut)
		ledOut = !ledOut
		time.Sleep(interval)
	}
}
