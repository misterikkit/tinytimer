// +build !js
// +build !wasm

package main

import (
	"device/sam"
	"image/color"
	"machine"

	"github.com/misterikkit/tinytimer/game"
	"github.com/misterikkit/tinytimer/ws2812"
)

type userInterface struct {
	neoPix    ws2812.Device
	led       machine.Pin
	btnCancel machine.Pin
	btn2Min   machine.Pin
	btn10Min  machine.Pin
}

func setup(g *game.Game) userInterface {
	enableFPU()
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
	g.PollInputs = func() {
		switch {
		case btnCancel.Get():
			g.Event(game.CANCEL)
		case btn2Min.Get():
			g.Event(game.TIMER_2M)
		case btn10Min.Get():
			g.Event(game.TIMER_10M)
		}
	}
	return userInterface{neoPix, machine.LED, btnCancel, btn2Min, btn10Min}
}

func (f *userInterface) DisplayLEDs(c []color.RGBA) {
	f.neoPix.WriteColors(c)
}

func makeInput(p machine.Pin)  { p.Configure(machine.PinConfig{Mode: machine.PinInputPulldown}) }
func makeOutput(p machine.Pin) { p.Configure(machine.PinConfig{Mode: machine.PinOutput}) }

func enableFPU() {
	// See section 7.3.1 in
	// http://infocenter.arm.com/help/topic/com.arm.doc.ddi0439b/DDI0439B_cortex_m4_r0p0_trm.pdf
	sam.SystemControl.CPACR.SetBits(0xF << 20)

}
