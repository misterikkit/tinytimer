package main

import (
	"device/sam"
	"machine"
)

var led_state bool

func handleExternalInterrupt(i int) {
	machine.LED.Set(led_state)
	led_state = !led_state
}

//go:export EIC_EXTINT_0_IRQHandler
func isrExtInt0() { handleExternalInterrupt(0) }

//go:export EIC_EXTINT_1_IRQHandler
func isrExtInt1() { handleExternalInterrupt(1) }

//go:export EIC_EXTINT_2_IRQHandler
func isrExtInt2() { handleExternalInterrupt(2) }

//go:export EIC_EXTINT_3_IRQHandler
func isrExtInt3() { handleExternalInterrupt(3) }

//go:export EIC_EXTINT_4_IRQHandler
func isrExtInt4() { handleExternalInterrupt(4) }

//go:export EIC_EXTINT_5_IRQHandler
func isrExtInt5() { handleExternalInterrupt(5) }

//go:export EIC_EXTINT_6_IRQHandler
func isrExtInt6() { handleExternalInterrupt(6) }

//go:export EIC_EXTINT_7_IRQHandler
func isrExtInt7() { handleExternalInterrupt(7) }

//go:export EIC_EXTINT_8_IRQHandler
func isrExtInt8() { handleExternalInterrupt(8) }

//go:export EIC_EXTINT_9_IRQHandler
func isrExtInt9() { handleExternalInterrupt(9) }

//go:export EIC_EXTINT_10_IRQHandler
func isrExtInt10() { handleExternalInterrupt(10) }

//go:export EIC_EXTINT_11_IRQHandler
func isrExtInt11() { handleExternalInterrupt(11) }

//go:export EIC_EXTINT_12_IRQHandler
func isrExtInt12() { handleExternalInterrupt(12) }

//go:export EIC_EXTINT_13_IRQHandler
func isrExtInt13() { handleExternalInterrupt(13) }

//go:export EIC_EXTINT_14_IRQHandler
func isrExtInt14() { handleExternalInterrupt(14) }

//go:export EIC_EXTINT_15_IRQHandler
func isrExtInt15() { handleExternalInterrupt(15) }

//go:export SysTick_Handler
func timer_isr() {
	machine.LED.Set(led_state)
	led_state = !led_state
}

func main() {
	println("is this thing on?")
	machine.LED.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.LED.Set(false)

	intPin := machine.PA00
	_ = intPin

	// 0. disable EIC?
	sam.EIC.CTRLA.ClearBits(sam.EIC_CTRLA_ENABLE)

	// 1. enable CLK_EIC_APB
	// it defaults to enabled. see section 15.6.2.6 table 15-1
	// 2. if required, enable NMI
	// SKIP
	// 3. Enable GCLK_EIC or CLK_ULP32K if needed
	//   - filtering, synch edge, or debouncing require clocks
	//   - CLK_ULP32K is for low power (but what about when I turn off all clocks?)
	// "The OSCULP32K is enabled by default after a Power-on Reset (POR), and will
	// always run except during POR" - section 29.6.4
	// Enable 32kHz output just to be safe
	sam.OSC32KCTRL.OSCULP32K.SetBits(sam.OSC32KCTRL_OSCULP32K_EN32K)
	// 3.a select clock source
	sam.EIC.CTRLA.SetBits(sam.EIC_CTRLA_CKSEL_CLK_ULP32K << sam.EIC_CTRLA_CKSEL_Pos)
	// 4. configure EIC input sense and filtering
	// 5. (optional) enable asynch mode
	// 6. (optional) enable debounce mode
	// 7. enable EIC! (Should I have disabled it for config changes?)
	sam.EIC.CTRLA.SetBits(sam.EIC_CTRLA_ENABLE)

	for {
	}
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
