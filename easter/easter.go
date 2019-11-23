package easter

import (
	"image/color"
	"time"

	"github.com/misterikkit/tinytimer/input"
	"github.com/misterikkit/tinytimer/rainbow"
)

var leftRightSpam = []input.Event{
	input.B_Fall, input.C_Fall,
	input.B_Fall, input.C_Fall,
	input.B_Fall, input.C_Fall,
	input.B_Fall, input.C_Fall,
	input.B_Fall, input.C_Fall,
}

// TODO: this needs a better home
type App interface {
	Update(time.Time)
	Frame() []color.RGBA
}

type Egger struct {
	rec     func(App)
	history []input.Event
	last    time.Time
}

func New(ui *input.Manager, receive func(App)) *Egger {
	e := new(Egger)
	ui.AddHandler(input.B_Fall, e.handle)
	ui.AddHandler(input.C_Fall, e.handle)
	return e
}

func (e *Egger) handle(evt input.Event) {
	if time.Since(e.last) > time.Second {
		e.history = e.history[:0] // this should hopefully reuse memory
	}
	e.history = append(e.history, evt)
	e.last = time.Now()
	if match(e.history, leftRightSpam) {
		e.rec(new(rainbow.Egg))
	}
}

func match(a, b []input.Event) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
