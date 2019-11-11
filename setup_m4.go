// +build !wasm

package main

import (
	"device/sam"
	"image/color"
	"machine"

	"github.com/misterikkit/tinytimer/ws2812" // TODO: fix upstream
	"tinygo.org/x/drivers/apa102"
)

type userInterface struct {
	neoPix    ws2812.Device
	btnCancel machine.Pin
	btn2Min   machine.Pin
	btn10Min  machine.Pin
}

func (ui *userInterface) BtnCancel() bool            { return ui.btnCancel.Get() }
func (ui *userInterface) Btn2Min() bool              { return ui.btn2Min.Get() }
func (ui *userInterface) Btn10Min() bool             { return ui.btn10Min.Get() }
func (ui *userInterface) DisplayLEDs(c []color.RGBA) { ui.neoPix.WriteColors(c) }

func setup() userInterface {
	enableFPU()
	turnOffDotStar()

	neoPin := machine.D5
	configureOutput(neoPin)
	neoPix := ws2812.New(neoPin)

	btnCancel := machine.D2
	btn2Min := machine.D11
	btn10Min := machine.D12
	configureInput(btnCancel)
	configureInput(btn2Min)
	configureInput(btn10Min)

	return userInterface{neoPix, btnCancel, btn2Min, btn10Min}
}

func configureInput(p machine.Pin)  { p.Configure(machine.PinConfig{Mode: machine.PinInputPulldown}) }
func configureOutput(p machine.Pin) { p.Configure(machine.PinConfig{Mode: machine.PinOutput}) }

func enableFPU() {
	// See section 7.3.1 in
	// http://infocenter.arm.com/help/topic/com.arm.doc.ddi0439b/DDI0439B_cortex_m4_r0p0_trm.pdf
	sam.SystemControl.CPACR.SetBits(0xF << 20)
}

// turnOffDotStar writes a zero value to the on-board RGB LED to save power.
func turnOffDotStar() {
	onboardDotStar := apa102.NewSoftwareSPI(machine.PB02, machine.PB03, 100)
	onboardDotStar.WriteColors([]color.RGBA{{}})
	// 50% of the time, it works all of the time. ¯\_(ツ)_/¯
	onboardDotStar.WriteColors([]color.RGBA{{}})
}
