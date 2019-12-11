// +build wasm

package simon

import (
	"math/rand"
	"time"
)

func initRandom() {
	rand.Seed(time.Now().UnixNano())
}

func randU32() uint32 {
	return rand.Uint32()
}
