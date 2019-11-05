package graphics

import (
	"image/color"
	"testing"

	"golang.org/x/image/math/fixed"
)

func fix(f float64) fixed.Int26_6 {
	return fixed.Int26_6(int32(f * 64.0))
}

func TestRender(t *testing.T) {
	// Outputs were manually inspected
	frame := make([]color.RGBA, FrameSize)
	sprite := Sprite{Color: color.RGBA{255, 0, 0, 0}, Position: 0, Size: PixelWidth}
	sprite.Render(frame)
	t.Logf("%v", frame)
	sprite = Sprite{Color: color.RGBA{0, 255, 0, 0}, Position: PixelWidth / 2, Size: PixelWidth}
	sprite.Render(frame)
	t.Logf("%v", frame)
}
