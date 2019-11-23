package rainbow

import (
	"image/color"
	"time"

	"github.com/lucasb-eyer/go-colorful"
)

type Egg struct{ offset float64 }

func (e *Egg) Update(t time.Time) {
	const period = 10 * time.Second
	e.offset = 360.0 * t.Sub(t.Truncate(period)).Seconds() / period.Seconds()
}

func (e *Egg) Frame() []color.RGBA {
	frame := make([]color.RGBA, 0, 24)
	h := e.offset
	for i := 0; i < 24; i++ {
		c := colorful.Hcl(h, 0.75, 0.25)
		r, g, b := c.RGB255()
		frame = append(frame, color.RGBA{r, g, b, 0})
		h += 360.0 / 24.0
	}
	return frame
}
