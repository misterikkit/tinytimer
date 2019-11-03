package main

import (
	"image/color"
	"machine"

	"github.com/misterikkit/tinytimer/ws2812"
)

type userInterface struct {
	neoPix    ws2812.Device
	led       machine.Pin
	btnCancel machine.Pin
	btn2Min   machine.Pin
	btn10Min  machine.Pin
}

func setup() userInterface {
	neoPin := machine.D5
	makeOutput(machine.LED)
	makeOutput(neoPin)
	neoPix := ws2812.New(neoPin)

	btnCancel := machine.D2
	btn2Min := machine.D11
	btn10Min := machine.D12
	makeInput(btnCancel)
	makeInput(btn2Min)
	makeInput(btn10Min)
	return userInterface{neoPix, machine.LED, btnCancel, btn2Min, btn10Min}
}

func (f *userInterface) DisplayLEDs(c []color.RGBA) {
	f.neoPix.WriteColors(c)
}

func makeInput(p machine.Pin)  { p.Configure(machine.PinConfig{Mode: machine.PinInputPulldown}) }
func makeOutput(p machine.Pin) { p.Configure(machine.PinConfig{Mode: machine.PinOutput}) }
