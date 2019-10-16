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

	FrameRate = 60 // per second

	FrameSize = 24 // pixels

	PixelWidth = Tau / FrameSize // radians
)

func main() {
	// g := NewGame()
	// // setup(&g)
	// nextTick := time.Now().Add(time.Second / FrameRate)
	// tick := func() time.Duration {
	// 	left := nextTick.Sub(time.Now())
	// 	nextTick = nextTick.Add(time.Second / FrameRate)
	// 	return left
	// }

	btn := machine.D9
	btn.Configure(machine.PinConfig{machine.PinInput})
	led := machine.D7
	led.Configure(machine.PinConfig{machine.PinOutput})
	led.Set(false)

	delay := 100 * time.Millisecond
	for {
		if btn.Get() {
			delay = time.Second
		}
		// machine.LED.Set(btn.Get())
		machine.LED.Set(true)
		led.Set(false)
		time.Sleep(delay)
		machine.LED.Set(false)
		led.Set(true)
		time.Sleep(delay)
	}

	// for {
	// 	g.update(time.Now())
	// 	DisplayLEDs(*g.animation.f)
	// 	// timeLeft := nextTick.Sub(time.Now())
	// 	// tick()
	// 	if time.Now().Second()%2 == 0 {

	// 	}
	// 	time.Sleep(tick())
	// }
}
