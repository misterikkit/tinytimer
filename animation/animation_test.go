package animation_test

import (
	"fmt"
	"image/color"
	"testing"
	"time"

	"github.com/misterikkit/tinytimer/animation"
	"github.com/misterikkit/tinytimer/graphics"
	"github.com/misterikkit/tinytimer/hack"
	"github.com/stretchr/testify/assert"
)

// newFrame returns a full FrameSize-sized newFrame populated with the given pixel values.
func newFrame(cs ...color.RGBA) frame {
	f := make([]color.RGBA, graphics.FrameSize)
	copy(f, cs)
	return f
}

// frame overrides GoString on a color slice for better test output.
type frame []color.RGBA

func (f frame) GoString() string {
	pxs := make([]string, len(f))
	for i := range f {
		pxs[i] = fmt.Sprintf("#%02x%02x%02x", f[i].R, f[i].G, f[i].B)
	}
	return fmt.Sprint(pxs)
}

func TestSpinner(t *testing.T) {
	s := animation.NewSpinner(color.RGBA{255, 255, 255, 0})
	var t0 time.Time
	s.Update(t0)
	f0 := newFrame()
	copy(f0, s.Frame())
	s.Update(t0.Add(hack.ScaleDuration(time.Second)))
	f1 := newFrame()
	copy(f1, s.Frame())
	assert.Equal(t, f0, f1, "frame is not identical to 1 second ago")
	assert.LessOrEqual(t, 7, countBlack(f0), "there are not enough black pixels between dots")
	assert.LessOrEqual(t, 7, countBlack(f1), "there are not enough black pixels between dots")
}

func countBlack(f frame) int {
	c := 0
	for i := range f {
		if f[i] == graphics.Black {
			c++
		}
	}
	return c
}
