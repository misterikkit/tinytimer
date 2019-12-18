package arcade

import (
	"image/color"
	"time"

	"github.com/misterikkit/tinytimer/graphics"
	"github.com/misterikkit/tinytimer/input"
)

const (
	minSpeed   float32 = 90  // degrees per sec
	maxSpeed   float32 = 360 // degrees per sec
	zoneSize   float32 = graphics.PixelWidth
	leftBound          = -zoneSize
	rightBound         = graphics.Circ + zoneSize
	maxScore           = 9
)

type state uint8

const (
	game state = iota
	victory
)

var paddleColors = []color.RGBA{
	color.RGBA{0xFF, 0, 0, 0},       //red
	color.RGBA{0xFF, 0x7F, 0, 0},    //orange
	color.RGBA{0xFF, 0xFF, 0, 0},    //yellow
	color.RGBA{0, 0xFF, 0, 0},       //green
	color.RGBA{0, 0xFF, 0xFF, 0},    //cyan
	color.RGBA{0, 0x7F, 0xFF, 0},    //light blue
	color.RGBA{0, 0, 0xFF, 0},       //blue
	color.RGBA{0x7F, 0, 0xFF, 0},    //indigo
	color.RGBA{0xFF, 0, 0xFF, 0},    //purple
	color.RGBA{0xFF, 0x14, 0x93, 0}, //pink-ish
}

type App struct {
	frame []color.RGBA
	state
	ball, paddle graphics.Sprite
	speed        float32
	bounces      int8
	lastUpdate   time.Time
	lastInput    time.Time
}

func New(ui *input.Manager) *App {
	a := &App{
		frame:  make([]color.RGBA, graphics.FrameSize),
		ball:   graphics.Sprite{Color: graphics.White, Size: graphics.PixelWidth},
		paddle: graphics.Sprite{Color: graphics.Red, Size: 2 * zoneSize},
	}
	ui.AddHandler(input.A_Rise, a.handle)
	ui.AddHandler(input.B_Rise, a.handle)
	ui.AddHandler(input.C_Rise, a.handle)
	return a
}

func (a *App) Frame() []color.RGBA { return a.frame }

func (a *App) Reset(now time.Time) {
	a.paddle.Size = 2 * zoneSize
	a.fail(now)
	if a.speed < 0 {
		a.ball.Position = graphics.Circ
	} else {
		a.ball.Position = 0
	}
}

func (a *App) fail(now time.Time) {
	a.bounces = 0
	if a.speed < 0 {
		a.speed = -minSpeed
	} else {
		a.speed = minSpeed
	}
	if a.ball.Position > graphics.Circ {
		a.ball.Position -= graphics.Circ
	}
	if a.ball.Position < 0 {
		a.ball.Position += graphics.Circ
	}
	a.state = game
	a.lastUpdate = now
}

func (a *App) Update(now time.Time) {
	a.ball.Position += float32(now.Sub(a.lastUpdate).Seconds()) * a.speed
	a.paddle.Color = paddleColors[int(a.bounces)%len(paddleColors)]

	switch a.state {
	case game:
		if a.ball.Position < leftBound || a.ball.Position > rightBound {
			// LOSE!
			a.fail(now)
		}
	case victory:
		if a.ball.Position > graphics.Circ {
			a.ball.Position -= graphics.Circ
		}
		if a.ball.Position < 0 {
			a.ball.Position += graphics.Circ
		}
		if a.paddle.Size < graphics.Circ {
			a.paddle.Size += float32(now.Sub(a.lastUpdate).Seconds()) * minSpeed
		}
	}

	graphics.Fill(a.frame, graphics.Black)
	a.paddle.Render(a.frame)
	a.ball.Render(a.frame)
	a.lastUpdate = now
}

func (a *App) handle(input.Event) {
	now := time.Now() // TODO: plumb this in
	if now.Sub(a.lastInput) < 250*time.Millisecond {
		return
	}
	a.lastInput = now
	if a.state == victory {
		if a.paddle.Size > graphics.Circ/2 {
			a.Reset(now)
		}
		return
	}
	// miss cases
	if a.speed > 0 && a.ball.Position < (graphics.Circ-zoneSize) {
		return
	}
	if a.speed < 0 && a.ball.Position > zoneSize {
		return
	}
	// hit!
	a.bounces++
	if a.bounces >= maxScore {
		a.state = victory
	}
	sign := float32(1.0)
	if a.speed > 0 {
		sign = -1.0
	}
	a.speed = (minSpeed + (maxSpeed-minSpeed)*float32(a.bounces)/float32(maxScore)) * sign
}
