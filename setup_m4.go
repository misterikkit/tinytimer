package main

import (
	_ "device/arm"
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

//export EIC_EXTINT_0_IRQHandler
func extint0_handler() {
	machine.LED.Set(true)
}

func setupExtInt() {
	machine.LED.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.LED.Set(false)
	intPin := machine.PA00

	// some bits copied from machine_atsamd51g19.go

	// CONFIGURE THE PIN
	sam.PORT.GROUP[0].DIRCLR.Set(1 << uint8(intPin)) // input
	sam.PORT.GROUP[0].OUTCLR.Set(1 << uint8(intPin)) // pulldown
	// sam.PORT.GROUP[0].OUTSET.Set(1 << uint8(intPin)) // pullup
	// enable input, pull, and pmux
	setPinCfg(intPin,
		sam.PORT_GROUP_PINCFG_INEN|
			sam.PORT_GROUP_PINCFG_PULLEN| // comment out to float pin
			sam.PORT_GROUP_PINCFG_PMUXEN,
	)

	const EIC_PERIPH_MUX = 0 // 0==A in the table 6-1 (I guess??)
	// even pin, so save the odd pins
	val := getPMux(intPin) & sam.PORT_GROUP_PMUX_PMUXO_Msk
	setPMux(intPin, val|(uint8(EIC_PERIPH_MUX)<<sam.PORT_GROUP_PMUX_PMUXE_Pos))

	// ENABLE EXTINT0
	sam.EIC.INTENSET.Set(1) // set 0th bit
	sam.EIC.CONFIG[0].SetBits(sam.EIC_CONFIG_SENSE0_HIGH)
	sam.EIC.CONFIG[0].SetBits(sam.EIC_CONFIG_FILTEN0)
}

func setup(g *game) {
	machine.LED.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.LED.Set(true)
	setupExtInt()
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
	btnCancel := machine.D2
	btnCancel.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})

	machine.PA00.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})

	g.pollInputs = func() {
		// machine.D13.Set(machine.PA00.Get())
		if btnCancel.Get() && btn10Min.Get() {
			goToSleep()
		}
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

func goToSleep() {
	////////////////////
	// Turn off all LEDs
	machine.D13.Set(false)
	DisplayLEDs(newFrame())

	//////////////////////////////
	// Configure wake-up interrupt

	// Port group PA is sam.PORT.GROUP[0]
	// Port PA00 is configured by sam.PORT.GROUP[0].PINCFG[0]
	// PA00 is wired to MOSI on the ItsyBitsy M4
	machine.PA00.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	// Enable peripheral mux
	sam.PORT.GROUP[0].PINCFG[0].SetBits(sam.PORT_GROUP_PINCFG_PMUXEN)
	// Select EIC peripheral for port PA00
	// Peripheral mux uses 4 bits. PA00 should be the low 4 bits of PMUX[0]
	sam.PORT.GROUP[0].PMUX[0].ClearBits(0xF)
	// mux value 1 => column A in table 6-1 of the SAMD51 guide (I guess??)
	// sam.PORT.GROUP[0].PMUX[0].SetBits(0x1)
	// PA00 is wired to EXTINT[0] according to table 6-1 of the SAMD51 guide
	// Configure EIC to trigger on EXTINT[0].
	// We want level-based sensing with no filters in order to do asynch
	// interrupts. Needed because we will be powering down all clocks.
	sam.EIC.CONFIG[0].SetBits(sam.EIC_CONFIG_SENSE0_HIGH)
	sam.EIC.CONFIG[0].ClearBits(sam.EIC_CONFIG_FILTEN0)
	// Enable EXTINT[0]
	sam.EIC.INTENSET.SetBits(0x01)

	return

	//////////////////////////
	// Enter BACKUP sleep mode
	sam.PM.SLEEPCFG.Set(sam.PM_SLEEPCFG_SLEEPMODE_BACKUP)
	// Need to wait for SLEEPCFG propagation?
	for !sam.PM.INTFLAG.HasBits(sam.PM_INTFLAG_SLEEPRDY) {
	}
	// arm.Asm("wfi")

	/////////////////////////////
	// After wake-up, reset board
	// arm.SystemReset()
}

// getPMux returns the value for the correct PMUX register for this pin.
// Copied from tinygo src
func getPMux(p machine.Pin) uint8 {
	switch {
	case p < 32:
		return sam.PORT.GROUP[0].PMUX[uint8(p)>>1].Get()
	case p >= 32 && p < 64:
		return sam.PORT.GROUP[1].PMUX[uint8(p-32)>>1].Get()
	default:
		return 0
	}
}

// setPMux sets the value for the correct PMUX register for this pin.
// Copied from tinygo src
func setPMux(p machine.Pin, val uint8) {
	switch {
	case p < 32:
		sam.PORT.GROUP[0].PMUX[uint8(p)>>1].Set(val)
	case p >= 32 && p < 64:
		sam.PORT.GROUP[1].PMUX[uint8(p-32)>>1].Set(val)
	}
}

// setPinCfg sets the value for the correct PINCFG register for this pin.
// Copied from tinygo src
func setPinCfg(p machine.Pin, val uint8) {
	switch {
	case p < 32:
		sam.PORT.GROUP[0].PINCFG[p].Set(val)
	case p >= 32 && p <= 64:
		sam.PORT.GROUP[1].PINCFG[p-32].Set(val)
	}
}
