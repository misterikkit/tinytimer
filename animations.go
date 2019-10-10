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
	s.f.reset()

	// compute fraction through the period
	progress := float32(now.Sub(now.Truncate(period))) / float32(period)
	for i := range s.dots {
		s.dots[i].Position = Tau*progress + float32(i)*divide
		s.dots[i].Render(s.f)
	}
}
