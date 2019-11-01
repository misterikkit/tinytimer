package main

import (
	"image/color"
	"machine"
	"time"

	"github.com/misterikkit/tinytimer/ws2812"
)

// colors
var (
	Red  = color.RGBA{1, 0, 0, 1}
	Blue = color.RGBA{0, 0, 1, 1}

	K8SBlue   = color.RGBA{0x32, 0x6C, 0xE5, 0}
	CSIOrange = color.RGBA{0xF5, 0x91, 0x1E, 0}

	AllRed  = []color.RGBA{CSIOrange, CSIOrange, CSIOrange, CSIOrange, CSIOrange, CSIOrange, CSIOrange, CSIOrange, CSIOrange, CSIOrange, CSIOrange, CSIOrange, CSIOrange, CSIOrange, CSIOrange, CSIOrange, CSIOrange, CSIOrange, CSIOrange, CSIOrange, CSIOrange, CSIOrange, CSIOrange, CSIOrange}
	AllBlue = []color.RGBA{K8SBlue, K8SBlue, K8SBlue, K8SBlue, K8SBlue, K8SBlue, K8SBlue, K8SBlue, K8SBlue, K8SBlue, K8SBlue, K8SBlue, K8SBlue, K8SBlue, K8SBlue, K8SBlue, K8SBlue, K8SBlue, K8SBlue, K8SBlue, K8SBlue, K8SBlue, K8SBlue, K8SBlue}
	Fun     = []color.RGBA{
		{0, 10, 10, 10},
		{10, 0, 10, 10},
		{10, 10, 0, 10},
		{10, 0, 0, 10},
		{0, 10, 0, 10},
		{0, 0, 10, 10},
		{10, 10, 10, 10},
		{10, 10, 10, 10},
		{10, 10, 10, 10},
		{10, 10, 10, 10},
		{10, 10, 10, 10},
		{10, 10, 10, 10},
		{10, 10, 10, 10},
		{10, 10, 10, 10},
		{10, 10, 10, 10},
		{10, 10, 10, 10},
		{10, 10, 10, 10},
		{10, 10, 10, 10},
		{10, 10, 10, 10},
		{10, 10, 10, 10},
		{10, 10, 10, 10},
		{10, 10, 10, 10},
		{10, 10, 10, 10},
		{10, 10, 10, 10},
	}
)

func main() {
	neo := machine.D5 // special level-shifted output pin
	neo.Configure(machine.PinConfig{Mode: machine.PinOutput})
	ws := ws2812.New(neo)

	for {
		ws.WriteColors(AllRed)
		time.Sleep(time.Second)

		ws.WriteColors(AllBlue)
		time.Sleep(time.Second)

		ws.WriteColors(Fun)
		time.Sleep(time.Second)
	}
}
