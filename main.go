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
	for {
		g.Update(time.Now())
		ui.DisplayLEDs(*g.Animation.Frame)
		time.Sleep(tickSize)
	}
}
