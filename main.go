package main

import (
	"image/color"
	"time"

	gfx "github.com/misterikkit/tinytimer/graphics"
)

const (
	FrameRate = 60
	TimeScale = 1.75
	FrameSize = 24
)

func main() {
	tickSize := time.Second / FrameRate
	ui := setup()
	black := make([]color.RGBA, FrameSize)
	white := func() []color.RGBA {
		w := []color.RGBA{}
		for i := 0; i < FrameSize; i++ {
			w = append(w, gfx.White)
		}
		return w
	}()

	sprs := []gfx.Sprite{
		{gfx.Red, 0.0, gfx.PixelWidth},
		{gfx.K8SBlue, gfx.Tau / 3, gfx.PixelWidth},
		{gfx.CSIOrange, 2 * gfx.Tau / 3, gfx.PixelWidth},
	}
	once := []bool{true, true, true}

	for {
		btns := []bool{ui.btnCancel.Get(), ui.btn2Min.Get(), ui.btn10Min.Get()}
		for i := range btns {
			if once[i] && btns[i] {
				sprs[i].Render(black)
				once[i] = false
			}
		}
		press := ui.btnCancel.Get() || ui.btn2Min.Get() || ui.btn10Min.Get()
		ui.led.Set(press)
		if press {
			ui.neoPix.WriteColors(white)
		} else {
			ui.neoPix.WriteColors(black)
		}
		time.Sleep(tickSize)
	}
}
