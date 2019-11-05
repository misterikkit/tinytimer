package animation

import (
	"image/color"
	"math"
	"time"

	"github.com/misterikkit/tinytimer/graphics"
	"github.com/misterikkit/tinytimer/hack"
)

// updateFn updates an animation based on the current time, and returns true if
// the animation is complete.
type updateFn = func(time.Time) bool

type Handle struct {
	Frame  *[]color.RGBA
	Update updateFn
}

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

var period = hack.ScaleDuration(time.Second * spinnerCount)

// Update computes the current frame of animation.
func (s *Spinner) Update(now time.Time) bool {
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

type Loader struct {
	Frame      []color.RGBA
	bar, dot   graphics.Sprite
	BG         color.RGBA
	start, end time.Time
	done       bool
}

func NewLoader(c color.RGBA, start, end time.Time) Loader {
	return Loader{
		Frame: make([]color.RGBA, graphics.FrameSize),
		bar:   graphics.Sprite{Color: c},
		dot:   graphics.Sprite{Color: graphics.White, Size: graphics.PixelWidth},
		start: start,
		end:   end,
	}
}

func (l *Loader) Update(now time.Time) bool {
	if l.done {
		return true
	}
	graphics.Fill(l.Frame, l.BG)
	progress := float32(1.0)
	if now.Before(l.end) {
		progress = float32(now.Sub(l.start)) / float32(l.end.Sub(l.start))
	} else {
		l.done = true
	}
	l.bar.Size = graphics.Circ * progress
	l.bar.Position = l.bar.Size / 2.0

	elapsed := float32(now.Sub(l.start).Seconds())
	// This is reverse scaled since it is supposed to match real second ticks.
	elapsed *= hack.TimeScale
	l.dot.Position = elapsed * graphics.Circ

	l.bar.Render(l.Frame)
	l.dot.Render(l.Frame)
	return false
}

type Flasher struct {
	Frame []color.RGBA
	c     color.RGBA
	end   time.Time
}

func NewFlasher(c color.RGBA, end time.Time) Flasher {
	return Flasher{make([]color.RGBA, graphics.FrameSize), c, end}
}

func (f *Flasher) Update(now time.Time) bool {
	progress := f.end.Sub(now).Seconds() * math.Pi * 2
	s := float32(math.Sin(progress))
	s = s * s // stay smooth. stay positive
	val := graphics.Scale(f.c, s)
	graphics.Fill(f.Frame, val)
	return now.After(f.end)
}

type Fader struct {
	Frame      []color.RGBA
	From, To   Handle
	start, end time.Time
}

func NewFader(start, end time.Time) Fader {
	return Fader{
		Frame: make([]color.RGBA, graphics.FrameSize),
		start: start,
		end:   end,
	}
}

func (f *Fader) Update(now time.Time) bool {
	// TODO: There is a weird blip at the beginning of each fade
	if now.After(f.end) {
		done := f.To.Update(now)
		copy(f.Frame, *f.To.Frame) // Unfortunate that copy is required
		return done
	}

	f.From.Update(now)
	f.To.Update(now)
	progress := float32(now.Sub(f.start)) / float32(f.end.Sub(f.start))
	for i := range f.Frame {
		f.Frame[i] = graphics.Add(
			graphics.Scale((*f.From.Frame)[i], 1.0-progress),
			graphics.Scale((*f.To.Frame)[i], progress),
		)
	}

	return false
}
