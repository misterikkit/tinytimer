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

	s := newSpinner()
	loader := newLoader(color.RGBA{255, 0, 0, 0})
	loader.start = time.Now()
	loader.end = time.Now().Add(2 * time.Second)

	t := time.NewTicker(time.Second / FrameRate)
	for now := range t.C {
		if loader.done {
			s.update(now)
			DisplayLEDs(s.f)
		} else {
			loader.update(now)
			DisplayLEDs(loader.f)
		}
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
