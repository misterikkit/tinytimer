package main

import (
	"time"

	"github.com/misterikkit/tinytimer/animation"
	"github.com/misterikkit/tinytimer/graphics"
)

const (
	FrameRate = 60
	TimeScale = 1.75
	FrameSize = 24
)

func main() {
	tickSize := time.Second / FrameRate
	ui := setup()
	spinner := animation.NewSpinner(graphics.K8SBlue)
	for {
		spinner.Update(time.Now())
		ui.neoPix.WriteColors(spinner.Frame)
		time.Sleep(tickSize)
	}
}
