package main

import (
	"image/color"
	"time"
)

const (
	FrameRate = 60
	TimeScale = 1.75
)

func main() {
	tickSize := time.Second / FrameRate
	ui := setup()
	black := make([]color.RGBA, 24)
	white := func() []color.RGBA {
		w := []color.RGBA{}
		for i := 0; i < 24; i++ {
			w = append(w, color.RGBA{64, 64, 64, 0})
		}
		return w
	}()

	blinkOn := true
	for {
		press := ui.btnCancel.Get() || ui.btn2Min.Get() || ui.btn10Min.Get()
		ui.led.Set(press)
		if press {
			ui.neoPix.WriteColors(white)
		} else {
			ui.neoPix.WriteColors(black)
		}
		blinkOn = !blinkOn
		time.Sleep(tickSize)
	}
}
