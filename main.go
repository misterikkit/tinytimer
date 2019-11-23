package main

import (
	"time"

	"github.com/misterikkit/tinytimer/game"
	"github.com/misterikkit/tinytimer/input"
)

const (
	FrameRate = 60
)

func main() {
	ui := setup()
	mgr := input.NewManager(ui.btnCancel.Get, ui.btn10Min.Get, ui.btn2Min.Get)
	mgr.AddHandler(input.ABC_Fall, func(input.Event) { ui.Sleepish() })
	g := game.New(mgr)
	for {
		mgr.Poll()
		g.Update(time.Now())
		ui.DisplayLEDs(g.Animation.Frame())
		// The effective frame rate is slightly less due to Update and DisplayLEDs,
		// but nobody will notice.
		time.Sleep(time.Second / FrameRate)
	}
}
