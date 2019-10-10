package main

import (
	"image/color"
	"syscall/js"
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
type Frame []RGB

// newFrame returns an all-black frame.
func newFrame() Frame { return make(Frame, FrameSize) }

// reset blanks a frame (to avoid allocating new frames)
func (f Frame) reset() {
	for i := range f {
		f[i] = RGB{}
	}
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
