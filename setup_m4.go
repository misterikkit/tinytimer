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
	btn2Min.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	btn10Min := machine.D12
	btn10Min.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	btnCancel := machine.D10
	btnCancel.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
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
	// dimColors(f)
	ws.WriteColors(f)
}

func dimColors(f Frame) {
	for i := range f {
		// Dim by 25%
		f[i] = scale(f[i], 0.25)
		// f[i].R >>= 2
		// f[i].G >>= 2
		// f[i].B >>= 2
	}
}
