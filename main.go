package main

import (
	"image/color"
	"time"

	"github.com/misterikkit/tinytimer/apps/rainbow"
	"github.com/misterikkit/tinytimer/apps/timer"
	"github.com/misterikkit/tinytimer/easter"
	"github.com/misterikkit/tinytimer/input"
)

const (
	FrameRate = 60
)

type App interface {
	Update(time.Time)
	Frame() []color.RGBA
}

func main() {
	ui := setup()
	mgr := input.NewManager(ui.btnCancel.Get, ui.btn10Min.Get, ui.btn2Min.Get)
	mgr.AddHandler(input.ABC_Fall, func(input.Event) { ui.Sleepish() })
	app := App(timer.New(mgr))
	eggs := easter.New(mgr)
	for {
		mgr.Poll()
		if eggs.Get() == easter.Rainbow {
			app = rainbow.New()
		}
		app.Update(time.Now())
		ui.DisplayLEDs(app.Frame())
		// The effective frame rate is slightly less due to Update and DisplayLEDs,
		// but nobody will notice.
		time.Sleep(time.Second / FrameRate)
	}
}
