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
	machine.LED.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.LED.Set(true)

	neo = machine.D5 // special level-shifted output pin
	neo.Configure(machine.PinConfig{Mode: machine.PinOutput})
	ws = ws2812.New(neo)

	btn2Min := machine.D11
	btn2Min.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	btn10Min := machine.D12
	btn10Min.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	btnCancel := machine.D2
	btnCancel.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})

	g.pollInputs = func() {
		if btnCancel.Get() {
			machine.LED.Set(true)

			g.event(CANCEL)
			return
		}
		if btn2Min.Get() {
			machine.LED.Set(false)
			g.event(TIMER_2M)
			return
		}
		if btn10Min.Get() {
			machine.LED.Set(false)
			g.event(TIMER_10M)
			return
		}

	}
}

func DisplayLEDs(f Frame) {
	ws.WriteColors(f)
}
