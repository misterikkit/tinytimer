package main

import (
	"image/color"
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
	f := newFrame()
	f.fill(CSIOrange)
	DisplayLEDs(f) // blank the LEDs

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
	// dd := newDumb()
	// fl := float32(7)
	start := time.Now()
	for {
		prog := time.Since(start)
		// fl = float32(prog.Seconds())
		// if fl > 10 {
		// 	fl = 0
		// }
		val := math.Sin(prog.Seconds())
		val = val * val
		// fl = float32(val)
		// println(fl)
		// fl += 0.1
		// g.update(time.Now())
		// dd.update()
		// DisplayLEDs(dd.f)

		var c color.RGBA
		if val > 0 {
			c = scale(K8SBlue, float32(val))
		} else {
			c = scale(CSIOrange, float32(-val))
		}
		// c = scale(c, fl)
		f.fill(c)
		DisplayLEDs(f)
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
