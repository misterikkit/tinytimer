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

	tickInterval := scaleDuration(time.Second)
	nextTick := time.Now().Add(tickInterval)
	// return the time until next timer tick, and update `nextTick`
	tick := func() time.Duration {
		// TODO: skip a tick if needed
		left := nextTick.Sub(time.Now())
		nextTick = nextTick.Add(tickInterval)
		return left
	}

	for {
		ws.WriteColors(AllRed)
		time.Sleep(tick())

		ws.WriteColors(AllBlue)
		time.Sleep(tick())

		ws.WriteColors(Fun)
		time.Sleep(tick())
	}
}

const TimeScale = 1.75

func scaleDuration(d time.Duration) time.Duration {
	fd := float32(d)
	fd /= TimeScale
	return time.Duration(fd)
}
