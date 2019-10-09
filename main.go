package main

import (
	"fmt"
	"syscall/js"
	"time"
)

type RGB struct {
	R, G, B uint8
}

func (c RGB) JSValue() js.Value {
	return js.ValueOf(map[string]interface{}{
		"R": c.R, "G": c.G, "B": c.B,
	})
}

type Frame []RGB

func (f Frame) Copy() Frame {
	f2 := make(Frame, len(f))
	for i := range f {
		f2[i] = f[i]
	}
	// copy(f2, f)
	return f2
}

var blank = Frame{
	{0, 0, 0},
	{0, 0, 0},
	{0, 0, 0},
	{0, 0, 0},
	{0, 0, 0},
	{0, 0, 0},
	{0, 0, 0},
	{0, 0, 0},
	{0, 0, 0},
	{0, 0, 0},
	{0, 0, 0},
	{0, 0, 0},
	{0, 0, 0},
	{0, 0, 0},
	{0, 0, 0},
	{0, 0, 0},
	{0, 0, 0},
	{0, 0, 0},
	{0, 0, 0},
	{0, 0, 0},
	{0, 0, 0},
	{0, 0, 0},
	{0, 0, 0},
	{0, 0, 0},
}

func main() {
	fmt.Println("hello from main.go")
	loaded := make(chan struct{})
	js.Global().Set("goLoad", js.FuncOf(func(js.Value, []js.Value) interface{} { close(loaded); return nil }))
	<-loaded
	DisplayLEDs(blank)

	// fadeIn := Tween{from: time.Now(), to: time.Now().Add(4 * time.Second), start: 30, end: -10}
	updaters := []UpdateFunc{
		// func(now time.Time) bool {
		// 	v, ok := fadeIn.Value(now)
		// 	fmt.Println(v)
		// 	return ok
		// },
		func(now time.Time) bool {
			// const period = 3 * time.Second
			// phase := now.Sub(now.Round(period))
			// frac := float64(phase) / float64(period)
			// fmt.Println(frac)
			step := int(time.Second.Nanoseconds()) / 24 // const
			pixel := float32(now.Nanosecond()) / float32(step)
			i := int(pixel)
			weight := 1.0 - (pixel - float32(i))
			val := uint8(255 * weight)
			f := blank.Copy()
			f[i] = RGB{val, val, val}
			f[(i+1)%len(f)] = RGB{255 - val, 255 - val, 255 - val}
			DisplayLEDs(f)
			return true
		},
	}
	t := time.NewTicker(time.Second / 60)
	for now := range t.C {
		updaters = RunUpdaters(updaters, now)
	}

	// <-context.Background().Done()
}

func DisplayLEDs(data Frame) {
	jsonData := make([]interface{}, len(data))
	for i := range data {
		jsonData[i] = data[i]
	}
	// copy(jsonData, data)
	f := js.Global().Get("DisplayLEDs")
	f.Invoke(jsonData)
}

// An UpdateFunc is any func meant to be called once per timer tick. An
// UpdateFunc that returns true will be called again on the next tick, while one
// that returns false will be discarded.
type UpdateFunc = func(time.Time) bool

// RunUpdaters runs each given UpdateFunc, and returns a slice of updaters that
// should be called at the next tick.
func RunUpdaters(fs []UpdateFunc, now time.Time) []UpdateFunc {
	ret := make([]UpdateFunc, 0, len(fs))
	for _, f := range fs {
		if !f(now) {
			continue
		}
		ret = append(ret, f)
	}
	return ret
}

// Tween computes a linear transition between two values over time.
type Tween struct {
	from, to   time.Time
	start, end int32
}

// Value gets the current value of a Tween. It returns true if the Tween is not
// yet finished.
func (t Tween) Value(now time.Time) (int32, bool) {
	if now.After(t.to) {
		return t.end, false
	}
	frac := float32(now.Sub(t.from)) / float32(t.to.Sub(t.from))
	return int32(float32(t.end-t.start)*frac) + t.start, true
}
