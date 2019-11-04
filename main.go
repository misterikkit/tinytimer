package main

import (
	"image/color"
	"time"

	"github.com/misterikkit/tinytimer/graphics"
	"golang.org/x/image/math/fixed"
)

const (
	FrameRate = 60
	TimeScale = 1.75
	FrameSize = 24
)

func main() {
	tickSize := time.Second / FrameRate
	ui := setup()
	frame := make([]color.RGBA, FrameSize)
	s := graphics.Sprite{Color: graphics.K8SBlue, Position: 0, Size: graphics.PixelWidth * 8 / 10}
	for i := 0; i < 7; i++ {
		s.Position = graphics.Circ.Mul(fixed.I(i)) / 7
		s.Render(frame)
	}
	for {
		ui.neoPix.WriteColors(frame)
		time.Sleep(tickSize)
	}
}
