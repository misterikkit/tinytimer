package main

import (
	"fmt"
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
	js.Global().Set("goHandleClick", js.FuncOf(handleClick))
	awaitWASMLoad()

	// s := newSpinner()
	// loader := newLoader(color.RGBA{255, 0, 0, 0})
	// loader.start = time.Now()
	// loader.end = time.Now().Add(20 * time.Second)

	g := NewGame()
	t := time.NewTicker(time.Second / FrameRate)
	for now := range t.C {
		g.update(now)
		// if loader.done {
		// 	s.update(now)
		// 	DisplayLEDs(s.f)
		// } else {
		// 	loader.update(now)
		// 	DisplayLEDs(loader.f)
		// }
	}
}

func handleClick(this js.Value, args []js.Value) interface{} {
	id := args[0].String()
	fmt.Println(id)
	return nil
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
		jsonData[i] = map[string]interface{}{"R": data[i].R, "G": data[i].G, "B": data[i].B}
	}
	f := js.Global().Get("DisplayLEDs")
	f.Invoke(jsonData)
}
