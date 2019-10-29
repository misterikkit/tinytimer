package main

import (
	"device/arm"
	"device/sam"
	"image/color"
	"machine"

	"github.com/misterikkit/tinytimer/ws2812"
	"tinygo.org/x/drivers/apa102"
)

var (
	neo machine.Pin
	ws  ws2812.Device
)

func setup(g *game) {
	machine.D13.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.D13.Set(true)
	// Disable DotStar LED
	// TODO: make this work.
	dsClock := machine.PB02
	dsData := machine.PB03
	dsClock.Configure(machine.PinConfig{Mode: machine.PinOutput})
	dsData.Configure(machine.PinConfig{Mode: machine.PinOutput})
	dotStar := apa102.NewSoftwareSPI(dsClock, dsData, 120000)
	dotStar.WriteColors([]color.RGBA{{}})

	neo = machine.D5 // special level-shifted output pin
	neo.Configure(machine.PinConfig{Mode: machine.PinOutput})
	ws = ws2812.New(neo)

	btn2Min := machine.D11
	btn2Min.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	btn10Min := machine.D12
	btn10Min.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	btnCancel := machine.D10
	btnCancel.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})

	machine.PA00.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})

	g.pollInputs = func() {
		machine.D13.Set(machine.PA00.Get())
		if btnCancel.Get() {
			goToSleep()
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

func goToSleep() {
	// Turn off all LEDs
	machine.D13.Set(false)
	DisplayLEDs(newFrame())

	// Configure wake-up interrupt
	machine.PA00.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	// Enable peripheral mux
	sam.PORT.GROUP[0].PINCFG[0].SetBits(sam.PORT_GROUP_PINCFG_PMUXEN)
	// Select EIC peripheral for pin PA00
	sam.PORT.GROUP[0].PMUX[0].ClearBits(0xF)
	sam.PORT.GROUP[0].PMUX[0].SetBits(0x1)
	// Configure EIC to trigger on EXTINT0, which I assume is the peripheral
	// pin we just selected in the PORT mux.
	sam.EIC.CONFIG[0].SetBits(sam.EIC_CONFIG_SENSE0_HIGH)
	sam.EIC.CONFIG[0].ClearBits(sam.EIC_CONFIG_FILTEN0)
	sam.EIC.INTENSET.SetBits(0x01)
	// Also trigger on PA21 and PA23, the 2min and 10min buttons
	sam.EIC.CONFIG[0].SetBits(sam.EIC_CONFIG_SENSE5_HIGH)
	sam.EIC.CONFIG[0].SetBits(sam.EIC_CONFIG_SENSE7_HIGH)

	// Enter BACKUP sleep mode
	sam.PM.SLEEPCFG.Set(sam.PM_SLEEPCFG_SLEEPMODE_STANDBY)
	// Need to wait for SLEEPCFG propagation?
	arm.Asm("wfi")
	arm.SystemReset()
}
