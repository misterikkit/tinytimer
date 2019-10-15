// +build !wasmsite

package main

import (
	"machine"

	"tinygo.org/x/drivers/ws2812"
)

var (
	neo Pin
	ws  ws2812.Device
)

func setup(g *game) {
	neo = machine.D13
	neo.Configure(machine.PinConfig{Mode: machine.PinOutput})
	ws = ws2812.New(new)
}

func DisplayLEDs(f Frame) {
	ws.WriteColors(f)
}
