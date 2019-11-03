package main

import (
	"image/color"
	"machine"

	"github.com/misterikkit/tinytimer/ws2812"
)

type userInterface struct {
	neoPix ws2812.Device
	led    machine.Pin
}

func setup() userInterface {
	neoPin := machine.D5
	machine.LED.Configure(machine.PinConfig{Mode: machine.PinOutput})
	neoPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	neoPix := ws2812.New(neoPin)
	return userInterface{neoPix, machine.LED}
}

func (f *userInterface) DisplayLEDs(c []color.RGBA) {
	f.neoPix.WriteColors(c)
}
