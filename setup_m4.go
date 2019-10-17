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
	neo = machine.D5 // special level-shifted output pin
	neo.Configure(machine.PinConfig{Mode: machine.PinOutput})
	ws = ws2812.New(neo)

	btn2Min := machine.D11
	btn2Min.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	btn10Min := machine.D12
	btn10Min.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	btnCancel := machine.D10
	btnCancel.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	g.pollInputs = func() {
		if btnCancel.Get() {
			g.event(CANCEL)
			return
		}
		if btn2Min.Get() {
			g.event(TIMER_2M)
			return
		}
		if btn10Min.Get() {
			g.event(TIMER_10M)
			return
		}

	}
}

func DisplayLEDs(f Frame) {
	ws.WriteColors(f)
}
