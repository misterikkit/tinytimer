package timer

import (
	"image/color"
	"time"

	"github.com/misterikkit/tinytimer/animation"
	"github.com/misterikkit/tinytimer/graphics"
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
		Animation: animation.NewLoader(graphics.White, graphics.Black, graphics.Black, time.Now(), time.Now().Add(time.Second/2)),
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
		t.handleEvent(timer2M)
	case input.C_Fall:
		t.handleEvent(timer10M)
	}
}

// handleEvent signals the game that an event has occurred. These are inputs into the
// state machine.
func (t *App) handleEvent(e event) {
	switch e {
	case animationDone:
		t.animationDone()
	case timer2M:
		t.startTimer(2 * time.Minute)
	case timer10M:
		// t.startTimer(10 * time.Second)
		t.startTimer(10 * time.Minute)
	case cancel:
		t.cancelTimer()
	}
}

func (t *App) animationDone() {
	switch t.state {
	case boot:
		t.toIdle(2 * time.Second)
	case countdown:
		t.state = timerPop
		t.Animation = animation.NewFlasher(graphics.Red, time.Now().Add(2*time.Second))
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
		t.Animation = animation.NewLoader(graphics.Black, graphics.Green, graphics.Red, time.Now(), time.Now().Add(d))
	}
}

func (t *App) cancelTimer() {
	switch t.state {
	case timerPop:
		fallthrough
	case countdown:
		t.toIdle(1 * time.Second)
	}
}

// toIdle fades from current animation to idle animation
func (t *App) toIdle(d time.Duration) {
	t.state = idle
	spin := animation.NewSpinner(color.RGBA{0, 0x7F, 0xFF, 0})
	t.Animation = animation.NewFader(t.Animation, spin, time.Now(), time.Now().Add(d))
}
