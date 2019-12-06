package main

import (
	"image/color"
	"time"

	"github.com/misterikkit/tinytimer/apps/debug"
	"github.com/misterikkit/tinytimer/apps/pong"
	"github.com/misterikkit/tinytimer/apps/rainbow"
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
	println("hello world")
	ui := setup() // initialize hardware
	mgr := input.NewManager(ui.BtnCancel, ui.Btn10Min, ui.Btn2Min)
	eggs := easter.New(mgr)

	var (
		timer   = timer.New(mgr)
		rainbow = rainbow.New()
		pong    = pong.New(mgr)
		debug   = debug.New(mgr)
	)
	app := App(timer)
	for {
		mgr.Poll() // Invokes appropriate handlers.
		switch eggs.Get() {
		case easter.Eggsit:
			if isTimer(app) {
				ui.Sleepish()
			}
			timer.Reset()
			app = timer
		case easter.Rainbow:
			if isTimer(app) {
				app = rainbow
			}
		case easter.Pong:
			if isTimer(app) {
				pong.Reset()
				app = pong
			}
		case easter.Debug:
			if isTimer(app) {
				app = debug
			}
		}

		// Compute and display next frame for the current app.
		app.Update(time.Now())
		ui.DisplayLEDs(app.Frame())

		// The effective frame rate is slightly less due to Update and DisplayLEDs,
		// but nobody will notice.
		time.Sleep(time.Second / frameRate)
	}
}

func isTimer(a App) bool {
	_, ok := a.(*timer.App)
	return ok
}
