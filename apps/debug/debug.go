package debug

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
	d := &App{
		frame: make([]color.RGBA, graphics.FrameSize),
	}
	ui.AddHandler(input.A_Fall, d.handle)
	ui.AddHandler(input.B_Fall, d.handle)
	ui.AddHandler(input.C_Fall, d.handle)
	ui.AddHandler(input.A_Rise, d.handle)
	ui.AddHandler(input.B_Rise, d.handle)
	ui.AddHandler(input.C_Rise, d.handle)
	return d
}
func (d *App) Frame() []color.RGBA { return d.frame }

func (d *App) Update(now time.Time) {
	graphics.Fill(d.frame, graphics.Black)
	rainbowDot(now, &d.frame[graphics.FrameSize-1])
	rainbowDot(now, &d.frame[8])

	// Count number of ticks since
	if now.Sub(d.lastTick) > tick {
		d.updateColorPicker()
		d.lastTick = now
	}

	//////////////
	// Render data
	r, g, b, a := d.picker.RGBA()
	actualColor := color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
	s := graphics.Sprite{
		Color:    actualColor,
		Position: graphics.Circ * 3 / 4,
		Size:     4 * graphics.PixelWidth,
	}
	s.Render(d.frame)
	switch d.currentComponent {
	case red:
		renderByte(d.frame, actualColor.R, graphics.Red)
	case green:
		renderByte(d.frame, actualColor.G, graphics.Green)
	case blue:
		renderByte(d.frame, actualColor.B, graphics.Blue)
	case bright:
		renderByte(d.frame, actualColor.A, graphics.White)
	}
}

func renderByte(frame []color.RGBA, value uint8, color color.RGBA) {
	for i := uint8(0); i < 8; i++ {
		if value&(1<<i) != 0 {
			frame[i] = color
		}
	}
}

func (d *App) updateColorPicker() {
	update := func(v uint8) uint8 { return v }
	if d.buttons.b && !d.buttons.c {
		update = func(v uint8) uint8 { return v - 1 }
	}
	if !d.buttons.b && d.buttons.c {
		update = func(v uint8) uint8 { return v + 1 }
	}
	switch d.currentComponent {
	case red:
		d.picker.R = update(d.picker.R)
	case green:
		d.picker.G = update(d.picker.G)
	case blue:
		d.picker.B = update(d.picker.B)
	case bright:
		d.picker.A = update(d.picker.A)
	}
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

func (d *App) handle(e input.Event) {
	switch e {
	case input.A_Fall:
		d.currentComponent = (d.currentComponent + 1) % componentCount
	case input.B_Rise:
		d.buttons.b = true
	case input.B_Fall:
		d.buttons.b = false
	case input.C_Rise:
		d.buttons.c = true
	case input.C_Fall:
		d.buttons.c = false
	}
}
