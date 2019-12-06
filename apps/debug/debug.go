package debug

import (
	"image/color"
	"time"

	"github.com/misterikkit/tinytimer/graphics"
	"github.com/misterikkit/tinytimer/input"
)

type App struct {
	frame []color.RGBA
}

func New(ui *input.Manager) *App {
	return &App{make([]color.RGBA, graphics.FrameSize)}
}
func (d *App) Frame() []color.RGBA  { return d.frame }
func (d *App) Update(now time.Time) {}
