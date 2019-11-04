package graphics

import (
	"image/color"

	"github.com/robaho/fixed"
)

// Geometries
var (
	// Circ is 360 degrees
	Circ = fixed.NewI(360, 0)

	// FrameSize is the number of pixels in one frame.
	FrameSize = 24

	// PixelWidth is the width in degrees of one pixel.
	PixelWidth = Circ.Div(fixed.NewI(int64(FrameSize), 0))
)

// Colors
var (
	// Scale factor for LED intensity
	toneItDown = fixed.NewI(1, 0) // / 5

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
	Position fixed.Fixed
	// Width of the sprite in degrees
	Size fixed.Fixed
}

// Render sets the sprite's pixels in the frame, rolling around the end of the
// buffer, and blending with existing color values.
func (s Sprite) Render(frame []color.RGBA) {
	half := s.Size.Div(fixed.NewI(2, 0))
	// Add Circ to position because math gets weird near zero.
	start := Circ.Add(s.Position).Sub(half) // degrees
	end := Circ.Add(s.Position).Add(half)   // degrees
	firstPx := int(start.Div(PixelWidth).Int())
	lastPx := 1 + int(end.Div(PixelWidth).Int()) // TODO: Do we need the extra 1?

	for i := firstPx; i <= lastPx; i++ {
		fi := fixed.NewI(int64(i), 0)
		// amount of overlap between sprite and current pixel in degrees
		pxStart := fi.Mul(PixelWidth)              // degrees
		pxEnd := pxStart.Add(PixelWidth)           // degrees
		amt := overlap(start, end, pxStart, pxEnd) // degrees
		coverage := amt.Div(PixelWidth)            // fraction of pixel covered
		// When rendering partial coverage, blend the color with the existing color.
		index := (len(frame) + i) % len(frame)
		frame[index] = add(
			scale(s.Color, coverage),
			scale(frame[index], fixed.NewI(1, 0).Sub(coverage)),
		)
	}
}

// overlap computes the size of overlap between ranges A and B. Returns 0.0 if
// there is no overlap.
func overlap(a1, a2, b1, b2 fixed.Fixed) fixed.Fixed {
	// Ensure that A starts to the left of B.
	if b1.LessThan(a1) {
		a1, a2, b1, b2 = b1, b2, a1, a2
	}
	// No overlap
	if b1.GreaterThan(a2) {
		return fixed.ZERO
	}
	// B is completely inside A; return size of B.
	if b2.LessThan(a2) {
		return b2.Sub(b1)
	}
	// Partial overlap
	return a2.Sub(b1)
}

// scale performs alpha multiplying on a scale of 0 to 1, rather than 0 to 255.
func scale(c color.RGBA, s fixed.Fixed) color.RGBA {
	r := s.Mul(fixed.NewI(int64(c.R), 0))
	g := s.Mul(fixed.NewI(int64(c.G), 0))
	b := s.Mul(fixed.NewI(int64(c.B), 0))
	return color.RGBA{
		R: uint8(r.Int()),
		G: uint8(g.Int()),
		B: uint8(b.Int()),
		A: c.A,
	}
}

func add(c1, c2 color.RGBA) color.RGBA {
	return color.RGBA{c1.R + c2.R, c1.G + c2.G, c1.B + c2.B, 0}
}

// Fill sets an entire frame to the given color.
func Fill(frame []color.RGBA, c color.RGBA) {
	for i := range frame {
		frame[i] = c
	}
}
