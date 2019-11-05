package main

import (
	"time"

	"github.com/misterikkit/tinytimer/game"
)

const (
	FrameRate = 60
	TimeScale = 1.75
	FrameSize = 24
)

func main() {
	tickSize := time.Second / FrameRate
	g := game.New()
	ui := setup(g)
	// spinner := animation.NewSpinner(graphics.K8SBlue)
	for {
		g.Update(time.Now())
		ui.DisplayLEDs(*g.Animation.Frame)
		// ui.neoPix.WriteColors(spinner.Frame)
		time.Sleep(tickSize)
	}
}
