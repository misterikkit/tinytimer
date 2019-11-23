package main

import (
	"image/color"
	"time"

	"github.com/misterikkit/tinytimer/game"
	"github.com/misterikkit/tinytimer/input"
	"github.com/misterikkit/tinytimer/rainbow"
)

const (
	FrameRate = 60
)

type app interface {
	Update(time.Time)
	Frame() []color.RGBA
}

func main() {
	ui := setup()
	mgr := input.NewManager(ui.btnCancel.Get, ui.btn10Min.Get, ui.btn2Min.Get)
	g := app(game.New(mgr))
	mgr.AddHandler(input.ABC_Fall, func(input.Event) { ui.Sleepish() })
	mgr.AddHandler(input.BC_Fall, func(input.Event) { g = &rainbow.Egg{} })
	for {
		mgr.Poll()
		g.Update(time.Now())
		ui.DisplayLEDs(g.Frame())
		// The effective frame rate is slightly less due to Update and DisplayLEDs,
		// but nobody will notice.
		time.Sleep(time.Second / FrameRate)
	}
}
