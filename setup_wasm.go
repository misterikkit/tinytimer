// +build wasm

package main

import (
	"image/color"
	"syscall/js"

	"github.com/misterikkit/tinytimer/apps/timer"
)

type userInterface struct {
	btnCancel bool
	btn2Min   bool
	btn10Min  bool
}

func (ui *userInterface) BtnCancel() bool { return ui.btnCancel }
func (ui *userInterface) Btn2Min() bool   { return ui.btn2Min }
func (ui *userInterface) Btn10Min() bool  { return ui.btn10Min }

func (ui *userInterface) Sleepish() { /*TODO*/ }

func setup() *userInterface {
	ui := new(userInterface)
	// js.Global().Set("goHandleClick", js.FuncOf(handleClick(g)))
	return ui
}

func handleClick(g *timer.App) func(js.Value, []js.Value) interface{} {
	return func(_ js.Value, args []js.Value) interface{} {
		id := args[0].String()
		switch id {
		case "timer_2m":
		case "timer_10m":
		case "cancel":
		}
		return nil
	}
}

// DisplayLEDs puts a Frame out into the real world.
func (userInterface) DisplayLEDs(data []color.RGBA) {
	jsonData := make([]interface{}, len(data))
	for i := range data {
		jsonData[i] = map[string]interface{}{"R": data[i].R, "G": data[i].G, "B": data[i].B}
	}
	f := js.Global().Get("DisplayLEDs")
	f.Invoke(jsonData)
}
