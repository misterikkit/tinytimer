package rainbow

import (
	"image/color"
	"time"

	"github.com/lucasb-eyer/go-colorful"
)

type App struct {
	frame  []color.RGBA
	offset float64
}

func New() *App {
	return &App{frame: make([]color.RGBA, 24)}
}

func (a *App) Update(t time.Time) {
	const period = 7 * time.Second
	offset := 360.0 * t.Sub(t.Truncate(period)).Seconds() / period.Seconds()

	for i := 0; i < 24; i++ {
		c := colorful.Hcl(offset, 0.75, 0.25)
		r, g, b := c.RGB255()
		a.frame[i] = color.RGBA{r, g, b, 0}
		offset += 360.0 / 24.0
	}
}

func (e *App) Frame() []color.RGBA { return e.frame }
