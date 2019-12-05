package rainbow

import (
	"image/color"
	"math"
	"time"

	"github.com/misterikkit/tinytimer/graphics"
)

type App struct {
	frame  []color.RGBA
	offset float64
}

func New() *App {
	return &App{frame: make([]color.RGBA, 24)}
}

func (a *App) Update(t time.Time) {
	const tau = 2 * math.Pi
	const period = 7 * time.Second
	offset := t.Sub(t.Truncate(period)).Seconds() / period.Seconds()

	for i := range a.frame {
		pos := float64(i)/float64(len(a.frame)) + offset
		// TODO: find a way to increase red/purple/blue and decrease green
		c := color.YCbCr{
			Y:  uint8(255 / 2),
			Cb: uint8(255 * math.Sin(tau*pos)),
			Cr: uint8(255 * math.Cos(tau*pos)),
		}
		r, g, b, _ := c.RGBA()
		a.frame[i] = graphics.Scale(color.RGBA{
			R: uint8(r >> (8)),
			G: uint8(g >> (8)),
			B: uint8(b >> (8)),
		}, graphics.MaxIntensity)
	}
}

func (e *App) Frame() []color.RGBA { return e.frame }
