package fixed_test

import (
	"testing"

	"github.com/misterikkit/tinytimer/fixed"
)

func TestToFrom32(t *testing.T) {
	for i := int32(-1000); i < 1000; i++ {
		if y := fixed.FromI32(i).ToI32(); i != y {
			t.Errorf("lossy conversion: %v -> %v", i, y)
		}
	}
}

func TestToFrom8(t *testing.T) {
	for i := uint8(0); i < 0xFF; i++ {
		if y := fixed.FromU8(i).ToU8(); i != y {
			t.Errorf("lossy conversion: %v -> %v", i, y)
		}
	}
}
