package simon

import (
	"image/color"
	"math"
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
	correct
	incorrect
)

var (
	bgWhite = []graphics.Sprite{
		{Size: graphics.PixelWidth, Position: graphics.Circ * 1 / 6, Color: graphics.White},
		{Size: graphics.PixelWidth, Position: graphics.Circ * 3 / 6, Color: graphics.White},
		{Size: graphics.PixelWidth, Position: graphics.Circ * 5 / 6, Color: graphics.White},
	}
	bgGreen = []graphics.Sprite{
		{Size: graphics.PixelWidth, Position: graphics.Circ * 1 / 6, Color: graphics.Green},
		{Size: graphics.PixelWidth, Position: graphics.Circ * 3 / 6, Color: graphics.Green},
		{Size: graphics.PixelWidth, Position: graphics.Circ * 5 / 6, Color: graphics.Green},
	}
	bgRed = []graphics.Sprite{
		{Size: graphics.PixelWidth, Position: graphics.Circ * 1 / 6, Color: graphics.Red},
		{Size: graphics.PixelWidth, Position: graphics.Circ * 3 / 6, Color: graphics.Red},
		{Size: graphics.PixelWidth, Position: graphics.Circ * 5 / 6, Color: graphics.Red},
	}
)

type App struct {
	// graphics
	frame   []color.RGBA
	sprites []graphics.Sprite
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
	initRandom()
	s := &App{
		frame: make([]color.RGBA, 24),
		sprites: []graphics.Sprite{
			{Size: 3 * graphics.PixelWidth, Position: graphics.Circ * 0 / 3, Color: graphics.Red},
			{Size: 3 * graphics.PixelWidth, Position: graphics.Circ * 2 / 3, Color: graphics.Green},
			{Size: 3 * graphics.PixelWidth, Position: graphics.Circ * 1 / 3, Color: graphics.Blue},
		},
		echo: make([]bool, 3),

		state:           intro,
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

func (s *App) Reset() {
	s.sequence = nil
	s.collectedInput = nil
	s.state = intro
	s.lastStateChange = time.Now()
}

func (s *App) Frame() []color.RGBA { return s.frame }

func (s *App) Update(now time.Time) {
	graphics.Fill(s.frame, graphics.Black)
	s.doBG(now)

	switch s.state {
	case intro:
		for i := range s.sprites {
			s.sprites[i].Render(s.frame)
		}
		if now.Sub(s.lastStateChange) > time.Second {
			s.state = correct
			s.lastStateChange = now
		}
	case displaying:
		s.doDisplay(now)
	case userInput:
		s.doEcho()
	case correct:
		if now.Sub(s.lastStateChange) > time.Second {
			s.collectedInput = nil
			s.sequence = append(s.sequence, token(randU32()%3))
			s.state = displaying
			s.lastStateChange = now
		}

	case incorrect:
		if now.Sub(s.lastStateChange) > time.Second {
			s.collectedInput = nil
			s.sequence = nil
			s.state = intro
			s.lastStateChange = now
		}
	}
}

func (s *App) doBG(now time.Time) {
	var bg []graphics.Sprite
	switch s.state {
	case correct:
		bg = bgGreen
	case incorrect:
		bg = bgRed
	default:
		bg = bgWhite
	}

	const period = 7 * time.Second
	progress := now.Sub(now.Truncate(period)).Seconds() / period.Seconds()
	size := float32(math.Sin(progress * math.Pi)) // use pi instead of 2*pi because we will square the result, halving the period
	size = graphics.PixelWidth * (size*size + 1.0)
	for i := range bg {
		bg[i].Size = size
		bg[i].Render(s.frame)
	}
}

func (s *App) doDisplay(now time.Time) {
	dwell := 750 * time.Millisecond // TODO: speed up as sequence gets longer
	const idle = 100 * time.Millisecond
	tokenDuration := dwell + idle
	progress := now.Sub(s.lastStateChange) // negative value should be impossible
	offset := int(progress.Seconds() / tokenDuration.Seconds())

	if offset >= len(s.sequence) {
		s.state = userInput
		s.lastStateChange = now
		return
	}

	// Brief blank period between tokens
	if (progress - progress.Truncate(tokenDuration)) < idle {
		return
	}
	s.sprites[s.sequence[offset]].Render(s.frame)
}

func (s *App) doEcho() {
	for i := range s.echo {
		if s.echo[i] {
			s.sprites[i].Render(s.frame)
		}
	}
}

// handleInput appends a user input to the collected input.
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

	// Check if last input was correct
	i := len(s.collectedInput) - 1
	if s.collectedInput[i] != s.sequence[i] {
		s.state = incorrect
		s.lastStateChange = time.Now() // TODO: plumb this in?
		return
	}
	// Check if we reached the end of sequence
	if len(s.collectedInput) == len(s.sequence) {
		s.state = correct
		s.lastStateChange = time.Now() // TODO: plumb this in?
	}
}

// handleEcho keeps track of current button state, since plumbing the pollable
// inputs into this app is not desirable.
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
