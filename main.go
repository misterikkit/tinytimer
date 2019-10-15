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
	g := NewGame()
	js.Global().Set("goHandleClick", js.FuncOf(handleClick(&g)))
	awaitWASMLoad()

	t := time.NewTicker(time.Second / FrameRate)
	for now := range t.C {
		g.update(now)
		DisplayLEDs(*g.animation.f)
	}
}

func handleClick(g *game) func(js.Value, []js.Value) interface{} {
	return func(_ js.Value, args []js.Value) interface{} {
		id := args[0].String()
		switch id {
		case "timer_2m":
			g.event(TIMER_2M)
		case "timer_10m":
			g.event(TIMER_10M)
		case "cancel":
			g.event(CANCEL)
		}
		return nil
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
		jsonData[i] = map[string]interface{}{"R": data[i].R, "G": data[i].G, "B": data[i].B}
	}
	f := js.Global().Get("DisplayLEDs")
	f.Invoke(jsonData)
}
