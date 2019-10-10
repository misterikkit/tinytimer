package main

import (
	"fmt"
	"image/color"
	"math"
	"syscall/js"
	"time"
)

// Constants used in rendering
const (
	// Tau is better than Pi.
	Tau = 2 * math.Pi

	FrameRate = 60 // per second

	FrameSize = 24 // pixels

	PixelWidth = Tau / FrameSize // radians
)

func main() {
	fmt.Println("hello from main.go")
	awaitWASMLoad()

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
		f.reset()

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

}

func awaitWASMLoad() {
	loaded := make(chan struct{})
	js.Global().Set("goLoad", js.FuncOf(func(js.Value, []js.Value) interface{} { close(loaded); return nil }))
	<-loaded
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
