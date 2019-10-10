package main

import (
	"fmt"
	"image/color"
	"math"
	"syscall/js"
	"time"
)

const (
	// Tau is better than Pi.
	Tau = 2 * math.Pi

	// Frames per second
	FrameRate = 60

	// Number of pixels in a frame
	FrameSize = 24

	// Size of a pixel in radians
	PixelWidth = Tau / FrameSize
)

// RGB is a color struct that plays nicely with wasm. Will need to move to
// `color.RGBA` on TinyGo.
type RGB struct {
	R, G, B uint8
}

// JSValue converts the given color to a JS object.
func (c RGB) JSValue() js.Value {
	return js.ValueOf(map[string]interface{}{
		"R": c.R, "G": c.G, "B": c.B,
	})
}

// Frame represents all the pixels in one frame of animation.
// Should this be [24]RGB?
type Frame = []RGB

// newFrame returns an all-black frame.
func newFrame() Frame { return make(Frame, FrameSize) }

// reset blanks a frame (to avoid allocating new frames)
func reset(f Frame) {
	for i := range f {
		f[i] = RGB{}
	}
}

func main() {
	fmt.Println("hello from main.go")
	loaded := make(chan struct{})
	js.Global().Set("goLoad", js.FuncOf(func(js.Value, []js.Value) interface{} { close(loaded); return nil }))
	<-loaded
	DisplayLEDs(newFrame())

	// Set up animation
	f := newFrame()
	var dots []sprite
	for i := 0; i < 7; i++ {
		dots = append(dots, sprite{Size: 1.5 * PixelWidth, Color: color.RGBA{0x32, 0x6C, 0xE5, 0}})
	}
	const (
		period = 10 * time.Second
		divide = Tau / 7.0
	)

	// Update animation
	update := func(now time.Time) {
		reset(f)

		// compute fraction through the period
		progress := float32(now.Sub(now.Truncate(period))) / float32(period)
		for i := range dots {
			dots[i].Position = Tau*progress + float32(i)*divide
			dots[i].Render(f)
		}

		DisplayLEDs(f)
	}

	t := time.NewTicker(time.Second / FrameRate)
	for now := range t.C {
		update(now)
	}

	// <-context.Background().Done()
}

// DisplayLEDs puts a Frame out into the real world.
func DisplayLEDs(data Frame) {
	jsonData := make([]interface{}, len(data))
	for i := range data {
		jsonData[i] = data[i]
	}
	// copy(jsonData, data)
	f := js.Global().Get("DisplayLEDs")
	f.Invoke(jsonData)
}

// sprite is a one-dimensional object to be rendered onto a Frame.
type sprite struct {
	Color color.RGBA
	// Center of the sprite in radians
	Position float32
	// Width of the sprite in radians
	Size float32
}

// Render will overwrite the requisite pixels.
// TODO: implement alpha channels?
func (s sprite) Render(f Frame) {

	start, end := s.Position-s.Size/2, s.Position+s.Size/2
	firstPx := int(start / PixelWidth)
	lastPx := int(0.5 + end/PixelWidth)
	// fmt.Printf("start=%0.2f end=%0.2f firstPx=%v lastPx=%v\n", start, end, firstPx, lastPx)
	// if lastPx >= FrameSize {
	// 	fmt.Println("wrapping around")
	// }
	for i := firstPx; i <= lastPx; i++ {
		fi := float32(i)
		// overlap amount in radians
		amt := overlap(start, end, fi*PixelWidth, (fi+1.0)*PixelWidth)
		coverage := amt / PixelWidth // fraction of pixel covered
		f[(len(f)+i)%len(f)] = RGB{
			R: uint8(float32(s.Color.R) * coverage),
			G: uint8(float32(s.Color.G) * coverage),
			B: uint8(float32(s.Color.B) * coverage),
		}
	}
}

// overlap computes the size of overlap between ranges A and B. Returns 0.0 if
// there is no overlap.
func overlap(a1, a2, b1, b2 float32) float32 {
	// Ensure that A starts to the left of B.
	if b1 < a1 {
		a1, a2, b1, b2 = b1, b2, a1, a2
	}
	// No overlap
	if b1 > a2 {
		return 0.0
	}
	// B is completely inside A; return size of B.
	if b2 < a2 {
		return b2 - b1
	}
	// Partial overlap
	return a2 - b1
}
