package main

import (
	"image/color"
	"time"

	"github.com/misterikkit/tinytimer/apps/arcade"
	"github.com/misterikkit/tinytimer/apps/colorpicker"
	"github.com/misterikkit/tinytimer/apps/pong"
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
	println("hello world")
	ui := setup() // initialize hardware
	mgr := input.NewManager(ui.BtnCancel, ui.Btn2Min, ui.Btn10Min)
	eggs := easter.New(mgr)

	var (
		timer   = timer.New(mgr)
		rainbow = rainbow.New()
		pong    = pong.New(mgr)
		picker  = colorpicker.New(mgr)
		simon   = simon.New(mgr)
		arcade  = arcade.New(mgr)
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
		case easter.ColorPicker:
			if isTimer(app) {
				picker.Reset()
				app = picker
			}
		case easter.Simon:
			if isTimer(app) {
				simon.Reset()
				app = simon
			}
		case easter.Arcade:
			if isTimer(app) {
				arcade.Reset(time.Now())
				app = arcade
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
