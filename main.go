package main

import (
	"fmt"
	"time"
)

func main() {
	t := time.NewTicker(time.Second / 24)
	fadeIn := Tween{from: time.Now(), to: time.Now().Add(4 * time.Second), start: 30, end: -10}
	updaters := []UpdateFunc{
		func(now time.Time) bool {
			v, ok := fadeIn.Value(now)
			fmt.Println(v)
			return ok
		},
	}
	for now := range t.C {
		updaters = RunUpdaters(updaters, now)
	}
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
