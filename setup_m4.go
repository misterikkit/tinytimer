// +build !wasm

package main

import (
	"device/arm"
	"device/sam"
	"image/color"
	"machine"
	"time"

	"tinygo.org/x/drivers/apa102"
	"tinygo.org/x/drivers/ws2812"
)

const frameSize = 24

type userInterface struct {
	neoPix      ws2812.Device
	btnCancel   machine.Pin
	btn2Min     machine.Pin
	btn10Min    machine.Pin
	gammaBuffer []color.RGBA
}

func (ui *userInterface) BtnCancel() bool { return ui.btnCancel.Get() }
func (ui *userInterface) Btn2Min() bool   { return ui.btn2Min.Get() }
func (ui *userInterface) Btn10Min() bool  { return ui.btn10Min.Get() }

func (ui *userInterface) Sleepish() { hibernate(ui) }

// DisplayLEDs writes the given color values to the LED, applying gamma
// correction in the process.
func (ui *userInterface) DisplayLEDs(c []color.RGBA) {
	for i := range ui.gammaBuffer {
		ui.gammaBuffer[i] = color.RGBA{
			R: gamma8[c[i].R/2], // extra red correction for maple
			G: gamma8[c[i].G],
			B: gamma8[c[i].B],
		}
	}
	mask := arm.DisableInterrupts()
	ui.neoPix.WriteColors(ui.gammaBuffer)
	arm.EnableInterrupts(mask)
}

func setup() *userInterface {
	enableFPU()
	turnOffDotStar()

	neoPin := machine.D5
	configureOutput(neoPin)
	ui := userInterface{
		neoPix:      ws2812.New(neoPin),
		btnCancel:   machine.D2,
		btn2Min:     machine.D12,
		btn10Min:    machine.D11,
		gammaBuffer: make([]color.RGBA, frameSize),
	}
	configureInput(ui.btnCancel)
	configureInput(ui.btn2Min)
	configureInput(ui.btn10Min)

	return &ui
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

// hibernate turns off the display and puts the device into the lowest available
// power mode. There is no waking up from this.
func hibernate(ui *userInterface) {
	/////////////
	// Blank LEDs
	ui.DisplayLEDs(make([]color.RGBA, frameSize))
	// LEDs fail to turn off if I don't sleep a moment.
	time.Sleep(time.Millisecond)

	//////////////////////////
	// Enter OFF sleep mode
	sam.PM.SLEEPCFG.Set(sam.PM_SLEEPCFG_SLEEPMODE_OFF)
	for !sam.PM.INTFLAG.HasBits(sam.PM_INTFLAG_SLEEPRDY) {
	}
	arm.Asm("wfi")

	////////////////////////////////////////////
	// The system is now halted until hard reset
}

// gamma correction for 8-bit color values yoinked from
// https://learn.adafruit.com/led-tricks-gamma-correction/the-quick-fix
var gamma8 = [...]uint8{
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2,
	2, 3, 3, 3, 3, 3, 3, 3, 4, 4, 4, 4, 4, 5, 5, 5,
	5, 6, 6, 6, 6, 7, 7, 7, 7, 8, 8, 8, 9, 9, 9, 10,
	10, 10, 11, 11, 11, 12, 12, 13, 13, 13, 14, 14, 15, 15, 16, 16,
	17, 17, 18, 18, 19, 19, 20, 20, 21, 21, 22, 22, 23, 24, 24, 25,
	25, 26, 27, 27, 28, 29, 29, 30, 31, 32, 32, 33, 34, 35, 35, 36,
	37, 38, 39, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 50,
	51, 52, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 66, 67, 68,
	69, 70, 72, 73, 74, 75, 77, 78, 79, 81, 82, 83, 85, 86, 87, 89,
	90, 92, 93, 95, 96, 98, 99, 101, 102, 104, 105, 107, 109, 110, 112, 114,
	115, 117, 119, 120, 122, 124, 126, 127, 129, 131, 133, 135, 137, 138, 140, 142,
	144, 146, 148, 150, 152, 154, 156, 158, 160, 162, 164, 167, 169, 171, 173, 175,
	177, 180, 182, 184, 186, 189, 191, 193, 196, 198, 200, 203, 205, 208, 210, 213,
	215, 218, 220, 223, 225, 228, 231, 233, 236, 239, 241, 244, 247, 249, 252, 255}
