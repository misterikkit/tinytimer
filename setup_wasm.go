// +build wasm

package main

import (
	"image/color"
	"syscall/js"
)

type userInterface struct {
	pollCancel js.Value
	poll2Min   js.Value
	poll10Min  js.Value
	display    js.Value
}

func setup() *userInterface {
	ui := &userInterface{
		pollCancel: js.Global().Get("ButtonCancel"),
		poll2Min:   js.Global().Get("Button2M"),
		poll10Min:  js.Global().Get("Button10M"),
		display:    js.Global().Get("DisplayLEDs"),
	}
	return ui
}

func (ui *userInterface) BtnCancel() bool { return ui.pollCancel.Invoke().String() == "true" }
func (ui *userInterface) Btn2Min() bool   { return ui.poll2Min.Invoke().String() == "true" }
func (ui *userInterface) Btn10Min() bool  { return ui.poll10Min.Invoke().String() == "true" }

func (ui *userInterface) Sleepish() { /*TODO*/ }

// DisplayLEDs puts a Frame out into the real world.
func (ui *userInterface) DisplayLEDs(data []color.RGBA) {
	jsonData := make([]interface{}, len(data))
	for i := range data {
		jsonData[i] = map[string]interface{}{"R": data[i].R, "G": data[i].G, "B": data[i].B}
	}
	ui.display.Invoke(jsonData)
}
