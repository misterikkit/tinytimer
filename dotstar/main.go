package main

import (
	"image/color"
	"machine"
	"time"

	"tinygo.org/x/drivers/apa102"
)

func main() {
	machine.LED.Configure(machine.PinConfig{Mode: machine.PinOutput})

	colors := []color.RGBA{
		{255, 255, 255, 255},
		{255, 0, 0, 255},
		{0, 0, 255, 255},
		{0, 255, 0, 255},
		{0, 0, 0, 0},
	}
	// machine.SPI0.Configure(machine.SPIConfig{})
	dsLED := apa102.NewSoftwareSPI(machine.PB02, machine.PB03, 0)
	for i := 0; ; i++ {
		machine.LED.Set(!machine.LED.Get())
		dsLED.WriteColors([]color.RGBA{colors[i%len(colors)]})
		time.Sleep(time.Second)
	}
}
