package main

import (
	"image/color"
	"time"
)

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

func (s *spinner) update(now time.Time) {
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
}

type loader struct {
	f          Frame
	s          sprite
	color      color.RGBA
	start, end time.Time
	done       bool
}

func newLoader(c color.RGBA) loader {
	return loader{
		f: newFrame(),
		s: sprite{Color: c},
	}
}

func (l *loader) update(now time.Time) {
	if l.done {
		return
	}
	l.f.fill(Black)
	progress := float32(1.0)
	if now.Before(l.end) {
		progress = float32(now.Sub(l.start)) / float32(l.end.Sub(l.start))
	} else {
		l.done = true
	}
	l.s.Size = Tau * progress
	l.s.Position = l.s.Size / 2.0

	l.s.Render(l.f)
}
