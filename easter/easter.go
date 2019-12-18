package easter

import (
	"fmt"
	"time"

	"github.com/misterikkit/tinytimer/input"
)

const (
	None uint8 = iota
	Eggsit
	Rainbow
	Pong
	ColorPicker
	Simon
	Arcade
)

type Egger struct {
	current uint8
	history []input.Event
	last    time.Time
}

func New(ui *input.Manager) *Egger {
	e := new(Egger)
	e.current = None
	ui.AddHandler(input.A_Rise, e.handle)
	ui.AddHandler(input.A_Fall, e.handle)
	ui.AddHandler(input.B_Fall, e.handle)
	ui.AddHandler(input.C_Fall, e.handle)
	ui.AddHandler(input.BC_Fall, e.handle)
	ui.AddHandler(input.ABC_Fall, e.handle)
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
	println(fmt.Sprintf("%v", e.history))
	switch {
	case evt == input.ABC_Fall:
		e.current = Eggsit
	case matchBCBC(e.history):
		e.current = Rainbow
	case matchCycle(e.history):
		e.current = Simon
	case matchKonami(e.history):
		e.current = Pong

	case matchA10(e.history):
		e.current = ColorPicker
	case matchBounce(e.history):
		e.current = Arcade
	}

}

// The sequence is (B, C) * 5, ignoring any accidental BC events.
func matchBCBC(h []input.Event) bool {
	if len(h) < 10 {
		return false
	}
	wantB := false // start from the end looking for C
	matchLen := 0
	for i := len(h) - 1; i >= 0; i-- {
		if matchLen >= 10 {
			break
		}
		switch h[i] {
		case input.A_Fall:
			return false
		case input.B_Fall:
			if !wantB {
				return false
			}
			matchLen++
			wantB = false
		case input.C_Fall:
			if wantB {
				return false
			}
			matchLen++
			wantB = true
		}
	}
	return matchLen >= 10
}

func matchCycle(h []input.Event) bool {
	if len(h) < 6 {
		return false
	}
	matchLen := 0
	want := 'c'
	for i := len(h) - 1; i >= 0; i-- {
		if matchLen >= 6 {
			break
		}
		switch h[i] {
		case input.A_Fall:
			if want != 'a' {
				return false
			}
			matchLen++
			want = 'c'
		case input.B_Fall:
			if want != 'b' {
				return false
			}
			matchLen++
			want = 'a'
		case input.C_Fall:
			if want != 'c' {
				return false
			}
			matchLen++
			want = 'b'
		}
	}
	return matchLen >= 6
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

// The sequence is A, A, BC, BC, B, C, B, C.
// But we have to ignore the extra B and C that come with each BC.
func matchKonami(h []input.Event) bool {
	if len(h) < 14 {
		return false
	}
	h = h[len(h)-14:]
	upup := h[:4]
	if !match(upup, []input.Event{input.A_Rise, input.A_Fall, input.A_Rise, input.A_Fall}) {
		return false
	}
	lrlr := h[10:]
	if !match(lrlr, []input.Event{input.B_Fall, input.C_Fall, input.B_Fall, input.C_Fall}) {
		return false
	}
	d1 := h[4:7]
	if !dumbMatch(d1) {
		return false
	}
	d2 := h[7:10]
	if !dumbMatch(d2) {
		return false
	}
	return true
}

// Checks if input matches one "down" event.
func dumbMatch(d []input.Event) bool {
	b, c, bc := false, false, false
	for i := range d {
		switch d[i] {
		case input.B_Fall:
			b = true
		case input.C_Fall:
			c = true
		case input.BC_Fall:
			bc = true
		}
	}
	return b && c && bc
}

func matchA10(d []input.Event) bool {
	if len(d) < 20 {
		return false
	}
	d = d[len(d)-20:]
	return match(d, []input.Event{
		input.A_Rise, input.A_Fall, input.A_Rise, input.A_Fall,
		input.A_Rise, input.A_Fall, input.A_Rise, input.A_Fall,
		input.A_Rise, input.A_Fall, input.A_Rise, input.A_Fall,
		input.A_Rise, input.A_Fall, input.A_Rise, input.A_Fall,
		input.A_Rise, input.A_Fall, input.A_Rise, input.A_Fall,
	})
}

func matchBounce(d []input.Event) bool {
	if len(d) < 6 {
		return false
	}
	return match(d[len(d)-6:], []input.Event{
		input.A_Rise, input.B_Fall, input.C_Fall,
		input.B_Fall, input.C_Fall, input.A_Fall,
	})
}
