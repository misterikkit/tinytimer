package main

import (
	"image/color"
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

func newSpinner() spinner {
	const size = 1.5 * PixelWidth
	var color = color.RGBA{0x32, 0x6C, 0xE5, 0}
	s := spinner{
		f:    newFrame(),
		dots: make([]sprite, 0, 7),
	}
	for i := 0; i < 7; i++ {
		s.dots = append(s.dots, sprite{Size: size, Color: color})
	}
	return s
}

func (s *spinner) update(now time.Time) bool {
	const (
		period = 10 * time.Second
		divide = Tau / 7.0
	)
	s.f.fill(Black)

	// compute fraction through the period
	progress := float32(now.Sub(now.Truncate(period))) / float32(period)
	for i := range s.dots {
		s.dots[i].Position = Tau*progress + float32(i)*divide
		s.dots[i].Render(s.f)
	}
	return false
}

type loader struct {
	f          Frame
	s          sprite
	bg         color.RGBA
	start, end time.Time
	done       bool
}

func newLoader(c color.RGBA, start, end time.Time) loader {
	return loader{
		f:     newFrame(),
		s:     sprite{Color: c},
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
	l.s.Size = Tau * progress
	l.s.Position = l.s.Size / 2.0

	l.s.Render(l.f)
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
	if now.After(f.end) {
		done := f.to.update(now)
		copy(f.f, *f.to.f)
		return done
	}

	f.from.update(now)
	f.to.update(now)
	progress := float32(now.Sub(f.start)) / float32(f.end.Sub(f.start))
	for i := range f.f {
		cFrom := scale((*f.from.f)[i], 1.0-progress)
		cTo := scale((*f.to.f)[i], progress)
		f.f[i] = color.RGBA{
			R: cFrom.R + cTo.R,
			G: cFrom.G + cTo.G,
			B: cFrom.B + cTo.B,
		}
	}

	return false
}
