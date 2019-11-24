package easter

import (
	"time"

	"github.com/misterikkit/tinytimer/input"
)

var leftRightSpam = []input.Event{
	input.B_Fall, input.C_Fall,
	input.B_Fall, input.C_Fall,
	input.B_Fall, input.C_Fall,
	input.B_Fall, input.C_Fall,
	input.B_Fall, input.C_Fall,
}

var cycle = []input.Event{
	input.A_Fall, input.B_Fall, input.C_Fall,
	input.A_Fall, input.B_Fall, input.C_Fall,
}

const (
	None uint8 = iota
	Rainbow
	Simon
)

type Egger struct {
	current uint8
	history []input.Event
	last    time.Time
}

func New(ui *input.Manager) *Egger {
	e := new(Egger)
	e.current = None
	ui.AddHandler(input.A_Fall, e.handle)
	ui.AddHandler(input.B_Fall, e.handle)
	ui.AddHandler(input.C_Fall, e.handle)
	return e
}

// Get returns the current easter egg and clears it so as not to trigger twice.
func (e *Egger) Get() uint8 {
	ret := e.current
	e.current = None
	return ret
}
func (e *Egger) handle(evt input.Event) {
	if time.Since(e.last) > time.Second {
		e.history = e.history[:0] // this should hopefully reuse memory
	}
	e.history = append(e.history, evt)
	e.last = time.Now()
	switch {
	case match(e.history, leftRightSpam):
		e.current = Rainbow
	case match(e.history, cycle):
		e.current = Simon
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
