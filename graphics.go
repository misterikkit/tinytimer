package main

import (
	"image/color"
)

// Colors
var (
	Black = color.RGBA{}
	White = scale(color.RGBA{0xFF, 0xFF, 0xFF, 0}, 0.2)
	Red   = color.RGBA{0xFF, 0, 0, 0}

	// K8SBlue is used in the k8s logo
	K8SBlue   = scale(color.RGBA{0x32, 0x6C, 0xE5, 0}, 0.2) // H=221deg S=78.2% V=89.8%
	CSIOrange = scale(color.RGBA{0xF5, 0x91, 0x1E, 0}, 0.2) // H=32deg  S=87.8% V=96.1%
)

// Frame represents all the pixels in one frame of animation.
// Should this be [24]RGB?
type Frame []color.RGBA

// newFrame returns an all-black frame.
func newFrame() Frame { return make(Frame, FrameSize) }

// fill blanks a frame (to avoid allocating new frames)
func (f Frame) fill(c color.RGBA) {
	for i := range f {
		f[i] = c
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
func (s sprite) Render(f Frame) {
	// fmt.Printf("RENDER: PixelWidth=%0.3f\n", PixelWidth)
	// fmt.Printf("RENDER: size=%0.3f pos=%0.3f\n", s.Size, s.Position)

	// Add Tau to position because math gets weird near zero.
	start, end := Tau+s.Position-s.Size/2, Tau+s.Position+s.Size/2
	// fmt.Printf("RENDER: start=%0.3f, end=%0.3f\n", start, end)
	firstPx := int(start / PixelWidth)
	lastPx := int(0.5 + end/PixelWidth)

	for i := firstPx; i <= lastPx; i++ {
		fi := float32(i)
		// amount of overlap between sprite and current pixel in radians
		// fmt.Printf("RENDER: pixel[%d] is {%0.3f, %0.3f}\n", i, fi*PixelWidth, (fi+1.0)*PixelWidth)
		amt := overlap(start, end, fi*PixelWidth, (fi+1.0)*PixelWidth)
		coverage := amt / PixelWidth // fraction of pixel covered
		// When rendering partial coverage, blend the color with the existing color.
		index := (len(f) + i) % len(f)
		// fmt.Printf("RENDER: coverage is %v\n", coverage)
		f[index] = add(scale(s.Color, coverage), scale(f[index], 1.0-coverage))
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

// scale performs alpha multiplying without a conversion through 32-bit values.
func scale(c color.RGBA, s float32) color.RGBA {
	scaler := uint16(s * 0xFF)
	r := uint16(c.R)
	r *= scaler
	r /= 0xFF
	g := uint16(c.G)
	g *= scaler
	g /= 0xFF
	b := uint16(c.B)
	b *= scaler
	b /= 0xFF
	return color.RGBA{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: c.A,
	}
}

func add(c1, c2 color.RGBA) color.RGBA {
	return color.RGBA{c1.R + c2.R, c1.G + c2.G, c1.B + c2.B, 0}
}
