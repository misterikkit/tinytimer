package main

import (
	"time"

	"github.com/misterikkit/tinytimer/game"
)

const (
	FrameRate = 60
)

func main() {
	ui := setup()
	g := game.New(ui)
	for {
		g.Update(time.Now())
		ui.DisplayLEDs(g.Animation.Frame())
		// The effective frame rate is slightly less due to Update and DisplayLEDs,
		// but nobody will notice.
		time.Sleep(time.Second / FrameRate)
	}
}
