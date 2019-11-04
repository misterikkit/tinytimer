package animation

import (
	"image/color"
	"time"

	"github.com/misterikkit/tinytimer/graphics"
	"github.com/robaho/fixed"
)

// updateFn updates an animation based on the current time, and returns true if
// the animation is complete.
type updateFn = func(time.Time) bool

type Spinner struct {
	Frame []color.RGBA
	dots  []graphics.Sprite
}

const spinnerCount = 7

var size = graphics.PixelWidth.Mul(fixed.MustParse("0.8"))
var divide = graphics.Circ.Div(fixed.NewI(spinnerCount, 0))

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
	elapsed := fixed.NewI(now.Sub(now.Truncate(period)).Nanoseconds(), 0)
	p := fixed.NewI(period.Nanoseconds(), 0)
	progress := elapsed.Div(p)
	// p := elapsed * 64 / period.Nanoseconds()
	// progress := fixed.Int26_6(p)
	// var progress fixed.Int26_6
	for i := range s.dots {
		s.dots[i].Position = (graphics.Circ.Mul(progress).Add(divide.Mul(fixed.NewI(int64(i), 0)))) // TODO: mod
		s.dots[i].Render(s.Frame)
	}
	return false
}
