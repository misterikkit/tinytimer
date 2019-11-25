package pong

import (
	"image/color"
	"time"

	"github.com/misterikkit/tinytimer/graphics"
	"github.com/misterikkit/tinytimer/input"
)

// measured in degrees
const (
	fieldLeft  = 270.0 + graphics.PixelWidth
	fieldRight = fieldLeft + 180.0 - 2*graphics.PixelWidth
	fieldMid   = (fieldLeft + fieldRight) / 2
	goalSize   = 2.0 * graphics.PixelWidth
	minSpeed   = 130.0 // per second
	maxSpeed   = 300.0 // per second
)

const maxPoints = 5 // and then you win

var (
	purple = graphics.Scale(color.RGBA{183, 95, 179, 0}, 0.5)
	green  = graphics.Scale(color.RGBA{105, 197, 91, 0}, 0.5)
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

type state uint8

const (
	volley state = iota
	score
	victory
)

// App is a game of pong.
type App struct {
	state           state
	lastStateChange time.Time
	p1, p2          player
	ball            ball
	scoreBG         graphics.Sprite
	frame           []color.RGBA
	lastUpdate      time.Time
	lastPlayerInput [2]time.Time
}

// New returns a fresh game of pong
func New(ui *input.Manager) *App {
	p := &App{
		p1: player{
			paddle: graphics.Sprite{
				Color:    purple,
				Position: graphics.Circ * 3 / 4,
				Size:     2.0 * graphics.PixelWidth,
			},
			scoreBar: graphics.Sprite{Color: purple},
		},
		p2: player{
			paddle: graphics.Sprite{
				Color:    green,
				Position: graphics.Circ / 4,
				Size:     2.0 * graphics.PixelWidth,
			},
			scoreBar: graphics.Sprite{Color: green},
		},
		ball: ball{
			speed: minSpeed,
			Sprite: graphics.Sprite{
				Color:    graphics.White,
				Position: fieldMid,
				Size:     1.0 * graphics.PixelWidth,
			},
		},
		scoreBG: graphics.Sprite{
			Color:    graphics.Scale(graphics.White, 0.5),
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

// Frame returns the current frame of the game.
func (p *App) Frame() []color.RGBA { return p.frame }

// Update computes the next frame of the game.
func (p *App) Update(now time.Time) {
	switch p.state {
	case volley:
		////////////////////
		// Update game state
		dt := now.Sub(p.lastUpdate)
		p.ball.Position += p.ball.speed * float32(dt.Seconds())
		if p.ball.Position >= fieldRight {
			p.ball.Position = fieldRight
			// p1 scores
			p.score(&p.p1)
		} else if p.ball.Position <= fieldLeft {
			p.ball.Position = fieldLeft
			// p2 scores
			p.score(&p.p2)
		}

	case score:
		/////////////////////
		// Start a new volley
		if now.Sub(p.lastStateChange) > time.Second {
			p.ball.Color = graphics.White
			if p.ball.speed < 0 {
				p.ball.speed = -minSpeed
			} else {
				p.ball.speed = minSpeed // does this need flip?
			}
			p.ball.Position = fieldMid // Maybe start further than halfway?
			p.state = volley
			p.lastStateChange = now
		}

	case victory:
		/////////////////////
		// Victory animation!
		progress := float32(now.Sub(p.lastStateChange)) / float32(2*time.Second)
		if progress > 1.0 {
			progress = 1.0
		}
		winner := &p.p1
		if p.p2.score >= maxPoints {
			winner = &p.p2
		}
		winner.paddle.Size = progress * 22 * graphics.PixelWidth
	}

	p.lastUpdate = now

	//////////////
	// Render game
	graphics.Fill(p.frame, graphics.Black)
	p.scoreBG.Render(p.frame)
	p.p1.paddle.Render(p.frame)
	p.p2.paddle.Render(p.frame)

	// Render higher scoring player last so victory animation looks right.
	pa, pb := &p.p1, &p.p2
	if p.p1.score < p.p2.score {
		pa, pb = pb, pa
	}
	pa.scoreBar.Render(p.frame)
	pb.scoreBar.Render(p.frame)

	p.ball.Render(p.frame)
}

func (p *App) score(player *player) {
	player.score++
	player.scoreBar.Size = float32(player.score) * graphics.PixelWidth
	if player == &p.p1 {
		player.scoreBar.Position = 180 + player.scoreBar.Size/2
	} else {
		player.scoreBar.Position = 180 - player.scoreBar.Size/2
	}
	p.ball.Color = graphics.Red
	if player.score >= maxPoints {
		p.state = victory
		p.ball.Size = 0
	} else {
		p.state = score
	}
	p.lastStateChange = time.Now() // plumb this in?
}

func (p *App) reset(input.Event) {
	p.state = volley
	p.lastStateChange = time.Now() // plumb this in?

	p.p1.score = 0
	p.p1.paddle.Size = 2 * graphics.PixelWidth
	p.p1.scoreBar.Size = 0
	p.p2.score = 0
	p.p2.paddle.Size = 2 * graphics.PixelWidth
	p.p2.scoreBar.Size = 0
	p.ball.speed = minSpeed
	p.ball.Position = fieldMid
	p.ball.Color = graphics.White
	p.ball.Size = graphics.PixelWidth
}

func (p *App) handle(e input.Event) {
	if p.state != volley {
		return
	}
	zone := struct{ l, r float32 }{}
	var lastInput time.Time
	switch e {
	case input.B_Rise:
		// player 1
		lastInput = p.lastPlayerInput[0]
		p.lastPlayerInput[0] = time.Now()
		zone.l = fieldLeft
		zone.r = fieldLeft + goalSize
		// Make sure player can't "catch" the ball as it leaves their zone.
		if p.ball.speed > 0 {
			return
		}
	case input.C_Rise:
		// player 2
		lastInput = p.lastPlayerInput[1]
		p.lastPlayerInput[1] = time.Now()
		zone.r = fieldRight
		zone.l = fieldRight - goalSize
		// Make sure player can't "catch" the ball as it leaves their zone.
		if p.ball.speed < 0 {
			return
		}
	}

	// Ignore user input if they've pressed their button too recently (spam).
	if time.Since(lastInput) < 100*time.Millisecond {
		return
	}
	if p.ball.Position >= zone.l && p.ball.Position <= zone.r {
		p.ball.speed = -p.ball.speed
	}
	// TODO: variable speed
}
