package main

import (
	"image/color"
	"time"

	"github.com/misterikkit/tinytimer/apps/rainbow"
	"github.com/misterikkit/tinytimer/apps/simon"
	"github.com/misterikkit/tinytimer/apps/timer"
	"github.com/misterikkit/tinytimer/easter"
	"github.com/misterikkit/tinytimer/input"
)

const frameRate = 60

type App interface {
	Update(time.Time)
	Frame() []color.RGBA
}

func main() {
	ui := setup() // initialize hardware
	mgr := input.NewManager(ui.btnCancel.Get, ui.btn10Min.Get, ui.btn2Min.Get)
	// Entering sleep mode is always available regardless of current app, so handle it here.
	mgr.AddHandler(input.ABC_Fall, func(input.Event) { ui.Sleepish() })
	// Set up initial app and easter egg launcher.
	app := App(timer.New(mgr))
	eggs := easter.New(mgr)
	for {
		mgr.Poll() // Invokes appropriate handlers.
		switch eggs.Get() {
		case easter.Rainbow:
			app = rainbow.New()
		case easter.Simon:
			app = simon.New(mgr)
		}

		// Compute and display next frame for the current app.
		app.Update(time.Now())
		ui.DisplayLEDs(app.Frame())

		// The effective frame rate is slightly less due to Update and DisplayLEDs,
		// but nobody will notice.
		time.Sleep(time.Second / frameRate)
	}
}
