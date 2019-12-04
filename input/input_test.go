package input_test

import (
	"testing"

	"github.com/misterikkit/tinytimer/input"
	"github.com/stretchr/testify/assert"
)

func TestInput(t *testing.T) {

	var a, b, c bool
	in := input.NewManager(func() bool { return a }, func() bool { return b }, func() bool { return c })
	var events []input.Event
	h := func(e input.Event) { events = append(events, e) }
	in.AddHandler(input.A_Rise, h)
	in.AddHandler(input.A_Fall, h)
	in.AddHandler(input.B_Rise, h)
	in.AddHandler(input.B_Fall, h)
	in.AddHandler(input.C_Rise, h)
	in.AddHandler(input.C_Fall, h)
	in.AddHandler(input.AB_Rise, h)
	in.AddHandler(input.AB_Fall, h)
	in.AddHandler(input.AC_Rise, h)
	in.AddHandler(input.AC_Fall, h)
	in.AddHandler(input.BC_Rise, h)
	in.AddHandler(input.BC_Fall, h)
	in.AddHandler(input.ABC_Rise, h)
	in.AddHandler(input.ABC_Fall, h)

	sequence := []struct{ a, b, c bool }{
		{false, false, false},
		// part 1
		{false, true, false},
		{true, true, false},
		{false, true, false},
		{false, false, false},
		// part 2 - all combos
		{true, true, true},
		{false, false, false},
	}

	for _, state := range sequence {
		a, b, c = state.a, state.b, state.c
		in.Poll()
	}
	expected := []input.Event{
		// part 1
		input.B_Rise,
		input.A_Rise,
		input.AB_Rise,
		input.A_Fall,
		input.B_Fall,
		input.AB_Fall,
		// part 2 - combos fall after their components
		input.A_Rise,
		input.B_Rise,
		input.C_Rise,
		input.AB_Rise,
		input.AC_Rise,
		input.BC_Rise,
		input.ABC_Rise,

		input.A_Fall,
		input.B_Fall,
		input.C_Fall,
		input.AB_Fall,
		input.AC_Fall,
		input.BC_Fall,
		input.ABC_Fall,
	}
	assert.Equal(t, expected, events)
}
