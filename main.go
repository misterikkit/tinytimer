package main

import (
	"context"
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

func main() {
	fmt.Println("hello from main.go")
	// t := time.NewTicker(time.Second / 24)
	// fadeIn := Tween{from: time.Now(), to: time.Now().Add(4 * time.Second), start: 30, end: -10}
	// updaters := []UpdateFunc{
	// 	func(now time.Time) bool {
	// 		v, ok := fadeIn.Value(now)
	// 		fmt.Println(v)
	// 		return ok
	// 	},
	// }
	// for now := range t.C {
	// 	updaters = RunUpdaters(updaters, now)
	// }
	c := RGB{}
	c.JSValue()
	time.Sleep(time.Second)
	DisplayLEDs([]interface{}{
		RGB{0, 0, 0},
		RGB{0, 0, 0},
		RGB{0, 0, 0},
		RGB{0, 0, 0},
		RGB{0, 0, 0},
		RGB{0, 0, 0},
		RGB{0, 0, 0},
		RGB{0, 0, 0},
		RGB{0, 0, 0},
		RGB{0, 0, 0},
		RGB{0, 0, 0},
		RGB{0, 0, 0},
		RGB{0, 0, 0},
		RGB{0, 0, 0},
		RGB{0, 0, 0},
		RGB{0, 0, 0},
		RGB{0, 0, 0},
		RGB{0, 0, 0},
		RGB{0, 0, 0},
		RGB{0, 0, 0},
		RGB{0, 0, 0},
		RGB{0, 0, 0},
		RGB{0, 0, 0},
		RGB{0, 0, 0},
	})
	<-context.Background().Done()
}

func DisplayLEDs(data []interface{}) {
	// jsonData := make([]interface{}, 0, len(data))
	// for _, d := range data {
	// 	jsonData = append(jsonData, js.ValueOf(map[string]uint8{"R": d.R, "G": d.G, "B": d.B}))
	// }
	f := js.Global().Get("DisplayLEDs")
	// fmt.Println(f.Type())
	f.Invoke(js.ValueOf(data))
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
