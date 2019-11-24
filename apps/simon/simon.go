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

type App struct {
	frame    []color.RGBA
	sequence []token
	sprites  []graphics.Sprite
}

func New(ui *input.Manager) *App {
	return &App{
		frame: make([]color.RGBA, 24),
		sprites: []graphics.Sprite{
			{Size: 3 * graphics.PixelWidth, Position: graphics.Circ * 0 / 3, Color: graphics.Red},
			{Size: 3 * graphics.PixelWidth, Position: graphics.Circ * 1 / 3, Color: graphics.Green},
			{Size: 3 * graphics.PixelWidth, Position: graphics.Circ * 2 / 3, Color: graphics.Blue},
		},
	}
}

func (s *App) Update(time.Time) {
	graphics.Fill(s.frame, graphics.Black)
	for i := range s.sprites {
		s.sprites[i].Render(s.frame)
	}
}

func (s *App) Frame() []color.RGBA { return s.frame }
