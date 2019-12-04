package animation

import (
	"image/color"
	"math"
	"time"

	"github.com/misterikkit/tinytimer/graphics"
	"github.com/misterikkit/tinytimer/hack"
)

// Interface provides a frame of pixels that changes over time.
type Interface interface {
	// Frame returns the animation's current pixel values.
	Frame() []color.RGBA
	// Update updates an animation's frame based on the current time, and returns
	// true if the animation is complete.
	Update(time.Time) bool
}

type spinner struct {
	frame []color.RGBA
	dots  []graphics.Sprite
}

const spinnerCount = 7

var size = graphics.PixelWidth * 1.2
var divide = graphics.Circ / spinnerCount

// NewSpinner initializes a spinner animation.
func NewSpinner(c color.RGBA) Interface {
	s := spinner{
		frame: make([]color.RGBA, graphics.FrameSize),
		dots:  make([]graphics.Sprite, 0, spinnerCount),
	}
	for i := 0; i < spinnerCount; i++ {
		s.dots = append(s.dots, graphics.Sprite{Size: size, Color: c})
	}
	return &s
}

func (s *spinner) Frame() []color.RGBA { return s.frame }

var period = hack.ScaleDuration(time.Second * spinnerCount)

// Update computes the current frame of animation.
func (s *spinner) Update(now time.Time) bool {
	graphics.Fill(s.frame, graphics.Black)

	// compute fraction through the period
	elapsed := float32(now.Sub(now.Truncate(period)).Nanoseconds())
	progress := elapsed / float32(period.Nanoseconds())
	for i := range s.dots {
		// The value of Position has an upper bound of `2*Circ`. The max progress is
		// 1.0 and divide*spinnerCount==Circ.
		s.dots[i].Position = graphics.Circ*progress + divide*float32(i)
		s.dots[i].Render(s.frame)
	}
	return false
}

type loader struct {
	frame      []color.RGBA
	bar, dot   graphics.Sprite
	bg         color.RGBA
	start, end time.Time
	done       bool
}

// NewLoader initializes a loader animation.
func NewLoader(fg, bg color.RGBA, start, end time.Time) Interface {
	return &loader{
		frame: make([]color.RGBA, graphics.FrameSize),
		bar:   graphics.Sprite{Color: fg},
		bg:    bg,
		dot:   graphics.Sprite{Color: graphics.White, Size: graphics.PixelWidth},
		start: start,
		end:   end,
	}
}

func (l *loader) Frame() []color.RGBA { return l.frame }

func (l *loader) Update(now time.Time) bool {
	if l.done {
		return true
	}
	graphics.Fill(l.frame, l.bg)
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

	l.bar.Render(l.frame)
	l.dot.Render(l.frame)
	return l.done
}

type flasher struct {
	frame []color.RGBA
	c     color.RGBA
	end   time.Time
}

// NewFlasher initializes a flasher animation.
func NewFlasher(c color.RGBA, end time.Time) Interface {
	return &flasher{make([]color.RGBA, graphics.FrameSize), c, end}
}

func (f *flasher) Frame() []color.RGBA { return f.frame }

func (f *flasher) Update(now time.Time) bool {
	progress := f.end.Sub(now).Seconds() * math.Pi * 2
	s := float32(math.Sin(progress))
	s = s * s // stay smooth. stay positive
	val := graphics.Scale(f.c, s)
	graphics.Fill(f.frame, val)
	return !now.Before(f.end) // This is double negated to return true on the exact frame.
}

type fader struct {
	frame      []color.RGBA
	from, to   Interface
	start, end time.Time
}

// NewFader initializes a fader animation that blends from and to.
func NewFader(from, to Interface, start, end time.Time) Interface {
	return &fader{
		frame: make([]color.RGBA, graphics.FrameSize),
		from:  from,
		to:    to,
		start: start,
		end:   end,
	}
}

func (f *fader) Frame() []color.RGBA { return f.frame }

func (f *fader) Update(now time.Time) bool {
	if now.After(f.end) {
		done := f.to.Update(now)
		copy(f.frame, f.to.Frame()) // Unfortunate that copy is required
		// TODO: f.frame = f.to.Frame()
		return done
	}

	f.from.Update(now)
	f.to.Update(now)
	progress := float32(now.Sub(f.start)) / float32(f.end.Sub(f.start))
	for i := range f.frame {
		f.frame[i] = graphics.Add(
			graphics.Scale(f.from.Frame()[i], 1.0-progress),
			graphics.Scale(f.to.Frame()[i], progress),
		)
	}

	return false
}
