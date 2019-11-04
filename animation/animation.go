package animation

import (
	"image/color"
	"time"

	"github.com/misterikkit/tinytimer/graphics"
	"golang.org/x/image/math/fixed"
)

// updateFn updates an animation based on the current time, and returns true if
// the animation is complete.
type updateFn = func(time.Time) bool

type spinner struct {
	Frame []color.RGBA
	dots  []graphics.Sprite
}

const spinnerCount = 7

// NewSpinner initializes a spinner animation.
func NewSpinner(c color.RGBA) spinner {
	const size = graphics.PixelWidth * 8 / 10
	s := spinner{
		Frame: make([]color.RGBA, graphics.FrameSize),
		dots:  make([]graphics.Sprite, 0, spinnerCount),
	}
	for i := 0; i < spinnerCount; i++ {
		s.dots = append(s.dots, graphics.Sprite{Size: size, Color: c})
	}
	return s
}

// Update computes the current frame of animation.
func (s *spinner) Update(now time.Time) bool {
	const period = time.Second * spinnerCount
	const divide = graphics.Circ / spinnerCount
	graphics.Fill(s.Frame, graphics.Black)

	// compute fraction through the period
	elapsed := now.Sub(now.Truncate(period)).Nanoseconds()
	period.Nanoseconds()
	p := elapsed * 64 / period.Nanoseconds()
	progress := fixed.Int26_6(p)
	// var progress fixed.Int26_6
	for i := range s.dots {
		s.dots[i].Position = (graphics.Circ.Mul(progress) + divide.Mul(fixed.I(i))) % graphics.Circ
		s.dots[i].Render(s.Frame)
	}
	return false
}
