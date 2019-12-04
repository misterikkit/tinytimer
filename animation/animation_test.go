package animation_test

import (
	"fmt"
	"image/color"
	"testing"
	"time"

	"github.com/misterikkit/tinytimer/animation"
	"github.com/misterikkit/tinytimer/graphics"
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

func count(f frame, c color.RGBA) int {
	ret := 0
	for i := range f {
		if f[i] == c {
			ret++
		}
	}
	return ret
}

func TestSpinner(t *testing.T) {
	s := animation.NewSpinner(color.RGBA{255, 255, 255, 0})
	var t0 time.Time
	s.Update(t0)
	f0 := newFrame()
	copy(f0, s.Frame())
	s.Update(t0.Add(time.Second))
	f1 := newFrame()
	copy(f1, s.Frame())
	assert.Equal(t, f0, f1, "frame is not identical to 1 second ago")
	assert.LessOrEqual(t, 7, count(f0, graphics.Black), "there are not enough black pixels between dots")
	assert.LessOrEqual(t, 7, count(f1, graphics.Black), "there are not enough black pixels between dots")
}

func TestLoader(t *testing.T) {
	yellow := color.RGBA{255, 255, 0, 0}
	cyan := color.RGBA{0, 255, 255, 0}
	var t0 time.Time
	tEnd := t0.Add(13 * time.Second)
	l := animation.NewLoader(yellow, cyan, t0, tEnd)
	// TODO: the below sequence is repetitive and could be captured in a loop.

	now := t0.Add(13 * time.Second / graphics.FrameSize * 5)
	done := l.Update(now)
	assert.Falsef(t, done, "loader should not be done at t=%v", now.Sub(t0))
	assert.Equal(t, 5, count(l.Frame(), yellow),
		"there should be 5 foreground pixels at t=%v", now.Sub(t0))
	assert.Equal(t, 19-2, count(l.Frame(), cyan)) // dot on the bg

	now = t0.Add(13 * time.Second / 2)
	done = l.Update(now)
	assert.Falsef(t, done, "loader should not be done at t=%v", now.Sub(t0))
	assert.Equal(t, 12-2, count(l.Frame(), yellow),
		"there should be 12-2 foreground pixels at t=%v", now.Sub(t0)) // dot on the fg
	assert.Equal(t, 12, count(l.Frame(), cyan))

	now = tEnd
	done = l.Update(now)
	assert.Truef(t, done, "loader should be done at t=%v", now.Sub(t0))
	assert.Equal(t, 24-2, count(l.Frame(), yellow),
		"there should be 24-2 foreground pixels at t=%v", now.Sub(t0)) // dot on the fg
	assert.Equal(t, 0, count(l.Frame(), cyan))

	now = tEnd.Add(2 * time.Second)
	done = l.Update(now)
	assert.Truef(t, done, "loader should be done at t=%v", now.Sub(t0))
	assert.Equal(t, 24-2, count(l.Frame(), yellow),
		"there should be 24-2 foreground pixels at t=%v", now.Sub(t0)) // dot on the fg
	assert.Equal(t, 0, count(l.Frame(), cyan))
}

func TestFlasher(t *testing.T) {
	var tEnd time.Time
	f := animation.NewFlasher(color.RGBA{0, 255, 0, 0}, tEnd)
	now := tEnd.Add(-time.Second)
	expected := newFrame()
	assert.Falsef(t, f.Update(now), "flasher should not be done at t=%v", now.Sub(tEnd))
	assert.Equalf(t, expected, frame(f.Frame()), "t=%v", now.Sub(tEnd))

	now = now.Add(time.Second / 4) // sin(tau/4)
	graphics.Fill(expected, color.RGBA{0, 255, 0, 0})
	assert.Falsef(t, f.Update(now), "flasher should not be done at t=%v", now.Sub(tEnd))
	assert.Equalf(t, expected, frame(f.Frame()), "t=%v", now.Sub(tEnd))

	now = now.Add(time.Second / 6)                   // sin(tau/3)
	graphics.Fill(expected, color.RGBA{0, 63, 0, 0}) // value is 1/2 squared
	assert.Falsef(t, f.Update(now), "flasher should not be done at t=%v", now.Sub(tEnd))
	assert.Equalf(t, expected, frame(f.Frame()), "t=%v", now.Sub(tEnd))

	now = tEnd
	expected = newFrame()
	assert.Truef(t, f.Update(now), "flasher should be done at t=%v", now.Sub(tEnd))
	assert.Equalf(t, expected, frame(f.Frame()), "t=%v", now.Sub(tEnd))

	now = now.Add(time.Second / 4)
	graphics.Fill(expected, color.RGBA{0, 255, 0, 0})
	assert.Truef(t, f.Update(now), "flasher should be done at t=%v", now.Sub(tEnd))
	assert.Equalf(t, expected, frame(f.Frame()), "t=%v", now.Sub(tEnd))
}
