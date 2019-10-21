package main

import (
	"image/color"
	"math"
	"time"
)

// updateFn updates an animation based on the current time, and returns true if
// the animation is complete.
type updateFn = func(time.Time) bool

// handle is an animation interface, but I am worried about using actual
// interfaces in TinyGo.
type handle struct {
	f      *Frame
	update updateFn
}

// spinner is the idle animation with 7 spinning dots
type spinner struct {
	f    Frame
	dots []sprite
}

const spinCount = 7

func newSpinner() spinner {
	const size = 0.8 * PixelWidth
	var color = color.RGBA{0x32, 0x6C, 0xE5, 0}
	s := spinner{
		f:    newFrame(),
		dots: make([]sprite, 0, spinCount),
	}
	for i := 0; i < spinCount; i++ {
		s.dots = append(s.dots, sprite{Size: size, Color: color})
	}
	return s
}

func (s *spinner) update(now time.Time) bool {
	var period = scaleDuration(spinCount * time.Second / 2)
	const divide = Tau / spinCount
	s.f.fill(Black)

	// compute fraction through the period
	progress := float32(now.Sub(now.Truncate(period))) / float32(period)
	for i := range s.dots {
		s.dots[i].Position = float32(math.Mod(float64(Tau*progress+float32(i)*divide), Tau))
		s.dots[i].Render(s.f)
	}
	return false
}

type loader struct {
	f          Frame
	bar, dot   sprite
	bg         color.RGBA
	start, end time.Time
	done       bool
}

func newLoader(c color.RGBA, start, end time.Time) loader {
	return loader{
		f:     newFrame(),
		bar:   sprite{Color: c},
		dot:   sprite{Color: White, Size: PixelWidth},
		start: start,
		end:   end,
	}
}

func (l *loader) update(now time.Time) bool {
	if l.done {
		return true
	}
	l.f.fill(l.bg)
	progress := float32(1.0)
	if now.Before(l.end) {
		progress = float32(now.Sub(l.start)) / float32(l.end.Sub(l.start))
	} else {
		l.done = true
	}
	l.bar.Size = Tau * progress
	l.bar.Position = l.bar.Size / 2.0

	elapsed := float32(now.Sub(l.start).Seconds())
	// This is reverse scaled since it is supposed to match real second ticks.
	elapsed *= TimeScale
	l.dot.Position = elapsed * Tau

	l.bar.Render(l.f)
	l.dot.Render(l.f)
	return false
}

type fader struct {
	f          Frame
	from, to   handle
	start, end time.Time
}

func newFader(start, end time.Time) fader {
	return fader{
		f:     newFrame(),
		start: start,
		end:   end,
	}
}

func (f *fader) update(now time.Time) bool {
	// TODO: There is a weird blip at the beginning of each fade
	if now.After(f.end) {
		done := f.to.update(now)
		copy(f.f, *f.to.f)
		return done
	}

	f.from.update(now)
	f.to.update(now)
	progress := float32(now.Sub(f.start)) / float32(f.end.Sub(f.start))
	for i := range f.f {
		f.f[i] = add(
			scale((*f.from.f)[i], 1.0-progress),
			scale((*f.to.f)[i], progress),
		)
	}

	return false
}

type flasher struct {
	f   Frame
	c   color.RGBA
	end time.Time
}

func newFlasher(c color.RGBA, end time.Time) flasher {
	return flasher{newFrame(), c, end}
}

func (f *flasher) update(now time.Time) bool {
	progress := f.end.Sub(now).Seconds() * Tau
	s := float32(math.Sin(progress))
	s = s * s // stay smooth. stay positive
	val := scale(f.c, s)
	f.f.fill(val)
	return now.After(f.end)
}
