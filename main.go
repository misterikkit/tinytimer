package main

import (
	"image/color"
	"machine"
	"time"

	"tinygo.org/x/drivers/ws2812"
)

const (
	FrameRate = 10
	TimeScale = 1.75
)

func main() {
	tickSize := time.Second / FrameRate
	machine.LED.Configure(machine.PinConfig{Mode: machine.PinOutput})

	neoPin := machine.D5
	neoPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	neoPix := ws2812.New(neoPin)
	black := make([]color.RGBA, 24)
	white := func() []color.RGBA {
		w := []color.RGBA{}
		for i := 0; i < 24; i++ {
			w = append(w, color.RGBA{64, 64, 64, 0})
		}
		return w
	}()

	blinkOn := true
	for {
		machine.LED.Set(blinkOn)
		if blinkOn {
			neoPix.WriteColors(white)
		} else {
			neoPix.WriteColors(black)
		}
		blinkOn = !blinkOn
		time.Sleep(tickSize)
	}
}
