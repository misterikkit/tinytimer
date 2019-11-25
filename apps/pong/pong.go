package pong

import (
	"image/color"
	"time"

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
	// TODO: increase chroma, decrease luminescence.
	purple = color.RGBA{172, 78, 200, 0}
	green  = color.RGBA{151, 211, 75, 0}
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
			speed: 150.0,
			Sprite: graphics.Sprite{
				Color:    graphics.White,
				Position: (fieldLeft + fieldRight) / 2,
				Size:     1.0 * graphics.PixelWidth,
			},
		},
		scoreBG: graphics.Sprite{
			Color:    graphics.Scale(graphics.White, 0.25),
			Position: 180,
			Size:     10 * graphics.PixelWidth,
		},
		frame:      make([]color.RGBA, graphics.FrameSize),
		lastUpdate: time.Now(),
	}
	ui.AddHandler(input.A_Fall, p.reset)
	ui.AddHandler(input.B_Rise, p.handle)
	ui.AddHandler(input.C_Rise, p.handle)
	return p
}

// Update computes the next frame of the game.
func (p *App) Update(now time.Time) {
	////////////////////
	// Update game state
	dt := now.Sub(p.lastUpdate)
	p.ball.Position += p.ball.speed * float32(dt.Seconds())
	if p.ball.RightEdge() >= p.p2.paddle.LeftEdge() {
		p.ball.Position = p.p2.paddle.LeftEdge() - p.ball.Size/2
		p.ball.speed = -p.ball.speed // this is temporary
		// p1 scores
		p.p1.score++
		p.p1.scoreBar.Size = float32(p.p1.score) * graphics.PixelWidth
		p.p1.scoreBar.Position = 180 + p.p1.scoreBar.Size/2
		// TODO: check for win
	}
	if p.ball.LeftEdge() <= p.p1.paddle.RightEdge() {
		p.ball.Position = p.p1.paddle.RightEdge() + p.ball.Size/2
		p.ball.speed = -p.ball.speed // this is temporary
		// p2 scores
		p.p2.score++
		p.p2.scoreBar.Size = float32(p.p2.score) * graphics.PixelWidth
		p.p2.scoreBar.Position = 180 - p.p2.scoreBar.Size/2
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

func (p *App) reset(input.Event) {
	p.p1.score = 0
	p.p2.score = 0
	p.ball.Position = (fieldLeft + fieldRight) / 2
	p.ball.speed = 10.0
}

func (p *App) handle(e input.Event) {
	zone := struct{ l, r float32 }{}
	switch e {
	case input.B_Rise:
		// player 1
		zone.l = fieldLeft
		zone.r = fieldLeft + goalSize
	case input.C_Rise:
		// player 2
		zone.r = fieldRight
		zone.l = fieldRight - goalSize
	}
}
