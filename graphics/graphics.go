package graphics

import (
	"image/color"

	"golang.org/x/image/math/fixed"
)

// Geometries
const (
	// Circ is 360 degrees
	Circ = fixed.Int26_6(360 << 6)

	// FrameSize is the number of pixels in one frame.
	FrameSize = 24

	// PixelWidth is the width in degrees of one pixel.
	PixelWidth = Circ / FrameSize
)

// Scale factor for LED intensity
const toneItDown = fixed.Int26_6(1 << 6) // / 5

// Colors
var (
	Black = color.RGBA{}
	White = scale(color.RGBA{0xFF, 0xFF, 0xFF, 0}, toneItDown)
	Red   = scale(color.RGBA{0xFF, 0, 0, 0}, toneItDown)

	// K8SBlue is used in the k8s logo
	K8SBlue   = scale(color.RGBA{0x32, 0x6C, 0xE5, 0}, toneItDown) // H=221deg S=78.2% V=89.8%
	CSIOrange = scale(color.RGBA{0xF5, 0x91, 0x1E, 0}, toneItDown) // H=32deg  S=87.8% V=96.1%
)

// Sprite is a one-dimensional object to be rendered onto a Frame.
type Sprite struct {
	Color color.RGBA
	// Center of the sprite in degrees
	Position fixed.Int26_6
	// Width of the sprite in degrees
	Size fixed.Int26_6
}

// Render sets the sprite's pixels in the frame, rolling around the end of the
// buffer, and blending with existing color values.
func (s Sprite) Render(frame []color.RGBA) {
	// Add Circ to position because math gets weird near zero.
	start, end := Circ+s.Position-s.Size/2, Circ+s.Position+s.Size/2 // degrees
	firstPx := int(start / PixelWidth)
	lastPx := 1 + int(end/PixelWidth) // TODO: Do we need the extra 1?

	for i := firstPx; i <= lastPx; i++ {
		// amount of overlap between sprite and current pixel in degrees
		pxStart := fixed.I(i).Mul(PixelWidth)      // degrees
		pxEnd := fixed.I(i + 1).Mul(PixelWidth)    // degrees
		amt := overlap(start, end, pxStart, pxEnd) // degrees
		coverage := amt * 64 / PixelWidth          // fraction of pixel covered
		// When rendering partial coverage, blend the color with the existing color.
		index := (len(frame) + i) % len(frame)
		frame[index] = add(scale(s.Color, coverage), scale(frame[index], fixed.I(1)-coverage))
	}
}

// overlap computes the size of overlap between ranges A and B. Returns 0.0 if
// there is no overlap.
func overlap(a1, a2, b1, b2 fixed.Int26_6) fixed.Int26_6 {
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
		return b2 - b1
	}
	// Partial overlap
	return a2 - b1
}

// scale performs alpha multiplying on a scale of 0 to 1, rather than 0 to 255.
func scale(c color.RGBA, s fixed.Int26_6) color.RGBA {
	r := s.Mul(fixed.I(int(c.R)))
	g := s.Mul(fixed.I(int(c.G)))
	b := s.Mul(fixed.I(int(c.B)))
	return color.RGBA{
		R: uint8(r.Round()),
		G: uint8(g.Round()),
		B: uint8(b.Round()),
		A: c.A,
	}
}

func add(c1, c2 color.RGBA) color.RGBA {
	return color.RGBA{c1.R + c2.R, c1.G + c2.G, c1.B + c2.B, 0}
}
