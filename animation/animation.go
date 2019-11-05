package animation

import (
	"image/color"
	"time"

	"github.com/misterikkit/tinytimer/graphics"
)

// updateFn updates an animation based on the current time, and returns true if
// the animation is complete.
type updateFn = func(time.Time) bool

type Spinner struct {
	Frame []color.RGBA
	dots  []graphics.Sprite
}

const spinnerCount = 7

var size = graphics.PixelWidth * 0.8
var divide = graphics.Circ / spinnerCount

// NewSpinner initializes a spinner animation.
func NewSpinner(c color.RGBA) Spinner {
	s := Spinner{
		Frame: make([]color.RGBA, graphics.FrameSize),
		dots:  make([]graphics.Sprite, 0, spinnerCount),
	}
	for i := 0; i < spinnerCount; i++ {
		s.dots = append(s.dots, graphics.Sprite{Size: size, Color: c})
	}
	return s
}

// Update computes the current frame of animation.
func (s *Spinner) Update(now time.Time) bool {
	const period = time.Second * spinnerCount
	graphics.Fill(s.Frame, graphics.Black)

	// compute fraction through the period
	elapsed := float32(now.Sub(now.Truncate(period)).Nanoseconds())
	p := float32(period.Nanoseconds())
	progress := elapsed / (p)
	// p := elapsed * 64 / period.Nanoseconds()
	// progress := fixed.Int26_6(p)
	// var progress fixed.Int26_6
	for i := range s.dots {
		s.dots[i].Position = graphics.Circ*progress + divide*float32(i) // TODO: mod
		s.dots[i].Render(s.Frame)
	}
	return false
}
