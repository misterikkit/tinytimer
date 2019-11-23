package main

import (
	"time"

	"github.com/misterikkit/tinytimer/easter"
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
	g := easter.App(game.New(mgr))
	easter.New(mgr, func(egg easter.App) { ui.Sleepish(); g = egg })
	// mgr.AddHandler(input.BC_Fall, func(input.Event) { g = &rainbow.Egg{} })
	for {
		mgr.Poll()
		g.Update(time.Now())
		ui.DisplayLEDs(g.Frame())
		// The effective frame rate is slightly less due to Update and DisplayLEDs,
		// but nobody will notice.
		time.Sleep(time.Second / FrameRate)
	}
}
