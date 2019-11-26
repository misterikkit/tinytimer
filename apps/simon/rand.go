// +build !wasm

package simon

import (
	"device/sam"
)

func initRandom() {
	sam.MCLK.APBCMASK.SetBits(sam.MCLK_APBCMASK_TRNG_)
	sam.TRNG.CTRLA.SetBits(sam.TRNG_CTRLA_ENABLE)
}

func rand() uint32 {
	// TODO: block until a new random is ready.
	return sam.TRNG.DATA.Get()
}
