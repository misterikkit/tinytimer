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

func configureExternalInterrupt() {
	// 0. Disable EIC
	sam.EIC.CTRLA.ClearBits(sam.EIC_CTRLA_ENABLE)
	for sam.EIC.SYNCBUSY.Get() != 0 {
	}

	// 1. Enable CLK_EIC_APB
	// Should be enabled by default, but let's be safe
	sam.MCLK.APBAMASK.SetBits(sam.MCLK_APBAMASK_EIC_)

	// 2. If required, configure the NMI by writing the Non-Maskable Interrupt Control register (NMICTRL)

	// 3. Enable GCLK_EIC or CLK_ULP32K when one of the following configuration is selected:
	// – the NMI uses edge detection or filtering.
	// – one EXTINT uses filtering.
	// – one EXTINT uses synchronous edge detection.
	// – one EXTINT uses debouncing.
	// GCLK_EIC is used when a frequency higher than 32KHz is required for filtering.
	// CLK_ULP32K is recommended when power consumption is the priority. For CLK_ULP32K write a
	// '1' to the Clock Selection bit in the Control A register (CTRLA.CKSEL).
	sam.EIC.CTRLA.SetBits(sam.EIC_CTRLA_CKSEL_CLK_ULP32K << sam.EIC_CTRLA_CKSEL_Pos)

	// 4. Configure the EIC input sense and filtering by writing the Configuration n register (CONFIG).
	// TODO:
	// D2==PA07==EXTINT[7], D11==PA21==EXTINT[5], D12==PA23==EXTINT[7]

	// 5. Optionally, enable the asynchronous mode.
	// 6. Optionally, enable the debouncer mode.

	// 7. Enable the EIC by writing a ‘1’ to CTRLA.ENABLE.
	sam.EIC.CTRLA.SetBits(sam.EIC_CTRLA_ENABLE)
	for sam.EIC.SYNCBUSY.Get() != 0 {
	}
}
