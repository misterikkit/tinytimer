package timer

import (
	"image/color"
	"time"

	"github.com/misterikkit/tinytimer/animation"
	"github.com/misterikkit/tinytimer/graphics"
	"github.com/misterikkit/tinytimer/hack"
	"github.com/misterikkit/tinytimer/input"
)

type state uint8
type event uint8

// valid states
const (
	boot state = iota
	idle
	countdown
	timerPop
)

// events
const (
	animationDone event = iota
	timer2M
	timer10M
	cancel
)

// App keeps track of the current state of the timer.
type App struct {
	state     state
	Animation animation.Interface
	// TODO: Make animations reusable and store them here to avoid allocations.
}

// New creates a new game in the "boot animation" state.
func New(ui *input.Manager) *App {
	t := &App{
		state:     boot,
		Animation: animation.NewLoader(graphics.White, graphics.Black, time.Now(), time.Now().Add(hack.ScaleDuration(time.Second/2))),
	}
	ui.AddHandler(input.A_Fall, t.handleInput)
	ui.AddHandler(input.B_Fall, t.handleInput)
	ui.AddHandler(input.C_Fall, t.handleInput)
	return t
}

// Update checks inputs and updates game state & animations based on the current time.
func (t *App) Update(now time.Time) {
	done := t.Animation.Update(now)
	if done {
		t.handleEvent(animationDone)
	}
}

// Frame returns the current animation's frame.
func (t *App) Frame() []color.RGBA { return t.Animation.Frame() }

func (t *App) Reset() {
	t.toIdle(-time.Second) // negative value to avoid 1 frame of transition
}

// handleInput converts input events to game events.
func (t *App) handleInput(e input.Event) {
	switch e {
	case input.A_Fall:
		t.handleEvent(cancel)
	case input.B_Fall:
		t.handleEvent(timer10M)
	case input.C_Fall:
		t.handleEvent(timer2M)
	}
}

// handleEvent signals the game that an event has occurred. These are inputs into the
// state machine.
func (t *App) handleEvent(e event) {
	switch e {
	case animationDone:
		t.animationDone()
	case timer2M:
		t.startTimer(hack.ScaleDuration(2 * time.Minute))
	case timer10M:
		// t.startTimer(hack.ScaleDuration(10 * time.Second))
		t.startTimer(hack.ScaleDuration(10 * time.Minute))
	case cancel:
		t.cancelTimer()
	}
}

func (t *App) animationDone() {
	switch t.state {
	case boot:
		t.toIdle(hack.ScaleDuration(2 * time.Second))
	case countdown:
		t.state = timerPop
		t.Animation = animation.NewFlasher(graphics.Red, time.Now().Add(hack.ScaleDuration(2*time.Second)))
	case timerPop:
		t.toIdle(hack.ScaleDuration(1 * time.Second))
	}
}

func (t *App) startTimer(d time.Duration) {
	switch t.state {
	// Allow boot and timer pop animations to be interrupted.
	case boot:
		fallthrough
	case timerPop:
		fallthrough
	case idle:
		t.state = countdown
		// Don't scale this duration because it has been scaled in the caller.
		t.Animation = animation.NewLoader(graphics.Black, graphics.CSIOrange, time.Now(), time.Now().Add(d))
	}
}

func (t *App) cancelTimer() {
	switch t.state {
	case countdown:
		t.toIdle((1 * time.Second))
	}
}

// toIdle fades from current animation to idle animation
func (t *App) toIdle(d time.Duration) {
	t.state = idle
	spin := animation.NewSpinner(graphics.K8SBlue)
	// Don't scale this duration because it has been scaled in the caller.
	t.Animation = animation.NewFader(t.Animation, spin, time.Now(), time.Now().Add(d))
}
