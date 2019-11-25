package pong

import (
	"image/color"
	"time"

	colorful "github.com/lucasb-eyer/go-colorful"
	"github.com/misterikkit/tinytimer/graphics"
	"github.com/misterikkit/tinytimer/input"
)

// measured in degrees
const (
	fieldLeft  = 270.0
	fieldRight = fieldLeft + 180.0
	goalSize   = 2.0 * graphics.PixelWidth
)

const maxPoints = 7 // and then you win

var (
	purple = func() color.RGBA { r, g, b := colorful.Hcl(320, 73.5, 50).RGB255(); return color.RGBA{r, g, b, 0} }()
	green  = func() color.RGBA { r, g, b := colorful.Hcl(124, 71.5, 78).RGB255(); return color.RGBA{r, g, b, 0} }()
)

type player struct {
	score    uint8
	paddle   graphics.Sprite
	scoreBar graphics.Sprite
}

type ball struct {
	graphics.Sprite
	speed float32 // degrees per second
}

// App is a game of pong.
type App struct {
	p1, p2     player
	ball       ball
	scoreBG    graphics.Sprite
	frame      []color.RGBA
	lastUpdate time.Time
}

// New returns a fresh game of pong
func New(ui *input.Manager) *App {
	p := &App{
		p1: player{
			paddle: graphics.Sprite{
				Color:    purple,
				Position: fieldLeft,
				Size:     2.0 * graphics.PixelWidth,
			},
			scoreBar: graphics.Sprite{Color: purple},
		},
		p2: player{
			paddle: graphics.Sprite{
				Color:    green,
				Position: fieldRight,
				Size:     2.0 * graphics.PixelWidth,
			},
			scoreBar: graphics.Sprite{Color: green},
		},
		ball: ball{
			speed: 10.0,
			Sprite: graphics.Sprite{
				Color:    graphics.White,
				Position: (fieldLeft + fieldRight) / 2,
				Size:     1.0 * graphics.PixelWidth,
			},
		},
		scoreBG: graphics.Sprite{
			Color:    graphics.White,
			Position: 180,
			Size:     10 * graphics.PixelWidth,
		},
		frame:      make([]color.RGBA, graphics.FrameSize),
		lastUpdate: time.Now(),
	}
	return p
}

// Update computes the next frame of the game.
func (p *App) Update(now time.Time) {
	////////////////////
	// Update game state
	dt := now.Sub(p.lastUpdate)
	p.ball.Position += p.ball.speed * float32(dt.Seconds())
	if p.ball.Position >= fieldRight {
		p.ball.Position = fieldRight /*and score*/
	}
	if p.ball.Position <= fieldLeft {
		p.ball.Position = fieldLeft /*and score*/
	}
	p.lastUpdate = now

	//////////////
	// Render game
	graphics.Fill(p.frame, graphics.Black)
	p.scoreBG.Render(p.frame)
	p.p1.paddle.Render(p.frame)
	p.p1.scoreBar.Render(p.frame)
	p.p2.paddle.Render(p.frame)
	p.p2.scoreBar.Render(p.frame)
	p.ball.Render(p.frame)
}

// Frame returns the current frame of the game.
func (p *App) Frame() []color.RGBA { return p.frame }
