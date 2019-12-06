package colorpicker

import (
	"image/color"
	"math"
	"time"

	"github.com/misterikkit/tinytimer/graphics"
	"github.com/misterikkit/tinytimer/input"
)

type colorComponent uint8

const (
	red    colorComponent = 0
	green  colorComponent = 1
	blue   colorComponent = 2
	bright colorComponent = 3

	componentCount = 4
)

const tick = 100 * time.Millisecond

type App struct {
	frame            []color.RGBA
	picker           color.NRGBA
	currentComponent colorComponent
	buttons          struct{ a, b, c bool }
	lastTick         time.Time
}

func New(ui *input.Manager) *App {
	p := &App{
		frame: make([]color.RGBA, graphics.FrameSize),
	}
	ui.AddHandler(input.A_Fall, p.handle)
	ui.AddHandler(input.B_Fall, p.handle)
	ui.AddHandler(input.C_Fall, p.handle)
	ui.AddHandler(input.A_Rise, p.handle)
	ui.AddHandler(input.B_Rise, p.handle)
	ui.AddHandler(input.C_Rise, p.handle)
	return p
}

func (p *App) Frame() []color.RGBA { return p.frame }

func (p *App) Update(now time.Time) {
	rainbowDot(now, &p.frame[graphics.FrameSize-1])
	rainbowDot(now, &p.frame[8])

	// Count number of ticks since
	if now.Sub(p.lastTick) > tick {
		p.updateColorPicker()
		p.lastTick = now
	}

	//////////////
	// Render data
	r, g, b, a := p.picker.RGBA()
	actualColor := color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
	s := graphics.Sprite{
		Color:    actualColor,
		Position: graphics.Circ * 3 / 4,
		Size:     6 * graphics.PixelWidth,
	}
	s.Render(p.frame)
	switch p.currentComponent {
	case red:
		renderByte(p.frame, actualColor.R, graphics.Red)
	case green:
		renderByte(p.frame, actualColor.G, graphics.Green)
	case blue:
		renderByte(p.frame, actualColor.B, graphics.Blue)
	case bright:
		renderByte(p.frame, actualColor.A, graphics.White)
	}
}

func (p *App) Reset() { p.currentComponent = red; p.picker = color.NRGBA{31, 63, 127, 255} }

func renderByte(frame []color.RGBA, value uint8, color color.RGBA) {
	for i := uint8(0); i < 8; i++ {
		if value&(1<<i) != 0 {
			frame[i] = color
		} else {
			frame[i] = graphics.Black
		}
	}
}

func (p *App) updateColorPicker() {
	update := func(v uint8) uint8 { return v }
	if p.buttons.b && !p.buttons.c {
		update = func(v uint8) uint8 { return v - 1 }
	}
	if !p.buttons.b && p.buttons.c {
		update = func(v uint8) uint8 { return v + 1 }
	}
	switch p.currentComponent {
	case red:
		p.picker.R = update(p.picker.R)
	case green:
		p.picker.G = update(p.picker.G)
	case blue:
		p.picker.B = update(p.picker.B)
	case bright:
		p.picker.A = update(p.picker.A)
	}
	println(p.picker.R, p.picker.G, p.picker.B, p.picker.A)
}

// rainbowDot cycles through colors on a single pixel to indicate liveness.
func rainbowDot(now time.Time, pixel *color.RGBA) {
	const tau = 2 * math.Pi
	const period = 7 * time.Second
	pos := now.Sub(now.Truncate(period)).Seconds() / period.Seconds()
	c := color.YCbCr{
		Y:  uint8(255 / 2),
		Cb: uint8(255 * math.Sin(tau*pos)),
		Cr: uint8(255 * math.Cos(tau*pos)),
	}
	r, g, b, _ := c.RGBA()
	*pixel = graphics.Scale(color.RGBA{
		R: uint8(r >> (8)),
		G: uint8(g >> (8)),
		B: uint8(b >> (8)),
	}, graphics.MaxIntensity)
}

func (p *App) handle(e input.Event) {
	switch e {
	case input.A_Fall:
		p.currentComponent = (p.currentComponent + 1) % componentCount
	case input.B_Rise:
		p.buttons.b = true
	case input.B_Fall:
		p.buttons.b = false
	case input.C_Rise:
		p.buttons.c = true
	case input.C_Fall:
		p.buttons.c = false
	}
}
