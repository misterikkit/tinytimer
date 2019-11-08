// This program demonstrates external interrupts on the ItsyBitsy M4 featuring
// ATSAMD51. The ItsyBitsy pins D2, D11, and D12 are configured as external
// interrupts. The EIC is configured to detect a rising edge.
// | pin | pad  | interrupt |
// | D2  | PA07 | EXTINT[7] |
// | D11 | PA21 | EXTINT[5] |
// | D12 | PA23 | EXTINT[7] |
package main

import (
	"device/sam"
	"machine"
	"time"
)

func main() {
	machine.LED.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.LED.Set(false)
	configureExternalInterrupt()
	configureExternalIntPins()
	for {
		if sam.EIC.INTFLAG.HasBits(1 << 5) {
			sam.EIC.INTFLAG.SetBits(1 << 5) // Clear int flag by writing 1 to it
			machine.LED.Set(true)
		}

		if sam.EIC.INTFLAG.HasBits(1 << 7) {
			sam.EIC.INTFLAG.SetBits(1 << 7) // Clear int flag by writing 1 to it
			machine.LED.Set(true)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func configureExternalInterrupt() {
	// 0. Disable EIC
	sam.EIC.CTRLA.ClearBits(sam.EIC_CTRLA_ENABLE)
	for sam.EIC.SYNCBUSY.Get() != 0 {
		// wait for disable
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

	// Set EXTINT[5] and EXTINT[7] to rising edge, filtered
	sam.EIC.CONFIG[0].SetBits(
		sam.EIC_CONFIG_SENSE5_RISE | sam.EIC_CONFIG_FILTEN5 |
			sam.EIC_CONFIG_SENSE7_RISE | sam.EIC_CONFIG_FILTEN7)
	// Enable EXTINT[5] adn EXTINT[7]
	sam.EIC.INTENSET.SetBits((1 << 5) | (1 << 7))

	// 5. Optionally, enable the asynchronous mode.
	// 6. Optionally, enable the debouncer mode.

	// 7. Enable the EIC by writing a ‘1’ to CTRLA.ENABLE.
	sam.EIC.CTRLA.SetBits(sam.EIC_CTRLA_ENABLE)
	for sam.EIC.SYNCBUSY.Get() != 0 {
		// wait for enable
	}
}

func configureExternalIntPins() {
	// Set each pin, PA07 PA21 PA23, to a pull-down input, and connect it to the
	// EIC through the peripheral mux.
	const A = 0
	for _, pin := range []uint{7, 21, 23} {
		sam.PORT.GROUP[A].DIRCLR.Set(1 << pin) // direction=input
		sam.PORT.GROUP[A].OUTCLR.Set(1 << pin) // pull=down
		sam.PORT.GROUP[A].PINCFG[pin].Set(
			sam.PORT_GROUP_PINCFG_INEN | sam.PORT_GROUP_PINCFG_PULLEN | sam.PORT_GROUP_PINCFG_PMUXEN)
		// Setting the peripheral mux to 0 selects peripheral A which is always the EIC.
		// Pins are all odd so clear the PMUXO half of the register.
		sam.PORT.GROUP[A].PMUX[pin/2].ClearBits(sam.PORT_GROUP_PMUX_PMUXO_Msk)
	}
}
