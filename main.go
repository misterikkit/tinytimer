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
	g := NewGame()
	setup(&g)
	DisplayLEDs(newFrame()) // blank the LEDs

	interval := time.Second / FrameRate
	// nextTick := time.Now().Add(time.Second / FrameRate)
	// // return the time until next timer tick, and update `nextTick`
	// tick := func() time.Duration {
	// 	// TODO: skip a tick if needed
	// 	left := nextTick.Sub(time.Now())
	// 	nextTick = nextTick.Add(time.Second / FrameRate)
	// 	return left
	// }

	ledOut := true
	// l := newLoader(K8SBlue, time.Now(), time.Now().Add(3*time.Second))
	// fl := newFlasher(K8SBlue, time.Now())
	dd := newDumb()
	for {
		// g.update(time.Now())
		dd.update()
		DisplayLEDs(dd.f)

		machine.LED.Set(ledOut)
		ledOut = !ledOut
		time.Sleep(interval)

		// 	DisplayLEDs(*g.animation.f)
		// 	// timeLeft := nextTick.Sub(time.Now())
		// 	// tick()
		// 	// if time.Now().Second()%2 == 0 {

		// 	// }
		// 	time.Sleep(time.Second / FrameRate)
	}
}
