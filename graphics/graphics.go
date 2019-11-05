package graphics

import (
	"image/color"
)

// Geometries
var (
	// Circ is 360 degrees
	Circ = float32(360.0)

	// FrameSize is the number of pixels in one frame.
	FrameSize = 24

	// PixelWidth is the width in degrees of one pixel.
	PixelWidth = Circ / float32(FrameSize)
)

// Colors
var (
	// Scale factor for LED intensity
	toneItDown = float32(0.3)

	Black = color.RGBA{}
	White = Scale(color.RGBA{0xFF, 0xFF, 0xFF, 0}, toneItDown)
	Red   = Scale(color.RGBA{0xFF, 0, 0, 0}, toneItDown)

	// K8SBlue is used in the k8s logo
	K8SBlue   = Scale(color.RGBA{0x32, 0x6C, 0xE5, 0}, toneItDown) // H=221deg S=78.2% V=89.8%
	CSIOrange = Scale(color.RGBA{0xF5, 0x91, 0x1E, 0}, toneItDown) // H=32deg  S=87.8% V=96.1%
)

// Sprite is a one-dimensional object to be rendered onto a Frame.
type Sprite struct {
	Color color.RGBA
	// Center of the sprite in degrees
	Position float32
	// Width of the sprite in degrees
	Size float32
}

// Render sets the sprite's pixels in the frame, rolling around the end of the
// buffer, and blending with existing color values.
func (s Sprite) Render(frame []color.RGBA) {
	half := s.Size / 2
	// Add Circ to position because math gets weird near zero.
	start := Circ + (s.Position) - (half) // degrees
	end := Circ + (s.Position) + (half)   // degrees
	firstPx := int(start / (PixelWidth))
	lastPx := 1 + int(end/(PixelWidth)) // TODO: Do we need the extra 1?

	for i := firstPx; i <= lastPx; i++ {
		fi := float32(i)
		// amount of overlap between sprite and current pixel in degrees
		pxStart := fi * (PixelWidth)               // degrees
		pxEnd := pxStart + (PixelWidth)            // degrees
		amt := overlap(start, end, pxStart, pxEnd) // degrees
		coverage := amt / (PixelWidth)             // fraction of pixel covered
		// When rendering partial coverage, blend the color with the existing color.
		index := (len(frame) + i) % len(frame)
		frame[index] = Add(
			Scale(s.Color, coverage),
			Scale(frame[index], 1.0-(coverage)),
		)
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
		return 0
	}
	// B is completely inside A; return size of B.
	if b2 < a2 {
		return b2 - (b1)
	}
	// Partial overlap
	return a2 - (b1)
}

// Scale performs alpha multiplying on a Scale of 0 to 1, rather than 0 to 255.
func Scale(c color.RGBA, s float32) color.RGBA {
	r := s * (float32(c.R))
	g := s * (float32(c.G))
	b := s * (float32(c.B))
	return color.RGBA{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: c.A,
	}
}

func Add(c1, c2 color.RGBA) color.RGBA {
	return color.RGBA{c1.R + c2.R, c1.G + c2.G, c1.B + c2.B, 0}
}

// Fill sets an entire frame to the given color.
func Fill(frame []color.RGBA, c color.RGBA) {
	for i := range frame {
		frame[i] = c
	}
}
