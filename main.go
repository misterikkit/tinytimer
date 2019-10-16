package main

import (
	"math"
	"time"
)

// Constants used in rendering
const (
	// Tau is better than Pi.
	Tau = 2 * math.Pi

	FrameRate = 10 // per second

	FrameSize = 24 // pixels

	PixelWidth = Tau / FrameSize // radians
)

func main() {
	g := NewGame()
	setup(&g)
	f2 := newFrame()
	DisplayLEDs(f2) // blank the LEDs
	interval := time.Second / FrameRate
	var now time.Time
	for {
		g.update(now)
		DisplayLEDs(*g.animation.f)
		now = now.Add(interval)
		time.Sleep(interval)
	}
	// nextTick := time.Now().Add(time.Second / FrameRate)
	// // return the time until next timer tick, and update `nextTick`
	// tick := func() time.Duration {
	// 	// TODO: skip a tick if needed
	// 	left := nextTick.Sub(time.Now())
	// 	nextTick = nextTick.Add(time.Second / FrameRate)
	// 	return left
	// }
	// <-context.Background().Done()
	// ledOut := true
	// for {
	// 	machine.LED.Set(ledOut)
	// 	ledOut = !ledOut
	// 	g.update(time.Now())
	// 	DisplayLEDs(*g.animation.f)
	// 	// timeLeft := nextTick.Sub(time.Now())
	// 	// tick()
	// 	// if time.Now().Second()%2 == 0 {

	// 	// }
	// 	time.Sleep(time.Second / FrameRate)
	// }
}
