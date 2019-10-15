package main

import (
	"machine"

	"github.com/misterikkit/tinytimer/ws2812"
)

var (
	neo machine.Pin
	ws  ws2812.Device
)

func setup(g *game) {
	neo = machine.D13
	neo.Configure(machine.PinConfig{Mode: machine.PinOutput})
	ws = ws2812.New(neo)
}

func DisplayLEDs(f Frame) {
	ws.WriteColors(f)
}
