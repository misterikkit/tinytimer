// +build wasm

package main

import (
	"fmt"
	"image/color"
	"syscall/js"

	"github.com/misterikkit/tinytimer/game"
)

func setup(g *game.Game) userInterface {
	fmt.Println("hello from main.go")
	js.Global().Set("goHandleClick", js.FuncOf(handleClick(g)))
	awaitWASMLoad()
	return userInterface{}
}

func handleClick(g *game.Game) func(js.Value, []js.Value) interface{} {
	return func(_ js.Value, args []js.Value) interface{} {
		id := args[0].String()
		switch id {
		case "timer_2m":
			g.Event(game.TIMER_2M)
		case "timer_10m":
			g.Event(game.TIMER_10M)
		case "cancel":
			g.Event(game.CANCEL)
		}
		return nil
	}
}

func awaitWASMLoad() {
	loaded := make(chan struct{})
	js.Global().Set("goLoad", js.FuncOf(func(js.Value, []js.Value) interface{} { close(loaded); return nil }))
	<-loaded
}

type userInterface struct{}

// DisplayLEDs puts a Frame out into the real world.
func (userInterface) DisplayLEDs(data []color.RGBA) {
	jsonData := make([]interface{}, len(data))
	for i := range data {
		jsonData[i] = map[string]interface{}{"R": data[i].R, "G": data[i].G, "B": data[i].B}
	}
	f := js.Global().Get("DisplayLEDs")
	f.Invoke(jsonData)
}
