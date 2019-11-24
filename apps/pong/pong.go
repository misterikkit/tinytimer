package pong

import (
	"image/color"
	"time"

	"github.com/misterikkit/tinytimer/animation"
	"github.com/misterikkit/tinytimer/graphics"
)

type App struct {
	s animation.Interface
}

func New() *App {
	return &App{animation.NewSpinner(graphics.Red)}
}

func (p *App) Update(now time.Time) { p.s.Update(now) }
func (p *App) Frame() []color.RGBA  { return p.s.Frame() }
