package simon

import (
	"image/color"
	"time"

	"github.com/misterikkit/tinytimer/graphics"
	"github.com/misterikkit/tinytimer/input"
)

type token uint8

const (
	a token = 0
	b token = 1
	c token = 2
)

type state uint8

const (
	intro state = iota
	displaying
	userInput
)

type App struct {
	// graphics
	frame   []color.RGBA
	sprites []graphics.Sprite
	bg      []graphics.Sprite
	// in input mode, indicates which buttons are down
	echo []bool

	// the correct sequence
	sequence []token
	// collected user input so far
	collectedInput []token

	state           state
	lastStateChange time.Time
}

func New(ui *input.Manager) *App {
	s := &App{
		frame: make([]color.RGBA, 24),
		sprites: []graphics.Sprite{
			{Size: 3 * graphics.PixelWidth, Position: graphics.Circ * 0 / 3, Color: graphics.Red},
			{Size: 3 * graphics.PixelWidth, Position: graphics.Circ * 2 / 3, Color: graphics.Green},
			{Size: 3 * graphics.PixelWidth, Position: graphics.Circ * 1 / 3, Color: graphics.Blue},
		},
		bg: []graphics.Sprite{
			{Size: graphics.PixelWidth, Position: graphics.Circ * 1 / 6, Color: graphics.White},
			{Size: graphics.PixelWidth, Position: graphics.Circ * 3 / 6, Color: graphics.White},
			{Size: graphics.PixelWidth, Position: graphics.Circ * 5 / 6, Color: graphics.White},
		},
		echo: make([]bool, 3),

		state:           userInput,
		lastStateChange: time.Now(),
	}
	ui.AddHandler(input.A_Fall, s.handleInput)
	ui.AddHandler(input.B_Fall, s.handleInput)
	ui.AddHandler(input.C_Fall, s.handleInput)
	ui.AddHandler(input.A_Fall, s.handleEcho)
	ui.AddHandler(input.B_Fall, s.handleEcho)
	ui.AddHandler(input.C_Fall, s.handleEcho)
	ui.AddHandler(input.A_Rise, s.handleEcho)
	ui.AddHandler(input.B_Rise, s.handleEcho)
	ui.AddHandler(input.C_Rise, s.handleEcho)
	return s
}

func (s *App) Update(time.Time) {
	graphics.Fill(s.frame, graphics.Black)
	for i := range s.bg {
		s.bg[i].Render(s.frame)
	}
	if s.state == userInput {
		for i := range s.echo {
			if !s.echo[i] {
				continue
			}
			s.sprites[i].Render(s.frame)
		}
	}

}

func (s *App) Frame() []color.RGBA { return s.frame }

func (s *App) handleInput(e input.Event) {
	if s.state != userInput {
		return
	}
	switch e {
	case input.A_Fall:
		s.collectedInput = append(s.collectedInput, a)
	case input.B_Fall:
		s.collectedInput = append(s.collectedInput, b)
	case input.C_Fall:
		s.collectedInput = append(s.collectedInput, c)
	}
}

func (s *App) handleEcho(e input.Event) {
	switch e {
	case input.A_Fall:
		s.echo[a] = false
	case input.B_Fall:
		s.echo[b] = false
	case input.C_Fall:
		s.echo[c] = false

	case input.A_Rise:
		s.echo[a] = true
	case input.B_Rise:
		s.echo[b] = true
	case input.C_Rise:
		s.echo[c] = true
	}
}
