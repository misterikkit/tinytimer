package graphics

import (
	"fmt"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

// rgb is a shorthand for constructing RGBA literals.
func rgb(r, g, b uint8) color.RGBA { return color.RGBA{r, g, b, 0} }

// newFrame returns a full FrameSize-sized newFrame populated with the given pixel values.
func newFrame(cs ...color.RGBA) frame {
	f := make([]color.RGBA, FrameSize)
	copy(f, cs)
	return f
}

// frame overrides GoString on a color slice for better test output.
type frame []color.RGBA

func (f frame) GoString() string {
	pxs := make([]string, len(f))
	for i := range f {
		pxs[i] = fmt.Sprintf("#%02x%02x%02x", f[i].R, f[i].G, f[i].B)
	}
	return fmt.Sprint(pxs)
}

func TestSprite_Render(t *testing.T) {
	tests := []struct {
		name     string
		sprites  []Sprite
		expected frame
	}{
		{
			"default black sprites",
			[]Sprite{{}, {}},
			newFrame(),
		},
		{
			"single pixel",
			[]Sprite{{Color: Red, Size: PixelWidth, Position: PixelWidth / 2}},
			newFrame(Red),
		},
		{
			"half pixel",
			[]Sprite{{Color: Red, Size: PixelWidth / 2, Position: PixelWidth / 2}},
			newFrame(rgb(127, 0, 0)),
		},
		{
			"misaligned pixel",
			[]Sprite{{Color: Red, Size: PixelWidth, Position: 0.5 * PixelWidth}},
			newFrame(rgb(127, 0, 0), rgb(127, 0, 0)),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := newFrame()
			for _, s := range test.sprites {
				s.Render(f)
			}
			assert.Equal(t, test.expected, f)
		})
	}
}

func TestOverlap(t *testing.T) {
	tests := []struct {
		a1, a2, b1, b2 float32
		expected       float32
	}{
		{0, 0, 0, 0, 0},
		{0, 1, 1, 2, 0},
		{0, 1, 0.5, 5, 0.5},
		{0, 4, 1, 2, 1},
		{0.3, 4.3, 1.3, 2.3, 1},
	}
	for _, test := range tests {
		actual := overlap(test.a1, test.a2, test.b1, test.b2)
		if actual != test.expected {
			t.Errorf("overlap(%v, %v, %v, %v) = %v, expected %v",
				test.a1, test.a2, test.b1, test.b2, actual, test.expected)
		}
		// symmetric test
		actual = overlap(test.b1, test.b2, test.a1, test.a2)
		if actual != test.expected {
			t.Errorf("overlap(%v, %v, %v, %v) = %v, expected %v",
				test.b1, test.b2, test.a1, test.a2, actual, test.expected)
		}
	}
}
